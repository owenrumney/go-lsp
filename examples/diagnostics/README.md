# Diagnostics Example

LSP server that publishes diagnostics as you type. Scans for `TODO` markers and reports them as warnings. Demonstrates `server.Diagnoser` with `server.WithDocumentStore`: the server keeps the document store in sync, debounces `Diagnose` after each edit, drops stale results, and clears diagnostics on close.

## Build & Run

```sh
go build -o diagnostics-server ./examples/diagnostics/
```

Configure your editor to use `./diagnostics-server` as an LSP server over stdio.

## What to Expect

1. Open a file. Warning diagnostics appear on each existing `TODO` marker.
2. Type a line containing `TODO`. The warning appears once you stop typing.
3. Close the file. The diagnostics are cleared.
