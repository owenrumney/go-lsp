package lsp

// LinkedEditingRangeParams is sent to find ranges that should be edited simultaneously (e.g. matching HTML open/close tags).
type LinkedEditingRangeParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// LinkedEditingRanges is the result of a linked editing range request.
//
// Since 3.16.0.
type LinkedEditingRanges struct {
	// A list of ranges that can be edited together. The ranges must have
	// identical length and contain identical text content. The ranges cannot overlap.
	Ranges []Range `json:"ranges"`
	// An optional word pattern (regular expression) that describes valid contents for
	// the given ranges. If no pattern is provided, the client configuration's word
	// pattern will be used.
	WordPattern string `json:"wordPattern,omitempty"`
}
