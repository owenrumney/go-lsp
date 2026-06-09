package lsp

import "encoding/json"

// WorkspaceFolder is a folder opened inside a client.
type WorkspaceFolder struct {
	// The associated URI for this workspace folder.
	URI DocumentURI `json:"uri"`
	// The name of the workspace folder. Used to refer to this
	// workspace folder in the user interface.
	Name string `json:"name"`
}

// WorkspaceEdit represents changes to many resources managed in the workspace. The edit
// should either provide changes or documentChanges. If documentChanges are present
// they are preferred over changes if the client can handle versioned document edits.
//
// Since version 3.13.0 a workspace edit can contain resource operations as well. If resource
// operations are present clients need to execute the operations in the order in which they
// are provided. So a workspace edit for example can consist of the following two changes:
// (1) a create file a.txt and (2) a text document edit which insert text into file a.txt.
//
// An invalid sequence (e.g. (1) delete file a.txt and (2) insert text into file a.txt) will
// cause failure of the operation. How the client recovers from the failure is described by
// the client capability: `workspace.workspaceEdit.failureHandling`.
type WorkspaceEdit struct {
	// Holds changes to existing resources.
	Changes map[DocumentURI][]TextEdit `json:"changes,omitempty"`
	// Depending on the client capability `workspace.workspaceEdit.resourceOperations` document changes
	// are either an array of TextDocumentEdits to express changes to n different text documents
	// where each text document edit addresses a specific version of a text document. Or it can contain
	// above TextDocumentEdits mixed with create, rename and delete file / folder operations.
	//
	// Whether a client supports versioned document edits is expressed via
	// `workspace.workspaceEdit.documentChanges` client capability.
	//
	// If a client neither supports documentChanges nor `workspace.workspaceEdit.resourceOperations` then
	// only plain TextEdits using the changes property are supported.
	DocumentChanges []TextDocumentEdit `json:"documentChanges,omitempty"`
	// A map of change annotations that can be referenced in AnnotatedTextEdits or create, rename and
	// delete file / folder operations.
	//
	// Whether clients honor this property depends on the client capability `workspace.changeAnnotationSupport`.
	//
	// Since 3.16.0
	ChangeAnnotations map[ChangeAnnotationIdentifier]ChangeAnnotation `json:"changeAnnotations,omitempty"`
}

// TextDocumentEdit describes textual changes on a text document. A TextDocumentEdit describes all changes
// on a document version Si and after they are applied move the document to version Si+1.
// So the creator of a TextDocumentEdit doesn't need to sort the array of edits or do any
// kind of ordering. However the edits must be non overlapping.
type TextDocumentEdit struct {
	// The text document to change.
	TextDocument OptionalVersionedTextDocumentIdentifier `json:"textDocument"`
	// The edits to be applied.
	//
	// Since 3.16.0 - support for AnnotatedTextEdit. This is guarded using a
	// client capability.
	Edits []TextEdit `json:"edits"`
}

// CreateFileOptions is used to create a file.
type CreateFileOptions struct {
	// Overwrite existing file. Overwrite wins over ignoreIfExists
	Overwrite *bool `json:"overwrite,omitempty"`
	// Ignore if exists.
	IgnoreIfExists *bool `json:"ignoreIfExists,omitempty"`
}

// CreateFile is an operation.
type CreateFile struct {
	// A create operation.
	Kind string `json:"kind"` // "create"
	// The resource to create.
	URI DocumentURI `json:"uri"`
	// Additional options
	Options *CreateFileOptions `json:"options,omitempty"`
}

// RenameFileOptions configures a rename-file workspace edit.
type RenameFileOptions struct {
	// Overwrite target if existing. Overwrite wins over ignoreIfExists
	Overwrite *bool `json:"overwrite,omitempty"`
	// Ignores if target exists.
	IgnoreIfExists *bool `json:"ignoreIfExists,omitempty"`
}

// RenameFile is an operation.
type RenameFile struct {
	// A rename operation.
	Kind string `json:"kind"` // "rename"
	// The old (existing) location.
	OldURI DocumentURI `json:"oldUri"`
	// The new location.
	NewURI DocumentURI `json:"newUri"`
	// Rename options.
	Options *RenameFileOptions `json:"options,omitempty"`
}

// DeleteFileOptions configures a delete-file workspace edit.
type DeleteFileOptions struct {
	// Delete the content recursively if a folder is denoted.
	Recursive *bool `json:"recursive,omitempty"`
	// Ignore the operation if the file doesn't exist.
	IgnoreIfNotExists *bool `json:"ignoreIfNotExists,omitempty"`
}

// DeleteFile is an operation.
type DeleteFile struct {
	// A delete operation.
	Kind string `json:"kind"` // "delete"
	// The file to delete.
	URI DocumentURI `json:"uri"`
	// Delete options.
	Options *DeleteFileOptions `json:"options,omitempty"`
}

