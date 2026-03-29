package lsp

// SelectionRangeParams is sent to request nested selection ranges (expand/shrink selection) at given positions.
type SelectionRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Positions    []Position             `json:"positions"`
}

// SelectionRange is a nested range tree enabling incremental expand/shrink selection (e.g. expression -> statement -> block -> function).
type SelectionRange struct {
	Range  Range           `json:"range"`
	Parent *SelectionRange `json:"parent,omitempty"`
}
