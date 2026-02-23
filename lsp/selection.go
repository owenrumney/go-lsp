package lsp

// SelectionRangeParams contains the params for textDocument/selectionRange.
type SelectionRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Positions    []Position             `json:"positions"`
}

// SelectionRange represents a selection range.
type SelectionRange struct {
	Range  Range           `json:"range"`
	Parent *SelectionRange `json:"parent,omitempty"`
}
