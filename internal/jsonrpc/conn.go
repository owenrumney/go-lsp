package jsonrpc

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ErrExit is returned by a notification handler to make Serve return.
var ErrExit = errors.New("jsonrpc: exit notification received")

// maxContentLength caps a message body so a bad header can't cause a huge allocation.
const maxContentLength = 128 << 20

// Conn is a JSON-RPC 2.0 connection over a Content-Length framed stream.
type Conn struct {
	reader                *bufio.Reader
	writer                io.Writer
	closer                io.Closer
	writeMu               sync.Mutex
	dispatcher            *Dispatcher
	cancelMu              sync.Mutex
	cancels               map[string]context.CancelFunc
	nextID                atomic.Int64
	pendingMu             sync.Mutex
	pending               map[string]chan *Response
	requestTimeout        time.Duration
	callTimeout           time.Duration
	maxConcurrentRequests int
	requestSem            chan struct{}
	handlers              sync.WaitGroup
	queueMu               sync.Mutex
	queue                 []any
	queueReady            chan struct{}
	exitOnce              sync.Once
	exitCh                chan struct{}
}

func NewConn(rw io.ReadWriteCloser, dispatcher *Dispatcher) *Conn {
	return &Conn{
		reader:     bufio.NewReader(rw),
		writer:     rw,
		closer:     rw,
		dispatcher: dispatcher,
		cancels:    make(map[string]context.CancelFunc),
		pending:    make(map[string]chan *Response),
		queueReady: make(chan struct{}, 1),
		exitCh:     make(chan struct{}),
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
	contentLen := -1
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
			if err != nil || contentLen <= 0 {
				return 0, fmt.Errorf("jsonrpc: invalid Content-Length: %s", val)
			}
		}
	}
	switch {
	case contentLen < 0:
		return 0, fmt.Errorf("jsonrpc: missing Content-Length header")
	case contentLen > maxContentLength:
		return 0, fmt.Errorf("jsonrpc: Content-Length %d exceeds maximum %d", contentLen, maxContentLength)
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

// Serve reads messages in a loop and dispatches them in receipt order:
// notifications run serially on a worker goroutine, requests run concurrently
// but never start before an earlier notification completes. A blocking
// notification handler therefore delays everything behind it — long-running
// work belongs in a goroutine. $/cancelRequest, exit, and response routing are
// handled inline so cancellation, termination, and server-to-client calls
// always make progress.
//
// Serve returns ErrExit on the LSP exit notification and waits for in-flight
// request handlers before returning; the transport is closed on return. On a
// transport whose Close cannot interrupt a blocked Read (e.g. stdin), the
// reader goroutine remains parked until process exit.
func (c *Conn) Serve(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer c.handlers.Wait()
	defer cancel()

	go func() {
		select {
		case <-ctx.Done():
		case <-c.exitCh:
		}
		if c.closer != nil {
			_ = c.closer.Close()
		}
	}()

	c.handlers.Add(1)
	go func() {
		defer c.handlers.Done()
		c.dispatchWorker(ctx)
	}()

	type readResult struct {
		msg any
		err error
	}
	readCh := make(chan readResult)
	go func() {
		for {
			msg, err := c.ReadMessage()
			select {
			case readCh <- readResult{msg, err}:
			case <-ctx.Done():
				return
			}
			if err != nil && !errors.Is(err, ErrParse) && !errors.Is(err, ErrInvalidMessage) {
				return
			}
		}
	}()

	for {
		if err := c.serveErr(ctx); err != nil {
			return err
		}

		var r readResult
		select {
		case <-ctx.Done():
			return c.serveErr(ctx)
		case <-c.exitCh:
			return ErrExit
		case r = <-readCh:
		}
		if r.err != nil {
			if serveErr := c.serveErr(ctx); serveErr != nil {
				return serveErr
			}
			// Decode errors leave the framing intact, so keep serving.
			if errors.Is(r.err, ErrParse) {
				_ = c.WriteMessage(NewErrorResponse(ID{}, NewError(CodeParseError, r.err.Error())))
				continue
			}
			if errors.Is(r.err, ErrInvalidMessage) {
				_ = c.WriteMessage(NewErrorResponse(ID{}, NewError(CodeInvalidRequest, r.err.Error())))
				continue
			}
			return r.err
		}

		switch m := r.msg.(type) {
		case *Request:
			reqCtx, reqCancel := context.WithCancel(ctx) //nolint:gosec // released via finishRequest
			c.cancelMu.Lock()
			c.cancels[m.ID.String()] = reqCancel
			c.cancelMu.Unlock()
			c.enqueue(queuedRequest{req: m, ctx: reqCtx})
		case *Notification:
			switch m.Method {
			case "$/cancelRequest":
				c.handleCancel(m)
			case "exit":
				// Inline so the server can always terminate, even with the
				// dispatch queue backed up.
				c.handleNotification(ctx, m)
			default:
				c.enqueue(m)
			}
		case *Response:
			c.routeResponse(m)
		}
	}
}

func (c *Conn) serveErr(ctx context.Context) error {
	select {
	case <-c.exitCh:
		return ErrExit
	default:
	}
	return ctx.Err()
}

func (c *Conn) enqueue(msg any) {
	c.queueMu.Lock()
	c.queue = append(c.queue, msg)
	c.queueMu.Unlock()
	select {
	case c.queueReady <- struct{}{}:
	default:
	}
}

// queuedRequest pairs a request with the receipt-scoped context that
// $/cancelRequest cancels.
type queuedRequest struct {
	req *Request
	ctx context.Context
}

func (c *Conn) dispatchWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-c.queueReady:
		}
		for {
			c.queueMu.Lock()
			if len(c.queue) == 0 {
				c.queueMu.Unlock()
				break
			}
			msg := c.queue[0]
			c.queue = c.queue[1:]
			c.queueMu.Unlock()

			switch m := msg.(type) {
			case queuedRequest:
				c.handlers.Add(1)
				go func() {
					defer c.handlers.Done()
					c.handleRequestWithLimit(m.ctx, m.req)
				}()
			case *Notification:
				c.handleNotification(ctx, m)
			}
		}
	}
}

