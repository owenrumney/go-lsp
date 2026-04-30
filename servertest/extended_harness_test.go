package servertest_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
	"github.com/owenrumney/go-lsp/servertest"
)

type fullHandler struct {
	client *server.Client
	events chan string

	workspaceFoldersChanged bool
	configurationChanged    bool
	watchedFilesChanged     bool
	trace                   lsp.TraceValue
}

func (h *fullHandler) SetClient(c *server.Client) { h.client = c }

func (h *fullHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{ServerInfo: &lsp.ServerInfo{Name: "full-test"}}, nil
}

func (h *fullHandler) Shutdown(_ context.Context) error { return nil }

func (h *fullHandler) DidOpen(_ context.Context, _ *lsp.DidOpenTextDocumentParams) error {
	return nil
}
func (h *fullHandler) DidChange(_ context.Context, _ *lsp.DidChangeTextDocumentParams) error {
	return nil
}
func (h *fullHandler) DidClose(_ context.Context, _ *lsp.DidCloseTextDocumentParams) error {
	return nil
}
func (h *fullHandler) WillSave(_ context.Context, _ *lsp.WillSaveTextDocumentParams) error {
	return nil
}
func (h *fullHandler) DidChangeWorkspaceFolders(_ context.Context, _ *lsp.DidChangeWorkspaceFoldersParams) error {
	h.workspaceFoldersChanged = true
	h.record("workspaceFolders")
	return nil
}
func (h *fullHandler) DidChangeConfiguration(_ context.Context, _ *lsp.DidChangeConfigurationParams) error {
	h.configurationChanged = true
	h.record("configuration")
	return nil
}
func (h *fullHandler) DidChangeWatchedFiles(_ context.Context, _ *lsp.DidChangeWatchedFilesParams) error {
	h.watchedFilesChanged = true
	h.record("watchedFiles")
	return nil
}
func (h *fullHandler) SetTrace(_ context.Context, params *lsp.SetTraceParams) error {
	h.trace = params.Value
	h.record("trace")
	return nil
}

func (h *fullHandler) record(event string) {
	if h.events == nil {
		return
	}
	select {
	case h.events <- event:
	default:
	}
}

