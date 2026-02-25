package lsp

import "encoding/json"

// InitializeParams contains the params for the initialize request.
type InitializeParams struct {
	WorkDoneProgressParams
	ProcessID             *int               `json:"processId"`
	ClientInfo            *ClientInfo        `json:"clientInfo,omitempty"`
	Locale                string             `json:"locale,omitempty"`
	RootPath              *string            `json:"rootPath,omitempty"`
	RootURI               *DocumentURI       `json:"rootUri"`
	InitializationOptions json.RawMessage    `json:"initializationOptions,omitempty"`
	Capabilities          ClientCapabilities `json:"capabilities"`
	Trace                 string             `json:"trace,omitempty"`
	WorkspaceFolders      []WorkspaceFolder  `json:"workspaceFolders,omitempty"`
}

// ClientInfo contains information about the client.
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// InitializeResult contains the result of the initialize request.
type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   *ServerInfo        `json:"serverInfo,omitempty"`
}

// ServerInfo contains information about the server.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// InitializedParams contains the params for the initialized notification.
type InitializedParams struct{}

// ReferenceParams contains the params for textDocument/references.
type ReferenceParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
	Context ReferenceContext `json:"context"`
}

// ReferenceContext contains additional information for reference requests.
type ReferenceContext struct {
	IncludeDeclaration bool `json:"includeDeclaration"`
}

// DeclarationParams contains the params for textDocument/declaration.
type DeclarationParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// DefinitionParams contains the params for textDocument/definition.
type DefinitionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// TypeDefinitionParams contains the params for textDocument/typeDefinition.
type TypeDefinitionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// ImplementationParams contains the params for textDocument/implementation.
type ImplementationParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// SetTraceParams contains the params for the $/setTrace notification.
type SetTraceParams struct {
	Value TraceValue `json:"value"`
}

// DocumentHighlightParams contains the params for textDocument/documentHighlight.
type DocumentHighlightParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// DocumentHighlight represents a highlight in a document.
type DocumentHighlight struct {
	Range Range                  `json:"range"`
	Kind  *DocumentHighlightKind `json:"kind,omitempty"`
}

// DocumentHighlightKind represents the kind of a document highlight.
type DocumentHighlightKind int

const (
	DocumentHighlightKindText  DocumentHighlightKind = 1
	DocumentHighlightKindRead  DocumentHighlightKind = 2
	DocumentHighlightKindWrite DocumentHighlightKind = 3
)
