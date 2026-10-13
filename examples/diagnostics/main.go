package main

import (
	"context"
	"log"
	"strings"

	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

type handler struct{}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		ServerInfo: &lsp.ServerInfo{Name: "diagnostics-example", Version: "0.1.0"},
	}, nil
}

func (h *handler) Shutdown(_ context.Context) error { return nil }

// Diagnose flags every TODO marker in the document.
func (h *handler) Diagnose(_ context.Context, doc *document.Document) ([]lsp.Diagnostic, error) {
	var diags []lsp.Diagnostic
	sev := lsp.SeverityWarning

	for i, line := range strings.Split(doc.Text(), "\n") {
		idx := strings.Index(line, "TODO")
		if idx < 0 {
			continue
		}
		diags = append(diags, lsp.Diagnostic{
			Range: lsp.Range{
				Start: lsp.Position{Line: i, Character: idx},
				End:   lsp.Position{Line: i, Character: idx + 4},
			},
			Severity: &sev,
			Source:   "todo-checker",
			Message:  "TODO found",
		})
	}
	return diags, nil
}

func main() {
	srv := server.NewServer(&handler{}, server.WithDocumentStore(document.NewStore()))
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
