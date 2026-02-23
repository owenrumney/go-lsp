package lsp

import "encoding/json"

// InlayHintKind defines the kind of an inlay hint.
type InlayHintKind int

const (
	InlayHintKindType      InlayHintKind = 1
	InlayHintKindParameter InlayHintKind = 2
)

// InlayHintParams contains the params for textDocument/inlayHint.
type InlayHintParams struct {
	WorkDoneProgressParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
}

// InlayHintLabelPart represents a part of an inlay hint label.
type InlayHintLabelPart struct {
	Value    string         `json:"value"`
	Tooltip  *MarkupContent `json:"tooltip,omitempty"`
	Location *Location      `json:"location,omitempty"`
	Command  *Command       `json:"command,omitempty"`
}

// InlayHint represents an inlay hint.
type InlayHint struct {
	Position     Position        `json:"position"`
	Label        json.RawMessage `json:"label"`
	Kind         *InlayHintKind  `json:"kind,omitempty"`
	TextEdits    []TextEdit      `json:"textEdits,omitempty"`
	Tooltip      *MarkupContent  `json:"tooltip,omitempty"`
	PaddingLeft  *bool           `json:"paddingLeft,omitempty"`
	PaddingRight *bool           `json:"paddingRight,omitempty"`
	Data         json.RawMessage `json:"data,omitempty"`
}

// InlayHintOptions describes inlay hint options.
type InlayHintOptions struct {
	WorkDoneProgressOptions
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

// InlayHintClientCapabilities defines capabilities for inlay hints.
type InlayHintClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	ResolveSupport      *struct {
		Properties []string `json:"properties"`
	} `json:"resolveSupport,omitempty"`
}

// InlayHintWorkspaceClientCapabilities defines workspace capabilities for inlay hints.
type InlayHintWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