func (h *fullHandler) Completion(_ context.Context, _ *lsp.CompletionParams) (*lsp.CompletionList, error) {
	return &lsp.CompletionList{Items: []lsp.CompletionItem{{Label: "complete"}}}, nil
}
func (h *fullHandler) ResolveCompletionItem(_ context.Context, item *lsp.CompletionItem) (*lsp.CompletionItem, error) {
	item.Detail = "resolved"
	return item, nil
}
func (h *fullHandler) Hover(_ context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	return &lsp.Hover{Contents: lsp.MarkupContent{Kind: lsp.Markdown, Value: "hover"}}, nil
}
func (h *fullHandler) SignatureHelp(_ context.Context, _ *lsp.SignatureHelpParams) (*lsp.SignatureHelp, error) {
	return &lsp.SignatureHelp{Signatures: []lsp.SignatureInformation{{Label: "fn(a int)"}}}, nil
}
func (h *fullHandler) Declaration(_ context.Context, _ *lsp.DeclarationParams) ([]lsp.Location, error) {
	return []lsp.Location{testLocation("file:///declaration.go")}, nil
}
func (h *fullHandler) Definition(_ context.Context, _ *lsp.DefinitionParams) ([]lsp.Location, error) {
	return []lsp.Location{testLocation("file:///definition.go")}, nil
}
func (h *fullHandler) TypeDefinition(_ context.Context, _ *lsp.TypeDefinitionParams) ([]lsp.Location, error) {
	return []lsp.Location{testLocation("file:///type.go")}, nil
}
func (h *fullHandler) Implementation(_ context.Context, _ *lsp.ImplementationParams) ([]lsp.Location, error) {
	return []lsp.Location{testLocation("file:///implementation.go")}, nil
}
func (h *fullHandler) References(_ context.Context, _ *lsp.ReferenceParams) ([]lsp.Location, error) {
	return []lsp.Location{testLocation("file:///reference.go")}, nil
}
func (h *fullHandler) DocumentHighlight(_ context.Context, _ *lsp.DocumentHighlightParams) ([]lsp.DocumentHighlight, error) {
	kind := lsp.DocumentHighlightKindRead
	return []lsp.DocumentHighlight{{Range: testRange(), Kind: &kind}}, nil
}
func (h *fullHandler) DocumentSymbol(_ context.Context, _ *lsp.DocumentSymbolParams) ([]lsp.DocumentSymbol, error) {
	return []lsp.DocumentSymbol{{Name: "Symbol", Kind: lsp.SymbolKindFunction, Range: testRange(), SelectionRange: testRange()}}, nil
}
func (h *fullHandler) CodeAction(_ context.Context, _ *lsp.CodeActionParams) ([]lsp.CodeAction, error) {
	kind := lsp.CodeActionQuickFix
	return []lsp.CodeAction{{Title: "Fix", Kind: &kind}}, nil
}
func (h *fullHandler) ResolveCodeAction(_ context.Context, action *lsp.CodeAction) (*lsp.CodeAction, error) {
	action.Command = &lsp.Command{Title: "Run", Command: "test.run"}
	return action, nil
}
func (h *fullHandler) CodeLens(_ context.Context, _ *lsp.CodeLensParams) ([]lsp.CodeLens, error) {
	return []lsp.CodeLens{{Range: testRange()}}, nil
}
func (h *fullHandler) ResolveCodeLens(_ context.Context, lens *lsp.CodeLens) (*lsp.CodeLens, error) {
	lens.Command = &lsp.Command{Title: "Lens", Command: "lens.run"}
	return lens, nil
}
func (h *fullHandler) DocumentLink(_ context.Context, _ *lsp.DocumentLinkParams) ([]lsp.DocumentLink, error) {
	target := lsp.DocumentURI("https://example.com")
	return []lsp.DocumentLink{{Range: testRange(), Target: &target}}, nil
}
func (h *fullHandler) ResolveDocumentLink(_ context.Context, link *lsp.DocumentLink) (*lsp.DocumentLink, error) {
	link.Tooltip = "resolved"
	return link, nil
}
func (h *fullHandler) DocumentColor(_ context.Context, _ *lsp.DocumentColorParams) ([]lsp.ColorInformation, error) {
	return []lsp.ColorInformation{{Range: testRange(), Color: lsp.Color{Red: 1, Alpha: 1}}}, nil
}
func (h *fullHandler) ColorPresentation(_ context.Context, _ *lsp.ColorPresentationParams) ([]lsp.ColorPresentation, error) {
	return []lsp.ColorPresentation{{Label: "#ff0000"}}, nil
}
func (h *fullHandler) Formatting(_ context.Context, _ *lsp.DocumentFormattingParams) ([]lsp.TextEdit, error) {
	return []lsp.TextEdit{{Range: testRange(), NewText: "formatted"}}, nil
}
func (h *fullHandler) RangeFormatting(_ context.Context, _ *lsp.DocumentRangeFormattingParams) ([]lsp.TextEdit, error) {
	return []lsp.TextEdit{{Range: testRange(), NewText: "range"}}, nil
}
func (h *fullHandler) OnTypeFormatting(_ context.Context, _ *lsp.DocumentOnTypeFormattingParams) ([]lsp.TextEdit, error) {
	return []lsp.TextEdit{{Range: testRange(), NewText: "type"}}, nil
}
func (h *fullHandler) Rename(_ context.Context, _ *lsp.RenameParams) (*lsp.WorkspaceEdit, error) {
	return &lsp.WorkspaceEdit{Changes: map[lsp.DocumentURI][]lsp.TextEdit{"file:///rename.go": {{Range: testRange(), NewText: "new"}}}}, nil
}
func (h *fullHandler) PrepareRename(_ context.Context, _ *lsp.PrepareRenameParams) (*lsp.PrepareRenameResult, error) {
	return &lsp.PrepareRenameResult{Range: testRange(), Placeholder: "old"}, nil
}
func (h *fullHandler) FoldingRange(_ context.Context, _ *lsp.FoldingRangeParams) ([]lsp.FoldingRange, error) {
	return []lsp.FoldingRange{{StartLine: 0, EndLine: 2}}, nil
}
func (h *fullHandler) SelectionRange(_ context.Context, _ *lsp.SelectionRangeParams) ([]lsp.SelectionRange, error) {
	return []lsp.SelectionRange{{Range: testRange()}}, nil
}
func (h *fullHandler) LinkedEditingRange(_ context.Context, _ *lsp.LinkedEditingRangeParams) (*lsp.LinkedEditingRanges, error) {
	return &lsp.LinkedEditingRanges{Ranges: []lsp.Range{testRange()}, WordPattern: "[a-z]+"}, nil
}
func (h *fullHandler) Moniker(_ context.Context, _ *lsp.MonikerParams) ([]lsp.Moniker, error) {
	kind := lsp.MonikerKindExport
	return []lsp.Moniker{{Scheme: "test", Identifier: "id", Unique: lsp.UniquenessLevelGlobal, Kind: &kind}}, nil
}
func (h *fullHandler) WillSaveWaitUntil(_ context.Context, _ *lsp.WillSaveTextDocumentParams) ([]lsp.TextEdit, error) {
	return []lsp.TextEdit{{Range: testRange(), NewText: "save"}}, nil
}
func (h *fullHandler) WorkspaceSymbol(_ context.Context, _ *lsp.WorkspaceSymbolParams) ([]lsp.SymbolInformation, error) {
	return []lsp.SymbolInformation{{Name: "Workspace", Kind: lsp.SymbolKindFunction, Location: testLocation("file:///workspace.go")}}, nil
}
func (h *fullHandler) ExecuteCommand(ctx context.Context, params *lsp.ExecuteCommandParams) (any, error) {
	if params.Command == "ask" {
		item, err := h.client.ShowMessageRequest(ctx, &lsp.ShowMessageRequestParams{
			Type:    lsp.MessageTypeInfo,
			Message: "choose",
			Actions: []lsp.MessageActionItem{
				{Title: "ok"},
			},
		})
		if err != nil {
			return nil, err
		}
		return item.Title, nil
	}
	return params.Command, nil
}
func (h *fullHandler) WillCreateFiles(_ context.Context, _ *lsp.CreateFilesParams) (*lsp.WorkspaceEdit, error) {
	return &lsp.WorkspaceEdit{Changes: map[lsp.DocumentURI][]lsp.TextEdit{"file:///created.go": {{Range: testRange(), NewText: "create"}}}}, nil
}
func (h *fullHandler) WillRenameFiles(_ context.Context, _ *lsp.RenameFilesParams) (*lsp.WorkspaceEdit, error) {
	return &lsp.WorkspaceEdit{Changes: map[lsp.DocumentURI][]lsp.TextEdit{"file:///renamed.go": {{Range: testRange(), NewText: "rename"}}}}, nil
}
func (h *fullHandler) WillDeleteFiles(_ context.Context, _ *lsp.DeleteFilesParams) (*lsp.WorkspaceEdit, error) {
	return &lsp.WorkspaceEdit{Changes: map[lsp.DocumentURI][]lsp.TextEdit{"file:///deleted.go": {{Range: testRange(), NewText: "delete"}}}}, nil
}
func (h *fullHandler) PrepareCallHierarchy(_ context.Context, _ *lsp.CallHierarchyPrepareParams) ([]lsp.CallHierarchyItem, error) {
	return []lsp.CallHierarchyItem{testCallItem("caller")}, nil
}
func (h *fullHandler) IncomingCalls(_ context.Context, _ *lsp.CallHierarchyIncomingCallsParams) ([]lsp.CallHierarchyIncomingCall, error) {
	return []lsp.CallHierarchyIncomingCall{{From: testCallItem("from"), FromRanges: []lsp.Range{testRange()}}}, nil
}
func (h *fullHandler) OutgoingCalls(_ context.Context, _ *lsp.CallHierarchyOutgoingCallsParams) ([]lsp.CallHierarchyOutgoingCall, error) {
	return []lsp.CallHierarchyOutgoingCall{{To: testCallItem("to"), FromRanges: []lsp.Range{testRange()}}}, nil
}
func (h *fullHandler) SemanticTokensFull(_ context.Context, _ *lsp.SemanticTokensParams) (*lsp.SemanticTokens, error) {
	return &lsp.SemanticTokens{ResultID: "full", Data: []int{0, 0, 1, 0, 0}}, nil
}
func (h *fullHandler) SemanticTokensDelta(_ context.Context, _ *lsp.SemanticTokensDeltaParams) (*lsp.SemanticTokensDelta, error) {
	return &lsp.SemanticTokensDelta{ResultID: "delta", Edits: []lsp.SemanticTokensEdit{{Start: 0, DeleteCount: 1}}}, nil
}
func (h *fullHandler) SemanticTokensRange(_ context.Context, _ *lsp.SemanticTokensRangeParams) (*lsp.SemanticTokens, error) {
	return &lsp.SemanticTokens{Data: []int{0, 0, 1, 0, 0}}, nil
}
func (h *fullHandler) PrepareTypeHierarchy(_ context.Context, _ *lsp.TypeHierarchyPrepareParams) ([]lsp.TypeHierarchyItem, error) {
	return []lsp.TypeHierarchyItem{testTypeItem("type")}, nil
}
func (h *fullHandler) Supertypes(_ context.Context, _ *lsp.TypeHierarchySupertypesParams) ([]lsp.TypeHierarchyItem, error) {
	return []lsp.TypeHierarchyItem{testTypeItem("super")}, nil
}
func (h *fullHandler) Subtypes(_ context.Context, _ *lsp.TypeHierarchySubtypesParams) ([]lsp.TypeHierarchyItem, error) {
	return []lsp.TypeHierarchyItem{testTypeItem("sub")}, nil
}
func (h *fullHandler) InlayHint(_ context.Context, _ *lsp.InlayHintParams) ([]lsp.InlayHint, error) {
	return []lsp.InlayHint{{Position: lsp.Position{}, Label: json.RawMessage(`"hint"`)}}, nil
}
func (h *fullHandler) ResolveInlayHint(_ context.Context, hint *lsp.InlayHint) (*lsp.InlayHint, error) {
	hint.Tooltip = &lsp.MarkupContent{Kind: lsp.PlainText, Value: "resolved"}
	return hint, nil
}
func (h *fullHandler) InlineValue(_ context.Context, _ *lsp.InlineValueParams) ([]json.RawMessage, error) {
	return []json.RawMessage{json.RawMessage(`{"text":"value","range":{"start":{"line":0,"character":0},"end":{"line":0,"character":1}}}`)}, nil
}
func (h *fullHandler) DocumentDiagnostic(_ context.Context, _ *lsp.DocumentDiagnosticParams) (any, error) {
	return lsp.FullDocumentDiagnosticReport{Kind: string(lsp.DiagnosticReportFull), Items: []lsp.Diagnostic{{Range: testRange(), Message: "diag"}}}, nil
}
func (h *fullHandler) WorkspaceDiagnostic(_ context.Context, _ *lsp.WorkspaceDiagnosticParams) (*lsp.WorkspaceDiagnosticReport, error) {
	return &lsp.WorkspaceDiagnosticReport{Items: []json.RawMessage{json.RawMessage(`{"kind":"full","uri":"file:///test.go","items":[]}`)}}, nil
}

