package server

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/owenrumney/go-lsp/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

// mockHandler implements LifecycleHandler, TextDocumentSyncHandler, and HoverHandler.
type mockHandler struct {
	initialized bool
	opened      []lsp.DocumentURI
}

func (m *mockHandler) Initialize(_ context.Context, params *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	m.initialized = true
	return &lsp.InitializeResult{
		ServerInfo: &lsp.ServerInfo{Name: "test-server", Version: "0.1.0"},
	}, nil
}

func (m *mockHandler) Shutdown(_ context.Context) error {
	return nil
}

func (m *mockHandler) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
	m.opened = append(m.opened, params.TextDocument.URI)
	return nil
}

func (m *mockHandler) DidChange(_ context.Context, _ *lsp.DidChangeTextDocumentParams) error {
	return nil
}

func (m *mockHandler) DidClose(_ context.Context, _ *lsp.DidCloseTextDocumentParams) error {
	return nil
}

func (m *mockHandler) Hover(_ context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
	return &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.Markdown,
			Value: "Hello from hover",
		},
	}, nil
}

func TestBuildCapabilities(t *testing.T) {
	h := &mockHandler{}
	caps := buildCapabilities(h)

	if caps.HoverProvider == nil || !*caps.HoverProvider {
		t.Error("expected HoverProvider to be true")
	}
	if caps.TextDocumentSync == nil {
		t.Fatal("expected TextDocumentSync to be set")
	}
	if caps.TextDocumentSync.Change != lsp.SyncIncremental {
		t.Errorf("expected SyncIncremental, got %d", caps.TextDocumentSync.Change)
	}
	if caps.CompletionProvider != nil {
		t.Error("expected CompletionProvider to be nil (not implemented)")
	}
}

type pipeRWC struct {
	io.Reader
	io.Writer
}

func (pipeRWC) Close() error { return nil }

func TestServerInitializeHandshake(t *testing.T) {
	// Create a pipe pair: client writes to serverIn, reads from serverOut.
	clientReader, serverWriter := io.Pipe()
	serverReader, clientWriter := io.Pipe()

	h := &mockHandler{}
	s := NewServer(h)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.Run(ctx, pipeRWC{Reader: serverReader, Writer: serverWriter})
	}()

	// Send initialize request from "client" side.
	clientConn := jsonrpc.NewConn(pipeRWC{Reader: clientReader, Writer: clientWriter}, jsonrpc.NewDispatcher())

	initParams := lsp.InitializeParams{
		Capabilities: lsp.ClientCapabilities{},
	}
	req, err := jsonrpc.NewRequest(jsonrpc.IntID(1), "initialize", initParams)
	if err != nil {
		t.Fatal(err)
	}
	if err := clientConn.WriteMessage(req); err != nil {
		t.Fatal(err)
	}

	// Read the response.
	msg, err := clientConn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	resp, ok := msg.(*jsonrpc.Response)
	if !ok {
		t.Fatalf("expected Response, got %T", msg)
	}
	if resp.Error != nil {
		t.Fatalf("initialize failed: %s", resp.Error.Message)
	}

	var result lsp.InitializeResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatal(err)
	}

	if result.ServerInfo == nil || result.ServerInfo.Name != "test-server" {
		t.Error("expected server info to be set")
	}
	if result.Capabilities.HoverProvider == nil || !*result.Capabilities.HoverProvider {
		t.Error("expected HoverProvider capability")
	}
	if result.Capabilities.TextDocumentSync == nil {
		t.Error("expected TextDocumentSync capability")
	}

	if !h.initialized {
		t.Error("handler was not initialized")
	}

	// Send textDocument/hover request.
	hoverParams := lsp.HoverParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: "file:///test.go"},
			Position:     lsp.Position{Line: 0, Character: 0},
		},
	}
	hoverReq, err := jsonrpc.NewRequest(jsonrpc.IntID(2), "textDocument/hover", hoverParams)
	if err != nil {
		t.Fatal(err)
	}
	if err := clientConn.WriteMessage(hoverReq); err != nil {
		t.Fatal(err)
	}

	msg, err = clientConn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	hoverResp, ok := msg.(*jsonrpc.Response)
	if !ok {
		t.Fatalf("expected Response, got %T", msg)
	}
	if hoverResp.Error != nil {
		t.Fatalf("hover failed: %s", hoverResp.Error.Message)
	}

	var hover lsp.Hover
	if err := json.Unmarshal(hoverResp.Result, &hover); err != nil {
		t.Fatal(err)
	}
	if hover.Contents.Value != "Hello from hover" {
		t.Errorf("hover contents = %q, want %q", hover.Contents.Value, "Hello from hover")
	}

	cancel()
}
