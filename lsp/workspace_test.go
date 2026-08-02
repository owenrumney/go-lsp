package lsp

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestDocumentChangeRoundTrip(t *testing.T) {
	version := 3
	edit := WorkspaceEdit{
		DocumentChanges: []DocumentChange{
			NewCreateFileChange(CreateFile{URI: "file:///new.go"}),
			NewTextDocumentEditChange(TextDocumentEdit{
				TextDocument: OptionalVersionedTextDocumentIdentifier{
					TextDocumentIdentifier: TextDocumentIdentifier{URI: "file:///new.go"},
					Version:                &version,
				},
				Edits: []TextEdit{{NewText: "package main\n"}},
			}),
			NewRenameFileChange(RenameFile{OldURI: "file:///old.go", NewURI: "file:///renamed.go"}),
			NewDeleteFileChange(DeleteFile{URI: "file:///gone.go"}),
		},
	}

	data, err := json.Marshal(edit)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`"kind":"create"`, `"kind":"rename"`, `"kind":"delete"`, `"newText":"package main\n"`} {
		if !strings.Contains(string(data), want) {
			t.Errorf("marshaled edit missing %s: %s", want, data)
		}
	}

	var got WorkspaceEdit
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if len(got.DocumentChanges) != 4 {
		t.Fatalf("got %d document changes, want 4", len(got.DocumentChanges))
	}
	if got.DocumentChanges[0].CreateFile == nil || got.DocumentChanges[0].CreateFile.URI != "file:///new.go" {
		t.Errorf("change 0 = %+v, want CreateFile", got.DocumentChanges[0])
	}
	if got.DocumentChanges[1].TextDocumentEdit == nil || len(got.DocumentChanges[1].TextDocumentEdit.Edits) != 1 {
		t.Errorf("change 1 = %+v, want TextDocumentEdit", got.DocumentChanges[1])
	}
	if got.DocumentChanges[2].RenameFile == nil || got.DocumentChanges[2].RenameFile.NewURI != "file:///renamed.go" {
		t.Errorf("change 2 = %+v, want RenameFile", got.DocumentChanges[2])
	}
	if got.DocumentChanges[3].DeleteFile == nil || got.DocumentChanges[3].DeleteFile.URI != "file:///gone.go" {
		t.Errorf("change 3 = %+v, want DeleteFile", got.DocumentChanges[3])
	}
}

func TestDocumentChangeUnknownKind(t *testing.T) {
	var c DocumentChange
	if err := json.Unmarshal([]byte(`{"kind":"explode"}`), &c); err == nil {
		t.Fatal("expected error for unknown kind")
	}
}
