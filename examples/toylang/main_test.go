package main

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/servertest"
)

func TestToylangExample(t *testing.T) {
	h := servertest.New(t, &handler{docs: document.NewStore()})
	uri := lsp.DocumentURI("file:///app.env")
	text := "PORT=8080\nHOST=localhost\nPORT=9090"

	if err := h.DidOpen(uri, "env", text); err != nil {
		t.Fatal(err)
	}

	hover, err := h.Hover(uri, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if hover == nil || !strings.Contains(hover.Contents.Value, "8080") {
		t.Fatalf("hover = %+v", hover)
	}

	list, err := h.Completion(uri, 0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if list == nil || len(list.Items) != 2 {
		t.Fatalf("completion = %+v", list)
	}

	locs, err := h.Definition(uri, 2, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(locs) != 2 {
		t.Fatalf("definitions = %d, want 2", len(locs))
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
	if len(diags) != 1 || !strings.Contains(diags[0].Message, "duplicate key") {
		t.Fatalf("diagnostics = %+v", diags)
	}
}
