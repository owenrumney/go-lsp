package lsp

// LinkedEditingRangeParams is sent to find ranges that should be edited simultaneously (e.g. matching HTML open/close tags).
type LinkedEditingRangeParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// LinkedEditingRanges contains the set of ranges that change together and an optional word pattern constraining valid edits.
type LinkedEditingRanges struct {
	Ranges      []Range `json:"ranges"`
	WordPattern string  `json:"wordPattern,omitempty"`
}
