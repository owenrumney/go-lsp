package lsp

// UniquenessLevel is a string enum indicating the scope at which a moniker is unique: document, project, group, scheme, or global.
type UniquenessLevel string

const (
	UniquenessLevelDocument UniquenessLevel = "document"
	UniquenessLevelProject  UniquenessLevel = "project"
	UniquenessLevelGroup    UniquenessLevel = "group"
	UniquenessLevelScheme   UniquenessLevel = "scheme"
	UniquenessLevelGlobal   UniquenessLevel = "global"
)

// MonikerKind is a string enum: "import" (referencing external), "export" (visible to others), or "local" (internal).
type MonikerKind string

const (
	MonikerKindImport MonikerKind = "import"
	MonikerKindExport MonikerKind = "export"
	MonikerKindLocal  MonikerKind = "local"
)

// MonikerParams is sent to request stable, cross-project identifiers for a symbol (used for cross-repo navigation and indexing).
type MonikerParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// Moniker is a stable, cross-project identifier for a symbol, enabling features like cross-repository go-to-definition.
type Moniker struct {
	Scheme     string          `json:"scheme"`
	Identifier string          `json:"identifier"`
	Unique     UniquenessLevel `json:"unique"`
	Kind       *MonikerKind    `json:"kind,omitempty"`
}
