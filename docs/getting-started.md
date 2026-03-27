# Getting Started with go-lsp

This guide walks you through building a language server from scratch using `go-lsp`. By the end, you'll have a working LSP server that tracks open documents, publishes diagnostics on save, and provides hover information.

## Prerequisites

- Go 1.25+
- An editor that supports LSP (VS Code, Neovim, etc.)

## Quick Start with Scaffold

The fastest way to get started is to generate a project:

```bash
go run github.com/owenrumney/go-lsp/cmd/scaffold@latest
```

This walks you through a few prompts and generates a working LSP server with the features you pick — including a passing test. You can also pass flags for scripting:

```bash
go run github.com/owenrumney/go-lsp/cmd/scaffold@latest \
  --name mylang \
  --module github.com/you/mylang-lsp \
  --lang mylang \
  --features hover,completion,diagnostics
```

Available features: `hover`, `completion`, `diagnostics`, `definition`, `references`, `formatting`, `codeactions`, `symbols`.

The generated project includes:
- `main.go` — server entrypoint wired to stdio
- `handler/handler.go` — handler struct with the interfaces you selected
- `handler/handler_test.go` — passing tests using the `servertest` harness
- `go.mod` with the go-lsp dependency

Build it, point your editor at the binary, and you're running.

If you'd rather build from scratch to understand how things work, read on.

## Install

```bash
go get github.com/owenrumney/go-lsp@latest
```

## How go-lsp Works

The core idea is simple: you write a handler struct, implement the interfaces for the LSP features you want, and `go-lsp` does the rest. It auto-detects which interfaces your handler implements, registers the corresponding JSON-RPC methods, and advertises the right capabilities to the client.

```
Your handler struct
    ↓ implements interfaces
go-lsp auto-detects capabilities
    ↓
JSON-RPC dispatch over stdio/TCP/WebSocket
```

The only required interface is `LifecycleHandler`:

```go
type LifecycleHandler interface {
    Initialize(ctx context.Context, params *lsp.InitializeParams) (*lsp.InitializeResult, error)
    Shutdown(ctx context.Context) error
}
```

Everything else is opt-in. Want hover? Implement `HoverHandler`. Want completions? Implement `CompletionHandler`. The server figures out the rest.

## Step 1: Minimal Server

Start with the smallest possible server — one that initializes and shuts down:

```go
package main

import (
    "context"
    "log"

    "github.com/owenrumney/go-lsp/lsp"
    "github.com/owenrumney/go-lsp/server"
)

type handler struct{}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
    return &lsp.InitializeResult{
        ServerInfo: &lsp.ServerInfo{Name: "my-lsp", Version: "0.1.0"},
    }, nil
}

func (h *handler) Shutdown(_ context.Context) error {
    return nil
}

func main() {
    srv := server.NewServer(&handler{})
    if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
        log.Fatal(err)
    }
}
```

That's it. This is a valid LSP server. It doesn't do anything useful yet, but it will respond to `initialize` and `shutdown` correctly.

## Step 2: Track Open Documents

Most language servers need to know what files are open and what they contain. Implement `TextDocumentSyncHandler` to get notified when documents are opened, changed, and closed:

```go
type handler struct {
    documents map[lsp.DocumentURI]string
}

func newHandler() *handler {
    return &handler{documents: make(map[lsp.DocumentURI]string)}
}

func (h *handler) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
    h.documents[params.TextDocument.URI] = params.TextDocument.Text
    return nil
}

func (h *handler) DidChange(_ context.Context, params *lsp.DidChangeTextDocumentParams) error {
    if len(params.ContentChanges) > 0 {
        h.documents[params.TextDocument.URI] = params.ContentChanges[len(params.ContentChanges)-1].Text
    }
    return nil
}

func (h *handler) DidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
    delete(h.documents, params.TextDocument.URI)
    return nil
}
```

Because you implemented `TextDocumentSyncHandler`, `go-lsp` will automatically advertise `textDocumentSync` capabilities to the client. You don't need to set them manually (though you can override them in `Initialize` if you want fine-grained control).

