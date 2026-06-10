package lsp

// Color represents a color in RGBA space.
type Color struct {
	// The red component of this color in the range [0-1].
	Red float64 `json:"red"`
	// The green component of this color in the range [0-1].
	Green float64 `json:"green"`
	// The blue component of this color in the range [0-1].
	Blue float64 `json:"blue"`
	// The alpha component of this color in the range [0-1].
	Alpha float64 `json:"alpha"`
}

// ColorInformation represents a color range from a document.
type ColorInformation struct {
	// The range in the document where this color appears.
	Range Range `json:"range"`
	// The actual color value for this color range.
	Color Color `json:"color"`
}

// DocumentColorParams holds the parameters for a [DocumentColorRequest].
type DocumentColorParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// ColorPresentationParams holds the parameters for a [ColorPresentationRequest].
type ColorPresentationParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The color to request presentations for.
	Color Color `json:"color"`
	// The range where the color would be inserted. Serves as a context.
	Range Range `json:"range"`
}

// ColorPresentation represents how a color should be presented.
type ColorPresentation struct {
	// The label of this color presentation. It will be shown on the color
	// picker header. By default this is also the text that is inserted when selecting
	// this color presentation.
	Label string `json:"label"`
	// An [TextEdit] which is applied to a document when selecting
	// this presentation for the color.  When nil the [ColorPresentation.Label]
	// is used.
	TextEdit *TextEdit `json:"textEdit,omitempty"`
	// An optional array of additional [TextEdit] that are applied when
	// selecting this color presentation. Edits must not overlap with the main [ColorPresentation.TextEdit] nor with themselves.
	AdditionalTextEdits []TextEdit `json:"additionalTextEdits,omitempty"`
}
