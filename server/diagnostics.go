package server

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"
	"sync"
	"time"

	"github.com/owenrumney/go-lsp/document"
	"github.com/owenrumney/go-lsp/lsp"
)

// Diagnoser is implemented by handlers that want the server to run
// diagnostics for them. The server keeps the document store in sync, calls
// Diagnose after didOpen and after the debounce window following didChange,
// drops results for versions that have since changed, and clears diagnostics
// on didClose. Requires WithDocumentStore; NewServer panics otherwise.
//
// doc is a snapshot and can be read without locking. ctx is cancelled when a
// newer edit arrives, the document closes, or the server shuts down.
// Returning ctx.Err() for a cancelled run is normal and is not logged.
type Diagnoser interface {
	Diagnose(ctx context.Context, doc *document.Document) ([]lsp.Diagnostic, error)
}

// DefaultDiagnosticsDebounce is the idle time after the last didChange before
// Diagnose runs when WithDiagnosticsDebounce is not set.
const DefaultDiagnosticsDebounce = 300 * time.Millisecond

// diagnosticsStopTimeout bounds how long shutdown waits for in-flight runs.
const diagnosticsStopTimeout = 2 * time.Second

type diagnosticsEngine struct {
	store    *document.Store
	diagnose Diagnoser
	client   *Client
	logger   *slog.Logger
	debounce time.Duration

	mu      sync.Mutex
	docs    map[lsp.DocumentURI]*docState
	stopped bool
	wg      sync.WaitGroup
}

type docState struct {
	timer  *time.Timer
	cancel context.CancelFunc
}

func newDiagnosticsEngine(store *document.Store, d Diagnoser, client *Client, logger *slog.Logger, debounce time.Duration) *diagnosticsEngine {
	return &diagnosticsEngine{
		store:    store,
		diagnose: d,
		client:   client,
		logger:   logger,
		debounce: debounce,
		docs:     make(map[lsp.DocumentURI]*docState),
	}
}

// opened runs Diagnose immediately.
func (e *diagnosticsEngine) opened(uri lsp.DocumentURI) {
	e.schedule(uri, 0)
}

// changed cancels any in-flight run and restarts the debounce timer.
func (e *diagnosticsEngine) changed(uri lsp.DocumentURI) {
	e.schedule(uri, e.debounce)
}

func (e *diagnosticsEngine) schedule(uri lsp.DocumentURI, delay time.Duration) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.stopped {
		return
	}
	e.reset(uri)

	ctx, cancel := context.WithCancel(context.Background())
	state := &docState{cancel: cancel}
	e.docs[uri] = state

	e.wg.Add(1)
	state.timer = time.AfterFunc(delay, func() {
		defer e.wg.Done()
		// This runs outside the dispatcher, so a Diagnose panic would kill
		// the process without this recover.
		defer func() {
			if r := recover(); r != nil && e.logger != nil {
				e.logger.Error("diagnose panic", "uri", uri, "panic", r, "stack", string(debug.Stack()))
			}
		}()
		e.run(ctx, uri, state)
	})
}

// reset cancels and stops the state for uri. Caller holds mu.
func (e *diagnosticsEngine) reset(uri lsp.DocumentURI) {
	state, ok := e.docs[uri]
	if !ok {
		return
	}
	state.cancel()
	// Stop reports false once the callback has started; the callback then
	// owns the wg.Done.
	if state.timer.Stop() {
		e.wg.Done()
	}
	delete(e.docs, uri)
}

func (e *diagnosticsEngine) run(ctx context.Context, uri lsp.DocumentURI, state *docState) {
	doc, ok := e.store.Get(uri)
	if !ok {
		return
	}
	diags, err := e.diagnose.Diagnose(ctx, doc)
	if err != nil {
		if !errors.Is(err, context.Canceled) && e.logger != nil {
			e.logger.Error("diagnose error", "uri", uri, "version", doc.Version(), "error", err)
		}
		return
	}
	if diags == nil {
		diags = []lsp.Diagnostic{}
	}
	version := doc.Version()

	// Publish under mu so a concurrent closed() cannot clear diagnostics
	// between these checks and the notification.
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.stopped || ctx.Err() != nil || e.docs[uri] != state {
		return
	}
	if v, ok := e.store.Version(uri); !ok || v != version {
		return
	}
	_ = e.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
		URI:         uri,
		Version:     &version,
		Diagnostics: diags,
	})
}

// closed cancels any pending run and clears the client's diagnostics for uri.
func (e *diagnosticsEngine) closed(ctx context.Context, uri lsp.DocumentURI) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.reset(uri)
	if e.stopped {
		return
	}
	_ = e.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
		URI:         uri,
		Diagnostics: []lsp.Diagnostic{},
	})
}

// stop cancels every run, stops every timer, and waits up to timeout for the
// in-flight goroutines. A Diagnose that ignores ctx must not hang shutdown.
func (e *diagnosticsEngine) stop(timeout time.Duration) {
	e.mu.Lock()
	e.stopped = true
	for uri := range e.docs {
		e.reset(uri)
	}
	e.mu.Unlock()

	done := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(timeout):
		if e.logger != nil {
			e.logger.Warn("diagnostics did not stop in time", "timeout", timeout)
		}
	}
}
