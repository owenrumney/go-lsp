package lsp

// Color represents a color in RGBA space.
type Color struct {
	Red   float64 `json:"red"`
	Green float64 `json:"green"`
	Blue  float64 `json:"blue"`
	Alpha float64 `json:"alpha"`
}

// ColorInformation represents a color range from a document.
type ColorInformation struct {
	Range Range `json:"range"`
	Color Color `json:"color"`
}

// DocumentColorParams contains the params for textDocument/documentColor.
type DocumentColorParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// ColorPresentationParams contains the params for textDocument/colorPresentation.
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
