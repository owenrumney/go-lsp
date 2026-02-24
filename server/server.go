package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/owenrumney/go-lsp/internal/debugui"
	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

// Server is an LSP server that dispatches JSON-RPC messages to handler interfaces.
type Server struct {
	handler       any
	conn          *jsonrpc.Conn
	Client        *Client
	initialized   bool
	shutdown      bool
	customMethods map[string]jsonrpc.MethodHandler
	customNotifs  map[string]jsonrpc.NotificationHandler
	debugAddr     string
	debugUI       *debugui.DebugUI
}

// NewServer creates a new LSP server with the given handler.
// The handler must implement LifecycleHandler at minimum.
func NewServer(handler LifecycleHandler, opts ...Option) *Server {
	s := &Server{
		handler:       handler,
		customMethods: make(map[string]jsonrpc.MethodHandler),
		customNotifs:  make(map[string]jsonrpc.NotificationHandler),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// HandleMethod registers a custom JSON-RPC method handler.
// This must be called before Run.
func (s *Server) HandleMethod(method string, handler jsonrpc.MethodHandler) {
	s.customMethods[method] = handler
}

// HandleNotification registers a custom JSON-RPC notification handler.
// This must be called before Run.
func (s *Server) HandleNotification(method string, handler jsonrpc.NotificationHandler) {
	s.customNotifs[method] = handler
}

// Run starts the server, reading from and writing to rw.
func (s *Server) Run(ctx context.Context, rw io.ReadWriteCloser) error {
	if s.debugAddr != "" {
		store := debugui.NewStore()
		s.debugUI = debugui.New(s.debugAddr, store)
		if err := s.debugUI.ListenAndServe(ctx); err != nil {
			return fmt.Errorf("debugui: %w", err)
		}
		rw = debugui.NewTap(rw, store)
	}

	dispatcher := jsonrpc.NewDispatcher()
	s.conn = jsonrpc.NewConn(rw, dispatcher)
	s.Client = newClient(s.conn)

	if h, ok := s.handler.(ClientHandler); ok {
		h.SetClient(s.Client)
	}

	s.registerMethods(dispatcher)
	s.registerNotifications(dispatcher)

	for method, handler := range s.customMethods {
		dispatcher.RegisterMethod(method, handler)
	}
	for method, handler := range s.customNotifs {
		dispatcher.RegisterNotification(method, handler)
	}

	return s.conn.Serve(ctx)
}

func (s *Server) registerMethods(d *jsonrpc.Dispatcher) {
	d.RegisterMethod("initialize", s.handleInitialize)
	d.RegisterMethod("shutdown", s.handleShutdown)

	registerIf(d, s.handler, "textDocument/completion", handleCompletion)
	registerIf(d, s.handler, "completionItem/resolve", handleCompletionResolve)
	registerIf(d, s.handler, "textDocument/hover", handleHover)
	registerIf(d, s.handler, "textDocument/signatureHelp", handleSignatureHelp)
	registerIf(d, s.handler, "textDocument/declaration", handleDeclaration)
	registerIf(d, s.handler, "textDocument/definition", handleDefinition)
	registerIf(d, s.handler, "textDocument/typeDefinition", handleTypeDefinition)
	registerIf(d, s.handler, "textDocument/implementation", handleImplementation)
	registerIf(d, s.handler, "textDocument/references", handleReferences)
	registerIf(d, s.handler, "textDocument/documentHighlight", handleDocumentHighlight)
	registerIf(d, s.handler, "textDocument/documentSymbol", handleDocumentSymbol)
	registerIf(d, s.handler, "textDocument/codeAction", handleCodeAction)
	registerIf(d, s.handler, "codeAction/resolve", handleCodeActionResolve)
	registerIf(d, s.handler, "textDocument/codeLens", handleCodeLens)
	registerIf(d, s.handler, "codeLens/resolve", handleCodeLensResolve)
	registerIf(d, s.handler, "textDocument/documentLink", handleDocumentLink)
	registerIf(d, s.handler, "documentLink/resolve", handleDocumentLinkResolve)
	registerIf(d, s.handler, "textDocument/documentColor", handleDocumentColor)
	registerIf(d, s.handler, "textDocument/colorPresentation", handleColorPresentation)
	registerIf(d, s.handler, "textDocument/formatting", handleFormatting)
	registerIf(d, s.handler, "textDocument/rangeFormatting", handleRangeFormatting)
	registerIf(d, s.handler, "textDocument/onTypeFormatting", handleOnTypeFormatting)
	registerIf(d, s.handler, "textDocument/rename", handleRename)
	registerIf(d, s.handler, "textDocument/prepareRename", handlePrepareRename)
	registerIf(d, s.handler, "textDocument/foldingRange", handleFoldingRange)
	registerIf(d, s.handler, "textDocument/selectionRange", handleSelectionRange)
	registerIf(d, s.handler, "textDocument/linkedEditingRange", handleLinkedEditingRange)
	registerIf(d, s.handler, "textDocument/moniker", handleMoniker)
	registerIf(d, s.handler, "textDocument/willSaveWaitUntil", handleWillSaveWaitUntil)
	registerIf(d, s.handler, "workspace/symbol", handleWorkspaceSymbol)
	registerIf(d, s.handler, "workspace/executeCommand", handleExecuteCommand)
	registerIf(d, s.handler, "workspace/willCreateFiles", handleWillCreateFiles)
	registerIf(d, s.handler, "workspace/willRenameFiles", handleWillRenameFiles)
	registerIf(d, s.handler, "workspace/willDeleteFiles", handleWillDeleteFiles)

	if h, ok := s.handler.(CallHierarchyHandler); ok {
		d.RegisterMethod("textDocument/prepareCallHierarchy", typedHandler(h, CallHierarchyHandler.PrepareCallHierarchy))
		d.RegisterMethod("callHierarchy/incomingCalls", typedHandler(h, CallHierarchyHandler.IncomingCalls))
		d.RegisterMethod("callHierarchy/outgoingCalls", typedHandler(h, CallHierarchyHandler.OutgoingCalls))
	}

	if h, ok := s.handler.(TypeHierarchyHandler); ok {
		d.RegisterMethod("textDocument/prepareTypeHierarchy", typedHandler(h, TypeHierarchyHandler.PrepareTypeHierarchy))
		d.RegisterMethod("typeHierarchy/supertypes", typedHandler(h, TypeHierarchyHandler.Supertypes))
		d.RegisterMethod("typeHierarchy/subtypes", typedHandler(h, TypeHierarchyHandler.Subtypes))
	}

	registerIf(d, s.handler, "textDocument/inlayHint", handleInlayHint)
	registerIf(d, s.handler, "inlayHint/resolve", handleInlayHintResolve)
	registerIf(d, s.handler, "textDocument/inlineValue", handleInlineValue)
	registerIf(d, s.handler, "textDocument/diagnostic", handleDocumentDiagnostic)
	registerIf(d, s.handler, "workspace/diagnostic", handleWorkspaceDiagnostic)

	if h, ok := s.handler.(SemanticTokensFullHandler); ok {
		d.RegisterMethod("textDocument/semanticTokens/full", typedHandler(h, SemanticTokensFullHandler.SemanticTokensFull))
	}
	if h, ok := s.handler.(SemanticTokensDeltaHandler); ok {
		d.RegisterMethod("textDocument/semanticTokens/full/delta", typedHandler(h, SemanticTokensDeltaHandler.SemanticTokensDelta))
	}
	if h, ok := s.handler.(SemanticTokensRangeHandler); ok {
		d.RegisterMethod("textDocument/semanticTokens/range", typedHandler(h, SemanticTokensRangeHandler.SemanticTokensRange))
	}
}

func (s *Server) registerNotifications(d *jsonrpc.Dispatcher) {
	d.RegisterNotification("initialized", func(_ context.Context, _ json.RawMessage) error {
		return nil
	})

	d.RegisterNotification("exit", func(_ context.Context, _ json.RawMessage) error {
		return fmt.Errorf("exit")
	})

	if h, ok := s.handler.(TextDocumentSyncHandler); ok {
		d.RegisterNotification("textDocument/didOpen", notifHandler(h, TextDocumentSyncHandler.DidOpen))
		d.RegisterNotification("textDocument/didChange", notifHandler(h, TextDocumentSyncHandler.DidChange))
		d.RegisterNotification("textDocument/didClose", notifHandler(h, TextDocumentSyncHandler.DidClose))
	}

	if h, ok := s.handler.(TextDocumentSaveHandler); ok {
		d.RegisterNotification("textDocument/didSave", notifHandler(h, TextDocumentSaveHandler.DidSave))
	}

	if h, ok := s.handler.(TextDocumentWillSaveHandler); ok {
		d.RegisterNotification("textDocument/willSave", notifHandler(h, TextDocumentWillSaveHandler.WillSave))
	}

	if h, ok := s.handler.(WorkspaceFoldersHandler); ok {
		d.RegisterNotification("workspace/didChangeWorkspaceFolders", notifHandler(h, WorkspaceFoldersHandler.DidChangeWorkspaceFolders))
	}

	if h, ok := s.handler.(DidChangeConfigurationHandler); ok {
		d.RegisterNotification("workspace/didChangeConfiguration", notifHandler(h, DidChangeConfigurationHandler.DidChangeConfiguration))
	}

	if h, ok := s.handler.(DidChangeWatchedFilesHandler); ok {
		d.RegisterNotification("workspace/didChangeWatchedFiles", notifHandler(h, DidChangeWatchedFilesHandler.DidChangeWatchedFiles))
	}

	if h, ok := s.handler.(SetTraceHandler); ok {
		d.RegisterNotification("$/setTrace", notifHandler(h, SetTraceHandler.SetTrace))
	}
}

func (s *Server) handleInitialize(ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.InitializeParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}

	h := s.handler.(LifecycleHandler)
	result, err := h.Initialize(ctx, &p)
	if err != nil {
		return nil, err
	}

	// Merge auto-detected capabilities
	autoCaps := buildCapabilities(s.handler)
	mergeCapabilities(&result.Capabilities, &autoCaps)

	s.initialized = true
	return result, nil
}

func (s *Server) handleShutdown(ctx context.Context, _ json.RawMessage) (any, error) {
	h := s.handler.(LifecycleHandler)
	err := h.Shutdown(ctx)
	s.shutdown = true
	return nil, err
}

// mergeCapabilities fills in any auto-detected capabilities that weren't explicitly set.
func mergeCapabilities(dst, src *lsp.ServerCapabilities) {
	if dst.TextDocumentSync == nil {
		dst.TextDocumentSync = src.TextDocumentSync
	}
	if dst.CompletionProvider == nil {
		dst.CompletionProvider = src.CompletionProvider
	}
	if dst.HoverProvider == nil {
		dst.HoverProvider = src.HoverProvider
	}
	if dst.SignatureHelpProvider == nil {
		dst.SignatureHelpProvider = src.SignatureHelpProvider
	}
	if dst.DeclarationProvider == nil {
		dst.DeclarationProvider = src.DeclarationProvider
	}
	if dst.DefinitionProvider == nil {
		dst.DefinitionProvider = src.DefinitionProvider
	}
	if dst.TypeDefinitionProvider == nil {
		dst.TypeDefinitionProvider = src.TypeDefinitionProvider
	}
	if dst.ImplementationProvider == nil {
		dst.ImplementationProvider = src.ImplementationProvider
	}
	if dst.ReferencesProvider == nil {
		dst.ReferencesProvider = src.ReferencesProvider
	}
	if dst.DocumentHighlightProvider == nil {
		dst.DocumentHighlightProvider = src.DocumentHighlightProvider
	}
	if dst.DocumentSymbolProvider == nil {
		dst.DocumentSymbolProvider = src.DocumentSymbolProvider
	}
	if dst.CodeActionProvider == nil {
		dst.CodeActionProvider = src.CodeActionProvider
	}
	if dst.CodeLensProvider == nil {
		dst.CodeLensProvider = src.CodeLensProvider
	}
	if dst.DocumentLinkProvider == nil {
		dst.DocumentLinkProvider = src.DocumentLinkProvider
	}
	if dst.ColorProvider == nil {
		dst.ColorProvider = src.ColorProvider
	}
	if dst.DocumentFormattingProvider == nil {
		dst.DocumentFormattingProvider = src.DocumentFormattingProvider
	}
	if dst.DocumentRangeFormattingProvider == nil {
		dst.DocumentRangeFormattingProvider = src.DocumentRangeFormattingProvider
	}
	if dst.DocumentOnTypeFormattingProvider == nil {
		dst.DocumentOnTypeFormattingProvider = src.DocumentOnTypeFormattingProvider
	}
	if dst.RenameProvider == nil {
		dst.RenameProvider = src.RenameProvider
	}
	if dst.FoldingRangeProvider == nil {
		dst.FoldingRangeProvider = src.FoldingRangeProvider
	}
	if dst.SelectionRangeProvider == nil {
		dst.SelectionRangeProvider = src.SelectionRangeProvider
	}
	if dst.LinkedEditingRangeProvider == nil {
		dst.LinkedEditingRangeProvider = src.LinkedEditingRangeProvider
	}
	if dst.CallHierarchyProvider == nil {
		dst.CallHierarchyProvider = src.CallHierarchyProvider
	}
	if dst.SemanticTokensProvider == nil {
		dst.SemanticTokensProvider = src.SemanticTokensProvider
	}
	if dst.MonikerProvider == nil {
		dst.MonikerProvider = src.MonikerProvider
	}
	if dst.TypeHierarchyProvider == nil {
		dst.TypeHierarchyProvider = src.TypeHierarchyProvider
	}
	if dst.InlayHintProvider == nil {
		dst.InlayHintProvider = src.InlayHintProvider
	}
	if dst.InlineValueProvider == nil {
		dst.InlineValueProvider = src.InlineValueProvider
	}
	if dst.DiagnosticProvider == nil {
		dst.DiagnosticProvider = src.DiagnosticProvider
	}
	if dst.WorkspaceSymbolProvider == nil {
		dst.WorkspaceSymbolProvider = src.WorkspaceSymbolProvider
	}
	if dst.Workspace == nil {
		dst.Workspace = src.Workspace
	}
}

// registerIf registers a method handler only if the server handler implements the interface.
func registerIf[H any](d *jsonrpc.Dispatcher, handler any, method string, fn func(H, context.Context, json.RawMessage) (any, error)) {
	if h, ok := handler.(H); ok {
		d.RegisterMethod(method, func(ctx context.Context, params json.RawMessage) (any, error) {
			return fn(h, ctx, params)
		})
	}
}

// typedHandler creates a MethodHandler from a typed handler method.
func typedHandler[H any, P any, R any](h H, fn func(H, context.Context, *P) (R, error)) jsonrpc.MethodHandler {
	return func(ctx context.Context, params json.RawMessage) (any, error) {
		var p P
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
		}
		return fn(h, ctx, &p)
	}
}

