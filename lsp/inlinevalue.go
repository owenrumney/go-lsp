package lsp

// InlineValueContext describes the context for inline values.
type InlineValueContext struct {
	FrameID         int   `json:"frameId"`
	StoppedLocation Range `json:"stoppedLocation"`
}

// InlineValueText displays a static text value inline at a source position during debugging.
type InlineValueText struct {
	Range Range  `json:"range"`
	Text  string `json:"text"`
}

// InlineValueVariableLookup tells the debugger to look up a variable's value and display it inline at a source position.
type InlineValueVariableLookup struct {
	Range               Range  `json:"range"`
	VariableName        string `json:"variableName,omitempty"`
	CaseSensitiveLookup bool   `json:"caseSensitiveLookup"`
}

// InlineValueEvaluatableExpression tells the debugger to evaluate an expression and display the result inline at a source position.
type InlineValueEvaluatableExpression struct {
	Range      Range  `json:"range"`
	Expression string `json:"expression,omitempty"`
}

// InlineValueParams is sent to request inline debug values for a document in a stopped debug session.
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

// InlineValueWorkspaceClientCapabilities declares whether the editor will refresh inline values when the server requests it.
type InlineValueWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
