package lsp

import "encoding/json"

// InlayHintKind is an int enum: type hint (1) or parameter name hint (2).
type InlayHintKind int

const (
	// An inlay hint that is for a type annotation.
	InlayHintKindType InlayHintKind = 1
	// An inlay hint that is for a parameter.
	InlayHintKindParameter InlayHintKind = 2
)

// InlayHintParams is a parameter literal used in inlay hint requests.
//
// Since 3.17.0.
type InlayHintParams struct {
	WorkDoneProgressParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The document range for which inlay hints should be computed.
	Range Range `json:"range"`
}

// InlayHintLabelPart allows for interactive and composite labels
// of inlay hints.
//
// Since 3.17.0.
type InlayHintLabelPart struct {
	// The value of this label part.
	Value string `json:"value"`
	// The tooltip text when you hover over this label part. Depending on
	// the client capability `inlayHint.resolveSupport` clients might resolve
	// this property late using the resolve request.
	Tooltip *MarkupContent `json:"tooltip,omitempty"`
	// An optional source code location that represents this
	// label part.
	//
	// The editor will use this location for the hover and for code navigation
	// features: This part will become a clickable link that resolves to the
	// definition of the symbol at the given location (not necessarily the
	// location itself), it shows the hover that shows at the given location,
	// and it shows a context menu with further code navigation commands.
	//
	// Depending on the client capability `inlayHint.resolveSupport` clients
	// might resolve this property late using the resolve request.
	Location *Location `json:"location,omitempty"`
	// An optional command for this label part.
	//
	// Depending on the client capability `inlayHint.resolveSupport` clients
	// might resolve this property late using the resolve request.
	Command *Command `json:"command,omitempty"`
}

// InlayHint represents inlay hint information.
//
// Since 3.17.0.
type InlayHint struct {
	// The position of this hint.
	//
	// If multiple hints have the same position, they will be shown in the order
	// they appear in the response.
	Position Position `json:"position"`
	// The label of this hint. A human readable string or an array of
	// InlayHintLabelPart label parts.
	//
	// *Note* that neither the string nor the label part can be empty.
	Label json.RawMessage `json:"label"`
	// The kind of this hint. Can be omitted in which case the client
	// should fall back to a reasonable default.
	Kind *InlayHintKind `json:"kind,omitempty"`
	// Optional text edits that are performed when accepting this inlay hint.
	//
	// *Note* that edits are expected to change the document so that the inlay
	// hint (or its nearest variant) is now part of the document and the inlay
	// hint itself is now obsolete.
	TextEdits []TextEdit `json:"textEdits,omitempty"`
	// The tooltip text when you hover over this item.
	Tooltip *MarkupContent `json:"tooltip,omitempty"`
	// Render padding before the hint.
	//
	// Note: Padding should use the editor's background color, not the
	// background color of the hint itself. That means padding can be used
	// to visually align/separate an inlay hint.
	PaddingLeft *bool `json:"paddingLeft,omitempty"`
	// Render padding after the hint.
	//
	// Note: Padding should use the editor's background color, not the
	// background color of the hint itself. That means padding can be used
	// to visually align/separate an inlay hint.
	PaddingRight *bool `json:"paddingRight,omitempty"`
	// A data entry field that is preserved on an inlay hint between
	// a `textDocument/inlayHint` and a `inlayHint/resolve` request.
	Data json.RawMessage `json:"data,omitempty"`
}

// InlayHintOptions is used during static registration.
//
// Since 3.17.0.
type InlayHintOptions struct {
	WorkDoneProgressOptions
	// The server provides support to resolve additional
	// information for an inlay hint item.
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

// InlayHintClientCapabilities declares client support for inlay-hint requests.
//
// Since 3.17.0.
type InlayHintClientCapabilities struct {
	// Whether inlay hints support dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// Indicates which properties a client can resolve lazily on an inlay
	// hint.
	ResolveSupport *struct {
		Properties []string `json:"properties"`
	} `json:"resolveSupport,omitempty"`
}

// InlayHintWorkspaceClientCapabilities is specific to inlay hints.
//
// Since 3.17.0.
type InlayHintWorkspaceClientCapabilities struct {
	// Whether the client implementation supports a refresh request sent from
	// the server to the client.
	//
	// Note that this event is global and will force the client to refresh all
	// inlay hints currently shown. It should be used with absolute care and
	// is useful for situation where a server for example detects a project wide
	// change that requires such a calculation.
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
