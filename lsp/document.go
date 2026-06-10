package lsp

// TextDocumentIdentifier is a literal to identify a text document in the client.
type TextDocumentIdentifier struct {
	// The text document's uri.
	URI DocumentURI `json:"uri"`
}

// VersionedTextDocumentIdentifier is a text document identifier to denote a specific version of a text document.
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	// The version number of this document.
	Version int `json:"version"`
}

// OptionalVersionedTextDocumentIdentifier is a text document identifier to optionally denote a specific version of a text document.
type OptionalVersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	// The version number of this document. If a versioned text document identifier
	// is sent from the server to the client and the file is not open in the editor
	// (the server has not received an open notification before) the server can send
	// null to indicate that the version is unknown and the content on disk is the
	// truth (as specified with document content ownership).
	Version *int `json:"version"`
}

// TextDocumentItem is an item to transfer a text document from the client to the
// server.
type TextDocumentItem struct {
	// The text document's uri.
	URI DocumentURI `json:"uri"`
	// The text document's language identifier.
	LanguageID string `json:"languageId"`
	// The version number of this document (it will increase after each
	// change, including undo/redo).
	Version int `json:"version"`
	// The content of the opened text document.
	Text string `json:"text"`
}

// TextDocumentPositionParams is a parameter literal used in requests to pass a text document and a position inside that
// document.
type TextDocumentPositionParams struct {
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The position inside the text document.
	Position Position `json:"position"`
}

// TextEdit is an edit applicable to a text document.
type TextEdit struct {
	// The range of the text document to be manipulated. To insert
	// text into a document, create a range where start == end.
	Range Range `json:"range"`
	// The string to be inserted. For delete operations, use an
	// empty string.
	NewText string `json:"newText"`
}

// AnnotatedTextEdit is a special text edit with an additional change annotation.
//
// Since 3.16.0.
type AnnotatedTextEdit struct {
	TextEdit
	// The actual identifier of the change annotation
	AnnotationID ChangeAnnotationIdentifier `json:"annotationId"`
}

// ChangeAnnotationIdentifier is an identifier for a change annotation.
type ChangeAnnotationIdentifier string

// ChangeAnnotation is the additional information that describes document changes.
//
// Since 3.16.0.
type ChangeAnnotation struct {
	// A human-readable string describing the actual change. The string
	// is rendered prominently in the user interface.
	Label string `json:"label"`
	// A flag which indicates that user confirmation is needed
	// before applying the change.
	NeedsConfirmation *bool `json:"needsConfirmation,omitempty"`
	// A human-readable string which is rendered less prominently in
	// the user interface.
	Description string `json:"description,omitempty"`
}

// TextDocumentContentChangeEvent describes a content change event.
type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength *int   `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// DidOpenTextDocumentParams holds the parameters sent in an open text document notification.
type DidOpenTextDocumentParams struct {
	// The document that was opened.
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidChangeTextDocumentParams holds the parameters of a text document change notification.
type DidChangeTextDocumentParams struct {
	// The document that did change. The version number points
	// to the version after all provided content changes have
	// been applied.
	TextDocument VersionedTextDocumentIdentifier `json:"textDocument"`
	// The actual content changes. The content changes describe single state changes
	// to the document. So if there are two content changes c1 (at array index 0) and
	// c2 (at array index 1) for a document in state S then c1 moves the document from
	// S to S' and c2 from S' to S''. So c1 is computed on the state S and c2 is computed
	// on the state S'.
	//
	// To mirror the content of a document using change events use the following approach:
	// - start with the same initial content
	// - apply the 'textDocument/didChange' notifications in the order you receive them.
	// - apply the TextDocumentContentChangeEvents in a single notification in the order
	//   you receive them.
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// DidCloseTextDocumentParams holds the parameters sent in a close text document notification.
type DidCloseTextDocumentParams struct {
	// The document that was closed.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DidSaveTextDocumentParams holds the parameters sent in a save text document notification.
type DidSaveTextDocumentParams struct {
	// The document that was saved.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// Optionally, the content when saved. Depends on the includeText value
	// when the save notification was requested.
	Text *string `json:"text,omitempty"`
}

// WillSaveTextDocumentParams holds the parameters sent in a will save text document notification.
type WillSaveTextDocumentParams struct {
	// The document that will be saved.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The 'TextDocumentSaveReason'.
	Reason TextDocumentSaveReason `json:"reason"`
}

// TextDocumentSaveReason is an int enum indicating why a save occurred (manual, afterDelay, or focusOut).
type TextDocumentSaveReason int

const (
	// Manually triggered, e.g. by the user pressing save, by starting debugging,
	// or by an API call.
	SaveManual TextDocumentSaveReason = 1
	// Automatic after a delay.
	SaveAfterDelay TextDocumentSaveReason = 2
	// When the editor lost focus.
	SaveFocusOut TextDocumentSaveReason = 3
)
