package client

import "github.com/owenrumney/go-lsp/lsp"

// DefaultCapabilities declares what this package supports.
//
// dynamicRegistration is advertised only where Client actually honours a
// later client/registerCapability, since over-declaring makes a server
// withhold static capabilities and then register features the client ignores.
// Position encoding is UTF-16 only: it is the one encoding every server
// supports, and the document package implements no other.
func DefaultCapabilities() lsp.ClientCapabilities {
	return lsp.ClientCapabilities{
		General: &lsp.GeneralClientCapabilities{
			PositionEncodings: []lsp.PositionEncodingKind{lsp.PositionEncodingUTF16},
		},
		Workspace: &lsp.WorkspaceClientCapabilities{
			ApplyEdit:     boolPtr(true),
			Configuration: boolPtr(true),
			WorkspaceEdit: &lsp.WorkspaceEditClientCapabilities{
				DocumentChanges: boolPtr(true),
			},
			DidChangeConfiguration: dynamic(),
			DidChangeWatchedFiles:  dynamic(),
			ExecuteCommand:         dynamic(),
			WorkspaceFolders:       boolPtr(true),
		},
		TextDocument: &lsp.TextDocumentClientCapabilities{
			Synchronization: &lsp.TextDocumentSyncClientCapabilities{
				DynamicRegistration: boolPtr(true),
			},
			Completion: &lsp.CompletionClientCapabilities{
				DynamicRegistration: boolPtr(true),
			},
			Hover: &lsp.HoverClientCapabilities{
				DynamicRegistration: boolPtr(true),
				ContentFormat:       []lsp.MarkupKind{lsp.Markdown, lsp.PlainText},
			},
			Declaration:       &lsp.DeclarationClientCapabilities{DynamicRegistration: boolPtr(true)},
			Definition:        &lsp.DefinitionClientCapabilities{DynamicRegistration: boolPtr(true)},
			TypeDefinition:    &lsp.TypeDefinitionClientCapabilities{DynamicRegistration: boolPtr(true)},
			Implementation:    &lsp.ImplementationClientCapabilities{DynamicRegistration: boolPtr(true)},
			References:        dynamic(),
			DocumentHighlight: dynamic(),
			DocumentSymbol: &lsp.DocumentSymbolClientCapabilities{
				DynamicRegistration:               boolPtr(true),
				HierarchicalDocumentSymbolSupport: boolPtr(true),
			},
			CodeAction: &lsp.CodeActionClientCapabilities{
				DynamicRegistration: boolPtr(true),
				// Without this, servers answer with Command[] instead of CodeAction[].
				CodeActionLiteralSupport: &struct {
					CodeActionKind struct {
						ValueSet []lsp.CodeActionKind `json:"valueSet"`
					} `json:"codeActionKind"`
				}{
					CodeActionKind: struct {
						ValueSet []lsp.CodeActionKind `json:"valueSet"`
					}{ValueSet: []lsp.CodeActionKind{
						lsp.CodeActionQuickFix,
						lsp.CodeActionRefactor,
						lsp.CodeActionRefactorExtract,
						lsp.CodeActionRefactorInline,
						lsp.CodeActionRefactorRewrite,
						lsp.CodeActionSource,
						lsp.CodeActionSourceOrganizeImports,
					}},
				},
			},
			Rename: &lsp.RenameClientCapabilities{DynamicRegistration: boolPtr(true)},
			Diagnostic: &lsp.DiagnosticClientCapabilities{
				DynamicRegistration: boolPtr(true),
			},
			PublishDiagnostics: &lsp.PublishDiagnosticsClientCapabilities{
				RelatedInformation: boolPtr(true),
				VersionSupport:     boolPtr(true),
			},
		},
		Window: &lsp.WindowClientCapabilities{
			WorkDoneProgress: boolPtr(true),
		},
	}
}

func boolPtr(b bool) *bool { return &b }

func dynamic() *lsp.DynamicRegistrationCapability {
	return &lsp.DynamicRegistrationCapability{DynamicRegistration: boolPtr(true)}
}
