package framing

import (
	"bytes"
	"encoding/binary"
	"errors"
	"hash/crc32"
	"io"
	"reflect"
	"strings"
	"testing"
)

func mustFrame(t *testing.T, f Frame) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := WriteFrame(&buf, f); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestWriteFrameLayout(t *testing.T) {
	got := mustFrame(t, Frame{Type: 7, Payload: []byte("hi")})
	want := []byte{Magic0, Magic1, 7, 0, 0, 0, 2, 'h', 'i'}
	sum := crc32.ChecksumIEEE([]byte("hi"))
	want = binary.BigEndian.AppendUint32(want, sum)
	if !bytes.Equal(got, want) {
		t.Errorf("frame = % x\nwant     % x", got, want)
	}
	if n := len(mustFrame(t, Frame{Type: 1})); n != HeaderLen+4 {
		t.Errorf("empty payload frame is %d bytes, want %d", n, HeaderLen+4)
	}
}

type countingWriter struct {
	bytes.Buffer
	writes int
}

func (c *countingWriter) Write(p []byte) (int, error) {
	c.writes++
	return c.Buffer.Write(p)
}

func TestWriteFrameIsOneWrite(t *testing.T) {
	var cw countingWriter
	WriteFrame(&cw, Frame{Type: 1, Payload: bytes.Repeat([]byte("x"), 100)})
	if cw.writes > 2 {
		t.Errorf("WriteFrame made %d Write calls; build the frame then write it once", cw.writes)
	}
}

func TestRoundTrip(t *testing.T) {
	frames := []Frame{
		{Type: 1, Payload: []byte("hello")},
		{Type: 2, Payload: nil},
		{Type: 255, Payload: bytes.Repeat([]byte{0}, 1000)},
	}
	var buf bytes.Buffer
	for _, f := range frames {
		if err := WriteFrame(&buf, f); err != nil {
			t.Fatal(err)
		}
	}
	for i, want := range frames {
		got, err := ReadFrame(&buf, 1<<20)
		if err != nil {
			t.Fatalf("frame %d: %v", i, err)
		}
		if got.Type != want.Type || !bytes.Equal(got.Payload, want.Payload) {
			t.Errorf("frame %d = %+v, want %+v", i, got, want)
		}
	}
	if _, err := ReadFrame(&buf, 1<<20); !errors.Is(err, io.EOF) {
		t.Errorf("at the end of the stream: %v, want io.EOF", err)
	}
}

func TestReadFrameErrors(t *testing.T) {
	good := mustFrame(t, Frame{Type: 1, Payload: []byte("hello")})

	tests := []struct {
		name string
		data []byte
		max  int
		want error
	}{
		{"truncated header", good[:3], 1 << 20, io.ErrUnexpectedEOF},
		{"truncated payload", good[:9], 1 << 20, io.ErrUnexpectedEOF},
		{"truncated checksum", good[:len(good)-2], 1 << 20, io.ErrUnexpectedEOF},
		{"bad magic", append([]byte{0, 0}, good[2:]...), 1 << 20, ErrBadMagic},
		{"too big", good, 2, ErrFrameTooBig},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ReadFrame(bytes.NewReader(tt.data), tt.max)
			if !errors.Is(err, tt.want) {
				t.Errorf("err = %v, want %v", err, tt.want)
			}
		})
	}

	corrupt := append([]byte(nil), good...)
	corrupt[8] ^= 0xFF // flip a payload bit
	if _, err := ReadFrame(bytes.NewReader(corrupt), 1<<20); !errors.Is(err, ErrChecksum) {
		t.Errorf("corrupt payload gave %v, want ErrChecksum", err)
	}
}

func TestReadFrameDoesNotTrustLength(t *testing.T) {
	// A hostile 4 GB length must be rejected without allocating 4 GB.
	hostile := []byte{Magic0, Magic1, 1, 0xFF, 0xFF, 0xFF, 0xFF}
	_, err := ReadFrame(bytes.NewReader(hostile), 1024)
	if !errors.Is(err, ErrFrameTooBig) {
		t.Errorf("err = %v, want ErrFrameTooBig", err)
	}
}

