package lsp

// MessageType is an int enum: error (1), warning (2), info (3), or log (4).
type MessageType int

const (
	MessageTypeError   MessageType = 1
	MessageTypeWarning MessageType = 2
	MessageTypeInfo    MessageType = 3
	MessageTypeLog     MessageType = 4
)

// ShowMessageParams is sent from server to client to display a notification message in the editor.
type ShowMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

// ShowMessageRequestParams is sent from server to client to show a message with clickable action buttons.
type ShowMessageRequestParams struct {
	Type    MessageType         `json:"type"`
	Message string              `json:"message"`
	Actions []MessageActionItem `json:"actions,omitempty"`
}

// MessageActionItem is a button label offered in a showMessageRequest that the user can click.
type MessageActionItem struct {
	Title string `json:"title"`
}

// LogMessageParams is sent from server to client to log a message to the editor's output channel.
type LogMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

// ShowDocumentParams is sent from server to client to open a URI (file, URL) in the editor or external browser.
type ShowDocumentParams struct {
	URI       URI    `json:"uri"`
	External  *bool  `json:"external,omitempty"`
	TakeFocus *bool  `json:"takeFocus,omitempty"`
	Selection *Range `json:"selection,omitempty"`
}

// ShowDocumentResult indicates whether the editor successfully showed the requested document.
type ShowDocumentResult struct {
	Success bool `json:"success"`
}

// WorkDoneProgressOptions represents options for work done progress.
type WorkDoneProgressOptions struct {
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
}
