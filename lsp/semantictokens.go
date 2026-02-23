package lsp

// SemanticTokensLegend describes the semantic token types and modifiers.
type SemanticTokensLegend struct {
	TokenTypes     []string `json:"tokenTypes"`
	TokenModifiers []string `json:"tokenModifiers"`
}

// SemanticTokensOptions describes options for semantic tokens.
type SemanticTokensOptions struct {
	WorkDoneProgressOptions
	Legend SemanticTokensLegend `json:"legend"`
	Range  *bool                `json:"range,omitempty"`
	Full   *SemanticTokensFull  `json:"full,omitempty"`
}

// SemanticTokensFull describes options for full semantic tokens.
type SemanticTokensFull struct {
	Delta *bool `json:"delta,omitempty"`
}

// SemanticTokensParams contains the params for textDocument/semanticTokens/full.
type SemanticTokensParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// SemanticTokensDeltaParams contains the params for textDocument/semanticTokens/full/delta.
type SemanticTokensDeltaParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument   TextDocumentIdentifier `json:"textDocument"`
	PreviousResultID string               `json:"previousResultId"`
}

// SemanticTokensRangeParams contains the params for textDocument/semanticTokens/range.
type SemanticTokensRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
}

// SemanticTokens represents the result of a semantic tokens request.
type SemanticTokens struct {
	ResultID string `json:"resultId,omitempty"`
	Data     []int  `json:"data"`
}

// SemanticTokensDelta represents a delta update for semantic tokens.
type SemanticTokensDelta struct {
	ResultID string               `json:"resultId,omitempty"`
	Edits    []SemanticTokensEdit `json:"edits"`
}

// SemanticTokensEdit represents a single edit in a semantic tokens delta.
type SemanticTokensEdit struct {
	Start       int   `json:"start"`
	DeleteCount int   `json:"deleteCount"`
	Data        []int `json:"data,omitempty"`
}
