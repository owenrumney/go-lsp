package server

import (
	"context"
	"encoding/json"

	"github.com/owenrumney/go-lsp/lsp"
)

// LifecycleHandler must be implemented by all servers.
type LifecycleHandler interface {
	Initialize(ctx context.Context, params *lsp.InitializeParams) (*lsp.InitializeResult, error)
	Shutdown(ctx context.Context) error
}

// SetTraceHandler handles $/setTrace notifications.
type SetTraceHandler interface {
	SetTrace(ctx context.Context, params *lsp.SetTraceParams) error
}

// TextDocumentSyncHandler handles document open/change/close notifications.
type TextDocumentSyncHandler interface {
	DidOpen(ctx context.Context, params *lsp.DidOpenTextDocumentParams) error
	DidChange(ctx context.Context, params *lsp.DidChangeTextDocumentParams) error
	DidClose(ctx context.Context, params *lsp.DidCloseTextDocumentParams) error
}

// TextDocumentSaveHandler handles document save notifications.
type TextDocumentSaveHandler interface {
	DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error
}

// CompletionHandler handles textDocument/completion.
type CompletionHandler interface {
	Completion(ctx context.Context, params *lsp.CompletionParams) (*lsp.CompletionList, error)
}

// CompletionResolveHandler handles completionItem/resolve.
type CompletionResolveHandler interface {
	ResolveCompletionItem(ctx context.Context, params *lsp.CompletionItem) (*lsp.CompletionItem, error)
}

// HoverHandler handles textDocument/hover.
type HoverHandler interface {
	Hover(ctx context.Context, params *lsp.HoverParams) (*lsp.Hover, error)
}

// SignatureHelpHandler handles textDocument/signatureHelp.
type SignatureHelpHandler interface {
	SignatureHelp(ctx context.Context, params *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error)
}

// DeclarationHandler handles textDocument/declaration.
type DeclarationHandler interface {
	Declaration(ctx context.Context, params *lsp.DeclarationParams) ([]lsp.Location, error)
}

// DefinitionHandler handles textDocument/definition.
type DefinitionHandler interface {
	Definition(ctx context.Context, params *lsp.DefinitionParams) ([]lsp.Location, error)
}

// TypeDefinitionHandler handles textDocument/typeDefinition.
type TypeDefinitionHandler interface {
	TypeDefinition(ctx context.Context, params *lsp.TypeDefinitionParams) ([]lsp.Location, error)
}

// ImplementationHandler handles textDocument/implementation.
type ImplementationHandler interface {
	Implementation(ctx context.Context, params *lsp.ImplementationParams) ([]lsp.Location, error)
}

// ReferencesHandler handles textDocument/references.
type ReferencesHandler interface {
	References(ctx context.Context, params *lsp.ReferenceParams) ([]lsp.Location, error)
}

// DocumentHighlightHandler handles textDocument/documentHighlight.
type DocumentHighlightHandler interface {
	DocumentHighlight(ctx context.Context, params *lsp.DocumentHighlightParams) ([]lsp.DocumentHighlight, error)
}

// DocumentSymbolHandler handles textDocument/documentSymbol.
type DocumentSymbolHandler interface {
	DocumentSymbol(ctx context.Context, params *lsp.DocumentSymbolParams) ([]lsp.DocumentSymbol, error)
}

// CodeActionHandler handles textDocument/codeAction.
type CodeActionHandler interface {
	CodeAction(ctx context.Context, params *lsp.CodeActionParams) ([]lsp.CodeAction, error)
}

// CodeLensHandler handles textDocument/codeLens.
type CodeLensHandler interface {
	CodeLens(ctx context.Context, params *lsp.CodeLensParams) ([]lsp.CodeLens, error)
}

// CodeLensResolveHandler handles codeLens/resolve.
type CodeLensResolveHandler interface {
	ResolveCodeLens(ctx context.Context, params *lsp.CodeLens) (*lsp.CodeLens, error)
}

// DocumentLinkHandler handles textDocument/documentLink.
type DocumentLinkHandler interface {
	DocumentLink(ctx context.Context, params *lsp.DocumentLinkParams) ([]lsp.DocumentLink, error)
}

