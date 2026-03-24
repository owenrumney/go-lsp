package servertest

import (
	"context"
	"encoding/json"
	"net"
	"sync"
	"testing"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
)

// Harness provides a test harness that simulates an LSP client over in-memory pipes.
type Harness struct {
	t      testing.TB
	conn   *rpcConn
	notifs *notifStore
	cancel context.CancelFunc
	ctx    context.Context

	// InitResult holds the result from the initialize request.
	InitResult *lsp.InitializeResult

	// versions tracks document versions for auto-incrementing DidChange.
	versions   map[lsp.DocumentURI]int
	versionsMu sync.Mutex
}

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

	rpc := newRPCConn(clientConn)

	notifs := newNotifStore()
	rpc.notifHandler = func(method string, params json.RawMessage) {
		switch method {
		case "textDocument/publishDiagnostics":
			var p lsp.PublishDiagnosticsParams
			if err := json.Unmarshal(params, &p); err == nil {
				notifs.addDiagnostics(p)
			}
		case "window/showMessage":
			var p lsp.ShowMessageParams
			if err := json.Unmarshal(params, &p); err == nil {
				notifs.addMessage(p)
			}
		case "window/logMessage":
			var p lsp.LogMessageParams
			if err := json.Unmarshal(params, &p); err == nil {
				notifs.addLogMessage(p)
			}
		}
	}

	// Handle server-to-client requests with default success responses.
	rpc.requestHandler = func(_ string, _ json.RawMessage) (any, error) {
		return nil, nil
	}

	// Start the read loop.
	go rpc.readLoop()

	h := &Harness{
		t:        t,
		conn:     rpc,
		notifs:   notifs,
		cancel:   cancel,
		ctx:      ctx,
		versions: make(map[lsp.DocumentURI]int),
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
		// Send shutdown request (ignore errors on already-closed connections).
		_, _ = rpc.call(context.Background(), "shutdown", nil)

		// Send exit notification.
		_ = rpc.notify(context.Background(), "exit", nil)

		// Cancel the context to stop the server.
		cancel()

		// Close the client side of the pipe.
		_ = clientConn.Close()

		// Wait for server to finish.
		<-serverDone
	})

	return h
}
