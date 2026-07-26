package runes

import (
	"testing"
	"unicode/utf8" // used only to check your answers, never in your solution
)

func TestFirstRune(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantRune rune
		wantSize int
	}{
		{"ascii", "abc", 'a', 1},
		{"two byte", "éx", 'é', 2},
		{"three byte", "€", '€', 3},
		{"four byte", "\U0001D11E", 0x1D11E, 4},
		{"empty", "", RuneError, 0},
		{"lone continuation byte", "\x80", RuneError, 1},
		{"invalid ff", "\xff", RuneError, 1},
		{"truncated two byte", "\xc3", RuneError, 1},
		{"truncated three byte", "\xe2\x82", RuneError, 1},
		{"bad continuation", "\xc3\x28", RuneError, 1},
		{"overlong nul", "\xc0\x80", RuneError, 1},
		{"overlong slash", "\xe0\x80\xaf", RuneError, 1},
		{"surrogate", "\xed\xa0\x80", RuneError, 1},
		{"too large", "\xf5\x80\x80\x80", RuneError, 1},
		{"nul", "\x00", 0, 1},
		{"max", "\U0010FFFF", 0x10FFFF, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, size := FirstRune(tt.in)
			if r != tt.wantRune || size != tt.wantSize {
				t.Errorf("FirstRune(%q) = %U, %d; want %U, %d", tt.in, r, size, tt.wantRune, tt.wantSize)
			}
		})
	}
}

// FuzzFirstRune compares your decoder against the standard library's on random
// bytes. Run it: go test -run xxx -fuzz FuzzFirstRune ./02-strings-bytes/runes/
func FuzzFirstRune(f *testing.F) {
	f.Add("hello")
	f.Add("\xc3\x28")
	f.Add("\U0001F600")
	f.Fuzz(func(t *testing.T, s string) {
		gotR, gotN := FirstRune(s)
		wantR, wantN := utf8.DecodeRuneInString(s)
		if gotR != wantR || gotN != wantN {
			t.Fatalf("FirstRune(%q) = %U, %d; stdlib says %U, %d", s, gotR, gotN, wantR, wantN)
		}
	})
}

func TestEncodeRune(t *testing.T) {
	for _, r := range []rune{'a', 0, 'é', '€', 0x1D11E, 0x10FFFF, -1, 0xD800, 0x110000} {
		buf := make([]byte, 4)
		n := EncodeRune(buf, r)
		want := make([]byte, 4)
		wantN := utf8.EncodeRune(want, r)
		if n != wantN || string(buf[:n]) != string(want[:wantN]) {
			t.Errorf("EncodeRune(%U) = % x (%d bytes), want % x (%d bytes)", r, buf[:n], n, want[:wantN], wantN)
		}
	}
}

func TestEncodeDecodeRoundTrip(t *testing.T) {
	buf := make([]byte, 4)
	for r := rune(0); r < 0x11000; r++ {
		if r >= 0xD800 && r <= 0xDFFF {
			continue
		}
		n := EncodeRune(buf, r)
		got, size := FirstRune(string(buf[:n]))
		if got != r || size != n {
			t.Fatalf("round trip of %U gave %U (%d bytes)", r, got, size)
		}
	}
}

func TestCountRunes(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"abc", 3},
		{"héllo", 5},
		{"日本語", 3},
		{"a\xffb", 3},
		{"\U0001F600!", 2},
	}
	for _, tt := range tests {
		if got := CountRunes(tt.in); got != tt.want {
			t.Errorf("CountRunes(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
	if n := testing.AllocsPerRun(100, func() { CountRunes("héllo wörld 日本語") }); n > 0 {
		t.Errorf("CountRunes allocated %.0f times; don't convert to []rune", n)
	}
}

func TestReverse(t *testing.T) {
	tests := []struct{ in, want string }{
		{"", ""},
		{"abc", "cba"},
		{"héllo", "olléh"},
		{"日本語", "語本日"},
		{"\U0001F600ab", "ba\U0001F600"},
	}
	for _, tt := range tests {
		if got := Reverse(tt.in); got != tt.want {
			t.Errorf("Reverse(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestReverseClusters(t *testing.T) {
	tests := []struct{ in, want string }{
		{"café", "éfac"},
		{"abc", "cba"},
		{"à́b", "bà́"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := ReverseClusters(tt.in); got != tt.want {
			t.Errorf("ReverseClusters(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestTruncateRunes(t *testing.T) {
	tests := []struct {
		in   string
		n    int
		want string
	}{
		{"hello", 10, "hello"},
		{"hello", 5, "hello"},
		{"hello", 4, "hel…"},
		{"hello", 1, "…"},
		{"hello", 0, ""},
		{"日本語です", 3, "日本…"},
		{"", 3, ""},
	}
	for _, tt := range tests {
		got := TruncateRunes(tt.in, tt.n)
		if got != tt.want {
			t.Errorf("TruncateRunes(%q, %d) = %q, want %q", tt.in, tt.n, got, tt.want)
		}
		if CountRunes(got) > tt.n {
			t.Errorf("TruncateRunes(%q, %d) = %q which is longer than %d runes", tt.in, tt.n, got, tt.n)
		}
	}
}

func TestByteToRuneIndex(t *testing.T) {
	s := "héllo"
	tests := []struct{ b, want int }{
		{0, 0}, {1, 1}, {2, 1}, {3, 2}, {4, 3}, {5, 4}, {6, -1}, {-1, -1},
	}
	for _, tt := range tests {
		if got := ByteToRuneIndex(s, tt.b); got != tt.want {
			t.Errorf("ByteToRuneIndex(%q, %d) = %d, want %d", s, tt.b, got, tt.want)
		}
	}
}
