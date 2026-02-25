package jsonrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"testing"
)

func TestID_MarshalUnmarshal(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want string
	}{
		{"int", IntID(42), "42"},
		{"string", StringID("abc"), `"abc"`},
		{"zero", ID{}, "null"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.id)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tt.want {
				t.Errorf("got %s, want %s", data, tt.want)
			}

			var got ID
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if got.String() != tt.id.String() {
				t.Errorf("round-trip: got %s, want %s", got.String(), tt.id.String())
			}
		})
	}
}

func TestDecodeMessage(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantType string
	}{
		{
			"request",
			`{"jsonrpc":"2.0","id":1,"method":"test","params":{"a":1}}`,
			"*jsonrpc.Request",
		},
		{
			"notification",
			`{"jsonrpc":"2.0","method":"notify","params":{}}`,
			"*jsonrpc.Notification",
		},
		{
			"response with result",
			`{"jsonrpc":"2.0","id":1,"result":{"b":2}}`,
			"*jsonrpc.Response",
		},
		{
			"response with error",
			`{"jsonrpc":"2.0","id":1,"error":{"code":-32601,"message":"not found"}}`,
			"*jsonrpc.Response",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := DecodeMessage([]byte(tt.input))
			if err != nil {
				t.Fatal(err)
			}
			got := fmt.Sprintf("%T", msg)
			if got != tt.wantType {
				t.Errorf("got type %s, want %s", got, tt.wantType)
			}
		})
	}
}

type nopCloser struct {
	io.Reader
	io.Writer
}

func (nopCloser) Close() error { return nil }

func TestConn_ReadWriteMessage(t *testing.T) {
	// Write a framed message to a buffer
	var buf bytes.Buffer
	conn := NewConn(nopCloser{Reader: &buf, Writer: &buf}, NewDispatcher())

	req := &Request{
		JSONRPC: Version,
		ID:      IntID(1),
		Method:  "test",
		Params:  json.RawMessage(`{"x":1}`),
	}

	if err := conn.WriteMessage(req); err != nil {
		t.Fatal(err)
	}

	// Now read it back using a new conn on the same buffer
	conn2 := NewConn(nopCloser{Reader: &buf, Writer: io.Discard}, NewDispatcher())
	msg, err := conn2.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}

	got, ok := msg.(*Request)
	if !ok {
		t.Fatalf("got %T, want *Request", msg)
	}
	if got.Method != "test" {
		t.Errorf("method = %s, want test", got.Method)
	}
	if got.ID.String() != "1" {
		t.Errorf("id = %s, want 1", got.ID.String())
	}
}

func TestDispatcher_HandleRequest(t *testing.T) {
	d := NewDispatcher()
	d.RegisterMethod("add", func(_ context.Context, params json.RawMessage) (any, error) {
		var p struct{ A, B int }
		if err := json.Unmarshal(params, &p); err != nil {
			return nil, err
		}
		return map[string]int{"sum": p.A + p.B}, nil
	})

	req := &Request{
		JSONRPC: Version,
		ID:      IntID(1),
		Method:  "add",
		Params:  json.RawMessage(`{"A":2,"B":3}`),
	}

	resp := d.HandleRequest(context.Background(), req)
	if resp.Error != nil {
		t.Fatal(resp.Error)
	}

	var result struct{ Sum int }
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		t.Fatal(err)
	}
	if result.Sum != 5 {
		t.Errorf("sum = %d, want 5", result.Sum)
	}
}

func TestDispatcher_MethodNotFound(t *testing.T) {
	d := NewDispatcher()
	req := &Request{JSONRPC: Version, ID: IntID(1), Method: "missing"}

	resp := d.HandleRequest(context.Background(), req)
	if resp.Error == nil {
		t.Fatal("expected error")
	}
	if resp.Error.Code != CodeMethodNotFound {
		t.Errorf("code = %d, want %d", resp.Error.Code, CodeMethodNotFound)
	}
}
