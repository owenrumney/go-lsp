package server

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

// lifecycleClient drives a server over a net.Pipe with raw JSON-RPC messages.
type lifecycleClient struct {
	t       *testing.T
	conn    *jsonrpc.Conn
	end     net.Conn
	runDone <-chan error
}

func startLifecycleServer(t *testing.T, h LifecycleHandler, opts ...Option) *lifecycleClient {
	t.Helper()
	clientEnd, serverEnd := net.Pipe()

	s := NewServer(h, opts...)
	runDone := make(chan error, 1)
	go func() {
		runDone <- s.Run(t.Context(), serverEnd)
	}()
	t.Cleanup(func() { _ = clientEnd.Close() })

	return &lifecycleClient{
		t:       t,
		conn:    jsonrpc.NewConn(clientEnd, jsonrpc.NewDispatcher()),
		end:     clientEnd,
		runDone: runDone,
	}
}

func (c *lifecycleClient) call(id int64, method string, params any) *jsonrpc.Response {
	c.t.Helper()
	req, err := jsonrpc.NewRequest(jsonrpc.IntID(id), method, params)
	if err != nil {
		c.t.Fatal(err)
	}
	if err := c.conn.WriteMessage(req); err != nil {
		c.t.Fatal(err)
	}
	msg, err := c.conn.ReadMessage()
	if err != nil {
		c.t.Fatal(err)
	}
	resp, ok := msg.(*jsonrpc.Response)
	if !ok {
		c.t.Fatalf("expected *Response, got %T", msg)
	}
	return resp
}

func (c *lifecycleClient) notify(method string, params any) {
	c.t.Helper()
	notif, err := jsonrpc.NewNotification(method, params)
	if err != nil {
		c.t.Fatal(err)
	}
	if err := c.conn.WriteMessage(notif); err != nil {
		c.t.Fatal(err)
	}
}

func (c *lifecycleClient) waitRun() error {
	c.t.Helper()
	select {
	case err := <-c.runDone:
		return err
	case <-time.After(5 * time.Second):
		c.t.Fatal("Run did not return")
		return nil
	}
}

func TestLifecycle_RequestBeforeInitialize(t *testing.T) {
	c := startLifecycleServer(t, &mockHandler{})

	resp := c.call(1, "textDocument/hover", &lsp.HoverParams{})
	if resp.Error == nil || resp.Error.Code != jsonrpc.CodeServerNotInitialized {
		t.Fatalf("expected ServerNotInitialized, got %+v", resp)
	}
}

func TestLifecycle_DoubleInitialize(t *testing.T) {
	c := startLifecycleServer(t, &mockHandler{})

	if resp := c.call(1, "initialize", &lsp.InitializeParams{}); resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}
	resp := c.call(2, "initialize", &lsp.InitializeParams{})
	if resp.Error == nil || resp.Error.Code != jsonrpc.CodeInvalidRequest {
		t.Fatalf("expected InvalidRequest for second initialize, got %+v", resp)
	}
}

func TestLifecycle_NotificationBeforeInitializeDropped(t *testing.T) {
	h := &mockHandler{}
	c := startLifecycleServer(t, h)

	c.notify("textDocument/didOpen", &lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{URI: "file:///early.txt"},
	})
	// The initialize round-trip is the sync barrier for the notification.
	if resp := c.call(1, "initialize", &lsp.InitializeParams{}); resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}
	if len(h.opened) != 0 {
		t.Fatalf("expected pre-initialize didOpen to be dropped, got %v", h.opened)
	}
}

func TestLifecycle_RequestAfterShutdown(t *testing.T) {
	c := startLifecycleServer(t, &mockHandler{})

	if resp := c.call(1, "initialize", &lsp.InitializeParams{}); resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}
	if resp := c.call(2, "shutdown", nil); resp.Error != nil {
		t.Fatalf("shutdown failed: %v", resp.Error)
	}
	resp := c.call(3, "textDocument/hover", &lsp.HoverParams{})
	if resp.Error == nil || resp.Error.Code != jsonrpc.CodeInvalidRequest {
		t.Fatalf("expected InvalidRequest after shutdown, got %+v", resp)
	}
}

