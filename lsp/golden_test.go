package lsp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestGoldenJSONContracts(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		target   any
		assert   func(*testing.T, any)
	}{
		{
			name:     "initialize result",
			filename: "initialize_result.json",
			target:   &InitializeResult{},
			assert: func(t *testing.T, v any) {
				got := v.(*InitializeResult)
				if got.ServerInfo == nil || got.ServerInfo.Name != "go-lsp-test" {
					t.Fatalf("serverInfo = %+v", got.ServerInfo)
				}
				if got.Capabilities.PositionEncoding == nil || *got.Capabilities.PositionEncoding != PositionEncodingUTF16 {
					t.Fatalf("positionEncoding = %v", got.Capabilities.PositionEncoding)
				}
				if got.Capabilities.CompletionProvider == nil || len(got.Capabilities.CompletionProvider.TriggerCharacters) != 1 {
					t.Fatalf("completionProvider = %+v", got.Capabilities.CompletionProvider)
				}
			},
		},
		{
			name:     "publish diagnostics",
			filename: "publish_diagnostics.json",
			target:   &PublishDiagnosticsParams{},
			assert: func(t *testing.T, v any) {
				got := v.(*PublishDiagnosticsParams)
				if got.URI != "file:///workspace/main.go" || got.Version == nil || *got.Version != 7 {
					t.Fatalf("params = %+v", got)
				}
				if len(got.Diagnostics) != 1 || got.Diagnostics[0].Code == nil || got.Diagnostics[0].Message != "unused value" {
					t.Fatalf("diagnostics = %+v", got.Diagnostics)
				}
			},
		},
		{
			name:     "completion list",
			filename: "completion_list.json",
			target:   &CompletionList{},
			assert: func(t *testing.T, v any) {
				got := v.(*CompletionList)
				if got.IsIncomplete || len(got.Items) != 1 {
					t.Fatalf("completion list = %+v", got)
				}
				if got.Items[0].Kind == nil || *got.Items[0].Kind != CompletionItemKindFunction {
					t.Fatalf("completion item kind = %v", got.Items[0].Kind)
				}
			},
		},
		{
			name:     "code action",
			filename: "code_action.json",
			target:   &CodeAction{},
			assert: func(t *testing.T, v any) {
				got := v.(*CodeAction)
				if got.Kind == nil || *got.Kind != CodeActionQuickFix {
					t.Fatalf("kind = %v", got.Kind)
				}
				if got.Edit == nil || len(got.Edit.Changes["file:///workspace/main.go"]) != 1 {
					t.Fatalf("edit = %+v", got.Edit)
				}
				if got.Command == nil || got.Command.Command != "go.test" {
					t.Fatalf("command = %+v", got.Command)
				}
			},
		},
		{
			name:     "workspace edit",
			filename: "workspace_edit.json",
			target:   &WorkspaceEdit{},
			assert: func(t *testing.T, v any) {
				got := v.(*WorkspaceEdit)
				if len(got.Changes) != 1 || len(got.DocumentChanges) != 1 {
					t.Fatalf("workspace edit = %+v", got)
				}
				if got.DocumentChanges[0].TextDocument.Version == nil || *got.DocumentChanges[0].TextDocument.Version != 4 {
					t.Fatalf("document version = %+v", got.DocumentChanges[0].TextDocument.Version)
				}
			},
		},
		{
			name:     "semantic tokens",
			filename: "semantic_tokens.json",
			target:   &SemanticTokens{},
			assert: func(t *testing.T, v any) {
				got := v.(*SemanticTokens)
				if got.ResultID != "tokens-2" || len(got.Data) != 10 {
					t.Fatalf("semantic tokens = %+v", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := readGolden(t, tt.filename)
			if err := json.Unmarshal(data, tt.target); err != nil {
				t.Fatalf("unmarshal %s: %v", tt.filename, err)
			}
			tt.assert(t, tt.target)

			roundTrip, err := json.Marshal(tt.target)
			if err != nil {
				t.Fatalf("marshal %s: %v", tt.filename, err)
			}
			assertJSONSemanticallyEqual(t, data, roundTrip)
		})
	}
}

func readGolden(t *testing.T, filename string) []byte {
	t.Helper()
	data, err := os.ReadFile(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func assertJSONSemanticallyEqual(t *testing.T, want, got []byte) {
	t.Helper()
	var wantAny any
	if err := json.Unmarshal(want, &wantAny); err != nil {
		t.Fatalf("unmarshal want JSON: %v", err)
	}
	var gotAny any
	if err := json.Unmarshal(got, &gotAny); err != nil {
		t.Fatalf("unmarshal got JSON: %v", err)
	}
	if !reflect.DeepEqual(wantAny, gotAny) {
		t.Fatalf("JSON mismatch\ngot:  %s\nwant: %s", got, want)
	}
}
