# Getting Started with go-lsp

This guide walks you through building a language server using `go-lsp`.

## Prerequisites

- Go 1.25+
- An editor that supports LSP (VS Code, Neovim, etc.)

## Quick Start

Generate a project with the scaffold command:

```bash
go run github.com/owenrumney/go-lsp/cmd/scaffold@latest \
  --name mylang \
  --module github.com/you/mylang-lsp \
  --lang mylang \
  --features hover,completion,diagnostics
```

Or run it without flags for interactive prompts:

```bash
go run github.com/owenrumney/go-lsp/cmd/scaffold@latest
```

Available features: `hover`, `completion`, `diagnostics`, `definition`, `references`, `formatting`, `codeactions`, `symbols`.

This generates a ready-to-run project:

```
mylang-lsp/
  main.go                   — server entrypoint
  handler/handler.go        — handler with your selected interfaces
  handler/handler_test.go   — passing tests using the servertest harness
  go.mod
```

Build and run:

```bash
cd mylang-lsp
go build -o mylang-lsp .
```

Run the tests:

```bash
go test ./...
```

That's it — you have a working LSP server. Point your editor at the binary (see [Connect Your Editor](#connect-your-editor) below) and start iterating.

## What the Scaffold Generates

Here's what the generated code looks like with `hover,completion,diagnostics` selected.

### `main.go`

```go
package main

import (
    "context"
    "log"

    "github.com/you/mylang-lsp/handler"
    "github.com/owenrumney/go-lsp/server"
)

func main() {
    h := handler.New()
    srv := server.NewServer(h)
    if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
        log.Fatal(err)
    }
}
```

### `handler/handler.go`

```go
package handler

import (
    "context"
    "strings"

    "github.com/owenrumney/go-lsp/document"
    "github.com/owenrumney/go-lsp/lsp"
    "github.com/owenrumney/go-lsp/server"
)

type Handler struct {
    documents *document.Store
    client    *server.Client
}

func New() *Handler {
    return &Handler{
        documents: document.NewStore(),
    }
}

func (h *Handler) SetClient(client *server.Client) {
    h.client = client
}

func (h *Handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
    return &lsp.InitializeResult{
        ServerInfo: &lsp.ServerInfo{Name: "mylang", Version: "0.1.0"},
    }, nil
}

func (h *Handler) Shutdown(_ context.Context) error { return nil }

func (h *Handler) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
    _, err := h.documents.Open(params)
    return err
}

func (h *Handler) DidChange(_ context.Context, params *lsp.DidChangeTextDocumentParams) error {
    _, err := h.documents.Change(params)
    return err
}

func (h *Handler) DidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
    h.documents.Close(params)
    return nil
}

func (h *Handler) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error {
    var diags []lsp.Diagnostic

    text, ok := h.documents.Text(params.TextDocument.URI)
    if ok {
        for i, line := range strings.Split(text, "\n") {
            idx := strings.Index(line, "TODO")
            if idx < 0 {
                continue
            }
            sev := lsp.SeverityWarning
            diags = append(diags, lsp.Diagnostic{
                Range: lsp.Range{
                    Start: lsp.Position{Line: i, Character: idx},
                    End:   lsp.Position{Line: i, Character: idx + 4},
                },
                Severity: &sev,
                Source:   "mylang",
                Message:  "TODO found",
            })
        }
    }

    return h.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
        URI:         params.TextDocument.URI,
        Diagnostics: diags,
    })
}

func (h *Handler) Hover(_ context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
    return &lsp.Hover{
        Contents: lsp.MarkupContent{
            Kind:  lsp.Markdown,
            Value: "**mylang** hover",
        },
    }, nil
}

func (h *Handler) Completion(_ context.Context, _ *lsp.CompletionParams) (*lsp.CompletionList, error) {
    return &lsp.CompletionList{
        Items: []lsp.CompletionItem{
            {Label: "example"},
        },
    }, nil
}
```

### `handler/handler_test.go`

