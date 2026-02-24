package debugui

import (
	"bytes"
	"io"
	"testing"
)

type bufRWC struct {
	*bytes.Buffer
}

func (b bufRWC) Close() error { return nil }

func TestTapRead(t *testing.T) {
	store := NewStore()

	msg1 := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}`
	msg2 := `{"jsonrpc":"2.0","method":"initialized","params":{}}`
	input := FormatMessage([]byte(msg1))
	input = append(input, FormatMessage([]byte(msg2))...)

	inner := bufRWC{bytes.NewBuffer(input)}
	tap := NewTap(inner, store)

	out, err := io.ReadAll(tap)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(out, input) {
		t.Error("tap altered read data")
	}

	entries := store.Entries(0, 10)
	if len(entries) != 2 {
		t.Fatalf("got %d entries, want 2", len(entries))
	}
	if entries[0].Direction != "client→server" {
		t.Errorf("direction = %q, want client→server", entries[0].Direction)
	}
	if entries[0].Method != "initialize" {
		t.Errorf("method = %q, want initialize", entries[0].Method)
	}
	if entries[1].Method != "initialized" {
		t.Errorf("method = %q, want initialized", entries[1].Method)
	}
}

func TestTapWrite(t *testing.T) {
	store := NewStore()

	msg := `{"jsonrpc":"2.0","id":1,"result":{}}`
	framed := FormatMessage([]byte(msg))

	var outBuf bytes.Buffer
	inner := struct {
		io.Reader
		io.Writer
		io.Closer
	}{
		Reader: bytes.NewReader(nil),
		Writer: &outBuf,
		Closer: io.NopCloser(nil),
	}
	tap := NewTap(rwc{inner.Reader, inner.Writer, inner.Closer}, store)

	_, err := tap.Write(framed)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(outBuf.Bytes(), framed) {
		t.Error("tap altered write data")
	}

	entries := store.Entries(0, 10)
	if len(entries) != 1 {
		t.Fatalf("got %d entries, want 1", len(entries))
	}
	if entries[0].Direction != "server→client" {
		t.Errorf("direction = %q, want server→client", entries[0].Direction)
	}
	if entries[0].MsgType != "response" {
		t.Errorf("msgType = %q, want response", entries[0].MsgType)
	}
}

func TestTapPartialRead(t *testing.T) {
	store := NewStore()

	msg := `{"jsonrpc":"2.0","id":1,"method":"test","params":{}}`
	framed := FormatMessage([]byte(msg))

	// Split the message into two halves to simulate partial reads.
	half := len(framed) / 2
	part1 := framed[:half]
	part2 := framed[half:]

	inner := bufRWC{bytes.NewBuffer(nil)}
	tap := NewTap(inner, store)

	// Write first half — no complete message yet.
	inner.Buffer.Write(part1)
	buf := make([]byte, 4096)
	n, _ := tap.Read(buf)
	if n != len(part1) {
		t.Fatalf("read %d bytes, want %d", n, len(part1))
	}

	entries := store.Entries(0, 10)
	if len(entries) != 0 {
		t.Fatalf("got %d entries after partial read, want 0", len(entries))
	}

	// Write second half — now the message is complete.
	inner.Buffer.Write(part2)
	n, _ = tap.Read(buf)
	if n != len(part2) {
		t.Fatalf("read %d bytes, want %d", n, len(part2))
	}

	entries = store.Entries(0, 10)
	if len(entries) != 1 {
		t.Fatalf("got %d entries after full read, want 1", len(entries))
	}
}

// rwc combines separate Reader, Writer, Closer into an io.ReadWriteCloser.
type rwc struct {
	io.Reader
	io.Writer
	io.Closer
}
