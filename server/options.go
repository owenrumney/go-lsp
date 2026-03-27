package server

import (
	"log/slog"
	"time"
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