## Step 3: Publish Diagnostics on Save

Now let's do something useful. We'll scan saved files for lines containing "TODO" and report them as warnings. This requires two things:

1. Implement `TextDocumentSaveHandler` for the save notification
2. Use the `Client` to push diagnostics back to the editor

To get access to the `Client`, implement the `ClientHandler` interface:

```go
type handler struct {
    documents map[lsp.DocumentURI]string
    client    *server.Client
}

func (h *handler) SetClient(client *server.Client) {
    h.client = client
}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
    return &lsp.InitializeResult{
        Capabilities: lsp.ServerCapabilities{
            TextDocumentSync: &lsp.TextDocumentSyncOptions{
                OpenClose: boolPtr(true),
                Change:    lsp.SyncFull,
                Save:      &lsp.SaveOptions{IncludeText: boolPtr(true)},
            },
        },
        ServerInfo: &lsp.ServerInfo{Name: "my-lsp", Version: "0.1.0"},
    }, nil
}

func (h *handler) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error {
    var diags []lsp.Diagnostic

    if params.Text != nil {
        for i, line := range strings.Split(*params.Text, "\n") {
            if idx := strings.Index(line, "TODO"); idx >= 0 {
                sev := lsp.SeverityWarning
                diags = append(diags, lsp.Diagnostic{
                    Range: lsp.Range{
                        Start: lsp.Position{Line: i, Character: idx},
                        End:   lsp.Position{Line: i, Character: idx + 4},
                    },
                    Severity: &sev,
                    Source:   "my-lsp",
                    Message:  "TODO found",
                })
            }
        }
    }

    return h.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
        URI:         params.TextDocument.URI,
        Diagnostics: diags,
    })
}

func boolPtr(b bool) *bool { return &b }
```

## Step 4: Add Hover Support

Let's add hover information. When the user hovers over any position, we'll show them the line number and column. In a real server, you'd look up symbol information here.

Just implement `HoverHandler`:

```go
func (h *handler) Hover(_ context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
    content, ok := h.documents[params.TextDocument.URI]
    if !ok {
        return nil, nil // document not tracked
    }

    lines := strings.Split(content, "\n")
    line := params.Position.Line
    if line >= len(lines) {
        return nil, nil
    }

    return &lsp.Hover{
        Contents: lsp.MarkupContent{
            Kind:  lsp.Markdown,
            Value: fmt.Sprintf("**Line %d, Col %d**\n\n```\n%s\n```", line+1, params.Position.Character+1, lines[line]),
        },
    }, nil
}
```

## Step 5: Put It All Together

Here's the complete server:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/owenrumney/go-lsp/lsp"
    "github.com/owenrumney/go-lsp/server"
)

type handler struct {
    documents map[lsp.DocumentURI]string
    client    *server.Client
}

func newHandler() *handler {
    return &handler{documents: make(map[lsp.DocumentURI]string)}
}

func (h *handler) SetClient(client *server.Client) {
    h.client = client
}

func (h *handler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
    return &lsp.InitializeResult{
        Capabilities: lsp.ServerCapabilities{
            TextDocumentSync: &lsp.TextDocumentSyncOptions{
                OpenClose: boolPtr(true),
                Change:    lsp.SyncFull,
                Save:      &lsp.SaveOptions{IncludeText: boolPtr(true)},
            },
        },
        ServerInfo: &lsp.ServerInfo{Name: "my-lsp", Version: "0.1.0"},
    }, nil
}

func (h *handler) Shutdown(_ context.Context) error { return nil }

func (h *handler) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
    h.documents[params.TextDocument.URI] = params.TextDocument.Text
    return nil
}

func (h *handler) DidChange(_ context.Context, params *lsp.DidChangeTextDocumentParams) error {
    if len(params.ContentChanges) > 0 {
        h.documents[params.TextDocument.URI] = params.ContentChanges[len(params.ContentChanges)-1].Text
    }
    return nil
}

