package lsp

// FoldingRangeKind represents the kind of a folding range.
type FoldingRangeKind string

const (
	FoldingRangeKindComment FoldingRangeKind = "comment"
	FoldingRangeKindImports FoldingRangeKind = "imports"
	FoldingRangeKindRegion  FoldingRangeKind = "region"
)

// FoldingRangeParams contains the params for textDocument/foldingRange.
type FoldingRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// FoldingRange represents a folding range.
type FoldingRange struct {
	StartLine      int               `json:"startLine"`
	StartCharacter *int              `json:"startCharacter,omitempty"`
	EndLine        int               `json:"endLine"`
	EndCharacter   *int              `json:"endCharacter,omitempty"`
	Kind           *FoldingRangeKind `json:"kind,omitempty"`
}
