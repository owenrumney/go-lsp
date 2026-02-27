package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

// Client provides methods for server-to-client communication.
type Client struct {
	conn *jsonrpc.Conn
}

func newClient(conn *jsonrpc.Conn) *Client {
	return &Client{conn: conn}
}

func (c *Client) PublishDiagnostics(ctx context.Context, params *lsp.PublishDiagnosticsParams) error {
	return c.conn.Notify(ctx, "textDocument/publishDiagnostics", params)
}

func (c *Client) ShowMessage(ctx context.Context, params *lsp.ShowMessageParams) error {
	return c.conn.Notify(ctx, "window/showMessage", params)
}

func (c *Client) LogMessage(ctx context.Context, params *lsp.LogMessageParams) error {
	return c.conn.Notify(ctx, "window/logMessage", params)
}

func (c *Client) Progress(ctx context.Context, params *lsp.ProgressParams) error {
	return c.conn.Notify(ctx, "$/progress", params)
}

// ShowMessageRequest sends a window/showMessageRequest to the client and waits for a response.
func (c *Client) ShowMessageRequest(ctx context.Context, params *lsp.ShowMessageRequestParams) (*lsp.MessageActionItem, error) {
	resp, err := c.conn.Call(ctx, "window/showMessageRequest", params)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("showMessageRequest: %s", resp.Error.Message)
	}
	if resp.Result == nil || string(resp.Result) == "null" {
		return nil, nil
	}
	var item lsp.MessageActionItem
	if err := json.Unmarshal(resp.Result, &item); err != nil {
		return nil, err
	}
	return &item, nil
}

// CreateWorkDoneProgress sends a window/workDoneProgress/create request to the client.
// This must be called before sending $/progress notifications with the same token.
func (c *Client) CreateWorkDoneProgress(ctx context.Context, params *lsp.WorkDoneProgressCreateParams) error {
	_, err := c.conn.Call(ctx, "window/workDoneProgress/create", params)
	return err
}

// InlayHintRefresh sends a workspace/inlayHint/refresh request to the client.
func (c *Client) InlayHintRefresh(ctx context.Context) error {
	_, err := c.conn.Call(ctx, "workspace/inlayHint/refresh", nil)
	return err
}

// InlineValueRefresh sends a workspace/inlineValue/refresh request to the client.
func (c *Client) InlineValueRefresh(ctx context.Context) error {
	_, err := c.conn.Call(ctx, "workspace/inlineValue/refresh", nil)
	return err
}

// DiagnosticRefresh sends a workspace/diagnostic/refresh request to the client.
func (c *Client) DiagnosticRefresh(ctx context.Context) error {
	_, err := c.conn.Call(ctx, "workspace/diagnostic/refresh", nil)
	return err
}

// CodeLensRefresh sends a workspace/codeLens/refresh request to the client.
func (c *Client) CodeLensRefresh(ctx context.Context) error {
	_, err := c.conn.Call(ctx, "workspace/codeLens/refresh", nil)
	return err
}

// SemanticTokensRefresh sends a workspace/semanticTokens/refresh request to the client.
func (c *Client) SemanticTokensRefresh(ctx context.Context) error {
	_, err := c.conn.Call(ctx, "workspace/semanticTokens/refresh", nil)
	return err
}

// ShowDocument sends a window/showDocument request to the client and waits for a response.
func (c *Client) ShowDocument(ctx context.Context, params *lsp.ShowDocumentParams) (*lsp.ShowDocumentResult, error) {
	resp, err := c.conn.Call(ctx, "window/showDocument", params)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("showDocument: %s", resp.Error.Message)
	}
	var result lsp.ShowDocumentResult
	if err := json.Unmarshal(resp.Result, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
