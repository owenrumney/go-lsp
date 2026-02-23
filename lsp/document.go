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

// TextDocumentItem represents an open text document.
type TextDocumentItem struct {
	URI        DocumentURI `json:"uri"`
	LanguageID string      `json:"languageId"`
	Version    int         `json:"version"`
	Text       string      `json:"text"`
}

// TextDocumentPositionParams is a parameter literal for requests that take a position inside a text document.
type TextDocumentPositionParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Position     Position               `json:"position"`
}

// TextEdit represents a textual edit applicable to a text document.
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

// DidOpenTextDocumentParams contains the params for textDocument/didOpen.
type DidOpenTextDocumentParams struct {
	TextDocument TextDocumentItem `json:"textDocument"`
}

// DidChangeTextDocumentParams contains the params for textDocument/didChange.
type DidChangeTextDocumentParams struct {
	TextDocument   VersionedTextDocumentIdentifier  `json:"textDocument"`
	ContentChanges []TextDocumentContentChangeEvent `json:"contentChanges"`
}

// DidCloseTextDocumentParams contains the params for textDocument/didClose.
type DidCloseTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
}

// DidSaveTextDocumentParams contains the params for textDocument/didSave.
type DidSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Text         *string                `json:"text,omitempty"`
}

// WillSaveTextDocumentParams contains the params for textDocument/willSave.
type WillSaveTextDocumentParams struct {
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Reason       TextDocumentSaveReason `json:"reason"`
}

// TextDocumentSaveReason represents the reason why a text document is saved.
type TextDocumentSaveReason int

const (
	SaveManual     TextDocumentSaveReason = 1
	SaveAfterDelay TextDocumentSaveReason = 2
	SaveFocusOut   TextDocumentSaveReason = 3
)
