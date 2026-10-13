// Command lspclient drives another language server and prints what it knows
// about a file.
//
//	go run ./examples/lspclient -server gopls -file ./client/client.go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/owenrumney/go-lsp/client"
	"github.com/owenrumney/go-lsp/lsp"
)

// collector keeps the diagnostics the server pushes, which arrive as
// notifications rather than as a response to anything.
type collector struct {
	client.NopHandler

	mu    sync.Mutex
	diags map[lsp.DocumentURI][]lsp.Diagnostic
}

func (c *collector) Diagnostics(params *lsp.PublishDiagnosticsParams) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.diags == nil {
		c.diags = make(map[lsp.DocumentURI][]lsp.Diagnostic)
	}
	c.diags[params.URI] = params.Diagnostics
}

func (c *collector) forURI(uri lsp.DocumentURI) []lsp.Diagnostic {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.diags[uri]
}

// Configuration answers the settings request gopls makes during startup.
func (c *collector) Configuration(_ context.Context, items []lsp.ConfigurationItem) ([]json.RawMessage, error) {
	out := make([]json.RawMessage, len(items))
	for i := range out {
		out[i] = json.RawMessage(`{}`)
	}
	return out, nil
}

func main() {
	serverCmd := flag.String("server", "gopls", "language server command")
	file := flag.String("file", "", "file to inspect")
	line := flag.Int("line", 0, "zero-based line for the hover query")
	char := flag.Int("char", 0, "zero-based character for the hover query")
	flag.Parse()

	if *file == "" {
		log.Fatal("-file is required")
	}
	if err := run(*serverCmd, *file, *line, *char); err != nil {
		log.Fatal(err)
	}
}

func run(serverCmd, file string, line, char int) error {
	abs, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	text, err := os.ReadFile(abs) // #nosec G304 -- the user names the file to inspect
	if err != nil {
		return err
	}

	root := filepath.Dir(abs)
	uri := lsp.DocumentURI("file://" + abs)
	handler := &collector{}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	c, err := client.Dial(ctx, client.Config{
		Command: strings.Fields(serverCmd),
		Dir:     root,
		RootURI: lsp.DocumentURI("file://" + root),
		Handler: handler,
		Stderr:  os.Stderr,
	})
	if err != nil {
		return fmt.Errorf("starting %s: %w", serverCmd, err)
	}
	defer func() { _ = c.Close() }()

	if info := c.ServerInfo(); info != nil {
		fmt.Printf("connected to %s\n", info.Name)
	}

	if err := c.DidOpen(ctx, uri, languageID(abs), 1, string(text)); err != nil {
		return err
	}

	// Indexing has to finish first, or every answer below is empty but valid.
	if err := c.Ready(ctx); err != nil {
		return fmt.Errorf("waiting for the server to be ready: %w", err)
	}

	if err := printHover(ctx, c, uri, line, char); err != nil {
		return err
	}
	if err := printSymbols(ctx, c, uri); err != nil {
		return err
	}

	for _, d := range handler.forURI(uri) {
		fmt.Printf("diagnostic %d:%d %s\n", d.Range.Start.Line+1, d.Range.Start.Character+1, d.Message)
	}
	return nil
}

func printHover(ctx context.Context, c *client.Client, uri lsp.DocumentURI, line, char int) error {
	if !c.Supports("textDocument/hover", uri, "go") {
		fmt.Println("server provides no hover")
		return nil
	}

	hover, err := c.Hover(ctx, uri, lsp.Position{Line: line, Character: char})
	if err != nil {
		return err
	}
	if hover == nil {
		fmt.Printf("no hover at %d:%d\n", line, char)
		return nil
	}

	fmt.Printf("hover at %d:%d: %s\n", line, char, summarise(hover.Contents.Value()))
	return nil
}

func printSymbols(ctx context.Context, c *client.Client, uri lsp.DocumentURI) error {
	if !c.Supports("textDocument/documentSymbol", uri, "go") {
		fmt.Println("server provides no document symbols")
		return nil
	}

	symbols, err := c.DocumentSymbol(ctx, uri)
	if err != nil {
		return err
	}
	for _, s := range symbols {
		fmt.Printf("symbol %s (kind %d)\n", s.Name, s.Kind)
	}
	return nil
}

func languageID(path string) string {
	switch filepath.Ext(path) {
	case ".go":
		return "go"
	case ".ts":
		return "typescript"
	case ".py":
		return "python"
	default:
		return "plaintext"
	}
}

// summarise picks the first meaningful line of hover markdown, skipping the
// code fences servers wrap signatures in.
func summarise(s string) string {
	for line := range strings.SplitSeq(s, "\n") {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "```") {
			return line
		}
	}
	return "(empty)"
}
