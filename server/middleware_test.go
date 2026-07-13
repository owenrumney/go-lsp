package server_test

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
	"github.com/owenrumney/go-lsp/servertest"
)

type middlewareHandler struct {
	events chan string
}

func (h *middlewareHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{}, nil
}

func (h *middlewareHandler) Shutdown(_ context.Context) error { return nil }

func (h *middlewareHandler) Hover(_ context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	h.events <- "handler:hover"
	return &lsp.Hover{Contents: lsp.MarkupContent{Kind: lsp.PlainText, Value: "ok"}}, nil
}

func (h *middlewareHandler) DidOpen(_ context.Context, _ *lsp.DidOpenTextDocumentParams) error {
	h.events <- "handler:didOpen"
	return nil
}

func (h *middlewareHandler) DidChange(_ context.Context, _ *lsp.DidChangeTextDocumentParams) error {
	return nil
}

func (h *middlewareHandler) DidClose(_ context.Context, _ *lsp.DidCloseTextDocumentParams) error {
	return nil
}

func TestMethodMiddlewareOrder(t *testing.T) {
	handler := &middlewareHandler{events: make(chan string, 10)}
	h := servertest.New(t, handler, servertest.WithServerOptions(
		server.WithMethodMiddleware(
			func(method string, next server.MethodHandlerFunc) server.MethodHandlerFunc {
				return func(ctx context.Context, params json.RawMessage) (any, error) {
					handler.events <- "mw1:before:" + method
					result, err := next(ctx, params)
					handler.events <- "mw1:after:" + method
					return result, err
				}
			},
			func(method string, next server.MethodHandlerFunc) server.MethodHandlerFunc {
				return func(ctx context.Context, params json.RawMessage) (any, error) {
					handler.events <- "mw2:before:" + method
					result, err := next(ctx, params)
					handler.events <- "mw2:after:" + method
					return result, err
				}
			},
		),
	))
	handler.events = make(chan string, 10)

	if _, err := h.Hover("file:///test.txt", 0, 0); err != nil {
		t.Fatal(err)
	}

	got := collectEvents(t, handler.events, 5)
	want := []string{
		"mw1:before:textDocument/hover",
		"mw2:before:textDocument/hover",
		"handler:hover",
		"mw2:after:textDocument/hover",
		"mw1:after:textDocument/hover",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("events = %v, want %v", got, want)
	}
}

func TestNotificationMiddlewareCanShortCircuit(t *testing.T) {
	handler := &middlewareHandler{events: make(chan string, 10)}
	h := servertest.New(t, handler, servertest.WithServerOptions(
		server.WithNotificationMiddleware(
			func(method string, next server.NotificationHandlerFunc) server.NotificationHandlerFunc {
				return func(ctx context.Context, params json.RawMessage) error {
					handler.events <- "mw1:before:" + method
					return next(ctx, params)
				}
			},
			func(method string, next server.NotificationHandlerFunc) server.NotificationHandlerFunc {
				return func(ctx context.Context, params json.RawMessage) error {
					handler.events <- "mw2:block:" + method
					return nil
				}
			},
		),
	))
	handler.events = make(chan string, 10)

	if err := h.DidOpen("file:///test.txt", "plaintext", "hello"); err != nil {
		t.Fatal(err)
	}

	got := collectEvents(t, handler.events, 2)
	want := []string{
		"mw1:before:textDocument/didOpen",
		"mw2:block:textDocument/didOpen",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("events = %v, want %v", got, want)
	}
}

func collectEvents(t *testing.T, ch <-chan string, n int) []string {
	t.Helper()
	out := make([]string, 0, n)
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	for len(out) < n {
		select {
		case event := <-ch:
			out = append(out, event)
		case <-timer.C:
			t.Fatalf("timed out waiting for %d events, got %v", n, out)
		}
	}
	return out
}
