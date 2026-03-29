package lsp

// TextDocumentIdentifier identifies a text document using a URI.
type TextDocumentIdentifier struct {
	URI DocumentURI `json:"uri"`
}

// VersionedTextDocumentIdentifier extends TextDocumentIdentifier with a version.
type VersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version int `json:"version"`
}

// OptionalVersionedTextDocumentIdentifier allows a null version.
type OptionalVersionedTextDocumentIdentifier struct {
	TextDocumentIdentifier
	Version *int `json:"version"`
}

// TextDocumentItem carries the full content of a newly opened document along with its URI, language, and version.
type TextDocumentItem struct {
	URI        DocumentURI `json:"uri"`
	LanguageID string      `json:"languageId"`
	Version    int         `json:"version"`
	Text       string      `json:"text"`
}

// TextDocumentPositionParams identifies a specific cursor position in a document, used as the base for hover, completion, go-to-definition, etc.
type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// TextEdit replaces a range in a document with new text (or inserts if the range is empty).
type TextEdit struct {
	Range   Range  `json:"range"`
	NewText string `json:"newText"`
}

// AnnotatedTextEdit extends TextEdit with a change annotation.
type AnnotatedTextEdit struct {
	TextEdit
	AnnotationID ChangeAnnotationIdentifier `json:"annotationId"`
}

// ChangeAnnotationIdentifier is an identifier for a change annotation.
type ChangeAnnotationIdentifier string

// ChangeAnnotation describes additional metadata for a change.
type ChangeAnnotation struct {
	Label             string `json:"label"`
	NeedsConfirmation *bool  `json:"needsConfirmation,omitempty"`
	Description       string `json:"description,omitempty"`
}

// TextDocumentContentChangeEvent describes a content change event.
type TextDocumentContentChangeEvent struct {
	Range       *Range `json:"range,omitempty"`
	RangeLength *int   `json:"rangeLength,omitempty"`
	Text        string `json:"text"`
}

// DidOpenTextDocumentParams notifies the server that a document was opened in the editor, carrying its full content.
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidChangeTextDocumentParams notifies the server of edits made to an open document.
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// DidCloseTextDocumentParams notifies the server that a document is no longer open in the editor.
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DidSaveTextDocumentParams notifies the server that a document was saved, optionally including its text.
type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Text         *string                `json:"text,omitempty"`
}

// WillSaveTextDocumentParams notifies the server that a document is about to be saved, allowing it to compute pre-save edits.
type WillSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Reason       TextDocumentSaveReason `json:"reason"`
}

// TextDocumentSaveReason is an int enum indicating why a save occurred (manual, afterDelay, or focusOut).
type TextDocumentSaveReason int

const (
	SaveManual     TextDocumentSaveReason = 1
	SaveAfterDelay TextDocumentSaveReason = 2
	SaveFocusOut   TextDocumentSaveReason = 3
)
