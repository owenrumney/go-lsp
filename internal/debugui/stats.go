package debugui

import (
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// StatsSnapshot is the JSON-serializable runtime stats returned by the API.
type StatsSnapshot struct {
	Uptime        string `json:"uptime"`
	HeapAlloc     uint64 `json:"heapAlloc"`
	HeapSys       uint64 `json:"heapSys"`
	NumGC         uint32 `json:"numGC"`
	NumGoroutine  int    `json:"numGoroutine"`
	Requests      int64  `json:"requests"`
	Responses     int64  `json:"responses"`
	Notifications int64  `json:"notifications"`
	AvgLatencyMs  float64 `json:"avgLatencyMs"`
}

// Stats collects runtime and message-level metrics for the debug UI.
type Stats struct {
	startTime time.Time
	store     *Store

	// Runtime stats (updated by background goroutine).
	mu           sync.RWMutex
	heapAlloc    uint64
	heapSys      uint64
	numGC        uint32
	numGoroutine int

	// Message counters (updated via Store subscriber).
	requests      atomic.Int64
	responses     atomic.Int64
	notifications atomic.Int64

	// Latency tracking.
	latencyMu    sync.Mutex
	latencySum   float64
	latencyCount int64
}

// NewStats creates a Stats collector that subscribes to the store for
// message-level metrics. Call StartPolling to begin sampling runtime stats.
func NewStats(store *Store) *Stats {
	s := &Stats{
		startTime: time.Now(),
		store:     store,
	}

	store.Subscribe(s.onEntry)

	return s
}

// StartPolling begins the background goroutine that samples runtime stats.
// It samples immediately, then every 2 seconds until stop is closed.
func (s *Stats) StartPolling(stop <-chan struct{}) {
	s.sampleRuntime()
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.sampleRuntime()
			case <-stop:
				return
			}
		}
	}()
}

func (s *Stats) sampleRuntime() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	s.mu.Lock()
	s.heapAlloc = m.HeapAlloc
	s.heapSys = m.HeapSys
	s.numGC = m.NumGC
	s.numGoroutine = runtime.NumGoroutine()
	s.mu.Unlock()
}

func (s *Stats) onEntry(e Entry) {
	switch e.MsgType {
	case "request":
		s.requests.Add(1)
	case "response":
		s.responses.Add(1)
		if e.PairedWith >= 0 {
			if req := s.store.Entry(e.PairedWith); req != nil {
				ms := float64(e.Timestamp.Sub(req.Timestamp)) / float64(time.Millisecond)
				s.latencyMu.Lock()
				s.latencySum += ms
				s.latencyCount++
				s.latencyMu.Unlock()
			}
		}
	case "notification":
		s.notifications.Add(1)
	}
}

// Snapshot returns the current stats as a JSON-serializable struct.
func (s *Stats) Snapshot() StatsSnapshot {
	s.mu.RLock()
	heapAlloc := s.heapAlloc
	heapSys := s.heapSys
	numGC := s.numGC
	numGoroutine := s.numGoroutine
	s.mu.RUnlock()

	var avgLatency float64
	s.latencyMu.Lock()
	if s.latencyCount > 0 {
		avgLatency = s.latencySum / float64(s.latencyCount)
	}
	s.latencyMu.Unlock()

	return StatsSnapshot{
		Uptime:        time.Since(s.startTime).Truncate(time.Second).String(),
		HeapAlloc:     heapAlloc,
		HeapSys:       heapSys,
		NumGC:         numGC,
		NumGoroutine:  numGoroutine,
		Requests:      s.requests.Load(),
		Responses:     s.responses.Load(),
		Notifications: s.notifications.Load(),
		AvgLatencyMs:  avgLatency,
	}
}
