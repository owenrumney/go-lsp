package lsp

// CodeActionKind is a string enum classifying code actions (e.g. "quickfix", "refactor.extract", "source.organizeImports").
type CodeActionKind string

const (
	// Base kind for quickfix actions: 'quickfix'.
	CodeActionQuickFix CodeActionKind = "quickfix"
	// Base kind for refactoring actions: 'refactor'.
	CodeActionRefactor CodeActionKind = "refactor"
	// Base kind for refactoring extraction actions: 'refactor.extract'
	//
	// Example extract actions:
	//
	// - Extract method
	// - Extract function
	// - Extract variable
	// - Extract interface from class
	// - ...
	CodeActionRefactorExtract CodeActionKind = "refactor.extract"
	// Base kind for refactoring inline actions: 'refactor.inline'
	//
	// Example inline actions:
	//
	// - Inline function
	// - Inline variable
	// - Inline constant
	// - ...
	CodeActionRefactorInline CodeActionKind = "refactor.inline"
	// Base kind for refactoring rewrite actions: 'refactor.rewrite'
	//
	// Example rewrite actions:
	//
	// - Convert JavaScript function to class
	// - Add or remove parameter
	// - Encapsulate field
	// - Make method static
	// - Move method to base class
	// - ...
	CodeActionRefactorRewrite CodeActionKind = "refactor.rewrite"
	// Base kind for source actions: source
	//
	// Source code actions apply to the entire file.
	CodeActionSource CodeActionKind = "source"
	// Base kind for an organize imports source action: `source.organizeImports`.
	CodeActionSourceOrganizeImports CodeActionKind = "source.organizeImports"
)

// CodeActionContext contains additional diagnostic information about the context in which
// a [CodeActionProvider.ProvideCodeActions] is run.
type CodeActionContext struct {
	// An array of diagnostics known on the client side overlapping the range provided to the
	// `textDocument/codeAction` request. They are provided so that the server knows which
	// errors are currently presented to the user for the given range. There is no guarantee
	// that these accurately reflect the error state of the resource. The primary parameter
	// to compute code actions is the provided range.
	Diagnostics []Diagnostic `json:"diagnostics"`
	// Requested kind of actions to return.
	//
	// Actions not of this kind are filtered out by the client before being shown. So servers
	// can omit computing them.
	Only []CodeActionKind `json:"only,omitempty"`
}

// CodeActionParams holds the parameters of a [CodeActionRequest].
type CodeActionParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The document in which the command was invoked.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The range for which the command was invoked.
	Range Range `json:"range"`
	// Context carrying additional information.
	Context CodeActionContext `json:"context"`
}

// CodeAction represents a change that can be performed in code, e.g. to fix a problem or
// to refactor code.
//
// A CodeAction must set either edit and/or a command. If both are supplied, the edit is applied first, then the command is executed.
type CodeAction struct {
	// A short, human-readable, title for this code action.
	Title string `json:"title"`
	// The kind of the code action.
	//
	// Used to filter code actions.
	Kind *CodeActionKind `json:"kind,omitempty"`
	// The diagnostics that this code action resolves.
	Diagnostics []Diagnostic `json:"diagnostics,omitempty"`
	// Marks this as a preferred action. Preferred actions are used by the `auto fix` command and can be targeted
	// by keybindings.
	//
	// A quick fix should be marked preferred if it properly addresses the underlying error.
	// A refactoring should be marked preferred if it is the most reasonable choice of actions to take.
	//
	// Since 3.15.0
	IsPreferred *bool `json:"isPreferred,omitempty"`
	// Marks that the code action cannot currently be applied.
	//
	// Clients should follow the following guidelines regarding disabled code actions:
	//
	//   - Disabled code actions are not shown in automatic [lightbulbs]
	//     code action menus.
	//
	//   - Disabled actions are shown as faded out in the code action menu when the user requests a more specific type
	//     of code action, such as refactorings.
	//
	//   - If the user has a [keybinding]
	//     that auto applies a code action and only disabled code actions are returned, the client should show the user an
	//     error message with reason in the editor.
	//
	// Since 3.16.0
	//
	// [lightbulbs]: https://code.visualstudio.com/docs/editor/editingevolved#_code-action
	// [keybinding]: https://code.visualstudio.com/docs/editor/refactoring#_keybindings-for-code-actions
	Disabled *struct {
		Reason string `json:"reason"`
	} `json:"disabled,omitempty"`
	// The workspace edit this code action performs.
	Edit *WorkspaceEdit `json:"edit,omitempty"`
	// A command this code action executes. If a code action
	// provides an edit and a command, first the edit is
	// executed and then the command.
	Command *Command `json:"command,omitempty"`
	// A data entry field that is preserved on a code action between
	// a `textDocument/codeAction` and a `codeAction/resolve` request.
	//
	// Since 3.16.0
	Data any `json:"data,omitempty"`
}
