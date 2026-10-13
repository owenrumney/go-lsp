# toylang-env

A toy language server for `.env`-style `KEY=VALUE` config files, demonstrating multiple go-lsp features working together.

## Features

- **Document sync** -- tracks open file contents via `server.WithDocumentStore`
- **Hover** -- shows the value when hovering over a key
- **Completion** -- suggests known keys collected from all open documents
- **Diagnostics** -- flags duplicate keys as you type, via `server.Diagnoser`
- **Go-to-definition** -- jumps to where a key is defined, across all open files

## Build & Run

```
go build -o toylang-env ./examples/toylang/
```

Configure your editor to use `toylang-env` as the language server for `.env` files (or any plain-text file with `KEY=VALUE` lines), communicating over stdio.

## What to expect

1. Open two `.env` files. Type a key in one -- completion will suggest keys from both.
2. Hover over a key to see its value.
3. Add a duplicate key -- a warning diagnostic appears once you stop typing.
4. Use go-to-definition on a key to find all definitions across open files.
