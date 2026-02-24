package server

// Option configures a Server.
type Option func(*Server)

// WithDebugUI enables the debug web UI on the given address (e.g. ":7100").
func WithDebugUI(addr string) Option {
	return func(s *Server) {
		s.debugAddr = addr
	}
}
