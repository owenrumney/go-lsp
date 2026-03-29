package lsp

import "encoding/json"

// ClientCapabilities declares which LSP features the editor supports, sent during initialization so the server can adapt its behavior.
type ClientCapabilities struct {
	Workspace    *WorkspaceClientCapabilities    `json:"workspace,omitempty"`
	TextDocument *TextDocumentClientCapabilities `json:"textDocument,omitempty"`
	Window       *WindowClientCapabilities       `json:"window,omitempty"`
	General      *GeneralClientCapabilities      `json:"general,omitempty"`
	Experimental json.RawMessage                 `json:"experimental,omitempty"`
}

// WorkspaceClientCapabilities declares editor support for workspace-wide operations like symbol search, file watching, and workspace edits.
type WorkspaceClientCapabilities struct {
	ApplyEdit              *bool                                      `json:"applyEdit,omitempty"`
	WorkspaceEdit          *WorkspaceEditClientCapabilities           `json:"workspaceEdit,omitempty"`
	DidChangeConfiguration *DynamicRegistrationCapability             `json:"didChangeConfiguration,omitempty"`
	DidChangeWatchedFiles  *DynamicRegistrationCapability             `json:"didChangeWatchedFiles,omitempty"`
	Symbol                 *WorkspaceSymbolClientCapabilities         `json:"symbol,omitempty"`
	ExecuteCommand         *DynamicRegistrationCapability             `json:"executeCommand,omitempty"`
	WorkspaceFolders       *bool                                      `json:"workspaceFolders,omitempty"`
	Configuration          *bool                                      `json:"configuration,omitempty"`
	SemanticTokens         *SemanticTokensWorkspaceClientCapabilities `json:"semanticTokens,omitempty"`
	CodeLens               *CodeLensWorkspaceClientCapabilities       `json:"codeLens,omitempty"`
	FileOperations         *FileOperationClientCapabilities           `json:"fileOperations,omitempty"`
	InlayHint              *InlayHintWorkspaceClientCapabilities      `json:"inlayHint,omitempty"`
	InlineValue            *InlineValueWorkspaceClientCapabilities    `json:"inlineValue,omitempty"`
	Diagnostics            *DiagnosticWorkspaceClientCapabilities     `json:"diagnostics,omitempty"`
}

// DynamicRegistrationCapability indicates the editor can register/unregister capabilities at runtime rather than only at initialization.
type DynamicRegistrationCapability struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
}

// WorkspaceEditClientCapabilities declares which workspace edit features the editor supports (document changes, resource operations, change annotations).
type WorkspaceEditClientCapabilities struct {
	DocumentChanges         *bool                   `json:"documentChanges,omitempty"`
	ResourceOperations      []ResourceOperationKind `json:"resourceOperations,omitempty"`
	FailureHandling         *FailureHandlingKind    `json:"failureHandling,omitempty"`
	NormalizesLineEndings   *bool                   `json:"normalizesLineEndings,omitempty"`
	ChangeAnnotationSupport *struct {
		GroupsOnLabel *bool `json:"groupsOnLabel,omitempty"`
	} `json:"changeAnnotationSupport,omitempty"`
}

// WorkspaceSymbolClientCapabilities declares editor support for workspace-wide symbol search, including supported symbol kinds and tags.
type WorkspaceSymbolClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	SymbolKind          *struct {
		ValueSet []SymbolKind `json:"valueSet,omitempty"`
	} `json:"symbolKind,omitempty"`
	TagSupport *struct {
		ValueSet []SymbolTag `json:"valueSet,omitempty"`
	} `json:"tagSupport,omitempty"`
}

// SemanticTokensWorkspaceClientCapabilities declares whether the editor will refresh semantic tokens when the server requests it.
type SemanticTokensWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}

// CodeLensWorkspaceClientCapabilities declares whether the editor will refresh code lenses when the server requests it.
type CodeLensWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}

