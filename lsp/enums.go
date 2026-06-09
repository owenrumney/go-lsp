package lsp

// TextDocumentSyncKind is an int enum controlling whether the editor sends no content (0), full content (1), or incremental diffs (2) on each change.
type TextDocumentSyncKind int

const (
	// Documents should not be synced at all.
	SyncNone TextDocumentSyncKind = 0
	// Documents are synced by always sending the full content
	// of the document.
	SyncFull TextDocumentSyncKind = 1
	// Documents are synced by sending the full content on open.
	// After that only incremental updates to the document are
	// send.
	SyncIncremental TextDocumentSyncKind = 2
)

// FileChangeType is an int enum: created (1), changed (2), or deleted (3).
type FileChangeType int

const (
	// The file got created.
	FileCreated FileChangeType = 1
	// The file got changed.
	FileChanged FileChangeType = 2
	// The file got deleted.
	FileDeleted FileChangeType = 3
)

// WatchKind is a bitmask for subscribing to file create (1), change (2), and/or delete (4) events.
type WatchKind int

const (
	// Interested in create events.
	WatchCreate WatchKind = 1
	// Interested in change events.
	WatchChange WatchKind = 2
	// Interested in delete events.
	WatchDelete WatchKind = 4
)

// CompletionTriggerKind is an int enum: invoked manually (1), by a trigger character (2), or for incomplete completions (3).
type CompletionTriggerKind int

const (
	// Completion was triggered by typing an identifier (24x7 code
	// complete), manual invocation (e.g Ctrl+Space) or via API.
	CompletionTriggerInvoked CompletionTriggerKind = 1
	// Completion was triggered by a trigger character specified by
	// the triggerCharacters properties of the CompletionRegistrationOptions.
	CompletionTriggerCharacter CompletionTriggerKind = 2
	// Completion was re-triggered as current completion list is incomplete.
	CompletionTriggerForIncompleteCompletions CompletionTriggerKind = 3
)

// SignatureHelpTriggerKind is an int enum: invoked manually (1), by a trigger character (2), or by cursor movement within a signature (3).
type SignatureHelpTriggerKind int

const (
	// Signature help was invoked manually by the user or by a command.
	SignatureHelpTriggerInvoked SignatureHelpTriggerKind = 1
	// Signature help was triggered by a trigger character.
	SignatureHelpTriggerCharacter SignatureHelpTriggerKind = 2
	// Signature help was triggered by the cursor moving or by the document content changing.
	SignatureHelpTriggerContentChange SignatureHelpTriggerKind = 3
)

// InsertTextFormat is an int enum: plain text (1) or a snippet with tab stops and placeholders (2).
type InsertTextFormat int

const (
	// The primary text to be inserted is treated as a plain string.
	InsertTextFormatPlainText InsertTextFormat = 1
	// The primary text to be inserted is treated as a snippet.
	//
	// A snippet can define tab stops and placeholders with `$1`, `$2`
	// and `${3:foo}`. `$0` defines the final tab stop, it defaults to
	// the end of the snippet. Placeholders with equal identifiers are linked,
	// that is typing in one will update others too.
	//
	// See also: https://microsoft.github.io/language-server-protocol/specifications/specification-current/#snippet_syntax
	InsertTextFormatSnippet InsertTextFormat = 2
)

// InsertTextMode is an int enum: asIs (1) keeps original whitespace, adjustIndentation (2) adapts leading whitespace to the insertion context.
type InsertTextMode int

const (
	// The insertion or replace strings is taken as it is. If the
	// value is multi line the lines below the cursor will be
	// inserted using the indentation defined in the string value.
	// The client will not apply any kind of adjustments to the
	// string.
	InsertTextModeAsIs InsertTextMode = 1
	// The editor adjusts leading whitespace of new lines so that
	// they match the indentation up to the cursor of the line for
	// which the item is accepted.
	//
	// Consider a line like this: <2tabs><cursor><3tabs>foo. Accepting a
	// multi line completion item is indented using 2 tabs and all
	// following lines inserted will be indented using 2 tabs as well.
	InsertTextModeAdjustIndentation InsertTextMode = 2
)

// ResourceOperationKind is a string enum: "create", "rename", or "delete" for workspace edit file operations.
type ResourceOperationKind string

const (
	// Supports creating new files and folders.
	ResourceOperationCreate ResourceOperationKind = "create"
	// Supports renaming existing files and folders.
	ResourceOperationRename ResourceOperationKind = "rename"
	// Supports deleting existing files and folders.
	ResourceOperationDelete ResourceOperationKind = "delete"
)

// FailureHandlingKind represents how the client should handle failures.
type FailureHandlingKind string

const (
	// Applying the workspace change is simply aborted if one of the changes provided
	// fails. All operations executed before the failing operation stay executed.
	FailureHandlingAbort FailureHandlingKind = "abort"
	// All operations are executed transactional. That means they either all
	// succeed or no changes at all are applied to the workspace.
	FailureHandlingTransactional FailureHandlingKind = "transactional"
	// The client tries to undo the operations already executed. But there is no
	// guarantee that this is succeeding.
	FailureHandlingUndo FailureHandlingKind = "undo"
	// If the workspace edit contains only textual file changes they are executed transactional.
	// If resource changes (create, rename or delete file) are part of the change the failure
	// handling strategy is abort.
	FailureHandlingTextOnlyTransactional FailureHandlingKind = "textOnlyTransactional"
)

// PrepareSupportDefaultBehavior is an int enum defining fallback behavior when the server doesn't support prepareRename (1 = use identifier under cursor).
type PrepareSupportDefaultBehavior int

const (
	// The client's default behavior is to select the identifier
	// according to the language's syntax rule.
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
