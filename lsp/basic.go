package lsp

// DocumentURI is a URI identifying a text document, typically using the file:// scheme, but other schemes are permitted.
type DocumentURI string

// URI is a string-encoded URI as defined by RFC 3986.
type URI string

// Position is a zero-based line and character offset in a text document.
// Prior to 3.17 the offsets were always based on a UTF-16 string
// representation. So for a string of the form `a𐐀b`, the character offset of the
// character a is 0, the character offset of `𐐀` is 1 and the character
// offset of b is 3 since `𐐀` is represented using two code units in UTF-16.
// Since 3.17 clients and servers can agree on a different string encoding
// representation (e.g. UTF-8). The client announces its supported encoding
// via the client capability [general.positionEncodings].
// The value is an array of position encodings the client supports, with
// decreasing preference (e.g. the encoding at index `0` is the most preferred
// one). To stay backwards compatible the only mandatory encoding is UTF-16
// represented via the string `utf-16`. The server can pick one of the
// encodings offered by the client and signals that encoding back to the
// client via the initialize result's property
// [capabilities.positionEncoding]. If the string value
// `utf-16` is missing from the client's capability `general.positionEncodings`
// servers can safely assume that the client supports UTF-16. If the server
// omits the position encoding in its initialize result the encoding defaults
// to the string value `utf-16`. Implementation considerations: since the
// conversion from one encoding into another requires the content of the
// file / line the conversion is best done where the file is read which is
// usually on the server side.
//
// Positions are line end character agnostic. So you can not specify a position
// that denotes `\r|\n` or `\n|` where `|` represents the character offset.
//
// Since 3.17.0 - support for negotiated position encoding.
//
// [general.positionEncodings]: https://microsoft.github.io/language-server-protocol/specifications/specification-current/#clientCapabilities
// [capabilities.positionEncoding]: https://microsoft.github.io/language-server-protocol/specifications/specification-current/#serverCapabilities
type Position struct {
	// Line position in a document (zero-based).
	//
	// If a line number is greater than the number of lines in a document, it defaults back to the number of lines in the document.
	// If a line number is negative, it defaults to 0.
	Line int `json:"line"`
	// Character offset on a line in a document (zero-based).
	//
	// The meaning of this offset is determined by the negotiated
	// PositionEncodingKind.
	//
	// If the character value is greater than the line length it defaults back to the
	// line length.
	Character int `json:"character"`
}

// Range is a span in a text document, expressed as (zero-based) start and end positions.
//
// If you want to specify a range that contains a line including the line ending
// character(s) then use an end position denoting the start of the next line.
// For example:
//
//	{
//	    start: { line: 5, character: 23 }
//	    end : { line 6, character : 0 }
//	}
type Range struct {
	// The range's start position.
	Start Position `json:"start"`
	// The range's end position.
	End Position `json:"end"`
}

// Location represents a location inside a resource, such as a line
// inside a text file.
type Location struct {
	URI   DocumentURI `json:"uri"`
	Range Range       `json:"range"`
}

// LocationLink represents the connection of two locations. Provides additional metadata over normal [Location],
// including an origin range.
type LocationLink struct {
	// Span of the origin of this link.
	//
	// Used as the underlined span for mouse interaction. Defaults to the word range at
	// the definition position.
	OriginSelectionRange *Range `json:"originSelectionRange,omitempty"`
	// The target resource identifier of this link.
	TargetURI DocumentURI `json:"targetUri"`
	// The full target range of this link. If the target for example is a symbol then target range is the
	// range enclosing this symbol not including leading/trailing whitespace but everything else
	// like comments. This information is typically used to highlight the range in the editor.
	TargetRange Range `json:"targetRange"`
	// The range that should be selected and revealed when this link is being followed, e.g. the name of a function.
	// Must be contained by the targetRange. See also `DocumentSymbol#range`
	TargetSelectionRange Range `json:"targetSelectionRange"`
}
