package server

import (
	"context"
	"encoding/json"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
)

// MethodHandlerFunc is a raw JSON-RPC request handler used by middleware.
type MethodHandlerFunc func(context.Context, json.RawMessage) (any, error)

// NotificationHandlerFunc is a raw JSON-RPC notification handler used by middleware.
type NotificationHandlerFunc func(context.Context, json.RawMessage) error

// MethodMiddleware wraps inbound JSON-RPC request handling for a method.
// Middleware runs inside the LSP lifecycle gate: requests rejected for
// lifecycle reasons (e.g. before initialize) never reach it.
type MethodMiddleware func(method string, next MethodHandlerFunc) MethodHandlerFunc

// NotificationMiddleware wraps inbound JSON-RPC notification handling for a method.
// Middleware runs inside the LSP lifecycle gate: notifications dropped for
// lifecycle reasons never reach it.
type NotificationMiddleware func(method string, next NotificationHandlerFunc) NotificationHandlerFunc

func (s *Server) wrapMethod(method string, handler jsonrpc.MethodHandler) jsonrpc.MethodHandler {
	wrapped := MethodHandlerFunc(handler)
	for i := len(s.methodMiddleware) - 1; i >= 0; i-- {
		wrapped = s.methodMiddleware[i](method, wrapped)
	}
	withLifecycle := func(ctx context.Context, params json.RawMessage) (any, error) {
		if respErr := s.checkRequestLifecycle(method); respErr != nil {
			return nil, respErr
		}
		return wrapped(ctx, params)
	}
	return s.logMethod(method, withLifecycle)
}

func (s *Server) wrapNotification(method string, handler jsonrpc.NotificationHandler) jsonrpc.NotificationHandler {
	wrapped := NotificationHandlerFunc(handler)
	for i := len(s.notificationMiddleware) - 1; i >= 0; i-- {
		wrapped = s.notificationMiddleware[i](method, wrapped)
	}
	withLifecycle := func(ctx context.Context, params json.RawMessage) error {
		if !s.allowNotification(method) {
			return nil
		}
		return wrapped(ctx, params)
	}
	return s.logNotification(method, withLifecycle)
}
