package lsp

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestHoverContentsRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		in   HoverContents
		json string
	}{
		{
			"markup",
			NewHoverContents(Markdown, "**hi**"),
			`{"kind":"markdown","value":"**hi**"}`,
		},
		{
			"single bare marked string",
			NewHoverMarkedStrings(MarkedString{Value: "plain"}),
			`"plain"`,
		},
		{
			"single code marked string",
			NewHoverMarkedStrings(MarkedString{Language: "go", Value: "x := 1"}),
			`{"language":"go","value":"x := 1"}`,
		},
		{
			"multiple marked strings",
			NewHoverMarkedStrings(MarkedString{Value: "a"}, MarkedString{Language: "go", Value: "b"}),
			`["a",{"language":"go","value":"b"}]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tt.json {
				t.Fatalf("marshal = %s, want %s", data, tt.json)
			}

			var got HoverContents
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			back, err := json.Marshal(got)
			if err != nil {
				t.Fatal(err)
			}
			if string(back) != tt.json {
				t.Fatalf("round-trip = %s, want %s", back, tt.json)
			}
		})
	}
}

func TestHoverContentsValue(t *testing.T) {
	if v := NewHoverContents(Markdown, "m").Value(); v != "m" {
		t.Errorf("markup Value() = %q", v)
	}
	if v := NewHoverMarkedStrings(MarkedString{Value: "s"}).Value(); v != "s" {
		t.Errorf("marked string Value() = %q", v)
	}
}

func TestCompletionTextEditRoundTrip(t *testing.T) {
	plain := NewCompletionTextEdit(TextEdit{
		Range:   Range{Start: Position{0, 0}, End: Position{0, 3}},
		NewText: "abc",
	})
	data, err := json.Marshal(plain)
	if err != nil {
		t.Fatal(err)
	}
	var gotPlain CompletionTextEdit
	if err := json.Unmarshal(data, &gotPlain); err != nil {
		t.Fatal(err)
	}
	if gotPlain.TextEdit == nil || gotPlain.TextEdit.NewText != "abc" {
		t.Fatalf("got %+v, want TextEdit", gotPlain)
	}

	ire := NewCompletionInsertReplaceEdit(InsertReplaceEdit{
		NewText: "abc",
		Insert:  Range{Start: Position{0, 0}, End: Position{0, 1}},
		Replace: Range{Start: Position{0, 0}, End: Position{0, 3}},
	})
	data, err = json.Marshal(ire)
	if err != nil {
		t.Fatal(err)
	}
	var gotIRE CompletionTextEdit
	if err := json.Unmarshal(data, &gotIRE); err != nil {
		t.Fatal(err)
	}
	if gotIRE.InsertReplaceEdit == nil || gotIRE.TextEdit != nil {
		t.Fatalf("got %+v, want InsertReplaceEdit", gotIRE)
	}
	if gotIRE.InsertReplaceEdit.Replace.End.Character != 3 {
		t.Fatalf("replace range lost: %+v", gotIRE.InsertReplaceEdit)
	}
}

func TestProviderOptionsAcceptBareBool(t *testing.T) {
	var caps ServerCapabilities
	err := json.Unmarshal([]byte(`{"documentSymbolProvider":true,"workspaceSymbolProvider":true}`), &caps)
	if err != nil {
		t.Fatal(err)
	}
	if caps.DocumentSymbolProvider == nil {
		t.Error("documentSymbolProvider bool form not accepted")
	}
	if caps.WorkspaceSymbolProvider == nil {
		t.Error("workspaceSymbolProvider bool form not accepted")
	}

	err = json.Unmarshal([]byte(`{"workspaceSymbolProvider":{"resolveProvider":true}}`), &caps)
	if err != nil {
		t.Fatal(err)
	}
	if caps.WorkspaceSymbolProvider.ResolveProvider == nil || !*caps.WorkspaceSymbolProvider.ResolveProvider {
		t.Error("workspaceSymbolProvider options form lost resolveProvider")
	}
}

func TestUnionZeroValuesMarshal(t *testing.T) {
	tests := map[string]struct {
		v    any
		want string
	}{
		"hover":              {&Hover{}, `{"contents":{"kind":"","value":""}}`},
		"documentChange":     {DocumentChange{}, `{"textDocument":{"uri":"","version":null},"edits":null}`},
		"completionTextEdit": {&CompletionTextEdit{}, `{"range":{"start":{"line":0,"character":0},"end":{"line":0,"character":0}},"newText":""}`},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			data, err := json.Marshal(tt.v)
			if err != nil {
				t.Fatalf("zero value failed to marshal: %v", err)
			}
			if string(data) != tt.want {
				t.Fatalf("marshal = %s, want %s", data, tt.want)
			}
		})
	}
}

func TestProviderOptionsFalseRoundTrip(t *testing.T) {
	var caps ServerCapabilities
	if err := json.Unmarshal([]byte(`{"documentSymbolProvider":false,"workspaceSymbolProvider":false}`), &caps); err != nil {
		t.Fatal(err)
	}
	if caps.DocumentSymbolProvider.Enabled() {
		t.Error("documentSymbolProvider false must not report enabled")
	}
	if caps.WorkspaceSymbolProvider.Enabled() {
		t.Error("workspaceSymbolProvider false must not report enabled")
	}

	data, err := json.Marshal(caps)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"documentSymbolProvider":false`, `"workspaceSymbolProvider":false`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("round-trip lost %s: %s", want, data)
		}
	}

	if err := json.Unmarshal([]byte(`{"documentSymbolProvider":true}`), &caps); err != nil {
		t.Fatal(err)
	}
	if !caps.DocumentSymbolProvider.Enabled() {
		t.Error("documentSymbolProvider true must report enabled")
	}
	var nilOpts *DocumentSymbolOptions
	if nilOpts.Enabled() {
		t.Error("nil options must not report enabled")
	}
}
