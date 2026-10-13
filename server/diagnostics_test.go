package server_test

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
	"github.com/owenrumney/go-lsp/servertest"
)

type diagnoseFunc func(ctx context.Context, doc *document.Document) ([]lsp.Diagnostic, error)

type diagHandler struct {
	diagnose diagnoseFunc
}

func (h *diagHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{}, nil
}

func (h *diagHandler) Shutdown(_ context.Context) error { return nil }

func (h *diagHandler) Diagnose(ctx context.Context, doc *document.Document) ([]lsp.Diagnostic, error) {
	return h.diagnose(ctx, doc)
}

type lockedBuffer struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (b *lockedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.Write(p)
}

func (b *lockedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buf.String()
}

func oneDiag(msg string) []lsp.Diagnostic {
	return []lsp.Diagnostic{{Message: msg}}
}

func versionOf(p lsp.PublishDiagnosticsParams) int {
	if p.Version == nil {
		return -1
	}
	return *p.Version
}

func newDiagHarness(t *testing.T, fn diagnoseFunc, debounce time.Duration, opts ...server.Option) *servertest.Harness {
	t.Helper()
	opts = append(opts, server.WithDocumentStore(document.NewStore()), server.WithDiagnosticsDebounce(debounce))
	return servertest.New(t, &diagHandler{diagnose: fn}, servertest.WithServerOptions(opts...))
}

