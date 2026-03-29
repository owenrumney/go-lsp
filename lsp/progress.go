package lsp

import "encoding/json"

// ProgressToken is a token used to report progress (string or int).
type ProgressToken = json.RawMessage

// WorkDoneProgressParams embeds a progress token so the server can report progress for a long-running request.
type WorkDoneProgressParams struct {
	WorkDoneToken *ProgressToken `json:"workDoneToken,omitempty"`
}

// PartialResultParams embeds a progress token so the server can stream partial results for a request.
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

// ProgressParams carries a progress token and a value payload (begin, report, or end) for the $/progress notification.
type ProgressParams struct {
	Token ProgressToken   `json:"token"`
	Value json.RawMessage `json:"value"`
}

// WorkDoneProgressCreateParams is sent from server to client to create a new progress indicator in the UI.
type WorkDoneProgressCreateParams struct {
	Token ProgressToken `json:"token"`
}
