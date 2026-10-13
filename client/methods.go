package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/owenrumney/go-lsp/lsp"
)

// ErrClosed is returned when the connection ended without a specific cause.
var ErrClosed = errors.New("client: connection closed")

// DidOpen tells the server a document is now open and owned by the client.
func (c *Client) DidOpen(ctx context.Context, uri lsp.DocumentURI, languageID string, version int, text string) error {
	return c.conn.Notify(ctx, "textDocument/didOpen", &lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: languageID,
			Version:    version,
			Text:       text,
		},
	})
}

// DidChangeFull replaces the whole document. Incremental sync is the caller's
// job: build the params and use Notify.
func (c *Client) DidChangeFull(ctx context.Context, uri lsp.DocumentURI, version int, text string) error {
	return c.conn.Notify(ctx, "textDocument/didChange", &lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{URI: uri},
			Version:                version,
		},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{{Text: text}},
	})
}

func (c *Client) DidSave(ctx context.Context, uri lsp.DocumentURI) error {
	return c.conn.Notify(ctx, "textDocument/didSave", &lsp.DidSaveTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

func (c *Client) DidClose(ctx context.Context, uri lsp.DocumentURI) error {
	return c.conn.Notify(ctx, "textDocument/didClose", &lsp.DidCloseTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

func (c *Client) DidChangeConfiguration(ctx context.Context, settings any) error {
	raw, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("client: marshalling settings: %w", err)
	}
	return c.conn.Notify(ctx, "workspace/didChangeConfiguration", &lsp.DidChangeConfigurationParams{
		Settings: raw,
	})
}

func (c *Client) DidChangeWatchedFiles(ctx context.Context, events []lsp.FileEvent) error {
	return c.conn.Notify(ctx, "workspace/didChangeWatchedFiles", &lsp.DidChangeWatchedFilesParams{
		Changes: events,
	})
}

func (c *Client) Hover(ctx context.Context, uri lsp.DocumentURI, pos lsp.Position) (*lsp.Hover, error) {
	var out *lsp.Hover
	err := c.call(ctx, "textDocument/hover", &lsp.HoverParams{
		TextDocumentPositionParams: textDocumentPosition(uri, pos),
	}, &out)
	return out, err
}

func (c *Client) Completion(ctx context.Context, uri lsp.DocumentURI, pos lsp.Position) (*lsp.CompletionList, error) {
	raw, err := c.callRaw(ctx, "textDocument/completion", &lsp.CompletionParams{
		TextDocumentPositionParams: textDocumentPosition(uri, pos),
	})
	if err != nil || isNull(raw) {
		return nil, err
	}

	var list lsp.CompletionList
	if err := json.Unmarshal(raw, &list); err == nil && list.Items != nil {
		return &list, nil
	}

	var items []lsp.CompletionItem
	if err := json.Unmarshal(raw, &items); err != nil {
		return nil, fmt.Errorf("client: completion: %w", err)
	}
	return &lsp.CompletionList{Items: items}, nil
}

func (c *Client) Definition(ctx context.Context, uri lsp.DocumentURI, pos lsp.Position) ([]lsp.Location, error) {
	return c.locations(ctx, "textDocument/definition", uri, pos)
}

func (c *Client) Declaration(ctx context.Context, uri lsp.DocumentURI, pos lsp.Position) ([]lsp.Location, error) {
	return c.locations(ctx, "textDocument/declaration", uri, pos)
}

func (c *Client) TypeDefinition(ctx context.Context, uri lsp.DocumentURI, pos lsp.Position) ([]lsp.Location, error) {
	return c.locations(ctx, "textDocument/typeDefinition", uri, pos)
}

func (c *Client) Implementation(ctx context.Context, uri lsp.DocumentURI, pos lsp.Position) ([]lsp.Location, error) {
	return c.locations(ctx, "textDocument/implementation", uri, pos)
}

func (c *Client) References(ctx context.Context, uri lsp.DocumentURI, pos lsp.Position, includeDeclaration bool) ([]lsp.Location, error) {
	var out []lsp.Location
	err := c.call(ctx, "textDocument/references", &lsp.ReferenceParams{
		TextDocumentPositionParams: textDocumentPosition(uri, pos),
		Context:                    lsp.ReferenceContext{IncludeDeclaration: includeDeclaration},
	}, &out)
	return out, err
}

func (c *Client) DocumentSymbol(ctx context.Context, uri lsp.DocumentURI) ([]lsp.DocumentSymbol, error) {
	var out []lsp.DocumentSymbol
	err := c.call(ctx, "textDocument/documentSymbol", &lsp.DocumentSymbolParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	}, &out)
	return out, err
}

func (c *Client) WorkspaceSymbol(ctx context.Context, query string) ([]lsp.SymbolInformation, error) {
	var out []lsp.SymbolInformation
	err := c.call(ctx, "workspace/symbol", &lsp.WorkspaceSymbolParams{Query: query}, &out)
	return out, err
}

func (c *Client) CodeAction(ctx context.Context, params *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	var out []lsp.CodeAction
	err := c.call(ctx, "textDocument/codeAction", params, &out)
	return out, err
}

func (c *Client) Formatting(ctx context.Context, uri lsp.DocumentURI, opts lsp.FormattingOptions) ([]lsp.TextEdit, error) {
	var out []lsp.TextEdit
	err := c.call(ctx, "textDocument/formatting", &lsp.DocumentFormattingParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Options:      opts,
	}, &out)
	return out, err
}

func (c *Client) Rename(ctx context.Context, uri lsp.DocumentURI, pos lsp.Position, newName string) (*lsp.WorkspaceEdit, error) {
	var out *lsp.WorkspaceEdit
	err := c.call(ctx, "textDocument/rename", &lsp.RenameParams{
		TextDocumentPositionParams: textDocumentPosition(uri, pos),
		NewName:                    newName,
	}, &out)
	return out, err
}

// DocumentDiagnostic issues a pull diagnostic request. The result is returned
// raw because the report is a union of full and unchanged forms.
func (c *Client) DocumentDiagnostic(ctx context.Context, uri lsp.DocumentURI, previousResultID string) (json.RawMessage, error) {
	params := &lsp.DocumentDiagnosticParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	}
	if previousResultID != "" {
		params.PreviousResultID = &previousResultID
	}
	return c.callRaw(ctx, "textDocument/diagnostic", params)
}

