package document

import (
	"errors"
	"testing"

	"github.com/owenrumney/go-lsp/lsp"
)

func TestOffsetAtAndPositionAtUseUTF16(t *testing.T) {
	doc := newDocument(lsp.TextDocumentItem{
		URI:        "file:///test.txt",
		LanguageID: "txt",
		Version:    1,
		Text:       "a😀b\néx",
	})

	tests := []struct {
		name   string
		pos    lsp.Position
		offset int
	}{
		{name: "start", pos: lsp.Position{Line: 0, Character: 0}, offset: 0},
		{name: "after ascii", pos: lsp.Position{Line: 0, Character: 1}, offset: 1},
		{name: "after emoji", pos: lsp.Position{Line: 0, Character: 3}, offset: 5},
		{name: "after b", pos: lsp.Position{Line: 0, Character: 4}, offset: 6},
		{name: "second line accented", pos: lsp.Position{Line: 1, Character: 1}, offset: 9},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			offset, err := doc.OffsetAt(tt.pos)
			if err != nil {
				t.Fatalf("OffsetAt returned error: %v", err)
			}
			if offset != tt.offset {
				t.Fatalf("OffsetAt = %d, want %d", offset, tt.offset)
			}

			pos, err := doc.PositionAt(tt.offset)
			if err != nil {
				t.Fatalf("PositionAt returned error: %v", err)
			}
			if pos != tt.pos {
				t.Fatalf("PositionAt = %+v, want %+v", pos, tt.pos)
			}
		})
	}
}

func TestOffsetAtRejectsInsideSurrogatePair(t *testing.T) {
	doc := newDocument(lsp.TextDocumentItem{Text: "😀"})
	if _, err := doc.OffsetAt(lsp.Position{Line: 0, Character: 1}); !errors.Is(err, ErrInvalidPosition) {
		t.Fatalf("OffsetAt error = %v, want ErrInvalidPosition", err)
	}
}

func TestPositionAtRejectsInsideUTF8Sequence(t *testing.T) {
	doc := newDocument(lsp.TextDocumentItem{Text: "é"})
	if _, err := doc.PositionAt(1); !errors.Is(err, ErrInvalidPosition) {
		t.Fatalf("PositionAt error = %v, want ErrInvalidPosition", err)
	}
}

func TestApplyFullChange(t *testing.T) {
	doc := newDocument(lsp.TextDocumentItem{Version: 1, Text: "old"})
	err := doc.ApplyChange(lsp.TextDocumentContentChangeEvent{Text: "new\ntext"}, 2)
	if err != nil {
		t.Fatal(err)
	}
	if doc.Text() != "new\ntext" {
		t.Fatalf("text = %q", doc.Text())
	}
	if doc.Version() != 2 {
		t.Fatalf("version = %d", doc.Version())
	}
}

func TestApplyIncrementalInsertReplaceDelete(t *testing.T) {
	doc := newDocument(lsp.TextDocumentItem{Version: 1, Text: "hello world"})

	err := doc.ApplyChange(lsp.TextDocumentContentChangeEvent{
		Range: &lsp.Range{
			Start: lsp.Position{Line: 0, Character: 5},
			End:   lsp.Position{Line: 0, Character: 5},
		},
		Text: ",",
	}, 2)
	if err != nil {
		t.Fatal(err)
	}

	err = doc.ApplyChange(lsp.TextDocumentContentChangeEvent{
		Range: &lsp.Range{
			Start: lsp.Position{Line: 0, Character: 7},
			End:   lsp.Position{Line: 0, Character: 12},
		},
		Text: "gopher",
	}, 3)
	if err != nil {
		t.Fatal(err)
	}

	err = doc.ApplyChange(lsp.TextDocumentContentChangeEvent{
		Range: &lsp.Range{
			Start: lsp.Position{Line: 0, Character: 5},
			End:   lsp.Position{Line: 0, Character: 6},
		},
		Text: "",
	}, 4)
	if err != nil {
		t.Fatal(err)
	}

	if doc.Text() != "hello gopher" {
		t.Fatalf("text = %q", doc.Text())
	}
}

func TestApplyMultilineChangeWithCRLF(t *testing.T) {
	doc := newDocument(lsp.TextDocumentItem{Version: 1, Text: "one\r\ntwo\r\nthree"})

	err := doc.ApplyChange(lsp.TextDocumentContentChangeEvent{
		Range: &lsp.Range{
			Start: lsp.Position{Line: 0, Character: 3},
			End:   lsp.Position{Line: 2, Character: 5},
		},
		Text: "\n2",
	}, 2)
	if err != nil {
		t.Fatal(err)
	}

	if doc.Text() != "one\n2" {
		t.Fatalf("text = %q", doc.Text())
	}
}

func TestApplyInvalidRange(t *testing.T) {
	doc := newDocument(lsp.TextDocumentItem{Text: "abc"})
	err := doc.ApplyChange(lsp.TextDocumentContentChangeEvent{
		Range: &lsp.Range{
			Start: lsp.Position{Line: 0, Character: 2},
			End:   lsp.Position{Line: 0, Character: 1},
		},
		Text: "x",
	}, 1)
	if !errors.Is(err, ErrInvalidRange) {
		t.Fatalf("error = %v, want ErrInvalidRange", err)
	}
}

func TestApplyRejectsVersionRegression(t *testing.T) {
	doc := newDocument(lsp.TextDocumentItem{Version: 2, Text: "abc"})
	err := doc.ApplyChange(lsp.TextDocumentContentChangeEvent{Text: "x"}, 1)
	if !errors.Is(err, ErrVersionRegression) {
		t.Fatalf("error = %v, want ErrVersionRegression", err)
	}
}
