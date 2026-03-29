package lsp

// FormattingOptions describes options to be used during formatting.
type FormattingOptions struct {
	TabSize                int   `json:"tabSize"`
	InsertSpaces           bool  `json:"insertSpaces"`
	TrimTrailingWhitespace *bool `json:"trimTrailingWhitespace,omitempty"`
	InsertFinalNewline     *bool `json:"insertFinalNewline,omitempty"`
	TrimFinalNewlines      *bool `json:"trimFinalNewlines,omitempty"`
}

// DocumentFormattingParams is sent to request formatting of an entire document.
type DocumentFormattingParams struct {
	WorkDoneProgressParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

// DocumentRangeFormattingParams is sent to request formatting of a selected range within a document.
type DocumentRangeFormattingParams struct {
	WorkDoneProgressParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Options      FormattingOptions      `json:"options"`
}

// DocumentOnTypeFormattingParams is sent to request formatting triggered by typing a character (e.g. closing brace or semicolon).
type DocumentOnTypeFormattingParams struct {
	TextDocumentPositionParams
	Character string            `json:"ch"`
	Options   FormattingOptions `json:"options"`
}
