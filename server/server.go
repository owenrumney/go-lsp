package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/owenrumney/go-lsp/internal/debugui"
	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

// Server is an LSP server that dispatches JSON-RPC messages to handler interfaces.
type Server struct {
	handler             any
	conn                *jsonrpc.Conn
	Client              *Client
	initialized         bool
	shutdown            bool
	customMethods       map[string]jsonrpc.MethodHandler
	customNotifications map[string]jsonrpc.NotificationHandler
	debugAddr           string
	debugUI             *debugui.DebugUI
	logger              *slog.Logger
	requestTimeout      time.Duration
}

// NewServer creates a new LSP server with the given handler.
// The handler must implement LifecycleHandler at minimum.
func NewServer(handler LifecycleHandler, opts ...Option) *Server {
	s := &Server{
		handler:             handler,
		customMethods:       make(map[string]jsonrpc.MethodHandler),
		customNotifications: make(map[string]jsonrpc.NotificationHandler),
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// DebugHandler returns a slog.Handler that sends log records to the debug UI's log tab.
// Returns nil if the debug UI is not enabled. Must be called after Run has started.
//
// Usage: logger := slog.New(srv.DebugHandler())
func (s *Server) DebugHandler() slog.Handler {
	if s.debugUI == nil {
		return nil
	}
	return s.debugUI.SlogHandler()
}

// HandleMethod registers a custom JSON-RPC method handler.
// This must be called before Run.
func (s *Server) HandleMethod(method string, handler jsonrpc.MethodHandler) {
	s.customMethods[method] = handler
}

// HandleNotification registers a custom JSON-RPC notification handler.
// This must be called before Run.
func (s *Server) HandleNotification(method string, handler jsonrpc.NotificationHandler) {
	s.customNotifications[method] = handler
}

// Run starts the server, reading from and writing to rw.
func (s *Server) Run(ctx context.Context, rw io.ReadWriteCloser) error {
	if s.debugAddr != "" {
		logStore := debugui.NewLogStore()
		store := debugui.NewStore(logStore)
		s.debugUI = debugui.New(s.debugAddr, store, logStore)
		if err := s.debugUI.ListenAndServe(ctx); err != nil {
			return fmt.Errorf("debugui: %w", err)
		}
		rw = debugui.NewTap(rw, store)
	}

	dispatcher := jsonrpc.NewDispatcher()
	s.conn = jsonrpc.NewConn(rw, dispatcher)
	if s.requestTimeout > 0 {
		s.conn.SetRequestTimeout(s.requestTimeout)
	}
	s.Client = newClient(s.conn)

	if h, ok := s.handler.(ClientHandler); ok {
		h.SetClient(s.Client)
	}

	s.registerMethods(dispatcher)
	s.registerNotifications(dispatcher)

	for method, handler := range s.customMethods {
		dispatcher.RegisterMethod(method, s.logMethod(method, handler))
	}
	for method, handler := range s.customNotifications {
		dispatcher.RegisterNotification(method, s.logNotification(method, handler))
	}

	if s.logger != nil {
		s.logger.Info("server starting")
	}

	return s.conn.Serve(ctx)
}

func (s *Server) registerMethods(d *jsonrpc.Dispatcher) {
	d.RegisterMethod("initialize", s.logMethod("initialize", s.handleInitialize))
	d.RegisterMethod("shutdown", s.logMethod("shutdown", s.handleShutdown))

	registerIf(d, s, "textDocument/completion", handleCompletion)
	registerIf(d, s, "completionItem/resolve", handleCompletionResolve)
	registerIf(d, s, "textDocument/hover", handleHover)
	registerIf(d, s, "textDocument/signatureHelp", handleSignatureHelp)
	registerIf(d, s, "textDocument/declaration", handleDeclaration)
	registerIf(d, s, "textDocument/definition", handleDefinition)
	registerIf(d, s, "textDocument/typeDefinition", handleTypeDefinition)
	registerIf(d, s, "textDocument/implementation", handleImplementation)
	registerIf(d, s, "textDocument/references", handleReferences)
	registerIf(d, s, "textDocument/documentHighlight", handleDocumentHighlight)
	registerIf(d, s, "textDocument/documentSymbol", handleDocumentSymbol)
	registerIf(d, s, "textDocument/codeAction", handleCodeAction)
	registerIf(d, s, "codeAction/resolve", handleCodeActionResolve)
	registerIf(d, s, "textDocument/codeLens", handleCodeLens)
	registerIf(d, s, "codeLens/resolve", handleCodeLensResolve)
	registerIf(d, s, "textDocument/documentLink", handleDocumentLink)
	registerIf(d, s, "documentLink/resolve", handleDocumentLinkResolve)
	registerIf(d, s, "textDocument/documentColor", handleDocumentColor)
	registerIf(d, s, "textDocument/colorPresentation", handleColorPresentation)
	registerIf(d, s, "textDocument/formatting", handleFormatting)
	registerIf(d, s, "textDocument/rangeFormatting", handleRangeFormatting)
	registerIf(d, s, "textDocument/onTypeFormatting", handleOnTypeFormatting)
	registerIf(d, s, "textDocument/rename", handleRename)
	registerIf(d, s, "textDocument/prepareRename", handlePrepareRename)
	registerIf(d, s, "textDocument/foldingRange", handleFoldingRange)
	registerIf(d, s, "textDocument/selectionRange", handleSelectionRange)
	registerIf(d, s, "textDocument/linkedEditingRange", handleLinkedEditingRange)
	registerIf(d, s, "textDocument/moniker", handleMoniker)
	registerIf(d, s, "textDocument/willSaveWaitUntil", handleWillSaveWaitUntil)
	registerIf(d, s, "workspace/symbol", handleWorkspaceSymbol)
	registerIf(d, s, "workspace/executeCommand", handleExecuteCommand)
	registerIf(d, s, "workspace/willCreateFiles", handleWillCreateFiles)
	registerIf(d, s, "workspace/willRenameFiles", handleWillRenameFiles)
	registerIf(d, s, "workspace/willDeleteFiles", handleWillDeleteFiles)

	if h, ok := s.handler.(CallHierarchyHandler); ok {
		d.RegisterMethod("textDocument/prepareCallHierarchy", s.logMethod("textDocument/prepareCallHierarchy", typedHandler(h, CallHierarchyHandler.PrepareCallHierarchy)))
		d.RegisterMethod("callHierarchy/incomingCalls", s.logMethod("callHierarchy/incomingCalls", typedHandler(h, CallHierarchyHandler.IncomingCalls)))
		d.RegisterMethod("callHierarchy/outgoingCalls", s.logMethod("callHierarchy/outgoingCalls", typedHandler(h, CallHierarchyHandler.OutgoingCalls)))
	}

	if h, ok := s.handler.(TypeHierarchyHandler); ok {
		d.RegisterMethod("textDocument/prepareTypeHierarchy", s.logMethod("textDocument/prepareTypeHierarchy", typedHandler(h, TypeHierarchyHandler.PrepareTypeHierarchy)))
		d.RegisterMethod("typeHierarchy/supertypes", s.logMethod("typeHierarchy/supertypes", typedHandler(h, TypeHierarchyHandler.Supertypes)))
		d.RegisterMethod("typeHierarchy/subtypes", s.logMethod("typeHierarchy/subtypes", typedHandler(h, TypeHierarchyHandler.Subtypes)))
	}

	registerIf(d, s, "textDocument/inlayHint", handleInlayHint)
	registerIf(d, s, "inlayHint/resolve", handleInlayHintResolve)
	registerIf(d, s, "textDocument/inlineValue", handleInlineValue)
	registerIf(d, s, "textDocument/diagnostic", handleDocumentDiagnostic)
	registerIf(d, s, "workspace/diagnostic", handleWorkspaceDiagnostic)

	if h, ok := s.handler.(SemanticTokensFullHandler); ok {
		d.RegisterMethod("textDocument/semanticTokens/full", s.logMethod("textDocument/semanticTokens/full", typedHandler(h, SemanticTokensFullHandler.SemanticTokensFull)))
	}
	if h, ok := s.handler.(SemanticTokensDeltaHandler); ok {
		d.RegisterMethod("textDocument/semanticTokens/full/delta", s.logMethod("textDocument/semanticTokens/full/delta", typedHandler(h, SemanticTokensDeltaHandler.SemanticTokensDelta)))
	}
	if h, ok := s.handler.(SemanticTokensRangeHandler); ok {
		d.RegisterMethod("textDocument/semanticTokens/range", s.logMethod("textDocument/semanticTokens/range", typedHandler(h, SemanticTokensRangeHandler.SemanticTokensRange)))
	}
}

func (s *Server) registerNotifications(d *jsonrpc.Dispatcher) {
	d.RegisterNotification("initialized", s.logNotification("initialized", func(_ context.Context, _ json.RawMessage) error {
		return nil
	}))

	d.RegisterNotification("exit", s.logNotification("exit", func(_ context.Context, _ json.RawMessage) error {
		return fmt.Errorf("exit")
	}))

	if h, ok := s.handler.(TextDocumentSyncHandler); ok {
		d.RegisterNotification("textDocument/didOpen", s.logNotification("textDocument/didOpen", notifHandler(h, TextDocumentSyncHandler.DidOpen)))
		d.RegisterNotification("textDocument/didChange", s.logNotification("textDocument/didChange", notifHandler(h, TextDocumentSyncHandler.DidChange)))
		d.RegisterNotification("textDocument/didClose", s.logNotification("textDocument/didClose", notifHandler(h, TextDocumentSyncHandler.DidClose)))
	}

	if h, ok := s.handler.(TextDocumentSaveHandler); ok {
		d.RegisterNotification("textDocument/didSave", s.logNotification("textDocument/didSave", notifHandler(h, TextDocumentSaveHandler.DidSave)))
	}

	if h, ok := s.handler.(TextDocumentWillSaveHandler); ok {
		d.RegisterNotification("textDocument/willSave", s.logNotification("textDocument/willSave", notifHandler(h, TextDocumentWillSaveHandler.WillSave)))
	}

	if h, ok := s.handler.(WorkspaceFoldersHandler); ok {
		d.RegisterNotification("workspace/didChangeWorkspaceFolders", s.logNotification("workspace/didChangeWorkspaceFolders", notifHandler(h, WorkspaceFoldersHandler.DidChangeWorkspaceFolders)))
	}

	if h, ok := s.handler.(DidChangeConfigurationHandler); ok {
		d.RegisterNotification("workspace/didChangeConfiguration", s.logNotification("workspace/didChangeConfiguration", notifHandler(h, DidChangeConfigurationHandler.DidChangeConfiguration)))
	}

	if h, ok := s.handler.(DidChangeWatchedFilesHandler); ok {
		d.RegisterNotification("workspace/didChangeWatchedFiles", s.logNotification("workspace/didChangeWatchedFiles", notifHandler(h, DidChangeWatchedFilesHandler.DidChangeWatchedFiles)))
	}

	if h, ok := s.handler.(SetTraceHandler); ok {
		d.RegisterNotification("$/setTrace", s.logNotification("$/setTrace", notifHandler(h, SetTraceHandler.SetTrace)))
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

	if s.debugUI != nil {
		s.debugUI.SetCapabilities(result.Capabilities)
	}

	s.initialized = true

	if s.logger != nil {
		s.logger.Info("server initialized", "serverName", result.ServerInfo.Name)
	}

	return result, nil
}

func (s *Server) handleShutdown(ctx context.Context, _ json.RawMessage) (any, error) {
	h := s.handler.(LifecycleHandler)
	err := h.Shutdown(ctx)
	s.shutdown = true

	if s.logger != nil {
		s.logger.Info("server shutdown")
	}

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
func registerIf[H any](d *jsonrpc.Dispatcher, s *Server, method string, fn func(context.Context, H, json.RawMessage) (any, error)) {
	if h, ok := s.handler.(H); ok {
		handler := jsonrpc.MethodHandler(func(ctx context.Context, params json.RawMessage) (any, error) {
			return fn(ctx, h, params)
		})
		d.RegisterMethod(method, s.logMethod(method, handler))
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

// logMethod wraps a MethodHandler with logging for dispatch, duration, and errors.
func (s *Server) logMethod(method string, handler jsonrpc.MethodHandler) jsonrpc.MethodHandler {
	if s.logger == nil {
		return handler
	}
	return func(ctx context.Context, params json.RawMessage) (any, error) {
		start := time.Now()
		result, err := handler(ctx, params)
		duration := time.Since(start)
		if err != nil {
			s.logger.Error("method error", "method", method, "duration", duration, "error", err)
		} else {
			s.logger.Debug("method handled", "method", method, "duration", duration)
		}
		return result, err
	}
}

// logNotification wraps a NotificationHandler with logging.
func (s *Server) logNotification(method string, handler jsonrpc.NotificationHandler) jsonrpc.NotificationHandler {
	if s.logger == nil {
		return handler
	}
	return func(ctx context.Context, params json.RawMessage) error {
		start := time.Now()
		err := handler(ctx, params)
		duration := time.Since(start)
		if err != nil {
			s.logger.Error("notification error", "method", method, "duration", duration, "error", err)
		} else {
			s.logger.Debug("notification handled", "method", method, "duration", duration)
		}
		return err
	}
}

// Typed handler wrappers for registerIf.
func handleCompletion(ctx context.Context, h CompletionHandler, params json.RawMessage) (any, error) {
	var p lsp.CompletionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Completion(ctx, &p)
}

func handleCompletionResolve(ctx context.Context, h CompletionResolveHandler, params json.RawMessage) (any, error) {
	var p lsp.CompletionItem
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveCompletionItem(ctx, &p)
}

func handleHover(ctx context.Context, h HoverHandler, params json.RawMessage) (any, error) {
	var p lsp.HoverParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Hover(ctx, &p)
}

func handleSignatureHelp(ctx context.Context, h SignatureHelpHandler, params json.RawMessage) (any, error) {
	var p lsp.SignatureHelpParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.SignatureHelp(ctx, &p)
}

func handleDeclaration(ctx context.Context, h DeclarationHandler, params json.RawMessage) (any, error) {
	var p lsp.DeclarationParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Declaration(ctx, &p)
}

func handleDefinition(ctx context.Context, h DefinitionHandler, params json.RawMessage) (any, error) {
	var p lsp.DefinitionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Definition(ctx, &p)
}

func handleTypeDefinition(ctx context.Context, h TypeDefinitionHandler, params json.RawMessage) (any, error) {
	var p lsp.TypeDefinitionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.TypeDefinition(ctx, &p)
}

func handleImplementation(ctx context.Context, h ImplementationHandler, params json.RawMessage) (any, error) {
	var p lsp.ImplementationParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Implementation(ctx, &p)
}

func handleReferences(ctx context.Context, h ReferencesHandler, params json.RawMessage) (any, error) {
	var p lsp.ReferenceParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.References(ctx, &p)
}

func handleDocumentHighlight(ctx context.Context, h DocumentHighlightHandler, params json.RawMessage) (any, error) {
	var p lsp.DocumentHighlightParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentHighlight(ctx, &p)
}

func handleDocumentSymbol(ctx context.Context, h DocumentSymbolHandler, params json.RawMessage) (any, error) {
	var p lsp.DocumentSymbolParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentSymbol(ctx, &p)
}

func handleCodeAction(ctx context.Context, h CodeActionHandler, params json.RawMessage) (any, error) {
	var p lsp.CodeActionParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.CodeAction(ctx, &p)
}

func handleCodeLens(ctx context.Context, h CodeLensHandler, params json.RawMessage) (any, error) {
	var p lsp.CodeLensParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.CodeLens(ctx, &p)
}

func handleCodeLensResolve(ctx context.Context, h CodeLensResolveHandler, params json.RawMessage) (any, error) {
	var p lsp.CodeLens
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveCodeLens(ctx, &p)
}

func handleDocumentLink(ctx context.Context, h DocumentLinkHandler, params json.RawMessage) (any, error) {
	var p lsp.DocumentLinkParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentLink(ctx, &p)
}

func handleDocumentColor(ctx context.Context, h DocumentColorHandler, params json.RawMessage) (any, error) {
	var p lsp.DocumentColorParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentColor(ctx, &p)
}

func handleColorPresentation(ctx context.Context, h ColorPresentationHandler, params json.RawMessage) (any, error) {
	var p lsp.ColorPresentationParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ColorPresentation(ctx, &p)
}

func handleFormatting(ctx context.Context, h DocumentFormattingHandler, params json.RawMessage) (any, error) {
	var p lsp.DocumentFormattingParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Formatting(ctx, &p)
}

func handleRangeFormatting(ctx context.Context, h DocumentRangeFormattingHandler, params json.RawMessage) (any, error) {
	var p lsp.DocumentRangeFormattingParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.RangeFormatting(ctx, &p)
}

func handleOnTypeFormatting(ctx context.Context, h DocumentOnTypeFormattingHandler, params json.RawMessage) (any, error) {
	var p lsp.DocumentOnTypeFormattingParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.OnTypeFormatting(ctx, &p)
}

func handleRename(ctx context.Context, h RenameHandler, params json.RawMessage) (any, error) {
	var p lsp.RenameParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Rename(ctx, &p)
}

func handlePrepareRename(ctx context.Context, h PrepareRenameHandler, params json.RawMessage) (any, error) {
	var p lsp.PrepareRenameParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.PrepareRename(ctx, &p)
}

func handleFoldingRange(ctx context.Context, h FoldingRangeHandler, params json.RawMessage) (any, error) {
	var p lsp.FoldingRangeParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.FoldingRange(ctx, &p)
}

func handleSelectionRange(ctx context.Context, h SelectionRangeHandler, params json.RawMessage) (any, error) {
	var p lsp.SelectionRangeParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.SelectionRange(ctx, &p)
}

func handleLinkedEditingRange(ctx context.Context, h LinkedEditingRangeHandler, params json.RawMessage) (any, error) {
	var p lsp.LinkedEditingRangeParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.LinkedEditingRange(ctx, &p)
}

func handleMoniker(ctx context.Context, h MonikerHandler, params json.RawMessage) (any, error) {
	var p lsp.MonikerParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.Moniker(ctx, &p)
}

func handleWorkspaceSymbol(ctx context.Context, h WorkspaceSymbolHandler, params json.RawMessage) (any, error) {
	var p lsp.WorkspaceSymbolParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WorkspaceSymbol(ctx, &p)
}

func handleExecuteCommand(ctx context.Context, h ExecuteCommandHandler, params json.RawMessage) (any, error) {
	var p lsp.ExecuteCommandParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ExecuteCommand(ctx, &p)
}

func handleWillSaveWaitUntil(ctx context.Context, h TextDocumentWillSaveWaitUntilHandler, params json.RawMessage) (any, error) {
	var p lsp.WillSaveTextDocumentParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WillSaveWaitUntil(ctx, &p)
}

func handleCodeActionResolve(ctx context.Context, h CodeActionResolveHandler, params json.RawMessage) (any, error) {
	var p lsp.CodeAction
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveCodeAction(ctx, &p)
}

func handleDocumentLinkResolve(ctx context.Context, h DocumentLinkResolveHandler, params json.RawMessage) (any, error) {
	var p lsp.DocumentLink
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveDocumentLink(ctx, &p)
}

func handleWillCreateFiles(ctx context.Context, h WillCreateFilesHandler, params json.RawMessage) (any, error) {
	var p lsp.CreateFilesParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WillCreateFiles(ctx, &p)
}

func handleWillRenameFiles(ctx context.Context, h WillRenameFilesHandler, params json.RawMessage) (any, error) {
	var p lsp.RenameFilesParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WillRenameFiles(ctx, &p)
}

func handleWillDeleteFiles(ctx context.Context, h WillDeleteFilesHandler, params json.RawMessage) (any, error) {
	var p lsp.DeleteFilesParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WillDeleteFiles(ctx, &p)
}

func handleInlayHint(ctx context.Context, h InlayHintHandler, params json.RawMessage) (any, error) {
	var p lsp.InlayHintParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.InlayHint(ctx, &p)
}

func handleInlayHintResolve(ctx context.Context, h InlayHintResolveHandler, params json.RawMessage) (any, error) {
	var p lsp.InlayHint
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.ResolveInlayHint(ctx, &p)
}

func handleInlineValue(ctx context.Context, h InlineValueHandler, params json.RawMessage) (any, error) {
	var p lsp.InlineValueParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.InlineValue(ctx, &p)
}

func handleDocumentDiagnostic(ctx context.Context, h DocumentDiagnosticHandler, params json.RawMessage) (any, error) {
	var p lsp.DocumentDiagnosticParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.DocumentDiagnostic(ctx, &p)
}

func handleWorkspaceDiagnostic(ctx context.Context, h WorkspaceDiagnosticHandler, params json.RawMessage) (any, error) {
	var p lsp.WorkspaceDiagnosticParams
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
	}
	return h.WorkspaceDiagnostic(ctx, &p)
}
