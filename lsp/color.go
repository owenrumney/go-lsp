package lsp

// Color represents a color in RGBA space.
type Color struct {
	Red   float64 `json:"red"`
	Green float64 `json:"green"`
	Blue  float64 `json:"blue"`
	Alpha float64 `json:"alpha"`
}

// ColorInformation pairs a document range with the color value found there, used by the editor to show color decorators.
type ColorInformation struct {
	Range Range `json:"range"`
	Color Color `json:"color"`
}

// DocumentColorParams is sent to request all color literals in a document for inline color decoration.
type DocumentColorParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// ColorPresentationParams is sent to request how a color value can be represented as text (e.g. "#ff0000", "rgb(255,0,0)").
type ColorPresentationParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Color        Color                  `json:"color"`
	Range        Range                  `json:"range"`
}

// ColorPresentation represents how a color should be presented.
type ColorPresentation struct {
	Label               string     `json:"label"`
	TextEdit            *TextEdit  `json:"textEdit,omitempty"`
	AdditionalTextEdits []TextEdit `json:"additionalTextEdits,omitempty"`
}
