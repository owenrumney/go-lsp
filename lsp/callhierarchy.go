package lsp

import "encoding/json"

// CallHierarchyPrepareParams is the parameter of a `textDocument/prepareCallHierarchy` request.
//
// Since 3.16.0.
type CallHierarchyPrepareParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// CallHierarchyItem represents programming constructs like functions or constructors in the context
// of call hierarchy.
//
// Since 3.16.0.
type CallHierarchyItem struct {
	// The name of this item.
	Name string `json:"name"`
	// The kind of this item.
	Kind SymbolKind `json:"kind"`
	// Tags for this item.
	Tags []SymbolTag `json:"tags,omitempty"`
	// More detail for this item, e.g. the signature of a function.
	Detail string `json:"detail,omitempty"`
	// The resource identifier of this item.
	URI DocumentURI `json:"uri"`
	// The range enclosing this symbol not including leading/trailing whitespace but everything else, e.g. comments and code.
	Range Range `json:"range"`
	// The range that should be selected and revealed when this symbol is being picked, e.g. the name of a function.
	// Must be contained by the [CallHierarchyItem.Range].
	SelectionRange Range `json:"selectionRange"`
	// A data entry field that is preserved between a call hierarchy prepare and
	// incoming calls or outgoing calls requests.
	Data json.RawMessage `json:"data,omitempty"`
}

// CallHierarchyIncomingCallsParams is the parameter of a `callHierarchy/incomingCalls` request.
//
// Since 3.16.0.
type CallHierarchyIncomingCallsParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item CallHierarchyItem `json:"item"`
}

// CallHierarchyIncomingCall represents an incoming call, e.g. a caller of a method or constructor.
//
// Since 3.16.0.
type CallHierarchyIncomingCall struct {
	// The item that makes the call.
	From CallHierarchyItem `json:"from"`
	// The ranges at which the calls appear. This is relative to the caller
	// denoted by [CallHierarchyIncomingCall.From].
	FromRanges []Range `json:"fromRanges"`
}

// CallHierarchyOutgoingCallsParams is the parameter of a `callHierarchy/outgoingCalls` request.
//
// Since 3.16.0.
type CallHierarchyOutgoingCallsParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item CallHierarchyItem `json:"item"`
}

// CallHierarchyOutgoingCall represents an outgoing call, e.g. calling a getter from a method or a method from a constructor etc.
//
// Since 3.16.0.
type CallHierarchyOutgoingCall struct {
	// The item that is called.
	To CallHierarchyItem `json:"to"`
	// The range at which this item is called. This is the range relative to the caller, e.g the item
	// passed to [CallHierarchyItemProvider.ProvideCallHierarchyOutgoingCalls]
	// and not [CallHierarchyOutgoingCall.To].
	FromRanges []Range `json:"fromRanges"`
}