func TestAdditionalTypedRequestHelpers(t *testing.T) {
	h := servertest.New(t, &fullHandler{})
	uri := lsp.DocumentURI("file:///test.go")
	r := testRange()

	checkCall(t, "declaration", func() ([]lsp.Location, error) { return h.Declaration(uri, 0, 0) })
	checkCall(t, "definition", func() ([]lsp.Location, error) { return h.Definition(uri, 0, 0) })
	checkCall(t, "type definition", func() ([]lsp.Location, error) { return h.TypeDefinition(uri, 0, 0) })
	checkCall(t, "implementation", func() ([]lsp.Location, error) { return h.Implementation(uri, 0, 0) })
	checkCall(t, "references", func() ([]lsp.Location, error) { return h.References(uri, 0, 0, true) })
	checkCall(t, "document highlight", func() ([]lsp.DocumentHighlight, error) { return h.DocumentHighlight(uri, 0, 0) })
	checkCall(t, "code action", func() ([]lsp.CodeAction, error) {
		return h.CodeAction(&lsp.CodeActionParams{TextDocument: lsp.TextDocumentIdentifier{URI: uri}, Range: r})
	})
	checkCall(t, "code lens", func() ([]lsp.CodeLens, error) { return h.CodeLens(uri) })
	checkCall(t, "document link", func() ([]lsp.DocumentLink, error) { return h.DocumentLink(uri) })
	checkCall(t, "document color", func() ([]lsp.ColorInformation, error) { return h.DocumentColor(uri) })
	checkCall(t, "color presentation", func() ([]lsp.ColorPresentation, error) {
		return h.ColorPresentation(&lsp.ColorPresentationParams{TextDocument: lsp.TextDocumentIdentifier{URI: uri}, Range: r})
	})
	checkCall(t, "document symbol", func() ([]lsp.DocumentSymbol, error) { return h.DocumentSymbol(uri) })
	checkCall(t, "workspace symbol", func() ([]lsp.SymbolInformation, error) { return h.WorkspaceSymbol("work") })
	checkCall(t, "formatting", func() ([]lsp.TextEdit, error) { return h.Formatting(uri) })
	checkCall(t, "range formatting", func() ([]lsp.TextEdit, error) { return h.RangeFormatting(uri, r) })
	checkCall(t, "on type formatting", func() ([]lsp.TextEdit, error) { return h.OnTypeFormatting(uri, 0, 0, "}") })
	checkCall(t, "folding range", func() ([]lsp.FoldingRange, error) { return h.FoldingRange(uri) })
	checkCall(t, "selection range", func() ([]lsp.SelectionRange, error) {
		return h.SelectionRange(uri, []lsp.Position{{}})
	})
	checkCall(t, "moniker", func() ([]lsp.Moniker, error) { return h.Moniker(uri, 0, 0) })
	checkCall(t, "will save wait until", func() ([]lsp.TextEdit, error) { return h.WillSaveWaitUntil(uri, lsp.SaveManual) })
	checkCall(t, "inlay hint", func() ([]lsp.InlayHint, error) { return h.InlayHint(uri, r) })
	checkCall(t, "prepare call hierarchy", func() ([]lsp.CallHierarchyItem, error) { return h.PrepareCallHierarchy(uri, 0, 0) })
	checkCall(t, "incoming calls", func() ([]lsp.CallHierarchyIncomingCall, error) { return h.IncomingCalls(testCallItem("item")) })
	checkCall(t, "outgoing calls", func() ([]lsp.CallHierarchyOutgoingCall, error) { return h.OutgoingCalls(testCallItem("item")) })
	checkCall(t, "prepare type hierarchy", func() ([]lsp.TypeHierarchyItem, error) { return h.PrepareTypeHierarchy(uri, 0, 0) })
	checkCall(t, "supertypes", func() ([]lsp.TypeHierarchyItem, error) { return h.Supertypes(testTypeItem("item")) })
	checkCall(t, "subtypes", func() ([]lsp.TypeHierarchyItem, error) { return h.Subtypes(testTypeItem("item")) })
	checkCall(t, "inline value", func() ([]json.RawMessage, error) {
		return h.InlineValue(&lsp.InlineValueParams{
			TextDocument: lsp.TextDocumentIdentifier{URI: uri},
			Range:        r,
			Context:      lsp.InlineValueContext{StoppedLocation: r},
		})
	})

	if got, err := h.SignatureHelp(uri, 0, 0); err != nil || got == nil || len(got.Signatures) != 1 {
		t.Fatalf("signature help = %+v, %v", got, err)
	}
	if got, err := h.ResolveCompletionItem(&lsp.CompletionItem{Label: "complete"}); err != nil || got.Detail != "resolved" {
		t.Fatalf("resolve completion = %+v, %v", got, err)
	}
	if got, err := h.ResolveCodeAction(&lsp.CodeAction{Title: "Fix"}); err != nil || got.Command == nil {
		t.Fatalf("resolve code action = %+v, %v", got, err)
	}
	if got, err := h.ResolveCodeLens(&lsp.CodeLens{Range: r}); err != nil || got.Command == nil {
		t.Fatalf("resolve code lens = %+v, %v", got, err)
	}
	if got, err := h.ResolveDocumentLink(&lsp.DocumentLink{Range: r}); err != nil || got.Tooltip != "resolved" {
		t.Fatalf("resolve document link = %+v, %v", got, err)
	}
	if got, err := h.ResolveInlayHint(&lsp.InlayHint{Label: json.RawMessage(`"hint"`)}); err != nil || got.Tooltip == nil {
		t.Fatalf("resolve inlay hint = %+v, %v", got, err)
	}
	if got, err := h.Rename(uri, 0, 0, "new"); err != nil || got == nil || len(got.Changes) != 1 {
		t.Fatalf("rename = %+v, %v", got, err)
	}
	if got, err := h.PrepareRename(uri, 0, 0); err != nil || got.Placeholder != "old" {
		t.Fatalf("prepare rename = %+v, %v", got, err)
	}
	if got, err := h.LinkedEditingRange(uri, 0, 0); err != nil || got.WordPattern == "" {
		t.Fatalf("linked editing = %+v, %v", got, err)
	}
	if got, err := h.DocumentDiagnostic(&lsp.DocumentDiagnosticParams{TextDocument: lsp.TextDocumentIdentifier{URI: uri}}); err != nil || !strings.Contains(string(got), `"full"`) {
		t.Fatalf("document diagnostic = %s, %v", got, err)
	}
	if got, err := h.WorkspaceDiagnostic(&lsp.WorkspaceDiagnosticParams{}); err != nil || got == nil || len(got.Items) != 1 {
		t.Fatalf("workspace diagnostic = %+v, %v", got, err)
	}
	if got, err := h.SemanticTokensFull(uri); err != nil || got.ResultID != "full" {
		t.Fatalf("semantic full = %+v, %v", got, err)
	}
	if got, err := h.SemanticTokensDelta(uri, "full"); err != nil || got.ResultID != "delta" {
		t.Fatalf("semantic delta = %+v, %v", got, err)
	}
	if got, err := h.SemanticTokensRange(uri, r); err != nil || len(got.Data) == 0 {
		t.Fatalf("semantic range = %+v, %v", got, err)
	}
	if got, err := h.WillCreateFiles([]lsp.FileCreate{{URI: string(uri)}}); err != nil || got == nil || len(got.Changes) != 1 {
		t.Fatalf("will create = %+v, %v", got, err)
	}
	if got, err := h.WillRenameFiles([]lsp.FileRename{{OldURI: string(uri), NewURI: "file:///new.go"}}); err != nil || got == nil || len(got.Changes) != 1 {
		t.Fatalf("will rename = %+v, %v", got, err)
	}
	if got, err := h.WillDeleteFiles([]lsp.FileDelete{{URI: string(uri)}}); err != nil || got == nil || len(got.Changes) != 1 {
		t.Fatalf("will delete = %+v, %v", got, err)
	}
	if got, err := h.ExecuteCommand("echo", nil); err != nil || !strings.Contains(string(got), "echo") {
		t.Fatalf("execute command = %s, %v", got, err)
	}
}

