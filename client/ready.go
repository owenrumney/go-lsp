package client

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/owenrumney/go-lsp/lsp"
)

// readiness approximates whether a server has finished indexing, by tracking
// outstanding work-done progress tokens. Servers that never report progress
// are declared ready once grace has elapsed.
type readiness struct {
	mu      sync.Mutex
	active  map[string]struct{}
	pending map[string]time.Time
	seenAny bool
	grace   time.Duration
	start   time.Time
}

func newReadiness(grace time.Duration) *readiness {
	return &readiness{
		active:  make(map[string]struct{}),
		pending: make(map[string]time.Time),
		grace:   grace,
		start:   time.Now(),
	}
}

// begin starts the grace period, which runs from the end of initialize.
func (r *readiness) begin() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.start = time.Now()
}

// Ready blocks until the server appears to have finished its startup work.
//
// This is a heuristic: LSP has no readiness signal, so it watches work-done
// progress and falls back to Config.ReadyGrace for servers that report none.
// Requests are never blocked on it; call Ready when an empty result would be
// indistinguishable from a real answer.
func (c *Client) Ready(ctx context.Context) error {
	const poll = 25 * time.Millisecond

	for {
		if c.ready.isReady() {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-c.done:
			if err := c.Err(); err != nil {
				return err
			}
			return ErrClosed
		case <-time.After(poll):
		}
	}
}

func (r *readiness) isReady() bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	// A token that is created but never begun must not block forever.
	for token, created := range r.pending {
		if time.Since(created) >= r.grace {
			delete(r.pending, token)
		}
	}

	if !r.seenAny {
		return time.Since(r.start) >= r.grace
	}
	return len(r.active) == 0 && len(r.pending) == 0
}

// create records a token the server asked for before reporting against it.
func (r *readiness) create(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.seenAny = true
	r.pending[token] = time.Now()
}

func (r *readiness) observe(params *lsp.ProgressParams) {
	var value struct {
		Kind string `json:"kind"`
	}
	if err := json.Unmarshal(params.Value, &value); err != nil {
		return
	}

	token := string(params.Token)

	r.mu.Lock()
	defer r.mu.Unlock()

	switch value.Kind {
	case "begin":
		r.seenAny = true
		delete(r.pending, token)
		r.active[token] = struct{}{}
	case "end":
		delete(r.pending, token)
		delete(r.active, token)
	}
}
