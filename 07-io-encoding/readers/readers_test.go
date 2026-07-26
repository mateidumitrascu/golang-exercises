package readers

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

// dribbleReader returns at most one byte per call - a legal Reader that breaks
// any code assuming Read fills the buffer.
type dribbleReader struct {
	data string
	i    int
	// eofWithData makes the last read return (n>0, io.EOF) together, which is
	// legal and which trips up most naive loops.
	eofWithData bool
}

func (d *dribbleReader) Read(p []byte) (int, error) {
	if d.i >= len(d.data) {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}
	p[0] = d.data[d.i]
	d.i++
	if d.eofWithData && d.i == len(d.data) {
		return 1, io.EOF
	}
	return 1, nil
}

// shortWriter accepts at most limit bytes per Write and reports the truth.
type shortWriter struct {
	buf   bytes.Buffer
	limit int
}

func (s *shortWriter) Write(p []byte) (int, error) {
	if len(p) > s.limit {
		p = p[:s.limit]
	}
	return s.buf.Write(p)
}

type errWriter struct{ err error }

func (e errWriter) Write(p []byte) (int, error) { return 0, e.err }

func TestCountingWriter(t *testing.T) {
	var buf bytes.Buffer
	c := NewCountingWriter(&buf)
	n, err := c.Write([]byte("hello"))
	if n != 5 || err != nil {
		t.Fatalf("Write = %d, %v", n, err)
	}
	io.WriteString(c, " world")
	if c.N() != 11 {
		t.Errorf("N = %d, want 11", c.N())
	}
	if buf.String() != "hello world" {
		t.Errorf("underlying = %q", buf.String())
	}

	sw := &shortWriter{limit: 3}
	c2 := NewCountingWriter(sw)
	n, _ = c2.Write([]byte("hello"))
	if n != 3 || c2.N() != 3 {
		t.Errorf("short write: n = %d, N = %d; want 3, 3 (count what landed)", n, c2.N())
	}

	boom := errors.New("boom")
	c3 := NewCountingWriter(errWriter{boom})
	if _, err := c3.Write([]byte("x")); !errors.Is(err, boom) {
		t.Errorf("err = %v, want boom", err)
	}
}

func TestUpperReader(t *testing.T) {
	got, err := io.ReadAll(UpperReader(strings.NewReader("Hello, Wörld 123!")))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "HELLO, WöRLD 123!" {
		t.Errorf("= %q", got)
	}
	// One byte at a time must give the same answer.
	got, err = io.ReadAll(UpperReader(&dribbleReader{data: "abc", eofWithData: true}))
	if err != nil || string(got) != "ABC" {
		t.Errorf("dribbled = %q, %v", got, err)
	}
}

func TestExactReader(t *testing.T) {
	got, err := io.ReadAll(ExactReader(strings.NewReader("hello world"), 5))
	if err != nil || string(got) != "hello" {
		t.Errorf("= %q, %v", got, err)
	}
	_, err = io.ReadAll(ExactReader(strings.NewReader("abc"), 10))
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("truncated stream gave %v, want io.ErrUnexpectedEOF", err)
	}
	got, err = io.ReadAll(ExactReader(&dribbleReader{data: "abcdef"}, 3))
	if err != nil || string(got) != "abc" {
		t.Errorf("dribbled = %q, %v", got, err)
	}
	got, err = io.ReadAll(ExactReader(strings.NewReader("abc"), 0))
	if err != nil || len(got) != 0 {
		t.Errorf("n=0 gave %q, %v", got, err)
	}
}

func TestTeeReader(t *testing.T) {
	var side bytes.Buffer
	r := TeeReader(strings.NewReader("hello"), &side)
	got, err := io.ReadAll(r)
	if err != nil || string(got) != "hello" {
		t.Fatalf("= %q, %v", got, err)
	}
	if side.String() != "hello" {
		t.Errorf("side channel got %q", side.String())
	}
	boom := errors.New("boom")
	_, err = io.ReadAll(TeeReader(strings.NewReader("hello"), errWriter{boom}))
	if !errors.Is(err, boom) {
		t.Errorf("a write error must surface as a read error, got %v", err)
	}
}

