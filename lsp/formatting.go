package lsp

// FormattingOptions describes options to be used during formatting.
type FormattingOptions struct {
	TabSize                int  `json:"tabSize"`
	InsertSpaces           bool `json:"insertSpaces"`
	TrimTrailingWhitespace *bool `json:"trimTrailingWhitespace,omitempty"`
	InsertFinalNewline     *bool `json:"insertFinalNewline,omitempty"`
	TrimFinalNewlines      *bool `json:"trimFinalNewlines,omitempty"`
}

// DocumentFormattingParams contains the params for textDocument/formatting.
type DocumentFormattingParams struct {
	WorkDoneProgressParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Options      FormattingOptions      `json:"options"`
}

// DocumentRangeFormattingParams contains the params for textDocument/rangeFormatting.
type DocumentRangeFormattingParams struct {
	WorkDoneProgressParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Options      FormattingOptions      `json:"options"`
}

// DocumentOnTypeFormattingParams contains the params for textDocument/onTypeFormatting.
type DocumentOnTypeFormattingParams struct {
	TextDocumentPositionParams
	Character string            `json:"ch"`
	Options   FormattingOptions `json:"options"`
}
