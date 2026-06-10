package lsp

import "encoding/json"

// TypeHierarchyPrepareParams holds the parameters of a `textDocument/prepareTypeHierarchy` request.
//
// Since 3.17.0.
type TypeHierarchyPrepareParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// TypeHierarchyItem represents a node in a type-hierarchy graph.
//
// Since 3.17.0.
type TypeHierarchyItem struct {
	// The name of this item.
	Name string `json:"name"`
	// The kind of this item.
	Kind SymbolKind `json:"kind"`
	// Tags for this item.
	Tags       []SymbolTag `json:"tags,omitempty"`
	Deprecated *bool       `json:"deprecated,omitempty"`
	// More detail for this item, e.g. the signature of a function.
	Detail string `json:"detail,omitempty"`
	// The resource identifier of this item.
	URI DocumentURI `json:"uri"`
	// The range enclosing this symbol not including leading/trailing whitespace
	// but everything else, e.g. comments and code.
	Range Range `json:"range"`
	// The range that should be selected and revealed when this symbol is being
	// picked, e.g. the name of a function. Must be contained by the
	// [TypeHierarchyItem.Range].
	SelectionRange Range `json:"selectionRange"`
	// A data entry field that is preserved between a type hierarchy prepare and
	// supertypes or subtypes requests. It could also be used to identify the
	// type hierarchy in the server, helping improve the performance on
	// resolving supertypes and subtypes.
	Data json.RawMessage `json:"data,omitempty"`
}

// TypeHierarchySupertypesParams holds the parameters of a `typeHierarchy/supertypes` request.
//
// Since 3.17.0.
type TypeHierarchySupertypesParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item TypeHierarchyItem `json:"item"`
}

// TypeHierarchySubtypesParams holds the parameters of a `typeHierarchy/subtypes` request.
//
// Since 3.17.0.
type TypeHierarchySubtypesParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item TypeHierarchyItem `json:"item"`
}
