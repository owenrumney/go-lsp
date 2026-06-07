package lsp

// InlineValueContext is the additional context provided with an inline-value request, including the stack frame ID and stopped location.
//
// Since 3.17.0.
type InlineValueContext struct {
	// The stack frame (as a DAP Id) where the execution has stopped.
	FrameID int `json:"frameId"`
	// The document range where execution has stopped.
	// Typically the end position of the range denotes the line where the inline values are shown.
	StoppedLocation Range `json:"stoppedLocation"`
}

// InlineValueText provide inline value as text.
//
// Since 3.17.0.
type InlineValueText struct {
	// The document range for which the inline value applies.
	Range Range `json:"range"`
	// The text of the inline value.
	Text string `json:"text"`
}

// InlineValueVariableLookup provide inline value through a variable lookup.
// If only a range is specified, the variable name will be extracted from the underlying document.
// An optional variable name can be used to override the extracted name.
//
// Since 3.17.0.
type InlineValueVariableLookup struct {
	// The document range for which the inline value applies.
	// The range is used to extract the variable name from the underlying document.
	Range Range `json:"range"`
	// If specified the name of the variable to look up.
	VariableName string `json:"variableName,omitempty"`
	// How to perform the lookup.
	CaseSensitiveLookup bool `json:"caseSensitiveLookup"`
}

// InlineValueEvaluatableExpression provide an inline value through an expression evaluation.
// If only a range is specified, the expression will be extracted from the underlying document.
// An optional expression can be used to override the extracted expression.
//
// Since 3.17.0.
type InlineValueEvaluatableExpression struct {
	// The document range for which the inline value applies.
	// The range is used to extract the evaluatable expression from the underlying document.
	Range Range `json:"range"`
	// If specified the expression overrides the extracted expression.
	Expression string `json:"expression,omitempty"`
}

// InlineValueParams is a parameter literal used in inline value requests.
//
// Since 3.17.0.
type InlineValueParams struct {
	WorkDoneProgressParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The document range for which inline values should be computed.
	Range Range `json:"range"`
	// Additional information about the context in which inline values were
	// requested.
	Context InlineValueContext `json:"context"`
}

// InlineValueOptions is used during static registration.
//
// Since 3.17.0.
type InlineValueOptions struct {
	WorkDoneProgressOptions
}

// InlineValueWorkspaceClientCapabilities is specific to inline values.
//
// Since 3.17.0.
type InlineValueWorkspaceClientCapabilities struct {
	// Whether the client implementation supports a refresh request sent from the
	// server to the client.
	//
	// Note that this event is global and will force the client to refresh all
	// inline values currently shown. It should be used with absolute care and is
	// useful for situation where a server for example detects a project wide
	// change that requires such a calculation.
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
