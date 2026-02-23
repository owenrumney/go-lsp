package lsp

import "encoding/json"

// ProgressToken is a token used to report progress (string or int).
type ProgressToken = json.RawMessage

// WorkDoneProgressParams is a mixin for work done progress.
type WorkDoneProgressParams struct {
	WorkDoneToken *ProgressToken `json:"workDoneToken,omitempty"`
}

// PartialResultParams is a mixin for partial result progress.
type PartialResultParams struct {
	PartialResultToken *ProgressToken `json:"partialResultToken,omitempty"`
}

// WorkDoneProgressBegin is sent to start a work done progress.
type WorkDoneProgressBegin struct {
	Kind        string `json:"kind"` // "begin"
	Title       string `json:"title"`
	Cancellable *bool  `json:"cancellable,omitempty"`
	Message     string `json:"message,omitempty"`
	Percentage  *int   `json:"percentage,omitempty"`
}

// WorkDoneProgressReport is sent to report work done progress.
type WorkDoneProgressReport struct {
	Kind        string `json:"kind"` // "report"
	Cancellable *bool  `json:"cancellable,omitempty"`
	Message     string `json:"message,omitempty"`
	Percentage  *int   `json:"percentage,omitempty"`
}

// WorkDoneProgressEnd is sent to end a work done progress.
type WorkDoneProgressEnd struct {
	Kind    string `json:"kind"` // "end"
	Message string `json:"message,omitempty"`
}

// ProgressParams contains the params for $/progress.
type ProgressParams struct {
	Token ProgressToken   `json:"token"`
	Value json.RawMessage `json:"value"`
}

// WorkDoneProgressCreateParams contains the params for window/workDoneProgress/create.
type WorkDoneProgressCreateParams struct {
	Token ProgressToken `json:"token"`
}
