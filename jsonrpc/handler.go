package jsonrpc

import (
	"context"
	"encoding/json"
)

// MethodHandler handles a JSON-RPC request and returns a result or error.
type MethodHandler func(ctx context.Context, params json.RawMessage) (any, error)

// NotificationHandler handles a JSON-RPC notification.
type NotificationHandler func(ctx context.Context, params json.RawMessage) error

// Dispatcher routes JSON-RPC methods to their handlers.
type Dispatcher struct {
	methods       map[string]MethodHandler
	notifications map[string]NotificationHandler
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{
		methods:       make(map[string]MethodHandler),
		notifications: make(map[string]NotificationHandler),
	}
}

func (d *Dispatcher) RegisterMethod(method string, handler MethodHandler) {
	d.methods[method] = handler
}

func (d *Dispatcher) RegisterNotification(method string, handler NotificationHandler) {
	d.notifications[method] = handler
}

func (d *Dispatcher) HandleRequest(ctx context.Context, req *Request) *Response {
	handler, ok := d.methods[req.Method]
	if !ok {
		return NewErrorResponse(req.ID, NewError(CodeMethodNotFound, "method not found: "+req.Method))
	}

	result, err := handler(ctx, req.Params)
	if err != nil {
		if respErr, ok := err.(*ResponseError); ok {
			return NewErrorResponse(req.ID, respErr)
		}
		return NewErrorResponse(req.ID, NewError(CodeInternalError, err.Error()))
	}

	resp, err := NewResponse(req.ID, result)
	if err != nil {
		return NewErrorResponse(req.ID, NewError(CodeInternalError, "failed to marshal result: "+err.Error()))
	}
	return resp
}

func (d *Dispatcher) HandleNotification(ctx context.Context, notif *Notification) {
	handler, ok := d.notifications[notif.Method]
	if !ok {
		return
	}
	_ = handler(ctx, notif.Params)
}
