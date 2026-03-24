package server

import "os"

// StdRWC wraps stdin/stdout into a single io.ReadWriteCloser.
type StdRWC struct{}

func (StdRWC) Read(p []byte) (int, error)  { return os.Stdin.Read(p) }
func (StdRWC) Write(p []byte) (int, error) { return os.Stdout.Write(p) }
func (StdRWC) Close() error                { return nil }

// RunStdio returns a StdRWC for use with [Server.Run] to communicate over stdin/stdout.
func RunStdio() StdRWC {
	return StdRWC{}
}
