package lsp

import "encoding/json"

// CodeLensParams is sent to request code lenses (actionable annotations) for a document.
type CodeLensParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// CodeLens is an actionable annotation displayed inline (e.g. "3 references", "Run test"), bound to a source range and a command.
type CodeLens struct {
	Range   Range           `json:"range"`
	Command *Command        `json:"command,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}
