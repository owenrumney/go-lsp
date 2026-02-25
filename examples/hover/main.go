package main

import (
	"context"
	"log"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

type handler struct{}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{
		ServerInfo: &lsp.ServerInfo{Name: "hover-example", Version: "0.1.0"},
	}, nil
}

func (h *handler) Shutdown(_ context.Context) error { return nil }

func (h *handler) Hover(_ context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	return &lsp.Hover{
		Contents: lsp.MarkupContent{
			Kind:  lsp.Markdown,
			Value: "**go-lsp** hover example\n\nYou hovered on a symbol.",
		},
	}, nil
}

func main() {
	srv := server.NewServer(&handler{})
	if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
		log.Fatal(err)
	}
}
