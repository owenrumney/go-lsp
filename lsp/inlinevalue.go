package lsp

// InlineValueContext describes the context for inline values.
type InlineValueContext struct {
	FrameID         int   `json:"frameId"`
	StoppedLocation Range `json:"stoppedLocation"`
}

// InlineValueText represents an inline value as text.
type InlineValueText struct {
	Range Range  `json:"range"`
	Text  string `json:"text"`
}

// InlineValueVariableLookup represents an inline value via variable lookup.
type InlineValueVariableLookup struct {
	Range               Range  `json:"range"`
	VariableName        string `json:"variableName,omitempty"`
	CaseSensitiveLookup bool   `json:"caseSensitiveLookup"`
}

// InlineValueEvaluatableExpression represents an inline value via an evaluatable expression.
type InlineValueEvaluatableExpression struct {
	Range      Range  `json:"range"`
	Expression string `json:"expression,omitempty"`
}

// InlineValueParams contains the params for textDocument/inlineValue.
type InlineValueParams struct {
	WorkDoneProgressParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      InlineValueContext     `json:"context"`
}

// InlineValueOptions describes inline value options.
type InlineValueOptions struct {
	WorkDoneProgressOptions
}

// InlineValueWorkspaceClientCapabilities defines workspace capabilities for inline values.
type InlineValueWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
