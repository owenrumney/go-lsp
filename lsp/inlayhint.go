package lsp

import "encoding/json"

// InlayHintKind is an int enum: type hint (1) or parameter name hint (2).
type InlayHintKind int

const (
	InlayHintKindType      InlayHintKind = 1
	InlayHintKindParameter InlayHintKind = 2
)

// InlayHintParams is sent to request inlay hints (inline type/parameter annotations) for a document range.
type InlayHintParams struct {
	WorkDoneProgressParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
}

// InlayHintLabelPart is one segment of a multi-part inlay hint label, optionally clickable with a command or location.
type InlayHintLabelPart struct {
	Value    string         `json:"value"`
	Tooltip  *MarkupContent `json:"tooltip,omitempty"`
	Location *Location      `json:"location,omitempty"`
	Command  *Command       `json:"command,omitempty"`
}

// InlayHint is an inline annotation the editor renders in the source (e.g. inferred types, parameter names).
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

// InlayHintOptions configures whether the server supports resolving inlay hints lazily.
type InlayHintOptions struct {
	WorkDoneProgressOptions
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

// InlayHintClientCapabilities declares editor support for inlay hint features like lazy resolution.
type InlayHintClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	ResolveSupport      *struct {
		Properties []string `json:"properties"`
	} `json:"resolveSupport,omitempty"`
}

// InlayHintWorkspaceClientCapabilities declares whether the editor will refresh inlay hints when the server requests it.
type InlayHintWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
