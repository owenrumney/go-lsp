package server

import (
	"log/slog"
	"time"

	"github.com/owenrumney/go-lsp/lsp"
)

// Option configures a Server.
type Option func(*Server)

// WithDebugUI enables the debug web UI on the given address (e.g. ":7100").
func WithDebugUI(addr string) Option {
	return func(s *Server) {
		s.debugAddr = addr
	}
}

// WithLogger sets a logger for the server. The server logs lifecycle events,
// method dispatch, and errors. If not set, no logging is performed.
func WithLogger(logger *slog.Logger) Option {
	return func(s *Server) {
		s.logger = logger
	}
}

// WithRequestTimeout sets a default timeout for all incoming JSON-RPC requests.
// If a handler does not respond within the timeout, the request context is
// cancelled and the client receives a RequestCancelled error. A zero duration
// means no timeout (the default).
func WithRequestTimeout(d time.Duration) Option {
	return func(s *Server) {
		s.requestTimeout = d
	}
}

// CapabilityOptions configures detailed server capabilities that cannot be
// inferred from handler interfaces alone.
//
// These options enrich auto-detected capabilities. A feature is only advertised
// when the handler implements the corresponding handler interface.
type CapabilityOptions struct {
	Completion           *lsp.CompletionOptions
	SignatureHelp        *lsp.SignatureHelpOptions
	CodeAction           *lsp.CodeActionOptions
	ExecuteCommand       *lsp.ExecuteCommandOptions
	SemanticTokens       *lsp.SemanticTokensOptions
	FileOperationFilters []lsp.FileOperationFilter
	PositionEncoding     *lsp.PositionEncodingKind
}

// WithCapabilityOptions configures detailed LSP capability options.
func WithCapabilityOptions(opts CapabilityOptions) Option {
	return func(s *Server) {
		s.capabilityOptions = opts
	}
}

// WithCompletionOptions configures textDocument/completion capability options.
func WithCompletionOptions(opts lsp.CompletionOptions) Option {
	return func(s *Server) {
		s.capabilityOptions.Completion = &opts
	}
}

// WithSignatureHelpOptions configures textDocument/signatureHelp capability options.
func WithSignatureHelpOptions(opts lsp.SignatureHelpOptions) Option {
	return func(s *Server) {
		s.capabilityOptions.SignatureHelp = &opts
	}
}

// WithCodeActionOptions configures textDocument/codeAction capability options.
func WithCodeActionOptions(opts lsp.CodeActionOptions) Option {
	return func(s *Server) {
		s.capabilityOptions.CodeAction = &opts
	}
}

// WithExecuteCommandOptions configures workspace/executeCommand capability options.
func WithExecuteCommandOptions(opts lsp.ExecuteCommandOptions) Option {
	return func(s *Server) {
		s.capabilityOptions.ExecuteCommand = &opts
	}
}

// WithSemanticTokensOptions configures semantic token capability options,
// including the required token legend.
func WithSemanticTokensOptions(opts lsp.SemanticTokensOptions) Option {
	return func(s *Server) {
		s.capabilityOptions.SemanticTokens = &opts
	}
}

// WithFileOperationFilters configures the file filters used for workspace file
// operation capabilities such as willCreate, willRename, and willDelete.
func WithFileOperationFilters(filters []lsp.FileOperationFilter) Option {
	return func(s *Server) {
		s.capabilityOptions.FileOperationFilters = append([]lsp.FileOperationFilter(nil), filters...)
	}
}

// WithPositionEncoding advertises the position encoding used by this server.
// If unset, the LSP default of UTF-16 applies.
func WithPositionEncoding(encoding lsp.PositionEncodingKind) Option {
	return func(s *Server) {
		s.capabilityOptions.PositionEncoding = &encoding
	}
}
