package lsp

import "encoding/json"

// PositionEncodingKind is a string enum naming how positions are encoded.
type PositionEncodingKind string

const (
	// Character offsets count UTF-8 code units (e.g. bytes).
	PositionEncodingUTF8 PositionEncodingKind = "utf-8"
	// Character offsets count UTF-16 code units.
	//
	// This is the default and must always be supported
	// by servers.
	PositionEncodingUTF16 PositionEncodingKind = "utf-16"
	// Character offsets count UTF-32 code units.
	//
	// Implementation note: these are the same as Unicode codepoints,
	// so this PositionEncodingKind may also be used for an
	// encoding-agnostic representation of character offsets.
	PositionEncodingUTF32 PositionEncodingKind = "utf-32"
)

// ClientCapabilities defines the capabilities provided by the client.
type ClientCapabilities struct {
	// Workspace specific client capabilities.
	Workspace *WorkspaceClientCapabilities `json:"workspace,omitempty"`
	// Text document specific client capabilities.
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	// Window specific client capabilities.
	Window *WindowClientCapabilities `json:"window,omitempty"`
	// General client capabilities.
	//
	// Since 3.16.0
	General *GeneralClientCapabilities `json:"general,omitempty"`
	// Experimental client capabilities.
	Experimental json.RawMessage `json:"experimental,omitempty"`
}

// WorkspaceClientCapabilities declares workspace-specific client capabilities.
type WorkspaceClientCapabilities struct {
	// The client supports applying batch edits
	// to the workspace by supporting the request
	// 'workspace/applyEdit'
	ApplyEdit *bool `json:"applyEdit,omitempty"`
	// Capabilities specific to WorkspaceEdits.
	WorkspaceEdit *WorkspaceEditClientCapabilities `json:"workspaceEdit,omitempty"`
	// Capabilities specific to the `workspace/didChangeConfiguration` notification.
	DidChangeConfiguration *DynamicRegistrationCapability `json:"didChangeConfiguration,omitempty"`
	// Capabilities specific to the `workspace/didChangeWatchedFiles` notification.
	DidChangeWatchedFiles *DynamicRegistrationCapability `json:"didChangeWatchedFiles,omitempty"`
	// Capabilities specific to the `workspace/symbol` request.
	Symbol *WorkspaceSymbolClientCapabilities `json:"symbol,omitempty"`
	// Capabilities specific to the `workspace/executeCommand` request.
	ExecuteCommand *DynamicRegistrationCapability `json:"executeCommand,omitempty"`
	// The client has support for workspace folders.
	//
	// Since 3.6.0
	WorkspaceFolders *bool `json:"workspaceFolders,omitempty"`
	// The client supports `workspace/configuration` requests.
	//
	// Since 3.6.0
	Configuration *bool `json:"configuration,omitempty"`
	// Capabilities specific to the semantic token requests scoped to the
	// workspace.
	//
	// Since 3.16.0.
	SemanticTokens *SemanticTokensWorkspaceClientCapabilities `json:"semanticTokens,omitempty"`
	// Capabilities specific to the code lens requests scoped to the
	// workspace.
	//
	// Since 3.16.0.
	CodeLens *CodeLensWorkspaceClientCapabilities `json:"codeLens,omitempty"`
	// The client has support for file notifications/requests for user operations on files.
	//
	// Since 3.16.0
	FileOperations *FileOperationClientCapabilities `json:"fileOperations,omitempty"`
	// Capabilities specific to the inlay hint requests scoped to the
	// workspace.
	//
	// Since 3.17.0.
	InlayHint *InlayHintWorkspaceClientCapabilities `json:"inlayHint,omitempty"`
	// Capabilities specific to the inline values requests scoped to the
	// workspace.
	//
	// Since 3.17.0.
	InlineValue *InlineValueWorkspaceClientCapabilities `json:"inlineValue,omitempty"`
	// Capabilities specific to the diagnostic requests scoped to the
	// workspace.
	//
	// Since 3.17.0.
	Diagnostics *DiagnosticWorkspaceClientCapabilities `json:"diagnostics,omitempty"`
}

// DynamicRegistrationCapability indicates the editor can register/unregister capabilities at runtime rather than only at initialization.
type DynamicRegistrationCapability struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
}

// WorkspaceEditClientCapabilities declares which workspace edit features the editor supports (document changes, resource operations, change annotations).
type WorkspaceEditClientCapabilities struct {
	// The client supports versioned document changes in WorkspaceEdits
	DocumentChanges *bool `json:"documentChanges,omitempty"`
	// The resource operations the client supports. Clients should at least
	// support 'create', 'rename' and 'delete' files and folders.
	//
	// Since 3.13.0
	ResourceOperations []ResourceOperationKind `json:"resourceOperations,omitempty"`
	// The failure handling strategy of a client if applying the workspace edit
	// fails.
	//
	// Since 3.13.0
	FailureHandling *FailureHandlingKind `json:"failureHandling,omitempty"`
	// Whether the client normalizes line endings to the client specific
	// setting.
	// If set to true the client will normalize line ending characters
	// in a workspace edit to the client-specified new line
	// character.
	//
	// Since 3.16.0
	NormalizesLineEndings *bool `json:"normalizesLineEndings,omitempty"`
	// Whether the client in general supports change annotations on text edits,
	// create file, rename file and delete file changes.
	//
	// Since 3.16.0
	ChangeAnnotationSupport *struct {
		GroupsOnLabel *bool `json:"groupsOnLabel,omitempty"`
	} `json:"changeAnnotationSupport,omitempty"`
}