// DocumentColorHandler handles textDocument/documentColor.
type DocumentColorHandler interface {
	DocumentColor(ctx context.Context, params *lsp.DocumentColorParams) ([]lsp.ColorInformation, error)
}

// ColorPresentationHandler handles textDocument/colorPresentation.
type ColorPresentationHandler interface {
	ColorPresentation(ctx context.Context, params *lsp.ColorPresentationParams) ([]lsp.ColorPresentation, error)
}

// DocumentFormattingHandler handles textDocument/formatting.
type DocumentFormattingHandler interface {
	Formatting(ctx context.Context, params *lsp.DocumentFormattingParams) ([]lsp.TextEdit, error)
}

// DocumentRangeFormattingHandler handles textDocument/rangeFormatting.
type DocumentRangeFormattingHandler interface {
	RangeFormatting(ctx context.Context, params *lsp.DocumentRangeFormattingParams) ([]lsp.TextEdit, error)
}

// DocumentOnTypeFormattingHandler handles textDocument/onTypeFormatting.
type DocumentOnTypeFormattingHandler interface {
	OnTypeFormatting(ctx context.Context, params *lsp.DocumentOnTypeFormattingParams) ([]lsp.TextEdit, error)
}

// RenameHandler handles textDocument/rename.
type RenameHandler interface {
	Rename(ctx context.Context, params *lsp.RenameParams) (*lsp.WorkspaceEdit, error)
}

// PrepareRenameHandler handles textDocument/prepareRename.
type PrepareRenameHandler interface {
	PrepareRename(ctx context.Context, params *lsp.PrepareRenameParams) (*lsp.PrepareRenameResult, error)
}

// FoldingRangeHandler handles textDocument/foldingRange.
type FoldingRangeHandler interface {
	FoldingRange(ctx context.Context, params *lsp.FoldingRangeParams) ([]lsp.FoldingRange, error)
}

// SelectionRangeHandler handles textDocument/selectionRange.
type SelectionRangeHandler interface {
	SelectionRange(ctx context.Context, params *lsp.SelectionRangeParams) ([]lsp.SelectionRange, error)
}

// CallHierarchyHandler handles textDocument/prepareCallHierarchy.
type CallHierarchyHandler interface {
	PrepareCallHierarchy(ctx context.Context, params *lsp.CallHierarchyPrepareParams) ([]lsp.CallHierarchyItem, error)
	IncomingCalls(ctx context.Context, params *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error)
	OutgoingCalls(ctx context.Context, params *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error)
}

// SemanticTokensFullHandler handles textDocument/semanticTokens/full.
type SemanticTokensFullHandler interface {
	SemanticTokensFull(ctx context.Context, params *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error)
}

// SemanticTokensDeltaHandler handles textDocument/semanticTokens/full/delta.
type SemanticTokensDeltaHandler interface {
	SemanticTokensDelta(ctx context.Context, params *lsp.SemanticTokensDeltaParams) (*lsp.SemanticTokensDelta, error)
}

// SemanticTokensRangeHandler handles textDocument/semanticTokens/range.
type SemanticTokensRangeHandler interface {
	SemanticTokensRange(ctx context.Context, params *lsp.SemanticTokensRangeParams) (*lsp.SemanticTokens, error)
}

// LinkedEditingRangeHandler handles textDocument/linkedEditingRange.
type LinkedEditingRangeHandler interface {
	LinkedEditingRange(ctx context.Context, params *lsp.LinkedEditingRangeParams) (*lsp.LinkedEditingRanges, error)
}

// MonikerHandler handles textDocument/moniker.
type MonikerHandler interface {
	Moniker(ctx context.Context, params *lsp.MonikerParams) ([]lsp.Moniker, error)
}

// TypeHierarchyHandler handles textDocument/prepareTypeHierarchy and related methods.
type TypeHierarchyHandler interface {
	PrepareTypeHierarchy(ctx context.Context, params *lsp.TypeHierarchyPrepareParams) ([]lsp.TypeHierarchyItem, error)
	Supertypes(ctx context.Context, params *lsp.TypeHierarchySupertypesParams) ([]lsp.TypeHierarchyItem, error)
	Subtypes(ctx context.Context, params *lsp.TypeHierarchySubtypesParams) ([]lsp.TypeHierarchyItem, error)
}

