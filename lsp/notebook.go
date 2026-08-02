package lsp

import "encoding/json"

// NotebookDocument represents a notebook document.
//
// Since 3.17.0
type NotebookDocument struct {
	// The notebook document's URI.
	URI DocumentURI `json:"uri"`
	// The type of the notebook.
	NotebookType string `json:"notebookType"`
	// The version number of this document (it will increase after each
	// change, including undo/redo).
	Version int `json:"version"`
	// Additional metadata stored with the notebook document.
	Metadata json.RawMessage `json:"metadata,omitempty"`
	// The cells of a notebook.
	Cells []NotebookCell `json:"cells"`
}

// NotebookCellKind describes the kind of a notebook cell.
//
// Since 3.17.0
type NotebookCellKind int

const (
	// NotebookCellKindMarkup is a markup-cell, formatted source that is used
	// for display purposes.
	NotebookCellKindMarkup NotebookCellKind = 1
	// NotebookCellKindCode is a code cell.
	NotebookCellKindCode NotebookCellKind = 2
)

// NotebookCell holds the cell's kind, its backing text document URI, and
// metadata. The actual cell content is managed as a regular text document.
//
// Since 3.17.0
type NotebookCell struct {
	// The cell's kind.
	Kind NotebookCellKind `json:"kind"`
	// The URI of the cell's text document content.
	Document DocumentURI `json:"document"`
	// Additional metadata stored with the cell.
	Metadata json.RawMessage `json:"metadata,omitempty"`
	// Additional execution summary information if supported by the client.
	ExecutionSummary *ExecutionSummary `json:"executionSummary,omitempty"`
}

// ExecutionSummary describes the most recent execution of a notebook cell.
//
// Since 3.17.0
type ExecutionSummary struct {
	// A strictly monotonically increasing value indicating the execution
	// order of a cell inside a notebook.
	ExecutionOrder int `json:"executionOrder"`
	// Whether the execution was successful or not, if known by the client.
	Success *bool `json:"success,omitempty"`
}

// NotebookCellTextDocumentFilter is a filter that matches against a notebook
// cell's text document.
//
// Since 3.17.0
type NotebookCellTextDocumentFilter struct {
	// A filter that matches against the notebook containing the cell.
	Notebook NotebookDocumentFilter `json:"notebook"`
	// A language id like `python`. Will be matched against the language id of
	// the notebook cell document. '*' matches every language.
	Language string `json:"language,omitempty"`
}

// NotebookDocumentFilter is the spec union string | filter-object: a string
// is a notebook type. The object form matches on type, scheme, and/or pattern.
//
// Since 3.17.0
type NotebookDocumentFilter struct {
	// The type of the enclosing notebook.
	NotebookType string `json:"notebookType,omitempty"`
	// A URI scheme, like `file` or `untitled`.
	Scheme string `json:"scheme,omitempty"`
	// A glob pattern.
	Pattern string `json:"pattern,omitempty"`
}

func (f NotebookDocumentFilter) MarshalJSON() ([]byte, error) {
	if f.Scheme == "" && f.Pattern == "" && f.NotebookType != "" {
		return json.Marshal(f.NotebookType)
	}
	type alias NotebookDocumentFilter
	return json.Marshal(alias(f))
}

func (f *NotebookDocumentFilter) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		*f = NotebookDocumentFilter{}
		return json.Unmarshal(data, &f.NotebookType)
	}
	type alias NotebookDocumentFilter
	return json.Unmarshal(data, (*alias)(f))
}

// NotebookDocumentSyncOptions declares which notebooks (and which of their
// cells) the server wants to be synced.
//
// Since 3.17.0
type NotebookDocumentSyncOptions struct {
	// The notebooks to be synced.
	NotebookSelector []NotebookSelector `json:"notebookSelector"`
	// Whether save notifications should be forwarded to the server. Will only
	// be honored if mode === `notebook`.
	Save *bool `json:"save,omitempty"`
}

// NotebookSelector selects the notebook to be synced. Either Notebook or at
// least one entry in Cells must be provided.
//
// Since 3.17.0
type NotebookSelector struct {
	// The notebook to be synced.
	Notebook *NotebookDocumentFilter `json:"notebook,omitempty"`
	// The cells of the matching notebook to be synced.
	Cells []NotebookCellLanguage `json:"cells,omitempty"`
}

// NotebookCellLanguage selects notebook cells by language.
//
// Since 3.17.0
type NotebookCellLanguage struct {
	Language string `json:"language"`
}

