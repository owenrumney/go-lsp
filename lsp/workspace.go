package lsp

import "encoding/json"

// WorkspaceFolder is a root folder in a multi-root workspace, identified by URI and a display name.
type WorkspaceFolder struct {
	URI  DocumentURI `json:"uri"`
	Name string      `json:"name"`
}

// WorkspaceEdit represents changes to many resources managed in the workspace.
type WorkspaceEdit struct {
	Changes           map[DocumentURI][]TextEdit                      `json:"changes,omitempty"`
	DocumentChanges   []TextDocumentEdit                              `json:"documentChanges,omitempty"`
	ChangeAnnotations map[ChangeAnnotationIdentifier]ChangeAnnotation `json:"changeAnnotations,omitempty"`
}

// TextDocumentEdit represents edits to a single text document.
type TextDocumentEdit struct {
	TextDocument OptionalVersionedTextDocumentIdentifier `json:"textDocument"`
	Edits        []TextEdit                              `json:"edits"`
}

// CreateFileOptions represents options for creating a file.
type CreateFileOptions struct {
	Overwrite      *bool `json:"overwrite,omitempty"`
	IgnoreIfExists *bool `json:"ignoreIfExists,omitempty"`
}

// CreateFile is a workspace edit operation that creates a new file, with optional overwrite/ignoreIfExists flags.
type CreateFile struct {
	Kind    string             `json:"kind"` // "create"
	URI     DocumentURI        `json:"uri"`
	Options *CreateFileOptions `json:"options,omitempty"`
}

// RenameFileOptions represents options for renaming a file.
type RenameFileOptions struct {
	Overwrite      *bool `json:"overwrite,omitempty"`
	IgnoreIfExists *bool `json:"ignoreIfExists,omitempty"`
}

// RenameFile is a workspace edit operation that renames/moves a file, with optional overwrite/ignoreIfExists flags.
type RenameFile struct {
	Kind    string             `json:"kind"` // "rename"
	OldURI  DocumentURI        `json:"oldUri"`
	NewURI  DocumentURI        `json:"newUri"`
	Options *RenameFileOptions `json:"options,omitempty"`
}

// DeleteFileOptions represents options for deleting a file.
type DeleteFileOptions struct {
	Recursive         *bool `json:"recursive,omitempty"`
	IgnoreIfNotExists *bool `json:"ignoreIfNotExists,omitempty"`
}

// DeleteFile is a workspace edit operation that deletes a file, with optional recursive and ignoreIfNotExists flags.
type DeleteFile struct {
	Kind    string             `json:"kind"` // "delete"
	URI     DocumentURI        `json:"uri"`
	Options *DeleteFileOptions `json:"options,omitempty"`
}

// FileEvent pairs a URI with a change type (created/changed/deleted) in a didChangeWatchedFiles notification.
type FileEvent struct {
	URI  DocumentURI    `json:"uri"`
	Type FileChangeType `json:"type"`
}

// DidChangeWatchedFilesParams notifies the server that watched files have been created, changed, or deleted.
type DidChangeWatchedFilesParams struct {
	Changes []FileEvent `json:"changes"`
}

// FileSystemWatcher describes a file system watcher.
type FileSystemWatcher struct {
	GlobPattern string     `json:"globPattern"`
	Kind        *WatchKind `json:"kind,omitempty"`
}

// DidChangeConfigurationParams notifies the server that the client's configuration settings have changed.
type DidChangeConfigurationParams struct {
	Settings any `json:"settings"`
}

// ConfigurationParams is sent from server to client to fetch configuration values for one or more scopes.
type ConfigurationParams struct {
	Items []ConfigurationItem `json:"items"`
}

// ConfigurationItem identifies a configuration section to fetch, optionally scoped to a resource URI.
type ConfigurationItem struct {
	ScopeURI *DocumentURI `json:"scopeUri,omitempty"`
	Section  string       `json:"section,omitempty"`
}

// DidChangeWorkspaceFoldersParams notifies the server that workspace folders were added or removed.
type DidChangeWorkspaceFoldersParams struct {
	Event WorkspaceFoldersChangeEvent `json:"event"`
}

// WorkspaceFoldersChangeEvent describes workspace folder change events.
type WorkspaceFoldersChangeEvent struct {
	Added   []WorkspaceFolder `json:"added"`
	Removed []WorkspaceFolder `json:"removed"`
}

// ExecuteCommandParams is sent to ask the server to run a registered command by ID with the given arguments.
type ExecuteCommandParams struct {
	WorkDoneProgressParams
	Command   string            `json:"command"`
	Arguments []json.RawMessage `json:"arguments,omitempty"`
}

// FileCreate carries the URI of a file that was created, for willCreateFiles/didCreateFiles notifications.
type FileCreate struct {
	URI string `json:"uri"`
}

// CreateFilesParams notifies the server that files are about to be created, allowing it to return a workspace edit.
type CreateFilesParams struct {
	Files []FileCreate `json:"files"`
}

// FileRename carries the old and new URIs of a renamed file.
type FileRename struct {
	OldURI string `json:"oldUri"`
	NewURI string `json:"newUri"`
}

// RenameFilesParams notifies the server that files are about to be renamed, allowing it to return a workspace edit (e.g. import path updates).
type RenameFilesParams struct {
	Files []FileRename `json:"files"`
}

// FileDelete carries the URI of a file that was deleted.
type FileDelete struct {
	URI string `json:"uri"`
}

// DeleteFilesParams notifies the server that files are about to be deleted, allowing it to return a workspace edit.
type DeleteFilesParams struct {
	Files []FileDelete `json:"files"`
}

// ApplyWorkspaceEditParams is sent from server to client to apply a set of text edits and file operations across the workspace.
type ApplyWorkspaceEditParams struct {
	Label string        `json:"label,omitempty"`
	Edit  WorkspaceEdit `json:"edit"`
}

// ApplyWorkspaceEditResult contains the result for workspace/applyEdit.
type ApplyWorkspaceEditResult struct {
	Applied       bool   `json:"applied"`
	FailureReason string `json:"failureReason,omitempty"`
}
