package debugui

import (
	"encoding/json"
	"strings"
	"sync"
	"time"
)

const maxEntries = 10000

// Entry is a captured LSP message.
type Entry struct {
	ID         int             `json:"id"`
	Timestamp  time.Time       `json:"timestamp"`
	Direction  string          `json:"direction"`
	MsgType    string          `json:"msgType"`
	Method     string          `json:"method"`
	RPCID      string          `json:"rpcId"`
	Body       json.RawMessage `json:"body"`
	PairedWith int             `json:"pairedWith"`
}

// Subscriber receives new entries.
type Subscriber func(Entry)

// Store is a thread-safe ring buffer of captured LSP messages with request/response correlation.
type Store struct {
	mu          sync.RWMutex
	entries     []Entry
	nextID      int
	pending     map[string]int // rpcID -> entry ID for unmatched requests
	subscribers []Subscriber
	logStore    *LogStore // optional: cross-post window/logMessage notifications
}

// NewStore creates a new Store. If logStore is non-nil, window/logMessage
// notifications are automatically cross-posted to it.
func NewStore(logStore *LogStore) *Store {
	return &Store{
		entries:  make([]Entry, 0, 256),
		pending:  make(map[string]int),
		logStore: logStore,
	}
}

// Subscribe registers a callback for new entries.
func (s *Store) Subscribe(fn Subscriber) {
	s.mu.Lock()
	s.subscribers = append(s.subscribers, fn)
	s.mu.Unlock()
}

// Add decodes a raw JSON-RPC message, correlates it, stores it, and notifies subscribers.
func (s *Store) Add(direction string, raw []byte) {
	e := Entry{
		Timestamp:  time.Now(),
		Direction:  direction,
		Body:       json.RawMessage(append([]byte(nil), raw...)),
		PairedWith: -1,
	}

	// Decode to classify the message.
	var msg struct {
		ID     *json.RawMessage `json:"id,omitempty"`
		Method *string          `json:"method,omitempty"`
	}
	_ = json.Unmarshal(raw, &msg)

	hasID := msg.ID != nil && string(*msg.ID) != "null"

	switch {
	case hasID && msg.Method != nil:
		e.MsgType = "request"
		e.Method = *msg.Method
		e.RPCID = trimQuotes(string(*msg.ID))
	case hasID && msg.Method == nil:
		e.MsgType = "response"
		e.RPCID = trimQuotes(string(*msg.ID))
	case msg.Method != nil:
		e.MsgType = "notification"
		e.Method = *msg.Method
	}

	s.mu.Lock()

	e.ID = s.nextID
	s.nextID++

	// Correlate request/response pairs.
	if e.RPCID != "" {
		switch e.MsgType {
		case "request":
			s.pending[e.RPCID] = e.ID
		case "response":
			if reqID, ok := s.pending[e.RPCID]; ok {
				e.PairedWith = reqID
				// Update the request entry's PairedWith to point to this response.
				s.updatePairedWith(reqID, e.ID)
				delete(s.pending, e.RPCID)
			}
		}
	}

	if len(s.entries) < maxEntries {
		s.entries = append(s.entries, e)
	} else {
		s.entries[e.ID%maxEntries] = e
	}

	subs := make([]Subscriber, len(s.subscribers))
	copy(subs, s.subscribers)
	s.mu.Unlock()

	for _, fn := range subs {
		fn(e)
	}

	// Cross-post window/logMessage notifications to the log store.
	if s.logStore != nil && e.Method == "window/logMessage" {
		var params struct {
			Type    int    `json:"type"`
			Message string `json:"message"`
		}
		// The body is the full JSON-RPC message; extract params.
		var rpc struct {
			Params json.RawMessage `json:"params"`
		}
		if json.Unmarshal(e.Body, &rpc) == nil && rpc.Params != nil {
			if json.Unmarshal(rpc.Params, &params) == nil {
				level := "info"
				switch params.Type {
				case 1:
					level = "error"
				case 2:
					level = "warning"
				case 4:
					level = "log"
				}
				s.logStore.Add(level, params.Message)
			}
		}
	}
}

func (s *Store) updatePairedWith(entryID, pairedID int) {
	for i := range s.entries {
		if s.entries[i].ID == entryID {
			s.entries[i].PairedWith = pairedID
			return
		}
	}
}

// Clear removes all entries and resets correlation state.
func (s *Store) Clear() {
	s.mu.Lock()
	s.entries = s.entries[:0]
	s.nextID = 0
	s.pending = make(map[string]int)
	s.mu.Unlock()
}

// Entries returns a paginated slice of entries.
func (s *Store) Entries(offset, limit int) []Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := len(s.entries)
	if offset >= n {
		return nil
	}
	end := min(offset+limit, n)
	result := make([]Entry, end-offset)
	copy(result, s.entries[offset:end])
	return result
}

// Search returns entries where method or body contains the query substring.
func (s *Store) Search(query string) []Entry {
	query = strings.ToLower(query)
	s.mu.RLock()
	defer s.mu.RUnlock()

	var result []Entry
	for _, e := range s.entries {
		if strings.Contains(strings.ToLower(e.Method), query) ||
			strings.Contains(strings.ToLower(string(e.Body)), query) {
			result = append(result, e)
		}
	}
	return result
}

// Entry returns a single entry by ID, or nil if not found.
func (s *Store) Entry(id int) *Entry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.entries {
		if s.entries[i].ID == id {
			e := s.entries[i]
			return &e
		}
	}
	return nil
}

func trimQuotes(s string) string {
	if len(s) >= 2 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}
