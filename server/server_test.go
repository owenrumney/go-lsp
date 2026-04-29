package server

import (
	"context"
	"encoding/json"
	"io"
	"testing"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

// mockHandler implements LifecycleHandler, TextDocumentSyncHandler, and HoverHandler.
type mockHandler struct {
	initialized bool
	opened      []lsp.DocumentURI
}

func (m *mockHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
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

func (m *mockHandler) Hover(_ context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
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

type richCapabilityHandler struct {
	mockHandler
}

func (h *richCapabilityHandler) Completion(_ context.Context, _ *lsp.CompletionParams) (*lsp.CompletionList, error) {
	return &lsp.CompletionList{}, nil
}

func (h *richCapabilityHandler) ResolveCompletionItem(_ context.Context, item *lsp.CompletionItem) (*lsp.CompletionItem, error) {
	return item, nil
}

func (h *richCapabilityHandler) SignatureHelp(_ context.Context, _ *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error) {
	return &lsp.SignatureHelp{}, nil
}

func (h *richCapabilityHandler) CodeAction(_ context.Context, _ *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	return nil, nil
}

func (h *richCapabilityHandler) ResolveCodeAction(_ context.Context, action *lsp.CodeAction) (*lsp.CodeAction, error) {
	return action, nil
}

func (h *richCapabilityHandler) ExecuteCommand(_ context.Context, _ *lsp.ExecuteCommandParams) (any, error) {
	return nil, nil
}

func (h *richCapabilityHandler) SemanticTokensFull(_ context.Context, _ *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	return &lsp.SemanticTokens{}, nil
}

func (h *richCapabilityHandler) SemanticTokensDelta(_ context.Context, _ *lsp.SemanticTokensDeltaParams) (*lsp.SemanticTokensDelta, error) {
	return &lsp.SemanticTokensDelta{}, nil
}

func (h *richCapabilityHandler) SemanticTokensRange(_ context.Context, _ *lsp.SemanticTokensRangeParams) (*lsp.SemanticTokens, error) {
	return &lsp.SemanticTokens{}, nil
}

func (h *richCapabilityHandler) WillCreateFiles(_ context.Context, _ *lsp.CreateFilesParams) (*lsp.WorkspaceEdit, error) {
	return nil, nil
}

func TestApplyCapabilityOptions(t *testing.T) {
	h := &richCapabilityHandler{}
	caps := buildCapabilities(h)
	folder := lsp.FileOperationPatternKindFolder
	applyCapabilityOptions(&caps, h, CapabilityOptions{
		PositionEncoding: ptr(lsp.PositionEncodingUTF8),
		Completion: &lsp.CompletionOptions{
			TriggerCharacters: []string{"."},
		},
		SignatureHelp: &lsp.SignatureHelpOptions{
			TriggerCharacters: []string{"("},
		},
		CodeAction: &lsp.CodeActionOptions{
			CodeActionKinds: []lsp.CodeActionKind{lsp.CodeActionQuickFix},
		},
		ExecuteCommand: &lsp.ExecuteCommandOptions{
			Commands: []string{"test.generateDebugBundle"},
		},
		SemanticTokens: &lsp.SemanticTokensOptions{
			Legend: lsp.SemanticTokensLegend{
				TokenTypes:     []string{"type", "function"},
				TokenModifiers: []string{"declaration"},
			},
		},
		FileOperationFilters: []lsp.FileOperationFilter{
			{
				Pattern: lsp.FileOperationPattern{
					Glob:    "**/*.test",
					Matches: &folder,
				},
			},
		},
	})

	if got := caps.CompletionProvider.TriggerCharacters; len(got) != 1 || got[0] != "." {
		t.Fatalf("completion triggers = %v", got)
	}
	if caps.PositionEncoding == nil || *caps.PositionEncoding != lsp.PositionEncodingUTF8 {
		t.Fatalf("position encoding = %v", caps.PositionEncoding)
	}
	if caps.CompletionProvider.ResolveProvider == nil || !*caps.CompletionProvider.ResolveProvider {
		t.Fatal("expected completion resolve provider")
	}
	if got := caps.SignatureHelpProvider.TriggerCharacters; len(got) != 1 || got[0] != "(" {
		t.Fatalf("signature triggers = %v", got)
	}
	if got := caps.CodeActionProvider.CodeActionKinds; len(got) != 1 || got[0] != lsp.CodeActionQuickFix {
		t.Fatalf("code action kinds = %v", got)
	}
	if caps.CodeActionProvider.ResolveProvider == nil || !*caps.CodeActionProvider.ResolveProvider {
		t.Fatal("expected code action resolve provider")
	}
	if got := caps.ExecuteCommandProvider.Commands; len(got) != 1 || got[0] != "test.generateDebugBundle" {
		t.Fatalf("commands = %v", got)
	}
	if got := caps.SemanticTokensProvider.Legend.TokenTypes; len(got) != 2 || got[0] != "type" {
		t.Fatalf("semantic token legend = %v", got)
	}
	if caps.SemanticTokensProvider.Full == nil || caps.SemanticTokensProvider.Full.Delta == nil || !*caps.SemanticTokensProvider.Full.Delta {
		t.Fatal("expected semantic full delta support")
	}
	if caps.SemanticTokensProvider.Range == nil || !*caps.SemanticTokensProvider.Range {
		t.Fatal("expected semantic range support")
	}
	filters := caps.Workspace.FileOperations.WillCreate.Filters
	if len(filters) != 1 || filters[0].Pattern.Glob != "**/*.test" || filters[0].Pattern.Matches == nil || *filters[0].Pattern.Matches != folder {
		t.Fatalf("file operation filters = %#v", filters)
	}
}

func ptr[T any](v T) *T {
	return &v
}

func TestCapabilityOptionsDoNotAdvertiseMissingHandlers(t *testing.T) {
	h := &mockHandler{}
	caps := buildCapabilities(h)
	applyCapabilityOptions(&caps, h, CapabilityOptions{
		Completion: &lsp.CompletionOptions{TriggerCharacters: []string{"."}},
		ExecuteCommand: &lsp.ExecuteCommandOptions{
			Commands: []string{"test.generateDebugBundle"},
		},
	})

	if caps.CompletionProvider != nil {
		t.Fatal("did not expect completion provider without handler")
	}
	if caps.ExecuteCommandProvider != nil {
		t.Fatal("did not expect execute command provider without handler")
	}
}

func TestMergeCapabilitiesPreservesExplicitInitializeCapabilities(t *testing.T) {
	explicit := &lsp.ServerCapabilities{
		CompletionProvider: &lsp.CompletionOptions{
			TriggerCharacters: []string{"#"},
		},
	}
	auto := &lsp.ServerCapabilities{
		CompletionProvider: &lsp.CompletionOptions{
			TriggerCharacters: []string{"."},
			ResolveProvider:   &enabled,
		},
		ExecuteCommandProvider: &lsp.ExecuteCommandOptions{
			Commands: []string{"test.command"},
		},
	}

	mergeCapabilities(explicit, auto)

	if got := explicit.CompletionProvider.TriggerCharacters; len(got) != 1 || got[0] != "#" {
		t.Fatalf("completion triggers = %v", got)
	}
	if explicit.CompletionProvider.ResolveProvider != nil {
		t.Fatal("explicit completion provider should not be enriched by merge")
	}
	if got := explicit.ExecuteCommandProvider.Commands; len(got) != 1 || got[0] != "test.command" {
		t.Fatalf("commands = %v", got)
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

	ctx, cancel := context.WithCancel(t.Context())
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
