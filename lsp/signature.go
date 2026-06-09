package lsp

// SignatureHelpParams holds the parameters for a [SignatureHelpRequest].
type SignatureHelpParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	// The signature help context. This is only available if the client specifies
	// to send this using the client capability `textDocument.signatureHelp.contextSupport === true`
	//
	// Since 3.15.0
	Context *SignatureHelpContext `json:"context,omitempty"`
}

// SignatureHelpContext is the additional information about the context in which a signature help request was triggered.
//
// Since 3.15.0.
type SignatureHelpContext struct {
	// Action that caused signature help to be triggered.
	TriggerKind SignatureHelpTriggerKind `json:"triggerKind"`
	// Character that caused signature help to be triggered.
	//
	// This is empty when `triggerKind !== [SignatureHelpTriggerCharacter]`
	TriggerCharacter string `json:"triggerCharacter,omitempty"`
	// true if signature help was already showing when it was triggered.
	//
	// Retriggers occurs when the signature help is already active and can be caused by actions such as
	// typing a trigger character, a cursor move, or document content changes.
	IsRetrigger bool `json:"isRetrigger"`
	// The currently active SignatureHelp.
	//
	// The activeSignatureHelp has its `SignatureHelp.activeSignature` field updated based on
	// the user navigating through available signatures.
	ActiveSignatureHelp *SignatureHelp `json:"activeSignatureHelp,omitempty"`
}

// SignatureHelp represents the signature of something
// callable. There can be multiple signatures but only one
// active and only one active parameter.
type SignatureHelp struct {
	// One or more signatures.
	Signatures []SignatureInformation `json:"signatures"`
	// The active signature. If omitted or the value lies outside the
	// range of signatures the value defaults to zero or is ignored if
	// the SignatureHelp has no signatures.
	//
	// Whenever possible implementors should make an active decision about
	// the active signature and shouldn't rely on a default value.
	//
	// In future version of the protocol this property might become
	// mandatory to better express this.
	ActiveSignature *int `json:"activeSignature,omitempty"`
	// The active parameter of the active signature. If omitted or the value
	// lies outside the range of `signatures[activeSignature].parameters`
	// defaults to 0 if the active signature has parameters. If
	// the active signature has no parameters it is ignored.
	// In future version of the protocol this property might become
	// mandatory to better express the active parameter if the
	// active signature does have any.
	ActiveParameter *int `json:"activeParameter,omitempty"`
}

// SignatureInformation represents the signature of something callable. A signature
// can have a label, like a function-name, a doc-comment, and
// a set of parameters.
type SignatureInformation struct {
	// The label of this signature. Will be shown in
	// the UI.
	Label string `json:"label"`
	// The human-readable doc-comment of this signature. Will be shown
	// in the UI but can be omitted.
	Documentation *MarkupContent `json:"documentation,omitempty"`
	// The parameters of this signature.
	Parameters []ParameterInformation `json:"parameters,omitempty"`
	// The index of the active parameter.
	//
	// If provided, this is used in place of `SignatureHelp.activeParameter`.
	//
	// Since 3.16.0
	ActiveParameter *int `json:"activeParameter,omitempty"`
}

// ParameterInformation represents a parameter of a callable-signature. A parameter can
// have a label and a doc-comment.
type ParameterInformation struct {
	// The label of this parameter information.
	//
	// Either a string or an inclusive start and exclusive end offsets within its containing
	// signature label. (see SignatureInformation.label). The offsets are based on a UTF-16
	// string representation as Position and Range does.
	//
	// *Note*: a label of type string should be a substring of its containing signature label.
	// Its intended use case is to highlight the parameter label part in the `SignatureInformation.label`.
	Label any `json:"label"` // string | [int, int]
	// The human-readable doc-comment of this parameter. Will be shown
	// in the UI but can be omitted.
	Documentation *MarkupContent `json:"documentation,omitempty"`
}
