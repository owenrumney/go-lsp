package debugui

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
	"time"
)

const TraceVersion = 1

// TraceExportOptions controls how a debug trace is exported.
type TraceExportOptions struct {
	RedactDocumentText bool
	RedactFilePaths    bool
	RedactLogs         bool
	Pretty             bool
}

// Trace is a portable snapshot of a debug UI session.
type Trace struct {
	Version      int             `json:"version"`
	CreatedAt    time.Time       `json:"createdAt"`
	Messages     []Entry         `json:"messages"`
	Logs         []LogEntry      `json:"logs,omitempty"`
	Capabilities json.RawMessage `json:"capabilities,omitempty"`
}

// ExportTrace returns a JSON trace of the current debug UI state.
func (d *DebugUI) ExportTrace(opts TraceExportOptions) ([]byte, error) {
	trace := Trace{
		Version:      TraceVersion,
		CreatedAt:    time.Now().UTC(),
		Messages:     d.store.All(),
		Logs:         d.logStore.All(),
		Capabilities: d.capabilitiesSnapshot(),
	}

	redactTrace(&trace, opts)

	if opts.Pretty {
		return json.MarshalIndent(trace, "", "  ")
	}
	return json.Marshal(trace)
}

func (d *DebugUI) capabilitiesSnapshot() json.RawMessage {
	d.capsMu.RLock()
	defer d.capsMu.RUnlock()
	if d.capabilities == nil {
		return nil
	}
	return append(json.RawMessage(nil), d.capabilities...)
}

func redactTrace(trace *Trace, opts TraceExportOptions) {
	if opts.RedactLogs {
		trace.Logs = nil
	}

	for i := range trace.Messages {
		trace.Messages[i].Body = redactRawMessage(trace.Messages[i].Body, opts)
	}

	if opts.RedactFilePaths {
		for i := range trace.Logs {
			trace.Logs[i].Message = redactFilePaths(trace.Logs[i].Message)
		}
		trace.Capabilities = redactRawMessage(trace.Capabilities, opts)
	}
}

func redactRawMessage(raw json.RawMessage, opts TraceExportOptions) json.RawMessage {
	if len(raw) == 0 {
		return raw
	}

	var v any
	if err := json.Unmarshal(raw, &v); err != nil {
		if opts.RedactFilePaths {
			return json.RawMessage(strconvQuote(redactFilePaths(string(raw))))
		}
		return append(json.RawMessage(nil), raw...)
	}

	redactValue(v, opts, "")

	data, err := json.Marshal(v)
	if err != nil {
		return append(json.RawMessage(nil), raw...)
	}
	return data
}

func redactValue(v any, opts TraceExportOptions, key string) {
	switch x := v.(type) {
	case map[string]any:
		for k, child := range x {
			if opts.RedactDocumentText && isDocumentTextKey(k) {
				x[k] = "[redacted]"
				continue
			}
			redactValue(child, opts, k)
		}
		if opts.RedactFilePaths {
			for k, child := range x {
				redactedKey := redactFilePaths(k)
				if redactedKey != k {
					delete(x, k)
					x[redactedKey] = child
				}
			}
		}
	case []any:
		for _, child := range x {
			redactValue(child, opts, key)
		}
	case string:
		_ = key
	}

	if !opts.RedactFilePaths {
		return
	}

	switch x := v.(type) {
	case map[string]any:
		for k, child := range x {
			if s, ok := child.(string); ok {
				x[k] = redactFilePaths(s)
			}
		}
	case []any:
		for i, child := range x {
			if s, ok := child.(string); ok {
				x[i] = redactFilePaths(s)
			}
		}
	}
}

func isDocumentTextKey(key string) bool {
	switch key {
	case "text", "newText":
		return true
	default:
		return false
	}
}

var (
	fileURIRe = regexp.MustCompile(`file://[^\s"',)}\]]+`)
	pathRe    = regexp.MustCompile(`(?:/[A-Za-z0-9._@-]+){2,}`)
)

func redactFilePaths(s string) string {
	s = fileURIRe.ReplaceAllString(s, "file://[redacted]")
	return pathRe.ReplaceAllStringFunc(s, func(match string) string {
		if strings.HasPrefix(match, "/api/") {
			return match
		}
		return "/[redacted]"
	})
}

func strconvQuote(s string) string {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	_ = enc.Encode(s)
	return strings.TrimSpace(b.String())
}
