package lsp

import (
	"encoding/json"
	"fmt"
)

// HoverParams holds the parameters for a [HoverRequest].
type HoverParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
}

// Hover is the result of a hover request.
type Hover struct {
	// The hover's content
	Contents HoverContents `json:"contents"`
	// An optional range inside the text document that is used to
	// visualize the hover, e.g. by changing the background color.
	Range *Range `json:"range,omitempty"`
}

// HoverContents is the spec union MarkedString | MarkedString[] | MarkupContent.
// Use [NewHoverContents] for the preferred markup form.
type HoverContents struct {
	Markup        *MarkupContent
	MarkedStrings []MarkedString
}

// NewHoverContents returns hover contents in the preferred MarkupContent form.
func NewHoverContents(kind MarkupKind, value string) HoverContents {
	return HoverContents{Markup: &MarkupContent{Kind: kind, Value: value}}
}

// NewHoverMarkedStrings returns hover contents in the deprecated MarkedString form.
func NewHoverMarkedStrings(items ...MarkedString) HoverContents {
	return HoverContents{MarkedStrings: items}
}

// Value returns the content text regardless of variant (first item for
// multiple MarkedStrings).
func (h HoverContents) Value() string {
	switch {
	case h.Markup != nil:
		return h.Markup.Value
	case len(h.MarkedStrings) > 0:
		return h.MarkedStrings[0].Value
	}
	return ""
}

func (h HoverContents) MarshalJSON() ([]byte, error) {
	switch {
	case h.Markup != nil:
		return json.Marshal(h.Markup)
	case len(h.MarkedStrings) == 1:
		return json.Marshal(h.MarkedStrings[0])
	case h.MarkedStrings != nil:
		return json.Marshal(h.MarkedStrings)
	}
	return json.Marshal(MarkupContent{})
}

func (h *HoverContents) UnmarshalJSON(data []byte) error {
	*h = HoverContents{}
	if len(data) == 0 {
		return fmt.Errorf("lsp: empty hover contents")
	}
	switch data[0] {
	case '[':
		return json.Unmarshal(data, &h.MarkedStrings)
	case '"':
		var ms MarkedString
		if err := json.Unmarshal(data, &ms); err != nil {
			return err
		}
		h.MarkedStrings = []MarkedString{ms}
		return nil
	}
	// MarkupContent has a kind field, MarkedString does not.
	var probe struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(data, &probe); err != nil {
		return err
	}
	if probe.Kind != "" {
		h.Markup = &MarkupContent{}
		return json.Unmarshal(data, h.Markup)
	}
	var ms MarkedString
	if err := json.Unmarshal(data, &ms); err != nil {
		return err
	}
	h.MarkedStrings = []MarkedString{ms}
	return nil
}

// MarkedString is either a bare markdown string (empty Language) or a
// {language, value} code block.
//
// Deprecated: use MarkupContent instead.
type MarkedString struct {
	Language string `json:"language,omitempty"`
	Value    string `json:"value"`
}

func (m MarkedString) MarshalJSON() ([]byte, error) {
	if m.Language == "" {
		return json.Marshal(m.Value)
	}
	type alias MarkedString
	return json.Marshal(alias(m))
}

func (m *MarkedString) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		m.Language = ""
		return json.Unmarshal(data, &m.Value)
	}
	type alias MarkedString
	return json.Unmarshal(data, (*alias)(m))
}
