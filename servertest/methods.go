package servertest

import (
	"encoding/json"

	"github.com/owenrumney/go-lsp/lsp"
)

// DidOpen sends a textDocument/didOpen notification.
func (h *Harness) DidOpen(uri lsp.DocumentURI, languageID, text string) error {
	h.versionsMu.Lock()
	h.versions[uri] = 1
	h.versionsMu.Unlock()

	return h.conn.notify(h.ctx, "textDocument/didOpen", &lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: languageID,
			Version:    1,
			Text:       text,
		},
	})
}

// DidChange sends a textDocument/didChange notification with full document sync.
func (h *Harness) DidChange(uri lsp.DocumentURI, version int, text string) error {
	h.versionsMu.Lock()
	h.versions[uri] = version
	h.versionsMu.Unlock()

	return h.conn.notify(h.ctx, "textDocument/didChange", &lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{URI: uri},
			Version:                version,
		},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{
			{Text: text},
		},
	})
}

// DidSave sends a textDocument/didSave notification.
func (h *Harness) DidSave(uri lsp.DocumentURI) error {
	return h.conn.notify(h.ctx, "textDocument/didSave", &lsp.DidSaveTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// DidClose sends a textDocument/didClose notification.
func (h *Harness) DidClose(uri lsp.DocumentURI) error {
	h.versionsMu.Lock()
	delete(h.versions, uri)
	h.versionsMu.Unlock()

	return h.conn.notify(h.ctx, "textDocument/didClose", &lsp.DidCloseTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// WillSave sends a textDocument/willSave notification.
func (h *Harness) WillSave(uri lsp.DocumentURI, reason lsp.TextDocumentSaveReason) error {
	return h.conn.notify(h.ctx, "textDocument/willSave", &lsp.WillSaveTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Reason:       reason,
	})
}

// DidChangeWorkspaceFolders sends a workspace/didChangeWorkspaceFolders notification.
func (h *Harness) DidChangeWorkspaceFolders(params *lsp.DidChangeWorkspaceFoldersParams) error {
	return h.conn.notify(h.ctx, "workspace/didChangeWorkspaceFolders", params)
}

// DidChangeConfiguration sends a workspace/didChangeConfiguration notification.
func (h *Harness) DidChangeConfiguration(params *lsp.DidChangeConfigurationParams) error {
	return h.conn.notify(h.ctx, "workspace/didChangeConfiguration", params)
}

// DidChangeWatchedFiles sends a workspace/didChangeWatchedFiles notification.
func (h *Harness) DidChangeWatchedFiles(params *lsp.DidChangeWatchedFilesParams) error {
	return h.conn.notify(h.ctx, "workspace/didChangeWatchedFiles", params)
}

// SetTrace sends a $/setTrace notification.
func (h *Harness) SetTrace(value lsp.TraceValue) error {
	return h.conn.notify(h.ctx, "$/setTrace", &lsp.SetTraceParams{Value: value})
}

// Hover sends a textDocument/hover request.
func (h *Harness) Hover(uri lsp.DocumentURI, line, char int) (*lsp.Hover, error) {
	return callPtr[lsp.Hover](h, "textDocument/hover", &lsp.HoverParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Position:     lsp.Position{Line: line, Character: char},
		},
	})
}

// Completion sends a textDocument/completion request.
func (h *Harness) Completion(uri lsp.DocumentURI, line, char int) (*lsp.CompletionList, error) {
	return callPtr[lsp.CompletionList](h, "textDocument/completion", &lsp.CompletionParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Position:     lsp.Position{Line: line, Character: char},
		},
	})
}

// ResolveCompletionItem sends a completionItem/resolve request.
func (h *Harness) ResolveCompletionItem(item *lsp.CompletionItem) (*lsp.CompletionItem, error) {
	return callPtr[lsp.CompletionItem](h, "completionItem/resolve", item)
}