func (c *Client) ExecuteCommand(ctx context.Context, command string, args []json.RawMessage) (json.RawMessage, error) {
	return c.callRaw(ctx, "workspace/executeCommand", &lsp.ExecuteCommandParams{
		Command:   command,
		Arguments: args,
	})
}

// Call sends an arbitrary request and unmarshals the result into out.
func (c *Client) Call(ctx context.Context, method string, params, out any) error {
	return c.call(ctx, method, params, out)
}

// Notify sends an arbitrary notification.
func (c *Client) Notify(ctx context.Context, method string, params any) error {
	return c.conn.Notify(ctx, method, params)
}

func (c *Client) locations(ctx context.Context, method string, uri lsp.DocumentURI, pos lsp.Position) ([]lsp.Location, error) {
	raw, err := c.callRaw(ctx, method, &lsp.TextDocumentPositionParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Position:     pos,
	})
	if err != nil || isNull(raw) {
		return nil, err
	}

	var many []lsp.Location
	if err := json.Unmarshal(raw, &many); err == nil {
		return many, nil
	}

	var one lsp.Location
	if err := json.Unmarshal(raw, &one); err != nil {
		return nil, fmt.Errorf("client: %s: %w", method, err)
	}
	return []lsp.Location{one}, nil
}

func (c *Client) call(ctx context.Context, method string, params, out any) error {
	raw, err := c.callRaw(ctx, method, params)
	if err != nil || out == nil || isNull(raw) {
		return err
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("client: %s: %w", method, err)
	}
	return nil
}

func (c *Client) callRaw(ctx context.Context, method string, params any) (json.RawMessage, error) {
	resp, err := c.conn.Call(ctx, method, params)
	if err != nil {
		return nil, err
	}
	if resp.Error != nil {
		return nil, resp.Error
	}
	return resp.Result, nil
}

func isNull(raw json.RawMessage) bool {
	return len(raw) == 0 || string(raw) == "null"
}

func textDocumentPosition(uri lsp.DocumentURI, pos lsp.Position) lsp.TextDocumentPositionParams {
	return lsp.TextDocumentPositionParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Position:     pos,
	}
}
