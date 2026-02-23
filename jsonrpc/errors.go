package jsonrpc

import (
	"encoding/json"
	"fmt"
)

// Standard JSON-RPC 2.0 error codes.
const (
	CodeParseError     = -32700
	CodeInvalidRequest = -32600
	CodeMethodNotFound = -32601
	CodeInvalidParams  = -32602
	CodeInternalError  = -32603

	// LSP-specific error codes.
	CodeServerNotInitialized = -32002
	CodeRequestCancelled     = -32800
	CodeContentModified      = -32801
	CodeServerCancelled      = -32802
	CodeRequestFailed        = -32803
)

// ResponseError represents a JSON-RPC 2.0 error object.
type ResponseError struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data,omitempty"`
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("jsonrpc error %d: %s", e.Code, e.Message)
}

func NewError(code int, message string) *ResponseError {
	return &ResponseError{Code: code, Message: message}
}
