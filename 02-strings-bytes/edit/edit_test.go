package edit

import (
	"math"
	"strings"
	"testing"
)

func TestDistance(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "", 0},
		{"abc", "abc", 0},
		{"", "abc", 3},
		{"abc", "", 3},
		{"kitten", "sitting", 3},
		{"flaw", "lawn", 2},
		{"gumbo", "gambol", 2},
		{"café", "cafe", 1},
		{"日本", "日本語", 1},
		{"abc", "cba", 2},
	}
	for _, tt := range tests {
		if got := Distance(tt.a, tt.b); got != tt.want {
			t.Errorf("Distance(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
		}
		if got := Distance(tt.b, tt.a); got != tt.want {
			t.Errorf("Distance is not symmetric: (%q, %q) = %d, want %d", tt.b, tt.a, got, tt.want)
		}
	}
}

func TestDistanceMemory(t *testing.T) {
	a := strings.Repeat("ab", 400)
	b := strings.Repeat("ba", 400)
	// Two rows of 800 ints is ~4 allocations' worth; a full table is 800.
	if n := testing.AllocsPerRun(5, func() { Distance(a, b) }); n > 10 {
		t.Errorf("Distance made %.0f allocations; keep only two rows", n)
	}
}

func TestDistanceAtMost(t *testing.T) {
	if got := DistanceAtMost("kitten", "sitting", 10); got != 3 {
		t.Errorf("DistanceAtMost with a loose bound = %d, want 3", got)
	}
	if got := DistanceAtMost("kitten", "sitting", 3); got != 3 {
		t.Errorf("DistanceAtMost at the bound = %d, want 3", got)
	}
	if got := DistanceAtMost("kitten", "sitting", 2); got <= 2 {
		t.Errorf("DistanceAtMost(..., 2) = %d, want > 2", got)
	}
	if got := DistanceAtMost(strings.Repeat("a", 1000), strings.Repeat("b", 1000), 1); got <= 1 {
		t.Errorf("got %d, want > 1", got)
	}
}

func TestSimilarity(t *testing.T) {
	tests := []struct {
		a, b string
		want float64
	}{
		{"", "", 1},
		{"abc", "abc", 1},
		{"abcd", "abce", 0.75},
		{"abc", "xyz", 0},
		{"", "ab", 0},
	}
	for _, tt := range tests {
		if got := Similarity(tt.a, tt.b); math.Abs(got-tt.want) > 1e-9 {
			t.Errorf("Similarity(%q, %q) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func isSubsequence(sub, s string) bool {
	rs := []rune(s)
	i := 0
	for _, r := range sub {
		for i < len(rs) && rs[i] != r {
			i++
		}
		if i == len(rs) {
			return false
		}
		i++
	}
	return true
}

func TestLCS(t *testing.T) {
	tests := []struct {
		a, b    string
		wantLen int
	}{
		{"ABCBDAB", "BDCABA", 4},
		{"abc", "abc", 3},
		{"abc", "xyz", 0},
		{"", "abc", 0},
		{"AGGTAB", "GXTXAYB", 4},
	}
	for _, tt := range tests {
		got := LCS(tt.a, tt.b)
		if len([]rune(got)) != tt.wantLen {
			t.Errorf("LCS(%q, %q) = %q (len %d), want length %d", tt.a, tt.b, got, len([]rune(got)), tt.wantLen)
		}
		if !isSubsequence(got, tt.a) || !isSubsequence(got, tt.b) {
			t.Errorf("LCS(%q, %q) = %q, which is not a subsequence of both", tt.a, tt.b, got)
		}
	}
}

func TestScript(t *testing.T) {
	pairs := [][2]string{
		{"kitten", "sitting"},
		{"", "abc"},
		{"abc", ""},
		{"same", "same"},
		{"café", "cafés"},
		{"the quick fox", "a quick brown fox"},
	}
	for _, p := range pairs {
		ops := Script(p[0], p[1])
		if got := Apply(p[0], ops); got != p[1] {
			t.Errorf("Apply(%q, Script(%q, %q)) = %q, want %q", p[0], p[0], p[1], got, p[1])
		}
		edits := 0
		for _, op := range ops {
			if op.Kind != '=' {
				edits++
			}
			if op.Kind != '=' && op.Kind != '+' && op.Kind != '-' && op.Kind != '~' {
				t.Fatalf("unknown op kind %q", op.Kind)
			}
		}
		if want := Distance(p[0], p[1]); edits != want {
			t.Errorf("Script(%q, %q) has %d edits, want %d (minimal)", p[0], p[1], edits, want)
		}
	}
}

func FuzzScriptApplies(f *testing.F) {
	f.Add("kitten", "sitting")
	f.Add("", "x")
	f.Fuzz(func(t *testing.T, a, b string) {
		if len(a) > 64 || len(b) > 64 {
			t.Skip()
		}
		if got := Apply(a, Script(a, b)); got != b {
			t.Fatalf("Apply(%q, Script(%q,%q)) = %q, want %q", a, a, b, got, b)
		}
	})
}

func BenchmarkDistance(b *testing.B) {
	x, y := strings.Repeat("lorem ipsum ", 20), strings.Repeat("lorem gypsum ", 20)
	b.ReportAllocs()
	for b.Loop() {
		Distance(x, y)
	}
}
