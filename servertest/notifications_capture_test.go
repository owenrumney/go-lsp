package servertest_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/owenrumney/go-lsp/lsp"
	"github.com/owenrumney/go-lsp/server"
	"github.com/owenrumney/go-lsp/servertest"
)

type progressHandler struct {
	client *server.Client
}

func (h *progressHandler) Initialize(_ context.Context, _ *lsp.InitializeParams) (*lsp.InitializeResult, error) {
	return &lsp.InitializeResult{ServerInfo: &lsp.ServerInfo{Name: "progress-test"}}, nil
}

func (h *progressHandler) Shutdown(_ context.Context) error { return nil }

func (h *progressHandler) SetClient(client *server.Client) { h.client = client }

func (h *progressHandler) Hover(ctx context.Context, _ *lsp.HoverParams) (*lsp.Hover, error) {
	token := lsp.ProgressToken(json.RawMessage(`"indexing"`))
	err := h.client.Progress(ctx, &lsp.ProgressParams{
		Token: token,
		Value: json.RawMessage(`{"kind":"begin","title":"Indexing"}`),
	})
	if err != nil {
		return nil, err
	}
	return &lsp.Hover{Contents: lsp.NewHoverContents(lsp.PlainText, "hi")}, nil
}

func TestHarnessCapturesUntypedNotifications(t *testing.T) {
	h := servertest.New(t, &progressHandler{})

	if _, err := h.Hover("file:///a.txt", 0, 0); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
	defer cancel()

	n, err := h.WaitForNotification(ctx, "$/progress")
	if err != nil {
		t.Fatalf("waiting for $/progress: %v", err)
	}

	var params struct {
		Token json.RawMessage `json:"token"`
		Value struct {
			Kind  string `json:"kind"`
			Title string `json:"title"`
		} `json:"value"`
	}
	if err := json.Unmarshal(n.Params, &params); err != nil {
		t.Fatal(err)
	}
	if params.Value.Title != "Indexing" {
		t.Fatalf("title = %q, want Indexing", params.Value.Title)
	}

	var found bool
	for _, got := range h.Notifications() {
		if got.Method == "$/progress" {
			found = true
		}
	}
	if !found {
		t.Fatalf("Notifications() = %v, want a $/progress entry", h.Notifications())
	}
}