// FileEvent is an event describing a file change.
type FileEvent struct {
	// The file's uri.
	URI DocumentURI `json:"uri"`
	// The change type.
	Type FileChangeType `json:"type"`
}

// DidChangeWatchedFilesParams holds the watched files change notification's parameters.
type DidChangeWatchedFilesParams struct {
	// The actual file events.
	Changes []FileEvent `json:"changes"`
}

// FileSystemWatcher describes a file system watcher.
type FileSystemWatcher struct {
	// The glob pattern to watch. See [GlobPattern] for more detail.
	//
	// Since 3.17.0 support for relative patterns.
	GlobPattern string `json:"globPattern"`
	// The kind of events of interest. If omitted it defaults
	// to [WatchCreate] | [WatchChange] | [WatchDelete]
	// which is 7.
	Kind *WatchKind `json:"kind,omitempty"`
}

// DidChangeConfigurationParams holds the parameters of a change configuration notification.
type DidChangeConfigurationParams struct {
	// The actual changed settings
	Settings any `json:"settings"`
}

// ConfigurationParams holds the parameters of a configuration request.
type ConfigurationParams struct {
	Items []ConfigurationItem `json:"items"`
}

// ConfigurationItem identifies a configuration section to fetch, optionally scoped to a resource URI.
type ConfigurationItem struct {
	// The scope to get the configuration section for.
	ScopeURI *DocumentURI `json:"scopeUri,omitempty"`
	// The configuration section asked for.
	Section string `json:"section,omitempty"`
}

// DidChangeWorkspaceFoldersParams holds the parameters of a `workspace/didChangeWorkspaceFolders` notification.
type DidChangeWorkspaceFoldersParams struct {
	// The actual workspace folder change event.
	Event WorkspaceFoldersChangeEvent `json:"event"`
}

// WorkspaceFoldersChangeEvent is the workspace folder change event.
type WorkspaceFoldersChangeEvent struct {
	// The array of added workspace folders
	Added []WorkspaceFolder `json:"added"`
	// The array of the removed workspace folders
	Removed []WorkspaceFolder `json:"removed"`
}

// ExecuteCommandParams holds the parameters of a [ExecuteCommandRequest].
type ExecuteCommandParams struct {
	WorkDoneProgressParams
	// The identifier of the actual command handler.
	Command string `json:"command"`
	// Arguments that the command should be invoked with.
	Arguments []json.RawMessage `json:"arguments,omitempty"`
}

// FileCreate represents information on a file/folder create.
//
// Since 3.16.0.
type FileCreate struct {
	// A file:// URI for the location of the file/folder being created.
	URI string `json:"uri"`
}

// CreateFilesParams holds the parameters sent in notifications/requests for user-initiated creation of
// files.
//
// Since 3.16.0.
type CreateFilesParams struct {
	// An array of all files/folders created in this operation.
	Files []FileCreate `json:"files"`
}

// FileRename represents information on a file/folder rename.
//
// Since 3.16.0.
type FileRename struct {
	// A file:// URI for the original location of the file/folder being renamed.
	OldURI string `json:"oldUri"`
	// A file:// URI for the new location of the file/folder being renamed.
	NewURI string `json:"newUri"`
}

// RenameFilesParams holds the parameters sent in notifications/requests for user-initiated renames of
// files.
//
// Since 3.16.0.
type RenameFilesParams struct {
	// An array of all files/folders renamed in this operation. When a folder is renamed, only
	// the folder will be included, and not its children.
	Files []FileRename `json:"files"`
}

// FileDelete represents information on a file/folder delete.
//
// Since 3.16.0.
type FileDelete struct {
	// A file:// URI for the location of the file/folder being deleted.
	URI string `json:"uri"`
}

// DeleteFilesParams holds the parameters sent in notifications/requests for user-initiated deletes of
// files.
//
// Since 3.16.0.
type DeleteFilesParams struct {
	// An array of all files/folders deleted in this operation.
	Files []FileDelete `json:"files"`
}

// ApplyWorkspaceEditParams holds the parameters passed via an apply workspace edit request.
type ApplyWorkspaceEditParams struct {
	// An optional label of the workspace edit. This label is
	// presented in the user interface for example on an undo
	// stack to undo the workspace edit.
	Label string `json:"label,omitempty"`
	// The edits to apply.
	Edit WorkspaceEdit `json:"edit"`
}

// ApplyWorkspaceEditResult is the result returned from the apply workspace edit request.
//
// Since 3.17 renamed from ApplyWorkspaceEditResponse.
type ApplyWorkspaceEditResult struct {
	// Indicates whether the edit was applied or not.
	Applied bool `json:"applied"`
	// An optional textual description for why the edit was not applied.
	// This may be used by the server for diagnostic logging or to provide
	// a suitable error for a request that triggered the edit.
	FailureReason string `json:"failureReason,omitempty"`
}
