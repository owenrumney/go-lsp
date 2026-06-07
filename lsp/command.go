package lsp

import "encoding/json"

// Command represents a reference to a command. Provides a title which
// will be used to represent a command in the UI and, optionally,
// an array of arguments which will be passed to the command handler
// function when invoked.
type Command struct {
	// Title of the command, like save.
	Title string `json:"title"`
	// The identifier of the actual command handler.
	Command string `json:"command"`
	// Arguments that the command handler should be
	// invoked with.
	Arguments []json.RawMessage `json:"arguments,omitempty"`
}
