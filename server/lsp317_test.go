package server

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/owenrumney/go-lsp/lsp"
)

type lsp317Handler struct {
	mockHandler
	resolved      bool
	created       []lsp.DocumentURI
	notebookOpens []lsp.DocumentURI
	cancelled     bool
}

func (h *lsp317Handler) WorkspaceSymbols(_ context.Context, _ *lsp.WorkspaceSymbolParams) ([]lsp.WorkspaceSymbol, error) {
	return []lsp.WorkspaceSymbol{{
		Name:     "sym",
		Kind:     lsp.SymbolKindFunction,
		Location: lsp.WorkspaceSymbolLocation{URI: "file:///a.go"},
	}}, nil
}

func (h *lsp317Handler) ResolveWorkspaceSymbol(_ context.Context, s *lsp.WorkspaceSymbol) (*lsp.WorkspaceSymbol, error) {
	h.resolved = true
	s.Location.Range = &lsp.Range{}
	return s, nil
}

func (h *lsp317Handler) DidCreateFiles(_ context.Context, params *lsp.CreateFilesParams) error {
	for _, f := range params.Files {
		h.created = append(h.created, lsp.DocumentURI(f.URI))
	}
	return nil
}

func (h *lsp317Handler) DidRenameFiles(_ context.Context, _ *lsp.RenameFilesParams) error { return nil }
func (h *lsp317Handler) DidDeleteFiles(_ context.Context, _ *lsp.DeleteFilesParams) error { return nil }

func (h *lsp317Handler) WorkDoneProgressCancel(_ context.Context, _ *lsp.WorkDoneProgressCancelParams) error {
	h.cancelled = true
	return nil
}

func (h *lsp317Handler) DidOpenNotebookDocument(_ context.Context, params *lsp.DidOpenNotebookDocumentParams) error {
	h.notebookOpens = append(h.notebookOpens, params.NotebookDocument.URI)
	return nil
}

func (h *lsp317Handler) DidChangeNotebookDocument(_ context.Context, _ *lsp.DidChangeNotebookDocumentParams) error {
	return nil
}

func (h *lsp317Handler) DidSaveNotebookDocument(_ context.Context, _ *lsp.DidSaveNotebookDocumentParams) error {
	return nil
}

func (h *lsp317Handler) DidCloseNotebookDocument(_ context.Context, _ *lsp.DidCloseNotebookDocumentParams) error {
	return nil
}

func TestBuildCapabilities317(t *testing.T) {
	h := &lsp317Handler{}
	caps := buildCapabilities(h)

	if caps.WorkspaceSymbolProvider == nil {
		t.Fatal("expected workspaceSymbolProvider")
	}
	if caps.WorkspaceSymbolProvider.ResolveProvider == nil || !*caps.WorkspaceSymbolProvider.ResolveProvider {
		t.Error("expected workspaceSymbolProvider.resolveProvider")
	}

	fileOps := caps.Workspace.FileOperations
	if fileOps == nil || fileOps.DidCreate == nil || fileOps.DidRename == nil || fileOps.DidDelete == nil {
		t.Errorf("expected did* file operation capabilities, got %+v", fileOps)
	}

	if caps.NotebookDocumentSync != nil {
		t.Error("expected no notebook sync capability without options")
	}
	applyCapabilityOptions(&caps, h, CapabilityOptions{
		NotebookSync: &lsp.NotebookDocumentSyncOptions{
			NotebookSelector: []lsp.NotebookSelector{
				{Notebook: &lsp.NotebookDocumentFilter{NotebookType: "jupyter-notebook"}},
			},
		},
	})
	if caps.NotebookDocumentSync == nil {
		t.Error("expected notebook sync capability with options")
	}
}

func TestLSP317Dispatch(t *testing.T) {
	h := &lsp317Handler{}
	c := startLifecycleServer(t, h)

	if resp := c.call(1, "initialize", &lsp.InitializeParams{}); resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}

	if resp := c.call(2, "workspace/symbol", &lsp.WorkspaceSymbolParams{Query: "s"}); resp.Error != nil {
		t.Fatalf("workspace/symbol failed: %v", resp.Error)
	}
	if resp := c.call(3, "workspaceSymbol/resolve", &lsp.WorkspaceSymbol{Name: "sym"}); resp.Error != nil {
		t.Fatalf("workspaceSymbol/resolve failed: %v", resp.Error)
	}
	if !h.resolved {
		t.Error("expected ResolveWorkspaceSymbol to be called")
	}

	c.notify("workspace/didCreateFiles", &lsp.CreateFilesParams{
		Files: []lsp.FileCreate{{URI: "file:///new.go"}},
	})
	c.notify("window/workDoneProgress/cancel", &lsp.WorkDoneProgressCancelParams{Token: lsp.ProgressToken(`"t"`)})
	c.notify("notebookDocument/didOpen", &lsp.DidOpenNotebookDocumentParams{
		NotebookDocument: lsp.NotebookDocument{URI: "file:///nb.ipynb", NotebookType: "jupyter-notebook", Version: 1},
	})

	// The shutdown round-trip is the sync barrier for the notifications.
	if resp := c.call(4, "shutdown", nil); resp.Error != nil {
		t.Fatalf("shutdown failed: %v", resp.Error)
	}

	if len(h.created) != 1 || h.created[0] != "file:///new.go" {
		t.Errorf("didCreateFiles = %v", h.created)
	}
	if !h.cancelled {
		t.Error("expected WorkDoneProgressCancel to be called")
	}
	if len(h.notebookOpens) != 1 || h.notebookOpens[0] != "file:///nb.ipynb" {
		t.Errorf("notebook opens = %v", h.notebookOpens)
	}
}

func TestInitializeAdvertisesNotebookSync(t *testing.T) {
	c := startLifecycleServer(t, &lsp317Handler{}, WithNotebookSyncOptions(lsp.NotebookDocumentSyncOptions{
		NotebookSelector: []lsp.NotebookSelector{
			{Notebook: &lsp.NotebookDocumentFilter{NotebookType: "jupyter-notebook"}},
		},
	}))

	resp := c.call(1, "initialize", &lsp.InitializeParams{})
	if resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}
	var result lsp.InitializeResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatal(err)
	}
	sync := result.Capabilities.NotebookDocumentSync
	if sync == nil || len(sync.NotebookSelector) != 1 {
		t.Fatalf("notebookDocumentSync = %+v, want configured selector", sync)
	}
	if sync.NotebookSelector[0].Notebook.NotebookType != "jupyter-notebook" {
		t.Fatalf("selector = %+v", sync.NotebookSelector[0])
	}
}
