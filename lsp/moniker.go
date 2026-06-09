package lsp

// UniquenessLevel is a string enum indicating the scope at which a moniker is unique: document, project, group, scheme, or global.
type UniquenessLevel string

const (
	// The moniker is only unique inside a document.
	UniquenessLevelDocument UniquenessLevel = "document"
	// The moniker is unique inside a project for which a dump got created.
	UniquenessLevelProject UniquenessLevel = "project"
	// The moniker is unique inside the group to which a project belongs.
	UniquenessLevelGroup UniquenessLevel = "group"
	// The moniker is unique inside the moniker scheme.
	UniquenessLevelScheme UniquenessLevel = "scheme"
	// The moniker is globally unique.
	UniquenessLevelGlobal UniquenessLevel = "global"
)

// MonikerKind is a string enum: "import" (referencing external), "export" (visible to others), or "local" (internal).
type MonikerKind string

const (
	// The moniker represents a symbol that is imported into a project.
	MonikerKindImport MonikerKind = "import"
	// The moniker represents a symbol that is exported from a project.
	MonikerKindExport MonikerKind = "export"
	// The moniker represents a symbol that is local to a project (e.g. a local
	// variable of a function, a class not visible outside the project, ...)
	MonikerKindLocal MonikerKind = "local"
)

// MonikerParams is sent to request stable, cross-project identifiers for a symbol (used for cross-repo navigation and indexing).
type MonikerParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// Moniker is the LSP moniker definition, matching LSIF 0.5.
//
// Since 3.16.0.
type Moniker struct {
	// The scheme of the moniker. For example tsc or .Net
	Scheme string `json:"scheme"`
	// The identifier of the moniker. The value is opaque in LSIF however
	// schema owners are allowed to define the structure if they want.
	Identifier string `json:"identifier"`
	// The scope in which the moniker is unique
	Unique UniquenessLevel `json:"unique"`
	// The moniker kind if known.
	Kind *MonikerKind `json:"kind,omitempty"`
}
