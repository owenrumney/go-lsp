package document

import (
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/owenrumney/go-lsp/lsp"
)

func FuzzPositionOffsetRoundTrip(f *testing.F) {
	for _, seed := range []string{
		"",
		"hello",
		"a😀b\néx",
		"one\r\ntwo\r\nthree",
		"combining e\u0301\nsecond",
		"trailing\n",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		if !utf8.ValidString(text) {
			t.Skip()
		}

		doc := newDocument(lsp.TextDocumentItem{Text: text})
		for _, vp := range validDocumentPositions(text) {
			offset, err := doc.OffsetAt(vp.pos)
			if err != nil {
				t.Fatalf("OffsetAt(%+v) returned error for %q: %v", vp.pos, text, err)
			}
			if offset != vp.offset {
				t.Fatalf("OffsetAt(%+v) = %d, want %d for %q", vp.pos, offset, vp.offset, text)
			}

			pos, err := doc.PositionAt(vp.offset)
			if err != nil {
				t.Fatalf("PositionAt(%d) returned error for %q: %v", vp.offset, text, err)
			}
			if pos != vp.pos {
				t.Fatalf("PositionAt(%d) = %+v, want %+v for %q", vp.offset, pos, vp.pos, text)
			}
		}
	})
}

func FuzzInvalidPositions(f *testing.F) {
	for _, seed := range []string{
		"",
		"abc",
		"😀",
		"a😀b\néx",
		"one\r\ntwo",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, text string) {
		if !utf8.ValidString(text) {
			t.Skip()
		}

		doc := newDocument(lsp.TextDocumentItem{Text: text})
		invalid := []lsp.Position{
			{Line: -1, Character: 0},
			{Line: 0, Character: -1},
			{Line: len(strings.Split(text, "\n")), Character: 0},
		}

		for _, pos := range invalid {
			if _, err := doc.OffsetAt(pos); !errors.Is(err, ErrInvalidPosition) {
				t.Fatalf("OffsetAt(%+v) error = %v, want ErrInvalidPosition for %q", pos, err, text)
			}
		}

		lineStart := 0
		lines := strings.Split(text, "\n")
		for line, lineText := range lines {
			pos := lsp.Position{Line: line, Character: utf16Len(lineText) + 1}
			offset, err := doc.OffsetAt(pos)
			if err != nil {
				t.Fatalf("OffsetAt(%+v) returned error for %q: %v", pos, text, err)
			}
			// A trailing \r is stripped only when it terminates a \n line;
			// on the last line it is content.
			visible := lineText
			if line < len(lines)-1 {
				visible = strings.TrimSuffix(lineText, "\r")
			}
			want := lineStart + len(visible)
			if offset != want {
				t.Fatalf("OffsetAt(%+v) = %d, want clamped %d for %q", pos, offset, want, text)
			}

			for character := range invalidUTF16Characters(lineText) {
				pos := lsp.Position{Line: line, Character: character}
				offset, err := doc.OffsetAt(pos)
				if err != nil {
					t.Fatalf("OffsetAt(%+v) returned error for %q: %v", pos, text, err)
				}
				if !utf8.ValidString(text[:offset]) {
					t.Fatalf("OffsetAt(%+v) = %d splits a UTF-8 sequence in %q", pos, offset, text)
				}
			}

			lineStart += len(lineText) + 1
		}
	})
}

func FuzzApplyIncrementalEdit(f *testing.F) {
	seeds := []struct {
		text        string
		replacement string
		startIndex  int
		endIndex    int
	}{
		{"hello world", "gopher", 6, 11},
		{"a😀b\néx", "z", 1, 2},
		{"one\r\ntwo", "\n2", 1, 4},
		{"", "inserted", 0, 0},
	}
	for _, seed := range seeds {
		f.Add(seed.text, seed.replacement, seed.startIndex, seed.endIndex)
	}

	f.Fuzz(func(t *testing.T, text, replacement string, startIndex, endIndex int) {
		if !utf8.ValidString(text) || !utf8.ValidString(replacement) {
			t.Skip()
		}

		positions := validDocumentPositions(text)
		if len(positions) == 0 {
			t.Skip()
		}

		startIndex = positiveMod(startIndex, len(positions))
		endIndex = positiveMod(endIndex, len(positions))
		if startIndex > endIndex {
			startIndex, endIndex = endIndex, startIndex
		}

		start := positions[startIndex]
		end := positions[endIndex]
		doc := newDocument(lsp.TextDocumentItem{Version: 1, Text: text})

		err := doc.ApplyChange(lsp.TextDocumentContentChangeEvent{
			Range: &lsp.Range{
				Start: start.pos,
				End:   end.pos,
			},
			Text: replacement,
		}, 2)
		if err != nil {
			t.Fatalf("ApplyChange(%+v-%+v) returned error for %q: %v", start.pos, end.pos, text, err)
		}

		want := text[:start.offset] + replacement + text[end.offset:]
		if doc.Text() != want {
			t.Fatalf("text = %q, want %q", doc.Text(), want)
		}
		if !utf8.ValidString(doc.Text()) {
			t.Fatalf("document text is invalid UTF-8: %q", doc.Text())
		}
		if doc.Version() != 2 {
			t.Fatalf("version = %d, want 2", doc.Version())
		}
	})
}

type validPosition struct {
	pos    lsp.Position
	offset int
}

func validDocumentPositions(text string) []validPosition {
	positions := []validPosition{{pos: lsp.Position{}, offset: 0}}
	line := 0
	character := 0

	for offset, r := range text {
		// The position between a terminating \r and its \n is not distinct:
		// both map to the visible line end.
		betweenCRLF := r == '\n' && offset > 0 && text[offset-1] == '\r'
		if offset != 0 && !betweenCRLF {
			positions = append(positions, validPosition{
				pos:    lsp.Position{Line: line, Character: character},
				offset: offset,
			})
		}

		if r == '\n' {
			line++
			character = 0
			continue
		}
		if r == '\r' && offset+1 < len(text) && text[offset+1] == '\n' {
			continue
		}
		character += utf16RuneLen(r)
	}

	positions = append(positions, validPosition{
		pos:    lsp.Position{Line: line, Character: character},
		offset: len(text),
	})
	return positions
}

func invalidUTF16Characters(text string) map[int]struct{} {
	invalid := make(map[int]struct{})
	character := 0
	for _, r := range text {
		units := utf16RuneLen(r)
		if units > 1 {
			for i := 1; i < units; i++ {
				invalid[character+i] = struct{}{}
			}
		}
		character += units
	}
	return invalid
}

func positiveMod(n, mod int) int {
	n %= mod
	if n < 0 {
		n += mod
	}
	return n
}

func TestInvalidUTF16CharactersDetectsSurrogatePairInterior(t *testing.T) {
	got := invalidUTF16Characters("a😀b")
	if _, ok := got[2]; !ok {
		t.Fatalf("invalid UTF-16 characters = %v, want character 2", got)
	}
	if _, ok := got[1]; ok {
		t.Fatalf("invalid UTF-16 characters = %v, did not want character 1", got)
	}
}
