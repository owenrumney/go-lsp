// Package server provides a framework for building Language Server Protocol (LSP) 3.17 servers in Go.
//
// Create a server by implementing [LifecycleHandler] and any additional handler interfaces
// for the LSP features you need, then pass your handler to [NewServer]:
//
//	srv := server.NewServer(myHandler)
//	srv.Run(ctx, server.RunStdio())
//
// The server auto-detects which handler interfaces your struct implements, registers the
// corresponding JSON-RPC methods, and advertises the right capabilities to the client.
//
// For server-to-client communication (diagnostics, messages, progress), implement
// [ClientHandler] to receive a [Client] after the connection is established.
//
// See the handler interfaces in handlers.go for the full list of supported LSP features.
package server
