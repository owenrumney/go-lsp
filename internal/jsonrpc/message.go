package jsonrpc

import (
	"encoding/json"
	"fmt"
)

const Version = "2.0"

// ID represents a JSON-RPC 2.0 request ID, which can be a string or integer.
type ID struct {
	value any // string or int64
}

func IntID(v int64) ID    { return ID{value: v} }
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

// NewResponse creates a successful response.
func NewResponse(id ID, result any) (*Response, error) {
	var raw json.RawMessage
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
	var raw rawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("jsonrpc: failed to decode message: %w", err)
	}

	hasID := raw.ID != nil && string(*raw.ID) != "null"

	// Response: has ID but no method
	if hasID && raw.Method == nil {
		var id ID
		if err := json.Unmarshal(*raw.ID, &id); err != nil {
			return nil, err
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
			return nil, err
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

	return nil, fmt.Errorf("jsonrpc: cannot determine message type")
}
