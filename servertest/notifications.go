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
	s.cond.Broadcast()
}

func (s *notifStore) addLogMessage(params lsp.LogMessageParams) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logMessages = append(s.logMessages, params)
	s.cond.Broadcast()
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

// WaitForMessage waits until a window/showMessage notification is received.
func (h *Harness) WaitForMessage(ctx context.Context) (lsp.ShowMessageParams, error) {
	h.notifs.mu.Lock()
	defer h.notifs.mu.Unlock()

	for {
		if len(h.notifs.messages) > 0 {
			return h.notifs.messages[len(h.notifs.messages)-1], nil
		}
		if err := waitCond(ctx, h.notifs.cond); err != nil {
			return lsp.ShowMessageParams{}, err
		}
	}
}

// WaitForLogMessage waits until a window/logMessage notification is received.
func (h *Harness) WaitForLogMessage(ctx context.Context) (lsp.LogMessageParams, error) {
	h.notifs.mu.Lock()
	defer h.notifs.mu.Unlock()

	for {
		if len(h.notifs.logMessages) > 0 {
			return h.notifs.logMessages[len(h.notifs.logMessages)-1], nil
		}
		if err := waitCond(ctx, h.notifs.cond); err != nil {
			return lsp.LogMessageParams{}, err
		}
	}
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
		if err := waitCond(ctx, h.notifs.cond); err != nil {
			return nil, err
		}
	}
}

// ClearDiagnostics removes all collected diagnostics.
func (h *Harness) ClearDiagnostics() {
	h.notifs.mu.Lock()
	defer h.notifs.mu.Unlock()
	h.notifs.diagnostics = nil
}

func waitCond(ctx context.Context, cond *sync.Cond) error {
	stop := context.AfterFunc(ctx, func() {
		cond.L.Lock()
		cond.Broadcast()
		cond.L.Unlock()
	})
	defer stop()

	cond.Wait()
	if err := ctx.Err(); err != nil {
		return err
	}
	return nil
}