func (h *handler) DidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
    delete(h.documents, params.TextDocument.URI)
    return nil
}

func (h *handler) DidSave(ctx context.Context, params *lsp.DidSaveTextDocumentParams) error {
    var diags []lsp.Diagnostic

    if params.Text != nil {
        for i, line := range strings.Split(*params.Text, "\n") {
            if idx := strings.Index(line, "TODO"); idx >= 0 {
                sev := lsp.SeverityWarning
                diags = append(diags, lsp.Diagnostic{
                    Range: lsp.Range{
                        Start: lsp.Position{Line: i, Character: idx},
                        End:   lsp.Position{Line: i, Character: idx + 4},
                    },
                    Severity: &sev,
                    Source:   "my-lsp",
                    Message:  "TODO found",
                })
            }
        }
    }

    return h.client.PublishDiagnostics(ctx, &lsp.PublishDiagnosticsParams{
        URI:         params.TextDocument.URI,
        Diagnostics: diags,
    })
}

func (h *handler) Hover(_ context.Context, params *lsp.HoverParams) (*lsp.Hover, error) {
    content, ok := h.documents[params.TextDocument.URI]
    if !ok {
        return nil, nil
    }

    lines := strings.Split(content, "\n")
    line := params.Position.Line
    if line >= len(lines) {
        return nil, nil
    }

    return &lsp.Hover{
        Contents: lsp.MarkupContent{
            Kind:  lsp.Markdown,
            Value: fmt.Sprintf("**Line %d, Col %d**\n\n```\n%s\n```", line+1, params.Position.Character+1, lines[line]),
        },
    }, nil
}

func boolPtr(b bool) *bool { return &b }

func main() {
    h := newHandler()
    srv := server.NewServer(h)
    if err := srv.Run(context.Background(), server.RunStdio()); err != nil {
        log.Fatal(err)
    }
}
```

Build it:

```bash
go build -o my-lsp .
```

## Step 6: Connect Your Editor

### VS Code

Create a minimal extension or use a generic LSP client extension. Point it at your binary with stdio transport:

```json
{
    "my-lsp.server.path": "./my-lsp"
}
```

### Neovim (nvim-lspconfig)

```lua
vim.lsp.start({
    name = "my-lsp",
    cmd = { "./my-lsp" },
})
```

### Helix

Add to `languages.toml`:

```toml
[[language]]
name = "my-language"
language-servers = ["my-lsp"]

[language-server.my-lsp]
command = "./my-lsp"
```

## Logging

The server can log method dispatch, errors, and lifecycle events via `log/slog`:

```go
srv := server.NewServer(h, server.WithLogger(slog.Default()))
```

Methods and notifications are logged at `Debug` level with their duration. Errors are logged at `Error` level. Lifecycle events (init, shutdown) are logged at `Info` level. If no logger is set, nothing is logged.

For development, use a text handler with debug level to stderr:

```go
logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
srv := server.NewServer(h, server.WithLogger(logger))
```

## Using the Debug UI

`go-lsp` includes a built-in web UI for inspecting LSP traffic during development. Enable it with a single option:

```go
srv := server.NewServer(h, server.WithDebugUI(":7100"))
```

Then open `http://localhost:7100` in your browser to see all JSON-RPC messages flowing between client and server, along with the capabilities your server advertised.

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

    h.DidOpen("file:///test.env", "env", "FOO=bar")

    hover, err := h.Hover("file:///test.env", 0, 0)
    if err != nil {
        t.Fatal(err)
    }
    // assert on hover.Contents
}
```

See the [testing guide](testing.md) for the full API including diagnostics collection, code actions, and more.

## Next Steps

- Browse the [examples](../examples/) for focused feature demos
- Check the [LSP specification](https://microsoft.github.io/language-server-protocol/specifications/lsp/3.17/specification/) for protocol details
- See [`server/handlers.go`](../server/handlers.go) for all available handler interfaces
