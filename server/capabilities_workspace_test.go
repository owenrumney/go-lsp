package server

import (
	"context"
	"testing"

	"github.com/owenrumney/go-lsp/lsp"
)

type workspaceFoldersHandler struct {
	mockHandler
}

func (h *workspaceFoldersHandler) DidChangeWorkspaceFolders(_ context.Context, _ *lsp.DidChangeWorkspaceFoldersParams) error {
	return nil
}

func TestBuildCapabilitiesAdvertisesWorkspaceFolders(t *testing.T) {
	caps := buildCapabilities(&workspaceFoldersHandler{})

	if caps.Workspace == nil || caps.Workspace.WorkspaceFolders == nil {
		t.Fatal("expected workspace.workspaceFolders capability")
	}
	wf := caps.Workspace.WorkspaceFolders
	if wf.Supported == nil || !*wf.Supported {
		t.Error("expected workspaceFolders.supported = true")
	}
	if wf.ChangeNotifications == nil || !*wf.ChangeNotifications {
		t.Error("expected workspaceFolders.changeNotifications = true")
	}
}

type onTypeFormattingHandler struct {
	mockHandler
}

func (h *onTypeFormattingHandler) OnTypeFormatting(_ context.Context, _ *lsp.DocumentOnTypeFormattingParams) ([]lsp.TextEdit, error) {
	return nil, nil
}

func TestOnTypeFormattingCapabilityRequiresOptions(t *testing.T) {
	h := &onTypeFormattingHandler{}

	caps := buildCapabilities(h)
	applyCapabilityOptions(&caps, h, CapabilityOptions{})
	if caps.DocumentOnTypeFormattingProvider != nil {
		t.Fatal("expected no onTypeFormatting capability without options")
	}

	caps = buildCapabilities(h)
	applyCapabilityOptions(&caps, h, CapabilityOptions{
		OnTypeFormatting: &lsp.DocumentOnTypeFormattingOptions{FirstTriggerCharacter: "}"},
	})
	if caps.DocumentOnTypeFormattingProvider == nil {
		t.Fatal("expected onTypeFormatting capability with options")
	}
	if caps.DocumentOnTypeFormattingProvider.FirstTriggerCharacter != "}" {
		t.Errorf("firstTriggerCharacter = %q, want \"}\"", caps.DocumentOnTypeFormattingProvider.FirstTriggerCharacter)
	}

	plain := &mockHandler{}
	caps = buildCapabilities(plain)
	applyCapabilityOptions(&caps, plain, CapabilityOptions{
		OnTypeFormatting: &lsp.DocumentOnTypeFormattingOptions{FirstTriggerCharacter: "}"},
	})
	if caps.DocumentOnTypeFormattingProvider != nil {
		t.Fatal("expected no onTypeFormatting capability without handler")
	}
}
