package server

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

type slowHoverHandler struct {
	delay time.Duration
}

func (h *slowHoverHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		ServerInfo: &lsp.ServerInfo{Name: "timeout-test"},
	}, nil
}

func (h *slowHoverHandler) Shutdown(_ context.Context) error { return nil }

func (h *slowHoverHandler) Hover(ctx context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	select {
	case <-time.After(h.delay):
		return &lsp.Hover{
			Contents: lsp.MarkupContent{Kind: lsp.PlainText, Value: "done"},
		}, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func TestRequestTimeout(t *testing.T) {
	clientReader, serverWriter := io.Pipe()
	serverReader, clientWriter := io.Pipe()

	h := &slowHoverHandler{delay: 2 * time.Second}
	s := NewServer(h, WithRequestTimeout(50*time.Millisecond))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = s.Run(ctx, pipeRWC{Reader: serverReader, Writer: serverWriter})
	}()

	clientConn := jsonrpc.NewConn(pipeRWC{Reader: clientReader, Writer: clientWriter}, jsonrpc.NewDispatcher())

	// Initialize first.
	req, _ := jsonrpc.NewRequest(jsonrpc.IntID(1), "initialize", lsp.InitializeParams{})
	_ = clientConn.WriteMessage(req)
	_, _ = clientConn.ReadMessage()

	// Send hover request — should timeout.
	hoverReq, _ := jsonrpc.NewRequest(jsonrpc.IntID(2), "textDocument/hover", lsp.HoverParams{})
	_ = clientConn.WriteMessage(hoverReq)

	msg, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	resp, ok := msg.(*jsonrpc.Response)
	if !ok {
		t.Fatalf("expected Response, got %T", msg)
	}
	if resp.Error == nil {
		t.Fatal("expected error response due to timeout")
	}
	if resp.Error.Code != jsonrpc.CodeRequestCancelled {
		t.Errorf("expected code %d, got %d", jsonrpc.CodeRequestCancelled, resp.Error.Code)
	}
}

func TestNoTimeoutByDefault(t *testing.T) {
	clientReader, serverWriter := io.Pipe()
	serverReader, clientWriter := io.Pipe()

	h := &slowHoverHandler{delay: 50 * time.Millisecond}
	s := NewServer(h) // no timeout

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		_ = s.Run(ctx, pipeRWC{Reader: serverReader, Writer: serverWriter})
	}()

	clientConn := jsonrpc.NewConn(pipeRWC{Reader: clientReader, Writer: clientWriter}, jsonrpc.NewDispatcher())

	// Initialize first.
	req, _ := jsonrpc.NewRequest(jsonrpc.IntID(1), "initialize", lsp.InitializeParams{})
	_ = clientConn.WriteMessage(req)
	_, _ = clientConn.ReadMessage()

	// Send hover — should succeed.
	hoverReq, _ := jsonrpc.NewRequest(jsonrpc.IntID(2), "textDocument/hover", lsp.HoverParams{})
	_ = clientConn.WriteMessage(hoverReq)

	msg, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	resp, ok := msg.(*jsonrpc.Response)
	if !ok {
		t.Fatalf("expected Response, got %T", msg)
	}
	if resp.Error != nil {
		t.Fatalf("expected success, got error: %s", resp.Error.Message)
	}
}
