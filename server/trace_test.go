package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

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
	s := NewServer(&mockHandler{})
	s.recorder = debugui.NewRecorder()

	s.recorder.Store().Add("client→server", []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"rootUri":"file:///Users/owen/project"}}`))

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

func TestSaveDebugTraceRejectsSymlinkTarget(t *testing.T) {
	s := NewServer(&mockHandler{})
	s.recorder = debugui.NewRecorder()

	dir := t.TempDir()
	target := filepath.Join(dir, "target.json")
	if err := os.WriteFile(target, []byte("existing"), 0o600); err != nil {
		t.Fatal(err)
	}

	link := filepath.Join(dir, "trace.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}

	err := s.SaveDebugTrace(link, TraceExportOptions{})
	if !errors.Is(err, ErrInvalidDebugTracePath) {
		t.Fatalf("error = %v, want ErrInvalidDebugTracePath", err)
	}
}

func TestSaveDebugTraceRejectsDirectoryTarget(t *testing.T) {
	s := NewServer(&mockHandler{})
	s.recorder = debugui.NewRecorder()

	err := s.SaveDebugTrace(t.TempDir(), TraceExportOptions{})
	if !errors.Is(err, ErrInvalidDebugTracePath) {
		t.Fatalf("error = %v, want ErrInvalidDebugTracePath", err)
	}
}

// signalingPipe is an io.ReadWriteCloser that closes started on the first Read
// (so the test knows Run has reached its serve loop) and blocks Read until
// Close.
type signalingPipe struct {
	started chan struct{}
	closed  chan struct{}
	mu      sync.Mutex
	once    bool
}

func newSignalingPipe() *signalingPipe {
	return &signalingPipe{started: make(chan struct{}), closed: make(chan struct{})}
}

func (p *signalingPipe) Read(_ []byte) (int, error) {
	p.mu.Lock()
	if !p.once {
		p.once = true
		close(p.started)
	}
	p.mu.Unlock()
	<-p.closed
	return 0, io.EOF
}

func (p *signalingPipe) Write(b []byte) (int, error) { return len(b), nil }

func (p *signalingPipe) Close() error {
	select {
	case <-p.closed:
	default:
		close(p.closed)
	}
	return nil
}

func TestWithDebugCaptureEnablesTraceWithoutHTTP(t *testing.T) {
	s := NewServer(&mockHandler{}, WithDebugCapture())

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	rw := newSignalingPipe()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx, rw) }()

	waitForStart(t, rw)

	if s.debugUI != nil {
		t.Fatal("WithDebugCapture should not start the HTTP debug UI")
	}
	if _, err := s.ExportDebugTrace(TraceExportOptions{}); err != nil {
		t.Fatalf("ExportDebugTrace: %v", err)
	}

	cancel()
	_ = rw.Close()
	<-done
}

func TestWithDebugUIDegradesGracefullyOnBindFailure(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ln.Close() }()
	addr := ln.Addr().String()

	var logs lockedBuffer
	logger := slog.New(slog.NewTextHandler(&logs, &slog.HandlerOptions{Level: slog.LevelWarn}))

	s := NewServer(&mockHandler{}, WithDebugUI(addr), WithLogger(logger))

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	rw := newSignalingPipe()
	done := make(chan error, 1)
	go func() { done <- s.Run(ctx, rw) }()

	waitForStart(t, rw)

	if s.debugUI != nil {
		t.Fatal("debugUI should be nil after bind failure")
	}
	if !bytes.Contains(logs.Bytes(), []byte("HTTP UI unavailable")) {
		t.Fatalf("expected warning log, got: %s", logs.String())
	}
	if _, err := s.ExportDebugTrace(TraceExportOptions{}); err != nil {
		t.Fatalf("trace export should still work after bind failure: %v", err)
	}

	cancel()
	_ = rw.Close()
	<-done
}

func waitForStart(t *testing.T, p *signalingPipe) {
	t.Helper()
	select {
	case <-p.started:
	case <-time.After(2 * time.Second):
		t.Fatal("server did not enter serve loop within 2s")
	}
}

// lockedBuffer is a thread-safe bytes.Buffer for use as a log sink.
type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) Bytes() []byte {
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]byte(nil), b.buf.Bytes()...)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}
