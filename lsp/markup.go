package lsp

// MarkupKind describes the content type of a markup string.
type MarkupKind string

const (
	PlainText MarkupKind = "plaintext"
	Markdown  MarkupKind = "markdown"
)

// MarkupContent carries documentation or descriptive text tagged as either plaintext or markdown.
type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}