// notifHandler creates a NotificationHandler from a typed handler method.
func notifHandler[H any, P any](h H, fn func(H, context.Context, *P) error) jsonrpc.NotificationHandler {
	return func(ctx context.Context, params json.RawMessage) error {
		var p P
		if err := json.Unmarshal(params, &p); err != nil {
			return err
		}
		return fn(h, ctx, &p)
	}
}

// Typed handler wrappers for registerIf.
func handleCompletion(h CompletionHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.CompletionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Completion(ctx, &p)
}

func handleCompletionResolve(h CompletionResolveHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.CompletionItem
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveCompletionItem(ctx, &p)
}

func handleHover(h HoverHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.HoverParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Hover(ctx, &p)
}

func handleSignatureHelp(h SignatureHelpHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.SignatureHelpParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.SignatureHelp(ctx, &p)
}

func handleDeclaration(h DeclarationHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DeclarationParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Declaration(ctx, &p)
}

func handleDefinition(h DefinitionHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DefinitionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Definition(ctx, &p)
}

func handleTypeDefinition(h TypeDefinitionHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.TypeDefinitionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.TypeDefinition(ctx, &p)
}

func handleImplementation(h ImplementationHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.ImplementationParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Implementation(ctx, &p)
}

func handleReferences(h ReferencesHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.ReferenceParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.References(ctx, &p)
}

