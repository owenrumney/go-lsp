package lsp

// MessageType is an int enum: error (1), warning (2), info (3), or log (4).
type MessageType int

const (
	// An error message.
	MessageTypeError MessageType = 1
	// A warning message.
	MessageTypeWarning MessageType = 2
	// An information message.
	MessageTypeInfo MessageType = 3
	// A log message.
	MessageTypeLog MessageType = 4
)

// ShowMessageParams holds the parameters of a notification message.
type ShowMessageParams struct {
	// The message type. See [MessageType]
	Type MessageType `json:"type"`
	// The actual message.
	Message string `json:"message"`
}

// ShowMessageRequestParams is sent from server to client to show a message with clickable action buttons.
type ShowMessageRequestParams struct {
	// The message type. See [MessageType]
	Type MessageType `json:"type"`
	// The actual message.
	Message string `json:"message"`
	// The message action items to present.
	Actions []MessageActionItem `json:"actions,omitempty"`
}

// MessageActionItem is a button label offered in a showMessageRequest that the user can click.
type MessageActionItem struct {
	// A short title like 'Retry', 'Open Log' etc.
	Title string `json:"title"`
}

// LogMessageParams holds the log message parameters.
type LogMessageParams struct {
	// The message type. See [MessageType]
	Type MessageType `json:"type"`
	// The actual message.
	Message string `json:"message"`
}

// ShowDocumentParams is used to show a resource in the UI.
//
// Since 3.16.0.
type ShowDocumentParams struct {
	// The uri to show.
	URI URI `json:"uri"`
	// Indicates to show the resource in an external program.
	// To show, for example, `https://code.visualstudio.com/`
	// in the default WEB browser set external to true.
	External *bool `json:"external,omitempty"`
	// An optional property to indicate whether the editor
	// showing the document should take focus or not.
	// Clients might ignore this property if an external
	// program is started.
	TakeFocus *bool `json:"takeFocus,omitempty"`
	// An optional selection range if the document is a text
	// document. Clients might ignore the property if an
	// external program is started or the file is not a text
	// file.
	Selection *Range `json:"selection,omitempty"`
}

// ShowDocumentResult is the result of a showDocument request.
//
// Since 3.16.0.
type ShowDocumentResult struct {
	// A boolean indicating if the show was successful.
	Success bool `json:"success"`
}

// WorkDoneProgressOptions represents options for work done progress.
type WorkDoneProgressOptions struct {
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
}
