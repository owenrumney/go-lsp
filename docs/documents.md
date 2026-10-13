# Managing Documents

Language servers usually need the current text for each open file. The `document` package provides a small store for that state.

```go
type Handler struct {
    documents *document.Store
}

func NewHandler() *Handler {
    return &Handler{documents: document.NewStore()}
}
```

## Document Sync

Forward text document notifications directly to the store:

```go
func (h *Handler) DidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
    _, err := h.documents.Open(params)
    return err
}

func (h *Handler) DidChange(_ context.Context, params *lsp.DidChangeTextDocumentParams) error {
    _, err := h.documents.Change(params)
    return err
}

func (h *Handler) DidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
    h.documents.Close(params)
    return nil
}
```

If your handler implements `TextDocumentSyncHandler`, `go-lsp` advertises incremental sync by default. `document.Store` also accepts full-document changes, so it works with clients that send full sync.

### Letting The Server Sync The Store

`server.WithDocumentStore` removes the three forwarding methods above. The server updates the store on didOpen, didChange, and didClose, and advertises incremental sync even if your handler has no sync methods:

```go
store := document.NewStore()
srv := server.NewServer(h, server.WithDocumentStore(store))
```

Your own `DidOpen`, `DidChange`, and `DidClose` still run if you implement them. The store is updated first, so they see the new text.

A handler that also implements `server.Diagnoser` gets diagnostics run for it:

```go
func (h *Handler) Diagnose(ctx context.Context, doc *document.Document) ([]lsp.Diagnostic, error) {
    return check(doc.Text()), nil
}
```

The server runs `Diagnose` immediately on open and `WithDiagnosticsDebounce` after the last edit (300ms by default). `ctx` is cancelled when a newer edit arrives, the document closes, or the server shuts down, and a result for a version that has since changed is dropped rather than published. Closing the document publishes an empty list. `Diagnoser` requires `WithDocumentStore`; `NewServer` panics without it.

## Reading Text

Use `Text` when you just need the raw string:

```go
text, ok := h.documents.Text(params.TextDocument.URI)
if !ok {
    return nil
}
```

Use `Get` when you need metadata or position conversion:

```go
doc, ok := h.documents.Get(params.TextDocument.URI)
if !ok {
    return nil
}

text := doc.Text()
version := doc.Version()
```

`Get` returns a snapshot, so callers cannot mutate the store's internal state.

## Positions And Offsets

LSP positions use zero-based lines and UTF-16 character offsets. Go strings use byte offsets, and `range` iterates runes. These are not equivalent for all text.

Use `OffsetAt` to convert an LSP position to a byte offset:

```go
offset, err := doc.OffsetAt(params.Position)
if err != nil {
    return nil, err
}
```

Use `PositionAt` to convert a byte offset back to an LSP position:

```go
pos, err := doc.PositionAt(offset)
if err != nil {
    return nil, err
}
```

This matters for characters outside the BMP, such as emoji. In LSP, `"😀"` has character length 2 because it is represented by a UTF-16 surrogate pair.

## Errors

The package exposes sentinel errors for common failure cases:

- `document.ErrDocumentNotFound`
- `document.ErrInvalidPosition`
- `document.ErrInvalidRange`
- `document.ErrVersionRegression`

Use `errors.Is` when matching them:

```go
if errors.Is(err, document.ErrInvalidRange) {
    return nil, err
}
```
