package lsp

import "encoding/json"

// CodeLensParams contains the params for textDocument/codeLens.
type CodeLensParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// CodeLens represents a command that should be shown along with source text.
type CodeLens struct {
	Range   Range            `json:"range"`
	Command *Command         `json:"command,omitempty"`
	Data    json.RawMessage  `json:"data,omitempty"`
}
