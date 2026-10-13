package client_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/client"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

// helperEnv makes the test binary re-exec itself as a stdio LSP server, so the
// child-process transport is covered without depending on an external server.
const helperEnv = "GO_LSP_CLIENT_TEST_SERVER"

// realServerEnv names a real language server to run the conformance checks
// against, for example "gopls". Those tests skip when it is unset, because
// this repo has no other mechanism for optional tests.
const realServerEnv = "GO_LSP_TEST_SERVER"

func TestMain(m *testing.M) {
	if os.Getenv(helperEnv) == "1" {
		runHelperServer()
		return
	}
	os.Exit(m.Run())
}

func runHelperServer() {
	srv := server.NewServer(&hoverHandler{})
	err := srv.Run(context.Background(), server.RunStdio())
	if err != nil && !errors.Is(err, context.Canceled) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}

func TestDialSpawnsChildProcess(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}

	c, err := client.Dial(t.Context(), client.Config{
		Command: []string{exe, "-test.run=^$"},
		Env:     []string{helperEnv + "=1"},
		RootURI: "file:///work",
		Stderr:  os.Stderr,
	})
	if err != nil {
		t.Fatalf("Dial: %v", err)
	}

	if got := c.ServerInfo(); got == nil || got.Name != "test-server" {
		t.Fatalf("ServerInfo() = %+v, want test-server", got)
	}

	hover, err := c.Hover(t.Context(), "file:///a.go", lsp.Position{})
	if err != nil {
		t.Fatalf("Hover: %v", err)
	}
	if hover == nil {
		t.Fatal("Hover() = nil over a real process")
	}

	if err := c.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}

	select {
	case <-c.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("child process never exited after Close")
	}
}

func TestDialReportsMissingBinary(t *testing.T) {
	_, err := client.Dial(t.Context(), client.Config{
		Command: []string{"go-lsp-no-such-server-binary"},
	})
	if err == nil {
		t.Fatal("Dial() = nil error for a binary that does not exist")
	}
	if !strings.Contains(err.Error(), "go-lsp-no-such-server-binary") {
		t.Fatalf("error %q does not name the missing binary", err)
	}
}

func TestAgainstRealServer(t *testing.T) {
	command := os.Getenv(realServerEnv)
	if command == "" {
		t.Skipf("set %s to a language server command to run this", realServerEnv)
	}

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	c, err := client.Dial(t.Context(), client.Config{
		Command: strings.Fields(command),
		Dir:     dir,
		RootURI: lsp.DocumentURI("file://" + dir),
		Stderr:  os.Stderr,
	})
	if err != nil {
		t.Fatalf("Dial %q: %v", command, err)
	}
	t.Cleanup(func() { _ = c.Close() })

	ctx, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()

	if err := c.Ready(ctx); err != nil {
		t.Fatalf("Ready: %v", err)
	}

	if got := c.Capabilities(); got.TextDocumentSync == nil {
		t.Error("server advertised no textDocumentSync")
	}
	// gopls puts its whole build manifest in Version, so log the name only.
	if info := c.ServerInfo(); info != nil {
		t.Logf("server %q, position encoding %s", info.Name, c.PositionEncoding())
	}
}
