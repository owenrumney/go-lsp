package lsp

// FormattingOptions is the value-object describing what options formatting should use.
type FormattingOptions struct {
	// Size of a tab in spaces.
	TabSize int `json:"tabSize"`
	// Prefer spaces over tabs.
	InsertSpaces bool `json:"insertSpaces"`
	// Trim trailing whitespace on a line.
	//
	// Since 3.15.0
	TrimTrailingWhitespace *bool `json:"trimTrailingWhitespace,omitempty"`
	// Insert a newline character at the end of the file if one does not exist.
	//
	// Since 3.15.0
	InsertFinalNewline *bool `json:"insertFinalNewline,omitempty"`
	// Trim all newlines after the final newline at the end of the file.
	//
	// Since 3.15.0
	TrimFinalNewlines *bool `json:"trimFinalNewlines,omitempty"`
}

// DocumentFormattingParams is the parameters of a [DocumentFormattingRequest].
type DocumentFormattingParams struct {
	WorkDoneProgressParams
	// The document to format.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The format options.
	Options FormattingOptions `json:"options"`
}

// DocumentRangeFormattingParams is the parameters of a [DocumentRangeFormattingRequest].
type DocumentRangeFormattingParams struct {
	WorkDoneProgressParams
	// The document to format.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The range to format
	Range Range `json:"range"`
	// The format options
	Options FormattingOptions `json:"options"`
}

// DocumentOnTypeFormattingParams is the parameters of a [DocumentOnTypeFormattingRequest].
type DocumentOnTypeFormattingParams struct {
	TextDocumentPositionParams
	// The character that has been typed that triggered the formatting
	// on type request. That is not necessarily the last character that
	// got inserted into the document since the client could auto insert
	// characters as well (e.g. like automatic brace completion).
	Character string `json:"ch"`
	// The formatting options.
	Options FormattingOptions `json:"options"`
}
