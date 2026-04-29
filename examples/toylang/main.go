package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

// entry represents a key=value pair at a specific location.
type entry struct {
	key   string
	value string
	line  int
}

// handler implements a toy language server for .env-style key=value config files.
type handler struct {
	docs   *document.Store
	client *server.Client
}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			HoverProvider:      boolPtr(true),
			DefinitionProvider: boolPtr(true),
			CompletionProvider: &lsp.CompletionOptions{},
		},
		ServerInfo: &lsp.ServerInfo{Name: "toylang-env", Version: "0.1.0"},
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
	text, ok := h.docs.Text(params.TextDocument.URI)
	if !ok {
		return nil
	}

	diags := h.findDuplicates(text)
	return h.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
		URI:         params.TextDocument.URI,
		Diagnostics: diags,
	})
}

// Hover returns the value of the key under the cursor.
func (h *handler) Hover(_ context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
	text, ok := h.docs.Text(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}

	entry := h.entryAtLine(text, params.Position.Line)
	if entry == nil {
		return nil, nil
	}

	return &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.Markdown,
			Value: fmt.Sprintf("**%s** = `%s`", entry.key, entry.value),
		},
	}, nil
}

// Completion suggests known keys from all open documents.
func (h *handler) Completion(_ context.Context, _ *lsp.CompletionParams) (*lsp.CompletionList, error) {
	seen := map[string]string{}
	for _, uri := range h.docs.URIs() {
		text, _ := h.docs.Text(uri)
		for _, e := range h.parseEntries(text) {
			if _, exists := seen[e.key]; !exists {
				seen[e.key] = e.value
			}
		}
	}

	kind := lsp.CompletionItemKindVariable
	items := make([]lsp.CompletionItem, 0, len(seen))
	for key, val := range seen {
		items = append(items, lsp.CompletionItem{
			Label:  key,
			Kind:   &kind,
			Detail: fmt.Sprintf("= %s", val),
		})
	}
	return &lsp.CompletionList{Items: items}, nil
}

// Definition jumps to where a key is defined across all open files.
func (h *handler) Definition(_ context.Context, params *lsp.DefinitionParams) ([]lsp.Location, error) {
	text, ok := h.docs.Text(params.TextDocument.URI)
	if !ok {
		return nil, nil
	}

	entry := h.entryAtLine(text, params.Position.Line)
	if entry == nil {
		return nil, nil
	}

	var locations []lsp.Location
	for _, uri := range h.docs.URIs() {
		docText, _ := h.docs.Text(uri)
		for _, e := range h.parseEntries(docText) {
			if e.key == entry.key {
				locations = append(locations, lsp.Location{
					URI: uri,
					Range: lsp.Range{
						Start: lsp.Position{Line: e.line, Character: 0},
						End:   lsp.Position{Line: e.line, Character: len(e.key)},
					},
				})
			}
		}
	}
	return locations, nil
}

// parseEntries parses all KEY=VALUE lines from the given text.
func (h *handler) parseEntries(text string) []entry {
	var entries []entry
	for i, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, "="); idx > 0 {
			entries = append(entries, entry{
				key:   strings.TrimSpace(line[:idx]),
				value: strings.TrimSpace(line[idx+1:]),
				line:  i,
			})
		}
	}
	return entries
}

// entryAtLine returns the entry on the given line, or nil.
func (h *handler) entryAtLine(text string, line int) *entry {
	lines := strings.Split(text, "\n")
	if line < 0 || line >= len(lines) {
		return nil
	}
	l := strings.TrimSpace(lines[line])
	if idx := strings.Index(l, "="); idx > 0 {
		e := entry{
			key:   strings.TrimSpace(l[:idx]),
			value: strings.TrimSpace(l[idx+1:]),
			line:  line,
		}
		return &e
	}
	return nil
}

// findDuplicates returns diagnostics for duplicate keys in the document.
func (h *handler) findDuplicates(text string) []lsp.Diagnostic {
	entries := h.parseEntries(text)
	seen := map[string]int{} // key -> first line
	var diags []lsp.Diagnostic
	sev := lsp.SeverityWarning

	for _, e := range entries {
		if firstLine, exists := seen[e.key]; exists {
			diags = append(diags, lsp.Diagnostic{
				Range: lsp.Range{
					Start: lsp.Position{Line: e.line, Character: 0},
					End:   lsp.Position{Line: e.line, Character: len(e.key)},
				},
				Severity: &sev,
				Source:   "toylang-env",
				Message:  fmt.Sprintf("duplicate key %q (first defined on line %d)", e.key, firstLine+1),
			})
		} else {
			seen[e.key] = e.line
		}
	}
	return diags
}

func boolPtr(b bool) *bool { return &b }

func main() {
	h := &handler{docs: document.NewStore()}
	srv := server.NewServer(h)
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
