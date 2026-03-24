# Workspace Symbols Example

Demonstrates document and workspace symbol support using go-lsp. The server tracks open documents and parses simple headers to expose them as symbols.

Recognized patterns:
- `func ` / `def ` lines -> Function symbols
- `class ` lines -> Class symbols
- `## ` markdown headers -> String symbols

## Build & Run

```sh
go build -o symbols-server ./examples/symbols/
```

Connect the binary to an editor as an LSP server over stdio. Use "Go to Symbol in Workspace" (e.g. `Ctrl+T` in VS Code) to search across all open files, or "Go to Symbol in File" (`Ctrl+Shift+O`) for the current document.
