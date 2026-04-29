package debugui

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestExportTrace(t *testing.T) {
	logStore := NewLogStore()
	store := NewStore(logStore)
	ui := New(":0", store, logStore)
	ui.SetCapabilities(map[string]any{"hoverProvider": true})

	store.Add("client→server", []byte(`{"jsonrpc":"2.0","id":1,"method":"textDocument/didOpen","params":{"textDocument":{"uri":"file:///Users/owen/project/main.go","text":"secret source"}}}`))
	logStore.Add("info", "opened /Users/owen/project/main.go")

	data, err := ui.ExportTrace(TraceExportOptions{Pretty: true})
	if err != nil {
		t.Fatal(err)
	}

	var trace Trace
	if err := json.Unmarshal(data, &trace); err != nil {
		t.Fatal(err)
	}
	if trace.Version != TraceVersion {
		t.Fatalf("version = %d, want %d", trace.Version, TraceVersion)
	}
	if len(trace.Messages) != 1 {
		t.Fatalf("messages = %d, want 1", len(trace.Messages))
	}
	if len(trace.Logs) != 1 {
		t.Fatalf("logs = %d, want 1", len(trace.Logs))
	}
	if len(trace.Capabilities) == 0 {
		t.Fatal("expected capabilities")
	}
}

func TestExportTraceRedaction(t *testing.T) {
	logStore := NewLogStore()
	store := NewStore(logStore)
	ui := New(":0", store, logStore)
	ui.SetCapabilities(map[string]any{"root": "file:///Users/owen/project"})

	store.Add("client→server", []byte(`{"jsonrpc":"2.0","id":1,"method":"textDocument/didChange","params":{"textDocument":{"uri":"file:///Users/owen/project/main.go"},"contentChanges":[{"text":"secret source","range":{"start":{"line":0,"character":0},"end":{"line":0,"character":4}}}]}}`))
	store.Add("server→client", []byte(`{"jsonrpc":"2.0","method":"workspace/applyEdit","params":{"edit":{"changes":{"file:///Users/owen/project/main.go":[{"newText":"secret edit"}]}}}}`))
	logStore.Add("info", "opened /Users/owen/project/main.go")

	data, err := ui.ExportTrace(TraceExportOptions{
		RedactDocumentText: true,
		RedactFilePaths:    true,
		RedactLogs:         true,
	})
	if err != nil {
		t.Fatal(err)
	}

	out := string(data)
	for _, secret := range []string{"secret source", "secret edit", "/Users/owen/project", "file:///Users/owen"} {
		if strings.Contains(out, secret) {
			t.Fatalf("trace contains unredacted value %q: %s", secret, out)
		}
	}
	if !strings.Contains(out, "[redacted]") {
		t.Fatalf("trace does not contain redaction marker: %s", out)
	}

	var trace Trace
	if err := json.Unmarshal(data, &trace); err != nil {
		t.Fatal(err)
	}
	if trace.Logs != nil {
		t.Fatalf("logs = %#v, want nil", trace.Logs)
	}
}