// WorkspaceSymbolClientCapabilities declares client capabilities for a [WorkspaceSymbolRequest].
type WorkspaceSymbolClientCapabilities struct {
	// Symbol request supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// Specific capabilities for the SymbolKind in the `workspace/symbol` request.
	SymbolKind *struct {
		ValueSet []SymbolKind `json:"valueSet,omitempty"`
	} `json:"symbolKind,omitempty"`
	// The client supports tags on SymbolInformation.
	// Clients supporting tags have to handle unknown tags gracefully.
	//
	// Since 3.16.0
	TagSupport *struct {
		ValueSet []SymbolTag `json:"valueSet,omitempty"`
	} `json:"tagSupport,omitempty"`
}

// SemanticTokensWorkspaceClientCapabilities declares client support for workspace-wide semantic-tokens refreshes.
//
// Since 3.16.0.
type SemanticTokensWorkspaceClientCapabilities struct {
	// Whether the client implementation supports a refresh request sent from
	// the server to the client.
	//
	// Note that this event is global and will force the client to refresh all
	// semantic tokens currently shown. It should be used with absolute care
	// and is useful for situation where a server for example detects a project
	// wide change that requires such a calculation.
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}

// CodeLensWorkspaceClientCapabilities declares client support for workspace-wide code-lens refreshes.
//
// Since 3.16.0.
type CodeLensWorkspaceClientCapabilities struct {
	// Whether the client implementation supports a refresh request sent from the
	// server to the client.
	//
	// Note that this event is global and will force the client to refresh all
	// code lenses currently shown. It should be used with absolute care and is
	// useful for situation where a server for example detects a project wide
	// change that requires such a calculation.
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}

// FileOperationClientCapabilities relates to events from file operations by the user in the client.
//
// These events do not come from the file system, they come from user operations
// like renaming a file in the UI.
//
// Since 3.16.0.
type FileOperationClientCapabilities struct {
	// Whether the client supports dynamic registration for file requests/notifications.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client has support for sending didCreateFiles notifications.
	DidCreate *bool `json:"didCreate,omitempty"`
	// The client has support for sending willCreateFiles requests.
	WillCreate *bool `json:"willCreate,omitempty"`
	// The client has support for sending didRenameFiles notifications.
	DidRename *bool `json:"didRename,omitempty"`
	// The client has support for sending willRenameFiles requests.
	WillRename *bool `json:"willRename,omitempty"`
	// The client has support for sending didDeleteFiles notifications.
	DidDelete *bool `json:"didDelete,omitempty"`
	// The client has support for sending willDeleteFiles requests.
	WillDelete *bool `json:"willDelete,omitempty"`
}

// TextDocumentClientCapabilities declares text-document-specific client capabilities.
type TextDocumentClientCapabilities struct {
	// Defines which synchronization capabilities the client supports.
	Synchronization *TextDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
	// Capabilities specific to the `textDocument/completion` request.
	Completion *CompletionClientCapabilities `json:"completion,omitempty"`
	// Capabilities specific to the `textDocument/hover` request.
	Hover *HoverClientCapabilities `json:"hover,omitempty"`
	// Capabilities specific to the `textDocument/signatureHelp` request.
	SignatureHelp *SignatureHelpClientCapabilities `json:"signatureHelp,omitempty"`
	// Capabilities specific to the `textDocument/declaration` request.
	//
	// Since 3.14.0
	Declaration *DeclarationClientCapabilities `json:"declaration,omitempty"`
	// Capabilities specific to the `textDocument/definition` request.
	Definition *DefinitionClientCapabilities `json:"definition,omitempty"`
	// Capabilities specific to the `textDocument/typeDefinition` request.
	//
	// Since 3.6.0
	TypeDefinition *TypeDefinitionClientCapabilities `json:"typeDefinition,omitempty"`
	// Capabilities specific to the `textDocument/implementation` request.
	//
	// Since 3.6.0
	Implementation *ImplementationClientCapabilities `json:"implementation,omitempty"`
	// Capabilities specific to the `textDocument/references` request.
	References *DynamicRegistrationCapability `json:"references,omitempty"`
	// Capabilities specific to the `textDocument/documentHighlight` request.
	DocumentHighlight *DynamicRegistrationCapability `json:"documentHighlight,omitempty"`
	// Capabilities specific to the `textDocument/documentSymbol` request.
	DocumentSymbol *DocumentSymbolClientCapabilities `json:"documentSymbol,omitempty"`
	// Capabilities specific to the `textDocument/codeAction` request.
	CodeAction *CodeActionClientCapabilities `json:"codeAction,omitempty"`
	// Capabilities specific to the `textDocument/codeLens` request.
	CodeLens *DynamicRegistrationCapability `json:"codeLens,omitempty"`
	// Capabilities specific to the `textDocument/documentLink` request.
	DocumentLink *DocumentLinkClientCapabilities `json:"documentLink,omitempty"`
	// Capabilities specific to the `textDocument/documentColor` and the
	// `textDocument/colorPresentation` request.
	//
	// Since 3.6.0
	ColorProvider *DynamicRegistrationCapability `json:"colorProvider,omitempty"`
	// Capabilities specific to the `textDocument/formatting` request.
	Formatting *DynamicRegistrationCapability `json:"formatting,omitempty"`
	// Capabilities specific to the `textDocument/rangeFormatting` request.
	RangeFormatting *DynamicRegistrationCapability `json:"rangeFormatting,omitempty"`
	// Capabilities specific to the `textDocument/onTypeFormatting` request.
	OnTypeFormatting *DynamicRegistrationCapability `json:"onTypeFormatting,omitempty"`
	// Capabilities specific to the `textDocument/rename` request.
	Rename *RenameClientCapabilities `json:"rename,omitempty"`
	// Capabilities specific to the `textDocument/publishDiagnostics` notification.
	PublishDiagnostics *PublishDiagnosticsClientCapabilities `json:"publishDiagnostics,omitempty"`
	// Capabilities specific to the `textDocument/foldingRange` request.
	//
	// Since 3.10.0
	FoldingRange *FoldingRangeClientCapabilities `json:"foldingRange,omitempty"`
	// Capabilities specific to the `textDocument/selectionRange` request.
	//
	// Since 3.15.0
	SelectionRange *DynamicRegistrationCapability `json:"selectionRange,omitempty"`
	// Capabilities specific to the `textDocument/linkedEditingRange` request.
	//
	// Since 3.16.0
	LinkedEditingRange *DynamicRegistrationCapability `json:"linkedEditingRange,omitempty"`
	// Capabilities specific to the various call hierarchy requests.
	//
	// Since 3.16.0
	CallHierarchy *DynamicRegistrationCapability `json:"callHierarchy,omitempty"`
	// Capabilities specific to the various semantic token requests.
	//
	// Since 3.16.0
	SemanticTokens *SemanticTokensClientCapabilities `json:"semanticTokens,omitempty"`
	// Client capabilities specific to the `textDocument/moniker` request.
	//
	// Since 3.16.0
	Moniker *DynamicRegistrationCapability `json:"moniker,omitempty"`
	// Capabilities specific to the various type hierarchy requests.
	//
	// Since 3.17.0
	TypeHierarchy *DynamicRegistrationCapability `json:"typeHierarchy,omitempty"`
	// Capabilities specific to the `textDocument/inlayHint` request.
	//
	// Since 3.17.0
	InlayHint *InlayHintClientCapabilities `json:"inlayHint,omitempty"`
	// Capabilities specific to the `textDocument/inlineValue` request.
	//
	// Since 3.17.0
	InlineValue *DynamicRegistrationCapability `json:"inlineValue,omitempty"`
	// Capabilities specific to the diagnostic pull model.
	//
	// Since 3.17.0
	Diagnostic *DiagnosticClientCapabilities `json:"diagnostic,omitempty"`
}

