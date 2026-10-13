package client_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/client"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

// inProcess runs a server.Server on the other end of an in-memory pipe.
func inProcess(handler server.LifecycleHandler, opts ...server.Option) client.Transport {
	return func(ctx context.Context) (io.ReadWriteCloser, error) {
		clientConn, serverConn := net.Pipe()
		srv := server.NewServer(handler, opts...)
		go func() { _ = srv.Run(ctx, serverConn) }()
		return clientConn, nil
	}
}

type baseHandler struct {
	client *server.Client
}

func (h *baseHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{HoverProvider: &[]bool{true}[0]},
		ServerInfo:   &lsp.ServerInfo{Name: "test-server", Version: "1.2.3"},
	}, nil
}

func (h *baseHandler) Shutdown(_ context.Context) error { return nil }

func (h *baseHandler) SetClient(c *server.Client) { h.client = c }

type hoverHandler struct {
	baseHandler
	gotCapabilities lsp.ClientCapabilities
}

func (h *hoverHandler) Initialize(_ context.Context, params *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	h.gotCapabilities = params.Capabilities
	return h.baseHandler.Initialize(context.Background(), params)
}

func (h *hoverHandler) Hover(_ context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
	return &lsp.Hover{
		Contents: lsp.NewHoverContents(lsp.Markdown, "hovered "+string(params.TextDocument.URI)),
	}, nil
}

func dial(t *testing.T, handler server.LifecycleHandler, cfg client.Config) *client.Client {
	t.Helper()

	cfg.Transport = inProcess(handler)
	c, err := client.Dial(t.Context(), cfg)
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}
	t.Cleanup(func() { _ = c.Close() })
	return c
}

func TestDialInitializesAndQueries(t *testing.T) {
	h := &hoverHandler{}
	c := dial(t, h, client.Config{RootURI: "file:///work"})

	if got := c.ServerInfo(); got == nil || got.Name != "test-server" || got.Version != "1.2.3" {
		t.Fatalf("ServerInfo() = %+v, want test-server 1.2.3", got)
	}
	if got := c.PositionEncoding(); got != lsp.PositionEncodingUTF16 {
		t.Fatalf("PositionEncoding() = %q, want utf-16", got)
	}
	if h.gotCapabilities.TextDocument == nil || h.gotCapabilities.TextDocument.Hover == nil {
		t.Fatal("server saw no textDocument.hover capability; DefaultCapabilities was not sent")
	}

	hover, err := c.Hover(t.Context(), "file:///a.go", lsp.Position{Line: 1, Character: 2})
	if err != nil {
		t.Fatalf("Hover: %v", err)
	}
	if hover == nil {
		t.Fatal("Hover() = nil, want a result")
	}
}

func TestSupportsUsesStaticCapabilities(t *testing.T) {
	c := dial(t, &hoverHandler{}, client.Config{})

	if !c.Supports("textDocument/hover", "file:///a.go", "go") {
		t.Fatal("hover should be supported; the server advertised it")
	}
	if c.Supports("textDocument/rename", "file:///a.go", "go") {
		t.Fatal("rename should not be supported; the server never advertised it")
	}
}

type registeringHandler struct {
	baseHandler
}

func (h *registeringHandler) Hover(ctx context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	opts, err := json.Marshal(lsp.TextDocumentRegistrationOptions{
		DocumentSelector: &lsp.DocumentSelector{{Language: "go"}},
	})
	if err != nil {
		return nil, err
	}

	err = h.client.RegisterCapability(ctx, &lsp.RegistrationParams{
		Registrations: []lsp.Registration{{
			ID:              "rename-1",
			Method:          "textDocument/rename",
			RegisterOptions: opts,
		}},
	})
	if err != nil {
		return nil, err
	}
	return &lsp.Hover{Contents: lsp.NewHoverContents(lsp.PlainText, "ok")}, nil
}

func TestDynamicRegistrationChangesSupports(t *testing.T) {
	c := dial(t, &registeringHandler{}, client.Config{})

	if c.Supports("textDocument/rename", "file:///a.go", "go") {
		t.Fatal("rename supported before the server registered it")
	}

	if _, err := c.Hover(t.Context(), "file:///a.go", lsp.Position{}); err != nil {
		t.Fatalf("Hover: %v", err)
	}

	if !c.Supports("textDocument/rename", "file:///a.go", "go") {
		t.Fatal("rename not supported after client/registerCapability")
	}
	if c.Supports("textDocument/rename", "file:///a.ts", "typescript") {
		t.Fatal("registration selector was ignored; typescript should not match a go filter")
	}
	if got := c.Registrations("textDocument/rename"); len(got) != 1 {
		t.Fatalf("Registrations() = %v, want one entry", got)
	}
}

type configHandler struct {
	baseHandler
	got []json.RawMessage
}

func (h *configHandler) Hover(ctx context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	const section = "mylang"
	values, err := h.client.Configuration(ctx, &lsp.ConfigurationParams{
		Items: []lsp.ConfigurationItem{{Section: section}},
	})
	if err != nil {
		return nil, err
	}
	h.got = values
	return &lsp.Hover{Contents: lsp.NewHoverContents(lsp.PlainText, "ok")}, nil
}

