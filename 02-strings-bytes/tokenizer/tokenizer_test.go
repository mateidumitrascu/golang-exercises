package tokenizer

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

// dribble returns one byte per Read call, so that any SplitFunc that assumes it
// gets a whole token in one go will fail.
type dribble struct {
	s string
	i int
}

func (d *dribble) Read(p []byte) (int, error) {
	if d.i >= len(d.s) {
		return 0, io.EOF
	}
	if len(p) == 0 {
		return 0, nil
	}
	p[0] = d.s[d.i]
	d.i++
	return 1, nil
}

func collect(t *testing.T, r io.Reader, split bufio.SplitFunc) ([]string, error) {
	t.Helper()
	sc := bufio.NewScanner(r)
	sc.Split(split)
	var out []string
	for sc.Scan() {
		out = append(out, sc.Text())
	}
	return out, sc.Err()
}

func TestScanWords(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"Hello, wörld! 42x", []string{"Hello", "wörld", "42x"}},
		{"  spaced   out  ", []string{"spaced", "out"}},
		{"", nil},
		{"...,,,", nil},
		{"one", []string{"one"}},
		{"日本語 text", []string{"日本語", "text"}},
		{"trailing space ", []string{"trailing", "space"}},
	}
	for _, tt := range tests {
		got, err := collect(t, strings.NewReader(tt.in), ScanWords)
		if err != nil {
			t.Fatalf("scanning %q: %v", tt.in, err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ScanWords(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestScanWordsByteAtATime(t *testing.T) {
	const in = "héllo, wörld — 42"
	want := []string{"héllo", "wörld", "42"}
	got, err := collect(t, &dribble{s: in}, ScanWords)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("byte-at-a-time = %q, want %q\n(a multi-byte rune can straddle the buffer edge)", got, want)
	}
}

func TestScanFixed(t *testing.T) {
	got, err := collect(t, strings.NewReader("abcdefgh"), ScanFixed(3))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"abc", "def", "gh"}; !reflect.DeepEqual(got, want) {
		t.Errorf("ScanFixed(3) = %q, want %q", got, want)
	}
	got, _ = collect(t, strings.NewReader("abcdef"), ScanFixed(3))
	if want := []string{"abc", "def"}; !reflect.DeepEqual(got, want) {
		t.Errorf("exact multiple = %q, want %q (no trailing empty token)", got, want)
	}
	got, _ = collect(t, strings.NewReader(""), ScanFixed(3))
	if len(got) != 0 {
		t.Errorf("empty input = %q, want nothing", got)
	}
	got, _ = collect(t, &dribble{s: "abcdefgh"}, ScanFixed(3))
	if want := []string{"abc", "def", "gh"}; !reflect.DeepEqual(got, want) {
		t.Errorf("dribbled = %q, want %q", got, want)
	}
}

func TestScanDelimited(t *testing.T) {
	tests := []struct {
		in   string
		want []string
	}{
		{"a,b,c", []string{"a", "b", "c"}},
		{"a,,b", []string{"a", "", "b"}},
		{"a,b,", []string{"a", "b"}},
		{",a", []string{"", "a"}},
		{"", nil},
		{",", []string{""}},
	}
	for _, tt := range tests {
		got, err := collect(t, strings.NewReader(tt.in), ScanDelimited(','))
		if err != nil {
			t.Fatalf("%q: %v", tt.in, err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("ScanDelimited(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func frame(payloads ...string) []byte {
	var buf bytes.Buffer
	for _, p := range payloads {
		binary.Write(&buf, binary.BigEndian, uint32(len(p)))
		buf.WriteString(p)
	}
	return buf.Bytes()
}

func TestScanLengthPrefixed(t *testing.T) {
	data := frame("hello", "", "world!")
	got, err := collect(t, bytes.NewReader(data), ScanLengthPrefixed(1024))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"hello", "", "world!"}; !reflect.DeepEqual(got, want) {
		t.Errorf("frames = %q, want %q", got, want)
	}

	got, err = collect(t, &dribble{s: string(data)}, ScanLengthPrefixed(1024))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"hello", "", "world!"}; !reflect.DeepEqual(got, want) {
		t.Errorf("dribbled frames = %q, want %q", got, want)
	}
}

func TestScanLengthPrefixedTruncated(t *testing.T) {
	data := frame("hello")
	_, err := collect(t, bytes.NewReader(data[:7]), ScanLengthPrefixed(1024))
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("truncated payload gave err = %v, want io.ErrUnexpectedEOF", err)
	}
	_, err = collect(t, bytes.NewReader(data[:2]), ScanLengthPrefixed(1024))
	if !errors.Is(err, io.ErrUnexpectedEOF) {
		t.Errorf("truncated header gave err = %v, want io.ErrUnexpectedEOF", err)
	}
}

func TestScanLengthPrefixedTooBig(t *testing.T) {
	data := frame(strings.Repeat("x", 100))
	_, err := collect(t, bytes.NewReader(data), ScanLengthPrefixed(10))
	if err == nil {
		t.Error("a frame larger than maxFrame must produce an error")
	}
}

func TestWordFreq(t *testing.T) {
	got, err := WordFreq(strings.NewReader("Go go GO rust; Rust!"))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]int{"go": 3, "rust": 2}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("WordFreq = %v, want %v", got, want)
	}
}

type errReader struct{ err error }

func (e errReader) Read([]byte) (int, error) { return 0, e.err }

func TestWordFreqPropagatesError(t *testing.T) {
	sentinel := errors.New("boom")
	_, err := WordFreq(errReader{sentinel})
	if !errors.Is(err, sentinel) {
		t.Errorf("err = %v, want %v", err, sentinel)
	}
}