func TestWorkspaceNotificationHelpers(t *testing.T) {
	handler := &fullHandler{events: make(chan string, 4)}
	h := servertest.New(t, handler)

	if err := h.DidChangeWorkspaceFolders(&lsp.DidChangeWorkspaceFoldersParams{}); err != nil {
		t.Fatal(err)
	}
	if err := h.DidChangeConfiguration(&lsp.DidChangeConfigurationParams{Settings: map[string]any{"x": true}}); err != nil {
		t.Fatal(err)
	}
	if err := h.DidChangeWatchedFiles(&lsp.DidChangeWatchedFilesParams{Changes: []lsp.FileEvent{{URI: "file:///x.go"}}}); err != nil {
		t.Fatal(err)
	}
	if err := h.SetTrace(lsp.TraceVerbose); err != nil {
		t.Fatal(err)
	}

	seen := map[string]bool{}
	for len(seen) < 4 {
		select {
		case event := <-handler.events:
			seen[event] = true
		case <-time.After(2 * time.Second):
			t.Fatalf("timed out waiting for notification events: %v", seen)
		}
	}
	if !handler.workspaceFoldersChanged || !handler.configurationChanged || !handler.watchedFilesChanged || handler.trace != lsp.TraceVerbose {
		t.Fatalf("notifications not recorded: workspace=%v config=%v watched=%v trace=%q",
			handler.workspaceFoldersChanged, handler.configurationChanged, handler.watchedFilesChanged, handler.trace)
	}
}