// FileOperationClientCapabilities declares editor support for file create/rename/delete notifications and willCreate/willRename/willDelete requests.
type FileOperationClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	DidCreate           *bool `json:"didCreate,omitempty"`
	WillCreate          *bool `json:"willCreate,omitempty"`
	DidRename           *bool `json:"didRename,omitempty"`
	WillRename          *bool `json:"willRename,omitempty"`
	DidDelete           *bool `json:"didDelete,omitempty"`
	WillDelete          *bool `json:"willDelete,omitempty"`
}

// TextDocumentClientCapabilities declares editor support for document-level features like completion, hover, diagnostics, formatting, etc.
type TextDocumentClientCapabilities struct {
	Synchronization    *TextDocumentSyncClientCapabilities   `json:"synchronization,omitempty"`
	Completion         *CompletionClientCapabilities         `json:"completion,omitempty"`
	Hover              *HoverClientCapabilities              `json:"hover,omitempty"`
	SignatureHelp      *SignatureHelpClientCapabilities      `json:"signatureHelp,omitempty"`
	Declaration        *DeclarationClientCapabilities        `json:"declaration,omitempty"`
	Definition         *DefinitionClientCapabilities         `json:"definition,omitempty"`
	TypeDefinition     *TypeDefinitionClientCapabilities     `json:"typeDefinition,omitempty"`
	Implementation     *ImplementationClientCapabilities     `json:"implementation,omitempty"`
	References         *DynamicRegistrationCapability        `json:"references,omitempty"`
	DocumentHighlight  *DynamicRegistrationCapability        `json:"documentHighlight,omitempty"`
	DocumentSymbol     *DocumentSymbolClientCapabilities     `json:"documentSymbol,omitempty"`
	CodeAction         *CodeActionClientCapabilities         `json:"codeAction,omitempty"`
	CodeLens           *DynamicRegistrationCapability        `json:"codeLens,omitempty"`
	DocumentLink       *DocumentLinkClientCapabilities       `json:"documentLink,omitempty"`
	ColorProvider      *DynamicRegistrationCapability        `json:"colorProvider,omitempty"`
	Formatting         *DynamicRegistrationCapability        `json:"formatting,omitempty"`
	RangeFormatting    *DynamicRegistrationCapability        `json:"rangeFormatting,omitempty"`
	OnTypeFormatting   *DynamicRegistrationCapability        `json:"onTypeFormatting,omitempty"`
	Rename             *RenameClientCapabilities             `json:"rename,omitempty"`
	PublishDiagnostics *PublishDiagnosticsClientCapabilities `json:"publishDiagnostics,omitempty"`
	FoldingRange       *FoldingRangeClientCapabilities       `json:"foldingRange,omitempty"`
	SelectionRange     *DynamicRegistrationCapability        `json:"selectionRange,omitempty"`
	LinkedEditingRange *DynamicRegistrationCapability        `json:"linkedEditingRange,omitempty"`
	CallHierarchy      *DynamicRegistrationCapability        `json:"callHierarchy,omitempty"`
	SemanticTokens     *SemanticTokensClientCapabilities     `json:"semanticTokens,omitempty"`
	Moniker            *DynamicRegistrationCapability        `json:"moniker,omitempty"`
	TypeHierarchy      *DynamicRegistrationCapability        `json:"typeHierarchy,omitempty"`
	InlayHint          *InlayHintClientCapabilities          `json:"inlayHint,omitempty"`
	InlineValue        *DynamicRegistrationCapability        `json:"inlineValue,omitempty"`
	Diagnostic         *DiagnosticClientCapabilities         `json:"diagnostic,omitempty"`
}

// TextDocumentSyncClientCapabilities declares editor support for open/close/change/save document notifications.
type TextDocumentSyncClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	WillSave            *bool `json:"willSave,omitempty"`
	WillSaveWaitUntil   *bool `json:"willSaveWaitUntil,omitempty"`
	DidSave             *bool `json:"didSave,omitempty"`
}

