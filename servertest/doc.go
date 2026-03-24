// Package servertest provides a test harness for language servers built with
// the [github.com/owenrumney/go-lsp/server] package.
//
// The [Harness] simulates an LSP client over in-memory pipes, handling JSON-RPC
// framing, initialization, and cleanup automatically. Use it to write unit tests
// for handler logic without needing a real editor connection.
//
//	func TestHover(t *testing.T) {
//	    h := servertest.New(t, &myHandler{})
//	    h.DidOpen("file:///test.txt", "plaintext", "hello world")
//	    hover, err := h.Hover("file:///test.txt", 0, 5)
//	    // ...
//	}
package servertest
