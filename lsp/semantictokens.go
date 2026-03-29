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

// SemanticTokensParams is sent to request all semantic tokens for a document (for syntax highlighting beyond grammar-based tokenization).
type SemanticTokensParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// SemanticTokensDeltaParams is sent to request only the semantic tokens that changed since a previous response.
type SemanticTokensDeltaParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument     TextDocumentIdentifier `json:"textDocument"`
	PreviousResultID string                 `json:"previousResultId"`
}

// SemanticTokensRangeParams is sent to request semantic tokens for a visible range (for large files where full tokenization is too slow).
type SemanticTokensRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
}

// SemanticTokens contains encoded token data (position, length, type, modifiers) as a flat integer array.
type SemanticTokens struct {
	ResultID string `json:"resultId,omitempty"`
	Data     []int  `json:"data"`
}

// SemanticTokensDelta contains a set of edits to apply to previously returned token data, avoiding full retransmission.
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