func TestMultiWriter(t *testing.T) {
	var a, b bytes.Buffer
	w := MultiWriter(&a, &b)
	n, err := io.WriteString(w, "hello")
	if n != 5 || err != nil {
		t.Fatalf("= %d, %v", n, err)
	}
	if a.String() != "hello" || b.String() != "hello" {
		t.Errorf("a = %q, b = %q", a.String(), b.String())
	}

	boom := errors.New("boom")
	var c bytes.Buffer
	w = MultiWriter(&c, errWriter{boom}, &a)
	if _, err := io.WriteString(w, "x"); !errors.Is(err, boom) {
		t.Errorf("err = %v, want boom", err)
	}

	w = MultiWriter(&shortWriter{limit: 2})
	if _, err := io.WriteString(w, "hello"); !errors.Is(err, io.ErrShortWrite) {
		t.Errorf("err = %v, want io.ErrShortWrite", err)
	}
}

func TestReadAll(t *testing.T) {
	got, err := ReadAll(strings.NewReader("hello world"))
	if err != nil || string(got) != "hello world" {
		t.Errorf("= %q, %v", got, err)
	}
	got, err = ReadAll(strings.NewReader(""))
	if err != nil || len(got) != 0 {
		t.Errorf("empty = %q, %v", got, err)
	}
	// The last read returns data and io.EOF together.
	got, err = ReadAll(&dribbleReader{data: "abcdef", eofWithData: true})
	if err != nil {
		t.Fatalf("err = %v; io.EOF is not a failure, and the bytes that came with "+
			"it must not be dropped", err)
	}
	if string(got) != "abcdef" {
		t.Errorf("= %q, want abcdef", got)
	}
	big := strings.Repeat("x", 100000)
	got, err = ReadAll(strings.NewReader(big))
	if err != nil || len(got) != len(big) {
		t.Errorf("large read: %d bytes, %v", len(got), err)
	}
	boom := errors.New("boom")
	if _, err := ReadAll(io.MultiReader(strings.NewReader("ab"), errReader{boom})); !errors.Is(err, boom) {
		t.Errorf("err = %v, want boom", err)
	}
}

type errReader struct{ err error }

func (e errReader) Read([]byte) (int, error) { return 0, e.err }

func TestCopy(t *testing.T) {
	var dst bytes.Buffer
	n, err := Copy(&dst, strings.NewReader("hello"), nil)
	if n != 5 || err != nil || dst.String() != "hello" {
		t.Errorf("= %d, %v, %q", n, err, dst.String())
	}
	dst.Reset()
	n, err = Copy(&dst, &dribbleReader{data: "hello", eofWithData: true}, make([]byte, 2))
	if n != 5 || err != nil || dst.String() != "hello" {
		t.Errorf("small buffer: = %d, %v, %q", n, err, dst.String())
	}
	sw := &shortWriter{limit: 1}
	if _, err := Copy(sw, strings.NewReader("hello"), make([]byte, 4)); !errors.Is(err, io.ErrShortWrite) {
		t.Errorf("err = %v, want io.ErrShortWrite", err)
	}
}

func TestCopyProgress(t *testing.T) {
	var dst bytes.Buffer
	var updates []int64
	n, err := CopyProgress(&dst, &dribbleReader{data: "abcde"}, func(total int64) {
		updates = append(updates, total)
	})
	if n != 5 || err != nil {
		t.Fatalf("= %d, %v", n, err)
	}
	if len(updates) == 0 {
		t.Fatal("progress was never called")
	}
	if updates[len(updates)-1] != 5 {
		t.Errorf("final progress = %d, want 5", updates[len(updates)-1])
	}
	for i := 1; i < len(updates); i++ {
		if updates[i] <= updates[i-1] {
			t.Errorf("progress must increase: %v", updates)
			break
		}
	}
}

func TestLineWriter(t *testing.T) {
	var buf bytes.Buffer
	lw := NewLineWriter(&buf)

	io.WriteString(lw, "hel")
	if buf.Len() != 0 {
		t.Errorf("a partial line was passed through: %q", buf.String())
	}
	io.WriteString(lw, "lo\nwor")
	if buf.String() != "hello\n" {
		t.Errorf("= %q, want %q", buf.String(), "hello\n")
	}
	io.WriteString(lw, "ld\nand\nmore")
	if buf.String() != "hello\nworld\nand\n" {
		t.Errorf("= %q", buf.String())
	}
	if err := lw.Flush(); err != nil {
		t.Fatal(err)
	}
	if buf.String() != "hello\nworld\nand\nmore" {
		t.Errorf("after Flush = %q", buf.String())
	}

	// Write must report the number of bytes it accepted from p, not the number
	// it passed through.
	var b2 bytes.Buffer
	lw2 := NewLineWriter(&b2)
	n, err := lw2.Write([]byte("abc"))
	if n != 3 || err != nil {
		t.Errorf("Write = %d, %v; want 3, nil", n, err)
	}
}
