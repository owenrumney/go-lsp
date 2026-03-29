package lsp

import "encoding/json"

// DiagnosticSeverity is an int enum: Error (1), Warning (2), Information (3), Hint (4).
type DiagnosticSeverity int

const (
	SeverityError       DiagnosticSeverity = 1
	SeverityWarning     DiagnosticSeverity = 2
	SeverityInformation DiagnosticSeverity = 3
	SeverityHint        DiagnosticSeverity = 4
)

// DiagnosticTag represents additional metadata about a diagnostic.
type DiagnosticTag int

const (
	TagUnnecessary DiagnosticTag = 1
	TagDeprecated  DiagnosticTag = 2
)

// DiagnosticRelatedInformation points to a related source location that helps explain a diagnostic (e.g. "variable declared here").
type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

// CodeDescription describes a code with a URI to open.
type CodeDescription struct {
	Href URI `json:"href"`
}

// Diagnostic is an error, warning, or informational message attached to a source range, shown in the editor's problems panel.
type Diagnostic struct {
	Range              Range                          `json:"range"`
	Severity           *DiagnosticSeverity            `json:"severity,omitempty"`
	Code               json.RawMessage                `json:"code,omitempty"` // int | string
	CodeDescription    *CodeDescription               `json:"codeDescription,omitempty"`
	Source             string                         `json:"source,omitempty"`
	Message            string                         `json:"message"`
	Tags               []DiagnosticTag                `json:"tags,omitempty"`
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
	Data               json.RawMessage                `json:"data,omitempty"`
}

// PublishDiagnosticsParams is sent from server to client to publish diagnostics (errors, warnings) for a document.
type PublishDiagnosticsParams struct {
	URI         DocumentURI  `json:"uri"`
	Version     *int         `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}
