package servertest

import (
	"context"
	"encoding/json"
	"net"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

// Harness provides a test harness that simulates an LSP client over in-memory pipes.
type Harness struct {
	t              testing.TB
	conn           *rpcConn
	notifs         *notifStore
	clientRequests *clientRequestStore
	cancel         context.CancelFunc
	ctx            context.Context

	// InitResult holds the result from the initialize request.
	InitResult *lsp.InitializeResult
}

// cleanupTimeout bounds each stage of harness shutdown.
const cleanupTimeout = 5 * time.Second

// New creates a new test harness, starts the server, performs initialization,
// and registers cleanup to shut down gracefully.
func New(t testing.TB, handler server.LifecycleHandler, opts ...Option) *Harness {
	t.Helper()

	cfg := &config{}
	for _, opt := range opts {
		opt(cfg)
	}

	ctx, cancel := context.WithCancel(context.Background())

	clientConn, serverConn := net.Pipe()

	srv := server.NewServer(handler, cfg.serverOpts...)

	// Start the server in a goroutine.
	serverDone := make(chan error, 1)
	go func() {
		serverDone <- srv.Run(ctx, serverConn)
	}()

	notifs := newNotifStore()
	clientRequests := newClientRequestStore()

	dispatcher := jsonrpc.NewDispatcher()
	dispatcher.SetUnknownNotificationHandler(func(_ context.Context, method string, params json.RawMessage) error {
		notifs.add(method, params)
		return nil
	})
	dispatcher.SetUnknownMethodHandler(func(_ context.Context, method string, params json.RawMessage) (any, error) {
		return clientRequests.handle(method, params)
	})

	rpc := newRPCConn(clientConn, dispatcher)
	go func() { _ = rpc.serve(ctx) }()

	h := &Harness{
		t:              t,
		conn:           rpc,
		notifs:         notifs,
		clientRequests: clientRequests,
		cancel:         cancel,
		ctx:            ctx,
	}

	// Send initialize request.
	initParams := cfg.initParams
	if initParams == nil {
		pid := 1
		initParams = &lsp.InitializeParams{
			ProcessID: &pid,
			ClientInfo: &lsp.ClientInfo{
				Name:    "servertest",
				Version: "0.1.0",
			},
			Capabilities: lsp.ClientCapabilities{},
		}
	}

	result, err := rpc.call(ctx, "initialize", initParams)
	if err != nil {
		cancel()
		_ = clientConn.Close()
		t.Fatalf("initialize failed: %v", err)
	}

	var initResult lsp.InitializeResult
	if err := json.Unmarshal(result, &initResult); err != nil {
		cancel()
		_ = clientConn.Close()
		t.Fatalf("unmarshal InitializeResult: %v", err)
	}
	h.InitResult = &initResult

	// Send initialized notification.
	if err := rpc.notify(ctx, "initialized", &lsp.InitializedParams{}); err != nil {
		cancel()
		_ = clientConn.Close()
		t.Fatalf("initialized notification failed: %v", err)
	}

	t.Cleanup(func() {
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cleanupTimeout)
		defer shutdownCancel()

		// Send shutdown request (ignore errors on already-closed connections).
		_, _ = rpc.call(shutdownCtx, "shutdown", nil)

		// Send exit notification.
		_ = rpc.notify(shutdownCtx, "exit", nil)

		// Cancel the context to stop the server.
		cancel()

		// Close the client side of the pipe.
		_ = clientConn.Close()

		// Wait for server to finish, but fail the test rather than hang.
		select {
		case <-serverDone:
		case <-time.After(cleanupTimeout):
			t.Errorf("servertest: server did not exit within %s", cleanupTimeout)
		}
	})

	return h
}
