package jsonrpc

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

// Conn is a JSON-RPC 2.0 connection over a Content-Length framed stream.
type Conn struct {
	reader     *bufio.Reader
	writer     io.Writer
	writeMu    sync.Mutex
	dispatcher *Dispatcher
	cancelMu   sync.Mutex
	cancels    map[string]context.CancelFunc
	nextID     atomic.Int64
	pendingMu  sync.Mutex
	pending    map[string]chan *Response
}

func NewConn(rw io.ReadWriteCloser, dispatcher *Dispatcher) *Conn {
	return &Conn{
		reader:     bufio.NewReader(rw),
		writer:     rw,
		dispatcher: dispatcher,
		cancels:    make(map[string]context.CancelFunc),
		pending:    make(map[string]chan *Response),
	}
}

// ReadMessage reads and decodes a single Content-Length framed JSON-RPC message.
func (c *Conn) ReadMessage() (any, error) {
	contentLen, err := c.readHeaders()
	if err != nil {
		return nil, err
	}

	body := make([]byte, contentLen)
	if _, err := io.ReadFull(c.reader, body); err != nil {
		return nil, fmt.Errorf("jsonrpc: failed to read body: %w", err)
	}

	return DecodeMessage(body)
}

func (c *Conn) readHeaders() (int, error) {
	var contentLen int
	for {
		line, err := c.reader.ReadString('\n')
		if err != nil {
			return 0, fmt.Errorf("jsonrpc: failed to read header: %w", err)
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		if val, ok := strings.CutPrefix(line, "Content-Length:"); ok {
			val = strings.TrimSpace(val)
			contentLen, err = strconv.Atoi(val)
			if err != nil {
				return 0, fmt.Errorf("jsonrpc: invalid Content-Length: %s", val)
			}
		}
	}
	if contentLen == 0 {
		return 0, fmt.Errorf("jsonrpc: missing Content-Length header")
	}
	return contentLen, nil
}

// WriteMessage encodes and writes a JSON-RPC message with Content-Length framing.
func (c *Conn) WriteMessage(msg any) error {
	data, err := json.Marshal(msg)
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

// Serve reads messages in a loop and dispatches them.
func (c *Conn) Serve(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		msg, err := c.ReadMessage()
		if err != nil {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			return err
		}

		switch m := msg.(type) {
		case *Request:
			go c.handleRequest(ctx, m)
		case *Notification:
			c.handleNotification(ctx, m)
		case *Response:
			c.routeResponse(m)
		}
	}
}

func (c *Conn) handleRequest(ctx context.Context, req *Request) {
	reqCtx, cancel := context.WithCancel(ctx)
	idStr := req.ID.String()

	c.cancelMu.Lock()
	c.cancels[idStr] = cancel
	c.cancelMu.Unlock()

	defer func() {
		cancel()
		c.cancelMu.Lock()
		delete(c.cancels, idStr)
		c.cancelMu.Unlock()
	}()

	resp := c.dispatcher.HandleRequest(reqCtx, req)
	_ = c.WriteMessage(resp)
}

func (c *Conn) handleNotification(ctx context.Context, notif *Notification) {
	if notif.Method == "$/cancelRequest" {
		c.handleCancel(notif)
		return
	}
	c.dispatcher.HandleNotification(ctx, notif)
}

func (c *Conn) handleCancel(notif *Notification) {
	var params struct {
		ID ID `json:"id"`
	}
	if err := json.Unmarshal(notif.Params, &params); err != nil {
		return
	}

	c.cancelMu.Lock()
	cancel, ok := c.cancels[params.ID.String()]
	c.cancelMu.Unlock()

	if ok {
		cancel()
	}
}

// Call sends a request to the peer and waits for a response.
func (c *Conn) Call(ctx context.Context, method string, params any) (*Response, error) {
	id := IntID(c.nextID.Add(1))
	req, err := NewRequest(id, method, params)
	if err != nil {
		return nil, err
	}

	ch := make(chan *Response, 1)
	idStr := id.String()

	c.pendingMu.Lock()
	c.pending[idStr] = ch
	c.pendingMu.Unlock()

	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, idStr)
		c.pendingMu.Unlock()
	}()

	if err := c.WriteMessage(req); err != nil {
		return nil, err
	}

	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (c *Conn) routeResponse(resp *Response) {
	idStr := resp.ID.String()
	c.pendingMu.Lock()
	ch, ok := c.pending[idStr]
	c.pendingMu.Unlock()
	if ok {
		ch <- resp
	}
}

// Notify sends a notification to the peer.
func (c *Conn) Notify(ctx context.Context, method string, params any) error {
	notif, err := NewNotification(method, params)
	if err != nil {
		return err
	}
	return c.WriteMessage(notif)
}