func waitForPublish(t *testing.T, h *servertest.Harness, uri lsp.DocumentURI, version int) lsp.PublishDiagnosticsParams {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		for _, p := range h.AllDiagnostics() {
			if p.URI == uri && versionOf(p) == version {
				return p
			}
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("no publishDiagnostics for %s version %d; got %+v", uri, version, h.AllDiagnostics())
	return lsp.PublishDiagnosticsParams{}
}

func waitFor(t *testing.T, cond func() bool, what string) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(5 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", what)
}

func publishedVersions(h *servertest.Harness, uri lsp.DocumentURI) []int {
	var out []int
	for _, p := range h.AllDiagnostics() {
		if p.URI == uri {
			out = append(out, versionOf(p))
		}
	}
	return out
}

func TestDiagnosticsOpenPublishesImmediately(t *testing.T) {
	h := newDiagHarness(t, func(_ context.Context, doc *document.Document) ([]lsp.Diagnostic, error) {
		return oneDiag(doc.Text()), nil
	}, 10*time.Second)
	uri := lsp.DocumentURI("file:///a.txt")

	if err := h.DidOpen(uri, "txt", "hello"); err != nil {
		t.Fatal(err)
	}
	p := waitForPublish(t, h, uri, 1)
	if len(p.Diagnostics) != 1 || p.Diagnostics[0].Message != "hello" {
		t.Fatalf("diagnostics = %+v", p.Diagnostics)
	}
}

func TestDiagnosticsRapidChangesPublishOnce(t *testing.T) {
	var mu sync.Mutex
	var calls []int
	h := newDiagHarness(t, func(_ context.Context, doc *document.Document) ([]lsp.Diagnostic, error) {
		mu.Lock()
		calls = append(calls, doc.Version())
		mu.Unlock()
		return oneDiag(doc.Text()), nil
	}, 200*time.Millisecond)
	uri := lsp.DocumentURI("file:///a.txt")

	if err := h.DidOpen(uri, "txt", "v1"); err != nil {
		t.Fatal(err)
	}
	waitForPublish(t, h, uri, 1)

	for v := 2; v <= 4; v++ {
		if err := h.DidChange(uri, v, "v"+string(rune('0'+v))); err != nil {
			t.Fatal(err)
		}
	}
	p := waitForPublish(t, h, uri, 4)
	if p.Diagnostics[0].Message != "v4" {
		t.Fatalf("message = %q", p.Diagnostics[0].Message)
	}
	time.Sleep(400 * time.Millisecond)

	if got := publishedVersions(h, uri); len(got) != 2 || got[0] != 1 || got[1] != 4 {
		t.Fatalf("published versions = %v, want [1 4]", got)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(calls) != 2 || calls[1] != 4 {
		t.Fatalf("diagnose calls = %v, want [1 4]", calls)
	}
}

func TestDiagnosticsZeroDebounceRunsEveryChange(t *testing.T) {
	h := newDiagHarness(t, func(_ context.Context, doc *document.Document) ([]lsp.Diagnostic, error) {
		return oneDiag(doc.Text()), nil
	}, 0)
	uri := lsp.DocumentURI("file:///a.txt")

	if err := h.DidOpen(uri, "txt", "v1"); err != nil {
		t.Fatal(err)
	}
	waitForPublish(t, h, uri, 1)
	if err := h.DidChange(uri, 2, "v2"); err != nil {
		t.Fatal(err)
	}
	waitForPublish(t, h, uri, 2)
}

func TestDiagnosticsChangeCancelsInFlightRun(t *testing.T) {
	started := make(chan struct{}, 1)
	var cancelled atomic.Bool
	h := newDiagHarness(t, func(ctx context.Context, doc *document.Document) ([]lsp.Diagnostic, error) {
		if doc.Version() == 1 {
			started <- struct{}{}
			<-ctx.Done()
			cancelled.Store(true)
			return oneDiag("stale"), nil
		}
		return oneDiag("fresh"), nil
	}, 50*time.Millisecond)
	uri := lsp.DocumentURI("file:///a.txt")

	if err := h.DidOpen(uri, "txt", "v1"); err != nil {
		t.Fatal(err)
	}
	<-started
	if err := h.DidChange(uri, 2, "v2"); err != nil {
		t.Fatal(err)
	}
	p := waitForPublish(t, h, uri, 2)
	if p.Diagnostics[0].Message != "fresh" {
		t.Fatalf("message = %q", p.Diagnostics[0].Message)
	}
	if !cancelled.Load() {
		t.Fatal("in-flight Diagnose did not see ctx.Done()")
	}
	if got := publishedVersions(h, uri); len(got) != 1 || got[0] != 2 {
		t.Fatalf("published versions = %v, want [2]", got)
	}
}

func TestDiagnosticsClosePublishesEmpty(t *testing.T) {
	h := newDiagHarness(t, func(_ context.Context, _ *document.Document) ([]lsp.Diagnostic, error) {
		return oneDiag("x"), nil
	}, 0)
	uri := lsp.DocumentURI("file:///a.txt")

	if err := h.DidOpen(uri, "txt", "v1"); err != nil {
		t.Fatal(err)
	}
	waitForPublish(t, h, uri, 1)
	if err := h.DidClose(uri); err != nil {
		t.Fatal(err)
	}
	p := waitForPublish(t, h, uri, -1)
	if len(p.Diagnostics) != 0 {
		t.Fatalf("diagnostics = %+v, want empty", p.Diagnostics)
	}

	var raw string
	for _, n := range h.Notifications() {
		if n.Method == "textDocument/publishDiagnostics" && !strings.Contains(string(n.Params), `"version"`) {
			raw = string(n.Params)
		}
	}
	if !strings.Contains(raw, `"diagnostics":[]`) {
		t.Fatalf("clear notification = %s, want diagnostics:[] and no version", raw)
	}
}

func TestDiagnosticsErrorIsLoggedAndKeepsPrevious(t *testing.T) {
	var logs lockedBuffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	h := newDiagHarness(t, func(_ context.Context, doc *document.Document) ([]lsp.Diagnostic, error) {
		if doc.Version() == 1 {
			return oneDiag("first"), nil
		}
		return nil, errors.New("boom")
	}, 0, server.WithLogger(logger))
	uri := lsp.DocumentURI("file:///a.txt")

	if err := h.DidOpen(uri, "txt", "v1"); err != nil {
		t.Fatal(err)
	}
	waitForPublish(t, h, uri, 1)
	if err := h.DidChange(uri, 2, "v2"); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool {
		return strings.Contains(logs.String(), "diagnose error") && strings.Contains(logs.String(), "boom")
	}, "error log")

	if got := publishedVersions(h, uri); len(got) != 1 || got[0] != 1 {
		t.Fatalf("published versions = %v, want [1]", got)
	}
	if d := h.Diagnostics(uri); len(d) != 1 || d[0].Message != "first" {
		t.Fatalf("diagnostics = %+v, want previous", d)
	}
}

func TestDiagnosticsCancelledErrorIsNotLogged(t *testing.T) {
	var logs lockedBuffer
	logger := slog.New(slog.NewTextHandler(&logs, nil))
	started := make(chan struct{}, 1)
	h := newDiagHarness(t, func(ctx context.Context, doc *document.Document) ([]lsp.Diagnostic, error) {
		if doc.Version() == 1 {
			started <- struct{}{}
			<-ctx.Done()
			return nil, ctx.Err()
		}
		return oneDiag("ok"), nil
	}, 0, server.WithLogger(logger))
	uri := lsp.DocumentURI("file:///a.txt")

	if err := h.DidOpen(uri, "txt", "v1"); err != nil {
		t.Fatal(err)
	}
	<-started
	if err := h.DidChange(uri, 2, "v2"); err != nil {
		t.Fatal(err)
	}
	waitForPublish(t, h, uri, 2)
	if strings.Contains(logs.String(), "diagnose error") {
		t.Fatalf("cancelled run was logged:\n%s", logs.String())
	}
}

func TestDiagnosticsShutdownCancelsAndWaits(t *testing.T) {
	started := make(chan struct{}, 1)
	var cancelled atomic.Bool
	t.Cleanup(func() {
		if !cancelled.Load() {
			t.Error("shutdown did not cancel the in-flight Diagnose")
		}
	})
	h := newDiagHarness(t, func(ctx context.Context, _ *document.Document) ([]lsp.Diagnostic, error) {
		started <- struct{}{}
		<-ctx.Done()
		cancelled.Store(true)
		return nil, ctx.Err()
	}, 0)

	if err := h.DidOpen("file:///a.txt", "txt", "v1"); err != nil {
		t.Fatal(err)
	}
	<-started
}

func TestNewServerPanicsForDiagnoserWithoutStore(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("expected panic")
		}
		if !strings.Contains(r.(string), "WithDocumentStore") {
			t.Fatalf("panic = %v", r)
		}
	}()
	server.NewServer(&diagHandler{})
}