// CompletionClientCapabilities declares which completion features the editor supports (snippets, commit characters, markdown docs, etc).
type CompletionClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	CompletionItem      *struct {
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
	ContextSupport *bool `json:"contextSupport,omitempty"`
}

// HoverClientCapabilities declares which content formats (plaintext, markdown) the editor supports in hover results.
type HoverClientCapabilities struct {
	DynamicRegistration *bool        `json:"dynamicRegistration,omitempty"`
	ContentFormat       []MarkupKind `json:"contentFormat,omitempty"`
}

// SignatureHelpClientCapabilities declares editor support for signature help features like active parameter highlighting and context tracking.
type SignatureHelpClientCapabilities struct {
	DynamicRegistration  *bool `json:"dynamicRegistration,omitempty"`
	SignatureInformation *struct {
		DocumentationFormat  []MarkupKind `json:"documentationFormat,omitempty"`
		ParameterInformation *struct {
			LabelOffsetSupport *bool `json:"labelOffsetSupport,omitempty"`
		} `json:"parameterInformation,omitempty"`
		ActiveParameterSupport *bool `json:"activeParameterSupport,omitempty"`
	} `json:"signatureInformation,omitempty"`
	ContextSupport *bool `json:"contextSupport,omitempty"`
}

// DeclarationClientCapabilities declares editor support for go-to-declaration, including dynamic registration and LocationLink responses.
type DeclarationClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	LinkSupport         *bool `json:"linkSupport,omitempty"`
}

// DefinitionClientCapabilities declares editor support for go-to-definition, including dynamic registration and LocationLink responses.
type DefinitionClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	LinkSupport         *bool `json:"linkSupport,omitempty"`
}

// TypeDefinitionClientCapabilities declares editor support for go-to-type-definition, including dynamic registration and LocationLink responses.
type TypeDefinitionClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	LinkSupport         *bool `json:"linkSupport,omitempty"`
}

// ImplementationClientCapabilities declares editor support for go-to-implementation, including dynamic registration and LocationLink responses.
type ImplementationClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	LinkSupport         *bool `json:"linkSupport,omitempty"`
}

// DocumentSymbolClientCapabilities declares editor support for document symbol features like hierarchical symbols and symbol tags.
type DocumentSymbolClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	SymbolKind          *struct {
		ValueSet []SymbolKind `json:"valueSet,omitempty"`
	} `json:"symbolKind,omitempty"`
	HierarchicalDocumentSymbolSupport *bool `json:"hierarchicalDocumentSymbolSupport,omitempty"`
	TagSupport                        *struct {
		ValueSet []SymbolTag `json:"valueSet,omitempty"`
	} `json:"tagSupport,omitempty"`
	LabelSupport *bool `json:"labelSupport,omitempty"`
}

// CodeActionClientCapabilities declares editor support for code action features like action kinds, lazy resolution, and disabled action display.
type CodeActionClientCapabilities struct {
	DynamicRegistration      *bool `json:"dynamicRegistration,omitempty"`
	CodeActionLiteralSupport *struct {
		CodeActionKind struct {
			ValueSet []CodeActionKind `json:"valueSet"`
		} `json:"codeActionKind"`
	} `json:"codeActionLiteralSupport,omitempty"`
	IsPreferredSupport *bool `json:"isPreferredSupport,omitempty"`
	DisabledSupport    *bool `json:"disabledSupport,omitempty"`
	DataSupport        *bool `json:"dataSupport,omitempty"`
	ResolveSupport     *struct {
		Properties []string `json:"properties"`
	} `json:"resolveSupport,omitempty"`
	HonorsChangeAnnotations *bool `json:"honorsChangeAnnotations,omitempty"`
}

// DocumentLinkClientCapabilities declares editor support for document link features like tooltip display.
type DocumentLinkClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	TooltipSupport      *bool `json:"tooltipSupport,omitempty"`
}

