package lsp

import "encoding/json"

// TypeHierarchyPrepareParams is sent to resolve the type hierarchy item at a given cursor position.
type TypeHierarchyPrepareParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// TypeHierarchyItem represents a type (class, interface, etc.) with its location and detail, for navigating supertypes/subtypes.
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

// TypeHierarchySupertypesParams is sent to find all supertypes (parent classes, implemented interfaces) of a type.
type TypeHierarchySupertypesParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item TypeHierarchyItem `json:"item"`
}

// TypeHierarchySubtypesParams is sent to find all subtypes (child classes, implementors) of a type.
type TypeHierarchySubtypesParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item TypeHierarchyItem `json:"item"`
}
