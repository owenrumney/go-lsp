package lsp

// CodeActionKind represents the kind of a code action.
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

// CodeActionParams contains the params for textDocument/codeAction.
type CodeActionParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	Range        Range                  `json:"range"`
	Context      CodeActionContext      `json:"context"`
}

// CodeAction represents a change that can be performed in code.
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
