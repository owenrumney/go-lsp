package lsp

// RenameParams is sent to rename a symbol across the workspace.
type RenameParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	NewName string `json:"newName"`
}

// PrepareRenameParams is sent to validate that a rename is possible at a position before showing the rename UI.
type PrepareRenameParams struct {
	TextDocumentPositionParams
}

// PrepareRenameResult returns the range and placeholder text for the symbol to be renamed.
type PrepareRenameResult struct {
	Range       Range  `json:"range"`
	Placeholder string `json:"placeholder"`
}
