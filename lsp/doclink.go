package lsp

import "encoding/json"

// DocumentLinkParams contains the params for textDocument/documentLink.
type DocumentLinkParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DocumentLink represents a link found in a document.
type DocumentLink struct {
	Range   Range           `json:"range"`
	Target  *DocumentURI    `json:"target,omitempty"`
	Tooltip string          `json:"tooltip,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}
