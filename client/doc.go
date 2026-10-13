// Package client drives a language server over LSP.
//
// It is the counterpart to package server: where server turns a handler into a
// language server, client turns a language server into something a Go program
// can query. Dial spawns the server, performs the initialize handshake, and
// returns a Client whose methods map onto LSP requests.
//
// Client owns no document text. Callers hold their own buffers and pass
// explicit versions to DidOpen and DidChangeFull, because a client that owns
// text must also own version numbering, and that belongs to the application.
//
// The hard part of driving a server is knowing when to trust an empty answer.
// A server accepts initialize long before it has finished indexing, and
// requests answered in that window return empty results indistinguishable
// from real ones. Ready waits for startup work to finish; see its
// documentation for what that can and cannot promise.
package client
