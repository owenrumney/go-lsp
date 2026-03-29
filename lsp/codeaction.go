package lsp

// CodeActionKind is a string enum classifying code actions (e.g. "quickfix", "refactor.extract", "source.organizeImports").
type CodeActionKind string

const (
	CodeActionQuickFix              CodeActionKind = "quickfix"
	CodeActionRefactor              CodeActionKind = "refactor"
	CodeActionRefactorExtract       CodeActionKind = "refactor.extract"
	CodeActionRefactorInline        CodeActionKind = "refactor.inline"
	CodeActionRefactorRewrite       CodeActionKind = "refactor.rewrite"
	CodeActionSource                CodeActionKind = "source"
	CodeActionSourceOrganizeImports CodeActionKind = "source.organizeImports"
)

// CodeActionContext contains additional diagnostic information about the context of a code action.
type CodeActionContext struct {
	Diagnostics []Diagnostic     `json:"diagnostics"`
	Only        []CodeActionKind `json:"only,omitempty"`
}

// CodeActionParams is sent to request available fixes and refactorings for diagnostics or a selected range.
type CodeActionParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

// CodeAction is a fix, refactoring, or source action the server offers, optionally carrying a workspace edit and/or a command.
type CodeAction struct {
	Title       string          `json:"title"`
	Kind        *CodeActionKind `json:"kind,omitempty"`
	Diagnostics []Diagnostic    `json:"diagnostics,omitempty"`
	IsPreferred *bool           `json:"isPreferred,omitempty"`
	Disabled    *struct {
		Reason string `json:"reason"`
	} `json:"disabled,omitempty"`
	Edit    *WorkspaceEdit `json:"edit,omitempty"`
	Command *Command       `json:"command,omitempty"`
	Data    any            `json:"data,omitempty"`
}
