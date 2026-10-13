package servertest

import (
	"context"
	"encoding/json"
	"io"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
)

// PendingCall is an in-flight JSON-RPC request sent by a harness.
type PendingCall struct {
	pending *jsonrpc.Pending
}

// rpcConn adapts jsonrpc.Conn to the harness, which works in raw JSON results
// rather than Response values.
type rpcConn struct {
	conn *jsonrpc.Conn
}

func newRPCConn(rw io.ReadWriteCloser, d *jsonrpc.Dispatcher) *rpcConn {
	return &rpcConn{conn: jsonrpc.NewConn(rw, d)}
}

// ID returns the JSON-RPC request ID, for use with Harness.CancelRequest.
func (c *PendingCall) ID() int64 {
	id, _ := c.pending.ID().Int()
	return id
}

// Wait blocks until the response arrives, ctx is done, or the server exits.
func (c *PendingCall) Wait(ctx context.Context) (json.RawMessage, error) {
	resp, err := c.pending.Wait(ctx)
	if err != nil {
		return nil, err
	}
	return rawResult(resp)
}

func (c *rpcConn) serve(ctx context.Context) error {
	return c.conn.Serve(ctx)
}

func (c *rpcConn) call(ctx context.Context, method string, params any) (json.RawMessage, error) {
	resp, err := c.conn.Call(ctx, method, params)
	if err != nil {
		return nil, err
	}
	return rawResult(resp)
}

func (c *rpcConn) startCall(method string, params any) (*PendingCall, error) {
	p, err := c.conn.CallAsync(method, params)
	if err != nil {
		return nil, err
	}
	return &PendingCall{pending: p}, nil
}

func (c *rpcConn) notify(ctx context.Context, method string, params any) error {
	return c.conn.Notify(ctx, method, params)
}

func rawResult(resp *jsonrpc.Response) (json.RawMessage, error) {
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}
