package main

import (
	"context"
	"log"
	"strings"

	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

type handler struct {
	docs   *document.Store
	client *server.Client
}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		ServerInfo: &lsp.ServerInfo{Name: "diagnostics-example", Version: "0.1.0"},
	}, nil
}

func (h *handler) Shutdown(_ context.Context) error { return nil }

func (h *handler) SetClient(client *server.Client) {
	h.client = client
}

func (h *handler) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
	_, err := h.docs.Open(params)
	return err
}

func (h *handler) DidChange(_ context.Context, params *lsp.DidChangeTextDocumentParams) error {
	_, err := h.docs.Change(params)
	return err
}

func (h *handler) DidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
	h.docs.Close(params)
	return nil
}

func (h *handler) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error {
	_ = h.client.LogMessage(ctx, &lsp.LogMessageParams{
		Type:    lsp.MessageTypeInfo,
		Message: "document saved: " + string(params.TextDocument.URI),
	})

	var diags []lsp.Diagnostic
	text, ok := h.docs.Text(params.TextDocument.URI)
	if ok {
		for i, line := range strings.Split(text, "\n") {
			idx := strings.Index(line, "TODO")
			if idx < 0 {
				continue
			}
			sev := lsp.SeverityWarning
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
	}

	return h.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diags,
	})
}

func main() {
	h := &handler{docs: document.NewStore()}
	srv := server.NewServer(h)
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
