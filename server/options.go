package server

import "log/slog"

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
