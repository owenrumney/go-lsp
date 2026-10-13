# LSP client

Drives another language server and prints what it knows about a file. This is
the `client` package rather than `server`: the program here is the editor, not
the language server.

```bash
go install golang.org/x/tools/gopls@latest
go run ./examples/lspclient -server gopls -file ./client/ready.go -line 22 -char 7
```

```
connected to gopls
hover at 22:7: func newReadiness(grace time.Duration) *readiness
symbol readiness (kind 23)
symbol newReadiness (kind 12)
```

Three things it demonstrates that are easy to get wrong:

- **`Ready` is called before any query.** A server answers requests while it is
  still indexing, and those answers are empty but indistinguishable from real
  ones.
- **Diagnostics are collected in the `Handler`.** They arrive as notifications,
  so there is no response to wait for.
- **`Configuration` is answered.** gopls asks for settings during startup and
  the default handler would decline.
