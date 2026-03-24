package main

import (
	"context"
	"log"
	"strings"
	"sync"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

type handler struct {
	mu   sync.Mutex
	docs map[lsp.DocumentURI]string // URI -> full text
}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			TextDocumentSync: &lsp.TextDocumentSyncOptions{
				OpenClose: boolPtr(true),
				Change:    lsp.SyncFull,
			},
			DocumentSymbolProvider:  boolPtr(true),
			WorkspaceSymbolProvider: boolPtr(true),
		},
		ServerInfo: &lsp.ServerInfo{Name: "symbols-example", Version: "0.1.0"},
	}, nil
}

func (h *handler) Shutdown(_ context.Context) error { return nil }

// TextDocumentSyncHandler

func (h *handler) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.docs[params.TextDocument.URI] = params.TextDocument.Text
	return nil
}

func (h *handler) DidChange(_ context.Context, params *lsp.DidChangeTextDocumentParams) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(params.ContentChanges) > 0 {
		h.docs[params.TextDocument.URI] = params.ContentChanges[len(params.ContentChanges)-1].Text
	}
	return nil
}

func (h *handler) DidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.docs, params.TextDocument.URI)
	return nil
}

// DocumentSymbolHandler

func (h *handler) DocumentSymbol(_ context.Context, params *lsp.DocumentSymbolParams) ([]lsp.DocumentSymbol, error) {
	h.mu.Lock()
	text, ok := h.docs[params.TextDocument.URI]
	h.mu.Unlock()
	if !ok {
		return nil, nil
	}
	return parseDocumentSymbols(text), nil
}

// WorkspaceSymbolHandler

func (h *handler) WorkspaceSymbol(_ context.Context, params *lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error) {
	query := strings.ToLower(params.Query)

	h.mu.Lock()
	defer h.mu.Unlock()

	var results []lsp.SymbolInformation
	for uri, text := range h.docs {
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
	srv := server.NewServer(&handler{docs: make(map[lsp.DocumentURI]string)})
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
