package lsp

import "encoding/json"

// WorkspaceFolder represents a workspace folder.
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

// CreateFile represents a create file operation.
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

// RenameFile represents a rename file operation.
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

// DeleteFile represents a delete file operation.
type DeleteFile struct {
	Kind    string             `json:"kind"` // "delete"
	URI     DocumentURI        `json:"uri"`
	Options *DeleteFileOptions `json:"options,omitempty"`
}

// FileEvent represents an event describing a file change.
type FileEvent struct {
	URI  DocumentURI    `json:"uri"`
	Type FileChangeType `json:"type"`
}

// DidChangeWatchedFilesParams contains the params for workspace/didChangeWatchedFiles.
type DidChangeWatchedFilesParams struct {
	Changes []FileEvent `json:"changes"`
}

// FileSystemWatcher describes a file system watcher.
type FileSystemWatcher struct {
	GlobPattern string     `json:"globPattern"`
	Kind        *WatchKind `json:"kind,omitempty"`
}

// DidChangeConfigurationParams contains the params for workspace/didChangeConfiguration.
type DidChangeConfigurationParams struct {
	Settings any `json:"settings"`
}

// ConfigurationParams contains the params for workspace/configuration.
type ConfigurationParams struct {
	Items []ConfigurationItem `json:"items"`
}

// ConfigurationItem represents a configuration item.
type ConfigurationItem struct {
	ScopeURI *DocumentURI `json:"scopeUri,omitempty"`
	Section  string       `json:"section,omitempty"`
}

// DidChangeWorkspaceFoldersParams contains the params for workspace/didChangeWorkspaceFolders.
type DidChangeWorkspaceFoldersParams struct {
	Event WorkspaceFoldersChangeEvent `json:"event"`
}

// WorkspaceFoldersChangeEvent describes workspace folder change events.
type WorkspaceFoldersChangeEvent struct {
	Added   []WorkspaceFolder `json:"added"`
	Removed []WorkspaceFolder `json:"removed"`
}

// ExecuteCommandParams contains the params for workspace/executeCommand.
type ExecuteCommandParams struct {
	WorkDoneProgressParams
	Command   string            `json:"command"`
	Arguments []json.RawMessage `json:"arguments,omitempty"`
}

// FileCreate represents a file that was created.
type FileCreate struct {
	URI string `json:"uri"`
}

// CreateFilesParams contains the params for workspace/willCreateFiles.
type CreateFilesParams struct {
	Files []FileCreate `json:"files"`
}

// FileRename represents a file that was renamed.
type FileRename struct {
	OldURI string `json:"oldUri"`
	NewURI string `json:"newUri"`
}

// RenameFilesParams contains the params for workspace/willRenameFiles.
type RenameFilesParams struct {
	Files []FileRename `json:"files"`
}

// FileDelete represents a file that was deleted.
type FileDelete struct {
	URI string `json:"uri"`
}

// DeleteFilesParams contains the params for workspace/willDeleteFiles.
type DeleteFilesParams struct {
	Files []FileDelete `json:"files"`
}

// ApplyWorkspaceEditParams contains the params for workspace/applyEdit.
type ApplyWorkspaceEditParams struct {
	Label string        `json:"label,omitempty"`
	Edit  WorkspaceEdit `json:"edit"`
}

// ApplyWorkspaceEditResult contains the result for workspace/applyEdit.
type ApplyWorkspaceEditResult struct {
	Applied       bool   `json:"applied"`
	FailureReason string `json:"failureReason,omitempty"`
}
