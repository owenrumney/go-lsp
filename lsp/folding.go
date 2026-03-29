package lsp

// FoldingRangeKind is a string enum: "comment", "imports", or "region".
type FoldingRangeKind string

const (
	FoldingRangeKindComment FoldingRangeKind = "comment"
	FoldingRangeKindImports FoldingRangeKind = "imports"
	FoldingRangeKindRegion  FoldingRangeKind = "region"
)

// FoldingRangeParams is sent to request all foldable regions in a document.
type FoldingRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// FoldingRange defines a collapsible region in a document (e.g. a function body, import block, or comment block).
type FoldingRange struct {
	StartLine      int               `json:"startLine"`
	StartCharacter *int              `json:"startCharacter,omitempty"`
	EndLine        int               `json:"endLine"`
	EndCharacter   *int              `json:"endCharacter,omitempty"`
	Kind           *FoldingRangeKind `json:"kind,omitempty"`
}
