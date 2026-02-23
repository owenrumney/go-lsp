package lsp

// TextDocumentSyncKind defines how the host (editor) should sync document changes.
type TextDocumentSyncKind int

const (
	SyncNone        TextDocumentSyncKind = 0
	SyncFull        TextDocumentSyncKind = 1
	SyncIncremental TextDocumentSyncKind = 2
)

// FileChangeType represents the type of a file event.
type FileChangeType int

const (
	FileCreated FileChangeType = 1
	FileChanged FileChangeType = 2
	FileDeleted FileChangeType = 3
)

// WatchKind represents the kind of file events to watch.
type WatchKind int

const (
	WatchCreate WatchKind = 1
	WatchChange WatchKind = 2
	WatchDelete WatchKind = 4
)

// CompletionTriggerKind indicates how a completion was triggered.
type CompletionTriggerKind int

const (
	CompletionTriggerInvoked                         CompletionTriggerKind = 1
	CompletionTriggerCharacter                       CompletionTriggerKind = 2
	CompletionTriggerForIncompleteCompletions        CompletionTriggerKind = 3
)

// SignatureHelpTriggerKind indicates how signature help was triggered.
type SignatureHelpTriggerKind int

const (
	SignatureHelpTriggerInvoked         SignatureHelpTriggerKind = 1
	SignatureHelpTriggerCharacter       SignatureHelpTriggerKind = 2
	SignatureHelpTriggerContentChange   SignatureHelpTriggerKind = 3
)

// InsertTextFormat defines how an insert text is interpreted.
type InsertTextFormat int

const (
	InsertTextFormatPlainText InsertTextFormat = 1
	InsertTextFormatSnippet   InsertTextFormat = 2
)

// InsertTextMode indicates how whitespace and indentation is handled during completion.
type InsertTextMode int

const (
	InsertTextModeAsIs              InsertTextMode = 1
	InsertTextModeAdjustIndentation InsertTextMode = 2
)

// ResourceOperationKind represents the kind of resource operations.
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

// PrepareSupportDefaultBehavior represents the default behavior for prepare support.
type PrepareSupportDefaultBehavior int

const (
	PrepareSupportDefaultBehaviorIdentifier PrepareSupportDefaultBehavior = 1
)

// TokenFormat represents the token format.
type TokenFormat string

const (
	TokenFormatRelative TokenFormat = "relative"
)

// TraceValue represents the trace setting level.
type TraceValue string

const (
	TraceOff      TraceValue = "off"
	TraceMessages TraceValue = "messages"
	TraceVerbose  TraceValue = "verbose"
)
