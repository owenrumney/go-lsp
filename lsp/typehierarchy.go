package lsp

import "encoding/json"

// TypeHierarchyPrepareParams contains the params for textDocument/prepareTypeHierarchy.
type TypeHierarchyPrepareParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// TypeHierarchyItem represents an item in the type hierarchy.
type TypeHierarchyItem struct {
	Name           string          `json:"name"`
	Kind           SymbolKind      `json:"kind"`
	Tags           []SymbolTag     `json:"tags,omitempty"`
	Deprecated     *bool           `json:"deprecated,omitempty"`
	Detail         string          `json:"detail,omitempty"`
	URI            DocumentURI     `json:"uri"`
	Range          Range           `json:"range"`
	SelectionRange Range           `json:"selectionRange"`
	Data           json.RawMessage `json:"data,omitempty"`
}

// TypeHierarchySupertypesParams contains the params for typeHierarchy/supertypes.
type TypeHierarchySupertypesParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item TypeHierarchyItem `json:"item"`
}

// TypeHierarchySubtypesParams contains the params for typeHierarchy/subtypes.
type TypeHierarchySubtypesParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item TypeHierarchyItem `json:"item"`
}
