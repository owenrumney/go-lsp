package lsp

import "encoding/json"

// CallHierarchyPrepareParams is sent to resolve the call hierarchy item at a given cursor position before navigating callers/callees.
type CallHierarchyPrepareParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// CallHierarchyItem represents a function or method that can be navigated to in a call hierarchy view.
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

// CallHierarchyIncomingCallsParams is sent to find all callers of a given call hierarchy item.
type CallHierarchyIncomingCallsParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item CallHierarchyItem `json:"item"`
}

// CallHierarchyIncomingCall identifies a caller of a call hierarchy item, including the ranges where the call occurs.
type CallHierarchyIncomingCall struct {
	From       CallHierarchyItem `json:"from"`
	FromRanges []Range           `json:"fromRanges"`
}

// CallHierarchyOutgoingCallsParams is sent to find all functions/methods called from a given call hierarchy item.
type CallHierarchyOutgoingCallsParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Item CallHierarchyItem `json:"item"`
}

// CallHierarchyOutgoingCall identifies a function/method called from a call hierarchy item, including the call-site ranges.
type CallHierarchyOutgoingCall struct {
	To         CallHierarchyItem `json:"to"`
	FromRanges []Range           `json:"fromRanges"`
}
