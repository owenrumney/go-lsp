package jsonrpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

func startServe(t *testing.T, conn *Conn) <-chan error {
	t.Helper()
	done := make(chan error, 1)
	go func() {
		done <- conn.Serve(t.Context())
	}()
	return done
}

func waitServe(t *testing.T, done <-chan error) error {
	t.Helper()
	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		t.Fatal("Serve did not return")
		return nil
	}
}

func writeFrame(t *testing.T, w net.Conn, body string) {
	t.Helper()
	if _, err := fmt.Fprintf(w, "Content-Length: %d\r\n\r\n%s", len(body), body); err != nil {
		t.Fatal(err)
	}
}

func readResponse(t *testing.T, conn *Conn) *Response {
	t.Helper()
	msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatal(err)
	}
	resp, ok := msg.(*Response)
	if !ok {
		t.Fatalf("expected *Response, got %T", msg)
	}
	return resp
}

func TestConn_Serve_ExitNotification(t *testing.T) {
	clientEnd, serverEnd := net.Pipe()

	d := NewDispatcher()
	d.RegisterNotification("exit", func(_ context.Context, _ json.RawMessage) error {
		return ErrExit
	})
	conn := NewConn(serverEnd, d)
	done := startServe(t, conn)

	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","method":"exit"}`)

	if err := waitServe(t, done); !errors.Is(err, ErrExit) {
		t.Fatalf("Serve returned %v, want ErrExit", err)
	}
}

func TestConn_Serve_ParseErrorRecovery(t *testing.T) {
	clientEnd, serverEnd := net.Pipe()
	defer func() { _ = clientEnd.Close() }()

	d := NewDispatcher()
	d.RegisterMethod("ping", func(_ context.Context, _ json.RawMessage) (any, error) {
		return "pong", nil
	})
	serverConn := NewConn(serverEnd, d)
	done := startServe(t, serverConn)

	clientConn := NewConn(clientEnd, NewDispatcher())

	writeFrame(t, clientEnd, `{not json`)
	resp := readResponse(t, clientConn)
	if resp.Error == nil || resp.Error.Code != CodeParseError {
		t.Fatalf("expected ParseError response, got %+v", resp)
	}

	writeFrame(t, clientEnd, `{"jsonrpc":"2.0"}`)
	resp = readResponse(t, clientConn)
	if resp.Error == nil || resp.Error.Code != CodeInvalidRequest {
		t.Fatalf("expected InvalidRequest response, got %+v", resp)
	}

	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","id":1,"method":"ping"}`)
	resp = readResponse(t, clientConn)
	if resp.Error != nil {
		t.Fatalf("ping failed: %v", resp.Error)
	}
	if string(resp.Result) != `"pong"` {
		t.Fatalf("ping result = %s, want \"pong\"", resp.Result)
	}

	_ = clientEnd.Close()
	if err := waitServe(t, done); err == nil {
		t.Fatal("expected Serve to return an error after peer close")
	}
}

func TestConn_Serve_InvalidContentLength(t *testing.T) {
	tests := []struct {
		name   string
		header string
	}{
		{"negative", "Content-Length: -1\r\n\r\n"},
		{"zero", "Content-Length: 0\r\n\r\n"},
		{"huge", fmt.Sprintf("Content-Length: %d\r\n\r\n", maxContentLength+1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientEnd, serverEnd := net.Pipe()
			defer func() { _ = clientEnd.Close() }()

			conn := NewConn(serverEnd, NewDispatcher())
			done := startServe(t, conn)

			if _, err := clientEnd.Write([]byte(tt.header)); err != nil {
				t.Fatal(err)
			}

			if err := waitServe(t, done); err == nil {
				t.Fatal("expected Serve to return an error")
			}
		})
	}
}

func TestConn_Serve_ContextCancelUnblocksRead(t *testing.T) {
	clientEnd, serverEnd := net.Pipe()
	defer func() { _ = clientEnd.Close() }()

	conn := NewConn(serverEnd, NewDispatcher())
	ctx, cancel := context.WithCancel(t.Context())
	done := make(chan error, 1)
	go func() {
		done <- conn.Serve(ctx)
	}()

	cancel()
	if err := waitServe(t, done); !errors.Is(err, context.Canceled) {
		t.Fatalf("Serve returned %v, want context.Canceled", err)
	}
}

