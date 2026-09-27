package jsonrpc

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"
)

func TestCallTimeout(t *testing.T) {
	serverReader, clientWriter := io.Pipe()
	clientReader, serverWriter := io.Pipe()
	t.Cleanup(func() {
		_ = clientWriter.Close()
		_ = serverWriter.Close()
	})

	conn := NewConn(pipeRWC{Reader: serverReader, Writer: serverWriter}, NewDispatcher())
	conn.SetCallTimeout(50 * time.Millisecond)

	// Drain the request so the write does not block; never answer it.
	go func() { _, _ = io.Copy(io.Discard, clientReader) }()

	start := time.Now()
	_, err := conn.Call(context.Background(), "window/showMessageRequest", nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want DeadlineExceeded", err)
	}
	if time.Since(start) > time.Second {
		t.Fatalf("call took %s, expected to time out promptly", time.Since(start))
	}
}

type pipeRWC struct {
	io.Reader
	io.Writer
}

func (pipeRWC) Close() error { return nil }