// RenameClientCapabilities declares editor support for rename features like prepare-rename and honoring change annotations.
type RenameClientCapabilities struct {
	DynamicRegistration           *bool                          `json:"dynamicRegistration,omitempty"`
	PrepareSupport                *bool                          `json:"prepareSupport,omitempty"`
	PrepareSupportDefaultBehavior *PrepareSupportDefaultBehavior `json:"prepareSupportDefaultBehavior,omitempty"`
	HonorsChangeAnnotations       *bool                          `json:"honorsChangeAnnotations,omitempty"`
}

// PublishDiagnosticsClientCapabilities declares editor support for diagnostic features like related information, tags (deprecated/unnecessary), and versioning.
type PublishDiagnosticsClientCapabilities struct {
	RelatedInformation *bool `json:"relatedInformation,omitempty"`
	TagSupport         *struct {
		ValueSet []DiagnosticTag `json:"valueSet,omitempty"`
	} `json:"tagSupport,omitempty"`
	VersionSupport         *bool `json:"versionSupport,omitempty"`
	CodeDescriptionSupport *bool `json:"codeDescriptionSupport,omitempty"`
	DataSupport            *bool `json:"dataSupport,omitempty"`
}

// FoldingRangeClientCapabilities declares editor support for folding ranges, including range limits and line-only folding.
type FoldingRangeClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	RangeLimit          *int  `json:"rangeLimit,omitempty"`
	LineFoldingOnly     *bool `json:"lineFoldingOnly,omitempty"`
}

// SemanticTokensClientCapabilities declares which semantic token types, modifiers, and request styles (full, delta, range) the editor supports.
type SemanticTokensClientCapabilities struct {
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	Requests            struct {
		Range *bool               `json:"range,omitempty"`
		Full  *SemanticTokensFull `json:"full,omitempty"`
	} `json:"requests"`
	TokenTypes              []string      `json:"tokenTypes"`
	TokenModifiers          []string      `json:"tokenModifiers"`
	Formats                 []TokenFormat `json:"formats"`
	OverlappingTokenSupport *bool         `json:"overlappingTokenSupport,omitempty"`
	MultilineTokenSupport   *bool         `json:"multilineTokenSupport,omitempty"`
}

// WindowClientCapabilities declares editor support for window features like work-done progress, show-message requests, and show-document.
type WindowClientCapabilities struct {
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
	ShowMessage      *struct {
		MessageActionItem *struct {
			AdditionalPropertiesSupport *bool `json:"additionalPropertiesSupport,omitempty"`
		} `json:"messageActionItem,omitempty"`
	} `json:"showMessage,omitempty"`
	ShowDocument *struct {
		Support bool `json:"support"`
	} `json:"showDocument,omitempty"`
}

// GeneralClientCapabilities declares cross-cutting editor capabilities like supported markdown parser, regex engine, and position encodings.
type GeneralClientCapabilities struct {
	RegularExpressions *struct {
		Engine  string `json:"engine"`
		Version string `json:"version,omitempty"`
	} `json:"regularExpressions,omitempty"`
	Markdown *struct {
		Parser      string   `json:"parser"`
		Version     string   `json:"version,omitempty"`
		AllowedTags []string `json:"allowedTags,omitempty"`
	} `json:"markdown,omitempty"`
}

