package lsp

import "encoding/json"

// DiagnosticSeverity indicates the severity of a diagnostic.
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

// DiagnosticRelatedInformation represents a related message and source code location.
type DiagnosticRelatedInformation struct {
	Location Location `json:"location"`
	Message  string   `json:"message"`
}

// CodeDescription describes a code with a URI to open.
type CodeDescription struct {
	Href URI `json:"href"`
}

// Diagnostic represents a diagnostic such as a compiler error or warning.
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

// PublishDiagnosticsParams contains the params for textDocument/publishDiagnostics.
type PublishDiagnosticsParams struct {
	URI         DocumentURI  `json:"uri"`
	Version     *int         `json:"version,omitempty"`
	Diagnostics []Diagnostic `json:"diagnostics"`
}
