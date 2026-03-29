package lsp

// TextDocumentSyncKind is an int enum controlling whether the editor sends no content (0), full content (1), or incremental diffs (2) on each change.
type TextDocumentSyncKind int

const (
	SyncNone        TextDocumentSyncKind = 0
	SyncFull        TextDocumentSyncKind = 1
	SyncIncremental TextDocumentSyncKind = 2
)

// FileChangeType is an int enum: created (1), changed (2), or deleted (3).
type FileChangeType int

const (
	FileCreated FileChangeType = 1
	FileChanged FileChangeType = 2
	FileDeleted FileChangeType = 3
)

// WatchKind is a bitmask for subscribing to file create (1), change (2), and/or delete (4) events.
type WatchKind int

const (
	WatchCreate WatchKind = 1
	WatchChange WatchKind = 2
	WatchDelete WatchKind = 4
)

// CompletionTriggerKind is an int enum: invoked manually (1), by a trigger character (2), or for incomplete completions (3).
type CompletionTriggerKind int

const (
	CompletionTriggerInvoked                  CompletionTriggerKind = 1
	CompletionTriggerCharacter                CompletionTriggerKind = 2
	CompletionTriggerForIncompleteCompletions CompletionTriggerKind = 3
)

// SignatureHelpTriggerKind is an int enum: invoked manually (1), by a trigger character (2), or by cursor movement within a signature (3).
type SignatureHelpTriggerKind int

const (
	SignatureHelpTriggerInvoked       SignatureHelpTriggerKind = 1
	SignatureHelpTriggerCharacter     SignatureHelpTriggerKind = 2
	SignatureHelpTriggerContentChange SignatureHelpTriggerKind = 3
)

// InsertTextFormat is an int enum: plain text (1) or a snippet with tab stops and placeholders (2).
type InsertTextFormat int

const (
	InsertTextFormatPlainText InsertTextFormat = 1
	InsertTextFormatSnippet   InsertTextFormat = 2
)

// InsertTextMode is an int enum: asIs (1) keeps original whitespace, adjustIndentation (2) adapts leading whitespace to the insertion context.
type InsertTextMode int

const (
	InsertTextModeAsIs              InsertTextMode = 1
	InsertTextModeAdjustIndentation InsertTextMode = 2
)

// ResourceOperationKind is a string enum: "create", "rename", or "delete" for workspace edit file operations.
type ResourceOperationKind string

const (
	ResourceOperationCreate ResourceOperationKind = "create"
	ResourceOperationRename ResourceOperationKind = "rename"
	ResourceOperationDelete ResourceOperationKind = "delete"
)

// FailureHandlingKind represents how the client should handle failures.
type FailureHandlingKind string

const (
	FailureHandlingAbort                 FailureHandlingKind = "abort"
	FailureHandlingTransactional         FailureHandlingKind = "transactional"
	FailureHandlingUndo                  FailureHandlingKind = "undo"
	FailureHandlingTextOnlyTransactional FailureHandlingKind = "textOnlyTransactional"
)

// PrepareSupportDefaultBehavior is an int enum defining fallback behavior when the server doesn't support prepareRename (1 = use identifier under cursor).
type PrepareSupportDefaultBehavior int

const (
	PrepareSupportDefaultBehaviorIdentifier PrepareSupportDefaultBehavior = 1
)

// TokenFormat is a string enum for semantic token encoding; currently only "relative" is defined.
type TokenFormat string

const (
	TokenFormatRelative TokenFormat = "relative"
)

// TraceValue is a string enum ("off", "messages", "verbose") controlling LSP message tracing.
type TraceValue string

const (
	TraceOff      TraceValue = "off"
	TraceMessages TraceValue = "messages"
	TraceVerbose  TraceValue = "verbose"
)
