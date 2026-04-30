package debugui

import (
	"encoding/json"
	"io"
	"sync"
	"time"
)

// Recorder captures LSP traffic and logs without serving any HTTP UI. It owns
// the message Store, LogStore, and a snapshot of advertised capabilities, and
// is the foundation that both capture-only mode and the full DebugUI build on.
type Recorder struct {
	store    *Store
	logStore *LogStore

	capsMu       sync.RWMutex
	capabilities json.RawMessage
}

// NewRecorder creates a Recorder with empty stores.
func NewRecorder() *Recorder {
	logStore := NewLogStore()
	store := NewStore(logStore)
	return &Recorder{
		store:    store,
		logStore: logStore,
	}
}

// Store returns the underlying message store.
func (r *Recorder) Store() *Store { return r.store }

// LogStore returns the underlying log store.
func (r *Recorder) LogStore() *LogStore { return r.logStore }

// Tap wraps inner so all framed LSP messages flowing through it are captured.
func (r *Recorder) Tap(inner io.ReadWriteCloser) *Tap {
	return NewTap(inner, r.store)
}

// SlogHandler returns a slog.Handler that captures log records into the trace.
func (r *Recorder) SlogHandler() *SlogHandler {
	return NewSlogHandler(r.logStore)
}

// SetCapabilities stores the server's advertised capabilities so they are
// included in exported traces.
func (r *Recorder) SetCapabilities(caps any) {
	data, err := json.Marshal(caps)
	if err != nil {
		return
	}
	r.capsMu.Lock()
	r.capabilities = data
	r.capsMu.Unlock()
}

// ExportTrace returns a JSON snapshot of captured messages, logs, and
// capabilities, applying the requested redactions.
func (r *Recorder) ExportTrace(opts TraceExportOptions) ([]byte, error) {
	trace := Trace{
		Version:      TraceVersion,
		CreatedAt:    time.Now().UTC(),
		Messages:     r.store.All(),
		Logs:         r.logStore.All(),
		Capabilities: r.capabilitiesSnapshot(),
	}

	redactTrace(&trace, opts)

	if opts.Pretty {
		return json.MarshalIndent(trace, "", "  ")
	}
	return json.Marshal(trace)
}

func (r *Recorder) capabilitiesSnapshot() json.RawMessage {
	r.capsMu.RLock()
	defer r.capsMu.RUnlock()
	if r.capabilities == nil {
		return nil
	}
	return append(json.RawMessage(nil), r.capabilities...)
}
