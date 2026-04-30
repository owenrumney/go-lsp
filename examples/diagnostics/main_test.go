package main

import (
	"context"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/servertest"
)

func TestDiagnosticsExample(t *testing.T) {
	h := servertest.New(t, &handler{docs: document.NewStore()})
	uri := lsp.DocumentURI("file:///test.txt")

	if err := h.DidOpen(uri, "plaintext", "first\nTODO fix this"); err != nil {
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
	if diags[0].Source != "todo-checker" || diags[0].Message != "TODO found" {
		t.Fatalf("diagnostic = %+v", diags[0])
	}
	if logs := h.LogMessages(); len(logs) != 1 {
		t.Fatalf("log messages = %d, want 1", len(logs))
	}
}
