# Testing Your Language Server

The `servertest` package provides a test harness that simulates an LSP client over in-memory pipes. It handles JSON-RPC framing, initialization, and cleanup so you can write focused tests for your handler logic.

## Install

The package is included with go-lsp — no extra dependency needed.

```go
import "github.com/owenrumney/go-lsp/servertest"
```

## Basic Usage

Create a harness by passing your handler to `servertest.New`. It starts the server, performs the initialize/initialized handshake, and registers cleanup to shut down gracefully when the test ends.

```go
func TestHover(t *testing.T) {
    h := servertest.New(t, &myHandler{})

    h.DidOpen("file:///test.txt", "plaintext", "hello world")

    hover, err := h.Hover("file:///test.txt", 0, 5)
    if err != nil {
        t.Fatal(err)
    }
    if hover == nil {
        t.Fatal("expected hover result")
    }
}
```

No manual shutdown or cleanup needed — `t.Cleanup` handles it.

## Document Operations

The harness provides shortcuts for document sync notifications:

```go
// Open a document (sets version to 1)
h.DidOpen("file:///main.go", "go", sourceCode)

// Update with new content (full sync)
h.DidChange("file:///main.go", 2, updatedSource)

// Trigger save
h.DidSave("file:///main.go")

// Close
h.DidClose("file:///main.go")
```

## Request Methods

Typed methods for common LSP requests. These construct the params structs for you from minimal arguments:

```go
hover, err := h.Hover(uri, line, char)
list, err := h.Completion(uri, line, char)
item, err := h.ResolveCompletionItem(item)
sig, err := h.SignatureHelp(uri, line, char)
decls, err := h.Declaration(uri, line, char)
locs, err := h.Definition(uri, line, char)
types, err := h.TypeDefinition(uri, line, char)
impls, err := h.Implementation(uri, line, char)
locs, err := h.References(uri, line, char, includeDeclaration)
highlights, err := h.DocumentHighlight(uri, line, char)
syms, err := h.DocumentSymbol(uri)
syms, err := h.WorkspaceSymbol("query")
lenses, err := h.CodeLens(uri)
links, err := h.DocumentLink(uri)
colors, err := h.DocumentColor(uri)
edits, err := h.Formatting(uri)
edits, err := h.RangeFormatting(uri, theRange)
edits, err := h.OnTypeFormatting(uri, line, char, "}")
edit, err := h.Rename(uri, line, char, "newName")
prepare, err := h.PrepareRename(uri, line, char)
folds, err := h.FoldingRange(uri)
selection, err := h.SelectionRange(uri, positions)
linked, err := h.LinkedEditingRange(uri, line, char)
monikers, err := h.Moniker(uri, line, char)
hints, err := h.InlayHint(uri, theRange)
tokens, err := h.SemanticTokensFull(uri)
```

For code actions (which need richer params), pass the full struct:

```go
actions, err := h.CodeAction(&lsp.CodeActionParams{
    TextDocument: lsp.TextDocumentIdentifier{URI: uri},
    Range:        theRange,
    Context:      lsp.CodeActionContext{Diagnostics: diags},
})
```

The harness also has typed helpers for code action, code lens, document link, and inlay hint resolve requests; pull diagnostics; semantic token delta/range requests; call/type hierarchy; file operation requests; and `workspace/executeCommand`.

For anything not covered by a typed method, use the escape hatch:

```go
result, err := h.Call("textDocument/someMethod", params)
err := h.Notify("custom/notification", params)
```

Use `CallAsync` and `CancelRequest` when testing cancellation behavior:

```go
call, err := h.CallAsync("textDocument/hover", params)
if err != nil {
    t.Fatal(err)
}
_ = h.CancelRequest(call.ID())
_, err = call.Wait(ctx)
```

## Testing Diagnostics

Diagnostics arrive as server-to-client notifications, which are asynchronous. The harness collects them automatically. Use `WaitForDiagnostics` to block until they arrive:

```go
func TestDiagnosticsOnSave(t *testing.T) {
    h := servertest.New(t, newHandler())

    h.DidOpen("file:///test.txt", "plaintext", "TODO fix this")
    h.DidSave("file:///test.txt")

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    diags, err := h.WaitForDiagnostics(ctx, "file:///test.txt")
    if err != nil {
        t.Fatal(err)
    }
    if len(diags) != 1 {
        t.Fatalf("expected 1 diagnostic, got %d", len(diags))
    }
}
```

You can also check what's been collected so far without waiting:

```go
diags := h.Diagnostics("file:///test.txt")  // latest for this URI
all := h.AllDiagnostics()                    // all notifications received
h.ClearDiagnostics()                         // reset between test steps
```

## Testing Server Messages

The harness also collects `window/showMessage` and `window/logMessage` notifications:

```go
msgs := h.Messages()      // []lsp.ShowMessageParams
logs := h.LogMessages()   // []lsp.LogMessageParams
```

When the notification is asynchronous, wait for it:

```go
msg, err := h.WaitForMessage(ctx)
log, err := h.WaitForLogMessage(ctx)
```

## Testing Server-to-Client Requests

Handlers that call methods on `server.Client` can be tested without a real editor. Configure the client response, trigger the server behavior, then inspect the captured request:

```go
h.SetClientResponse("window/showMessageRequest", lsp.MessageActionItem{Title: "OK"})

_, err := h.ExecuteCommand("ask-user", nil)
if err != nil {
    t.Fatal(err)
}

req, err := h.WaitForClientRequest(ctx, "window/showMessageRequest")
if err != nil {
    t.Fatal(err)
}
// assert on req.Method or req.Params
```

Use `SetClientError` to force an error response, or `ClientRequests` to inspect everything captured so far.

## Repository Test Commands

The repository Makefile exposes the usual local verification commands:

```bash
make test               # go test -v ./...
make test-race          # go test -race ./...
make test-cover         # go test -cover ./...
make test-fuzz-document # short document fuzz run
```

Fuzz tests live in normal package test files. Ordinary `go test ./...` runs only their seed cases, while `make test-fuzz-document` actively fuzzes document position and edit invariants.

## Checking Capabilities

The initialize result is stored on the harness:

```go
h := servertest.New(t, &myHandler{})
if h.InitResult.Capabilities.HoverProvider == nil {
    t.Fatal("expected hover support")
}
```

## Options

Override the default initialize params or pass server options:

```go
h := servertest.New(t, &myHandler{},
    servertest.WithInitializeParams(&lsp.InitializeParams{
        // custom client capabilities, root URI, etc.
    }),
    servertest.WithServerOptions(
        server.WithLogger(slog.Default()),
    ),
)
```

## Full Example

Testing a handler that flags duplicate keys in .env files:

```go
func TestDuplicateKeys(t *testing.T) {
    handler := newEnvHandler()
    h := servertest.New(t, handler)

    uri := lsp.DocumentURI("file:///test.env")
    h.DidOpen(uri, "env", "FOO=1\nBAR=2\nFOO=3")
    h.DidSave(uri)

    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    diags, err := h.WaitForDiagnostics(ctx, uri)
    if err != nil {
        t.Fatal(err)
    }

    if len(diags) != 1 {
        t.Fatalf("expected 1 duplicate key diagnostic, got %d", len(diags))
    }
    if diags[0].Message == "" {
        t.Fatal("expected diagnostic message")
    }

    // Hover should show the value
    hover, err := h.Hover(uri, 0, 0)
    if err != nil {
        t.Fatal(err)
    }
    if hover == nil || !strings.Contains(hover.Contents.Value, "FOO") {
        t.Fatal("expected hover to show key name")
    }
}
```
