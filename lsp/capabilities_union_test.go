package lsp

import (
	"encoding/json"
	"testing"
)

// goplsCapabilities is the shape gopls v0.23.0 actually sends. Every union arm
// here broke parsing before the unmarshallers in capabilities_union.go.
const goplsCapabilities = `{
  "codeActionProvider": true,
  "renameProvider": true,
  "hoverProvider": true,
  "definitionProvider": true,
  "documentSymbolProvider": true,
  "workspaceSymbolProvider": true,
  "inlayHintProvider": true,
  "foldingRangeProvider": {"workDoneProgress": true},
  "semanticTokensProvider": {
    "legend": {"tokenTypes": ["keyword"], "tokenModifiers": []},
    "full": true,
    "range": true
  },
  "workspace": {
    "workspaceFolders": {"supported": true, "changeNotifications": "workspace/didChangeWorkspaceFolders"}
  }
}`

func TestServerCapabilitiesAcceptsGoplsUnions(t *testing.T) {
	var caps ServerCapabilities
	if err := json.Unmarshal([]byte(goplsCapabilities), &caps); err != nil {
		t.Fatalf("unmarshalling gopls capabilities: %v", err)
	}

	if !caps.CodeActionProvider.Enabled() {
		t.Error(`codeActionProvider: true should be enabled`)
	}
	if !caps.RenameProvider.Enabled() {
		t.Error(`renameProvider: true should be enabled`)
	}
	if !caps.InlayHintProvider.Enabled() {
		t.Error(`inlayHintProvider: true should be enabled`)
	}
	if caps.HoverProvider == nil || !*caps.HoverProvider {
		t.Error(`hoverProvider: true should be set`)
	}
	if caps.FoldingRangeProvider == nil || !*caps.FoldingRangeProvider {
		t.Error("an options object on a bool provider should read as true")
	}
	if caps.SemanticTokensProvider == nil {
		t.Fatal("semanticTokensProvider missing")
	}
	if caps.SemanticTokensProvider.Full == nil {
		t.Error(`semanticTokens full: true should become empty options`)
	}
	if caps.SemanticTokensProvider.Range == nil || !*caps.SemanticTokensProvider.Range {
		t.Error("semanticTokens range: true should be set")
	}
	if len(caps.SemanticTokensProvider.Legend.TokenTypes) != 1 {
		t.Error("the legend was dropped while coercing the unions around it")
	}

	folders := caps.Workspace.WorkspaceFolders
	if folders == nil {
		t.Fatal("workspaceFolders missing")
	}
	if folders.ChangeNotifications == nil || !*folders.ChangeNotifications {
		t.Error("a changeNotifications string should read as true")
	}
	if folders.ChangeNotificationsID != "workspace/didChangeWorkspaceFolders" {
		t.Errorf("ChangeNotificationsID = %q, want the registration ID", folders.ChangeNotificationsID)
	}
}

func TestServerCapabilitiesFalseMeansUnsupported(t *testing.T) {
	const body = `{
	  "codeActionProvider": false,
	  "hoverProvider": false,
	  "semanticTokensProvider": {"legend": {"tokenTypes": [], "tokenModifiers": []}, "full": false}
	}`

	var caps ServerCapabilities
	if err := json.Unmarshal([]byte(body), &caps); err != nil {
		t.Fatal(err)
	}

	if caps.CodeActionProvider.Enabled() {
		t.Error("codeActionProvider: false must not be enabled")
	}
	if caps.HoverProvider != nil {
		t.Error("hoverProvider: false should be nil, so nil and false checks agree")
	}
	if caps.SemanticTokensProvider.Full != nil {
		t.Error("semanticTokens full: false should be nil")
	}
}
