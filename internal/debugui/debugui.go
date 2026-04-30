package debugui

import (
	"context"
	"encoding/json"
	"io/fs"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(_ *http.Request) bool { return true },
}

// Hub manages websocket clients and broadcasts entries to them.
type Hub struct {
	mu      sync.RWMutex
	clients map[*websocket.Conn]chan []byte
}

func newHub() *Hub {
	return &Hub{clients: make(map[*websocket.Conn]chan []byte)}
}

func (h *Hub) add(conn *websocket.Conn) chan []byte {
	ch := make(chan []byte, 256)
	h.mu.Lock()
	h.clients[conn] = ch
	h.mu.Unlock()
	return ch
}

func (h *Hub) remove(conn *websocket.Conn) {
	h.mu.Lock()
	if ch, ok := h.clients[conn]; ok {
		close(ch)
		delete(h.clients, conn)
	}
	h.mu.Unlock()
}

// Broadcast sends a value to all connected clients. Non-blocking: slow clients drop messages.
func (h *Hub) Broadcast(v any) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, ch := range h.clients {
		select {
		case ch <- data:
		default:
		}
	}
}

// wsMessage wraps an entry or log for WebSocket broadcast with a kind discriminator.
type wsMessage struct {
	Kind string `json:"kind"`
	Data any    `json:"data"`
}

// DebugUI is an HTTP server that serves the debug web interface and websocket
// over a Recorder's captured data.
type DebugUI struct {
	recorder *Recorder
	hub      *Hub
	stats    *Stats
	srv      *http.Server
}

// New creates a DebugUI bound to addr that exposes recorder's captured data.
func New(addr string, recorder *Recorder) *DebugUI {
	d := &DebugUI{
		recorder: recorder,
		hub:      newHub(),
		stats:    NewStats(recorder.Store()),
	}

	recorder.Store().Subscribe(func(e Entry) {
		d.hub.Broadcast(wsMessage{Kind: "message", Data: e})
	})
	recorder.LogStore().Subscribe(func(e LogEntry) {
		d.hub.Broadcast(wsMessage{Kind: "log", Data: e})
	})

	mux := http.NewServeMux()
	mux.Handle("GET /", http.FileServerFS(staticFiles()))
	mux.HandleFunc("GET /ws", d.handleWS)
	mux.HandleFunc("GET /api/messages", d.handleMessages)
	mux.HandleFunc("DELETE /api/messages", d.handleMessagesClear)
	mux.HandleFunc("GET /api/messages/search", d.handleSearch)
	mux.HandleFunc("GET /api/logs", d.handleLogs)
	mux.HandleFunc("DELETE /api/logs", d.handleLogsClear)
	mux.HandleFunc("GET /api/logs/search", d.handleLogSearch)
	mux.HandleFunc("GET /api/stats", d.handleStats)
	mux.HandleFunc("GET /api/capabilities", d.handleCapabilities)

	d.srv = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 10 * time.Second,
	}

	return d
}

// ListenAndServe binds the port synchronously, then serves in the background.
// Returns an error if the port cannot be bound. The server shuts down when ctx is cancelled.
func (d *DebugUI) ListenAndServe(ctx context.Context) error {
	ln, err := net.Listen("tcp", d.srv.Addr)
	if err != nil {
		return err
	}

	stop := make(chan struct{})
	d.stats.StartPolling(stop)

	go func() {
		<-ctx.Done()
		close(stop)
		shutCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		_ = d.srv.Shutdown(shutCtx)
	}()

	log.Printf("debugui: listening on http://%s", ln.Addr())

	go func() {
		if err := d.srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("debugui: server error: %v", err)
		}
	}()

	return nil
}

func (d *DebugUI) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	ch := d.hub.add(conn)
	defer func() {
		d.hub.remove(conn)
		_ = conn.Close()
	}()

	// Discard incoming messages (just read to detect close).
	go func() {
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				_ = conn.Close()
				return
			}
		}
	}()

	for msg := range ch {
		if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

func (d *DebugUI) handleMessages(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	entries := d.recorder.store.Entries(offset, limit)
	if entries == nil {
		entries = []Entry{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}

func (d *DebugUI) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
		return
	}

	entries := d.recorder.store.Search(q)
	if entries == nil {
		entries = []Entry{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}

func (d *DebugUI) handleLogs(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}

	entries := d.recorder.logStore.Entries(offset, limit)
	if entries == nil {
		entries = []LogEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}

func (d *DebugUI) handleMessagesClear(w http.ResponseWriter, _ *http.Request) {
	d.recorder.store.Clear()
	d.hub.Broadcast(wsMessage{Kind: "clear-messages"})
	w.WriteHeader(http.StatusNoContent)
}

func (d *DebugUI) handleLogsClear(w http.ResponseWriter, _ *http.Request) {
	d.recorder.logStore.Clear()
	d.hub.Broadcast(wsMessage{Kind: "clear-logs"})
	w.WriteHeader(http.StatusNoContent)
}

func (d *DebugUI) handleLogSearch(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte("[]"))
		return
	}

	entries := d.recorder.logStore.Search(q)
	if entries == nil {
		entries = []LogEntry{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(entries)
}

func (d *DebugUI) handleStats(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(d.stats.Snapshot())
}

func (d *DebugUI) handleCapabilities(w http.ResponseWriter, _ *http.Request) {
	data := d.recorder.capabilitiesSnapshot()

	w.Header().Set("Content-Type", "application/json")
	if data == nil {
		_, _ = w.Write([]byte("null"))
		return
	}
	_, _ = w.Write(data)
}

// staticFiles returns the filesystem for the embedded static files.
func staticFiles() fs.FS {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	return sub
}
