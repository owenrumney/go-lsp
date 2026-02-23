package lsp

// LinkedEditingRangeParams contains the params for textDocument/linkedEditingRange.
type LinkedEditingRangeParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// LinkedEditingRanges represents the result of a linked editing range request.
type LinkedEditingRanges struct {
	Ranges      []Range `json:"ranges"`
	WordPattern string  `json:"wordPattern,omitempty"`
}