func handleDocumentHighlight(h DocumentHighlightHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DocumentHighlightParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentHighlight(ctx, &p)
}

func handleDocumentSymbol(h DocumentSymbolHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DocumentSymbolParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentSymbol(ctx, &p)
}

func handleCodeAction(h CodeActionHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.CodeActionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.CodeAction(ctx, &p)
}

func handleCodeLens(h CodeLensHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.CodeLensParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.CodeLens(ctx, &p)
}

func handleCodeLensResolve(h CodeLensResolveHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.CodeLens
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveCodeLens(ctx, &p)
}

func handleDocumentLink(h DocumentLinkHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DocumentLinkParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentLink(ctx, &p)
}

func handleDocumentColor(h DocumentColorHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DocumentColorParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentColor(ctx, &p)
}

func handleColorPresentation(h ColorPresentationHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.ColorPresentationParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ColorPresentation(ctx, &p)
}

func handleFormatting(h DocumentFormattingHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DocumentFormattingParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Formatting(ctx, &p)
}

func handleRangeFormatting(h DocumentRangeFormattingHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DocumentRangeFormattingParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.RangeFormatting(ctx, &p)
}

func handleOnTypeFormatting(h DocumentOnTypeFormattingHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DocumentOnTypeFormattingParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.OnTypeFormatting(ctx, &p)
}

