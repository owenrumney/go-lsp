# Driving a Language Server

Package `client` is the mirror of package `server`. Where `server` turns your handler into a language server, `client` turns a language server into something your program can ask questions of.

Use it to run `gopls` in CI, build a code-aware tool, or drive a server from an editor or agent.

## Connect

```go
c, err := client.Dial(ctx, client.Config{
    Command: []string{"gopls"},
    Dir:     root,
    RootURI: lsp.DocumentURI("file://" + root),
    Stderr:  os.Stderr,
})
if err != nil {
    log.Fatal(err)
}
defer c.Close()
```

The context passed to `Dial` bounds startup only. The connection lives until `Close`, so passing a request-scoped context will not kill the server underneath you.

## Wait before you trust an empty result

A server accepts `initialize` long before it has finished indexing. Ask for references in that window and you get an empty list that looks exactly like a real answer. A human sees a spinner and waits. A program records it as fact and reasons from it.

```go
if err := c.Ready(ctx); err != nil {
    return err
}

refs, err := c.References(ctx, uri, pos, true)
```

`Ready` is a heuristic, because LSP has no readiness signal. It tracks outstanding work-done progress tokens and, for servers that report none, waits out `Config.ReadyGrace` (2 seconds by default). Requests are never blocked on it, so call it when a wrong answer would cost more than the wait. For a server with a better signal, skip `Ready` and poll for that signal yourself.

## Query

```go
if err := c.DidOpen(ctx, uri, "go", 1, text); err != nil {
    return err
}

hover, err := c.Hover(ctx, uri, lsp.Position{Line: 10, Character: 4})
defs, err := c.Definition(ctx, uri, pos)
syms, err := c.DocumentSymbol(ctx, uri)
edits, err := c.Formatting(ctx, uri, lsp.FormattingOptions{TabSize: 4})
```

The client owns no document text. You hold the buffers and pass explicit versions, because whoever owns the text must own the version numbering.

`DidChangeFull` replaces a whole document, which is what a program usually wants. Incremental sync exists for human keystrokes; if you need it, build the params yourself and send them with `Notify`.

## Ask what the server can do

`Supports` considers both the capabilities from `initialize` and anything the server registered afterwards, matching document selectors for the dynamic ones.

```go
if c.Supports("textDocument/rename", uri, "go") {
    edit, err := c.Rename(ctx, uri, pos, "NewName")
}
```

Read provider fields through `Enabled()` rather than a nil check. A server sending `"codeActionProvider": false` decodes to a non-nil value whose `Enabled()` is false, so `!= nil` reports support that is not there.

## Answer the server

Servers ask their client for things: configuration, permission to edit files, messages to show. Implement `Handler` to answer them. Embed `NopHandler` and override only what you need; it declines everything by default.

```go
type handler struct {
    client.NopHandler
}

func (handler) Configuration(_ context.Context, items []lsp.ConfigurationItem) ([]json.RawMessage, error) {
    out := make([]json.RawMessage, len(items))
    for i := range out {
        out[i] = json.RawMessage(`{"staticcheck": true}`)
    }
    return out, nil
}

func (handler) Diagnostics(params *lsp.PublishDiagnosticsParams) {
    for _, d := range params.Diagnostics {
        log.Printf("%s:%d %s", params.URI, d.Range.Start.Line+1, d.Message)
    }
}
```

Diagnostics arrive as notifications, not as a response, so there is nothing to wait on. Collect them in your `Handler` and match them against the document version you sent. Where the server supports pull diagnostics, `DocumentDiagnostic` is the more direct route.

Capability registration is handled for you. `Client` answers `client/registerCapability`, `client/unregisterCapability`, progress creation and the refresh requests itself, so `Supports` stays accurate. Those never reach your `Handler`.

## Test without a real server

`Config.Transport` replaces the child process with any stream, so a test can run a `server.Server` in the same process over a pipe.

```go
cfg := client.Config{
    Transport: func(ctx context.Context) (io.ReadWriteCloser, error) {
        clientConn, serverConn := net.Pipe()
        srv := server.NewServer(&myHandler{})
        go func() { _ = srv.Run(ctx, serverConn) }()
        return clientConn, nil
    },
}
```

Test against a real server as well. Tests against your own server pass while real servers break you, because real servers are the ones with indexing delays, dynamic registration, and both arms of every union type in the spec.

## When the server dies

```go
select {
case <-c.Done():
    log.Printf("server exited: %v", c.Err())
case <-ctx.Done():
}
```

Requests still in flight fail as soon as the connection drops rather than waiting out their timeout. `Close` shuts the server down gracefully and kills it if it has not exited within `Config.ShutdownTimeout`.
