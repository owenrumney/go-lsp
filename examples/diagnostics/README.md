# Diagnostics Example

LSP server that publishes diagnostics on save. Scans for `TODO` markers and reports them as warnings. Demonstrates `TextDocumentSaveHandler` and server-to-client communication via `Client.PublishDiagnostics`.

## Build & Run

```sh
go build -o diagnostics-server ./examples/diagnostics/
```

Configure your editor to use `./diagnostics-server` as an LSP server over stdio.

## What to Expect

1. Open a file and add a line containing `TODO`.
2. Save the file.
3. Warning diagnostics appear on each `TODO` marker.
