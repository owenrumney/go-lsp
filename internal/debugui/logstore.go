package debugui

import (
	"context"
	"log/slog"
	"strings"
	"sync"
	"time"
)

const maxLogEntries = 5000

// LogEntry is a captured log message.
type LogEntry struct {
	ID        int       `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // "error", "warning", "info", "debug"
	Message   string    `json:"message"`
}

// LogSubscriber receives new log entries.
type LogSubscriber func(LogEntry)

// LogStore is a thread-safe ring buffer of log messages.
type LogStore struct {
	mu          sync.RWMutex
	entries     []LogEntry
	nextID      int
	subscribers []LogSubscriber
}

// NewLogStore creates a new LogStore.
func NewLogStore() *LogStore {
	return &LogStore{
		entries: make([]LogEntry, 0, 256),
	}
}

// Subscribe registers a callback for new log entries.
func (s *LogStore) Subscribe(fn LogSubscriber) {
	s.mu.Lock()
	s.subscribers = append(s.subscribers, fn)
	s.mu.Unlock()
}

// Add stores a log entry and notifies subscribers.
func (s *LogStore) Add(level, message string) {
	e := LogEntry{
		Timestamp: time.Now(),
		Level:     level,
		Message:   message,
	}

	s.mu.Lock()

	e.ID = s.nextID
	s.nextID++

	if len(s.entries) < maxLogEntries {
		s.entries = append(s.entries, e)
	} else {
		s.entries[e.ID%maxLogEntries] = e
	}

	subs := make([]LogSubscriber, len(s.subscribers))
	copy(subs, s.subscribers)
	s.mu.Unlock()

	for _, fn := range subs {
		fn(e)
	}
}

// Clear removes all log entries.
func (s *LogStore) Clear() {
	s.mu.Lock()
	s.entries = s.entries[:0]
	s.nextID = 0
	s.mu.Unlock()
}

// Entries returns a paginated slice of log entries.
func (s *LogStore) Entries(offset, limit int) []LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := len(s.entries)
	if offset >= n {
		return nil
	}
	end := min(offset+limit, n)
	result := make([]LogEntry, end-offset)
	copy(result, s.entries[offset:end])
	return result
}

// Search returns log entries where the message contains the query substring.
func (s *LogStore) Search(query string) []LogEntry {
	query = strings.ToLower(query)
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []LogEntry
	for _, e := range s.entries {
		if strings.Contains(strings.ToLower(e.Message), query) ||
			strings.Contains(strings.ToLower(e.Level), query) {
			result = append(result, e)
		}
	}
	return result
}

// SlogHandler is a slog.Handler that sends log records to a LogStore.
type SlogHandler struct {
	store *LogStore
	attrs []slog.Attr
	group string
}

// NewSlogHandler creates a slog.Handler that writes to the given LogStore.
func NewSlogHandler(store *LogStore) *SlogHandler {
	return &SlogHandler{store: store}
}

func (h *SlogHandler) Enabled(_ context.Context, _ slog.Level) bool {
	return true
}

func (h *SlogHandler) Handle(_ context.Context, r slog.Record) error {
	level := strings.ToLower(r.Level.String())

	msg := r.Message
	// Append attrs if any.
	var parts []string
	for _, a := range h.attrs {
		parts = append(parts, formatAttr(h.group, a))
	}
	r.Attrs(func(a slog.Attr) bool {
		parts = append(parts, formatAttr(h.group, a))
		return true
	})
	if len(parts) > 0 {
		msg += " " + strings.Join(parts, " ")
	}

	h.store.Add(level, msg)
	return nil
}

func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &SlogHandler{
		store: h.store,
		attrs: append(append([]slog.Attr{}, h.attrs...), attrs...),
		group: h.group,
	}
}

func (h *SlogHandler) WithGroup(name string) slog.Handler {
	prefix := name
	if h.group != "" {
		prefix = h.group + "." + name
	}
	return &SlogHandler{
		store: h.store,
		attrs: append([]slog.Attr{}, h.attrs...),
		group: prefix,
	}
}

func formatAttr(group string, a slog.Attr) string {
	key := a.Key
	if group != "" {
		key = group + "." + key
	}
	return key + "=" + a.Value.String()
}
