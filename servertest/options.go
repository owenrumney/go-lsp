package servertest

import (
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

type config struct {
	initParams *lsp.InitializeParams
	serverOpts []server.Option
}

// Option configures a Harness.
type Option func(*config)

// WithInitializeParams overrides the default initialize params.
func WithInitializeParams(params *lsp.InitializeParams) Option {
	return func(c *config) {
		c.initParams = params
	}
}

// WithServerOptions passes additional options to the server.
func WithServerOptions(opts ...server.Option) Option {
	return func(c *config) {
		c.serverOpts = append(c.serverOpts, opts...)
	}
}
