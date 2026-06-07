package lsp

import "encoding/json"

// Registration is the general parameters to register for a notification or to register a provider.
type Registration struct {
	// The id used to register the request. The id can be used to deregister
	// the request again.
	ID string `json:"id"`
	// The method / capability to register for.
	Method string `json:"method"`
	// Options necessary for the registration.
	RegisterOptions json.RawMessage `json:"registerOptions,omitempty"`
}

// RegistrationParams is sent from server to client to dynamically register one or more capabilities.
type RegistrationParams struct {
	Registrations []Registration `json:"registrations"`
}

// Unregistration is the general parameters to unregister a request or notification.
type Unregistration struct {
	// The id used to unregister the request or notification. Usually an id
	// provided during the register request.
	ID string `json:"id"`
	// The method to unregister for.
	Method string `json:"method"`
}

// UnregistrationParams is sent from server to client to remove previously registered capabilities.
type UnregistrationParams struct {
	Unregisterations []Unregistration `json:"unregisterations"`
}

// TextDocumentRegistrationOptions is general text document registration options.
type TextDocumentRegistrationOptions struct {
	// A document selector to identify the scope of the registration. If set to null
	// the document selector provided on the client side will be used.
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

// TextDocumentChangeRegistrationOptions is the describe options to be used when registered for text document change events.
type TextDocumentChangeRegistrationOptions struct {
	TextDocumentRegistrationOptions
	// How documents are synced to the server.
	SyncKind TextDocumentSyncKind `json:"syncKind"`
}

// TextDocumentSaveRegistrationOptions are the registration options for the textDocument/didSave notification.
type TextDocumentSaveRegistrationOptions struct {
	TextDocumentRegistrationOptions
	// The client is supposed to include the content on save.
	IncludeText *bool `json:"includeText,omitempty"`
}

// DidChangeWatchedFilesRegistrationOptions is the describe options to be used when registered for text document change events.
type DidChangeWatchedFilesRegistrationOptions struct {
	// The watchers to register.
	Watchers []FileSystemWatcher `json:"watchers"`
}
