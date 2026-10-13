package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

// Transport opens the stream carrying LSP traffic. It replaces Config.Command
// for in-process servers and tests.
type Transport func(ctx context.Context) (io.ReadWriteCloser, error)

// Config describes a language server to connect to.
type Config struct {
	// Command is the argv of the server to spawn; argv[0] is resolved on PATH.
	// Ignored when Transport is set.
	Command []string
	// Transport overrides Command with an already-open stream.
	Transport Transport

	Dir string
	Env []string

	RootURI          lsp.DocumentURI
	WorkspaceFolders []lsp.WorkspaceFolder
	ClientInfo       *lsp.ClientInfo

	// Capabilities declares what the client supports. DefaultCapabilities is
	// used when nil.
	Capabilities *lsp.ClientCapabilities
	// Settings becomes initializationOptions.
	Settings any

	// Handler answers server-to-client traffic. NopHandler is used when nil.
	Handler Handler

	Stderr io.Writer
	Logger *slog.Logger

	// ShutdownTimeout bounds the shutdown request and process exit during
	// Close. Defaults to 5s.
	ShutdownTimeout time.Duration
	// ReadyGrace is how long Ready waits for a server that never reports
	// progress before declaring it ready. Defaults to 2s.
	ReadyGrace time.Duration
}

// Client is a connection to one language server.
type Client struct {
	conn   *jsonrpc.Conn
	cmd    *exec.Cmd
	cancel context.CancelFunc
	logger *slog.Logger

	shutdownTimeout time.Duration

	caps       lsp.ServerCapabilities
	serverInfo *lsp.ServerInfo
	encoding   lsp.PositionEncodingKind

	registry *registry
	ready    *readiness

	done     chan struct{}
	errOnce  sync.Once
	errValue error

	closeOnce sync.Once
}

// Dial starts a language server, initializes it, and returns once the server
// has answered initialize. ctx bounds startup only; the connection lives until
// Close.
func Dial(ctx context.Context, cfg Config) (*Client, error) {
	cfg.applyDefaults()

	runCtx, cancel := context.WithCancel(context.Background())

	c := &Client{
		cancel:          cancel,
		logger:          cfg.Logger,
		shutdownTimeout: cfg.ShutdownTimeout,
		registry:        newRegistry(),
		ready:           newReadiness(cfg.ReadyGrace),
		done:            make(chan struct{}),
	}

	rwc, err := c.open(runCtx, cfg)
	if err != nil {
		cancel()
		return nil, err
	}

	c.conn = jsonrpc.NewConn(rwc, c.dispatcher(cfg.Handler))

	go func() {
		err := c.conn.Serve(runCtx)
		if waitErr := c.waitProcess(); waitErr != nil && err == nil {
			err = waitErr
		}
		c.setErr(err)
		close(c.done)
	}()

	if err := c.initialize(ctx, cfg); err != nil {
		_ = c.Close()
		return nil, err
	}
	c.ready.begin()
	return c, nil
}

// Capabilities returns the capabilities the server advertised at initialize.
// Capabilities added later by dynamic registration are not reflected here; use
// Supports instead.
func (c *Client) Capabilities() lsp.ServerCapabilities { return c.caps }

// ServerInfo returns the server's self-reported name and version, if it sent any.
func (c *Client) ServerInfo() *lsp.ServerInfo { return c.serverInfo }

// PositionEncoding returns the encoding the server selected, defaulting to UTF-16.
func (c *Client) PositionEncoding() lsp.PositionEncodingKind { return c.encoding }

// Done is closed when the connection ends, whether by Close or by the server exiting.
func (c *Client) Done() <-chan struct{} { return c.done }

// Err returns why the connection ended, or nil while it is live.
func (c *Client) Err() error {
	select {
	case <-c.done:
		return c.errValue
	default:
		return nil
	}
}

// Close shuts the server down gracefully, then kills it if it does not exit
// within ShutdownTimeout.
func (c *Client) Close() error {
	var err error
	c.closeOnce.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), c.shutdownTimeout)
		defer cancel()

		// Writes ignore ctx, so a send blocked on a full pipe would hold the
		// write lock past the timeout; cancelling the run context kills it.
		stop := context.AfterFunc(ctx, c.cancel)
		defer stop()

		if _, callErr := c.conn.Call(ctx, "shutdown", nil); callErr != nil {
			err = callErr
		}
		if notifyErr := c.conn.Notify(ctx, "exit", nil); notifyErr != nil && err == nil {
			err = notifyErr
		}

		select {
		case <-c.done:
		case <-ctx.Done():
			c.cancel()
			<-c.done
		}
		c.cancel()
	})
	return err
}

