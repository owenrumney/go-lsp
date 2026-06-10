package lsp

// FoldingRangeKind is a string enum: "comment", "imports", or "region".
type FoldingRangeKind string

const (
	// Folding range for a comment.
	FoldingRangeKindComment FoldingRangeKind = "comment"
	// Folding range for an import or include.
	FoldingRangeKindImports FoldingRangeKind = "imports"
	// Folding range for a region (e.g. `#region`).
	FoldingRangeKindRegion FoldingRangeKind = "region"
)

// FoldingRangeParams holds the parameters for a [FoldingRangeRequest].
type FoldingRangeParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// FoldingRange represents a folding range. To be valid, start and end line must be bigger than zero and smaller
// than the number of lines in the document. Clients are free to ignore invalid ranges.
type FoldingRange struct {
	// The zero-based start line of the range to fold. The folded area starts after the line's last character.
	// To be valid, the end must be zero or larger and smaller than the number of lines in the document.
	StartLine int `json:"startLine"`
	// The zero-based character offset from where the folded range starts. If not defined, defaults to the length of the start line.
	StartCharacter *int `json:"startCharacter,omitempty"`
	// The zero-based end line of the range to fold. The folded area ends with the line's last character.
	// To be valid, the end must be zero or larger and smaller than the number of lines in the document.
	EndLine int `json:"endLine"`
	// The zero-based character offset before the folded range ends. If not defined, defaults to the length of the end line.
	EndCharacter *int `json:"endCharacter,omitempty"`
	// Describes the kind of the folding range such as `comment` or `region`. The kind
	// is used to categorize folding ranges and used by commands like 'Fold all comments'.
	// See [FoldingRangeKind] for an enumeration of standardized kinds.
	Kind *FoldingRangeKind `json:"kind,omitempty"`
}