// SignatureHelp sends a textDocument/signatureHelp request.
func (h *Harness) SignatureHelp(uri lsp.DocumentURI, line, char int) (*lsp.SignatureHelp, error) {
	return callPtr[lsp.SignatureHelp](h, "textDocument/signatureHelp", &lsp.SignatureHelpParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// Declaration sends a textDocument/declaration request.
func (h *Harness) Declaration(uri lsp.DocumentURI, line, char int) ([]lsp.Location, error) {
	return callValue[[]lsp.Location](h, "textDocument/declaration", &lsp.DeclarationParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// Definition sends a textDocument/definition request.
func (h *Harness) Definition(uri lsp.DocumentURI, line, char int) ([]lsp.Location, error) {
	return callValue[[]lsp.Location](h, "textDocument/definition", &lsp.DefinitionParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// TypeDefinition sends a textDocument/typeDefinition request.
func (h *Harness) TypeDefinition(uri lsp.DocumentURI, line, char int) ([]lsp.Location, error) {
	return callValue[[]lsp.Location](h, "textDocument/typeDefinition", &lsp.TypeDefinitionParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// Implementation sends a textDocument/implementation request.
func (h *Harness) Implementation(uri lsp.DocumentURI, line, char int) ([]lsp.Location, error) {
	return callValue[[]lsp.Location](h, "textDocument/implementation", &lsp.ImplementationParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// References sends a textDocument/references request.
func (h *Harness) References(uri lsp.DocumentURI, line, char int, includeDecl bool) ([]lsp.Location, error) {
	return callValue[[]lsp.Location](h, "textDocument/references", &lsp.ReferenceParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
		Context:                    lsp.ReferenceContext{IncludeDeclaration: includeDecl},
	})
}

// DocumentHighlight sends a textDocument/documentHighlight request.
func (h *Harness) DocumentHighlight(uri lsp.DocumentURI, line, char int) ([]lsp.DocumentHighlight, error) {
	return callValue[[]lsp.DocumentHighlight](h, "textDocument/documentHighlight", &lsp.DocumentHighlightParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// CodeAction sends a textDocument/codeAction request.
func (h *Harness) CodeAction(params *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	return callValue[[]lsp.CodeAction](h, "textDocument/codeAction", params)
}

// ResolveCodeAction sends a codeAction/resolve request.
func (h *Harness) ResolveCodeAction(action *lsp.CodeAction) (*lsp.CodeAction, error) {
	return callPtr[lsp.CodeAction](h, "codeAction/resolve", action)
}

// CodeLens sends a textDocument/codeLens request.
func (h *Harness) CodeLens(uri lsp.DocumentURI) ([]lsp.CodeLens, error) {
	return callValue[[]lsp.CodeLens](h, "textDocument/codeLens", &lsp.CodeLensParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// ResolveCodeLens sends a codeLens/resolve request.
func (h *Harness) ResolveCodeLens(lens *lsp.CodeLens) (*lsp.CodeLens, error) {
	return callPtr[lsp.CodeLens](h, "codeLens/resolve", lens)
}

// DocumentLink sends a textDocument/documentLink request.
func (h *Harness) DocumentLink(uri lsp.DocumentURI) ([]lsp.DocumentLink, error) {
	return callValue[[]lsp.DocumentLink](h, "textDocument/documentLink", &lsp.DocumentLinkParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// ResolveDocumentLink sends a documentLink/resolve request.
func (h *Harness) ResolveDocumentLink(link *lsp.DocumentLink) (*lsp.DocumentLink, error) {
	return callPtr[lsp.DocumentLink](h, "documentLink/resolve", link)
}

// DocumentColor sends a textDocument/documentColor request.
func (h *Harness) DocumentColor(uri lsp.DocumentURI) ([]lsp.ColorInformation, error) {
	return callValue[[]lsp.ColorInformation](h, "textDocument/documentColor", &lsp.DocumentColorParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// ColorPresentation sends a textDocument/colorPresentation request.
func (h *Harness) ColorPresentation(params *lsp.ColorPresentationParams) ([]lsp.ColorPresentation, error) {
	return callValue[[]lsp.ColorPresentation](h, "textDocument/colorPresentation", params)
}

// DocumentSymbol sends a textDocument/documentSymbol request.
func (h *Harness) DocumentSymbol(uri lsp.DocumentURI) ([]lsp.DocumentSymbol, error) {
	return callValue[[]lsp.DocumentSymbol](h, "textDocument/documentSymbol", &lsp.DocumentSymbolParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// WorkspaceSymbol sends a workspace/symbol request.
func (h *Harness) WorkspaceSymbol(query string) ([]lsp.SymbolInformation, error) {
	return callValue[[]lsp.SymbolInformation](h, "workspace/symbol", &lsp.WorkspaceSymbolParams{
		Query: query,
	})
}

// Formatting sends a textDocument/formatting request.
func (h *Harness) Formatting(uri lsp.DocumentURI) ([]lsp.TextEdit, error) {
	return callValue[[]lsp.TextEdit](h, "textDocument/formatting", &lsp.DocumentFormattingParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Options:      defaultFormattingOptions(),
	})
}

// RangeFormatting sends a textDocument/rangeFormatting request.
func (h *Harness) RangeFormatting(uri lsp.DocumentURI, r lsp.Range) ([]lsp.TextEdit, error) {
	return callValue[[]lsp.TextEdit](h, "textDocument/rangeFormatting", &lsp.DocumentRangeFormattingParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Range:        r,
		Options:      defaultFormattingOptions(),
	})
}

// OnTypeFormatting sends a textDocument/onTypeFormatting request.
func (h *Harness) OnTypeFormatting(uri lsp.DocumentURI, line, char int, typed string) ([]lsp.TextEdit, error) {
	return callValue[[]lsp.TextEdit](h, "textDocument/onTypeFormatting", &lsp.DocumentOnTypeFormattingParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
		Character:                  typed,
		Options:                    defaultFormattingOptions(),
	})
}

// Rename sends a textDocument/rename request.
func (h *Harness) Rename(uri lsp.DocumentURI, line, char int, newName string) (*lsp.WorkspaceEdit, error) {
	return callPtr[lsp.WorkspaceEdit](h, "textDocument/rename", &lsp.RenameParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
		NewName:                    newName,
	})
}

// PrepareRename sends a textDocument/prepareRename request.
func (h *Harness) PrepareRename(uri lsp.DocumentURI, line, char int) (*lsp.PrepareRenameResult, error) {
	return callPtr[lsp.PrepareRenameResult](h, "textDocument/prepareRename", &lsp.PrepareRenameParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// FoldingRange sends a textDocument/foldingRange request.
func (h *Harness) FoldingRange(uri lsp.DocumentURI) ([]lsp.FoldingRange, error) {
	return callValue[[]lsp.FoldingRange](h, "textDocument/foldingRange", &lsp.FoldingRangeParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// SelectionRange sends a textDocument/selectionRange request.
func (h *Harness) SelectionRange(uri lsp.DocumentURI, positions []lsp.Position) ([]lsp.SelectionRange, error) {
	return callValue[[]lsp.SelectionRange](h, "textDocument/selectionRange", &lsp.SelectionRangeParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Positions:    positions,
	})
}

// LinkedEditingRange sends a textDocument/linkedEditingRange request.
func (h *Harness) LinkedEditingRange(uri lsp.DocumentURI, line, char int) (*lsp.LinkedEditingRanges, error) {
	return callPtr[lsp.LinkedEditingRanges](h, "textDocument/linkedEditingRange", &lsp.LinkedEditingRangeParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// Moniker sends a textDocument/moniker request.
func (h *Harness) Moniker(uri lsp.DocumentURI, line, char int) ([]lsp.Moniker, error) {
	return callValue[[]lsp.Moniker](h, "textDocument/moniker", &lsp.MonikerParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// WillSaveWaitUntil sends a textDocument/willSaveWaitUntil request.
func (h *Harness) WillSaveWaitUntil(uri lsp.DocumentURI, reason lsp.TextDocumentSaveReason) ([]lsp.TextEdit, error) {
	return callValue[[]lsp.TextEdit](h, "textDocument/willSaveWaitUntil", &lsp.WillSaveTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Reason:       reason,
	})
}

// InlayHint sends a textDocument/inlayHint request.
func (h *Harness) InlayHint(uri lsp.DocumentURI, r lsp.Range) ([]lsp.InlayHint, error) {
	return callValue[[]lsp.InlayHint](h, "textDocument/inlayHint", &lsp.InlayHintParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Range:        r,
	})
}

// ResolveInlayHint sends an inlayHint/resolve request.
func (h *Harness) ResolveInlayHint(hint *lsp.InlayHint) (*lsp.InlayHint, error) {
	return callPtr[lsp.InlayHint](h, "inlayHint/resolve", hint)
}

// InlineValue sends a textDocument/inlineValue request.
func (h *Harness) InlineValue(params *lsp.InlineValueParams) ([]json.RawMessage, error) {
	return callValue[[]json.RawMessage](h, "textDocument/inlineValue", params)
}

// DocumentDiagnostic sends a textDocument/diagnostic request and returns the raw report.
func (h *Harness) DocumentDiagnostic(params *lsp.DocumentDiagnosticParams) (json.RawMessage, error) {
	return h.conn.call(h.ctx, "textDocument/diagnostic", params)
}

// WorkspaceDiagnostic sends a workspace/diagnostic request.
func (h *Harness) WorkspaceDiagnostic(params *lsp.WorkspaceDiagnosticParams) (*lsp.WorkspaceDiagnosticReport, error) {
	return callPtr[lsp.WorkspaceDiagnosticReport](h, "workspace/diagnostic", params)
}

// SemanticTokensFull sends a textDocument/semanticTokens/full request.
func (h *Harness) SemanticTokensFull(uri lsp.DocumentURI) (*lsp.SemanticTokens, error) {
	return callPtr[lsp.SemanticTokens](h, "textDocument/semanticTokens/full", &lsp.SemanticTokensParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// SemanticTokensDelta sends a textDocument/semanticTokens/full/delta request.
func (h *Harness) SemanticTokensDelta(uri lsp.DocumentURI, previousResultID string) (*lsp.SemanticTokensDelta, error) {
	return callPtr[lsp.SemanticTokensDelta](h, "textDocument/semanticTokens/full/delta", &lsp.SemanticTokensDeltaParams{
		TextDocument:     lsp.TextDocumentIdentifier{URI: uri},
		PreviousResultID: previousResultID,
	})
}

// SemanticTokensRange sends a textDocument/semanticTokens/range request.
func (h *Harness) SemanticTokensRange(uri lsp.DocumentURI, r lsp.Range) (*lsp.SemanticTokens, error) {
	return callPtr[lsp.SemanticTokens](h, "textDocument/semanticTokens/range", &lsp.SemanticTokensRangeParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Range:        r,
	})
}

// PrepareCallHierarchy sends a textDocument/prepareCallHierarchy request.
func (h *Harness) PrepareCallHierarchy(uri lsp.DocumentURI, line, char int) ([]lsp.CallHierarchyItem, error) {
	return callValue[[]lsp.CallHierarchyItem](h, "textDocument/prepareCallHierarchy", &lsp.CallHierarchyPrepareParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// IncomingCalls sends a callHierarchy/incomingCalls request.
func (h *Harness) IncomingCalls(item lsp.CallHierarchyItem) ([]lsp.CallHierarchyIncomingCall, error) {
	return callValue[[]lsp.CallHierarchyIncomingCall](h, "callHierarchy/incomingCalls", &lsp.CallHierarchyIncomingCallsParams{Item: item})
}

// OutgoingCalls sends a callHierarchy/outgoingCalls request.
func (h *Harness) OutgoingCalls(item lsp.CallHierarchyItem) ([]lsp.CallHierarchyOutgoingCall, error) {
	return callValue[[]lsp.CallHierarchyOutgoingCall](h, "callHierarchy/outgoingCalls", &lsp.CallHierarchyOutgoingCallsParams{Item: item})
}

// PrepareTypeHierarchy sends a textDocument/prepareTypeHierarchy request.
func (h *Harness) PrepareTypeHierarchy(uri lsp.DocumentURI, line, char int) ([]lsp.TypeHierarchyItem, error) {
	return callValue[[]lsp.TypeHierarchyItem](h, "textDocument/prepareTypeHierarchy", &lsp.TypeHierarchyPrepareParams{
		TextDocumentPositionParams: textDocumentPosition(uri, line, char),
	})
}

// Supertypes sends a typeHierarchy/supertypes request.
func (h *Harness) Supertypes(item lsp.TypeHierarchyItem) ([]lsp.TypeHierarchyItem, error) {
	return callValue[[]lsp.TypeHierarchyItem](h, "typeHierarchy/supertypes", &lsp.TypeHierarchySupertypesParams{Item: item})
}

// Subtypes sends a typeHierarchy/subtypes request.
func (h *Harness) Subtypes(item lsp.TypeHierarchyItem) ([]lsp.TypeHierarchyItem, error) {
	return callValue[[]lsp.TypeHierarchyItem](h, "typeHierarchy/subtypes", &lsp.TypeHierarchySubtypesParams{Item: item})
}

// ExecuteCommand sends a workspace/executeCommand request and returns the raw result.
func (h *Harness) ExecuteCommand(command string, args []json.RawMessage) (json.RawMessage, error) {
	return h.conn.call(h.ctx, "workspace/executeCommand", &lsp.ExecuteCommandParams{
		Command:   command,
		Arguments: args,
	})
}

// WillCreateFiles sends a workspace/willCreateFiles request.
func (h *Harness) WillCreateFiles(files []lsp.FileCreate) (*lsp.WorkspaceEdit, error) {
	return callPtr[lsp.WorkspaceEdit](h, "workspace/willCreateFiles", &lsp.CreateFilesParams{Files: files})
}

// WillRenameFiles sends a workspace/willRenameFiles request.
func (h *Harness) WillRenameFiles(files []lsp.FileRename) (*lsp.WorkspaceEdit, error) {
	return callPtr[lsp.WorkspaceEdit](h, "workspace/willRenameFiles", &lsp.RenameFilesParams{Files: files})
}

// WillDeleteFiles sends a workspace/willDeleteFiles request.
func (h *Harness) WillDeleteFiles(files []lsp.FileDelete) (*lsp.WorkspaceEdit, error) {
	return callPtr[lsp.WorkspaceEdit](h, "workspace/willDeleteFiles", &lsp.DeleteFilesParams{Files: files})
}

// Call sends an arbitrary JSON-RPC request and returns the raw result.
func (h *Harness) Call(method string, params any) (json.RawMessage, error) {
	return h.conn.call(h.ctx, method, params)
}

// CallAsync sends an arbitrary JSON-RPC request and returns before the response arrives.
func (h *Harness) CallAsync(method string, params any) (*PendingCall, error) {
	return h.conn.startCall(method, params)
}

// Notify sends an arbitrary JSON-RPC notification.
func (h *Harness) Notify(method string, params any) error {
	return h.conn.notify(h.ctx, method, params)
}

// CancelRequest sends a $/cancelRequest notification for an in-flight request ID.
func (h *Harness) CancelRequest(id int64) error {
	return h.Notify("$/cancelRequest", map[string]any{"id": id})
}

func textDocumentPosition(uri lsp.DocumentURI, line, char int) lsp.TextDocumentPositionParams {
	return lsp.TextDocumentPositionParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Position:     lsp.Position{Line: line, Character: char},
	}
}

func defaultFormattingOptions() lsp.FormattingOptions {
	return lsp.FormattingOptions{TabSize: 4, InsertSpaces: true}
}

func callPtr[T any](h *Harness, method string, params any) (*T, error) {
	result, err := h.conn.call(h.ctx, method, params)
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var value T
	if err := json.Unmarshal(result, &value); err != nil {
		return nil, err
	}
	return &value, nil
}

func callValue[T any](h *Harness, method string, params any) (T, error) {
	var value T
	result, err := h.conn.call(h.ctx, method, params)
	if err != nil {
		return value, err
	}
	if string(result) == "null" {
		return value, nil
	}
	if err := json.Unmarshal(result, &value); err != nil {
		return value, err
	}
	return value, nil
}
