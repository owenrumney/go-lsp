package servertest

import (
	"encoding/json"

	"github.com/owenrumney/go-lsp/lsp"
)

// DidOpen sends a textDocument/didOpen notification.
func (h *Harness) DidOpen(uri lsp.DocumentURI, languageID, text string) error {
	h.versionsMu.Lock()
	h.versions[uri] = 1
	h.versionsMu.Unlock()

	return h.conn.notify(h.ctx, "textDocument/didOpen", &lsp.DidOpenTextDocumentParams{
		TextDocument: lsp.TextDocumentItem{
			URI:        uri,
			LanguageID: languageID,
			Version:    1,
			Text:       text,
		},
	})
}

// DidChange sends a textDocument/didChange notification with full document sync.
func (h *Harness) DidChange(uri lsp.DocumentURI, version int, text string) error {
	h.versionsMu.Lock()
	h.versions[uri] = version
	h.versionsMu.Unlock()

	return h.conn.notify(h.ctx, "textDocument/didChange", &lsp.DidChangeTextDocumentParams{
		TextDocument: lsp.VersionedTextDocumentIdentifier{
			TextDocumentIdentifier: lsp.TextDocumentIdentifier{URI: uri},
			Version:                version,
		},
		ContentChanges: []lsp.TextDocumentContentChangeEvent{
			{Text: text},
		},
	})
}

// DidSave sends a textDocument/didSave notification.
func (h *Harness) DidSave(uri lsp.DocumentURI) error {
	return h.conn.notify(h.ctx, "textDocument/didSave", &lsp.DidSaveTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// DidClose sends a textDocument/didClose notification.
func (h *Harness) DidClose(uri lsp.DocumentURI) error {
	h.versionsMu.Lock()
	delete(h.versions, uri)
	h.versionsMu.Unlock()

	return h.conn.notify(h.ctx, "textDocument/didClose", &lsp.DidCloseTextDocumentParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
}

// Hover sends a textDocument/hover request.
func (h *Harness) Hover(uri lsp.DocumentURI, line, char int) (*lsp.Hover, error) {
	result, err := h.conn.call(h.ctx, "textDocument/hover", &lsp.HoverParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Position:     lsp.Position{Line: line, Character: char},
		},
	})
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var hover lsp.Hover
	if err := json.Unmarshal(result, &hover); err != nil {
		return nil, err
	}
	return &hover, nil
}

// Completion sends a textDocument/completion request.
func (h *Harness) Completion(uri lsp.DocumentURI, line, char int) (*lsp.CompletionList, error) {
	result, err := h.conn.call(h.ctx, "textDocument/completion", &lsp.CompletionParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Position:     lsp.Position{Line: line, Character: char},
		},
	})
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var list lsp.CompletionList
	if err := json.Unmarshal(result, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

// Definition sends a textDocument/definition request.
func (h *Harness) Definition(uri lsp.DocumentURI, line, char int) ([]lsp.Location, error) {
	result, err := h.conn.call(h.ctx, "textDocument/definition", &lsp.DefinitionParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Position:     lsp.Position{Line: line, Character: char},
		},
	})
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var locs []lsp.Location
	if err := json.Unmarshal(result, &locs); err != nil {
		return nil, err
	}
	return locs, nil
}

// References sends a textDocument/references request.
func (h *Harness) References(uri lsp.DocumentURI, line, char int, includeDecl bool) ([]lsp.Location, error) {
	result, err := h.conn.call(h.ctx, "textDocument/references", &lsp.ReferenceParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Position:     lsp.Position{Line: line, Character: char},
		},
		Context: lsp.ReferenceContext{IncludeDeclaration: includeDecl},
	})
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var locs []lsp.Location
	if err := json.Unmarshal(result, &locs); err != nil {
		return nil, err
	}
	return locs, nil
}

// CodeAction sends a textDocument/codeAction request.
func (h *Harness) CodeAction(params *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	result, err := h.conn.call(h.ctx, "textDocument/codeAction", params)
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var actions []lsp.CodeAction
	if err := json.Unmarshal(result, &actions); err != nil {
		return nil, err
	}
	return actions, nil
}

// DocumentSymbol sends a textDocument/documentSymbol request.
func (h *Harness) DocumentSymbol(uri lsp.DocumentURI) ([]lsp.DocumentSymbol, error) {
	result, err := h.conn.call(h.ctx, "textDocument/documentSymbol", &lsp.DocumentSymbolParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
	})
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var symbols []lsp.DocumentSymbol
	if err := json.Unmarshal(result, &symbols); err != nil {
		return nil, err
	}
	return symbols, nil
}

// WorkspaceSymbol sends a workspace/symbol request.
func (h *Harness) WorkspaceSymbol(query string) ([]lsp.SymbolInformation, error) {
	result, err := h.conn.call(h.ctx, "workspace/symbol", &lsp.WorkspaceSymbolParams{
		Query: query,
	})
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var symbols []lsp.SymbolInformation
	if err := json.Unmarshal(result, &symbols); err != nil {
		return nil, err
	}
	return symbols, nil
}

// Formatting sends a textDocument/formatting request.
func (h *Harness) Formatting(uri lsp.DocumentURI) ([]lsp.TextEdit, error) {
	result, err := h.conn.call(h.ctx, "textDocument/formatting", &lsp.DocumentFormattingParams{
		TextDocument: lsp.TextDocumentIdentifier{URI: uri},
		Options: lsp.FormattingOptions{
			TabSize:      4,
			InsertSpaces: true,
		},
	})
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var edits []lsp.TextEdit
	if err := json.Unmarshal(result, &edits); err != nil {
		return nil, err
	}
	return edits, nil
}

// Rename sends a textDocument/rename request.
func (h *Harness) Rename(uri lsp.DocumentURI, line, char int, newName string) (*lsp.WorkspaceEdit, error) {
	result, err := h.conn.call(h.ctx, "textDocument/rename", &lsp.RenameParams{
		TextDocumentPositionParams: lsp.TextDocumentPositionParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Position:     lsp.Position{Line: line, Character: char},
		},
		NewName: newName,
	})
	if err != nil {
		return nil, err
	}
	if string(result) == "null" {
		return nil, nil
	}
	var edit lsp.WorkspaceEdit
	if err := json.Unmarshal(result, &edit); err != nil {
		return nil, err
	}
	return &edit, nil
}

// Call sends an arbitrary JSON-RPC request and returns the raw result.
func (h *Harness) Call(method string, params any) (json.RawMessage, error) {
	return h.conn.call(h.ctx, method, params)
}

// Notify sends an arbitrary JSON-RPC notification.
func (h *Harness) Notify(method string, params any) error {
	return h.conn.notify(h.ctx, method, params)
}
