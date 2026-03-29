package lsp

// DocumentURI is a URI identifying a text document, typically using the file:// scheme but other schemes are permitted.
type DocumentURI string

// URI is a string-encoded URI as defined by RFC 3986.
type URI string

// Position in a text document expressed as zero-based line and character offset.
type Position struct {
	Line      int `json:"line"`
	Character int `json:"character"`
}

// Range in a text document expressed as start and end positions.
type Range struct {
	Start Position `json:"start"`
	End   Position `json:"end"`
}

// Location identifies a range within a text document, used for go-to-definition results, references, etc.
type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

// LocationLink connects a source selection range to a target definition, used for rich go-to-definition results that highlight the origin.
type LocationLink struct {
	OriginSelectionRange *Range      `json:"originSelectionRange,omitempty"`
	TargetURI            DocumentURI `json:"targetUri"`
	TargetRange          Range       `json:"targetRange"`
	TargetSelectionRange Range       `json:"targetSelectionRange"`
}
