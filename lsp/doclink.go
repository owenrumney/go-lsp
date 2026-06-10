package lsp

import "encoding/json"

// DocumentLinkParams holds the parameters of a [DocumentLinkRequest].
type DocumentLinkParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The document to provide document links for.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DocumentLink is a range in a text document that links to an internal or external resource, like another
// text document or a web site.
type DocumentLink struct {
	// The range this link applies to.
	Range Range `json:"range"`
	// The uri this link points to. If missing, a resolve request is sent later.
	Target *DocumentURI `json:"target,omitempty"`
	// The tooltip text when you hover over this link.
	//
	// If a tooltip is provided, it will be displayed in a string that includes instructions on how to
	// trigger the link, such as `{0} (ctrl + click)`. The specific instructions vary depending on OS,
	// user settings, and localization.
	//
	// Since 3.15.0
	Tooltip string `json:"tooltip,omitempty"`
	// A data entry field that is preserved on a document link between a
	// DocumentLinkRequest and a DocumentLinkResolveRequest.
	Data json.RawMessage `json:"data,omitempty"`
}
