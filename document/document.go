package document

import (
	"fmt"
	"slices"
	"strings"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/owenrumney/go-lsp/lsp"
)

// Document is an open text document.
type Document struct {
	uri        lsp.DocumentURI
	languageID string
	version    int
	text       string
	lineStarts []int
}

func newDocument(item lsp.TextDocumentItem) *Document {
	d := &Document{
		uri:        item.URI,
		languageID: item.LanguageID,
		version:    item.Version,
		text:       item.Text,
	}
	d.reindex()
	return d
}

// URI returns the document URI.
func (d *Document) URI() lsp.DocumentURI {
	return d.uri
}

// LanguageID returns the document's language identifier.
func (d *Document) LanguageID() string {
	return d.languageID
}

// Version returns the latest document version seen by the store.
func (d *Document) Version() int {
	return d.version
}

// Text returns the full document text.
func (d *Document) Text() string {
	return d.text
}

// Lines returns the document split on "\n". Line endings are preserved except
// for the delimiter removed by strings.Split.
func (d *Document) Lines() []string {
	return strings.Split(d.text, "\n")
}

// Line returns a single zero-based line.
func (d *Document) Line(n int) (string, bool) {
	if n < 0 || n >= len(d.lineStarts) {
		return "", false
	}
	start := d.lineStarts[n]
	end := len(d.text)
	if n+1 < len(d.lineStarts) {
		end = d.lineStarts[n+1] - 1
	}
	return d.text[start:end], true
}

// OffsetAt converts an LSP UTF-16 position to a byte offset in Text().
func (d *Document) OffsetAt(pos lsp.Position) (int, error) {
	return d.offsetAt(pos)
}

// PositionAt converts a byte offset in Text() to an LSP UTF-16 position.
func (d *Document) PositionAt(offset int) (lsp.Position, error) {
	if offset < 0 || offset > len(d.text) {
		return lsp.Position{}, fmt.Errorf("%w: offset %d out of bounds", ErrInvalidPosition, offset)
	}
	if !utf8.ValidString(d.text[:offset]) {
		return lsp.Position{}, fmt.Errorf("%w: offset %d splits a UTF-8 sequence", ErrInvalidPosition, offset)
	}

	line := max(0, len(d.lineStarts)-1)
	idx, ok := slices.BinarySearchFunc(d.lineStarts, offset, func(start, target int) int {
		switch {
		case start <= target:
			return -1
		default:
			return 1
		}
	})
	if ok {
		line = idx
	} else if idx > 0 {
		line = idx - 1
	}

	char := utf16Len(d.text[d.lineStarts[line]:offset])
	return lsp.Position{Line: line, Character: char}, nil
}

// ApplyChange applies a single LSP content change to the document.
func (d *Document) ApplyChange(change lsp.TextDocumentContentChangeEvent, version int) error {
	if version < d.version {
		return fmt.Errorf("%w: current=%d new=%d", ErrVersionRegression, d.version, version)
	}

	if change.Range == nil {
		d.text = change.Text
		d.version = version
		d.reindex()
		return nil
	}

	start, end, err := d.offsetRange(*change.Range)
	if err != nil {
		return err
	}

	d.text = d.text[:start] + change.Text + d.text[end:]
	d.version = version
	d.reindex()
	return nil
}

func (d *Document) offsetRange(r lsp.Range) (int, int, error) {
	start, err := d.offsetAt(r.Start)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: start: %v", ErrInvalidRange, err)
	}
	end, err := d.offsetAt(r.End)
	if err != nil {
		return 0, 0, fmt.Errorf("%w: end: %v", ErrInvalidRange, err)
	}
	if start > end {
		return 0, 0, fmt.Errorf("%w: start after end", ErrInvalidRange)
	}
	return start, end, nil
}

func (d *Document) offsetAt(pos lsp.Position) (int, error) {
	if pos.Line < 0 || pos.Line >= len(d.lineStarts) {
		return 0, fmt.Errorf("%w: line %d out of bounds", ErrInvalidPosition, pos.Line)
	}
	if pos.Character < 0 {
		return 0, fmt.Errorf("%w: character %d out of bounds", ErrInvalidPosition, pos.Character)
	}

	start := d.lineStarts[pos.Line]
	end := len(d.text)
	if pos.Line+1 < len(d.lineStarts) {
		end = d.lineStarts[pos.Line+1] - 1
	}

	offset, ok := byteOffsetForUTF16Character(d.text[start:end], pos.Character)
	if !ok {
		return 0, fmt.Errorf("%w: character %d out of bounds", ErrInvalidPosition, pos.Character)
	}
	return start + offset, nil
}

func (d *Document) reindex() {
	d.lineStarts = []int{0}
	for i, b := range []byte(d.text) {
		if b == '\n' {
			d.lineStarts = append(d.lineStarts, i+1)
		}
	}
}

func byteOffsetForUTF16Character(s string, character int) (int, bool) {
	if character == 0 {
		return 0, true
	}

	units := 0
	for offset, r := range s {
		if units == character {
			return offset, true
		}
		units += utf16RuneLen(r)
		if units > character {
			return 0, false
		}
	}
	if units == character {
		return len(s), true
	}
	return 0, false
}

func utf16Len(s string) int {
	n := 0
	for _, r := range s {
		n += utf16RuneLen(r)
	}
	return n
}

func utf16RuneLen(r rune) int {
	if r1, r2 := utf16.EncodeRune(r); r1 != '\uFFFD' || r2 != '\uFFFD' {
		return 2
	}
	return 1
}
