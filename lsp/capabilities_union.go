package lsp

import "encoding/json"

// UnmarshalJSON tolerates servers that send an options object where this
// struct types a provider as a bool.
//
// The spec types these as `boolean | XOptions`. Fields typed here as options
// handle the union themselves through their own UnmarshalJSON and Enabled;
// the bool-typed fields cannot, so an options object is coerced to true and
// its contents discarded. A literal false decodes to nil, so a nil check and
// a false check agree.
func (c *ServerCapabilities) UnmarshalJSON(data []byte) error {
	type plain ServerCapabilities

	var shadow struct {
		plain
		HoverProvider                   json.RawMessage `json:"hoverProvider,omitempty"`
		DeclarationProvider             json.RawMessage `json:"declarationProvider,omitempty"`
		DefinitionProvider              json.RawMessage `json:"definitionProvider,omitempty"`
		TypeDefinitionProvider          json.RawMessage `json:"typeDefinitionProvider,omitempty"`
		ImplementationProvider          json.RawMessage `json:"implementationProvider,omitempty"`
		ReferencesProvider              json.RawMessage `json:"referencesProvider,omitempty"`
		DocumentHighlightProvider       json.RawMessage `json:"documentHighlightProvider,omitempty"`
		ColorProvider                   json.RawMessage `json:"colorProvider,omitempty"`
		DocumentFormattingProvider      json.RawMessage `json:"documentFormattingProvider,omitempty"`
		DocumentRangeFormattingProvider json.RawMessage `json:"documentRangeFormattingProvider,omitempty"`
		FoldingRangeProvider            json.RawMessage `json:"foldingRangeProvider,omitempty"`
		SelectionRangeProvider          json.RawMessage `json:"selectionRangeProvider,omitempty"`
		LinkedEditingRangeProvider      json.RawMessage `json:"linkedEditingRangeProvider,omitempty"`
		CallHierarchyProvider           json.RawMessage `json:"callHierarchyProvider,omitempty"`
		MonikerProvider                 json.RawMessage `json:"monikerProvider,omitempty"`
		TypeHierarchyProvider           json.RawMessage `json:"typeHierarchyProvider,omitempty"`
		InlineValueProvider             json.RawMessage `json:"inlineValueProvider,omitempty"`
	}

	if err := json.Unmarshal(data, &shadow); err != nil {
		return err
	}
	*c = ServerCapabilities(shadow.plain)

	c.HoverProvider = providerBool(shadow.HoverProvider)
	c.DeclarationProvider = providerBool(shadow.DeclarationProvider)
	c.DefinitionProvider = providerBool(shadow.DefinitionProvider)
	c.TypeDefinitionProvider = providerBool(shadow.TypeDefinitionProvider)
	c.ImplementationProvider = providerBool(shadow.ImplementationProvider)
	c.ReferencesProvider = providerBool(shadow.ReferencesProvider)
	c.DocumentHighlightProvider = providerBool(shadow.DocumentHighlightProvider)
	c.ColorProvider = providerBool(shadow.ColorProvider)
	c.DocumentFormattingProvider = providerBool(shadow.DocumentFormattingProvider)
	c.DocumentRangeFormattingProvider = providerBool(shadow.DocumentRangeFormattingProvider)
	c.FoldingRangeProvider = providerBool(shadow.FoldingRangeProvider)
	c.SelectionRangeProvider = providerBool(shadow.SelectionRangeProvider)
	c.LinkedEditingRangeProvider = providerBool(shadow.LinkedEditingRangeProvider)
	c.CallHierarchyProvider = providerBool(shadow.CallHierarchyProvider)
	c.MonikerProvider = providerBool(shadow.MonikerProvider)
	c.TypeHierarchyProvider = providerBool(shadow.TypeHierarchyProvider)
	c.InlineValueProvider = providerBool(shadow.InlineValueProvider)

	return nil
}

func providerBool(raw json.RawMessage) *bool {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		if !b {
			return nil
		}
		return &b
	}

	provided := true
	return &provided
}

// semanticTokensRangeFromJSON reads `range?: boolean | {}`. The object form is
// empty, so presence means supported.
func semanticTokensRangeFromJSON(raw json.RawMessage) *bool {
	if len(raw) == 0 || string(raw) == "null" {
		return nil
	}

	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		return &b
	}

	supported := true
	return &supported
}

// semanticTokensFullFromJSON reads `full?: boolean | {delta?: boolean}`.
// false yields nil, so unsupported reads the same as absent.
func semanticTokensFullFromJSON(raw json.RawMessage) (*SemanticTokensFull, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return nil, nil
	}

	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		if !b {
			return nil, nil
		}
		return &SemanticTokensFull{}, nil
	}

	var full SemanticTokensFull
	if err := json.Unmarshal(raw, &full); err != nil {
		return nil, err
	}
	return &full, nil
}
