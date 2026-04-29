package server

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/owenrumney/go-lsp/internal/debugui"
)

func TestExportDebugTraceUnavailable(t *testing.T) {
	s := NewServer(&mockHandler{})
	_, err := s.ExportDebugTrace(TraceExportOptions{})
	if !errors.Is(err, ErrDebugTraceUnavailable) {
		t.Fatalf("error = %v, want ErrDebugTraceUnavailable", err)
	}
}

func TestSaveDebugTrace(t *testing.T) {
	logStore := debugui.NewLogStore()
	store := debugui.NewStore(logStore)
	s := NewServer(&mockHandler{})
	s.debugUI = debugui.New(":0", store, logStore)

	store.Add("client→server", []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"rootUri":"file:///Users/owen/project"}}`))

	path := filepath.Join(t.TempDir(), "trace.json")
	if err := s.SaveDebugTrace(path, TraceExportOptions{Pretty: true, RedactFilePaths: true}); err != nil {
		t.Fatal(err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("permissions = %v, want 0600", got)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var trace struct {
		Version  int `json:"version"`
		Messages []struct {
			Body json.RawMessage `json:"body"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(data, &trace); err != nil {
		t.Fatal(err)
	}
	if trace.Version != debugui.TraceVersion {
		t.Fatalf("version = %d, want %d", trace.Version, debugui.TraceVersion)
	}
	if len(trace.Messages) != 1 {
		t.Fatalf("messages = %d, want 1", len(trace.Messages))
	}
	if string(trace.Messages[0].Body) == "" {
		t.Fatal("expected message body")
	}
}
