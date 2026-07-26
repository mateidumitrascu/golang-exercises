package wordwrap

import (
	"strings"
	"testing"
)

func TestWrap(t *testing.T) {
	tests := []struct {
		name  string
		in    string
		width int
		want  string
	}{
		{
			"simple",
			"the quick brown fox jumps over the lazy dog",
			10,
			"the quick\nbrown fox\njumps over\nthe lazy\ndog",
		},
		{
			"exact fit",
			"aaa bbb",
			7,
			"aaa bbb",
		},
		{
			"long word gets its own line",
			"hi supercalifragilistic bye",
			5,
			"hi\nsupercalifragilistic\nbye",
		},
		{"collapses whitespace", "  a \t\n b  ", 10, "a b"},
		{"empty", "", 10, ""},
		{"only whitespace", "  \n\t ", 10, ""},
		{"unicode counts runes", "héllo wörld", 5, "héllo\nwörld"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Wrap(tt.in, tt.width)
			if got != tt.want {
				t.Errorf("Wrap(%q, %d) =\n%q\nwant\n%q", tt.in, tt.width, got, tt.want)
			}
			for _, line := range strings.Split(got, "\n") {
				if strings.HasSuffix(line, " ") || strings.HasPrefix(line, " ") {
					t.Errorf("line %q has stray spaces", line)
				}
			}
		})
	}
}

func TestWrapPanicsOnBadWidth(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Wrap(s, 0) did not panic")
		}
	}()
	Wrap("x", 0)
}

func TestWrapHard(t *testing.T) {
	got := WrapHard("hi abcdefghij bye", 4)
	want := "hi\nabcd\nefgh\nij\nbye"
	if got != want {
		t.Errorf("WrapHard =\n%q\nwant\n%q", got, want)
	}
	if got := WrapHard("日本語です", 2); got != "日本\n語で\nす" {
		t.Errorf("WrapHard(unicode) = %q", got)
	}
}

func TestWrapParagraphs(t *testing.T) {
	in := "one two three\n\nfour five six seven\n\n\neight"
	got := WrapParagraphs(in, 9)
	want := "one two\nthree\n\nfour five\nsix seven\n\neight"
	if got != want {
		t.Errorf("WrapParagraphs =\n%q\nwant\n%q", got, want)
	}
}

func TestIndent(t *testing.T) {
	got := Indent("a\n\nb", "> ")
	want := "> a\n\n> b"
	if got != want {
		t.Errorf("Indent = %q, want %q", got, want)
	}
	if got := Indent("", "  "); got != "" {
		t.Errorf("Indent(\"\") = %q", got)
	}
}

func TestDedent(t *testing.T) {
	in := "    line one\n      line two\n\n    line three\n"
	want := "line one\n  line two\n\nline three\n"
	if got := Dedent(in); got != want {
		t.Errorf("Dedent =\n%q\nwant\n%q", got, want)
	}
	// Mixed indentation: the common prefix is the shorter one.
	if got := Dedent("\t\ta\n\tb"); got != "\ta\nb" {
		t.Errorf("Dedent(tabs) = %q, want %q", got, "\ta\nb")
	}
	// A line with no indentation means nothing can be removed.
	if got := Dedent("a\n    b"); got != "a\n    b" {
		t.Errorf("Dedent = %q, want unchanged", got)
	}
}

func BenchmarkWrap(b *testing.B) {
	text := strings.Repeat("the quick brown fox jumps over the lazy dog ", 200)
	b.ReportAllocs()
	for b.Loop() {
		_ = Wrap(text, 72)
	}
}
