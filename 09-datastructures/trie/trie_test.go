package trie

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func build(words ...string) *Trie {
	t := New()
	for _, w := range words {
		t.Insert(w)
	}
	return t
}

func TestInsertContains(t *testing.T) {
	tr := New()
	if tr.Len() != 0 || tr.Contains("") {
		t.Error("a fresh trie should be empty")
	}
	if !tr.Insert("go") {
		t.Error("Insert of a new word should report true")
	}
	if tr.Insert("go") {
		t.Error("Insert of an existing word should report false")
	}
	tr.Insert("gopher")
	tr.Insert("rust")
	if tr.Len() != 3 {
		t.Errorf("Len = %d, want 3", tr.Len())
	}
	for _, w := range []string{"go", "gopher", "rust"} {
		if !tr.Contains(w) {
			t.Errorf("Contains(%q) = false", w)
		}
	}
	for _, w := range []string{"g", "gop", "gophers", "", "python"} {
		if tr.Contains(w) {
			t.Errorf("Contains(%q) = true, want false", w)
		}
	}
	tr.Insert("")
	if !tr.Contains("") || tr.Len() != 4 {
		t.Error("the empty string is a valid word")
	}
}

func TestUnicode(t *testing.T) {
	tr := build("日本", "日本語", "héllo")
	if !tr.Contains("日本語") || !tr.Contains("héllo") {
		t.Error("multi-byte words must work")
	}
	if !tr.HasPrefix("日") {
		t.Error("HasPrefix on a rune boundary")
	}
	if tr.Contains("日") {
		t.Error("a prefix is not a word")
	}
}

func TestHasPrefix(t *testing.T) {
	tr := build("go", "gopher")
	for _, p := range []string{"", "g", "go", "gop", "gopher"} {
		if !tr.HasPrefix(p) {
			t.Errorf("HasPrefix(%q) = false", p)
		}
	}
	for _, p := range []string{"gophers", "x", "og"} {
		if tr.HasPrefix(p) {
			t.Errorf("HasPrefix(%q) = true", p)
		}
	}
	if !New().HasPrefix("") {
		t.Error(`HasPrefix("") on an empty trie should be true`)
	}
}

func TestDeletePrunes(t *testing.T) {
	tr := build("go")
	base := tr.Nodes()
	tr.Insert("gopher")
	grown := tr.Nodes()
	if grown <= base {
		t.Fatal("inserting a longer word should allocate nodes")
	}
	if !tr.Delete("gopher") {
		t.Error("Delete should report true")
	}
	if tr.Delete("gopher") {
		t.Error("deleting twice should report false")
	}
	if !tr.Contains("go") {
		t.Error("deleting gopher must not remove go")
	}
	if tr.Nodes() != base {
		t.Errorf("Nodes = %d after deleting, want %d - prune the dead branch", tr.Nodes(), base)
	}
	if tr.Delete("nonexistent") {
		t.Error("deleting a missing word should report false")
	}
	// Deleting a prefix word must not remove the longer one.
	tr2 := build("go", "gopher")
	tr2.Delete("go")
	if !tr2.Contains("gopher") || tr2.Contains("go") {
		t.Error("deleting a prefix word broke the longer word")
	}
}

func TestWithPrefix(t *testing.T) {
	tr := build("car", "cart", "cat", "dog", "carbon")
	got := tr.WithPrefix("car", 0)
	want := []string{"car", "carbon", "cart"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("= %v, want %v (lexical order)", got, want)
	}
	if got := tr.WithPrefix("car", 2); !reflect.DeepEqual(got, []string{"car", "carbon"}) {
		t.Errorf("limit = %v", got)
	}
	if got := tr.WithPrefix("", 0); len(got) != 5 {
		t.Errorf("empty prefix returned %d words, want all 5", len(got))
	}
	if got := tr.WithPrefix("zzz", 0); len(got) != 0 {
		t.Errorf("= %v, want none", got)
	}
}

func TestAllIterator(t *testing.T) {
	tr := build("b", "a", "c")
	var got []string
	for w := range tr.All() {
		got = append(got, w)
	}
	if !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("= %v", got)
	}
	n := 0
	for range tr.All() {
		n++
		break
	}
	if n != 1 {
		t.Errorf("iterator visited %d after break", n)
	}
}

func TestLongestPrefixOf(t *testing.T) {
	tr := build("go", "gopher", "golang", "")
	tests := []struct {
		in   string
		want string
	}{
		{"gophers", "gopher"},
		{"gopher", "gopher"},
		{"going", "go"},
		{"go", "go"},
		{"g", ""},
		{"zebra", ""},
	}
	for _, tt := range tests {
		got, ok := tr.LongestPrefixOf(tt.in)
		if !ok || got != tt.want {
			t.Errorf("LongestPrefixOf(%q) = %q, %v; want %q", tt.in, got, ok, tt.want)
		}
	}
	tr2 := build("xyz")
	if _, ok := tr2.LongestPrefixOf("abc"); ok {
		t.Error("no stored prefix should report false")
	}
}

func TestMatch(t *testing.T) {
	tr := build("cat", "cot", "cart", "cut", "dog", "ct")
	tests := []struct {
		pattern string
		want    []string
	}{
		{"c.t", []string{"cat", "cot", "cut"}},
		{"...", []string{"cat", "cot", "cut", "dog"}},
		{"cat", []string{"cat"}},
		{"c..t", []string{"cart"}},
		{"..", []string{"ct"}},
		{"z.", nil},
		{".", nil},
	}
	for _, tt := range tests {
		got := tr.Match(tt.pattern)
		if len(got) == 0 && len(tt.want) == 0 {
			continue
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Match(%q) = %v, want %v", tt.pattern, got, tt.want)
		}
	}
}

func TestCountPrefix(t *testing.T) {
	tr := build("a", "ab", "abc", "b")
	tests := []struct {
		prefix string
		want   int
	}{
		{"", 4}, {"a", 3}, {"ab", 2}, {"abc", 1}, {"abcd", 0}, {"z", 0},
	}
	for _, tt := range tests {
		if got := tr.CountPrefix(tt.prefix); got != tt.want {
			t.Errorf("CountPrefix(%q) = %d, want %d", tt.prefix, got, tt.want)
		}
	}
	tr.Delete("ab")
	if got := tr.CountPrefix("a"); got != 2 {
		t.Errorf("after delete = %d, want 2", got)
	}
}

func TestLargeVocabulary(t *testing.T) {
	tr := New()
	var words []string
	for a := 'a'; a <= 'z'; a++ {
		for b := 'a'; b <= 'z'; b++ {
			for c := 'a'; c <= 'z'; c++ {
				w := string([]rune{a, b, c})
				words = append(words, w)
				tr.Insert(w)
			}
		}
	}
	if tr.Len() != len(words) {
		t.Fatalf("Len = %d, want %d", tr.Len(), len(words))
	}
	sort.Strings(words)
	got := tr.WithPrefix("", 0)
	if !reflect.DeepEqual(got, words) {
		t.Error("full traversal is not in lexical order")
	}
	if got := len(tr.WithPrefix("ab", 0)); got != 26 {
		t.Errorf("WithPrefix(ab) = %d words, want 26", got)
	}
	if got := tr.Match("a.c"); len(got) != 26 {
		t.Errorf("Match = %d words, want 26", len(got))
	}
}

func BenchmarkWithPrefix(b *testing.B) {
	tr := New()
	for i := range 20000 {
		tr.Insert(strings.Repeat("x", i%10) + string(rune('a'+i%26)))
	}
	b.ReportAllocs()
	for b.Loop() {
		tr.WithPrefix("xxx", 10)
	}
}
