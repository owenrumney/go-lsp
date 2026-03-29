package lsp

import "encoding/json"

// InitializeParams is the first message from client to server, carrying client capabilities, root URI, and configuration.
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

// InitializedParams is sent after the client receives the initialize result; the struct is intentionally empty.
type InitializedParams struct{}

// ReferenceParams is sent to find all references to a symbol at a given position.
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

// DeclarationParams is sent to find the declaration of a symbol at a given position.
type DeclarationParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// DefinitionParams is sent to find the definition of a symbol at a given position.
type DefinitionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// TypeDefinitionParams is sent to find the type definition of a symbol at a given position.
type TypeDefinitionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// ImplementationParams is sent to find implementations of an interface or abstract method at a given position.
type ImplementationParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// SetTraceParams changes the server's trace verbosity level at runtime.
type SetTraceParams struct {
	Value TraceValue `json:"value"`
}

// DocumentHighlightParams is sent to request highlights for all occurrences of a symbol in a document.
type DocumentHighlightParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// DocumentHighlight identifies a range to highlight when the cursor is on a symbol (e.g. all usages of a variable).
type DocumentHighlight struct {
	Range Range                  `json:"range"`
	Kind  *DocumentHighlightKind `json:"kind,omitempty"`
}

// DocumentHighlightKind is an int enum: text (1), read-access (2), or write-access (3).
type DocumentHighlightKind int

const (
	DocumentHighlightKindText  DocumentHighlightKind = 1
	DocumentHighlightKindRead  DocumentHighlightKind = 2
	DocumentHighlightKindWrite DocumentHighlightKind = 3
)
