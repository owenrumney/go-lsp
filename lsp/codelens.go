package lsp

import "encoding/json"

// CodeLensParams is the parameters of a [CodeLensRequest].
type CodeLensParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The document to request code lens for.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// CodeLens represents a [Command] that should be shown along with
// source text, like the number of references, a way to run tests, etc.
//
// A code lens is _unresolved_ when no command is associated to it. For performance
// reasons the creation of a code lens and resolving should be done in two stages.
type CodeLens struct {
	// The range in which this code lens is valid. Should only span a single line.
	Range Range `json:"range"`
	// The command this code lens represents.
	Command *Command `json:"command,omitempty"`
	// A data entry field that is preserved on a code lens item between
	// a [CodeLensRequest] and a [CodeLensResolveRequest]
	Data json.RawMessage `json:"data,omitempty"`
}