func handleRename(h RenameHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.RenameParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Rename(ctx, &p)
}

func handlePrepareRename(h PrepareRenameHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.PrepareRenameParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.PrepareRename(ctx, &p)
}

func handleFoldingRange(h FoldingRangeHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.FoldingRangeParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.FoldingRange(ctx, &p)
}

func handleSelectionRange(h SelectionRangeHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.SelectionRangeParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.SelectionRange(ctx, &p)
}

func handleLinkedEditingRange(h LinkedEditingRangeHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.LinkedEditingRangeParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.LinkedEditingRange(ctx, &p)
}

func handleMoniker(h MonikerHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.MonikerParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Moniker(ctx, &p)
}

func handleWorkspaceSymbol(h WorkspaceSymbolHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.WorkspaceSymbolParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WorkspaceSymbol(ctx, &p)
}

func handleExecuteCommand(h ExecuteCommandHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.ExecuteCommandParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ExecuteCommand(ctx, &p)
}

func handleWillSaveWaitUntil(h TextDocumentWillSaveWaitUntilHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.WillSaveTextDocumentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WillSaveWaitUntil(ctx, &p)
}

func handleCodeActionResolve(h CodeActionResolveHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.CodeAction
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveCodeAction(ctx, &p)
}

func handleDocumentLinkResolve(h DocumentLinkResolveHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DocumentLink
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveDocumentLink(ctx, &p)
}

func handleWillCreateFiles(h WillCreateFilesHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.CreateFilesParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WillCreateFiles(ctx, &p)
}

func handleWillRenameFiles(h WillRenameFilesHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.RenameFilesParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WillRenameFiles(ctx, &p)
}

func handleWillDeleteFiles(h WillDeleteFilesHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DeleteFilesParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WillDeleteFiles(ctx, &p)
}

func handleInlayHint(h InlayHintHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.InlayHintParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.InlayHint(ctx, &p)
}

func handleInlayHintResolve(h InlayHintResolveHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.InlayHint
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveInlayHint(ctx, &p)
}

func handleInlineValue(h InlineValueHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.InlineValueParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.InlineValue(ctx, &p)
}

func handleDocumentDiagnostic(h DocumentDiagnosticHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.DocumentDiagnosticParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentDiagnostic(ctx, &p)
}

func handleWorkspaceDiagnostic(h WorkspaceDiagnosticHandler, ctx context.Context, params json.RawMessage) (any, error) {
	var p lsp.WorkspaceDiagnosticParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WorkspaceDiagnostic(ctx, &p)
}
