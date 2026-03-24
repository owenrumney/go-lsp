package servertest_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
	"github.com/owenrumney/go-lsp/servertest"
)

// testHandler implements enough handler interfaces for testing the harness.
type testHandler struct {
	client *server.Client
	docs   map[lsp.DocumentURI]string
}

func newTestHandler() *testHandler {
	return &testHandler{docs: make(map[lsp.DocumentURI]string)}
}

func (h *testHandler) SetClient(c *server.Client) { h.client = c }

func (h *testHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptions{
				OpenClose: boolPtr(true),
				Change:    lsp.SyncFull,
				Save:      &lsp.SaveOptions{IncludeText: boolPtr(true)},
			},
			HoverProvider:      boolPtr(true),
			CompletionProvider: &lsp.CompletionOptions{},
		},
		ServerInfo: &lsp.ServerInfo{Name: "test-server", Version: "0.1.0"},
	}, nil
}

func (h *testHandler) Shutdown(_ context.Context) error { return nil }

func (h *testHandler) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
	h.docs[params.TextDocument.URI] = params.TextDocument.Text
	return nil
}

func (h *testHandler) DidChange(_ context.Context, params *lsp.DidChangeTextDocumentParams) error {
	if len(params.ContentChanges) > 0 {
		h.docs[params.TextDocument.URI] = params.ContentChanges[len(params.ContentChanges)-1].Text
	}
	return nil
}

func (h *testHandler) DidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
	delete(h.docs, params.TextDocument.URI)
	return nil
}

func (h *testHandler) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error {
	if params.Text != nil {
		h.docs[params.TextDocument.URI] = *params.Text
	}
	text, ok := h.docs[params.TextDocument.URI]
	if !ok {
		return nil
	}

	var diags []lsp.Diagnostic
	sev := lsp.SeverityWarning
	for i, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "ERROR") {
			diags = append(diags, lsp.Diagnostic{
				Range: lsp.Range{
					Start: lsp.Position{Line: i, Character: 0},
					End:   lsp.Position{Line: i, Character: len(line)},
				},
				Severity: &sev,
				Source:   "test",
				Message:  fmt.Sprintf("error on line %d", i+1),
			})
		}
	}

	return h.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diags,
	})
}

func (h *testHandler) Hover(_ context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
	text, ok := h.docs[params.TextDocument.URI]
	if !ok {
		return nil, nil
	}
	lines := strings.Split(text, "\n")
	line := params.Position.Line
	if line < 0 || line >= len(lines) {
		return nil, nil
	}
	return &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.Markdown,
			Value: fmt.Sprintf("Line %d: `%s`", line, lines[line]),
		},
	}, nil
}

func (h *testHandler) Completion(_ context.Context, _ *lsp.CompletionParams) (*lsp.CompletionList, error) {
	kind := lsp.CompletionItemKindKeyword
	return &lsp.CompletionList{
		Items: []lsp.CompletionItem{
			{Label: "hello", Kind: &kind},
			{Label: "world", Kind: &kind},
		},
	}, nil
}

func boolPtr(b bool) *bool { return &b }

func TestHarness(t *testing.T) {
	tests := []struct {
		name string
		fn   func(t *testing.T, h *servertest.Harness)
	}{
		{
			name: "initialize returns capabilities",
			fn: func(t *testing.T, h *servertest.Harness) {
				if h.InitResult == nil {
					t.Fatal("InitResult is nil")
				}
				if h.InitResult.ServerInfo == nil || h.InitResult.ServerInfo.Name != "test-server" {
					t.Fatalf("unexpected server name: %+v", h.InitResult.ServerInfo)
				}
				if h.InitResult.Capabilities.HoverProvider == nil || !*h.InitResult.Capabilities.HoverProvider {
					t.Fatal("expected hover provider capability")
				}
			},
		},
		{
			name: "didOpen and hover",
			fn: func(t *testing.T, h *servertest.Harness) {
				uri := lsp.DocumentURI("file:///test.txt")
				if err := h.DidOpen(uri, "plaintext", "first line\nsecond line"); err != nil {
					t.Fatal(err)
				}

				hover, err := h.Hover(uri, 0, 0)
				if err != nil {
					t.Fatal(err)
				}
				if hover == nil {
					t.Fatal("hover is nil")
				}
				if !strings.Contains(hover.Contents.Value, "first line") {
					t.Fatalf("unexpected hover content: %s", hover.Contents.Value)
				}
			},
		},
		{
			name: "didSave publishes diagnostics",
			fn: func(t *testing.T, h *servertest.Harness) {
				uri := lsp.DocumentURI("file:///test.txt")
				if err := h.DidOpen(uri, "plaintext", "ERROR bad line\ngood line"); err != nil {
					t.Fatal(err)
				}
				if err := h.DidSave(uri); err != nil {
					t.Fatal(err)
				}

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				diags, err := h.WaitForDiagnostics(ctx, uri)
				if err != nil {
					t.Fatal(err)
				}
				if len(diags) != 1 {
					t.Fatalf("expected 1 diagnostic, got %d", len(diags))
				}
				if !strings.Contains(diags[0].Message, "error on line 1") {
					t.Fatalf("unexpected diagnostic message: %s", diags[0].Message)
				}
			},
		},
		{
			name: "didClose",
			fn: func(t *testing.T, h *servertest.Harness) {
				uri := lsp.DocumentURI("file:///test.txt")
				if err := h.DidOpen(uri, "plaintext", "content"); err != nil {
					t.Fatal(err)
				}
				if err := h.DidClose(uri); err != nil {
					t.Fatal(err)
				}
				// Hover on a closed doc should return nil.
				hover, err := h.Hover(uri, 0, 0)
				if err != nil {
					t.Fatal(err)
				}
				if hover != nil {
					t.Fatal("expected nil hover after close")
				}
			},
		},
		{
			name: "completion returns items",
			fn: func(t *testing.T, h *servertest.Harness) {
				uri := lsp.DocumentURI("file:///test.txt")
				if err := h.DidOpen(uri, "plaintext", ""); err != nil {
					t.Fatal(err)
				}

				list, err := h.Completion(uri, 0, 0)
				if err != nil {
					t.Fatal(err)
				}
				if list == nil {
					t.Fatal("completion list is nil")
				}
				if len(list.Items) != 2 {
					t.Fatalf("expected 2 items, got %d", len(list.Items))
				}
				if list.Items[0].Label != "hello" {
					t.Fatalf("unexpected first item: %s", list.Items[0].Label)
				}
			},
		},
		{
			name: "didChange updates content",
			fn: func(t *testing.T, h *servertest.Harness) {
				uri := lsp.DocumentURI("file:///test.txt")
				if err := h.DidOpen(uri, "plaintext", "original"); err != nil {
					t.Fatal(err)
				}
				if err := h.DidChange(uri, 2, "updated content"); err != nil {
					t.Fatal(err)
				}

				hover, err := h.Hover(uri, 0, 0)
				if err != nil {
					t.Fatal(err)
				}
				if hover == nil {
					t.Fatal("hover is nil after change")
				}
				if !strings.Contains(hover.Contents.Value, "updated content") {
					t.Fatalf("hover should reflect updated content: %s", hover.Contents.Value)
				}
			},
		},
		{
			name: "cleanup runs without error",
			fn: func(_ *testing.T, _ *servertest.Harness) {
				// The test cleanup (shutdown+exit) is automatically verified
				// by the harness t.Cleanup. If it panics or hangs, this test fails.
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := newTestHandler()
			h := servertest.New(t, handler)
			tt.fn(t, h)
		})
	}
}
