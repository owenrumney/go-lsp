package main

import (
	"context"
	"fmt"
	"log"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

var keywords = []string{"func", "var", "const", "type"}

type handler struct{}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		Capabilities: lsp.ServerCapabilities{
			CompletionProvider: &lsp.CompletionOptions{
				ResolveProvider: boolPtr(true),
			},
		},
		ServerInfo: &lsp.ServerInfo{Name: "completion-example", Version: "0.1.0"},
	}, nil
}

func (h *handler) Shutdown(_ context.Context) error { return nil }

func (h *handler) Completion(_ context.Context, _ *lsp.CompletionParams) (*lsp.CompletionList, error) {
	kind := lsp.CompletionItemKindKeyword
	items := make([]lsp.CompletionItem, len(keywords))
	for i, kw := range keywords {
		items[i] = lsp.CompletionItem{
			Label: kw,
			Kind:  &kind,
		}
	}
	return &lsp.CompletionList{Items: items}, nil
}

func (h *handler) ResolveCompletionItem(_ context.Context, item *lsp.CompletionItem) (*lsp.CompletionItem, error) {
	item.Documentation = &lsp.MarkupContent{
		Kind:  lsp.Markdown,
		Value: fmt.Sprintf("The **%s** keyword in Go.", item.Label),
	}
	return item, nil
}

func boolPtr(b bool) *bool { return &b }

func main() {
	srv := server.NewServer(&handler{})
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
