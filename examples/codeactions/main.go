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
	kind := lsp.CodeActionQuickFix
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptions{
				OpenClose: boolPtr(true),
				Change:    lsp.SyncFull,
				Save:      &lsp.SaveOptions{IncludeText: boolPtr(true)},
			},
			CodeActionProvider: &lsp.CodeActionOptions{
				CodeActionKinds: []lsp.CodeActionKind{kind},
			},
		},
		ServerInfo: &lsp.ServerInfo{Name: "codeactions-example", Version: "0.1.0"},
	}, nil
}

func (h *handler) Shutdown(_ context.Context) error { return nil }

func (h *handler) DidOpen(_ context.Context, _ *lsp.DidOpenTextDocumentParams) error     { return nil }
func (h *handler) DidChange(_ context.Context, _ *lsp.DidChangeTextDocumentParams) error { return nil }
func (h *handler) DidClose(_ context.Context, _ *lsp.DidCloseTextDocumentParams) error   { return nil }

func (h *handler) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error {
	diags := findTrailingWhitespace(params.Text)
	return h.srv.Client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diags,
	})
}

func (h *handler) CodeAction(_ context.Context, params *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	var actions []lsp.CodeAction
	kind := lsp.CodeActionQuickFix

	for _, diag := range params.Context.Diagnostics {
		if diag.Source != "trailing-whitespace" {
			continue
		}
		actions = append(actions, lsp.CodeAction{
			Title:       "Trim trailing whitespace",
			Kind:        &kind,
			Diagnostics: []lsp.Diagnostic{diag},
			Edit: &lsp.WorkspaceEdit{
				Changes: map[lsp.DocumentURI][]lsp.TextEdit{
					params.TextDocument.URI: {
						{
							Range:   diag.Range,
							NewText: "",
						},
					},
				},
			},
		})
	}

	return actions, nil
}

func findTrailingWhitespace(text *string) []lsp.Diagnostic {
	if text == nil {
		return nil
	}

	var diags []lsp.Diagnostic
	sev := lsp.SeverityWarning

	for i, line := range strings.Split(*text, "\n") {
		trimmed := strings.TrimRight(line, " \t")
		if len(trimmed) < len(line) {
			diags = append(diags, lsp.Diagnostic{
				Range: lsp.Range{
					Start: lsp.Position{Line: i, Character: len(trimmed)},
					End:   lsp.Position{Line: i, Character: len(line)},
				},
				Severity: &sev,
				Source:   "trailing-whitespace",
				Message:  "Trailing whitespace",
			})
		}
	}

	return diags
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
