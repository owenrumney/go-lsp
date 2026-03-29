package lsp

// HoverParams is sent to request hover information (docs, type info) at a cursor position.
type HoverParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// Hover contains the content (typically documentation or type info) and optional range to display in a hover tooltip.
type Hover struct {
	Contents MarkupContent `json:"contents"`
	Range    *Range        `json:"range,omitempty"`
}
