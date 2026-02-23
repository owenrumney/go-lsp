package lsp

// MarkupKind describes the content type of a markup string.
type MarkupKind string

const (
	PlainText MarkupKind = "plaintext"
	Markdown  MarkupKind = "markdown"
)

// MarkupContent represents a string value with a specific content type.
type MarkupContent struct {
	Kind  MarkupKind `json:"kind"`
	Value string     `json:"value"`
}
