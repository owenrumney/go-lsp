package server

import "github.com/owenrumney/go-lsp/lsp"

var enabled = true

func buildCapabilities(handler any) lsp.ServerCapabilities {
	var caps lsp.ServerCapabilities

	// Text document sync
	if _, ok := handler.(TextDocumentSyncHandler); ok {
		openClose := true
		caps.TextDocumentSync = &lsp.TextDocumentSyncOptions{
			OpenClose: &openClose,
			Change:    lsp.SyncIncremental,
		}
		if _, ok := handler.(TextDocumentSaveHandler); ok {
			caps.TextDocumentSync.Save = &lsp.SaveOptions{IncludeText: &enabled}
		}
		if _, ok := handler.(TextDocumentWillSaveHandler); ok {
			caps.TextDocumentSync.WillSave = &enabled
		}
		if _, ok := handler.(TextDocumentWillSaveWaitUntilHandler); ok {
			caps.TextDocumentSync.WillSaveWaitUntil = &enabled
		}
	}

	if _, ok := handler.(CompletionHandler); ok {
		opts := &lsp.CompletionOptions{}
		if _, ok := handler.(CompletionResolveHandler); ok {
			opts.ResolveProvider = &enabled
		}
		caps.CompletionProvider = opts
	}

	if _, ok := handler.(HoverHandler); ok {
		caps.HoverProvider = &enabled
	}

	if _, ok := handler.(SignatureHelpHandler); ok {
		caps.SignatureHelpProvider = &lsp.SignatureHelpOptions{}
	}

	if _, ok := handler.(DeclarationHandler); ok {
		caps.DeclarationProvider = &enabled
	}

	if _, ok := handler.(DefinitionHandler); ok {
		caps.DefinitionProvider = &enabled
	}

	if _, ok := handler.(TypeDefinitionHandler); ok {
		caps.TypeDefinitionProvider = &enabled
	}

	if _, ok := handler.(ImplementationHandler); ok {
		caps.ImplementationProvider = &enabled
	}

	if _, ok := handler.(ReferencesHandler); ok {
		caps.ReferencesProvider = &enabled
	}

	if _, ok := handler.(DocumentHighlightHandler); ok {
		caps.DocumentHighlightProvider = &enabled
	}

	if _, ok := handler.(DocumentSymbolHandler); ok {
		caps.DocumentSymbolProvider = &enabled
	}

	if _, ok := handler.(CodeActionHandler); ok {
		opts := &lsp.CodeActionOptions{}
		if _, ok := handler.(CodeActionResolveHandler); ok {
			opts.ResolveProvider = &enabled
		}
		caps.CodeActionProvider = opts
	}

	if _, ok := handler.(CodeLensHandler); ok {
		opts := &lsp.CodeLensOptions{}
		if _, ok := handler.(CodeLensResolveHandler); ok {
			opts.ResolveProvider = &enabled
		}
		caps.CodeLensProvider = opts
	}

	if _, ok := handler.(DocumentLinkHandler); ok {
		opts := &lsp.DocumentLinkOptions{}
		if _, ok := handler.(DocumentLinkResolveHandler); ok {
			opts.ResolveProvider = &enabled
		}
		caps.DocumentLinkProvider = opts
	}

	if _, ok := handler.(DocumentColorHandler); ok {
		caps.ColorProvider = &enabled
	}

	if _, ok := handler.(DocumentFormattingHandler); ok {
		caps.DocumentFormattingProvider = &enabled
	}

	if _, ok := handler.(DocumentRangeFormattingHandler); ok {
		caps.DocumentRangeFormattingProvider = &enabled
	}

	if _, ok := handler.(RenameHandler); ok {
		opts := &lsp.RenameOptions{}
		if _, ok := handler.(PrepareRenameHandler); ok {
			opts.PrepareProvider = &enabled
		}
		caps.RenameProvider = opts
	}

	if _, ok := handler.(FoldingRangeHandler); ok {
		caps.FoldingRangeProvider = &enabled
	}

	if _, ok := handler.(SelectionRangeHandler); ok {
		caps.SelectionRangeProvider = &enabled
	}

	if _, ok := handler.(LinkedEditingRangeHandler); ok {
		caps.LinkedEditingRangeProvider = &enabled
	}

	if _, ok := handler.(CallHierarchyHandler); ok {
		caps.CallHierarchyProvider = &enabled
	}

	if _, ok := handler.(MonikerHandler); ok {
		caps.MonikerProvider = &enabled
	}

	if _, ok := handler.(TypeHierarchyHandler); ok {
		caps.TypeHierarchyProvider = &enabled
	}

	if _, ok := handler.(InlayHintHandler); ok {
		opts := &lsp.InlayHintOptions{}
		if _, ok := handler.(InlayHintResolveHandler); ok {
			opts.ResolveProvider = &enabled
		}
		caps.InlayHintProvider = opts
	}

	if _, ok := handler.(InlineValueHandler); ok {
		caps.InlineValueProvider = &enabled
	}

	if _, ok := handler.(DocumentDiagnosticHandler); ok {
		opts := &lsp.DiagnosticOptions{}
		if _, ok := handler.(WorkspaceDiagnosticHandler); ok {
			opts.WorkspaceDiagnostics = true
		}
		caps.DiagnosticProvider = opts
	}

	if _, ok := handler.(WorkspaceSymbolHandler); ok {
		caps.WorkspaceSymbolProvider = &enabled
	}

	if _, ok := handler.(ExecuteCommandHandler); ok {
		caps.ExecuteCommandProvider = &lsp.ExecuteCommandOptions{}
	}

	if hasSemanticTokensHandler(handler) {
		caps.SemanticTokensProvider = buildSemanticTokensOptions(handler, nil)
	}

	// Workspace file operations
	allFiles := fileOperationRegistrationOptions(nil)
	fileOps := &lsp.FileOperationOptions{}
	hasFileOps := false
	if _, ok := handler.(WillCreateFilesHandler); ok {
		fileOps.WillCreate = &allFiles
		hasFileOps = true
	}
	if _, ok := handler.(WillRenameFilesHandler); ok {
		fileOps.WillRename = &allFiles
		hasFileOps = true
	}
	if _, ok := handler.(WillDeleteFilesHandler); ok {
		fileOps.WillDelete = &allFiles
		hasFileOps = true
	}
	if hasFileOps {
		caps.Workspace = &lsp.ServerWorkspaceCapabilities{
			FileOperations: fileOps,
		}
	}

	return caps
}