func TestLifecycle_ExitAfterShutdownReturnsNil(t *testing.T) {
	c := startLifecycleServer(t, &mockHandler{})

	if resp := c.call(1, "initialize", &lsp.InitializeParams{}); resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}
	if resp := c.call(2, "shutdown", nil); resp.Error != nil {
		t.Fatalf("shutdown failed: %v", resp.Error)
	}
	c.notify("exit", nil)

	if err := c.waitRun(); err != nil {
		t.Fatalf("Run returned %v, want nil after shutdown+exit", err)
	}
}

func TestLifecycle_ExitWithoutShutdown(t *testing.T) {
	c := startLifecycleServer(t, &mockHandler{})

	if resp := c.call(1, "initialize", &lsp.InitializeParams{}); resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}
	c.notify("exit", nil)

	if err := c.waitRun(); !errors.Is(err, ErrExitWithoutShutdown) {
		t.Fatalf("Run returned %v, want ErrExitWithoutShutdown", err)
	}
}

func TestInitialize_NilServerInfoDoesNotPanic(t *testing.T) {
	c := startLifecycleServer(t, &bareInitHandler{})

	resp := c.call(1, "initialize", &lsp.InitializeParams{})
	if resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}
	var result lsp.InitializeResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatal(err)
	}
}

type bareInitHandler struct{}

func (h *bareInitHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{}, nil
}

func (h *bareInitHandler) Shutdown(_ context.Context) error { return nil }

type blockingInitHandler struct {
	started chan struct{}
	release chan struct{}
	inits   atomic.Int32
}

func (h *blockingInitHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	h.inits.Add(1)
	h.started <- struct{}{}
	<-h.release
	return &lsp.InitializeResult{}, nil
}

func (h *blockingInitHandler) Shutdown(_ context.Context) error { return nil }

func TestLifecycle_ConcurrentInitializeRejected(t *testing.T) {
	h := &blockingInitHandler{started: make(chan struct{}, 1), release: make(chan struct{})}
	c := startLifecycleServer(t, h)

	req1, err := jsonrpc.NewRequest(jsonrpc.IntID(1), "initialize", &lsp.InitializeParams{})
	if err != nil {
		t.Fatal(err)
	}
	if err := c.conn.WriteMessage(req1); err != nil {
		t.Fatal(err)
	}
	<-h.started

	// Second initialize while the first is still in the handler.
	resp := c.call(2, "initialize", &lsp.InitializeParams{})
	if resp.Error == nil || resp.Error.Code != jsonrpc.CodeInvalidRequest {
		t.Fatalf("expected InvalidRequest for concurrent initialize, got %+v", resp)
	}

	close(h.release)
	msg, err := c.conn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	if resp, ok := msg.(*jsonrpc.Response); !ok || resp.Error != nil {
		t.Fatalf("first initialize failed: %+v", msg)
	}
	if got := h.inits.Load(); got != 1 {
		t.Fatalf("Initialize ran %d times, want 1", got)
	}
}

func TestLifecycle_UnknownMethodBeforeInitialize(t *testing.T) {
	c := startLifecycleServer(t, &mockHandler{})

	resp := c.call(1, "does/notExist", nil)
	if resp.Error == nil || resp.Error.Code != jsonrpc.CodeServerNotInitialized {
		t.Fatalf("expected ServerNotInitialized pre-init, got %+v", resp)
	}

	if resp := c.call(2, "initialize", &lsp.InitializeParams{}); resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}
	resp = c.call(3, "does/notExist", nil)
	if resp.Error == nil || resp.Error.Code != jsonrpc.CodeMethodNotFound {
		t.Fatalf("expected MethodNotFound post-init, got %+v", resp)
	}
}