// DidOpenNotebookDocumentParams holds the params sent in an open notebook
// document notification.
//
// Since 3.17.0
type DidOpenNotebookDocumentParams struct {
	// The notebook document that got opened.
	NotebookDocument NotebookDocument `json:"notebookDocument"`
	// The text documents that represent the content of a notebook cell.
	CellTextDocuments []TextDocumentItem `json:"cellTextDocuments"`
}

// VersionedNotebookDocumentIdentifier identifies a specific version of a
// notebook document.
//
// Since 3.17.0
type VersionedNotebookDocumentIdentifier struct {
	// The version number of this notebook document.
	Version int `json:"version"`
	// The notebook document's URI.
	URI DocumentURI `json:"uri"`
}

// NotebookDocumentIdentifier identifies a notebook document.
//
// Since 3.17.0
type NotebookDocumentIdentifier struct {
	// The notebook document's URI.
	URI DocumentURI `json:"uri"`
}

// DidChangeNotebookDocumentParams holds the params sent in a change notebook
// document notification.
//
// Since 3.17.0
type DidChangeNotebookDocumentParams struct {
	// The notebook document that did change.
	NotebookDocument VersionedNotebookDocumentIdentifier `json:"notebookDocument"`
	// The actual changes to the notebook document.
	Change NotebookDocumentChangeEvent `json:"change"`
}

// NotebookDocumentChangeEvent describes a change to a notebook document.
//
// Since 3.17.0
type NotebookDocumentChangeEvent struct {
	// The changed metadata if any.
	Metadata json.RawMessage `json:"metadata,omitempty"`
	// Changes to cells.
	Cells *NotebookDocumentCellChanges `json:"cells,omitempty"`
}

// NotebookDocumentCellChanges describes changes to the cells of a notebook.
//
// Since 3.17.0
type NotebookDocumentCellChanges struct {
	// Changes to the cell structure to add or remove cells.
	Structure *NotebookCellArrayChange `json:"structure,omitempty"`
	// Changes to notebook cells properties like its kind, execution summary
	// or metadata.
	Data []NotebookCell `json:"data,omitempty"`
	// Changes to the text content of notebook cells.
	TextContent []NotebookCellTextChange `json:"textContent,omitempty"`
}

// NotebookCellArrayChange describes a structural change to cells in a
// notebook document.
//
// Since 3.17.0
type NotebookCellArrayChange struct {
	// The change to the cell array.
	Array NotebookCellArrayDelta `json:"array"`
	// Additional opened cell text documents.
	DidOpen []TextDocumentItem `json:"didOpen,omitempty"`
	// Additional closed cell text documents.
	DidClose []TextDocumentIdentifier `json:"didClose,omitempty"`
}

// NotebookCellArrayDelta describes a splice-style change to a cell array.
//
// Since 3.17.0
type NotebookCellArrayDelta struct {
	// The start offset of the cell that changed.
	Start int `json:"start"`
	// The number of deleted cells.
	DeleteCount int `json:"deleteCount"`
	// The new cells, if any.
	Cells []NotebookCell `json:"cells,omitempty"`
}

// NotebookCellTextChange holds text content changes for a single cell
// document.
//
// Since 3.17.0
type NotebookCellTextChange struct {
	Document VersionedTextDocumentIdentifier  `json:"document"`
	Changes  []TextDocumentContentChangeEvent `json:"changes"`
}

// DidSaveNotebookDocumentParams holds the params sent in a save notebook
// document notification.
//
// Since 3.17.0
type DidSaveNotebookDocumentParams struct {
	// The notebook document that got saved.
	NotebookDocument NotebookDocumentIdentifier `json:"notebookDocument"`
}

// DidCloseNotebookDocumentParams holds the params sent in a close notebook
// document notification.
//
// Since 3.17.0
type DidCloseNotebookDocumentParams struct {
	// The notebook document that got closed.
	NotebookDocument NotebookDocumentIdentifier `json:"notebookDocument"`
	// The text documents that represent the content of a notebook cell that
	// got closed.
	CellTextDocuments []TextDocumentIdentifier `json:"cellTextDocuments"`
}

// NotebookDocumentClientCapabilities declares client capabilities specific to
// notebook documents.
//
// Since 3.17.0
type NotebookDocumentClientCapabilities struct {
	// Capabilities specific to notebook document synchronization.
	Synchronization *NotebookDocumentSyncClientCapabilities `json:"synchronization,omitempty"`
}

// NotebookDocumentSyncClientCapabilities declares client capabilities for
// notebook document synchronization.
//
// Since 3.17.0
type NotebookDocumentSyncClientCapabilities struct {
	// Whether implementation supports dynamic registration.
	DynamicRegistration *bool `json:"dynamicRegistration,omitempty"`
	// The client supports sending execution summary data per cell.
	ExecutionSummarySupport *bool `json:"executionSummarySupport,omitempty"`
}
