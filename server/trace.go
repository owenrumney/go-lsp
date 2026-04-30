package server

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/owenrumney/go-lsp/internal/debugui"
)

// ErrDebugTraceUnavailable is returned when trace export is requested without
// either WithDebugCapture or WithDebugUI being enabled.
var ErrDebugTraceUnavailable = errors.New("debug trace unavailable")

// ErrInvalidDebugTracePath is returned when the trace destination is not a
// writable regular file path.
var ErrInvalidDebugTracePath = errors.New("invalid debug trace path")

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
	return writeDebugTraceFile(path, data)
}

func writeDebugTraceFile(path string, data []byte) (err error) {
	cleanPath := filepath.Clean(path)
	dir, name := filepath.Split(cleanPath)
	if name == "" || name == "." || name == string(filepath.Separator) {
		return fmt.Errorf("%w: missing file name", ErrInvalidDebugTracePath)
	}
	if dir == "" {
		dir = "."
	}

	root, err := os.OpenRoot(dir)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := root.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if info, statErr := root.Lstat(name); statErr == nil {
		mode := info.Mode()
		if mode&fs.ModeSymlink != 0 || !mode.IsRegular() {
			return fmt.Errorf("%w: %s", ErrInvalidDebugTracePath, cleanPath)
		}
	} else if !errors.Is(statErr, fs.ErrNotExist) {
		return statErr
	}

	f, err := root.OpenFile(name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
