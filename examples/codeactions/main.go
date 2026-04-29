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
	kind := lsp.CodeActionQuickFix
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			CodeActionProvider: &lsp.CodeActionOptions{
				CodeActionKinds: []lsp.CodeActionKind{kind},
			},
		},
		ServerInfo: &lsp.ServerInfo{Name: "codeactions-example", Version: "0.1.0"},
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
	text, _ := h.docs.Text(params.TextDocument.URI)
	diags := findTrailingWhitespace(text)
	return h.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
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

func findTrailingWhitespace(text string) []lsp.Diagnostic {
	var diags []lsp.Diagnostic
	sev := lsp.SeverityWarning

	for i, line := range strings.Split(text, "\n") {
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

func main() {
	h := &handler{docs: document.NewStore()}
	srv := server.NewServer(h)
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
