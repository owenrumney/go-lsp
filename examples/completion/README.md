# Completion Example

LSP server that provides keyword completions with resolve support. Demonstrates `CompletionHandler` and `CompletionResolveHandler` interfaces.

## Build & Run

```sh
go build -o completion-server ./examples/completion/
```

Configure your editor to use `./completion-server` as an LSP server over stdio.

## What to Expect

1. Start typing in any file — the server suggests Go keywords (`func`, `var`, `const`, `type`).
2. Select a completion item — the server resolves it with a markdown documentation string.
