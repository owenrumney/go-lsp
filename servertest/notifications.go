package servertest

import (
	"context"
	"sync"

	"github.com/owenrumney/go-lsp/lsp"
)

type notifStore struct {
	mu          sync.Mutex
	cond        *sync.Cond
	diagnostics []lsp.PublishDiagnosticsParams
	messages    []lsp.ShowMessageParams
	logMessages []lsp.LogMessageParams
}

func newNotifStore() *notifStore {
	s := &notifStore{}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *notifStore) addDiagnostics(params lsp.PublishDiagnosticsParams) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.diagnostics = append(s.diagnostics, params)
	s.cond.Broadcast()
}

func (s *notifStore) addMessage(params lsp.ShowMessageParams) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages = append(s.messages, params)
}

func (s *notifStore) addLogMessage(params lsp.LogMessageParams) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logMessages = append(s.logMessages, params)
}

// Diagnostics returns the most recent diagnostics for the given URI.
func (h *Harness) Diagnostics(uri lsp.DocumentURI) []lsp.Diagnostic {
	h.notifs.mu.Lock()
	defer h.notifs.mu.Unlock()
	// Return the last published diagnostics for this URI.
	for i := len(h.notifs.diagnostics) - 1; i >= 0; i-- {
		if h.notifs.diagnostics[i].URI == uri {
			return h.notifs.diagnostics[i].Diagnostics
		}
	}
	return nil
}

// AllDiagnostics returns all collected diagnostics notifications.
func (h *Harness) AllDiagnostics() []lsp.PublishDiagnosticsParams {
	h.notifs.mu.Lock()
	defer h.notifs.mu.Unlock()
	result := make([]lsp.PublishDiagnosticsParams, len(h.notifs.diagnostics))
	copy(result, h.notifs.diagnostics)
	return result
}

// Messages returns all collected showMessage notifications.
func (h *Harness) Messages() []lsp.ShowMessageParams {
	h.notifs.mu.Lock()
	defer h.notifs.mu.Unlock()
	result := make([]lsp.ShowMessageParams, len(h.notifs.messages))
	copy(result, h.notifs.messages)
	return result
}

// LogMessages returns all collected logMessage notifications.
func (h *Harness) LogMessages() []lsp.LogMessageParams {
	h.notifs.mu.Lock()
	defer h.notifs.mu.Unlock()
	result := make([]lsp.LogMessageParams, len(h.notifs.logMessages))
	copy(result, h.notifs.logMessages)
	return result
}

// WaitForDiagnostics waits until diagnostics are published for the given URI.
// If diagnostics already exist, they are returned immediately.
func (h *Harness) WaitForDiagnostics(ctx context.Context, uri lsp.DocumentURI) ([]lsp.Diagnostic, error) {
	h.notifs.mu.Lock()
	defer h.notifs.mu.Unlock()

	for {
		// Check if we already have diagnostics for this URI.
		for i := len(h.notifs.diagnostics) - 1; i >= 0; i-- {
			if h.notifs.diagnostics[i].URI == uri {
				return h.notifs.diagnostics[i].Diagnostics, nil
			}
		}

		// Wait in a goroutine so we can also respect context cancellation.
		waitDone := make(chan struct{})
		go func() {
			h.notifs.mu.Lock()
			h.notifs.cond.Wait()
			h.notifs.mu.Unlock()
			close(waitDone)
		}()

		h.notifs.mu.Unlock()
		select {
		case <-waitDone:
			h.notifs.mu.Lock()
			continue
		case <-ctx.Done():
			// Wake up the waiting goroutine so it doesn't leak.
			h.notifs.cond.Broadcast()
			<-waitDone
			h.notifs.mu.Lock()
			return nil, ctx.Err()
		}
	}
}

// ClearDiagnostics removes all collected diagnostics.
func (h *Harness) ClearDiagnostics() {
	h.notifs.mu.Lock()
	defer h.notifs.mu.Unlock()
	h.notifs.diagnostics = nil
}
