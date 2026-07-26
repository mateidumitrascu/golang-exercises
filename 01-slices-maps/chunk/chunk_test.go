package chunk

import (
	"reflect"
	"testing"
)

func lens[T any](gs [][]T) []int {
	out := make([]int, len(gs))
	for i, g := range gs {
		out[i] = len(g)
	}
	return out
}

func TestChunk(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		size int
		want [][]int
	}{
		{"even", []int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {3, 4}}},
		{"ragged", []int{1, 2, 3, 4, 5}, 2, [][]int{{1, 2}, {3, 4}, {5}}},
		{"size bigger than input", []int{1, 2}, 10, [][]int{{1, 2}}},
		{"size one", []int{1, 2}, 1, [][]int{{1}, {2}}},
		{"empty", []int{}, 3, [][]int{}},
		{"nil", nil, 3, [][]int{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Chunk(tt.in, tt.size)
			if len(got) != len(tt.want) {
				t.Fatalf("Chunk(%v, %d) = %v, want %v", tt.in, tt.size, got, tt.want)
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("group %d = %v, want %v", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestChunkPanicsOnBadSize(t *testing.T) {
	for _, size := range []int{0, -1} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("Chunk(s, %d) did not panic", size)
				}
			}()
			Chunk([]int{1, 2, 3}, size)
		}()
	}
}

func TestChunkSharesMemory(t *testing.T) {
	s := []int{1, 2, 3, 4}
	gs := Chunk(s, 2)
	gs[1][0] = 99
	if s[2] != 99 {
		t.Fatalf("chunks must alias the input; s = %v, want s[2] == 99", s)
	}
}

func TestChunkAppendDoesNotClobber(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6}
	gs := Chunk(s, 2)
	_ = append(gs[0], 999) // must reallocate, not scribble on s[2]
	if s[2] != 3 {
		t.Fatalf("appending to a chunk overwrote the next chunk: s = %v\n"+
			"hint: full slice expressions, s[i:j:j]", s)
	}
}

func TestWindows(t *testing.T) {
	tests := []struct {
		in   []int
		size int
		want [][]int
	}{
		{[]int{1, 2, 3, 4}, 2, [][]int{{1, 2}, {2, 3}, {3, 4}}},
		{[]int{1, 2, 3}, 3, [][]int{{1, 2, 3}}},
		{[]int{1, 2}, 3, [][]int{}},
		{nil, 1, [][]int{}},
	}
	for _, tt := range tests {
		got := Windows(tt.in, tt.size)
		if len(got) != len(tt.want) {
			t.Fatalf("Windows(%v, %d) = %v, want %v", tt.in, tt.size, got, tt.want)
		}
		for i := range got {
			if !reflect.DeepEqual(got[i], tt.want[i]) {
				t.Errorf("Windows(%v, %d)[%d] = %v, want %v", tt.in, tt.size, i, got[i], tt.want[i])
			}
		}
	}
}

func TestWindowsAppendDoesNotClobber(t *testing.T) {
	s := []int{1, 2, 3, 4}
	w := Windows(s, 2)
	_ = append(w[0], 999)
	if s[2] != 3 {
		t.Fatalf("appending to a window overwrote the input: %v", s)
	}
}

func TestSplit(t *testing.T) {
	tests := []struct {
		in       []int
		n        int
		wantLens []int
		wantFlat []int
	}{
		{[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 3, []int{4, 3, 3}, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{[]int{1, 2, 3, 4}, 4, []int{1, 1, 1, 1}, []int{1, 2, 3, 4}},
		{[]int{1, 2}, 5, []int{1, 1, 0, 0, 0}, []int{1, 2}},
		{[]int{1, 2, 3, 4, 5}, 2, []int{3, 2}, []int{1, 2, 3, 4, 5}},
		{nil, 3, []int{0, 0, 0}, []int{}},
	}
	for _, tt := range tests {
		got := Split(tt.in, tt.n)
		if len(got) != tt.n {
			t.Fatalf("Split(%v, %d) returned %d groups, want %d", tt.in, tt.n, len(got), tt.n)
		}
		if !reflect.DeepEqual(lens(got), tt.wantLens) {
			t.Errorf("Split(%v, %d) group lengths = %v, want %v", tt.in, tt.n, lens(got), tt.wantLens)
		}
		flat := []int{}
		for _, g := range got {
			flat = append(flat, g...)
		}
		if !reflect.DeepEqual(flat, tt.wantFlat) {
			t.Errorf("Split(%v, %d) flattened = %v, want %v", tt.in, tt.n, flat, tt.wantFlat)
		}
	}
}

func TestSplitPanicsOnBadN(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Split(s, 0) did not panic")
		}
	}()
	Split([]int{1}, 0)
}

func TestGenericOverStrings(t *testing.T) {
	got := Chunk([]string{"a", "b", "c"}, 2)
	want := [][]string{{"a", "b"}, {"c"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func BenchmarkChunk(b *testing.B) {
	s := make([]int, 1<<16)
	b.ReportAllocs()
	for b.Loop() {
		_ = Chunk(s, 64)
	}
}
