package servertest

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// rpcConn is a minimal client-side JSON-RPC 2.0 connection over Content-Length framed streams.
type rpcConn struct {
	reader  *bufio.Reader
	writer  io.Writer
	writeMu sync.Mutex

	nextID  atomic.Int64
	pending sync.Map // map[string]chan *rpcResponse

	// notifHandler is called for server-to-client notifications.
	notifHandler func(method string, params json.RawMessage)
	// requestHandler is called for server-to-client requests (e.g. window/workDoneProgress/create).
	requestHandler func(method string, params json.RawMessage) (any, error)

	done chan struct{}
}

type rpcRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id,omitempty"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      any             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *rpcError       `json:"error,omitempty"`
}

type rpcNotification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *rpcError) Error() string {
	return fmt.Sprintf("jsonrpc error %d: %s", e.Code, e.Message)
}

// rawMsg is used for initial JSON parsing to determine message type.
type rawMsg struct {
	ID     *json.RawMessage `json:"id,omitempty"`
	Method *string          `json:"method,omitempty"`
	Params json.RawMessage  `json:"params,omitempty"`
	Result json.RawMessage  `json:"result,omitempty"`
	Error  *rpcError        `json:"error,omitempty"`
}

func newRPCConn(rw io.ReadWriter) *rpcConn {
	return &rpcConn{
		reader: bufio.NewReader(rw),
		writer: rw,
		done:   make(chan struct{}),
	}
}

// readLoop reads messages from the connection and routes them.
// It closes c.done when it returns.
func (c *rpcConn) readLoop() {
	defer close(c.done)
	for {
		data, err := c.readFrame()
		if err != nil {
			return
		}

		var raw rawMsg
		if err := json.Unmarshal(data, &raw); err != nil {
			continue
		}

		hasID := raw.ID != nil && string(*raw.ID) != "null"

		switch {
		case hasID && raw.Method == nil:
			// Response
			c.routeResponse(data, raw)
		case hasID && raw.Method != nil:
			// Server-to-client request
			c.handleIncomingRequest(data, raw)
		case raw.Method != nil:
			// Notification from server
			if c.notifHandler != nil {
				c.notifHandler(*raw.Method, raw.Params)
			}
		}
	}
}

func (c *rpcConn) routeResponse(_ []byte, raw rawMsg) {
	idStr := idToString(raw.ID)
	if v, ok := c.pending.Load(idStr); ok {
		ch := v.(chan *rpcResponse)
		resp := &rpcResponse{
			Result: raw.Result,
			Error:  raw.Error,
		}
		ch <- resp
	}
}

func (c *rpcConn) handleIncomingRequest(_ []byte, raw rawMsg) {
	// Parse the ID
	var id any
	if raw.ID != nil {
		_ = json.Unmarshal(*raw.ID, &id)
	}

	var result any
	var rpcErr *rpcError
	if c.requestHandler != nil {
		res, err := c.requestHandler(*raw.Method, raw.Params)
		if err != nil {
			rpcErr = &rpcError{Code: -32603, Message: err.Error()}
		} else {
			result = res
		}
	}

	// Send response
	var resultRaw json.RawMessage
	if rpcErr == nil {
		if result != nil {
			b, _ := json.Marshal(result)
			resultRaw = b
		} else {
			resultRaw = json.RawMessage("null")
		}
	}

	resp := struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      any             `json:"id"`
		Result  json.RawMessage `json:"result,omitempty"`
		Error   *rpcError       `json:"error,omitempty"`
	}{
		JSONRPC: "2.0",
		ID:      id,
		Result:  resultRaw,
		Error:   rpcErr,
	}
	_ = c.writeJSON(resp)
}

func (c *rpcConn) readFrame() ([]byte, error) {
	var contentLen int
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if val, ok := strings.CutPrefix(line, "Content-Length:"); ok {
			val = strings.TrimSpace(val)
			contentLen, err = strconv.Atoi(val)
			if err != nil {
				return nil, fmt.Errorf("invalid Content-Length: %s", val)
			}
		}
	}
	if contentLen == 0 {
		return nil, fmt.Errorf("missing Content-Length header")
	}

	body := make([]byte, contentLen)
	if _, err := io.ReadFull(c.reader, body); err != nil {
		return nil, err
	}
	return body, nil
}

func (c *rpcConn) writeJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	header := fmt.Sprintf("Content-Length: %d\r\n\r\n", len(data))
	if _, err := io.WriteString(c.writer, header); err != nil {
		return err
	}
	_, err = c.writer.Write(data)
	return err
}

// call sends a JSON-RPC request and waits for the response.
func (c *rpcConn) call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	id := c.nextID.Add(1)
	idStr := fmt.Sprintf("%d", id)

	var raw json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		raw = b
	}

	ch := make(chan *rpcResponse, 1)
	c.pending.Store(idStr, ch)
	defer c.pending.Delete(idStr)

	req := rpcRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  raw,
	}
	if err := c.writeJSON(req); err != nil {
		return nil, err
	}

	select {
	case resp := <-ch:
		if resp.Error != nil {
			return nil, resp.Error
		}
		return resp.Result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.done:
		return nil, fmt.Errorf("connection closed")
	}
}

// notify sends a JSON-RPC notification (no ID, no response expected).
func (c *rpcConn) notify(_ context.Context, method string, params any) error {
	var raw json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return err
		}
		raw = b
	}

	notif := rpcNotification{
		JSONRPC: "2.0",
		Method:  method,
		Params:  raw,
	}
	return c.writeJSON(notif)
}

func idToString(raw *json.RawMessage) string {
	if raw == nil {
		return ""
	}
	var n int64
	if err := json.Unmarshal(*raw, &n); err == nil {
		return strconv.FormatInt(n, 10)
	}
	var s string
	if err := json.Unmarshal(*raw, &s); err == nil {
		return s
	}
	return string(*raw)
}