func TestEncoderDecoder(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	for i := range 5 {
		if err := enc.Encode(Frame{Type: uint8(i), Payload: []byte(strings.Repeat("a", i))}); err != nil {
			t.Fatal(err)
		}
	}
	if err := enc.Flush(); err != nil {
		t.Fatal(err)
	}

	dec := NewDecoder(&buf, 1<<20)
	n := 0
	for f := range dec.All {
		if f.Type != uint8(n) || len(f.Payload) != n {
			t.Errorf("frame %d = %+v", n, f)
		}
		n++
	}
	if err := dec.Err(); err != nil {
		t.Errorf("Err = %v, want nil at a clean EOF", err)
	}
	if n != 5 {
		t.Errorf("decoded %d frames, want 5", n)
	}
}

func TestDecoderIterationStopsOnError(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.Encode(Frame{Type: 1, Payload: []byte("ok")})
	enc.Flush()
	buf.Write([]byte{0x00, 0x00, 0x00}) // garbage

	dec := NewDecoder(&buf, 1<<20)
	n := 0
	for range dec.All {
		n++
	}
	if n != 1 {
		t.Errorf("yielded %d frames, want 1", n)
	}
	if dec.Err() == nil {
		t.Error("Err must report why iteration stopped")
	}
}

func TestUvarint(t *testing.T) {
	tests := []struct {
		x    uint64
		want []byte
	}{
		{0, []byte{0x00}},
		{1, []byte{0x01}},
		{127, []byte{0x7F}},
		{128, []byte{0x80, 0x01}},
		{300, []byte{0xAC, 0x02}},
		{16384, []byte{0x80, 0x80, 0x01}},
		{^uint64(0), []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01}},
	}
	for _, tt := range tests {
		got := AppendUvarint(nil, tt.x)
		if !bytes.Equal(got, tt.want) {
			t.Errorf("AppendUvarint(%d) = % x, want % x", tt.x, got, tt.want)
		}
		x, n := ConsumeUvarint(got)
		if x != tt.x || n != len(tt.want) {
			t.Errorf("ConsumeUvarint(% x) = %d, %d; want %d, %d", got, x, n, tt.x, len(tt.want))
		}
	}
	if _, n := ConsumeUvarint(nil); n != 0 {
		t.Errorf("empty input: n = %d, want 0", n)
	}
	if _, n := ConsumeUvarint([]byte{0x80}); n != 0 {
		t.Errorf("truncated varint: n = %d, want 0", n)
	}
	if _, n := ConsumeUvarint(bytes.Repeat([]byte{0xFF}, 11)); n >= 0 {
		t.Errorf("overflowing varint: n = %d, want negative", n)
	}
	// Appending must extend rather than replace.
	buf := AppendUvarint([]byte{0xAA}, 1)
	if !bytes.Equal(buf, []byte{0xAA, 0x01}) {
		t.Errorf("append = % x", buf)
	}
}

func TestVarintStrings(t *testing.T) {
	buf := AppendString(nil, "hello")
	buf = AppendString(buf, "")
	buf = AppendString(buf, strings.Repeat("x", 200))

	var got []string
	for len(buf) > 0 {
		s, n, ok := ConsumeString(buf)
		if !ok {
			t.Fatalf("ConsumeString failed with % x left", buf)
		}
		got = append(got, s)
		buf = buf[n:]
	}
	want := []string{"hello", "", strings.Repeat("x", 200)}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("= %q", got)
	}
	if _, _, ok := ConsumeString([]byte{0x05, 'a'}); ok {
		t.Error("a truncated string must report ok = false")
	}
}

func FuzzFrameRoundTrip(f *testing.F) {
	f.Add(uint8(1), []byte("hello"))
	f.Add(uint8(0), []byte(""))
	f.Fuzz(func(t *testing.T, typ uint8, payload []byte) {
		var buf bytes.Buffer
		if err := WriteFrame(&buf, Frame{Type: typ, Payload: payload}); err != nil {
			t.Fatal(err)
		}
		got, err := ReadFrame(&buf, 1<<20)
		if err != nil {
			t.Fatalf("read back: %v", err)
		}
		if got.Type != typ || !bytes.Equal(got.Payload, payload) {
			t.Fatalf("round trip changed the frame")
		}
	})
}

func FuzzReadFrameDoesNotPanic(f *testing.F) {
	f.Add([]byte{Magic0, Magic1, 1, 0, 0, 0, 2, 'h', 'i'})
	f.Fuzz(func(t *testing.T, data []byte) {
		// Any input at all: an error is fine, a panic is not.
		ReadFrame(bytes.NewReader(data), 4096)
	})
}
