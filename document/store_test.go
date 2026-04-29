package document

import (
	"errors"
	"testing"

	"github.com/owenrumney/go-lsp/lsp"
)

func TestStoreLifecycle(t *testing.T) {
	store := NewStore()
	uri := lsp.DocumentURI("file:///test.txt")

	doc, err := store.Open(&lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: "txt",
			Version:    1,
			Text:       "hello",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if doc.URI() != uri || doc.LanguageID() != "txt" || doc.Text() != "hello" {
		t.Fatalf("unexpected doc: %+v", doc)
	}

	doc, err = store.Change(&lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{URI: uri},
			Version:                2,
		},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{{Text: "updated"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if doc.Text() != "updated" {
		t.Fatalf("text = %q", doc.Text())
	}

	text, ok := store.Text(uri)
	if !ok || text != "updated" {
		t.Fatalf("Text = %q, %v", text, ok)
	}

	version, ok := store.Version(uri)
	if !ok || version != 2 {
		t.Fatalf("Version = %d, %v", version, ok)
	}

	uris := store.URIs()
	if len(uris) != 1 || uris[0] != uri {
		t.Fatalf("URIs = %v", uris)
	}

	store.Close(&lsp.DidCloseTextDocumentParams{TextDocument: lsp.TextDocumentIdentifier{URI: uri}})
	if _, ok := store.Get(uri); ok {
		t.Fatal("expected document to be closed")
	}
}

func TestStoreChangeMissingDocument(t *testing.T) {
	store := NewStore()
	_, err := store.Change(&lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{URI: "file:///missing.txt"},
			Version:                2,
		},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{{Text: "updated"}},
	})
	if !errors.Is(err, ErrDocumentNotFound) {
		t.Fatalf("error = %v, want ErrDocumentNotFound", err)
	}
}

func TestStoreGetReturnsSnapshot(t *testing.T) {
	store := NewStore()
	uri := lsp.DocumentURI("file:///test.txt")
	doc, err := store.Open(&lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{URI: uri, Version: 1, Text: "hello"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := doc.ApplyChange(lsp.TextDocumentContentChangeEvent{Text: "mutated"}, 2); err != nil {
		t.Fatal(err)
	}

	text, _ := store.Text(uri)
	if text != "hello" {
		t.Fatalf("stored text changed through snapshot: %q", text)
	}
}