// ServerCapabilities declares which LSP features this server supports, returned during initialization so the editor knows what to request.
type ServerCapabilities struct {
	TextDocumentSync                 *TextDocumentSyncOptions         `json:"textDocumentSync,omitempty"`
	CompletionProvider               *CompletionOptions               `json:"completionProvider,omitempty"`
	HoverProvider                    *bool                            `json:"hoverProvider,omitempty"`
	SignatureHelpProvider            *SignatureHelpOptions            `json:"signatureHelpProvider,omitempty"`
	DeclarationProvider              *bool                            `json:"declarationProvider,omitempty"`
	DefinitionProvider               *bool                            `json:"definitionProvider,omitempty"`
	TypeDefinitionProvider           *bool                            `json:"typeDefinitionProvider,omitempty"`
	ImplementationProvider           *bool                            `json:"implementationProvider,omitempty"`
	ReferencesProvider               *bool                            `json:"referencesProvider,omitempty"`
	DocumentHighlightProvider        *bool                            `json:"documentHighlightProvider,omitempty"`
	DocumentSymbolProvider           *bool                            `json:"documentSymbolProvider,omitempty"`
	CodeActionProvider               *CodeActionOptions               `json:"codeActionProvider,omitempty"`
	CodeLensProvider                 *CodeLensOptions                 `json:"codeLensProvider,omitempty"`
	DocumentLinkProvider             *DocumentLinkOptions             `json:"documentLinkProvider,omitempty"`
	ColorProvider                    *bool                            `json:"colorProvider,omitempty"`
	DocumentFormattingProvider       *bool                            `json:"documentFormattingProvider,omitempty"`
	DocumentRangeFormattingProvider  *bool                            `json:"documentRangeFormattingProvider,omitempty"`
	DocumentOnTypeFormattingProvider *DocumentOnTypeFormattingOptions `json:"documentOnTypeFormattingProvider,omitempty"`
	RenameProvider                   *RenameOptions                   `json:"renameProvider,omitempty"`
	FoldingRangeProvider             *bool                            `json:"foldingRangeProvider,omitempty"`
	ExecuteCommandProvider           *ExecuteCommandOptions           `json:"executeCommandProvider,omitempty"`
	SelectionRangeProvider           *bool                            `json:"selectionRangeProvider,omitempty"`
	LinkedEditingRangeProvider       *bool                            `json:"linkedEditingRangeProvider,omitempty"`
	CallHierarchyProvider            *bool                            `json:"callHierarchyProvider,omitempty"`
	SemanticTokensProvider           *SemanticTokensOptions           `json:"semanticTokensProvider,omitempty"`
	MonikerProvider                  *bool                            `json:"monikerProvider,omitempty"`
	TypeHierarchyProvider            *bool                            `json:"typeHierarchyProvider,omitempty"`
	InlayHintProvider                *InlayHintOptions                `json:"inlayHintProvider,omitempty"`
	InlineValueProvider              *bool                            `json:"inlineValueProvider,omitempty"`
	DiagnosticProvider               *DiagnosticOptions               `json:"diagnosticProvider,omitempty"`
	WorkspaceSymbolProvider          *bool                            `json:"workspaceSymbolProvider,omitempty"`
	Workspace                        *ServerWorkspaceCapabilities     `json:"workspace,omitempty"`
	Experimental                     json.RawMessage                  `json:"experimental,omitempty"`
}

// ServerWorkspaceCapabilities declares server support for workspace features like workspace folders and file operations.
type ServerWorkspaceCapabilities struct {
	WorkspaceFolders *WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
	FileOperations   *FileOperationOptions               `json:"fileOperations,omitempty"`
}

// WorkspaceFoldersServerCapabilities declares whether the server supports multi-root workspaces and wants workspace folder change notifications.
type WorkspaceFoldersServerCapabilities struct {
	Supported           *bool `json:"supported,omitempty"`
	ChangeNotifications *bool `json:"changeNotifications,omitempty"`
}

// FileOperationOptions defines options for file operations.
type FileOperationOptions struct {
	DidCreate  *FileOperationRegistrationOptions `json:"didCreate,omitempty"`
	WillCreate *FileOperationRegistrationOptions `json:"willCreate,omitempty"`
	DidRename  *FileOperationRegistrationOptions `json:"didRename,omitempty"`
	WillRename *FileOperationRegistrationOptions `json:"willRename,omitempty"`
	DidDelete  *FileOperationRegistrationOptions `json:"didDelete,omitempty"`
	WillDelete *FileOperationRegistrationOptions `json:"willDelete,omitempty"`
}

// FileOperationRegistrationOptions defines registration options for file operations.
type FileOperationRegistrationOptions struct {
	Filters []FileOperationFilter `json:"filters"`
}