```go
package handler_test

import (
    "context"
    "testing"
    "time"

    "github.com/you/mylang-lsp/handler"
    "github.com/owenrumney/go-lsp/servertest"
)

func TestInitialize(t *testing.T) {
    h := servertest.New(t, handler.New())

    if h.InitResult.ServerInfo == nil || h.InitResult.ServerInfo.Name != "mylang" {
        t.Error("expected server info to be set")
    }
}

func TestHover(t *testing.T) {
    h := servertest.New(t, handler.New())

    h.DidOpen("file:///test.txt", "mylang", "hello world")

    hover, err := h.Hover("file:///test.txt", 0, 0)
    if err != nil {
        t.Fatal(err)
    }
    if hover == nil {
        t.Fatal("expected hover result")
    }
}

func TestCompletion(t *testing.T) {
    h := servertest.New(t, handler.New())

    h.DidOpen("file:///test.txt", "mylang", "hello")

    result, err := h.Completion("file:///test.txt", 0, 5)
    if err != nil {
        t.Fatal(err)
    }
    if len(result.Items) == 0 {
        t.Error("expected at least one completion item")
    }
}

func TestDiagnostics(t *testing.T) {
    h := servertest.New(t, handler.New())

    h.DidOpen("file:///test.txt", "mylang", "TODO fix this")
    h.DidSave("file:///test.txt")

    ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
    defer cancel()

    diags, err := h.WaitForDiagnostics(ctx, "file:///test.txt")
    if err != nil {
        t.Fatal(err)
    }
    if len(diags) == 0 {
        t.Error("expected diagnostics for TODO")
    }
}
```

The scaffold gives you working stubs — replace the placeholder implementations with your real logic.

## How go-lsp Works

You write a handler struct and implement interfaces for the LSP features you want. `go-lsp` auto-detects which interfaces your handler implements, registers the JSON-RPC methods, and advertises the right capabilities to the client.

The only required interface is `LifecycleHandler`:

```go
type LifecycleHandler interface {
    Initialize(ctx context.Context, params *lsp.InitializeParams) (*lsp.InitializeResult, error)
    Shutdown(ctx context.Context) error
}
```

Everything else is opt-in. Want hover? Implement `HoverHandler`. Want completions? Implement `CompletionHandler`. The server figures out the rest.

Some features need extra capability metadata. For example, if your completion handler should run when the user types `.`, or your server exposes commands for the editor to invoke, pass capability options when creating the server:

```go
srv := server.NewServer(h,
    server.WithCompletionOptions(lsp.CompletionOptions{
        TriggerCharacters: []string{"."},
    }),
    server.WithExecuteCommandOptions(lsp.ExecuteCommandOptions{
        Commands: []string{"mylang.generateDebugBundle"},
    }),
)
```

These options enrich auto-detected capabilities. They do not advertise a feature unless your handler implements the matching interface, and explicit capabilities returned from `Initialize` still take precedence.

LSP defaults to UTF-16 positions. If your server deliberately uses another LSP 3.17 position encoding, advertise it with:

```go
srv := server.NewServer(h,
    server.WithPositionEncoding(lsp.PositionEncodingUTF8),
)
```

## Tracking Documents

Most language features need the current text for an open file. Use `document.Store` rather than maintaining a raw `map[lsp.DocumentURI]string`:

```go
type Handler struct {
    documents *document.Store
}

func (h *Handler) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
    _, err := h.documents.Open(params)
    return err
}

func (h *Handler) DidChange(_ context.Context, params *lsp.DidChangeTextDocumentParams) error {
    _, err := h.documents.Change(params)
    return err
}

func (h *Handler) DidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
    h.documents.Close(params)
    return nil
}
```

The store applies full or incremental changes and keeps the document version updated. It also converts between byte offsets and LSP positions:

```go
doc, ok := h.documents.Get(uri)
if !ok {
    return nil
}

offset, err := doc.OffsetAt(lsp.Position{Line: 0, Character: 5})
pos, err := doc.PositionAt(offset)
```

LSP `character` values are UTF-16 code units. This matters for non-ASCII text, especially emoji and other characters represented as surrogate pairs in UTF-16.

## Connect Your Editor

### VS Code

Use a generic LSP client extension and point it at your binary with stdio transport:

```json
{
    "my-lsp.server.path": "./mylang-lsp"
}
```

### Neovim (nvim-lspconfig)

```lua
vim.lsp.start({
    name = "mylang-lsp",
    cmd = { "./mylang-lsp" },
})
```

### Helix

Add to `languages.toml`:

```toml
[[language]]
name = "my-language"
language-servers = ["mylang-lsp"]

[language-server.mylang-lsp]
command = "./mylang-lsp"
```

## Logging

The server logs method dispatch, errors, and lifecycle events via `log/slog`:

```go
srv := server.NewServer(h, server.WithLogger(slog.Default()))
```

