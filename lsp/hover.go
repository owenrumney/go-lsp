package lsp

// HoverParams contains the params for textDocument/hover.
type HoverParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// Hover represents the result of a hover request.
type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}
