package client

import (
	"context"
	"encoding/json"

	"github.com/owenrumney/go-lsp/internal/jsonrpc"
	"github.com/owenrumney/go-lsp/lsp"
)

// Handler answers the requests and notifications a server sends to its client.
// Dynamic registration and progress are handled by Client itself and are not
// surfaced here.
type Handler interface {
	Configuration(ctx context.Context, items []lsp.ConfigurationItem) ([]json.RawMessage, error)
	ApplyEdit(ctx context.Context, params *lsp.ApplyWorkspaceEditParams) (*lsp.ApplyWorkspaceEditResult, error)
	ShowMessageRequest(ctx context.Context, params *lsp.ShowMessageRequestParams) (*lsp.MessageActionItem, error)
	ShowDocument(ctx context.Context, params *lsp.ShowDocumentParams) (*lsp.ShowDocumentResult, error)

	Diagnostics(params *lsp.PublishDiagnosticsParams)
	ShowMessage(params *lsp.ShowMessageParams)
	LogMessage(params *lsp.LogMessageParams)
	Progress(params *lsp.ProgressParams)
	Notification(method string, params json.RawMessage)
}

// NopHandler declines everything a server asks and ignores what it reports.
// Embed it to implement only the parts you need.
type NopHandler struct{}

func (NopHandler) Configuration(_ context.Context, items []lsp.ConfigurationItem) ([]json.RawMessage, error) {
	out := make([]json.RawMessage, len(items))
	for i := range out {
		out[i] = json.RawMessage("null")
	}
	return out, nil
}

func (NopHandler) ApplyEdit(_ context.Context, _ *lsp.ApplyWorkspaceEditParams) (*lsp.ApplyWorkspaceEditResult, error) {
	return &lsp.ApplyWorkspaceEditResult{
		Applied:       false,
		FailureReason: "client does not apply edits",
	}, nil
}

func (NopHandler) ShowMessageRequest(_ context.Context, _ *lsp.ShowMessageRequestParams) (*lsp.MessageActionItem, error) {
	return nil, nil
}

func (NopHandler) ShowDocument(_ context.Context, _ *lsp.ShowDocumentParams) (*lsp.ShowDocumentResult, error) {
	return &lsp.ShowDocumentResult{Success: false}, nil
}

func (NopHandler) Diagnostics(_ *lsp.PublishDiagnosticsParams) {}
func (NopHandler) ShowMessage(_ *lsp.ShowMessageParams)        {}
func (NopHandler) LogMessage(_ *lsp.LogMessageParams)          {}
func (NopHandler) Progress(_ *lsp.ProgressParams)              {}
func (NopHandler) Notification(_ string, _ json.RawMessage)    {}

func (c *Client) dispatcher(h Handler) *jsonrpc.Dispatcher {
	d := jsonrpc.NewDispatcher()

	d.RegisterMethod("workspace/configuration", func(ctx context.Context, raw json.RawMessage) (any, error) {
		var params lsp.ConfigurationParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
		}
		return h.Configuration(ctx, params.Items)
	})

	d.RegisterMethod("workspace/applyEdit", func(ctx context.Context, raw json.RawMessage) (any, error) {
		var params lsp.ApplyWorkspaceEditParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
		}
		return h.ApplyEdit(ctx, &params)
	})

	d.RegisterMethod("window/showMessageRequest", func(ctx context.Context, raw json.RawMessage) (any, error) {
		var params lsp.ShowMessageRequestParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
		}
		return h.ShowMessageRequest(ctx, &params)
	})

	d.RegisterMethod("window/showDocument", func(ctx context.Context, raw json.RawMessage) (any, error) {
		var params lsp.ShowDocumentParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
		}
		return h.ShowDocument(ctx, &params)
	})

	d.RegisterMethod("client/registerCapability", func(_ context.Context, raw json.RawMessage) (any, error) {
		var params lsp.RegistrationParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
		}
		c.registry.register(params.Registrations)
		return nil, nil
	})

	d.RegisterMethod("client/unregisterCapability", func(_ context.Context, raw json.RawMessage) (any, error) {
		var params lsp.UnregistrationParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
		}
		c.registry.unregister(params.Unregisterations)
		return nil, nil
	})

	d.RegisterMethod("window/workDoneProgress/create", func(_ context.Context, raw json.RawMessage) (any, error) {
		var params lsp.WorkDoneProgressCreateParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil, jsonrpc.NewError(jsonrpc.CodeInvalidParams, err.Error())
		}
		c.ready.create(string(params.Token))
		return nil, nil
	})

	// Refresh requests have no client-side cache to invalidate; acknowledge them.
	for _, method := range []string{
		"workspace/semanticTokens/refresh",
		"workspace/inlayHint/refresh",
		"workspace/inlineValue/refresh",
		"workspace/codeLens/refresh",
		"workspace/diagnostic/refresh",
	} {
		d.RegisterMethod(method, func(_ context.Context, _ json.RawMessage) (any, error) {
			return nil, nil
		})
	}

	d.RegisterNotification("textDocument/publishDiagnostics", func(_ context.Context, raw json.RawMessage) error {
		var params lsp.PublishDiagnosticsParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil
		}
		h.Diagnostics(&params)
		return nil
	})

	d.RegisterNotification("window/showMessage", func(_ context.Context, raw json.RawMessage) error {
		var params lsp.ShowMessageParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil
		}
		h.ShowMessage(&params)
		return nil
	})

	d.RegisterNotification("window/logMessage", func(_ context.Context, raw json.RawMessage) error {
		var params lsp.LogMessageParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil
		}
		h.LogMessage(&params)
		return nil
	})

	d.RegisterNotification("$/progress", func(_ context.Context, raw json.RawMessage) error {
		var params lsp.ProgressParams
		if err := json.Unmarshal(raw, &params); err != nil {
			return nil
		}
		c.ready.observe(&params)
		h.Progress(&params)
		return nil
	})

	d.SetUnknownNotificationHandler(func(_ context.Context, method string, raw json.RawMessage) error {
		h.Notification(method, raw)
		return nil
	})

	return d
}
