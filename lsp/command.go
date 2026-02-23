package lsp

import "encoding/json"

// Command represents a reference to a command.
type Command struct {
	Title     string            `json:"title"`
	Command   string            `json:"command"`
	Arguments []json.RawMessage `json:"arguments,omitempty"`
}
