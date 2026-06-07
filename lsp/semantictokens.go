package lsp

// SemanticTokensLegend lists the token types and modifiers a server uses when encoding semantic tokens.
//
// Since 3.16.0.
type SemanticTokensLegend struct {
	// The token types a server uses.
	TokenTypes []string `json:"tokenTypes"`
	// The token modifiers a server uses.
	TokenModifiers []string `json:"tokenModifiers"`
}

// SemanticTokensOptions configures the server's semantic-tokens provider.
//
// Since 3.16.0.
type SemanticTokensOptions struct {
	WorkDoneProgressOptions
	// The legend used by the server
	Legend SemanticTokensLegend `json:"legend"`
	// Server supports providing semantic tokens for a specific range
	// of a document.
	Range *bool `json:"range,omitempty"`
	// Server supports providing semantic tokens for a full document.
	Full *SemanticTokensFull `json:"full,omitempty"`
}

// SemanticTokensFull describes options for full semantic tokens.
type SemanticTokensFull struct {
	Delta *bool `json:"delta,omitempty"`
}

// SemanticTokensParams is the parameters of a textDocument/semanticTokens/full request.
//
// Since 3.16.0.
type SemanticTokensParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// SemanticTokensDeltaParams is the parameters of a textDocument/semanticTokens/full/delta request.
//
// Since 3.16.0.
type SemanticTokensDeltaParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The result id of a previous response. The result Id can either point to a full response
	// or a delta response depending on what was received last.
	PreviousResultID string `json:"previousResultId"`
}

// SemanticTokensRangeParams is the parameters of a textDocument/semanticTokens/range request.
//
// Since 3.16.0.
type SemanticTokensRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The range the semantic tokens are requested for.
	Range Range `json:"range"`
}

// SemanticTokens is the result of a semantic-tokens request, an encoded array of token positions, types, and modifiers.
//
// Since 3.16.0.
type SemanticTokens struct {
	// An optional result id. If provided and clients support delta updating
	// the client will include the result id in the next semantic token request.
	// A server can then instead of computing all semantic tokens again simply
	// send a delta.
	ResultID string `json:"resultId,omitempty"`
	// The actual tokens.
	Data []int `json:"data"`
}

// SemanticTokensDelta is the delta result of a semantic-tokens delta request, expressed as edits over a prior [SemanticTokens] result.
//
// Since 3.16.0.
type SemanticTokensDelta struct {
	ResultID string `json:"resultId,omitempty"`
	// The semantic token edits to transform a previous result into a new result.
	Edits []SemanticTokensEdit `json:"edits"`
}

// SemanticTokensEdit describes an in-place edit to a previously returned [SemanticTokens] data array.
//
// Since 3.16.0.
type SemanticTokensEdit struct {
	// The start offset of the edit.
	Start int `json:"start"`
	// The count of elements to remove.
	DeleteCount int `json:"deleteCount"`
	// The elements to insert.
	Data []int `json:"data,omitempty"`
}
