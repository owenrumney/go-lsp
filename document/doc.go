// Package document manages open LSP text documents.
//
// Store tracks open documents and applies full or incremental
// textDocument/didChange updates. Positions are interpreted using LSP's UTF-16
// character offsets, not Go byte offsets or rune indexes.
package document
