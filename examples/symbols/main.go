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
	docs *document.Store
}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			DocumentSymbolProvider:  boolPtr(true),
			WorkspaceSymbolProvider: boolPtr(true),
		},
		ServerInfo: &lsp.ServerInfo{Name: "symbols-example", Version: "0.1.0"},
	}, nil
}

func (h *handler) Shutdown(_ context.Context) error { return nil }

// TextDocumentSyncHandler

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

// DocumentSymbolHandler

func (h *handler) DocumentSymbol(_ context.Context, params *lsp.DocumentSymbolParams) ([]lsp.DocumentSymbol, error) {
	text, ok := h.docs.Text(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}
	return parseDocumentSymbols(text), nil
}

// WorkspaceSymbolHandler

func (h *handler) WorkspaceSymbol(_ context.Context, params *lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error) {
	query := strings.ToLower(params.Query)

	var results []lsp.SymbolInformation
	for _, uri := range h.docs.URIs() {
		text, _ := h.docs.Text(uri)
		for _, sym := range parseSymbols(uri, text) {
			if query == "" || strings.Contains(strings.ToLower(sym.Name), query) {
				results = append(results, sym)
			}
		}
	}
	return results, nil
}

// Symbol parsing helpers

type symbolMatch struct {
	prefix string
	kind   lsp.SymbolKind
}

var matchers = []symbolMatch{
	{"func ", lsp.SymbolKindFunction},
	{"def ", lsp.SymbolKindFunction},
	{"class ", lsp.SymbolKindClass},
	{"## ", lsp.SymbolKindString},
}

func parseDocumentSymbols(text string) []lsp.DocumentSymbol {
	var syms []lsp.DocumentSymbol
	for i, line := range strings.Split(text, "\n") {
		name, kind, ok := matchLine(line)
		if !ok {
			continue
		}
		r := lsp.Range{
			Start: lsp.Position{Line: i, Character: 0},
			End:   lsp.Position{Line: i, Character: len(line)},
		}
		syms = append(syms, lsp.DocumentSymbol{
			Name:           name,
			Kind:           kind,
			Range:          r,
			SelectionRange: r,
		})
	}
	return syms
}

func parseSymbols(uri lsp.DocumentURI, text string) []lsp.SymbolInformation {
	var syms []lsp.SymbolInformation
	for i, line := range strings.Split(text, "\n") {
		name, kind, ok := matchLine(line)
		if !ok {
			continue
		}
		syms = append(syms, lsp.SymbolInformation{
			Name: name,
			Kind: kind,
			Location: lsp.Location{
				URI: uri,
				Range: lsp.Range{
					Start: lsp.Position{Line: i, Character: 0},
					End:   lsp.Position{Line: i, Character: len(line)},
				},
			},
		})
	}
	return syms
}

func matchLine(line string) (name string, kind lsp.SymbolKind, ok bool) {
	trimmed := strings.TrimSpace(line)
	for _, m := range matchers {
		if strings.HasPrefix(trimmed, m.prefix) {
			name = strings.TrimSpace(trimmed[len(m.prefix):])
			// Strip trailing punctuation like ":" or "{"
			name = strings.TrimRight(name, " :{(")
			if name != "" {
				return name, m.kind, true
			}
		}
	}
	return "", 0, false
}

func boolPtr(b bool) *bool { return &b }

func main() {
	srv := server.NewServer(&handler{docs: document.NewStore()})
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
