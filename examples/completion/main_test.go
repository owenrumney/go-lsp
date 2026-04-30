package main

import (
	"testing"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/servertest"
)

func TestCompletionExample(t *testing.T) {
	h := servertest.New(t, &handler{})

	if err := h.DidOpen("file:///test.go", "go", "f"); err != nil {
		t.Fatal(err)
	}

	list, err := h.Completion("file:///test.go", 0, 1)
	if err != nil {
		t.Fatal(err)
	}
	if list == nil {
		t.Fatal("expected completion list")
	}
	if len(list.Items) != len(keywords) {
		t.Fatalf("completion items = %d, want %d", len(list.Items), len(keywords))
	}
	if list.Items[0].Kind == nil || *list.Items[0].Kind != lsp.CompletionItemKindKeyword {
		t.Fatalf("first item kind = %v", list.Items[0].Kind)
	}
}
