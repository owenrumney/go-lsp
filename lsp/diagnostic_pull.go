package lsp

import "encoding/json"

// DocumentDiagnosticReportKind is a string enum ("full" or "unchanged") indicating whether a diagnostic response contains new results or is unchanged since the last request.
type DocumentDiagnosticReportKind string

const (
	// A diagnostic report with a full
	// set of problems.
	DiagnosticReportFull DocumentDiagnosticReportKind = "full"
	// A report indicating that the last
	// returned report is still accurate.
	DiagnosticReportUnchanged DocumentDiagnosticReportKind = "unchanged"
)

// DocumentDiagnosticParams holds the parameters of the document diagnostic request.
//
// Since 3.17.0.
type DocumentDiagnosticParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The text document.
	TextDocument TextDocumentIdentifier `json:"textDocument"`
	// The additional identifier  provided during registration.
	Identifier string `json:"identifier,omitempty"`
	// The result id of a previous response if provided.
	PreviousResultID *string `json:"previousResultId,omitempty"`
}

// FullDocumentDiagnosticReport is a diagnostic report with a full set of problems.
//
// Since 3.17.0.
type FullDocumentDiagnosticReport struct {
	// A full document diagnostic report.
	Kind string `json:"kind"` // always "full"
	// An optional result id. If provided it will
	// be sent on the next diagnostic request for the
	// same document.
	ResultID *string `json:"resultId,omitempty"`
	// The actual items.
	Items []Diagnostic `json:"items"`
}

// UnchangedDocumentDiagnosticReport is a diagnostic report indicating that the last returned
// report is still accurate.
//
// Since 3.17.0.
type UnchangedDocumentDiagnosticReport struct {
	// A document diagnostic report indicating
	// no changes to the last result. A server can
	// only return unchanged if result ids are
	// provided.
	Kind string `json:"kind"` // always "unchanged"
	// A result id which will be sent on the next
	// diagnostic request for the same document.
	ResultID string `json:"resultId"`
}

// RelatedFullDocumentDiagnosticReport is a full diagnostic report with a set of related documents.
//
// Since 3.17.0.
type RelatedFullDocumentDiagnosticReport struct {
	FullDocumentDiagnosticReport
	// Diagnostics of related documents. This information is useful
	// in programming languages where code in a file A can generate
	// diagnostics in a file B which A depends on. An example of
	// such a language is C/C++ where macro definitions in a file
	// a.cpp can result in errors in a header file b.hpp.
	//
	// Since 3.17.0
	RelatedDocuments map[DocumentURI]json.RawMessage `json:"relatedDocuments,omitempty"`
}

// RelatedUnchangedDocumentDiagnosticReport is an unchanged diagnostic report with a set of related documents.
//
// Since 3.17.0.
type RelatedUnchangedDocumentDiagnosticReport struct {
	UnchangedDocumentDiagnosticReport
	// Diagnostics of related documents. This information is useful
	// in programming languages where code in a file A can generate
	// diagnostics in a file B which A depends on. An example of
	// such a language is C/C++ where macro definitions in a file
	// a.cpp can result in errors in a header file b.hpp.
	//
	// Since 3.17.0
	RelatedDocuments map[DocumentURI]json.RawMessage `json:"relatedDocuments,omitempty"`
}

// PreviousResultID pairs a document URI with its last known result ID.
type PreviousResultID struct {
	URI   DocumentURI `json:"uri"`
	Value string      `json:"value"`
}

// WorkspaceDiagnosticParams holds the parameters of the workspace diagnostic request.
//
// Since 3.17.0.
type WorkspaceDiagnosticParams struct {
	WorkDoneProgressParams
	PartialResultParams
	// The additional identifier provided during registration.
	Identifier string `json:"identifier,omitempty"`
	// The currently known diagnostic reports with their
	// previous result ids.
	PreviousResultIDs []PreviousResultID `json:"previousResultIds"`
}

// WorkspaceFullDocumentDiagnosticReport is a full document diagnostic report for a workspace diagnostic result.
//
// Since 3.17.0.
type WorkspaceFullDocumentDiagnosticReport struct {
	FullDocumentDiagnosticReport
	// The URI for which diagnostic information is reported.
	URI DocumentURI `json:"uri"`
	// The version number for which the diagnostics are reported.
	// If the document is not marked as open null can be provided.
	Version *int `json:"version"`
}

// WorkspaceUnchangedDocumentDiagnosticReport is an unchanged document diagnostic report for a workspace diagnostic result.
//
// Since 3.17.0.
type WorkspaceUnchangedDocumentDiagnosticReport struct {
	UnchangedDocumentDiagnosticReport
	// The URI for which diagnostic information is reported.
	URI DocumentURI `json:"uri"`
	// The version number for which the diagnostics are reported.
	// If the document is not marked as open null can be provided.
	Version *int `json:"version"`
}

// WorkspaceDiagnosticReport is the result of a workspace/diagnostic request.
//
// Since 3.17.0.
type WorkspaceDiagnosticReport struct {
	Items []json.RawMessage `json:"items"`
}

// DiagnosticOptions configures the server's pull-diagnostic provider (textDocument/diagnostic and optionally workspace/diagnostic).
//
// Since 3.17.0.
type DiagnosticOptions struct {
	WorkDoneProgressOptions
	// An optional identifier under which the diagnostics are
	// managed by the client.
	Identifier string `json:"identifier,omitempty"`
	// Whether the language has inter file dependencies meaning that
	// editing code in one file can result in a different diagnostic
	// set in another file. Inter file dependencies are common for
	// most programming languages and typically uncommon for linters.
	InterFileDependencies bool `json:"interFileDependencies"`
	// The server provides support for workspace diagnostics as well.
	WorkspaceDiagnostics bool `json:"workspaceDiagnostics"`
}

// DiagnosticServerCancellationData is returned from a diagnostic request.
//
// Since 3.17.0.
type DiagnosticServerCancellationData struct {
	RetriggerRequest bool `json:"retriggerRequest"`
}

// DiagnosticClientCapabilities is specific to diagnostic pull requests.
//
// Since 3.17.0.
type DiagnosticClientCapabilities struct {
	// Whether implementation supports dynamic registration. If this is set to true
	// the client supports the new `(TextDocumentRegistrationOptions & StaticRegistrationOptions)`
	// return value for the corresponding server capability as well.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// Whether the client supports related documents for document diagnostic pulls.
	RelatedDocumentSupport *bool `json:"relatedDocumentSupport,omitempty"`
}

// DiagnosticWorkspaceClientCapabilities is specific to diagnostic pull requests.
//
// Since 3.17.0.
type DiagnosticWorkspaceClientCapabilities struct {
	// Whether the client implementation supports a refresh request sent from
	// the server to the client.
	//
	// Note that this event is global and will force the client to refresh all
	// pulled diagnostics currently shown. It should be used with absolute care and
	// is useful for situations where a server, for example, detects a project-wide
	// change that requires such a calculation.
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
