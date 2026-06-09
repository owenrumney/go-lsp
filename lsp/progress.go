package lsp

import "encoding/json"

// ProgressToken is a token used to report progress (string or int).
type ProgressToken = json.RawMessage

// WorkDoneProgressParams embeds a progress token so the server can report progress for a long-running request.
type WorkDoneProgressParams struct {
	// An optional token that a server can use to report work done progress.
	WorkDoneToken *ProgressToken `json:"workDoneToken,omitempty"`
}

// PartialResultParams embeds a progress token so the server can stream partial results for a request.
type PartialResultParams struct {
	// An optional token that a server can use to report partial results (e.g. streaming) to
	// the client.
	PartialResultToken *ProgressToken `json:"partialResultToken,omitempty"`
}

// WorkDoneProgressBegin is sent to start a work done progress.
type WorkDoneProgressBegin struct {
	Kind string `json:"kind"` // "begin"
	// Mandatory title of the progress operation. Used to briefly inform about
	// the kind of operation being performed.
	//
	// Examples: "Indexing" or "Linking dependencies".
	Title string `json:"title"`
	// Controls if a cancel button should show to allow the user to cancel the
	// long running operation. Clients that don't support cancellation are allowed
	// to ignore the setting.
	Cancellable *bool `json:"cancellable,omitempty"`
	// Optional, more detailed associated progress message. Contains
	// complementary information to the title.
	//
	// Examples: "3/25 files", "project/src/module2", "node_modules/some_dep".
	// If unset, the previous progress message (if any) is still valid.
	Message string `json:"message,omitempty"`
	// Optional progress percentage to display (value 100 is considered 100%).
	// If not provided infinite progress is assumed and clients are allowed
	// to ignore the percentage value in subsequent report notifications.
	//
	// The value should be steadily rising. Clients are free to ignore values
	// that are not following this rule. The value range is [0, 100].
	Percentage *int `json:"percentage,omitempty"`
}

// WorkDoneProgressReport is sent to report work done progress.
type WorkDoneProgressReport struct {
	Kind string `json:"kind"` // "report"
	// Controls enablement state of a cancel button.
	//
	// Clients that don't support cancellation or don't support controlling the button's
	// enablement state are allowed to ignore the property.
	Cancellable *bool `json:"cancellable,omitempty"`
	// Optional, more detailed associated progress message. Contains
	// complementary information to the title.
	//
	// Examples: "3/25 files", "project/src/module2", "node_modules/some_dep".
	// If unset, the previous progress message (if any) is still valid.
	Message string `json:"message,omitempty"`
	// Optional progress percentage to display (value 100 is considered 100%).
	// If not provided infinite progress is assumed and clients are allowed
	// to ignore the percentage value in subsequent report notifications.
	//
	// The value should be steadily rising. Clients are free to ignore values
	// that are not following this rule. The value range is [0, 100].
	Percentage *int `json:"percentage,omitempty"`
}

// WorkDoneProgressEnd is sent to end a work done progress.
type WorkDoneProgressEnd struct {
	Kind string `json:"kind"` // "end"
	// Optional, a final message indicating, for example, the outcome
	// of the operation.
	Message string `json:"message,omitempty"`
}

// ProgressParams carries a progress token and a value payload (begin, report, or end) for the $/progress notification.
type ProgressParams struct {
	// The progress token provided by the client or server.
	Token ProgressToken `json:"token"`
	// The progress data.
	Value json.RawMessage `json:"value"`
}

// WorkDoneProgressCreateParams is sent from server to client to create a new progress indicator in the UI.
type WorkDoneProgressCreateParams struct {
	// The token to be used to report progress.
	Token ProgressToken `json:"token"`
}
