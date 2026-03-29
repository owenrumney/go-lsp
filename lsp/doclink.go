package lsp

import "encoding/json"

// DocumentLinkParams is sent to request clickable links (URLs, file references) within a document.
type DocumentLinkParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DocumentLink is a clickable range in a document that resolves to a URI (e.g. an import path or URL).
type DocumentLink struct {
	Range   Range           `json:"range"`
	Target  *DocumentURI    `json:"target,omitempty"`
	Tooltip string          `json:"tooltip,omitempty"`
	Data    json.RawMessage `json:"data,omitempty"`
}