func (c *Conn) signalExit() {
	c.exitOnce.Do(func() { close(c.exitCh) })
}

// SetRequestTimeout sets a default timeout for all incoming requests.
// A zero duration means no timeout (the default).
func (c *Conn) SetRequestTimeout(d time.Duration) {
	c.requestTimeout = d
}

// SetCallTimeout bounds how long Call waits for the peer to respond.
// A zero duration means no timeout (the default).
func (c *Conn) SetCallTimeout(d time.Duration) {
	c.callTimeout = d
}

// SetMaxConcurrentRequests limits how many incoming requests may run at once.
// A value <= 0 keeps the default unlimited behavior.
func (c *Conn) SetMaxConcurrentRequests(n int) {
	c.maxConcurrentRequests = n
	if n <= 0 {
		c.requestSem = nil
		return
	}
	c.requestSem = make(chan struct{}, n)
}

func (c *Conn) handleRequestWithLimit(ctx context.Context, req *Request) {
	if c.requestSem != nil {
		select {
		case c.requestSem <- struct{}{}:
			defer func() { <-c.requestSem }()
		case <-ctx.Done():
			c.finishRequest(req.ID.String())
			_ = c.WriteMessage(NewErrorResponse(req.ID, NewError(CodeRequestCancelled, "request cancelled while queued")))
			return
		}
	}
	c.handleRequest(ctx, req)
}

func (c *Conn) handleRequest(ctx context.Context, req *Request) {
	idStr := req.ID.String()
	if c.requestTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.requestTimeout)
		defer cancel()
	}

	defer func() {
		if r := recover(); r != nil {
			resp := NewErrorResponse(req.ID, NewError(CodeInternalError, fmt.Sprintf("panic in handler %s: %v", req.Method, r)))
			_ = c.WriteMessage(resp)
		}
		c.finishRequest(idStr)
	}()

	if ctx.Err() != nil {
		_ = c.WriteMessage(NewErrorResponse(req.ID, NewError(CodeRequestCancelled, ctx.Err().Error())))
		return
	}
	resp := c.dispatcher.HandleRequest(ctx, req)
	_ = c.WriteMessage(resp)
}

// finishRequest releases a request's receipt-scoped cancel, if registered.
func (c *Conn) finishRequest(id string) {
	c.cancelMu.Lock()
	cancel := c.cancels[id]
	delete(c.cancels, id)
	c.cancelMu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (c *Conn) handleNotification(ctx context.Context, notif *Notification) {
	defer func() {
		if r := recover(); r != nil {
			slog.Error("panic in notification handler", "method", notif.Method, "panic", r)
		}
	}()

	if notif.Method == "$/cancelRequest" {
		c.handleCancel(notif)
		return
	}
	if err := c.dispatcher.HandleNotification(ctx, notif); errors.Is(err, ErrExit) {
		c.signalExit()
	}
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

	if c.callTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.callTimeout)
		defer cancel()
	}

	if err := c.WriteMessage(req); err != nil {
		return nil, err
	}

	select {
	case resp := <-ch:
		return resp, nil
	case <-ctx.Done():
		return nil, fmt.Errorf("%s: %w", method, ctx.Err())
	}
}

func (c *Conn) routeResponse(resp *Response) {
	idStr := resp.ID.String()
	c.pendingMu.Lock()
	ch, ok := c.pending[idStr]
	c.pendingMu.Unlock()
	if ok {
		// Non-blocking: a duplicate response must not wedge the read loop.
		select {
		case ch <- resp:
		default:
		}
	}
}

// Notify sends a notification to the peer.
func (c *Conn) Notify(_ context.Context, method string, params any) error {
	notif, err := NewNotification(method, params)
	if err != nil {
		return err
	}
	return c.WriteMessage(notif)
}
