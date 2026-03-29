package lsp

import "encoding/json"

// Command is an editor command with a title and arguments, typically shown in the UI and executed via workspace/executeCommand.
type Command struct {
	Title     string            `json:"title"`
	Command   string            `json:"command"`
	Arguments []json.RawMessage `json:"arguments,omitempty"`
}
