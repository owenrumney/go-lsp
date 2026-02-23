package lsp

import (
	"encoding/json"
	"testing"
)

func TestPositionRoundTrip(t *testing.T) {
	pos := Position{Line: 10, Character: 5}
	data, err := json.Marshal(pos)
	if err != nil {
		t.Fatal(err)
	}

	var got Position
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got != pos {
		t.Errorf("got %+v, want %+v", got, pos)
	}
}

func TestDiagnosticMarshal(t *testing.T) {
	sev := SeverityError
	diag := Diagnostic{
		Range: Range{
			Start: Position{Line: 1, Character: 0},
			End:   Position{Line: 1, Character: 10},
		},
		Severity: &sev,
		Message:  "undefined: foo",
		Source:   "compiler",
	}

	data, err := json.Marshal(diag)
	if err != nil {
		t.Fatal(err)
	}

	var got Diagnostic
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Message != diag.Message {
		t.Errorf("message = %q, want %q", got.Message, diag.Message)
	}
	if got.Severity == nil || *got.Severity != SeverityError {
		t.Error("severity not preserved")
	}
}

func TestCompletionItemMarshal(t *testing.T) {
	kind := CompletionItemKindFunction
	item := CompletionItem{
		Label:  "myFunc",
		Kind:   &kind,
		Detail: "func myFunc()",
		Documentation: &MarkupContent{
			Kind:  Markdown,
			Value: "Does stuff",
		},
	}

	data, err := json.Marshal(item)
	if err != nil {
		t.Fatal(err)
	}

	var got CompletionItem
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.Label != "myFunc" {
		t.Errorf("label = %q, want myFunc", got.Label)
	}
	if got.Kind == nil || *got.Kind != CompletionItemKindFunction {
		t.Error("kind not preserved")
	}
}

func TestInitializeParamsMarshal(t *testing.T) {
	rootURI := DocumentURI("file:///workspace")
	pid := 1234
	params := InitializeParams{
		ProcessID: &pid,
		RootURI:   &rootURI,
		Capabilities: ClientCapabilities{
			TextDocument: &TextDocumentClientCapabilities{
				Hover: &HoverClientCapabilities{
					ContentFormat: []MarkupKind{Markdown, PlainText},
				},
			},
		},
	}

	data, err := json.Marshal(params)
	if err != nil {
		t.Fatal(err)
	}

	var got InitializeParams
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.ProcessID == nil || *got.ProcessID != 1234 {
		t.Error("processId not preserved")
	}
	if got.RootURI == nil || *got.RootURI != rootURI {
		t.Error("rootUri not preserved")
	}
	if got.Capabilities.TextDocument == nil || got.Capabilities.TextDocument.Hover == nil {
		t.Fatal("text document hover capabilities not preserved")
	}
	if len(got.Capabilities.TextDocument.Hover.ContentFormat) != 2 {
		t.Error("content format not preserved")
	}
}

func TestServerCapabilitiesMarshal(t *testing.T) {
	tr := true
	caps := ServerCapabilities{
		HoverProvider:     &tr,
		DefinitionProvider: &tr,
		TextDocumentSync: &TextDocumentSyncOptions{
			OpenClose: &tr,
			Change:    SyncIncremental,
		},
		CompletionProvider: &CompletionOptions{
			TriggerCharacters: []string{"."},
			ResolveProvider:   &tr,
		},
	}

	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatal(err)
	}

	var got ServerCapabilities
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got.HoverProvider == nil || !*got.HoverProvider {
		t.Error("hoverProvider not preserved")
	}
	if got.CompletionProvider == nil {
		t.Fatal("completionProvider not preserved")
	}
	if len(got.CompletionProvider.TriggerCharacters) != 1 || got.CompletionProvider.TriggerCharacters[0] != "." {
		t.Error("trigger characters not preserved")
	}
}
