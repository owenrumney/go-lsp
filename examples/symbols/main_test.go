package main

import (
	"testing"

	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/servertest"
)

func TestSymbolsExample(t *testing.T) {
	h := servertest.New(t, &handler{docs: document.NewStore()})
	uri := lsp.DocumentURI("file:///test.toy")
	text := "class App:\nfunc run\n## Section"

	if err := h.DidOpen(uri, "toy", text); err != nil {
		t.Fatal(err)
	}

	docSymbols, err := h.DocumentSymbol(uri)
	if err != nil {
		t.Fatal(err)
	}
	if len(docSymbols) != 3 {
		t.Fatalf("document symbols = %d, want 3", len(docSymbols))
	}
	if docSymbols[0].Name != "App" {
		t.Fatalf("first symbol = %+v", docSymbols[0])
	}

	workspaceSymbols, err := h.WorkspaceSymbol("run")
	if err != nil {
		t.Fatal(err)
	}
	if len(workspaceSymbols) != 1 || workspaceSymbols[0].Name != "run" {
		t.Fatalf("workspace symbols = %+v", workspaceSymbols)
	}
}
