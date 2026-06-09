package lsp

import "encoding/json"

// InitializeParams is the first message from client to server, carrying client capabilities, root URI, and configuration.
type InitializeParams struct {
	WorkDoneProgressParams
	// The process Id of the parent process that started
	// the server.
	//
	// Is null if the process has not been started by another process.
	// If the parent process is not alive then the server should exit.
	ProcessID *int `json:"processId"`
	// Information about the client
	//
	// Since 3.15.0
	ClientInfo *ClientInfo `json:"clientInfo,omitempty"`
	// The locale the client is currently showing the user interface
	// in. This must not necessarily be the locale of the operating
	// system.
	//
	// Uses IETF language tags as the value's syntax
	// (See https://en.wikipedia.org/wiki/IETF_language_tag)
	//
	// Since 3.16.0
	Locale string `json:"locale,omitempty"`
	// The rootPath of the workspace. Is null
	// if no folder is open.
	//
	// Deprecated: In favour of rootUri.
	RootPath *string `json:"rootPath,omitempty"`
	// The rootUri of the workspace. Is null if no
	// folder is open. If both rootPath and rootUri are set
	// rootUri wins.
	//
	// Deprecated: In favour of workspaceFolders.
	RootURI *DocumentURI `json:"rootUri"`
	// User provided initialization options.
	InitializationOptions json.RawMessage `json:"initializationOptions,omitempty"`
	// The capabilities provided by the client (editor or tool)
	Capabilities ClientCapabilities `json:"capabilities"`
	// The initial trace setting. If omitted trace is disabled ('off').
	Trace string `json:"trace,omitempty"`
	// The workspace folders configured in the client when the server starts.
	//
	// This property is only available if the client supports workspace folders.
	// It can be null if the client supports workspace folders but none are
	// configured.
	//
	// Since 3.6.0
	WorkspaceFolders []WorkspaceFolder `json:"workspaceFolders,omitempty"`
}

// ClientInfo contains information about the client.
type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// InitializeResult is the result returned from an initialize request.
type InitializeResult struct {
	// The capabilities the language server provides.
	Capabilities ServerCapabilities `json:"capabilities"`
	// Information about the server.
	//
	// Since 3.15.0
	ServerInfo *ServerInfo `json:"serverInfo,omitempty"`
}

// ServerInfo contains information about the server.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// InitializedParams is sent after the client receives the initialize result; the struct is intentionally empty.
type InitializedParams struct{}

// ReferenceParams holds the parameters for a [ReferencesRequest].
type ReferenceParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
	Context ReferenceContext `json:"context"`
}

// ReferenceContext is the value-object that contains additional information when
// requesting references.
type ReferenceContext struct {
	// Include the declaration of the current symbol.
	IncludeDeclaration bool `json:"includeDeclaration"`
}

// DeclarationParams is sent to find the declaration of a symbol at a given position.
type DeclarationParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// DefinitionParams holds the parameters for a [DefinitionRequest].
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

// DocumentHighlightParams holds the parameters for a [DocumentHighlightRequest].
type DocumentHighlightParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
}

// DocumentHighlight is a range inside a text document which deserves
// special attention. Usually a document highlight is visualized by changing
// the background color of its range.
type DocumentHighlight struct {
	// The range this highlight applies to.
	Range Range `json:"range"`
	// The highlight kind, default is [DocumentHighlightKindText].
	Kind *DocumentHighlightKind `json:"kind,omitempty"`
}

// DocumentHighlightKind is an int enum: text (1), read-access (2), or write-access (3).
type DocumentHighlightKind int

const (
	// A textual occurrence.
	DocumentHighlightKindText DocumentHighlightKind = 1
	// Read-access of a symbol, like reading a variable.
	DocumentHighlightKindRead DocumentHighlightKind = 2
	// Write-access of a symbol, like writing to a variable.
	DocumentHighlightKindWrite DocumentHighlightKind = 3
)