type storeHookHandler struct {
	store  *document.Store
	opened chan string
}

func (h *storeHookHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{}, nil
}

func (h *storeHookHandler) Shutdown(_ context.Context) error { return nil }

func (h *storeHookHandler) DidOpen(_ context.Context, p *lsp.DidOpenTextDocumentParams) error {
	text, _ := h.store.Text(p.TextDocument.URI)
	h.opened <- text
	return nil
}

func (h *storeHookHandler) DidChange(_ context.Context, _ *lsp.DidChangeTextDocumentParams) error {
	return nil
}

func (h *storeHookHandler) DidClose(_ context.Context, _ *lsp.DidCloseTextDocumentParams) error {
	return nil
}

func TestWithDocumentStoreUpdatesStoreBeforeHook(t *testing.T) {
	store := document.NewStore()
	hh := &storeHookHandler{store: store, opened: make(chan string, 1)}
	h := servertest.New(t, hh, servertest.WithServerOptions(server.WithDocumentStore(store)))
	uri := lsp.DocumentURI("file:///a.txt")

	if err := h.DidOpen(uri, "txt", "hello"); err != nil {
		t.Fatal(err)
	}
	if got := <-hh.opened; got != "hello" {
		t.Fatalf("store text in DidOpen = %q, want hello", got)
	}
	if err := h.DidClose(uri); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool { _, ok := store.Get(uri); return !ok }, "store close")
}

type plainHandler struct{}

func (plainHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{}, nil
}

func (plainHandler) Shutdown(_ context.Context) error { return nil }

func TestWithDocumentStoreWithoutSyncHandler(t *testing.T) {
	store := document.NewStore()
	h := servertest.New(t, plainHandler{}, servertest.WithServerOptions(server.WithDocumentStore(store)))
	uri := lsp.DocumentURI("file:///a.txt")

	sync := h.InitResult.Capabilities.TextDocumentSync
	if sync == nil || sync.Change != lsp.SyncIncremental || sync.OpenClose == nil || !*sync.OpenClose {
		t.Fatalf("textDocumentSync = %+v, want incremental open/close", sync)
	}

	if err := h.DidOpen(uri, "txt", "hello"); err != nil {
		t.Fatal(err)
	}
	if err := h.DidChange(uri, 2, "changed"); err != nil {
		t.Fatal(err)
	}
	waitFor(t, func() bool { text, _ := store.Text(uri); return text == "changed" }, "store change")
	if v, _ := store.Version(uri); v != 2 {
		t.Fatalf("version = %d, want 2", v)
	}
}