// TextDocumentSyncClientCapabilities declares editor support for open/close/change/save document notifications.
type TextDocumentSyncClientCapabilities struct {
	// Whether text document synchronization supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client supports sending will save notifications.
	WillSave *bool `json:"willSave,omitempty"`
	// The client supports sending a will save request and
	// waits for a response providing text edits which will
	// be applied to the document before it is saved.
	WillSaveWaitUntil *bool `json:"willSaveWaitUntil,omitempty"`
	// The client supports did save notifications.
	DidSave *bool `json:"didSave,omitempty"`
}

// CompletionClientCapabilities declares client support for completion requests.
type CompletionClientCapabilities struct {
	// Whether completion supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client supports the following CompletionItem specific
	// capabilities.
	CompletionItem *struct {
		SnippetSupport          *bool        `json:"snippetSupport,omitempty"`
		CommitCharactersSupport *bool        `json:"commitCharactersSupport,omitempty"`
		DocumentationFormat     []MarkupKind `json:"documentationFormat,omitempty"`
		DeprecatedSupport       *bool        `json:"deprecatedSupport,omitempty"`
		PreselectSupport        *bool        `json:"preselectSupport,omitempty"`
		TagSupport              *struct {
			ValueSet []CompletionItemTag `json:"valueSet"`
		} `json:"tagSupport,omitempty"`
		InsertReplaceSupport *bool `json:"insertReplaceSupport,omitempty"`
		ResolveSupport       *struct {
			Properties []string `json:"properties"`
		} `json:"resolveSupport,omitempty"`
		InsertTextModeSupport *struct {
			ValueSet []InsertTextMode `json:"valueSet"`
		} `json:"insertTextModeSupport,omitempty"`
	} `json:"completionItem,omitempty"`
	CompletionItemKind *struct {
		ValueSet []CompletionItemKind `json:"valueSet,omitempty"`
	} `json:"completionItemKind,omitempty"`
	// The client supports to send additional context information for a
	// `textDocument/completion` request.
	ContextSupport *bool `json:"contextSupport,omitempty"`
}

// HoverClientCapabilities declares which content formats (plaintext, markdown) the editor supports in hover results.
type HoverClientCapabilities struct {
	// Whether hover supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// Client supports the following content formats for the content
	// property. The order describes the preferred format of the client.
	ContentFormat []MarkupKind `json:"contentFormat,omitempty"`
}