type configAnswering struct {
	client.NopHandler
}

func (configAnswering) Configuration(_ context.Context, items []lsp.ConfigurationItem) ([]json.RawMessage, error) {
	out := make([]json.RawMessage, len(items))
	for i := range out {
		out[i] = json.RawMessage(`{"tabSize":4}`)
	}
	return out, nil
}

func TestHandlerAnswersServerRequests(t *testing.T) {
	h := &configHandler{}
	c := dial(t, h, client.Config{Handler: configAnswering{}})

	if _, err := c.Hover(t.Context(), "file:///a.go", lsp.Position{}); err != nil {
		t.Fatalf("Hover: %v", err)
	}
	if len(h.got) != 1 || string(h.got[0]) != `{"tabSize":4}` {
		t.Fatalf("server received %v, want one {\"tabSize\":4}", h.got)
	}
}

type progressHandler struct {
	baseHandler
}

func (h *progressHandler) Hover(ctx context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	token := lsp.ProgressToken(json.RawMessage(`"indexing"`))
	for _, value := range []string{`{"kind":"begin","title":"Indexing"}`, `{"kind":"end"}`} {
		err := h.client.Progress(ctx, &lsp.ProgressParams{Token: token, Value: json.RawMessage(value)})
		if err != nil {
			return nil, err
		}
	}
	return &lsp.Hover{Contents: lsp.NewHoverContents(lsp.PlainText, "ok")}, nil
}

func TestReadyWaitsForProgressToEnd(t *testing.T) {
	// A long grace means only the progress path can make this ready.
	c := dial(t, &progressHandler{}, client.Config{ReadyGrace: time.Hour})

	if _, err := c.Hover(t.Context(), "file:///a.go", lsp.Position{}); err != nil {
		t.Fatalf("Hover: %v", err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	if err := c.Ready(ctx); err != nil {
		t.Fatalf("Ready: %v", err)
	}
}

func TestReadyFallsBackToGraceWhenNoProgress(t *testing.T) {
	c := dial(t, &hoverHandler{}, client.Config{ReadyGrace: 150 * time.Millisecond})

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	start := time.Now()
	if err := c.Ready(ctx); err != nil {
		t.Fatalf("Ready: %v", err)
	}
	if elapsed := time.Since(start); elapsed < 100*time.Millisecond {
		t.Fatalf("Ready returned after %s, want to wait out the grace period", elapsed)
	}
}

type stuckProgressHandler struct {
	baseHandler
}

func (h *stuckProgressHandler) Hover(ctx context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	err := h.client.Progress(ctx, &lsp.ProgressParams{
		Token: lsp.ProgressToken(json.RawMessage(`"stuck"`)),
		Value: json.RawMessage(`{"kind":"begin","title":"Indexing forever"}`),
	})
	if err != nil {
		return nil, err
	}
	return &lsp.Hover{Contents: lsp.NewHoverContents(lsp.PlainText, "ok")}, nil
}

func TestReadyGivesUpWhenServerDies(t *testing.T) {
	c := dial(t, &stuckProgressHandler{}, client.Config{ReadyGrace: time.Hour})

	// Indexing starts and never ends, so Ready can only return on failure.
	if _, err := c.Hover(t.Context(), "file:///a.go", lsp.Position{}); err != nil {
		t.Fatalf("Hover: %v", err)
	}

	go func() { _ = c.Close() }()

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	if err := c.Ready(ctx); err == nil {
		t.Fatal("Ready() = nil, want an error once the connection closed")
	} else if errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Ready blocked until the test deadline: %v", err)
	}
}

func TestCloseIsIdempotentAndClosesDone(t *testing.T) {
	h := &hoverHandler{}
	c, err := client.Dial(t.Context(), client.Config{Transport: inProcess(h)})
	if err != nil {
		t.Fatal(err)
	}

	if err := c.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if err := c.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}

	select {
	case <-c.Done():
	case <-time.After(5 * time.Second):
		t.Fatal("Done() never closed after Close")
	}
}

type disablingHandler struct {
	baseHandler
}

func (h *disablingHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	var caps lsp.ServerCapabilities
	body := `{"hoverProvider":true,"codeActionProvider":false,"renameProvider":true}`
	if err := json.Unmarshal([]byte(body), &caps); err != nil {
		return nil, err
	}
	return &lsp.InitializeResult{
		Capabilities: caps,
		ServerInfo:   &lsp.ServerInfo{Name: "test-server"},
	}, nil
}

func TestSupportsRejectsProvidersSentAsFalse(t *testing.T) {
	c := dial(t, &disablingHandler{}, client.Config{})

	if c.Supports("textDocument/codeAction", "file:///a.go", "go") {
		t.Error("codeActionProvider was false; a nil check alone would wrongly pass here")
	}
	if !c.Supports("textDocument/rename", "file:///a.go", "go") {
		t.Error("renameProvider was true and should be supported")
	}
	if !c.Supports("textDocument/hover", "file:///a.go", "go") {
		t.Error("hoverProvider was true and should be supported")
	}
}
