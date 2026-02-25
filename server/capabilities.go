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

	// Workspace file operations
	var fileOps *lsp.FileOperationOptions
	allFiles := lsp.FileOperationRegistrationOptions{
		Filters: []lsp.FileOperationFilter{{Pattern: lsp.FileOperationPattern{Glob: "**/*"}}},
	}
	if _, ok := handler.(WillCreateFilesHandler); ok {
		if fileOps == nil {
			fileOps = &lsp.FileOperationOptions{}
		}
		fileOps.WillCreate = &allFiles
	}
	if _, ok := handler.(WillRenameFilesHandler); ok {
		if fileOps == nil {
			fileOps = &lsp.FileOperationOptions{}
		}
		fileOps.WillRename = &allFiles
	}
	if _, ok := handler.(WillDeleteFilesHandler); ok {
		if fileOps == nil {
			fileOps = &lsp.FileOperationOptions{}
		}
		fileOps.WillDelete = &allFiles
	}
	if fileOps != nil {
		caps.Workspace = &lsp.ServerWorkspaceCapabilities{
			FileOperations: fileOps,
		}
	}

	return caps
}