func TestClientRequestRecordingAndResponses(t *testing.T) {
	h := servertest.New(t, &fullHandler{})
	h.SetClientResponse("window/showMessageRequest", lsp.MessageActionItem{Title: "ok"})

	result, err := h.ExecuteCommand("ask", nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(result), "ok") {
		t.Fatalf("execute command result = %s", result)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	req, err := h.WaitForClientRequest(ctx, "window/showMessageRequest")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(req.Params), "choose") {
		t.Fatalf("client request params = %s", req.Params)
	}
	if len(h.ClientRequests()) != 1 {
		t.Fatalf("client requests = %d, want 1", len(h.ClientRequests()))
	}
}

type messageHandler struct {
	client *server.Client
}

func (h *messageHandler) SetClient(c *server.Client) { h.client = c }
func (h *messageHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{}, nil
}
func (h *messageHandler) Shutdown(_ context.Context) error { return nil }
func (h *messageHandler) DidSave(ctx context.Context, _ *lsp.DidSaveTextDocumentParams) error {
	if err := h.client.ShowMessage(ctx, &lsp.ShowMessageParams{Type: lsp.MessageTypeInfo, Message: "shown"}); err != nil {
		return err
	}
	return h.client.LogMessage(ctx, &lsp.LogMessageParams{Type: lsp.MessageTypeLog, Message: "logged"})
}

