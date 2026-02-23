package lsp

// RenameParams contains the params for textDocument/rename.
type RenameParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	NewName string `json:"newName"`
}

// PrepareRenameParams contains the params for textDocument/prepareRename.
type PrepareRenameParams struct {
	TextDocumentPositionParams
}

// PrepareRenameResult represents the result of a prepare rename request.
type PrepareRenameResult struct {
	Range       Range  `json:"range"`
	Placeholder string `json:"placeholder"`
}
