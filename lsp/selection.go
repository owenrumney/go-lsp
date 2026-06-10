package lsp

// SelectionRangeParams is a parameter literal used in selection range requests.
type SelectionRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The positions inside the text document.
	Positions []Position `json:"positions"`
}

// SelectionRange represents a part of a selection hierarchy. A selection range
// may have a parent selection range that contains it.
type SelectionRange struct {
	// The [Range] of this selection range.
	Range Range `json:"range"`
	// The parent selection range containing this range. Therefore `parent.range` must contain `this.range`.
	Parent *SelectionRange `json:"parent,omitempty"`
}
