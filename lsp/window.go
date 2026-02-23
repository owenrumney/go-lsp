package lsp

// MessageType represents the type of a show message notification.
type MessageType int

const (
	MessageTypeError   MessageType = 1
	MessageTypeWarning MessageType = 2
	MessageTypeInfo    MessageType = 3
	MessageTypeLog     MessageType = 4
)

// ShowMessageParams contains the params for window/showMessage.
type ShowMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

// ShowMessageRequestParams contains the params for window/showMessageRequest.
type ShowMessageRequestParams struct {
	Type    MessageType          `json:"type"`
	Message string               `json:"message"`
	Actions []MessageActionItem  `json:"actions,omitempty"`
}

// MessageActionItem represents a message action item.
type MessageActionItem struct {
	Title string `json:"title"`
}

// LogMessageParams contains the params for window/logMessage.
type LogMessageParams struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
}

// ShowDocumentParams contains the params for window/showDocument.
type ShowDocumentParams struct {
	URI       URI    `json:"uri"`
	External  *bool  `json:"external,omitempty"`
	TakeFocus *bool  `json:"takeFocus,omitempty"`
	Selection *Range `json:"selection,omitempty"`
}

// ShowDocumentResult is the result of a show document request.
type ShowDocumentResult struct {
	Success bool `json:"success"`
}

// WorkDoneProgressOptions represents options for work done progress.
type WorkDoneProgressOptions struct {
	WorkDoneProgress *bool `json:"workDoneProgress,omitempty"`
}