// SignatureHelpClientCapabilities declares client capabilities for a [SignatureHelpRequest].
type SignatureHelpClientCapabilities struct {
	// Whether signature help supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client supports the following SignatureInformation
	// specific properties.
	SignatureInformation *struct {
		DocumentationFormat  []MarkupKind `json:"documentationFormat,omitempty"`
		ParameterInformation *struct {
			LabelOffsetSupport *bool `json:"labelOffsetSupport,omitempty"`
		} `json:"parameterInformation,omitempty"`
		ActiveParameterSupport *bool `json:"activeParameterSupport,omitempty"`
	} `json:"signatureInformation,omitempty"`
	// The client supports to send additional context information for a
	// `textDocument/signatureHelp` request. A client that opts into
	// contextSupport will also support the retriggerCharacters on
	// SignatureHelpOptions.
	//
	// Since 3.15.0
	ContextSupport *bool `json:"contextSupport,omitempty"`
}

// DeclarationClientCapabilities declares client support for go-to-declaration requests.
//
// Since 3.14.0.
type DeclarationClientCapabilities struct {
	// Whether declaration supports dynamic registration. If this is set to true
	// the client supports the new DeclarationRegistrationOptions return value
	// for the corresponding server capability as well.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client supports additional metadata in the form of declaration links.
	LinkSupport *bool `json:"linkSupport,omitempty"`
}

// DefinitionClientCapabilities declares client capabilities for a [DefinitionRequest].
type DefinitionClientCapabilities struct {
	// Whether definition supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client supports additional metadata in the form of definition links.
	//
	// Since 3.14.0
	LinkSupport *bool `json:"linkSupport,omitempty"`
}

// TypeDefinitionClientCapabilities declares client support for go-to-type-definition requests.
//
// Since 3.6.0.
type TypeDefinitionClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to true
	// the client supports the new TypeDefinitionRegistrationOptions return value
	// for the corresponding server capability as well.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client supports additional metadata in the form of definition links.
	//
	// Since 3.14.0
	LinkSupport *bool `json:"linkSupport,omitempty"`
}

// ImplementationClientCapabilities declares client support for go-to-implementation requests.
//
// Since 3.6.0.
type ImplementationClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to true
	// the client supports the new ImplementationRegistrationOptions return value
	// for the corresponding server capability as well.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client supports additional metadata in the form of definition links.
	//
	// Since 3.14.0
	LinkSupport *bool `json:"linkSupport,omitempty"`
}

// DocumentSymbolClientCapabilities declares client capabilities for a [DocumentSymbolRequest].
type DocumentSymbolClientCapabilities struct {
	// Whether document symbol supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// Specific capabilities for the SymbolKind in the
	// `textDocument/documentSymbol` request.
	SymbolKind *struct {
		ValueSet []SymbolKind `json:"valueSet,omitempty"`
	} `json:"symbolKind,omitempty"`
	// The client supports hierarchical document symbols.
	HierarchicalDocumentSymbolSupport *bool `json:"hierarchicalDocumentSymbolSupport,omitempty"`
	// The client supports tags on SymbolInformation. Tags are supported on
	// DocumentSymbol if hierarchicalDocumentSymbolSupport is set to true.
	// Clients supporting tags have to handle unknown tags gracefully.
	//
	// Since 3.16.0
	TagSupport *struct {
		ValueSet []SymbolTag `json:"valueSet,omitempty"`
	} `json:"tagSupport,omitempty"`
	// The client supports an additional label presented in the UI when
	// registering a document symbol provider.
	//
	// Since 3.16.0
	LabelSupport *bool `json:"labelSupport,omitempty"`
}

// CodeActionClientCapabilities defines the client capabilities of a [CodeActionRequest].
type CodeActionClientCapabilities struct {
	// Whether code action supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client support code action literals of type CodeAction as a valid
	// response of the `textDocument/codeAction` request. If the property is not
	// set the request can only return Command literals.
	//
	// Since 3.8.0
	CodeActionLiteralSupport *struct {
		CodeActionKind struct {
			ValueSet []CodeActionKind `json:"valueSet"`
		} `json:"codeActionKind"`
	} `json:"codeActionLiteralSupport,omitempty"`
	// Whether code action supports the isPreferred property.
	//
	// Since 3.15.0
	IsPreferredSupport *bool `json:"isPreferredSupport,omitempty"`
	// Whether code action supports the disabled property.
	//
	// Since 3.16.0
	DisabledSupport *bool `json:"disabledSupport,omitempty"`
	// Whether code action supports the data property which is
	// preserved between a `textDocument/codeAction` and a
	// `codeAction/resolve` request.
	//
	// Since 3.16.0
	DataSupport *bool `json:"dataSupport,omitempty"`
	// Whether the client supports resolving additional code action
	// properties via a separate `codeAction/resolve` request.
	//
	// Since 3.16.0
	ResolveSupport *struct {
		Properties []string `json:"properties"`
	} `json:"resolveSupport,omitempty"`
	// Whether the client honors the change annotations in
	// text edits and resource operations returned via the
	// `CodeAction#edit` property by for example presenting
	// the workspace edit in the user interface and asking
	// for confirmation.
	//
	// Since 3.16.0
	HonorsChangeAnnotations *bool `json:"honorsChangeAnnotations,omitempty"`
}

// DocumentLinkClientCapabilities defines the client capabilities of a [DocumentLinkRequest].
type DocumentLinkClientCapabilities struct {
	// Whether document link supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// Whether the client supports the tooltip property on DocumentLink.
	//
	// Since 3.15.0
	TooltipSupport *bool `json:"tooltipSupport,omitempty"`
}

