package server

import (
	"errors"
	"os"

	"github.com/owenrumney/go-lsp/internal/debugui"
)

// ErrDebugTraceUnavailable is returned when trace export is requested without
// either WithDebugCapture or WithDebugUI being enabled.
var ErrDebugTraceUnavailable = errors.New("debug trace unavailable")

// TraceExportOptions controls how a debug trace is exported.
type TraceExportOptions struct {
	// RedactDocumentText replaces LSP text payload fields such as text and
	// newText with a placeholder.
	RedactDocumentText bool

	// RedactFilePaths replaces file:// URIs and absolute-looking paths with a
	// placeholder.
	RedactFilePaths bool

	// RedactLogs omits captured log entries from the trace.
	RedactLogs bool

	// Pretty formats the exported JSON with indentation.
	Pretty bool
}

// ExportDebugTrace returns a JSON snapshot of the captured debug session.
//
// Capture must be enabled with WithDebugCapture or WithDebugUI, and Run must
// have started.
func (s *Server) ExportDebugTrace(opts TraceExportOptions) ([]byte, error) {
	if s.recorder == nil {
		return nil, ErrDebugTraceUnavailable
	}
	return s.recorder.ExportTrace(debugui.TraceExportOptions{
		RedactDocumentText: opts.RedactDocumentText,
		RedactFilePaths:    opts.RedactFilePaths,
		RedactLogs:         opts.RedactLogs,
		Pretty:             opts.Pretty,
	})
}

// SaveDebugTrace writes a JSON snapshot of the active debug UI session to path.
//
// The file is written with 0600 permissions because traces may contain source
// code, file paths, or other local project details.
func (s *Server) SaveDebugTrace(path string, opts TraceExportOptions) error {
	data, err := s.ExportDebugTrace(opts)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