func applyCapabilityOptions(caps *lsp.ServerCapabilities, handler any, opts CapabilityOptions) {
	if opts.PositionEncoding != nil {
		caps.PositionEncoding = opts.PositionEncoding
	}

	if caps.CompletionProvider != nil && opts.Completion != nil {
		completion := *opts.Completion
		if _, ok := handler.(CompletionResolveHandler); ok {
			completion.ResolveProvider = &enabled
		}
		caps.CompletionProvider = &completion
	}

	if caps.SignatureHelpProvider != nil && opts.SignatureHelp != nil {
		signatureHelp := *opts.SignatureHelp
		caps.SignatureHelpProvider = &signatureHelp
	}

	if caps.CodeActionProvider != nil && opts.CodeAction != nil {
		codeAction := *opts.CodeAction
		if _, ok := handler.(CodeActionResolveHandler); ok {
			codeAction.ResolveProvider = &enabled
		}
		caps.CodeActionProvider = &codeAction
	}

	if caps.ExecuteCommandProvider != nil && opts.ExecuteCommand != nil {
		executeCommand := *opts.ExecuteCommand
		caps.ExecuteCommandProvider = &executeCommand
	}

	if hasSemanticTokensHandler(handler) && opts.SemanticTokens != nil {
		caps.SemanticTokensProvider = buildSemanticTokensOptions(handler, opts.SemanticTokens)
	}

	if len(opts.FileOperationFilters) > 0 && caps.Workspace != nil && caps.Workspace.FileOperations != nil {
		reg := fileOperationRegistrationOptions(opts.FileOperationFilters)
		fileOps := caps.Workspace.FileOperations
		if fileOps.WillCreate != nil {
			fileOps.WillCreate = &reg
		}
		if fileOps.WillRename != nil {
			fileOps.WillRename = &reg
		}
		if fileOps.WillDelete != nil {
			fileOps.WillDelete = &reg
		}
	}
}

func hasSemanticTokensHandler(handler any) bool {
	if _, ok := handler.(SemanticTokensFullHandler); ok {
		return true
	}
	if _, ok := handler.(SemanticTokensDeltaHandler); ok {
		return true
	}
	if _, ok := handler.(SemanticTokensRangeHandler); ok {
		return true
	}
	return false
}

func buildSemanticTokensOptions(handler any, configured *lsp.SemanticTokensOptions) *lsp.SemanticTokensOptions {
	opts := &lsp.SemanticTokensOptions{}
	if configured != nil {
		copied := *configured
		opts = &copied
	}
	if _, ok := handler.(SemanticTokensFullHandler); ok {
		if opts.Full == nil {
			opts.Full = &lsp.SemanticTokensFull{}
		}
	}
	if _, ok := handler.(SemanticTokensDeltaHandler); ok {
		if opts.Full == nil {
			opts.Full = &lsp.SemanticTokensFull{}
		}
		opts.Full.Delta = &enabled
	}
	if _, ok := handler.(SemanticTokensRangeHandler); ok {
		opts.Range = &enabled
	}
	return opts
}

func fileOperationRegistrationOptions(filters []lsp.FileOperationFilter) lsp.FileOperationRegistrationOptions {
	if len(filters) == 0 {
		filters = []lsp.FileOperationFilter{{Pattern: lsp.FileOperationPattern{Glob: "**/*"}}}
	}
	return lsp.FileOperationRegistrationOptions{
		Filters: append([]lsp.FileOperationFilter(nil), filters...),
	}
}