func TestWaitForMessageAndLogMessage(t *testing.T) {
	h := servertest.New(t, &messageHandler{})
	if err := h.DidSave("file:///test.go"); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()

	msg, err := h.WaitForMessage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if msg.Message != "shown" {
		t.Fatalf("message = %+v", msg)
	}
	log, err := h.WaitForLogMessage(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if log.Message != "logged" {
		t.Fatalf("log = %+v", log)
	}
}

type cancelHandler struct {
	started chan struct{}
}

func (h *cancelHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{}, nil
}
func (h *cancelHandler) Shutdown(_ context.Context) error { return nil }
func (h *cancelHandler) Hover(ctx context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	close(h.started)
	<-ctx.Done()
	return nil, ctx.Err()
}

func TestCallAsyncAndCancelRequest(t *testing.T) {
	handler := &cancelHandler{started: make(chan struct{})}
	h := servertest.New(t, handler)

	call, err := h.CallAsync("textDocument/hover", &lsp.HoverParams{})
	if err != nil {
		t.Fatal(err)
	}
	<-handler.started
	if err := h.CancelRequest(call.ID()); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
	defer cancel()
	_, err = call.Wait(ctx)
	if err == nil {
		t.Fatal("expected cancellation error")
	}
	if !strings.Contains(err.Error(), "cancel") && !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v", err)
	}
}

func checkCall[T any](t *testing.T, name string, fn func() ([]T, error)) {
	t.Helper()
	got, err := fn()
	if err != nil {
		t.Fatalf("%s error: %v", name, err)
	}
	if len(got) == 0 {
		t.Fatalf("%s returned no results", name)
	}
}

func testRange() lsp.Range {
	return lsp.Range{
		Start: lsp.Position{Line: 0, Character: 0},
		End:   lsp.Position{Line: 0, Character: 1},
	}
}

func testLocation(uri lsp.DocumentURI) lsp.Location {
	return lsp.Location{URI: uri, Range: testRange()}
}

func testCallItem(name string) lsp.CallHierarchyItem {
	return lsp.CallHierarchyItem{
		Name:           name,
		Kind:           lsp.SymbolKindFunction,
		URI:            "file:///call.go",
		Range:          testRange(),
		SelectionRange: testRange(),
	}
}

func testTypeItem(name string) lsp.TypeHierarchyItem {
	return lsp.TypeHierarchyItem{
		Name:           name,
		Kind:           lsp.SymbolKindClass,
		URI:            "file:///type.go",
		Range:          testRange(),
		SelectionRange: testRange(),
	}
}
