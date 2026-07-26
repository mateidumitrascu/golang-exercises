package dedup

import (
	"reflect"
	"testing"
)

func TestDedup(t *testing.T) {
	tests := []struct {
		in   []int
		want []int
	}{
		{[]int{3, 1, 3, 2, 1}, []int{3, 1, 2}},
		{[]int{1, 1, 1}, []int{1}},
		{[]int{1, 2, 3}, []int{1, 2, 3}},
		{[]int{}, []int{}},
		{nil, nil},
	}
	for _, tt := range tests {
		in := append([]int(nil), tt.in...)
		got := Dedup(in)
		if len(got) != len(tt.want) {
			t.Fatalf("Dedup(%v) = %v, want %v", tt.in, got, tt.want)
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Fatalf("Dedup(%v) = %v, want %v", tt.in, got, tt.want)
			}
		}
	}
}

func TestDedupInPlace(t *testing.T) {
	in := []int{5, 5, 7}
	got := Dedup(in)
	if len(got) > 0 && &got[0] != &in[0] {
		t.Error("Dedup must reuse the input's backing array, not allocate a new one")
	}
}

func TestDedupZeroesTail(t *testing.T) {
	type box struct{ n int }
	a, b := &box{1}, &box{1}
	in := []*box{a, b, a}
	got := Dedup(in)
	if len(got) != 2 {
		t.Fatalf("Dedup returned %d elements, want 2", len(got))
	}
	full := in[:cap(in)]
	for i := len(got); i < len(full); i++ {
		if full[i] != nil {
			t.Errorf("element %d past the new length is still %v; zero the tail or you leak", i, full[i])
		}
	}
}

func TestDedupSorted(t *testing.T) {
	tests := []struct{ in, want []string }{
		{[]string{"a", "a", "b", "c", "c", "c"}, []string{"a", "b", "c"}},
		{[]string{"a"}, []string{"a"}},
		{[]string{}, []string{}},
	}
	for _, tt := range tests {
		got := DedupSorted(append([]string(nil), tt.in...))
		if !reflect.DeepEqual(got, tt.want) && !(len(got) == 0 && len(tt.want) == 0) {
			t.Errorf("DedupSorted(%v) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestDeleteIf(t *testing.T) {
	in := []int{1, 2, 3, 4, 5, 6}
	got := DeleteIf(in, func(n int) bool { return n%2 == 0 })
	want := []int{1, 3, 5}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("DeleteIf(evens) = %v, want %v", got, want)
	}
	full := in[:cap(in)]
	for i := len(got); i < len(full); i++ {
		if full[i] != 0 {
			t.Errorf("tail not zeroed at %d: %v", i, full)
		}
	}
}

func TestDeleteIfAll(t *testing.T) {
	got := DeleteIf([]int{1, 2, 3}, func(int) bool { return true })
	if len(got) != 0 {
		t.Errorf("DeleteIf(always) = %v, want empty", got)
	}
}

func TestCompact(t *testing.T) {
	got := Compact([]int{1, 1, 2, 2, 2, 1}, func(a, b int) bool { return a == b })
	want := []int{1, 2, 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Compact = %v, want %v", got, want)
	}

	// Case-insensitive compaction: the FIRST of each run survives.
	words := []string{"Go", "go", "GO", "rust", "Rust"}
	gotS := Compact(words, func(a, b string) bool {
		return len(a) == len(b) && lower(a) == lower(b)
	})
	wantS := []string{"Go", "rust"}
	if !reflect.DeepEqual(gotS, wantS) {
		t.Errorf("Compact(words) = %v, want %v", gotS, wantS)
	}
}

func lower(s string) string {
	b := []byte(s)
	for i := range b {
		if b[i] >= 'A' && b[i] <= 'Z' {
			b[i] += 'a' - 'A'
		}
	}
	return string(b)
}

func BenchmarkDedup(b *testing.B) {
	src := make([]int, 10000)
	for i := range src {
		src[i] = i % 100
	}
	buf := make([]int, len(src))
	b.ReportAllocs()
	for b.Loop() {
		copy(buf, src)
		_ = Dedup(buf)
	}
}