Methods and notifications are logged at `Debug` level with their duration. Errors at `Error` level. Lifecycle events at `Info` level. If no logger is set, nothing is logged.

For development, use a text handler with debug level to stderr:

```go
logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
srv := server.NewServer(h, server.WithLogger(logger))
```

## Using the Debug UI

`go-lsp` includes a built-in web UI for inspecting LSP traffic during development:

```go
srv := server.NewServer(h, server.WithDebugUI(":7100"))
```

Open `http://localhost:7100` to see all JSON-RPC messages flowing between client and server.

You can also save the captured session as a JSON trace from your own server code:

```go
err := srv.SaveDebugTrace("/tmp/mylang.trace.json", server.TraceExportOptions{
    RedactDocumentText: true,
    RedactFilePaths:    true,
    Pretty:             true,
})
```

Use this from a custom command, signal handler, or debug endpoint when you need a portable trace for a bug report or regression test.

## Adding More Features

Each LSP feature is an interface. Implement it and the server handles registration and capability advertisement automatically.

| Feature | Interface | Method |
|---------|-----------|--------|
| Completion | `CompletionHandler` | `Completion(ctx, *lsp.CompletionParams) (*lsp.CompletionList, error)` |
| Go to Definition | `DefinitionHandler` | `Definition(ctx, *lsp.DefinitionParams) ([]lsp.Location, error)` |
| Find References | `ReferencesHandler` | `References(ctx, *lsp.ReferenceParams) ([]lsp.Location, error)` |
| Code Actions | `CodeActionHandler` | `CodeAction(ctx, *lsp.CodeActionParams) ([]lsp.CodeAction, error)` |
| Document Symbols | `DocumentSymbolHandler` | `DocumentSymbol(ctx, *lsp.DocumentSymbolParams) ([]lsp.DocumentSymbol, error)` |
| Formatting | `DocumentFormattingHandler` | `Formatting(ctx, *lsp.DocumentFormattingParams) ([]lsp.TextEdit, error)` |
| Rename | `RenameHandler` | `Rename(ctx, *lsp.RenameParams) (*lsp.WorkspaceEdit, error)` |
| Semantic Tokens | `SemanticTokensFullHandler` | `SemanticTokensFull(ctx, *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error)` |
| Inlay Hints | `InlayHintHandler` | `InlayHint(ctx, *lsp.InlayHintParams) ([]lsp.InlayHint, error)` |

See the full list of handler interfaces in [`server/handlers.go`](../server/handlers.go).

## Server-to-Client Communication

The `Client` type provides methods for pushing information to the editor:

```go
// Push diagnostics (errors, warnings)
h.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{...})

// Show a message popup
h.client.ShowMessage(ctx, &lsp.ShowMessageParams{
    Type:    lsp.MessageTypeInfo,
    Message: "Indexing complete",
})

// Log to the output panel
h.client.LogMessage(ctx, &lsp.LogMessageParams{
    Type:    lsp.MessageTypeLog,
    Message: "Processing file...",
})

// Report progress
h.client.CreateWorkDoneProgress(ctx, &lsp.WorkDoneProgressCreateParams{Token: "indexing"})
h.client.Progress(ctx, &lsp.ProgressParams{Token: "indexing", Value: ...})
```

## Custom JSON-RPC Methods

If you need methods outside the LSP spec:

```go
srv := server.NewServer(h)

srv.HandleMethod("custom/myMethod", func(ctx context.Context, params json.RawMessage) (any, error) {
    return map[string]string{"status": "ok"}, nil
})

srv.HandleNotification("custom/myNotification", func(ctx context.Context, params json.RawMessage) error {
    // handle notification
    return nil
})
```

## Testing Your Server

The `servertest` package provides a test harness that simulates an LSP client over in-memory pipes:

```go
func TestHover(t *testing.T) {
    h := servertest.New(t, newHandler())

    h.DidOpen("file:///test.txt", "mylang", "hello world")

    hover, err := h.Hover("file:///test.txt", 0, 0)
    if err != nil {
        t.Fatal(err)
    }
    // assert on hover.Contents
}
```

See the [testing guide](testing.md) for the full API including diagnostics collection, code actions, and more.

## Next Steps

- Browse the [examples](../examples/) for focused feature demos
- Read the [document guide](documents.md) for document sync and UTF-16 position handling
- Check the [LSP specification](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/) for protocol details
- See [`server/handlers.go`](../server/handlers.go) for all available handler interfaces
