package server

import "os"

// stdRWC wraps stdin/stdout into a single io.ReadWriteCloser.
type stdRWC struct{}

func (stdRWC) Read(p []byte) (int, error)  { return os.Stdin.Read(p) }
func (stdRWC) Write(p []byte) (int, error) { return os.Stdout.Write(p) }
func (stdRWC) Close() error                { return nil }

func RunStdio() stdRWC {
	return stdRWC{}
}
