package lsp

// CompletionItemKind is an int enum classifying completions (Function, Variable, Class, etc.) so the editor can show appropriate icons.
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
	// Render a completion as obsolete, usually using a strike-out.
	CompletionItemTagDeprecated CompletionItemTag = 1
)

// CompletionContext contains the additional information about the context in which a completion request is triggered.
type CompletionContext struct {
	// How the completion was triggered.
	TriggerKind CompletionTriggerKind `json:"triggerKind"`
	// The trigger character (a single character) that triggered code completion.
	// Is empty if `triggerKind !== [CompletionTriggerCharacter]`
	TriggerCharacter string `json:"triggerCharacter,omitempty"`
}

// CompletionParams holds the parameters.
type CompletionParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	PartialResultParams
	// The completion context. This is only available if the client specifies
	// to send this using the client capability `textDocument.completion.contextSupport === true`
	Context *CompletionContext `json:"context,omitempty"`
}

// CompletionItem represents a text snippet that is
// proposed to complete text that is being typed.
type CompletionItem struct {
	// The label of this completion item.
	//
	// The label property is also by default the text that
	// is inserted when selecting this completion.
	//
	// If label details are provided the label itself should
	// be an unqualified name of the completion item.
	Label string `json:"label"`
	// The kind of this completion item. Based on the kind
	// an icon is chosen by the editor.
	Kind *CompletionItemKind `json:"kind,omitempty"`
	// Tags for this completion item.
	//
	// Since 3.15.0
	Tags []CompletionItemTag `json:"tags,omitempty"`
	// A human-readable string with additional information
	// about this item, like type or symbol information.
	Detail string `json:"detail,omitempty"`
	// A human-readable string that represents a doc-comment.
	Documentation *MarkupContent `json:"documentation,omitempty"`
	// Indicates if this item is deprecated.
	//
	// Deprecated: Use tags instead.
	Deprecated *bool `json:"deprecated,omitempty"`
	// Select this item when showing.
	//
	// *Note* that only one completion item can be selected and that the
	// tool / client decides which item that is. The rule is that the *first*
	// item of those that match best is selected.
	Preselect *bool `json:"preselect,omitempty"`
	// A string that should be used when comparing this item
	// with other items. When empty the [CompletionItem.Label]
	// is used.
	SortText string `json:"sortText,omitempty"`
	// A string that should be used when filtering a set of
	// completion items. When empty the [CompletionItem.Label]
	// is used.
	FilterText string `json:"filterText,omitempty"`
	// A string that should be inserted into a document when selecting
	// this completion. When empty the [CompletionItem.Label]
	// is used.
	//
	// The insertText is subject to interpretation by the client side.
	// Some tools might not take the string literally. For example
	// VS Code when code complete is requested in this example
	// `con<cursor position>` and a completion item with an insertText of
	// console is provided it will only insert sole. Therefore it is
	// recommended to use textEdit instead since it avoids additional client
	// side interpretation.
	InsertText string `json:"insertText,omitempty"`
	// The format of the insert text. The format applies to both the
	// insertText property and the newText property of a provided
	// textEdit. If omitted defaults to `[InsertTextFormatPlainText]`.
	//
	// Please note that the insertTextFormat doesn't apply to
	// additionalTextEdits.
	InsertTextFormat *InsertTextFormat `json:"insertTextFormat,omitempty"`
	// How whitespace and indentation is handled during completion
	// item insertion. If not provided the clients default value depends on
	// the `textDocument.completion.insertTextMode` client capability.
	//
	// Since 3.16.0
	InsertTextMode *InsertTextMode `json:"insertTextMode,omitempty"`
	// An [TextEdit] which is applied to a document when selecting
	// this completion. When an edit is provided the value of
	// [CompletionItem.InsertText] is ignored.
	//
	// Most editors support two different operations when accepting a completion
	// item. One is to insert a completion text and the other is to replace an
	// existing text with a completion text. Since this can usually not be
	// predetermined by a server it can report both ranges. Clients need to
	// signal support for InsertReplaceEdits via the
	// `textDocument.completion.insertReplaceSupport` client capability
	// property.
	//
	// *Note 1:* The text edit's range as well as both ranges from an insert
	// replace edit must be a [single line] and they must contain the position
	// at which completion has been requested.
	// *Note 2:* If an InsertReplaceEdit is returned the edit's insert range
	// must be a prefix of the edit's replace range, that means it must be
	// contained and starting at the same position.
	//
	// Since 3.16.0 additional type InsertReplaceEdit
	TextEdit *TextEdit `json:"textEdit,omitempty"`
	// An optional array of additional [TextEdit] that are applied when
	// selecting this completion. Edits must not overlap (including the same insert position)
	// with the main [CompletionItem.TextEdit] nor with themselves.
	//
	// Additional text edits should be used to change text unrelated to the current cursor position
	// (for example adding an import statement at the top of the file if the completion item will
	// insert an unqualified type).
	AdditionalTextEdits []TextEdit `json:"additionalTextEdits,omitempty"`
	// An optional set of characters that when pressed while this completion is active will accept it first and
	// then type that character. *Note* that all commit characters should have `length=1` and that superfluous
	// characters will be ignored.
	CommitCharacters []string `json:"commitCharacters,omitempty"`
	// An optional [Command] that is executed *after* inserting this completion. *Note* that
	// additional modifications to the current document should be described with the
	// [CompletionItem.AdditionalTextEdits]-property.
	Command *Command `json:"command,omitempty"`
	// A data entry field that is preserved on a completion item between a
	// [CompletionRequest] and a [CompletionResolveRequest].
	Data any `json:"data,omitempty"`
}

// CompletionList represents a collection of [CompletionItem] to be presented
// in the editor.
type CompletionList struct {
	// This list is not complete. Further typing results in recomputing this list.
	//
	// Recomputed lists have all their items replaced (not appended) in the
	// incomplete completion sessions.
	IsIncomplete bool `json:"isIncomplete"`
	// The completion items.
	Items []CompletionItem `json:"items"`
}
