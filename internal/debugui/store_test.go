package debugui

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestStoreAddAndEntries(t *testing.T) {
	s := NewStore(nil)

	s.Add("client→server", []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`))
	s.Add("server→client", []byte(`{"jsonrpc":"2.0","id":1,"result":{}}`))
	s.Add("client→server", []byte(`{"jsonrpc":"2.0","method":"initialized","params":{}}`))

	entries := s.Entries(0, 10)
	if len(entries) != 3 {
		t.Fatalf("got %d entries, want 3", len(entries))
	}

	tests := []struct {
		idx        int
		msgType    string
		method     string
		rpcID      string
		pairedWith int
	}{
		{0, "request", "initialize", "1", 1},
		{1, "response", "", "1", 0},
		{2, "notification", "initialized", "", -1},
	}

	for _, tt := range tests {
		e := entries[tt.idx]
		if e.MsgType != tt.msgType {
			t.Errorf("entry %d: msgType = %q, want %q", tt.idx, e.MsgType, tt.msgType)
		}
		if e.Method != tt.method {
			t.Errorf("entry %d: method = %q, want %q", tt.idx, e.Method, tt.method)
		}
		if e.RPCID != tt.rpcID {
			t.Errorf("entry %d: rpcID = %q, want %q", tt.idx, e.RPCID, tt.rpcID)
		}
		if e.PairedWith != tt.pairedWith {
			t.Errorf("entry %d: pairedWith = %d, want %d", tt.idx, e.PairedWith, tt.pairedWith)
		}
	}
}

func TestStoreSubscriber(t *testing.T) {
	s := NewStore(nil)

	var got []Entry
	s.Subscribe(func(e Entry) {
		got = append(got, e)
	})

	s.Add("client→server", []byte(`{"jsonrpc":"2.0","id":1,"method":"hover","params":{}}`))

	if len(got) != 1 {
		t.Fatalf("subscriber got %d entries, want 1", len(got))
	}
	if got[0].Method != "hover" {
		t.Errorf("method = %q, want hover", got[0].Method)
	}
}

func TestStoreSearch(t *testing.T) {
	s := NewStore(nil)

	s.Add("client→server", []byte(`{"jsonrpc":"2.0","id":1,"method":"textDocument/hover","params":{}}`))
	s.Add("client→server", []byte(`{"jsonrpc":"2.0","id":2,"method":"textDocument/completion","params":{}}`))
	s.Add("client→server", []byte(`{"jsonrpc":"2.0","method":"initialized","params":{}}`))

	results := s.Search("hover")
	if len(results) != 1 {
		t.Fatalf("got %d results, want 1", len(results))
	}
	if results[0].Method != "textDocument/hover" {
		t.Errorf("method = %q, want textDocument/hover", results[0].Method)
	}
}

func TestStoreRingBuffer(t *testing.T) {
	s := NewStore(nil)

	for i := 0; i < maxEntries+100; i++ {
		msg := fmt.Sprintf(`{"jsonrpc":"2.0","method":"test/%d","params":{}}`, i)
		s.Add("client→server", []byte(msg))
	}

	entries := s.Entries(0, maxEntries)
	if len(entries) != maxEntries {
		t.Fatalf("got %d entries, want %d", len(entries), maxEntries)
	}
}

func TestStorePagination(t *testing.T) {
	s := NewStore(nil)

	for i := 0; i < 10; i++ {
		msg := fmt.Sprintf(`{"jsonrpc":"2.0","method":"test/%d","params":{}}`, i)
		s.Add("client→server", []byte(msg))
	}

	entries := s.Entries(5, 3)
	if len(entries) != 3 {
		t.Fatalf("got %d entries, want 3", len(entries))
	}
	if entries[0].ID != 5 {
		t.Errorf("first entry ID = %d, want 5", entries[0].ID)
	}

	entries = s.Entries(8, 10)
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}

	entries = s.Entries(100, 10)
	if entries != nil {
		t.Fatalf("got entries for out-of-range offset")
	}
}

func TestStoreCrossPostLogMessage(t *testing.T) {
	logStore := NewLogStore()
	s := NewStore(logStore)

	// Simulate a window/logMessage notification.
	s.Add("server→client", []byte(`{"jsonrpc":"2.0","method":"window/logMessage","params":{"type":3,"message":"hello from server"}}`))
	s.Add("server→client", []byte(`{"jsonrpc":"2.0","method":"window/logMessage","params":{"type":1,"message":"something broke"}}`))
	// Non-logMessage notification should not cross-post.
	s.Add("server→client", []byte(`{"jsonrpc":"2.0","method":"textDocument/publishDiagnostics","params":{}}`))

	logs := logStore.Entries(0, 10)
	if len(logs) != 2 {
		t.Fatalf("got %d log entries, want 2", len(logs))
	}

	tests := []struct {
		idx     int
		level   string
		message string
	}{
		{0, "info", "hello from server"},
		{1, "error", "something broke"},
	}
	for _, tt := range tests {
		l := logs[tt.idx]
		if l.Level != tt.level {
			t.Errorf("log %d: level = %q, want %q", tt.idx, l.Level, tt.level)
		}
		if l.Message != tt.message {
			t.Errorf("log %d: message = %q, want %q", tt.idx, l.Message, tt.message)
		}
	}
}

func TestStoreBodyCopy(t *testing.T) {
	s := NewStore(nil)

	raw := []byte(`{"jsonrpc":"2.0","method":"test","params":{}}`)
	s.Add("client→server", raw)

	// Mutate original — store should be unaffected.
	raw[0] = 'X'

	e := s.Entries(0, 1)[0]
	var msg struct {
		Method string `json:"method"`
	}
	if err := json.Unmarshal(e.Body, &msg); err != nil {
		t.Fatal(err)
	}
	if msg.Method != "test" {
		t.Errorf("method = %q, want test (body was mutated)", msg.Method)
	}
}
