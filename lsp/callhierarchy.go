package lsp

import "encoding/json"

// CallHierarchyPrepareParams contains the params for textDocument/prepareCallHierarchy.
type CallHierarchyPrepareParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// CallHierarchyItem represents an item in the call hierarchy.
type CallHierarchyItem struct {
	Name           string          `json:"name"`
	Kind           SymbolKind      `json:"kind"`
	Tags           []SymbolTag     `json:"tags,omitempty"`
	Detail         string          `json:"detail,omitempty"`
	URI            DocumentURI     `json:"uri"`
	Range          Range           `json:"range"`
	SelectionRange Range           `json:"selectionRange"`
	Data           json.RawMessage `json:"data,omitempty"`
}

// CallHierarchyIncomingCallsParams contains the params for callHierarchy/incomingCalls.
type CallHierarchyIncomingCallsParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item CallHierarchyItem `json:"item"`
}

// CallHierarchyIncomingCall represents an incoming call.
type CallHierarchyIncomingCall struct {
	From       CallHierarchyItem `json:"from"`
	FromRanges []Range           `json:"fromRanges"`
}

// CallHierarchyOutgoingCallsParams contains the params for callHierarchy/outgoingCalls.
type CallHierarchyOutgoingCallsParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item CallHierarchyItem `json:"item"`
}

// CallHierarchyOutgoingCall represents an outgoing call.
type CallHierarchyOutgoingCall struct {
	To         CallHierarchyItem `json:"to"`
	FromRanges []Range           `json:"fromRanges"`
}
