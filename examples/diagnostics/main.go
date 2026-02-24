package main

import (
	"context"
	"log"
	"strings"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

type handler struct {
	srv *server.Server
}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptions{
				OpenClose: boolPtr(true),
				Change:    lsp.SyncFull,
				Save:      &lsp.SaveOptions{IncludeText: boolPtr(true)},
			},
		},
		ServerInfo: &lsp.ServerInfo{Name: "diagnostics-example", Version: "0.1.0"},
	}, nil
}

func (h *handler) Shutdown(_ context.Context) error { return nil }

func (h *handler) DidOpen(_ context.Context, _ *lsp.DidOpenTextDocumentParams) error  { return nil }
func (h *handler) DidChange(_ context.Context, _ *lsp.DidChangeTextDocumentParams) error { return nil }
func (h *handler) DidClose(_ context.Context, _ *lsp.DidCloseTextDocumentParams) error   { return nil }

func (h *handler) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error {
	_ = h.srv.Client.LogMessage(ctx, &lsp.LogMessageParams{
		Type:    lsp.MessageTypeInfo,
		Message: "document saved: " + string(params.TextDocument.URI),
	})

	var diags []lsp.Diagnostic
	if params.Text != nil {
		for i, line := range strings.Split(*params.Text, "\n") {
			if idx := strings.Index(line, "TODO"); idx >= 0 {
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
	}

	return h.srv.Client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diags,
	})
}

func boolPtr(b bool) *bool { return &b }

func main() {
	h := &handler{}
	srv := server.NewServer(h)
	h.srv = srv
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
