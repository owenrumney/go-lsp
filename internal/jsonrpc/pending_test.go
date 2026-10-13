package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"sync/atomic"
	"testing"
	"time"
)

func TestCallSendsCancelRequestOnContextExpiry(t *testing.T) {
	serverReader, clientWriter := io.Pipe()
	clientReader, serverWriter := io.Pipe()
	t.Cleanup(func() {
		_ = clientWriter.Close()
		_ = serverWriter.Close()
	})

	conn := NewConn(pipeRWC{Reader: serverReader, Writer: serverWriter}, NewDispatcher())

	peer := NewConn(pipeRWC{Reader: clientReader, Writer: clientWriter}, NewDispatcher())
	seen := make(chan string, 2)
	go func() {
		for {
			msg, err := peer.ReadMessage()
			if err != nil {
				return
			}
			switch m := msg.(type) {
			case *Request:
				seen <- m.Method
			case *Notification:
				seen <- m.Method
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	_, err := conn.Call(ctx, "workspace/applyEdit", nil)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want DeadlineExceeded", err)
	}

	want := []string{"workspace/applyEdit", "$/cancelRequest"}
	for _, w := range want {
		select {
		case got := <-seen:
			if got != w {
				t.Fatalf("method = %q, want %q", got, w)
			}
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out waiting for %q", w)
		}
	}
}

func TestPendingCancelCarriesRequestID(t *testing.T) {
	serverReader, clientWriter := io.Pipe()
	clientReader, serverWriter := io.Pipe()
	t.Cleanup(func() {
		_ = clientWriter.Close()
		_ = serverWriter.Close()
	})

	conn := NewConn(pipeRWC{Reader: serverReader, Writer: serverWriter}, NewDispatcher())

	peer := NewConn(pipeRWC{Reader: clientReader, Writer: clientWriter}, NewDispatcher())
	notifs := make(chan *Notification, 1)
	go func() {
		for {
			msg, err := peer.ReadMessage()
			if err != nil {
				return
			}
			if n, ok := msg.(*Notification); ok {
				notifs <- n
			}
		}
	}()

	p, err := conn.CallAsync("textDocument/completion", nil)
	if err != nil {
		t.Fatal(err)
	}

	id, ok := p.ID().Int()
	if !ok {
		t.Fatalf("ID() = %v, want an integer", p.ID())
	}

	p.Cancel()

	select {
	case n := <-notifs:
		if n.Method != "$/cancelRequest" {
			t.Fatalf("method = %q, want $/cancelRequest", n.Method)
		}
		var params struct{ ID int64 }
		if err := json.Unmarshal(n.Params, &params); err != nil {
			t.Fatal(err)
		}
		if params.ID != id {
			t.Fatalf("cancelled id = %d, want %d", params.ID, id)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("no $/cancelRequest received")
	}
}

func TestWaitFailsWhenConnectionCloses(t *testing.T) {
	serverReader, clientWriter := io.Pipe()
	clientReader, serverWriter := io.Pipe()

	conn := NewConn(pipeRWC{Reader: serverReader, Writer: serverWriter}, NewDispatcher())

	go func() { _, _ = io.Copy(io.Discard, clientReader) }()

	serveDone := make(chan struct{})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer close(serveDone)
		_ = conn.Serve(ctx)
	}()

	p, err := conn.CallAsync("workspace/configuration", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Stop Serve, standing in for a server process that has died.
	cancel()
	<-serveDone

	start := time.Now()
	_, err = p.Wait(context.Background())
	if !errors.Is(err, ErrConnClosed) {
		t.Fatalf("err = %v, want ErrConnClosed", err)
	}
	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Fatalf("Wait took %s, expected to fail immediately", elapsed)
	}

	_ = clientWriter.Close()
	_ = serverWriter.Close()
}

// slowFirstWrite delays only the first write, so a cancel notification sent
// afterwards does not distort timings.
type slowFirstWrite struct {
	delay time.Duration
	wrote atomic.Bool
	block chan struct{}
}

func (s *slowFirstWrite) Read(_ []byte) (int, error) {
	<-s.block
	return 0, io.EOF
}

func (s *slowFirstWrite) Write(p []byte) (int, error) {
	if s.wrote.CompareAndSwap(false, true) {
		time.Sleep(s.delay)
	}
	return len(p), nil
}

func (s *slowFirstWrite) Close() error { return nil }

func TestCallTimeoutStartsBeforeTheWrite(t *testing.T) {
	const writeDelay = 300 * time.Millisecond
	const callTimeout = 400 * time.Millisecond

	rwc := &slowFirstWrite{delay: writeDelay, block: make(chan struct{})}
	t.Cleanup(func() { close(rwc.block) })

	conn := NewConn(rwc, NewDispatcher())
	conn.SetCallTimeout(callTimeout)

	start := time.Now()
	_, err := conn.Call(context.Background(), "workspace/applyEdit", nil)
	elapsed := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("err = %v, want DeadlineExceeded", err)
	}

	// Starting the deadline after the write would push this to writeDelay+callTimeout.
	if limit := callTimeout + writeDelay/2; elapsed > limit {
		t.Fatalf("call took %s, want under %s: deadline did not start until after the write", elapsed, limit)
	}
}
