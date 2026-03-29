package lsp

// SignatureHelpParams is sent to request parameter hints for a function call at the cursor position.
type SignatureHelpParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	Context *SignatureHelpContext `json:"context,omitempty"`
}

// SignatureHelpContext contains additional information about the context of a signature help request.
type SignatureHelpContext struct {
	TriggerKind         SignatureHelpTriggerKind `json:"triggerKind"`
	TriggerCharacter    string                   `json:"triggerCharacter,omitempty"`
	IsRetrigger         bool                     `json:"isRetrigger"`
	ActiveSignatureHelp *SignatureHelp           `json:"activeSignatureHelp,omitempty"`
}

// SignatureHelp contains one or more function signatures and indicates which signature and parameter are currently active.
type SignatureHelp struct {
	Signatures      []SignatureInformation `json:"signatures"`
	ActiveSignature *int                   `json:"activeSignature,omitempty"`
	ActiveParameter *int                   `json:"activeParameter,omitempty"`
}

// SignatureInformation describes a single callable signature: its label, documentation, and parameter list.
type SignatureInformation struct {
	Label           string                 `json:"label"`
	Documentation   *MarkupContent         `json:"documentation,omitempty"`
	Parameters      []ParameterInformation `json:"parameters,omitempty"`
	ActiveParameter *int                   `json:"activeParameter,omitempty"`
}

// ParameterInformation describes one parameter in a signature: its label (or label range) and optional documentation.
type ParameterInformation struct {
	Label         any            `json:"label"` // string | [int, int]
	Documentation *MarkupContent `json:"documentation,omitempty"`
}
