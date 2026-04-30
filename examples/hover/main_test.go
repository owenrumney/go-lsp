package main

import (
	"strings"
	"testing"

	"github.com/owenrumney/go-lsp/servertest"
)

func TestHoverExample(t *testing.T) {
	h := servertest.New(t, &handler{})

	if err := h.DidOpen("file:///test.txt", "plaintext", "hello"); err != nil {
		t.Fatal(err)
	}

	hover, err := h.Hover("file:///test.txt", 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if hover == nil {
		t.Fatal("expected hover result")
	}
	if !strings.Contains(hover.Contents.Value, "go-lsp") {
		t.Fatalf("hover = %q", hover.Contents.Value)
	}
}