// FileOperationFilter defines a filter for file operations.
type FileOperationFilter struct {
	Scheme  string               `json:"scheme,omitempty"`
	Pattern FileOperationPattern `json:"pattern"`
}

// FileOperationPattern defines a pattern for file operations.
type FileOperationPattern struct {
	Glob    string                       `json:"glob"`
	Matches *FileOperationPatternKind    `json:"matches,omitempty"`
	Options *FileOperationPatternOptions `json:"options,omitempty"`
}

// FileOperationPatternKind is a string enum ("file" or "folder") filtering which filesystem entries a file operation pattern matches.
type FileOperationPatternKind string

const (
	FileOperationPatternKindFile   FileOperationPatternKind = "file"
	FileOperationPatternKindFolder FileOperationPatternKind = "folder"
)

// FileOperationPatternOptions defines options for file operation patterns.
type FileOperationPatternOptions struct {
	IgnoreCase *bool `json:"ignoreCase,omitempty"`
}

// TextDocumentSyncOptions configures how the server receives document content: open/close notifications, incremental vs full change events, and save behavior.
type TextDocumentSyncOptions struct {
	OpenClose         *bool                `json:"openClose,omitempty"`
	Change            TextDocumentSyncKind `json:"change,omitempty"`
	WillSave          *bool                `json:"willSave,omitempty"`
	WillSaveWaitUntil *bool                `json:"willSaveWaitUntil,omitempty"`
	Save              *SaveOptions         `json:"save,omitempty"`
}

// SaveOptions configures whether the server receives document text on save.
type SaveOptions struct {
	IncludeText *bool `json:"includeText,omitempty"`
}

// CompletionOptions configures server completion behavior: trigger characters, resolve support, and work-done progress.
type CompletionOptions struct {
	WorkDoneProgressOptions
	TriggerCharacters   []string `json:"triggerCharacters,omitempty"`
	AllCommitCharacters []string `json:"allCommitCharacters,omitempty"`
	ResolveProvider     *bool    `json:"resolveProvider,omitempty"`
}

// SignatureHelpOptions configures which characters trigger signature help and which re-trigger it inside an active signature.
type SignatureHelpOptions struct {
	WorkDoneProgressOptions
	TriggerCharacters   []string `json:"triggerCharacters,omitempty"`
	RetriggerCharacters []string `json:"retriggerCharacters,omitempty"`
}

// CodeActionOptions configures which code action kinds the server supports and whether actions can be lazily resolved.
type CodeActionOptions struct {
	WorkDoneProgressOptions
	CodeActionKinds []CodeActionKind `json:"codeActionKinds,omitempty"`
	ResolveProvider *bool            `json:"resolveProvider,omitempty"`
}

// CodeLensOptions configures whether code lenses can be lazily resolved via codeLens/resolve.
type CodeLensOptions struct {
	WorkDoneProgressOptions
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

// DocumentLinkOptions configures whether document links can be lazily resolved via documentLink/resolve.
type DocumentLinkOptions struct {
	WorkDoneProgressOptions
	ResolveProvider *bool `json:"resolveProvider,omitempty"`
}

// DocumentOnTypeFormattingOptions configures which characters trigger on-type formatting (e.g. "}", ";", newline).
type DocumentOnTypeFormattingOptions struct {
	FirstTriggerCharacter string   `json:"firstTriggerCharacter"`
	MoreTriggerCharacter  []string `json:"moreTriggerCharacter,omitempty"`
}

// RenameOptions configures whether the server supports prepareRename for validating renames before applying them.
type RenameOptions struct {
	WorkDoneProgressOptions
	PrepareProvider *bool `json:"prepareProvider,omitempty"`
}

// ExecuteCommandOptions lists the command IDs the server can execute via workspace/executeCommand.
type ExecuteCommandOptions struct {
	WorkDoneProgressOptions
	Commands []string `json:"commands"`
}
