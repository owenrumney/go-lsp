# Code Actions Example

Demonstrates diagnostics combined with code actions — the most common real-world LSP pattern.

On save, the server flags lines with trailing whitespace as warnings. Each diagnostic has an associated quick-fix code action that trims the trailing whitespace.

## Build & Run

```sh
go build -o codeactions-server ./examples/codeactions/
```

Configure your editor to use `./codeactions-server` as an LSP server over stdio.

## What to Expect

1. Open any file and add trailing spaces to a line.
2. Save the file.
3. Warning diagnostics appear on lines with trailing whitespace.
4. Place your cursor on a warning and invoke "Quick Fix" — the trailing whitespace is removed.
