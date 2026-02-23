package lsp

import "encoding/json"

// Registration represents a general registration request.
type Registration struct {
	ID              string          `json:"id"`
	Method          string          `json:"method"`
	RegisterOptions json.RawMessage `json:"registerOptions,omitempty"`
}

// RegistrationParams contains the params for client/registerCapability.
type RegistrationParams struct {
	Registrations []Registration `json:"registrations"`
}

// Unregistration describes a request to unregister a capability.
type Unregistration struct {
	ID     string `json:"id"`
	Method string `json:"method"`
}

// UnregistrationParams contains the params for client/unregisterCapability.
type UnregistrationParams struct {
	Unregisterations []Unregistration `json:"unregisterations"`
}

// TextDocumentRegistrationOptions describes options for text document registration.
type TextDocumentRegistrationOptions struct {
	DocumentSelector *DocumentSelector `json:"documentSelector"`
}

// DocumentSelector is an array of document filters.
type DocumentSelector []DocumentFilter

// DocumentFilter denotes a document through properties like language, scheme, or pattern.
type DocumentFilter struct {
	Language string `json:"language,omitempty"`
	Scheme   string `json:"scheme,omitempty"`
	Pattern  string `json:"pattern,omitempty"`
}

// TextDocumentChangeRegistrationOptions describes options for text document change registration.
type TextDocumentChangeRegistrationOptions struct {
	TextDocumentRegistrationOptions
	SyncKind TextDocumentSyncKind `json:"syncKind"`
}

// TextDocumentSaveRegistrationOptions describes options for text document save registration.
type TextDocumentSaveRegistrationOptions struct {
	TextDocumentRegistrationOptions
	IncludeText *bool `json:"includeText,omitempty"`
}

// DidChangeWatchedFilesRegistrationOptions describes options for watching file changes.
type DidChangeWatchedFilesRegistrationOptions struct {
	Watchers []FileSystemWatcher `json:"watchers"`
}
