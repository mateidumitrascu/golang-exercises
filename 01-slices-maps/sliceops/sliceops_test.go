package sliceops

import (
	"reflect"
	"testing"
)

func TestInsertAt(t *testing.T) {
	tests := []struct {
		s    []int
		i    int
		v    []int
		want []int
	}{
		{[]int{1, 2, 5}, 2, []int{3, 4}, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3}, 0, []int{0}, []int{0, 1, 2, 3}},
		{[]int{1, 2, 3}, 3, []int{4}, []int{1, 2, 3, 4}},
		{[]int{}, 0, []int{1}, []int{1}},
		{nil, 0, []int{1, 2}, []int{1, 2}},
		{[]int{1, 2}, 1, nil, []int{1, 2}},
	}
	for _, tt := range tests {
		got := InsertAt(append([]int(nil), tt.s...), tt.i, tt.v...)
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("InsertAt(%v, %d, %v) = %v, want %v", tt.s, tt.i, tt.v, got, tt.want)
		}
	}
}

func TestInsertAtPanics(t *testing.T) {
	for _, i := range []int{-1, 4} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("InsertAt with i=%d did not panic", i)
				}
			}()
			InsertAt([]int{1, 2, 3}, i, 9)
		}()
	}
}

func TestDeleteRange(t *testing.T) {
	tests := []struct {
		s    []int
		i, j int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 1, 3, []int{1, 4, 5}},
		{[]int{1, 2, 3}, 0, 3, []int{}},
		{[]int{1, 2, 3}, 1, 1, []int{1, 2, 3}},
		{[]int{1, 2, 3}, 2, 3, []int{1, 2}},
	}
	for _, tt := range tests {
		in := append([]int(nil), tt.s...)
		got := DeleteRange(in, tt.i, tt.j)
		if len(got) != len(tt.want) {
			t.Fatalf("DeleteRange(%v, %d, %d) = %v, want %v", tt.s, tt.i, tt.j, got, tt.want)
		}
		for k := range got {
			if got[k] != tt.want[k] {
				t.Fatalf("DeleteRange(%v, %d, %d) = %v, want %v", tt.s, tt.i, tt.j, got, tt.want)
			}
		}
		full := in[:cap(in)]
		for k := len(got); k < len(full); k++ {
			if full[k] != 0 {
				t.Errorf("DeleteRange left junk at index %d: %v", k, full)
			}
		}
	}
}

func TestDeleteUnordered(t *testing.T) {
	in := []int{10, 20, 30, 40}
	got := DeleteUnordered(in, 1)
	want := []int{10, 40, 30}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DeleteUnordered(%v, 1) = %v, want %v (swap with the last element)", []int{10, 20, 30, 40}, got, want)
	}
	if in[:cap(in)][3] != 0 {
		t.Error("DeleteUnordered must zero the vacated slot")
	}

	single := DeleteUnordered([]int{7}, 0)
	if len(single) != 0 {
		t.Errorf("deleting the only element gave %v", single)
	}
}

func TestClone(t *testing.T) {
	if Clone[int](nil) != nil {
		t.Error("Clone(nil) must be nil")
	}
	empty := Clone([]int{})
	if empty == nil || len(empty) != 0 {
		t.Error("Clone of an empty slice must be non-nil and empty")
	}
	src := []int{1, 2, 3}
	dst := Clone(src)
	dst[0] = 99
	if src[0] != 1 {
		t.Error("Clone must not share memory with the source")
	}
}

func TestGrow(t *testing.T) {
	s := make([]int, 3, 10)
	got := Grow(s, 5)
	if cap(got) != 10 || &got[0] != &s[0] {
		t.Error("Grow must not reallocate when there is already room")
	}
	got = Grow(s, 100)
	if cap(got) < 103 {
		t.Errorf("Grow(s, 100) gave capacity %d, want at least %d", cap(got), 103)
	}
	if len(got) != 3 || got[0] != s[0] {
		t.Error("Grow must preserve length and contents")
	}
}

func TestEqual(t *testing.T) {
	tests := []struct {
		a, b []string
		want bool
	}{
		{[]string{"a", "b"}, []string{"a", "b"}, true},
		{[]string{"a"}, []string{"b"}, false},
		{[]string{"a"}, []string{"a", "b"}, false},
		{nil, []string{}, true},
		{nil, nil, true},
	}
	for _, tt := range tests {
		if got := Equal(tt.a, tt.b); got != tt.want {
			t.Errorf("Equal(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
	}
}

func TestReslice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	r := Reslice(s, 1, 3)
	if !reflect.DeepEqual(r, []int{2, 3}) {
		t.Fatalf("Reslice(s, 1, 3) = %v, want [2 3]", r)
	}
	if cap(r) != 2 {
		t.Errorf("cap = %d, want 2", cap(r))
	}
	_ = append(r, 99)
	if s[3] != 4 {
		t.Errorf("appending to the reslice leaked into s: %v", s)
	}
}
