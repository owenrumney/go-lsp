package lsp

import "encoding/json"

// DocumentDiagnosticReportKind defines the kind of a document diagnostic report.
type DocumentDiagnosticReportKind string

const (
	DiagnosticReportFull      DocumentDiagnosticReportKind = "full"
	DiagnosticReportUnchanged DocumentDiagnosticReportKind = "unchanged"
)

// DocumentDiagnosticParams contains the params for textDocument/diagnostic.
type DocumentDiagnosticParams struct {
	WorkDoneProgressParams
	PartialResultParams
	TextDocument     TextDocumentIdentifier `json:"textDocument"`
	Identifier       string                 `json:"identifier,omitempty"`
	PreviousResultID *string                `json:"previousResultId,omitempty"`
}

// FullDocumentDiagnosticReport represents a full diagnostic report.
type FullDocumentDiagnosticReport struct {
	Kind     string       `json:"kind"` // always "full"
	ResultID *string      `json:"resultId,omitempty"`
	Items    []Diagnostic `json:"items"`
}

// UnchangedDocumentDiagnosticReport represents an unchanged diagnostic report.
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

// WorkspaceDiagnosticParams contains the params for workspace/diagnostic.
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

// DiagnosticOptions describes diagnostic options.
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

// DiagnosticClientCapabilities defines capabilities for pull diagnostics.
type DiagnosticClientCapabilities struct {
	DynamicRegistration    *bool `json:"dynamicRegistration,omitempty"`
	RelatedDocumentSupport *bool `json:"relatedDocumentSupport,omitempty"`
}

// DiagnosticWorkspaceClientCapabilities defines workspace capabilities for diagnostics.
type DiagnosticWorkspaceClientCapabilities struct {
	RefreshSupport *bool `json:"refreshSupport,omitempty"`
}
