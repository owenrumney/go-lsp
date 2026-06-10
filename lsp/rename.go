package lsp

// RenameParams holds the parameters of a [RenameRequest].
type RenameParams struct {
	TextDocumentPositionParams
	WorkDoneProgressParams
	// The new name of the symbol. If the given name is not valid the
	// request must return a [ResponseError] with an
	// appropriate message set.
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
