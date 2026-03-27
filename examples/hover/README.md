# Hover Example

Minimal LSP server that returns hover information for every symbol. This is the simplest possible go-lsp example — a good starting point.

## Build & Run

```sh
go build -o hover-server ./examples/hover/
```

Configure your editor to use `./hover-server` as an LSP server over stdio.

## What to Expect

Hover over any token in an open file to see a markdown tooltip from the server.
