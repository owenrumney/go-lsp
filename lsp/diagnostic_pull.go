package lsp

import "encoding/json"

// DocumentDiagnosticReportKind is a string enum ("full" or "unchanged") indicating whether a diagnostic response contains new results or is unchanged since the last request.
type DocumentDiagnosticReportKind string

const (
	DiagnosticReportFull      DocumentDiagnosticReportKind = "full"
	DiagnosticReportUnchanged DocumentDiagnosticReportKind = "unchanged"
)

// DocumentDiagnosticParams is sent to pull diagnostics for a document (as opposed to the server pushing them).
type DocumentDiagnosticParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument     TextDocumentIdentifier `json:"textDocument"`
	Identifier       string                 `json:"identifier,omitempty"`
	PreviousResultID *string                `json:"previousResultId,omitempty"`
}

// FullDocumentDiagnosticReport contains the complete set of current diagnostics for a document.
type FullDocumentDiagnosticReport struct {
	Kind     string       `json:"kind"` // always "full"
	ResultID *string      `json:"resultId,omitempty"`
	Items    []Diagnostic `json:"items"`
}

// UnchangedDocumentDiagnosticReport indicates diagnostics have not changed since the previous request, identified by a result ID.
type UnchangedDocumentDiagnosticReport struct {
	Kind     string `json:"kind"` // always "unchanged"
	ResultID string `json:"resultId"`
}

// RelatedFullDocumentDiagnosticReport extends FullDocumentDiagnosticReport with related documents.
type RelatedFullDocumentDiagnosticReport struct {
	FullDocumentDiagnosticReport
	RelatedDocuments map[DocumentURI]json.RawMessage `json:"relatedDocuments,omitempty"`
}

// RelatedUnchangedDocumentDiagnosticReport extends UnchangedDocumentDiagnosticReport with related documents.
type RelatedUnchangedDocumentDiagnosticReport struct {
	UnchangedDocumentDiagnosticReport
	RelatedDocuments map[DocumentURI]json.RawMessage `json:"relatedDocuments,omitempty"`
}

// PreviousResultID pairs a document URI with its last known result ID.
type PreviousResultID struct {
	URI   DocumentURI `json:"uri"`
	Value string      `json:"value"`
}

// WorkspaceDiagnosticParams is sent to pull diagnostics across the entire workspace.
type WorkspaceDiagnosticParams struct {
	WorkDoneProgressParams
	PartialResultParams
	Identifier        string             `json:"identifier,omitempty"`
	PreviousResultIDs []PreviousResultID `json:"previousResultIds"`
}

// WorkspaceFullDocumentDiagnosticReport extends FullDocumentDiagnosticReport with workspace info.
type WorkspaceFullDocumentDiagnosticReport struct {
	FullDocumentDiagnosticReport
	URI     DocumentURI `json:"uri"`
	Version *int        `json:"version"`
}

// WorkspaceUnchangedDocumentDiagnosticReport extends UnchangedDocumentDiagnosticReport with workspace info.
type WorkspaceUnchangedDocumentDiagnosticReport struct {
	UnchangedDocumentDiagnosticReport
	URI     DocumentURI `json:"uri"`
	Version *int        `json:"version"`
}

// WorkspaceDiagnosticReport contains the results of a workspace diagnostic request.
type WorkspaceDiagnosticReport struct {
	Items []json.RawMessage `json:"items"`
}

// DiagnosticOptions configures pull-based diagnostics: whether the server supports inter-file dependencies and workspace-wide diagnostics.
type DiagnosticOptions struct {
	WorkDoneProgressOptions
	Identifier            string `json:"identifier,omitempty"`
	InterFileDependencies bool   `json:"interFileDependencies"`
	WorkspaceDiagnostics  bool   `json:"workspaceDiagnostics"`
}

// DiagnosticServerCancellationData is returned when the server cancels a diagnostic request.
type DiagnosticServerCancellationData struct {
	RetriggerRequest bool `json:"retriggerRequest"`
}

// DiagnosticClientCapabilities declares editor support for pull-based diagnostics (textDocument/diagnostic).
type DiagnosticClientCapabilities struct {
	DynamicRegistration    *bool `json:"dynamicRegistration,omitempty"`
	RelatedDocumentSupport *bool `json:"relatedDocumentSupport,omitempty"`
}

// DiagnosticWorkspaceClientCapabilities declares whether the editor will refresh diagnostics when the server requests it.
type DiagnosticWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
