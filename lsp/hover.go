package lsp

// HoverParams holds the parameters for a [HoverRequest].
type HoverParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// Hover is the result of a hover request.
type Hover struct {
	// The hover's content
	Contents MarkupContent `json:"contents"`
	// An optional range inside the text document that is used to
	// visualize the hover, e.g. by changing the background color.
	Range *Range `json:"range,omitempty"`
}
