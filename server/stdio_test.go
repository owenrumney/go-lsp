package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

// TestMain re-execs the test binary as a stdio LSP server when the helper
// env var is set, so TestRunStdioExitsOnExitNotification can drive a real
// process over real pipes.
func TestMain(m *testing.M) {
	if os.Getenv("GO_LSP_STDIO_SERVER") == "1" {
		runStdioHelper()
		return
	}
	os.Exit(m.Run())
}

func runStdioHelper() {
	ctx, cancel := context.WithCancel(context.Background())
	srv := NewServer(&bareInitHandler{})
	// Lets the ctx-cancellation test cancel Run from inside the process
	// while the read loop is blocked on stdin.
	srv.HandleNotification("test/stop", func(context.Context, json.RawMessage) error {
		cancel()
		return nil
	})
	err := srv.Run(ctx, RunStdio())
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

type pipePair struct {
	io.Reader
	io.Writer
}

func (pipePair) Close() error { return nil }

type stdioServer struct {
	cmd  *exec.Cmd
	conn *jsonrpc.Conn
}

func startStdioServer(t *testing.T) *stdioServer {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	cmd := exec.Command(exe, "-test.run=^$")
	cmd.Env = append(os.Environ(), "GO_LSP_STDIO_SERVER=1")
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_ = stdin.Close()
	})

	return &stdioServer{
		cmd:  cmd,
		conn: jsonrpc.NewConn(pipePair{Reader: stdout, Writer: stdin}, jsonrpc.NewDispatcher()),
	}
}

func (s *stdioServer) call(t *testing.T, id int64, method string, params any) {
	t.Helper()
	req, err := jsonrpc.NewRequest(jsonrpc.IntID(id), method, params)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.conn.WriteMessage(req); err != nil {
		t.Fatal(err)
	}
	msg, err := s.conn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	resp, ok := msg.(*jsonrpc.Response)
	if !ok || resp.Error != nil {
		t.Fatalf("%s failed: %+v", method, msg)
	}
}

func (s *stdioServer) notify(t *testing.T, method string) {
	t.Helper()
	notif, err := jsonrpc.NewNotification(method, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.conn.WriteMessage(notif); err != nil {
		t.Fatal(err)
	}
}

// assertExits fails unless the process terminates cleanly while stdin is
// still held open.
func (s *stdioServer) assertExits(t *testing.T, within time.Duration) {
	t.Helper()
	waitDone := make(chan error, 1)
	go func() { waitDone <- s.cmd.Wait() }()
	select {
	case err := <-waitDone:
		if err != nil {
			t.Fatalf("server exited with error: %v", err)
		}
	case <-time.After(within):
		t.Fatal("server did not exit with stdin held open")
	}
}

func TestRunStdioExitsOnExitNotification(t *testing.T) {
	s := startStdioServer(t)
	s.call(t, 1, "initialize", &lsp.InitializeParams{})
	s.call(t, 2, "shutdown", nil)
	s.notify(t, "exit")
	s.assertExits(t, 10*time.Second)
}

// The regression this guards: cancelling Run's context must unblock the
// stdin read, which requires StdRWC.Close to actually close os.Stdin.
func TestRunStdioReturnsOnContextCancel(t *testing.T) {
	s := startStdioServer(t)
	s.call(t, 1, "initialize", &lsp.InitializeParams{})
	s.notify(t, "test/stop")
	s.assertExits(t, 10*time.Second)
}