// RenameClientCapabilities declares editor support for rename features like prepare-rename and honoring change annotations.
type RenameClientCapabilities struct {
	// Whether rename supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// Client supports testing for validity of rename operations
	// before execution.
	//
	// Since 3.12.0
	PrepareSupport *bool `json:"prepareSupport,omitempty"`
	// Client supports the default behavior result.
	//
	// The value indicates the default behavior used by the
	// client.
	//
	// Since 3.16.0
	PrepareSupportDefaultBehavior *PrepareSupportDefaultBehavior `json:"prepareSupportDefaultBehavior,omitempty"`
	// Whether the client honors the change annotations in
	// text edits and resource operations returned via the
	// rename request's workspace edit by for example presenting
	// the workspace edit in the user interface and asking
	// for confirmation.
	//
	// Since 3.16.0
	HonorsChangeAnnotations *bool `json:"honorsChangeAnnotations,omitempty"`
}

// PublishDiagnosticsClientCapabilities declares client capabilities for published diagnostics.
type PublishDiagnosticsClientCapabilities struct {
	// Whether the clients accepts diagnostics with related information.
	RelatedInformation *bool `json:"relatedInformation,omitempty"`
	// Client supports the tag property to provide meta data about a diagnostic.
	// Clients supporting tags have to handle unknown tags gracefully.
	//
	// Since 3.15.0
	TagSupport *struct {
		ValueSet []DiagnosticTag `json:"valueSet,omitempty"`
	} `json:"tagSupport,omitempty"`
	// Whether the client interprets the version property of the
	// `textDocument/publishDiagnostics` notification's parameter.
	//
	// Since 3.15.0
	VersionSupport *bool `json:"versionSupport,omitempty"`
	// Client supports a codeDescription property
	//
	// Since 3.16.0
	CodeDescriptionSupport *bool `json:"codeDescriptionSupport,omitempty"`
	// Whether code action supports the data property which is
	// preserved between a `textDocument/publishDiagnostics` and
	// `textDocument/codeAction` request.
	//
	// Since 3.16.0
	DataSupport *bool `json:"dataSupport,omitempty"`
}

// FoldingRangeClientCapabilities declares editor support for folding ranges, including range limits and line-only folding.
type FoldingRangeClientCapabilities struct {
	// Whether implementation supports dynamic registration for folding range
	// providers. If this is set to true the client supports the new
	// FoldingRangeRegistrationOptions return value for the corresponding
	// server capability as well.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The maximum number of folding ranges that the client prefers to receive
	// per document. The value serves as a hint, servers are free to follow the
	// limit.
	RangeLimit *int `json:"rangeLimit,omitempty"`
	// If set, the client signals that it only supports folding complete lines.
	// If set, client will ignore specified startCharacter and endCharacter
	// properties in a FoldingRange.
	LineFoldingOnly *bool `json:"lineFoldingOnly,omitempty"`
}

// SemanticTokensClientCapabilities declares client support for semantic-tokens requests.
//
// Since 3.16.0.
type SemanticTokensClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to true
	// the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
	// return value for the corresponding server capability as well.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// Which requests the client supports and might send to the server
	// depending on the server's capability. Please note that clients might not
	// show semantic tokens or degrade some of the user experience if a range
	// or full request is advertised by the client but not provided by the
	// server. If for example the client capability `requests.full` and
	// `request.range` are both set to true but the server only provides a
	// range provider the client might not render a minimap correctly or might
	// even decide to not show any semantic tokens at all.
	Requests SemanticTokensRequestsCapabilities `json:"requests"`
	// The token types that the client supports.
	TokenTypes []string `json:"tokenTypes"`
	// The token modifiers that the client supports.
	TokenModifiers []string `json:"tokenModifiers"`
	// The token formats the client supports.
	Formats []TokenFormat `json:"formats"`
	// Whether the client supports tokens that can overlap each other.
	OverlappingTokenSupport *bool `json:"overlappingTokenSupport,omitempty"`
	// Whether the client supports tokens that can span multiple lines.
	MultilineTokenSupport *bool `json:"multilineTokenSupport,omitempty"`
}

// SemanticTokensRequestsCapabilities describes the semantic token request styles the client supports.
// Per the LSP spec, "range" is `boolean | {}` and "full" is `boolean | { delta?: boolean }`;
// UnmarshalJSON accepts either form.
type SemanticTokensRequestsCapabilities struct {
	Range *bool               `json:"range,omitempty"`
	Full  *SemanticTokensFull `json:"full,omitempty"`
}

