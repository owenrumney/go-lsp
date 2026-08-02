package server

import (
	"context"
	"testing"

	"github.com/owenrumney/go-lsp/lsp"
)

type linkDefinitionHandler struct {
	mockHandler
	linkCalled bool
}

func (h *linkDefinitionHandler) Definition(_ context.Context, _ *lsp.DefinitionParams) ([]lsp.Location, error) {
	return []lsp.Location{{URI: "file:///loc.go"}}, nil
}

func (h *linkDefinitionHandler) DefinitionLinks(_ context.Context, _ *lsp.DefinitionParams) ([]lsp.LocationLink, error) {
	h.linkCalled = true
	return []lsp.LocationLink{{TargetURI: "file:///link.go"}}, nil
}

func TestDefinitionLinkHandlerTakesPrecedence(t *testing.T) {
	h := &linkDefinitionHandler{}
	c := startLifecycleServer(t, h)

	if resp := c.call(1, "initialize", &lsp.InitializeParams{}); resp.Error != nil {
		t.Fatalf("initialize failed: %v", resp.Error)
	}
	resp := c.call(2, "textDocument/definition", &lsp.DefinitionParams{})
	if resp.Error != nil {
		t.Fatalf("definition failed: %v", resp.Error)
	}
	if !h.linkCalled {
		t.Error("expected DefinitionLinks to be called, not Definition")
	}
}

func TestLinkHandlerAdvertisesCapability(t *testing.T) {
	type linkOnly struct {
		mockHandler
		DefinitionLinkHandler
	}
	caps := buildCapabilities(&linkOnly{})
	if caps.DefinitionProvider == nil || !*caps.DefinitionProvider {
		t.Error("expected definitionProvider for link-only handler")
	}
}
