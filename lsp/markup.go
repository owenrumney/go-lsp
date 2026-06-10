package lsp

// MarkupKind describes the content type of a markup string.
type MarkupKind string

const (
	// Plain text is supported as a content format.
	PlainText MarkupKind = "plaintext"
	// Markdown is supported as a content format.
	Markdown MarkupKind = "markdown"
)

// MarkupContent is a literal that represents a string value whose content is interpreted based on its
// kind flag. Currently the protocol supports plaintext and markdown as markup kinds.
//
// If the kind is markdown then the value can contain fenced code blocks like in GitHub issues.
// See https://help.github.com/articles/creating-and-highlighting-code-blocks/#syntax-highlighting
//
// Here is an example how such a string can be constructed using JavaScript / TypeScript:
//
//	let markdown: MarkdownContent = {
//	 kind: "markdown",
//	 value: [
//	   '# Header',
//	   'Some text',
//	   '```typescript',
//	   'someCode();',
//	   '```'
//	 ].join('\n')
//	};
//
// *Please Note* that clients might sanitize the return markdown. A client could decide to
// remove HTML from the markdown to avoid script execution.
type MarkupContent struct {
	// The type of the Markup
	Kind MarkupKind `json:"kind"`
	// The content itself
	Value string `json:"value"`
}