// UnmarshalJSON accepts both the boolean and object forms for "range" and "full" per the LSP spec.
func (s *SemanticTokensRequestsCapabilities) UnmarshalJSON(data []byte) error {
	var raw struct {
		Range json.RawMessage `json:"range"`
		Full  json.RawMessage `json:"full"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	if len(raw.Range) > 0 && string(raw.Range) != "null" {
		var b bool
		if err := json.Unmarshal(raw.Range, &b); err == nil {
			s.Range = &b
		} else {
			// Empty object form ({}) — presence indicates support.
			t := true
			s.Range = &t
		}
	}

	if len(raw.Full) > 0 && string(raw.Full) != "null" {
		var b bool
		if err := json.Unmarshal(raw.Full, &b); err == nil {
			if b {
				s.Full = &SemanticTokensFull{}
			}
			// false → leave Full nil to signal unsupported.
		} else {
			var f SemanticTokensFull
			if err := json.Unmarshal(raw.Full, &f); err != nil {
				return err
			}
			s.Full = &f
		}
	}

	return nil
}

// WindowClientCapabilities declares editor support for window features like work-done progress, show-message requests, and show-document.
type WindowClientCapabilities struct {
	// It indicates whether the client supports server initiated
	// progress using the `window/workDoneProgress/create` request.
	//
	// The capability also controls whether the client supports handling
	// of progress notifications. If set servers are allowed to report a
	// workDoneProgress property in the request specific server
	// capabilities.
	//
	// Since 3.15.0
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
	// Capabilities specific to the showMessage request.
	//
	// Since 3.16.0
	ShowMessage *struct {
		MessageActionItem *struct {
			AdditionalPropertiesSupport *bool `json:"additionalPropertiesSupport,omitempty"`
		} `json:"messageActionItem,omitempty"`
	} `json:"showMessage,omitempty"`
	// Capabilities specific to the showDocument request.
	//
	// Since 3.16.0
	ShowDocument *struct {
		Support bool `json:"support"`
	} `json:"showDocument,omitempty"`
}

// GeneralClientCapabilities declares general client capabilities not tied to a single feature.
//
// Since 3.16.0.
type GeneralClientCapabilities struct {
	// The position encodings supported by the client. Client and server
	// have to agree on the same position encoding to ensure that offsets
	// (e.g. character position in a line) are interpreted the same on both
	// sides.
	//
	// To keep the protocol backwards compatible the following applies: if
	// the value 'utf-16' is missing from the array of position encodings
	// servers can assume that the client supports UTF-16. UTF-16 is
	// therefore a mandatory encoding.
	//
	// If omitted it defaults to ['utf-16'].
	//
	// Implementation considerations: since the conversion from one encoding
	// into another requires the content of the file / line the conversion
	// is best done where the file is read which is usually on the server
	// side.
	//
	// Since 3.17.0
	PositionEncodings []PositionEncodingKind `json:"positionEncodings,omitempty"`
	// Client capabilities specific to regular expressions.
	//
	// Since 3.16.0
	RegularExpressions *struct {
		Engine  string `json:"engine"`
		Version string `json:"version,omitempty"`
	} `json:"regularExpressions,omitempty"`
	// Client capabilities specific to the client's markdown parser.
	//
	// Since 3.16.0
	Markdown *struct {
		Parser      string   `json:"parser"`
		Version     string   `json:"version,omitempty"`
		AllowedTags []string `json:"allowedTags,omitempty"`
	} `json:"markdown,omitempty"`
}

// ServerCapabilities defines the capabilities provided by a language
// server.
type ServerCapabilities struct {
	// The position encoding the server picked from the encodings offered
	// by the client via the client capability `general.positionEncodings`.
	//
	// If the client didn't provide any position encodings the only valid
	// value that a server can return is 'utf-16'.
	//
	// If omitted it defaults to 'utf-16'.
	//
	// Since 3.17.0
	PositionEncoding *PositionEncodingKind `json:"positionEncoding,omitempty"`
	// Defines how text documents are synced. It is either a detailed structure
	// defining each notification or for backwards compatibility the
	// TextDocumentSyncKind number.
	TextDocumentSync *TextDocumentSyncOptions `json:"textDocumentSync,omitempty"`
	// The server provides completion support.
	CompletionProvider *CompletionOptions `json:"completionProvider,omitempty"`
	// The server provides hover support.
	HoverProvider *bool `json:"hoverProvider,omitempty"`
	// The server provides signature help support.
	SignatureHelpProvider *SignatureHelpOptions `json:"signatureHelpProvider,omitempty"`
	// The server provides Goto declaration support.
	DeclarationProvider *bool `json:"declarationProvider,omitempty"`
	// The server provides goto definition support.
	DefinitionProvider *bool `json:"definitionProvider,omitempty"`
	// The server provides Goto type definition support.
	TypeDefinitionProvider *bool `json:"typeDefinitionProvider,omitempty"`
	// The server provides Goto implementation support.
	ImplementationProvider *bool `json:"implementationProvider,omitempty"`
	// The server provides find references support.
	ReferencesProvider *bool `json:"referencesProvider,omitempty"`
	// The server provides document highlight support.
	DocumentHighlightProvider *bool `json:"documentHighlightProvider,omitempty"`
	// The server provides document symbol support.
	DocumentSymbolProvider *bool `json:"documentSymbolProvider,omitempty"`
	// The server provides code actions. CodeActionOptions may only be
	// specified if the client states that it supports
	// codeActionLiteralSupport in its initial initialize request.
	CodeActionProvider *CodeActionOptions `json:"codeActionProvider,omitempty"`
	// The server provides code lens.
	CodeLensProvider *CodeLensOptions `json:"codeLensProvider,omitempty"`
	// The server provides document link support.
	DocumentLinkProvider *DocumentLinkOptions `json:"documentLinkProvider,omitempty"`
	// The server provides color provider support.
	ColorProvider *bool `json:"colorProvider,omitempty"`
	// The server provides document formatting.
	DocumentFormattingProvider *bool `json:"documentFormattingProvider,omitempty"`
	// The server provides document range formatting.
	DocumentRangeFormattingProvider *bool `json:"documentRangeFormattingProvider,omitempty"`
	// The server provides document formatting on typing.
	DocumentOnTypeFormattingProvider *DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`
	// The server provides rename support. RenameOptions may only be
	// specified if the client states that it supports
	// prepareSupport in its initial initialize request.
	RenameProvider *RenameOptions `json:"renameProvider,omitempty"`
	// The server provides folding provider support.
	FoldingRangeProvider *bool `json:"foldingRangeProvider,omitempty"`
	// The server provides execute command support.
	ExecuteCommandProvider *ExecuteCommandOptions `json:"executeCommandProvider,omitempty"`
	// The server provides selection range support.
	SelectionRangeProvider *bool `json:"selectionRangeProvider,omitempty"`
	// The server provides linked editing range support.
	//
	// Since 3.16.0
	LinkedEditingRangeProvider *bool `json:"linkedEditingRangeProvider,omitempty"`
	// The server provides call hierarchy support.
	//
	// Since 3.16.0
	CallHierarchyProvider *bool `json:"callHierarchyProvider,omitempty"`
	// The server provides semantic tokens support.
	//
	// Since 3.16.0
	SemanticTokensProvider *SemanticTokensOptions `json:"semanticTokensProvider,omitempty"`
	// The server provides moniker support.
	//
	// Since 3.16.0
	MonikerProvider *bool `json:"monikerProvider,omitempty"`
	// The server provides type hierarchy support.
	//
	// Since 3.17.0
	TypeHierarchyProvider *bool `json:"typeHierarchyProvider,omitempty"`
	// The server provides inlay hints.
	//
	// Since 3.17.0
	InlayHintProvider *InlayHintOptions `json:"inlayHintProvider,omitempty"`
	// The server provides inline values.
	//
	// Since 3.17.0
	InlineValueProvider *bool `json:"inlineValueProvider,omitempty"`
	// The server has support for pull model diagnostics.
	//
	// Since 3.17.0
	DiagnosticProvider *DiagnosticOptions `json:"diagnosticProvider,omitempty"`
	// The server provides workspace symbol support.
	WorkspaceSymbolProvider *bool `json:"workspaceSymbolProvider,omitempty"`
	// Workspace specific server capabilities.
	Workspace *ServerWorkspaceCapabilities `json:"workspace,omitempty"`
	// Experimental server capabilities.
	Experimental json.RawMessage `json:"experimental,omitempty"`
}

// ServerWorkspaceCapabilities declares server support for workspace features like workspace folders and file operations.
type ServerWorkspaceCapabilities struct {
	WorkspaceFolders *WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
	FileOperations   *FileOperationOptions               `json:"fileOperations,omitempty"`
}

// WorkspaceFoldersServerCapabilities declares whether the server supports multi-root workspaces and wants workspace folder change notifications.
type WorkspaceFoldersServerCapabilities struct {
	// The server has support for workspace folders
	Supported *bool `json:"supported,omitempty"`
	// Whether the server wants to receive workspace folder
	// change notifications.
	//
	// If a string is provided the string is treated as an ID
	// under which the notification is registered on the client
	// side. The ID can be used to unregister for these events
	// using the `client/unregisterCapability` request.
	ChangeNotifications *bool `json:"changeNotifications,omitempty"`
}

// FileOperationOptions is for notifications/requests for user operations on files.
//
// Since 3.16.0.
type FileOperationOptions struct {
	// The server is interested in receiving didCreateFiles notifications.
	DidCreate *FileOperationRegistrationOptions `json:"didCreate,omitempty"`
	// The server is interested in receiving willCreateFiles requests.
	WillCreate *FileOperationRegistrationOptions `json:"willCreate,omitempty"`
	// The server is interested in receiving didRenameFiles notifications.
	DidRename *FileOperationRegistrationOptions `json:"didRename,omitempty"`
	// The server is interested in receiving willRenameFiles requests.
	WillRename *FileOperationRegistrationOptions `json:"willRename,omitempty"`
	// The server is interested in receiving didDeleteFiles notifications.
	DidDelete *FileOperationRegistrationOptions `json:"didDelete,omitempty"`
	// The server is interested in receiving willDeleteFiles requests.
	WillDelete *FileOperationRegistrationOptions `json:"willDelete,omitempty"`
}

// FileOperationRegistrationOptions holds the options to register for file operations.
//
// Since 3.16.0.
type FileOperationRegistrationOptions struct {
	// The actual filters.
	Filters []FileOperationFilter `json:"filters"`
}

// FileOperationFilter is a filter to describe in which file operation requests or notifications
// the server is interested in receiving.
//
// Since 3.16.0.
type FileOperationFilter struct {
	// A Uri scheme like file or untitled.
	Scheme string `json:"scheme,omitempty"`
	// The actual file operation pattern.
	Pattern FileOperationPattern `json:"pattern"`
}

// FileOperationPattern is a pattern to describe in which file operation requests or notifications
// the server is interested in receiving.
//
// Since 3.16.0.
type FileOperationPattern struct {
	// The glob pattern to match. Glob patterns can have the following syntax:
	// - `*` to match zero or more characters in a path segment
	// - `?` to match on one character in a path segment
	// - `**` to match any number of path segments, including none
	// - `{}` to group sub patterns into an OR expression. (e.g. `**​/*.{ts,js}` matches all TypeScript and JavaScript files)
	// - `[]` to declare a range of characters to match in a path segment (e.g., `example.[0-9]` to match on `example.0`, `example.1`, …)
	// - `[!...]` to negate a range of characters to match in a path segment (e.g., `example.[!0-9]` to match on `example.a`, `example.b`, but not `example.0`)
	Glob string `json:"glob"`
	// Whether to match files or folders with this pattern.
	//
	// Matches both if nil.
	Matches *FileOperationPatternKind `json:"matches,omitempty"`
	// Additional options used during matching.
	Options *FileOperationPatternOptions `json:"options,omitempty"`
}

// FileOperationPatternKind is a string enum ("file" or "folder") filtering which filesystem entries a file operation pattern matches.
type FileOperationPatternKind string

const (
	// The pattern matches a file only.
	FileOperationPatternKindFile FileOperationPatternKind = "file"
	// The pattern matches a folder only.
	FileOperationPatternKindFolder FileOperationPatternKind = "folder"
)

// FileOperationPatternOptions is the matching options for the file operation pattern.
//
// Since 3.16.0.
type FileOperationPatternOptions struct {
	// The pattern should be matched ignoring casing.
	IgnoreCase *bool `json:"ignoreCase,omitempty"`
}

// TextDocumentSyncOptions configures how the server receives document content: open/close notifications, incremental vs full change events, and save behavior.
type TextDocumentSyncOptions struct {
	// Open and close notifications are sent to the server. If omitted, open and close notifications should not
	// be sent.
	OpenClose *bool `json:"openClose,omitempty"`
	// Change notifications are sent to the server. See [SyncNone], [SyncFull]
	// and [SyncIncremental]. If omitted it defaults to [SyncNone].
	Change TextDocumentSyncKind `json:"change,omitempty"`
	// If present, will save notifications are sent to the server. If omitted, the notifications should not be
	// sent.
	WillSave *bool `json:"willSave,omitempty"`
	// If present, will save wait until requests are sent to the server. If omitted, the requests should not be
	// sent.
	WillSaveWaitUntil *bool `json:"willSaveWaitUntil,omitempty"`
	// If present, save notifications are sent to the server. If omitted, the notifications should not be
	// sent.
	Save *SaveOptions `json:"save,omitempty"`
}

// SaveOptions configures save-notification behavior.
type SaveOptions struct {
	// The client is supposed to include the content on save.
	IncludeText *bool `json:"includeText,omitempty"`
}

// CompletionOptions configures the server's completion provider.
type CompletionOptions struct {
	WorkDoneProgressOptions
	// Most tools trigger completion request automatically without explicitly requesting
	// it using a keyboard shortcut (e.g. Ctrl+Space). Typically they do so when the user
	// starts to type an identifier. For example if the user types c in a JavaScript file
	// code complete will automatically pop up present console besides others as a
	// completion item. Characters that make up identifiers don't need to be listed here.
	//
	// If code complete should automatically be triggered on characters not being valid inside
	// an identifier (for example `.` in JavaScript), list them in triggerCharacters.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
	// The list of all possible characters that commit a completion. This field can be used
	// if clients don't support individual commit characters per completion item. See
	// `ClientCapabilities.textDocument.completion.completionItem.commitCharactersSupport`
	//
	// If a server provides both allCommitCharacters and commit characters on an individual
	// completion item, the ones on the completion item win.
	//
	// Since 3.2.0
	AllCommitCharacters []string `json:"allCommitCharacters,omitempty"`
	// The server provides support to resolve additional
	// information for a completion item.
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

// SignatureHelpOptions holds the server capabilities for a [SignatureHelpRequest].
type SignatureHelpOptions struct {
	WorkDoneProgressOptions
	// List of characters that trigger signature help automatically.
	TriggerCharacters []string `json:"triggerCharacters,omitempty"`
	// List of characters that re-trigger signature help.
	//
	// These trigger characters are only active when signature help is already showing. All trigger characters
	// are also counted as re-trigger characters.
	//
	// Since 3.15.0
	RetriggerCharacters []string `json:"retriggerCharacters,omitempty"`
}

// CodeActionOptions holds the provider options for a [CodeActionRequest].
type CodeActionOptions struct {
	WorkDoneProgressOptions
	// CodeActionKinds that this server may return.
	//
	// The list of kinds may be generic, such as `[CodeActionRefactor]`, or the server
	// may list out every specific kind they provide.
	CodeActionKinds []CodeActionKind `json:"codeActionKinds,omitempty"`
	// The server provides support to resolve additional
	// information for a code action.
	//
	// Since 3.16.0
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

// CodeLensOptions holds the provider options of a [CodeLensRequest].
type CodeLensOptions struct {
	WorkDoneProgressOptions
	// Code lens has a resolve provider as well.
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

// DocumentLinkOptions holds the provider options for a [DocumentLinkRequest].
type DocumentLinkOptions struct {
	WorkDoneProgressOptions
	// Document links have a resolve provider as well.
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

// DocumentOnTypeFormattingOptions holds the provider options for a [DocumentOnTypeFormattingRequest].
type DocumentOnTypeFormattingOptions struct {
	// A character on which formatting should be triggered, like `{`.
	FirstTriggerCharacter string `json:"firstTriggerCharacter"`
	// More trigger characters.
	MoreTriggerCharacter []string `json:"moreTriggerCharacter,omitempty"`
}

// RenameOptions holds the provider options for a [RenameRequest].
type RenameOptions struct {
	WorkDoneProgressOptions
	// Renames should be checked and tested before being executed.
	//
	// Since version 3.12.0
	PrepareProvider *bool `json:"prepareProvider,omitempty"`
}

// ExecuteCommandOptions holds the server capabilities for a [ExecuteCommandRequest].
type ExecuteCommandOptions struct {
	WorkDoneProgressOptions
	// The commands to be executed on the server
	Commands []string `json:"commands"`
}
