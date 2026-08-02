package debugui

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAllowLoopbackOrigin(t *testing.T) {
	tests := []struct {
		origin string
		want   bool
	}{
		{"", true}, // non-browser clients send no Origin
		{"http://localhost:7100", true},
		{"http://127.0.0.1:7100", true},
		{"http://[::1]:7100", true},
		{"http://evil.example.com", false},
		{"http://192.168.1.20:7100", false},
		{"not a url\x7f://", false},
	}
	for _, tt := range tests {
		r := httptest.NewRequest(http.MethodGet, "/ws", nil)
		if tt.origin != "" {
			r.Header.Set("Origin", tt.origin)
		}
		if got := allowLoopbackOrigin(r); got != tt.want {
			t.Errorf("allowLoopbackOrigin(origin=%q) = %v, want %v", tt.origin, got, tt.want)
		}
	}
}

func TestNewNormalizesWildcardAddrToLoopback(t *testing.T) {
	d := New(":0", NewRecorder())
	if d.srv.Addr != "127.0.0.1:0" {
		t.Errorf("addr = %q, want 127.0.0.1:0", d.srv.Addr)
	}

	d = New("0.0.0.0:7100", NewRecorder())
	if d.srv.Addr != "0.0.0.0:7100" {
		t.Errorf("explicit host must be preserved, got %q", d.srv.Addr)
	}
}