func (cfg *Config) applyDefaults() {
	if cfg.Handler == nil {
		cfg.Handler = NopHandler{}
	}
	if cfg.Capabilities == nil {
		caps := DefaultCapabilities()
		cfg.Capabilities = &caps
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.New(slog.DiscardHandler)
	}
	if cfg.ShutdownTimeout <= 0 {
		cfg.ShutdownTimeout = 5 * time.Second
	}
	if cfg.ReadyGrace <= 0 {
		cfg.ReadyGrace = 2 * time.Second
	}
}

func (c *Client) open(ctx context.Context, cfg Config) (io.ReadWriteCloser, error) {
	if cfg.Transport != nil {
		return cfg.Transport(ctx)
	}
	if len(cfg.Command) == 0 {
		return nil, errors.New("client: Config needs Command or Transport")
	}

	cmd := exec.CommandContext(ctx, cfg.Command[0], cfg.Command[1:]...) // #nosec G204 -- the caller chooses the server to run
	cmd.Dir = cfg.Dir
	if len(cfg.Env) > 0 {
		cmd.Env = append(os.Environ(), cfg.Env...)
	}
	cmd.Stderr = cfg.Stderr
	// Without this, Wait blocks on a grandchild holding the stderr pipe open.
	cmd.WaitDelay = cfg.ShutdownTimeout

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("client: stdin pipe: %w", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("client: stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("client: starting %s: %w", cfg.Command[0], err)
	}

	c.cmd = cmd
	return procRWC{in: stdin, out: stdout}, nil
}

func (c *Client) initialize(ctx context.Context, cfg Config) error {
	params := &lsp.InitializeParams{
		ClientInfo:       cfg.ClientInfo,
		Capabilities:     *cfg.Capabilities,
		WorkspaceFolders: cfg.WorkspaceFolders,
	}
	if params.ClientInfo == nil {
		params.ClientInfo = &lsp.ClientInfo{Name: "go-lsp"}
	}
	if cfg.RootURI != "" {
		root := cfg.RootURI
		// rootUri is deprecated in favour of workspaceFolders, but gopls and
		// others still read it, so send both.
		params.RootURI = &root //nolint:staticcheck
		if len(params.WorkspaceFolders) == 0 {
			params.WorkspaceFolders = []lsp.WorkspaceFolder{{URI: root, Name: path.Base(string(root))}}
		}
	}
	if pid := os.Getpid(); pid != 0 {
		params.ProcessID = &pid
	}
	if cfg.Settings != nil {
		opts, err := json.Marshal(cfg.Settings)
		if err != nil {
			return fmt.Errorf("client: marshalling settings: %w", err)
		}
		params.InitializationOptions = opts
	}

	var result lsp.InitializeResult
	if err := c.call(ctx, "initialize", params, &result); err != nil {
		return fmt.Errorf("client: initialize: %w", err)
	}

	c.caps = result.Capabilities
	c.serverInfo = result.ServerInfo
	c.encoding = lsp.PositionEncodingUTF16
	if result.Capabilities.PositionEncoding != nil {
		c.encoding = *result.Capabilities.PositionEncoding
	}

	return c.conn.Notify(ctx, "initialized", &lsp.InitializedParams{})
}

func (c *Client) waitProcess() error {
	if c.cmd == nil {
		return nil
	}
	if err := c.cmd.Wait(); err != nil {
		return fmt.Errorf("client: server exited: %w", err)
	}
	return nil
}

func (c *Client) setErr(err error) {
	c.errOnce.Do(func() { c.errValue = err })
}

type procRWC struct {
	in  io.WriteCloser
	out io.ReadCloser
}

func (p procRWC) Read(b []byte) (int, error)  { return p.out.Read(b) }
func (p procRWC) Write(b []byte) (int, error) { return p.in.Write(b) }

func (p procRWC) Close() error {
	err := p.in.Close()
	if cerr := p.out.Close(); err == nil {
		err = cerr
	}
	return err
}
