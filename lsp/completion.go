package lsp

// CompletionItemKind represents the kind of a completion item.
type CompletionItemKind int

const (
	CompletionItemKindText          CompletionItemKind = 1
	CompletionItemKindMethod        CompletionItemKind = 2
	CompletionItemKindFunction      CompletionItemKind = 3
	CompletionItemKindConstructor   CompletionItemKind = 4
	CompletionItemKindField         CompletionItemKind = 5
	CompletionItemKindVariable      CompletionItemKind = 6
	CompletionItemKindClass         CompletionItemKind = 7
	CompletionItemKindInterface     CompletionItemKind = 8
	CompletionItemKindModule        CompletionItemKind = 9
	CompletionItemKindProperty      CompletionItemKind = 10
	CompletionItemKindUnit          CompletionItemKind = 11
	CompletionItemKindValue         CompletionItemKind = 12
	CompletionItemKindEnum          CompletionItemKind = 13
	CompletionItemKindKeyword       CompletionItemKind = 14
	CompletionItemKindSnippet       CompletionItemKind = 15
	CompletionItemKindColor         CompletionItemKind = 16
	CompletionItemKindFile          CompletionItemKind = 17
	CompletionItemKindReference     CompletionItemKind = 18
	CompletionItemKindFolder        CompletionItemKind = 19
	CompletionItemKindEnumMember    CompletionItemKind = 20
	CompletionItemKindConstant      CompletionItemKind = 21
	CompletionItemKindStruct        CompletionItemKind = 22
	CompletionItemKindEvent         CompletionItemKind = 23
	CompletionItemKindOperator      CompletionItemKind = 24
	CompletionItemKindTypeParameter CompletionItemKind = 25
)

// CompletionItemTag represents extra annotations for a completion item.
type CompletionItemTag int

const (
	CompletionItemTagDeprecated CompletionItemTag = 1
)

// CompletionContext contains additional information about the completion request.
type CompletionContext struct {
	TriggerKind      CompletionTriggerKind `json:"triggerKind"`
	TriggerCharacter string                `json:"triggerCharacter,omitempty"`
}

// CompletionParams contains the params for textDocument/completion.
type CompletionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
	Context *CompletionContext `json:"context,omitempty"`
}

// CompletionItem represents a completion suggestion.
type CompletionItem struct {
	Label               string              `json:"label"`
	Kind                *CompletionItemKind `json:"kind,omitempty"`
	Tags                []CompletionItemTag `json:"tags,omitempty"`
	Detail              string              `json:"detail,omitempty"`
	Documentation       *MarkupContent      `json:"documentation,omitempty"`
	Deprecated          *bool               `json:"deprecated,omitempty"`
	Preselect           *bool               `json:"preselect,omitempty"`
	SortText            string              `json:"sortText,omitempty"`
	FilterText          string              `json:"filterText,omitempty"`
	InsertText          string              `json:"insertText,omitempty"`
	InsertTextFormat    *InsertTextFormat    `json:"insertTextFormat,omitempty"`
	InsertTextMode      *InsertTextMode      `json:"insertTextMode,omitempty"`
	TextEdit            *TextEdit            `json:"textEdit,omitempty"`
	AdditionalTextEdits []TextEdit           `json:"additionalTextEdits,omitempty"`
	CommitCharacters    []string             `json:"commitCharacters,omitempty"`
	Command             *Command             `json:"command,omitempty"`
	Data                any                  `json:"data,omitempty"`
}

// CompletionList represents a collection of completion items.
type CompletionList struct {
	IsIncomplete bool             `json:"isIncomplete"`
	Items        []CompletionItem `json:"items"`
}
