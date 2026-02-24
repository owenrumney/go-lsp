package debugui

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
)

// Tap wraps an io.ReadWriteCloser and intercepts Content-Length framed messages,
// sending copies to the Store. All bytes pass through unmodified.
type Tap struct {
	inner io.ReadWriteCloser
	store *Store

	readMu  sync.Mutex
	readBuf bytes.Buffer

	writeMu  sync.Mutex
	writeBuf bytes.Buffer
}

// NewTap creates a Tap that wraps inner and sends captured messages to store.
func NewTap(inner io.ReadWriteCloser, store *Store) *Tap {
	return &Tap{inner: inner, store: store}
}

func (t *Tap) Read(p []byte) (int, error) {
	n, err := t.inner.Read(p)
	if n > 0 {
		t.readMu.Lock()
		t.readBuf.Write(p[:n])
		t.extractMessages(&t.readBuf, "client→server")
		t.readMu.Unlock()
	}
	return n, err
}

func (t *Tap) Write(p []byte) (int, error) {
	n, err := t.inner.Write(p)
	if n > 0 {
		t.writeMu.Lock()
		t.writeBuf.Write(p[:n])
		t.extractMessages(&t.writeBuf, "server→client")
		t.writeMu.Unlock()
	}
	return n, err
}

func (t *Tap) Close() error {
	return t.inner.Close()
}

// extractMessages parses Content-Length framed messages from buf, sending complete
// messages to the store. Incomplete data is left in buf for the next call.
func (t *Tap) extractMessages(buf *bytes.Buffer, direction string) {
	for {
		data := buf.Bytes()
		contentLen, headerEnd, ok := parseContentLength(data)
		if !ok {
			return
		}
		totalLen := headerEnd + contentLen
		if len(data) < totalLen {
			return
		}
		body := make([]byte, contentLen)
		copy(body, data[headerEnd:totalLen])
		// Consume the parsed message from the buffer.
		buf.Next(totalLen)
		t.store.Add(direction, body)
	}
}

// parseContentLength scans for "Content-Length: N\r\n\r\n" in data.
// Returns the content length, the byte offset where the body starts, and whether parsing succeeded.
func parseContentLength(data []byte) (contentLen int, headerEnd int, ok bool) {
	// Find the end of headers (double CRLF).
	idx := bytes.Index(data, []byte("\r\n\r\n"))
	if idx < 0 {
		return 0, 0, false
	}
	headerEnd = idx + 4

	// Parse headers.
	headers := string(data[:idx])
	for _, line := range strings.Split(headers, "\r\n") {
		if val, found := strings.CutPrefix(line, "Content-Length:"); found {
			val = strings.TrimSpace(val)
			n, err := strconv.Atoi(val)
			if err != nil {
				return 0, 0, false
			}
			contentLen = n
		}
	}
	if contentLen == 0 {
		return 0, 0, false
	}
	return contentLen, headerEnd, true
}

// FormatMessage wraps a JSON body in a Content-Length framed message.
// Useful for testing.
func FormatMessage(body []byte) []byte {
	return []byte(fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(body), body))
}
