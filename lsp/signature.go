package lsp

// SignatureHelpParams contains the params for textDocument/signatureHelp.
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

// SignatureHelp represents the signature of something callable.
type SignatureHelp struct {
	Signatures      []SignatureInformation `json:"signatures"`
	ActiveSignature *int                   `json:"activeSignature,omitempty"`
	ActiveParameter *int                   `json:"activeParameter,omitempty"`
}

// SignatureInformation represents the signature of something callable.
type SignatureInformation struct {
	Label           string                 `json:"label"`
	Documentation   *MarkupContent         `json:"documentation,omitempty"`
	Parameters      []ParameterInformation `json:"parameters,omitempty"`
	ActiveParameter *int                   `json:"activeParameter,omitempty"`
}

// ParameterInformation represents a parameter of a callable-signature.
type ParameterInformation struct {
	Label         any            `json:"label"` // string | [int, int]
	Documentation *MarkupContent `json:"documentation,omitempty"`
}
