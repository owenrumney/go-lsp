package lsp

// UniquenessLevel represents the uniqueness level of a moniker.
type UniquenessLevel string

const (
	UniquenessLevelDocument UniquenessLevel = "document"
	UniquenessLevelProject  UniquenessLevel = "project"
	UniquenessLevelGroup    UniquenessLevel = "group"
	UniquenessLevelScheme   UniquenessLevel = "scheme"
	UniquenessLevelGlobal   UniquenessLevel = "global"
)

// MonikerKind represents the kind of a moniker.
type MonikerKind string

const (
	MonikerKindImport MonikerKind = "import"
	MonikerKindExport MonikerKind = "export"
	MonikerKindLocal  MonikerKind = "local"
)

// MonikerParams contains the params for textDocument/moniker.
type MonikerParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// Moniker represents a moniker.
type Moniker struct {
	Scheme     string          `json:"scheme"`
	Identifier string          `json:"identifier"`
	Unique     UniquenessLevel `json:"unique"`
	Kind       *MonikerKind    `json:"kind,omitempty"`
}