// InlayHintHandler handles textDocument/inlayHint.
type InlayHintHandler interface {
	InlayHint(ctx context.Context, params *lsp.InlayHintParams) ([]lsp.InlayHint, error)
}

// InlayHintResolveHandler handles inlayHint/resolve.
type InlayHintResolveHandler interface {
	ResolveInlayHint(ctx context.Context, params *lsp.InlayHint) (*lsp.InlayHint, error)
}

// InlineValueHandler handles textDocument/inlineValue.
type InlineValueHandler interface {
	InlineValue(ctx context.Context, params *lsp.InlineValueParams) ([]json.RawMessage, error)
}

// DocumentDiagnosticHandler handles textDocument/diagnostic.
type DocumentDiagnosticHandler interface {
	DocumentDiagnostic(ctx context.Context, params *lsp.DocumentDiagnosticParams) (any, error)
}

// WorkspaceDiagnosticHandler handles workspace/diagnostic.
type WorkspaceDiagnosticHandler interface {
	WorkspaceDiagnostic(ctx context.Context, params *lsp.WorkspaceDiagnosticParams) (*lsp.WorkspaceDiagnosticReport, error)
}

// WorkspaceSymbolHandler handles workspace/symbol.
type WorkspaceSymbolHandler interface {
	WorkspaceSymbol(ctx context.Context, params *lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error)
}

// ExecuteCommandHandler handles workspace/executeCommand.
type ExecuteCommandHandler interface {
	ExecuteCommand(ctx context.Context, params *lsp.ExecuteCommandParams) (any, error)
}

// WorkspaceFoldersHandler handles workspace folder notifications.
type WorkspaceFoldersHandler interface {
	DidChangeWorkspaceFolders(ctx context.Context, params *lsp.DidChangeWorkspaceFoldersParams) error
}

// DidChangeConfigurationHandler handles workspace/didChangeConfiguration.
type DidChangeConfigurationHandler interface {
	DidChangeConfiguration(ctx context.Context, params *lsp.DidChangeConfigurationParams) error
}

// DidChangeWatchedFilesHandler handles workspace/didChangeWatchedFiles.
type DidChangeWatchedFilesHandler interface {
	DidChangeWatchedFiles(ctx context.Context, params *lsp.DidChangeWatchedFilesParams) error
}

// TextDocumentWillSaveHandler handles textDocument/willSave notifications.
type TextDocumentWillSaveHandler interface {
	WillSave(ctx context.Context, params *lsp.WillSaveTextDocumentParams) error
}

// TextDocumentWillSaveWaitUntilHandler handles textDocument/willSaveWaitUntil requests.
type TextDocumentWillSaveWaitUntilHandler interface {
	WillSaveWaitUntil(ctx context.Context, params *lsp.WillSaveTextDocumentParams) ([]lsp.TextEdit, error)
}

// CodeActionResolveHandler handles codeAction/resolve.
type CodeActionResolveHandler interface {
	ResolveCodeAction(ctx context.Context, params *lsp.CodeAction) (*lsp.CodeAction, error)
}

// DocumentLinkResolveHandler handles documentLink/resolve.
type DocumentLinkResolveHandler interface {
	ResolveDocumentLink(ctx context.Context, params *lsp.DocumentLink) (*lsp.DocumentLink, error)
}

// WillCreateFilesHandler handles workspace/willCreateFiles.
type WillCreateFilesHandler interface {
	WillCreateFiles(ctx context.Context, params *lsp.CreateFilesParams) (*lsp.WorkspaceEdit, error)
}

// WillRenameFilesHandler handles workspace/willRenameFiles.
type WillRenameFilesHandler interface {
	WillRenameFiles(ctx context.Context, params *lsp.RenameFilesParams) (*lsp.WorkspaceEdit, error)
}

// WillDeleteFilesHandler handles workspace/willDeleteFiles.
type WillDeleteFilesHandler interface {
	WillDeleteFiles(ctx context.Context, params *lsp.DeleteFilesParams) (*lsp.WorkspaceEdit, error)
}
