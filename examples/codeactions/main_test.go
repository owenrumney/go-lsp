package main

import (
	"context"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/servertest"
)

func TestCodeActionsExample(t *testing.T) {
	h := servertest.New(t, &handler{docs: document.NewStore()})
	uri := lsp.DocumentURI("file:///test.txt")

	if err := h.DidOpen(uri, "plaintext", "trim me   \nclean"); err != nil {
		t.Fatal(err)
	}
	if err := h.DidSave(uri); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	diags, err := h.WaitForDiagnostics(ctx, uri)
	if err != nil {
		t.Fatal(err)
	}
	if len(diags) != 1 {
		t.Fatalf("diagnostics = %d, want 1", len(diags))
	}

	actions, err := h.CodeAction(&lsp.CodeActionParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Range:        diags[0].Range,
		Context:      lsp.CodeActionContext{Diagnostics: diags},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(actions) != 1 {
		t.Fatalf("actions = %d, want 1", len(actions))
	}
	if actions[0].Edit == nil || len(actions[0].Edit.Changes[uri]) != 1 {
		t.Fatalf("action edit = %+v", actions[0].Edit)
	}
}
