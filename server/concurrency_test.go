package server_test

import (
	"context"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
	"github.com/owenrumney/go-lsp/servertest"
)

type blockingHoverHandler struct {
	started chan struct{}
	release chan struct{}
}

func (h *blockingHoverHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{}, nil
}

func (h *blockingHoverHandler) Shutdown(_ context.Context) error { return nil }

func (h *blockingHoverHandler) Hover(_ context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	h.started <- struct{}{}
	<-h.release
	return &lsp.Hover{Contents: lsp.NewHoverContents(lsp.PlainText, "ok")}, nil
}

func TestDefaultRequestHandlingRemainsConcurrent(t *testing.T) {
	handler := &blockingHoverHandler{
		started: make(chan struct{}, 2),
		release: make(chan struct{}, 2),
	}
	h := servertest.New(t, handler)

	call1, err := h.CallAsync("textDocument/hover", &lsp.HoverParams{})
	if err != nil {
		t.Fatal(err)
	}
	call2, err := h.CallAsync("textDocument/hover", &lsp.HoverParams{})
	if err != nil {
		t.Fatal(err)
	}

	waitStarted(t, handler.started)
	waitStarted(t, handler.started)

	handler.release <- struct{}{}
	handler.release <- struct{}{}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	if _, err := call1.Wait(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := call2.Wait(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestMaxConcurrentRequestsLimitsExecution(t *testing.T) {
	handler := &blockingHoverHandler{
		started: make(chan struct{}, 2),
		release: make(chan struct{}, 2),
	}
	h := servertest.New(t, handler, servertest.WithServerOptions(
		server.WithMaxConcurrentRequests(1),
	))

	call1, err := h.CallAsync("textDocument/hover", &lsp.HoverParams{})
	if err != nil {
		t.Fatal(err)
	}
	waitStarted(t, handler.started)

	call2, err := h.CallAsync("textDocument/hover", &lsp.HoverParams{})
	if err != nil {
		t.Fatal(err)
	}

	select {
	case <-handler.started:
		t.Fatal("second request started before first finished")
	case <-time.After(150 * time.Millisecond):
	}

	handler.release <- struct{}{}
	waitStarted(t, handler.started)
	handler.release <- struct{}{}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	if _, err := call1.Wait(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := call2.Wait(ctx); err != nil {
		t.Fatal(err)
	}
}

func waitStarted(t *testing.T, ch <-chan struct{}) {
	t.Helper()
	select {
	case <-ch:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for request to start")
	}
}
