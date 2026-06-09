package lsp

import "encoding/json"

// DiagnosticSeverity is an int enum: Error (1), Warning (2), Information (3), Hint (4).
type DiagnosticSeverity int

const (
	// Reports an error.
	SeverityError DiagnosticSeverity = 1
	// Reports a warning.
	SeverityWarning DiagnosticSeverity = 2
	// Reports information.
	SeverityInformation DiagnosticSeverity = 3
	// Reports a hint.
	SeverityHint DiagnosticSeverity = 4
)

// DiagnosticTag represents additional metadata about a diagnostic.
type DiagnosticTag int

const (
	// Unused or unnecessary code.
	//
	// Clients are allowed to render diagnostics with this tag faded out instead of having
	// an error squiggle.
	TagUnnecessary DiagnosticTag = 1
	// Deprecated or obsolete code.
	//
	// Clients are allowed to render diagnostics with this tag struck through.
	TagDeprecated DiagnosticTag = 2
)

// DiagnosticRelatedInformation represents a related message and source code location for a diagnostic. This should be
// used to point to code locations that cause or related to a diagnostic, e.g when duplicating
// a symbol in a scope.
type DiagnosticRelatedInformation struct {
	// The location of this related diagnostic information.
	Location Location `json:"location"`
	// The message of this related diagnostic information.
	Message string `json:"message"`
}

// CodeDescription is the structure to capture a description for an error code.
//
// Since 3.16.0.
type CodeDescription struct {
	// An URI to open with more information about the diagnostic error.
	Href URI `json:"href"`
}

// Diagnostic represents a diagnostic, such as a compiler error or warning. Diagnostic objects
// are only valid in the scope of a resource.
type Diagnostic struct {
	// The range at which the message applies
	Range Range `json:"range"`
	// The diagnostic's severity. Can be omitted. If omitted it is up to the
	// client to interpret diagnostics as error, warning, info or hint.
	Severity *DiagnosticSeverity `json:"severity,omitempty"`
	// The diagnostic's code, which usually appear in the user interface.
	Code json.RawMessage `json:"code,omitempty"` // int | string
	// An optional property to describe the error code.
	// Requires the code field (above) to be present/not null.
	//
	// Since 3.16.0
	CodeDescription *CodeDescription `json:"codeDescription,omitempty"`
	// A human-readable string describing the source of this
	// diagnostic, e.g. 'typescript' or 'super lint'. It usually
	// appears in the user interface.
	Source string `json:"source,omitempty"`
	// The diagnostic's message. It usually appears in the user interface
	Message string `json:"message"`
	// Additional metadata about the diagnostic.
	//
	// Since 3.15.0
	Tags []DiagnosticTag `json:"tags,omitempty"`
	// An array of related diagnostic information, e.g. when symbol-names within
	// a scope collide all definitions can be marked via this property.
	RelatedInformation []DiagnosticRelatedInformation `json:"relatedInformation,omitempty"`
	// A data entry field that is preserved between a `textDocument/publishDiagnostics`
	// notification and `textDocument/codeAction` request.
	//
	// Since 3.16.0
	Data json.RawMessage `json:"data,omitempty"`
}

// PublishDiagnosticsParams holds the publish diagnostic notification's parameters.
type PublishDiagnosticsParams struct {
	// The URI for which diagnostic information is reported.
	URI DocumentURI `json:"uri"`
	// Optional the version number of the document the diagnostics are published for.
	//
	// Since 3.15.0
	Version *int `json:"version,omitempty"`
	// An array of diagnostic information items.
	Diagnostics []Diagnostic `json:"diagnostics"`
}
