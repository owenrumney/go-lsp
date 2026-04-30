package servertest

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// ClientRequest is a server-to-client request captured by the harness.
type ClientRequest struct {
	Method string
	Params json.RawMessage
}

type clientResponse struct {
	result any
	err    error
}

type clientRequestStore struct {
	mu        sync.Mutex
	cond      *sync.Cond
	requests  []ClientRequest
	responses map[string]clientResponse
}

func newClientRequestStore() *clientRequestStore {
	s := &clientRequestStore{
		responses: make(map[string]clientResponse),
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

func (s *clientRequestStore) handle(method string, params json.RawMessage) (any, error) {
	s.mu.Lock()
	s.requests = append(s.requests, ClientRequest{Method: method, Params: cloneRawMessage(params)})
	resp := s.responses[method]
	s.cond.Broadcast()
	s.mu.Unlock()
	return resp.result, resp.err
}

func (s *clientRequestStore) setResponse(method string, result any, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.responses[method] = clientResponse{result: result, err: err}
}

func (s *clientRequestStore) all() []ClientRequest {
	s.mu.Lock()
	defer s.mu.Unlock()
	result := make([]ClientRequest, len(s.requests))
	copy(result, s.requests)
	for i := range result {
		result[i].Params = cloneRawMessage(result[i].Params)
	}
	return result
}

func (s *clientRequestStore) wait(ctx context.Context, method string) (ClientRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for {
		for _, req := range s.requests {
			if method == "" || req.Method == method {
				req.Params = cloneRawMessage(req.Params)
				return req, nil
			}
		}

		if err := waitCond(ctx, s.cond); err != nil {
			return ClientRequest{}, err
		}
	}
}

// SetClientResponse configures the response returned for a server-to-client request method.
func (h *Harness) SetClientResponse(method string, result any) {
	h.clientRequests.setResponse(method, result, nil)
}

// SetClientError configures an error response for a server-to-client request method.
func (h *Harness) SetClientError(method string, err error) {
	if err == nil {
		err = fmt.Errorf("client request failed")
	}
	h.clientRequests.setResponse(method, nil, err)
}

// ClientRequests returns all server-to-client requests captured so far.
func (h *Harness) ClientRequests() []ClientRequest {
	return h.clientRequests.all()
}

// WaitForClientRequest waits until the server sends a matching server-to-client request.
// Pass an empty method to wait for any request.
func (h *Harness) WaitForClientRequest(ctx context.Context, method string) (ClientRequest, error) {
	return h.clientRequests.wait(ctx, method)
}

func cloneRawMessage(raw json.RawMessage) json.RawMessage {
	if raw == nil {
		return nil
	}
	clone := make(json.RawMessage, len(raw))
	copy(clone, raw)
	return clone
}
