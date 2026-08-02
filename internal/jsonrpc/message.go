package jsonrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
)

const Version = "2.0"

// ErrParse marks a body that is not valid JSON; framing is intact so the
// connection can recover.
var ErrParse = errors.New("jsonrpc: parse error")

// ErrInvalidMessage marks valid JSON that is not a JSON-RPC message; recoverable.
var ErrInvalidMessage = errors.New("jsonrpc: invalid message")

// ID represents a JSON-RPC 2.0 request ID, which can be a string or integer.
type ID struct {
	value any // string or int64
}

func IntID(v int64) ID     { return ID{value: v} }
func StringID(v string) ID { return ID{value: v} }

func (id ID) IsZero() bool { return id.value == nil }

func (id ID) String() string {
	switch v := id.value.(type) {
	case string:
		return v
	case int64:
		return fmt.Sprintf("%d", v)
	default:
		return ""
	}
}

func (id ID) MarshalJSON() ([]byte, error) {
	if id.value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(id.value)
}

func (id *ID) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		id.value = nil
		return nil
	}

	var n int64
	if err := json.Unmarshal(data, &n); err == nil {
		id.value = n
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		id.value = s
		return nil
	}

	return fmt.Errorf("jsonrpc: ID must be a string or number, got %s", string(data))
}

// Request is a JSON-RPC 2.0 request.
type Request struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      ID              `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// Response is a JSON-RPC 2.0 response.
type Response struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      ID              `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *ResponseError  `json:"error,omitempty"`
}

// Notification is a JSON-RPC 2.0 notification (no ID).
type Notification struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// NewRequest creates a request with the given method and params.
func NewRequest(id ID, method string, params any) (*Request, error) {
	var raw json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		raw = b
	}
	return &Request{JSONRPC: Version, ID: id, Method: method, Params: raw}, nil
}

// NewResponse creates a successful response; a nil result encodes as an
// explicit null since the spec requires the result member.
func NewResponse(id ID, result any) (*Response, error) {
	raw := json.RawMessage("null")
	if result != nil {
		b, err := json.Marshal(result)
		if err != nil {
			return nil, err
		}
		raw = b
	}
	return &Response{JSONRPC: Version, ID: id, Result: raw}, nil
}

// NewErrorResponse creates an error response.
func NewErrorResponse(id ID, respErr *ResponseError) *Response {
	return &Response{JSONRPC: Version, ID: id, Error: respErr}
}

// NewNotification creates a notification.
func NewNotification(method string, params any) (*Notification, error) {
	var raw json.RawMessage
	if params != nil {
		b, err := json.Marshal(params)
		if err != nil {
			return nil, err
		}
		raw = b
	}
	return &Notification{JSONRPC: Version, Method: method, Params: raw}, nil
}

// rawMessage is used for initial JSON parsing to determine the message type.
type rawMessage struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  *string          `json:"method,omitempty"`
	Params  json.RawMessage  `json:"params,omitempty"`
	Result  json.RawMessage  `json:"result,omitempty"`
	Error   *ResponseError   `json:"error,omitempty"`
}

// DecodeMessage decodes a JSON-RPC message into a Request, Response, or Notification.
func DecodeMessage(data []byte) (any, error) {
	if trimmed := bytes.TrimLeft(data, " \t\r\n"); len(trimmed) > 0 && trimmed[0] == '[' {
		return nil, fmt.Errorf("%w: batch messages are not supported", ErrInvalidMessage)
	}

	var raw rawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrParse, err)
	}

	hasID := raw.ID != nil && string(*raw.ID) != "null"

	// A null id decodes raw.ID to nil, indistinguishable from absent; probe
	// with a non-pointer RawMessage, which preserves the literal null.
	if raw.Method != nil && raw.ID == nil {
		var probe struct {
			ID json.RawMessage `json:"id"`
		}
		if err := json.Unmarshal(data, &probe); err == nil && string(probe.ID) == "null" {
			return nil, fmt.Errorf("%w: request id must not be null", ErrInvalidMessage)
		}
	}

	// Response: no method. The ID may be null for ParseError responses.
	if raw.Method == nil && (hasID || raw.Result != nil || raw.Error != nil) {
		var id ID
		if hasID {
			if err := json.Unmarshal(*raw.ID, &id); err != nil {
				return nil, fmt.Errorf("%w: %v", ErrInvalidMessage, err)
			}
		}
		return &Response{
			JSONRPC: raw.JSONRPC,
			ID:      id,
			Result:  raw.Result,
			Error:   raw.Error,
		}, nil
	}

	// Request: has ID and method
	if hasID && raw.Method != nil {
		var id ID
		if err := json.Unmarshal(*raw.ID, &id); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrInvalidMessage, err)
		}
		return &Request{
			JSONRPC: raw.JSONRPC,
			ID:      id,
			Method:  *raw.Method,
			Params:  raw.Params,
		}, nil
	}

	// Notification: has method but no ID
	if raw.Method != nil {
		return &Notification{
			JSONRPC: raw.JSONRPC,
			Method:  *raw.Method,
			Params:  raw.Params,
		}, nil
	}

	return nil, fmt.Errorf("%w: cannot determine message type", ErrInvalidMessage)
}