func TestConn_NotificationHandlerCanCallPeer(t *testing.T) {
	clientEnd, serverEnd := net.Pipe()
	defer func() { _ = clientEnd.Close() }()

	handled := make(chan error, 1)

	serverDisp := NewDispatcher()
	serverConn := NewConn(serverEnd, serverDisp)
	serverDisp.RegisterNotification("kick", func(ctx context.Context, _ json.RawMessage) error {
		// Must not deadlock the connection.
		callCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		_, err := serverConn.Call(callCtx, "peer/method", nil)
		handled <- err
		return nil
	})

	clientDisp := NewDispatcher()
	clientDisp.RegisterMethod("peer/method", func(_ context.Context, _ json.RawMessage) (any, error) {
		return "ok", nil
	})
	clientConn := NewConn(clientEnd, clientDisp)

	startServe(t, serverConn)
	startServe(t, clientConn)

	if err := clientConn.Notify(t.Context(), "kick", nil); err != nil {
		t.Fatal(err)
	}

	select {
	case err := <-handled:
		if err != nil {
			t.Fatalf("call from notification handler failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("deadlock: notification handler's call never completed")
	}
}

func TestConn_NotificationOrderPreserved(t *testing.T) {
	clientEnd, serverEnd := net.Pipe()
	defer func() { _ = clientEnd.Close() }()

	const count = 20
	got := make(chan int, count)

	d := NewDispatcher()
	d.RegisterNotification("seq", func(_ context.Context, params json.RawMessage) error {
		var p struct{ N int }
		if err := json.Unmarshal(params, &p); err != nil {
			return err
		}
		got <- p.N
		return nil
	})
	serverConn := NewConn(serverEnd, d)
	startServe(t, serverConn)

	clientConn := NewConn(clientEnd, NewDispatcher())
	for i := 0; i < count; i++ {
		if err := clientConn.Notify(t.Context(), "seq", map[string]int{"N": i}); err != nil {
			t.Fatal(err)
		}
	}

	for i := 0; i < count; i++ {
		select {
		case n := <-got:
			if n != i {
				t.Fatalf("notification %d arrived out of order (got %d)", i, n)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out waiting for notification %d", i)
		}
	}
}

func TestConn_Serve_ExitWinsOverEOF(t *testing.T) {
	for range 25 {
		clientEnd, serverEnd := net.Pipe()

		d := NewDispatcher()
		d.RegisterNotification("exit", func(_ context.Context, _ json.RawMessage) error {
			return ErrExit
		})
		conn := NewConn(serverEnd, d)
		done := startServe(t, conn)

		writeFrame(t, clientEnd, `{"jsonrpc":"2.0","method":"exit"}`)
		_ = clientEnd.Close()

		if err := waitServe(t, done); !errors.Is(err, ErrExit) {
			t.Fatalf("Serve returned %v, want ErrExit when exit precedes close", err)
		}
	}
}

func TestConn_Serve_BadRequestIDRecoverable(t *testing.T) {
	clientEnd, serverEnd := net.Pipe()
	defer func() { _ = clientEnd.Close() }()

	d := NewDispatcher()
	d.RegisterMethod("ping", func(_ context.Context, _ json.RawMessage) (any, error) {
		return "pong", nil
	})
	serverConn := NewConn(serverEnd, d)
	startServe(t, serverConn)

	clientConn := NewConn(clientEnd, NewDispatcher())

	for _, body := range []string{
		`{"jsonrpc":"2.0","id":{"a":1},"method":"x"}`,
		`{"jsonrpc":"2.0","id":1.5,"method":"x"}`,
		`{"jsonrpc":"2.0","id":true,"method":"x"}`,
	} {
		writeFrame(t, clientEnd, body)
		resp := readResponse(t, clientConn)
		if resp.Error == nil || resp.Error.Code != CodeInvalidRequest {
			t.Fatalf("body %s: expected InvalidRequest, got %+v", body, resp)
		}
	}

	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","id":1,"method":"ping"}`)
	if resp := readResponse(t, clientConn); resp.Error != nil {
		t.Fatalf("connection did not survive bad IDs: %v", resp.Error)
	}
}

func TestConn_CancelQueuedRequestGetsResponse(t *testing.T) {
	clientEnd, serverEnd := net.Pipe()
	defer func() { _ = clientEnd.Close() }()

	started := make(chan struct{})
	release := make(chan struct{})
	ran := make(chan string, 2)

	d := NewDispatcher()
	d.RegisterMethod("work", func(_ context.Context, params json.RawMessage) (any, error) {
		ran <- string(params)
		started <- struct{}{}
		<-release
		return "done", nil
	})
	serverConn := NewConn(serverEnd, d)
	serverConn.SetMaxConcurrentRequests(1)
	startServe(t, serverConn)

	clientConn := NewConn(clientEnd, NewDispatcher())

	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","id":1,"method":"work","params":"1"}`)
	<-started
	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","id":2,"method":"work","params":"2"}`)
	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","method":"$/cancelRequest","params":{"id":2}}`)

	resp := readResponse(t, clientConn)
	if resp.ID.String() != "2" || resp.Error == nil || resp.Error.Code != CodeRequestCancelled {
		t.Fatalf("expected RequestCancelled for id 2, got %+v", resp)
	}

	close(release)
	resp = readResponse(t, clientConn)
	if resp.ID.String() != "1" || resp.Error != nil {
		t.Fatalf("first request failed: %+v", resp)
	}
	if len(ran) != 1 {
		t.Fatalf("cancelled request still executed")
	}
}

func TestConn_ServeWaitsForHandlers(t *testing.T) {
	clientEnd, serverEnd := net.Pipe()

	started := make(chan struct{})
	release := make(chan struct{})
	var finished atomic.Bool

	d := NewDispatcher()
	d.RegisterMethod("work", func(_ context.Context, _ json.RawMessage) (any, error) {
		started <- struct{}{}
		<-release
		finished.Store(true)
		return nil, nil
	})
	conn := NewConn(serverEnd, d)
	done := startServe(t, conn)

	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","id":1,"method":"work"}`)
	<-started
	_ = clientEnd.Close()
	close(release)

	_ = waitServe(t, done)
	if !finished.Load() {
		t.Fatal("Serve returned before the request handler finished")
	}
}

func TestConn_CancelRequestInDispatchQueue(t *testing.T) {
	clientEnd, serverEnd := net.Pipe()
	defer func() { _ = clientEnd.Close() }()

	started := make(chan struct{})
	release := make(chan struct{})
	ran := make(chan struct{}, 1)

	d := NewDispatcher()
	d.RegisterNotification("pause", func(_ context.Context, _ json.RawMessage) error {
		started <- struct{}{}
		<-release
		return nil
	})
	d.RegisterMethod("work", func(_ context.Context, _ json.RawMessage) (any, error) {
		ran <- struct{}{}
		return "done", nil
	})
	serverConn := NewConn(serverEnd, d)
	startServe(t, serverConn)

	clientConn := NewConn(clientEnd, NewDispatcher())

	// Park the worker so the request stays in the dispatch queue while the
	// cancel is processed inline.
	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","method":"pause"}`)
	<-started
	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","id":1,"method":"work"}`)
	writeFrame(t, clientEnd, `{"jsonrpc":"2.0","method":"$/cancelRequest","params":{"id":1}}`)
	// Frames are processed in order: the ParseError reply to this garbage
	// proves the cancel has been handled before the worker is released.
	writeFrame(t, clientEnd, `{garbage`)
	if resp := readResponse(t, clientConn); resp.Error == nil || resp.Error.Code != CodeParseError {
		t.Fatalf("expected ParseError barrier, got %+v", resp)
	}
	close(release)

	resp := readResponse(t, clientConn)
	if resp.Error == nil || resp.Error.Code != CodeRequestCancelled {
		t.Fatalf("expected RequestCancelled, got %+v", resp)
	}
	select {
	case <-ran:
		t.Fatal("cancelled request still executed")
	default:
	}
}

func TestConn_ClosesTransportOnReturn(t *testing.T) {
	tests := map[string]func(t *testing.T, clientEnd net.Conn, done <-chan error){
		"exit": func(t *testing.T, clientEnd net.Conn, done <-chan error) {
			writeFrame(t, clientEnd, `{"jsonrpc":"2.0","method":"exit"}`)
			if err := waitServe(t, done); !errors.Is(err, ErrExit) {
				t.Fatalf("Serve returned %v", err)
			}
		},
	}
	for name, drive := range tests {
		t.Run(name, func(t *testing.T) {
			clientEnd, serverEnd := net.Pipe()
			defer func() { _ = clientEnd.Close() }()

			d := NewDispatcher()
			d.RegisterNotification("exit", func(_ context.Context, _ json.RawMessage) error {
				return ErrExit
			})
			conn := NewConn(serverEnd, d)
			done := startServe(t, conn)

			drive(t, clientEnd, done)

			// The peer must observe the transport closing.
			_ = clientEnd.SetReadDeadline(time.Now().Add(5 * time.Second))
			buf := make([]byte, 1)
			if _, err := clientEnd.Read(buf); err == nil {
				t.Fatal("transport still open after Serve returned")
			} else if errors.Is(err, os.ErrDeadlineExceeded) {
				t.Fatal("transport was not closed after Serve returned")
			}
		})
	}
}
