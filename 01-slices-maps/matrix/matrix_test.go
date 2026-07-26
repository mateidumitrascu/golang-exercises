package matrix

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	m := New(3, 4)
	if len(m) != 3 {
		t.Fatalf("rows = %d, want 3", len(m))
	}
	for _, r := range m {
		if len(r) != 4 {
			t.Fatalf("cols = %d, want 4", len(r))
		}
	}
	m[1][2] = 7
	if m[0][2] == 7 {
		t.Fatal("rows must not alias each other")
	}
	// One flat array plus one array of row headers: two allocations, not rows+1.
	got := testing.AllocsPerRun(50, func() { _ = New(64, 64) })
	if got > 2 {
		t.Errorf("New allocated %.0f times; use one flat array plus one slice header array", got)
	}
	if m2 := New(0, 5); len(m2) != 0 {
		t.Errorf("New(0,5) = %v, want empty", m2)
	}
}

func TestTranspose(t *testing.T) {
	m := [][]int{{1, 2, 3}, {4, 5, 6}}
	want := [][]int{{1, 4}, {2, 5}, {3, 6}}
	if got := Transpose(m); !reflect.DeepEqual(got, want) {
		t.Errorf("Transpose = %v, want %v", got, want)
	}
	if !reflect.DeepEqual(m, [][]int{{1, 2, 3}, {4, 5, 6}}) {
		t.Error("Transpose must not modify its input")
	}
	if got := Transpose(nil); len(got) != 0 {
		t.Errorf("Transpose(nil) = %v, want empty", got)
	}
}

func TestTransposePanicsOnRagged(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("Transpose of a ragged matrix did not panic")
		}
	}()
	Transpose([][]int{{1, 2}, {3}})
}

func TestRotateCW(t *testing.T) {
	m := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	RotateCW(m)
	want := [][]int{
		{7, 4, 1},
		{8, 5, 2},
		{9, 6, 3},
	}
	if !reflect.DeepEqual(m, want) {
		t.Fatalf("RotateCW =\n%v\nwant\n%v", m, want)
	}

	four := [][]int{
		{1, 2, 3, 4},
		{5, 6, 7, 8},
		{9, 10, 11, 12},
		{13, 14, 15, 16},
	}
	RotateCW(four)
	want4 := [][]int{
		{13, 9, 5, 1},
		{14, 10, 6, 2},
		{15, 11, 7, 3},
		{16, 12, 8, 4},
	}
	if !reflect.DeepEqual(four, want4) {
		t.Fatalf("RotateCW 4x4 =\n%v\nwant\n%v", four, want4)
	}

	// Four rotations return to the original.
	orig := [][]int{{1, 2}, {3, 4}}
	cp := [][]int{{1, 2}, {3, 4}}
	for range 4 {
		RotateCW(cp)
	}
	if !reflect.DeepEqual(cp, orig) {
		t.Errorf("four rotations gave %v, want %v", cp, orig)
	}
}

func TestRotateCWNoAlloc(t *testing.T) {
	m := New(32, 32)
	if got := testing.AllocsPerRun(20, func() { RotateCW(m) }); got > 0 {
		t.Errorf("RotateCW allocated %.0f times; rotate in place", got)
	}
}

func TestSpiral(t *testing.T) {
	tests := []struct {
		m    [][]int
		want []int
	}{
		{[][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, []int{1, 2, 3, 6, 9, 8, 7, 4, 5}},
		{[][]int{{1, 2, 3, 4}}, []int{1, 2, 3, 4}},
		{[][]int{{1}, {2}, {3}}, []int{1, 2, 3}},
		{[][]int{{1, 2}, {3, 4}}, []int{1, 2, 4, 3}},
		{[][]int{
			{1, 2, 3, 4},
			{5, 6, 7, 8},
			{9, 10, 11, 12},
		}, []int{1, 2, 3, 4, 8, 12, 11, 10, 9, 5, 6, 7}},
		{nil, []int{}},
	}
	for _, tt := range tests {
		got := Spiral(tt.m)
		if len(got) != len(tt.want) {
			t.Fatalf("Spiral(%v) = %v, want %v", tt.m, got, tt.want)
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Fatalf("Spiral(%v) = %v, want %v", tt.m, got, tt.want)
			}
		}
	}
}

func TestNeighbours(t *testing.T) {
	m := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	tests := []struct {
		r, c int
		want []int
	}{
		{1, 1, []int{1, 2, 3, 4, 6, 7, 8, 9}},
		{0, 0, []int{2, 4, 5}},
		{2, 2, []int{5, 6, 8}},
		{0, 1, []int{1, 3, 4, 5, 6}},
	}
	for _, tt := range tests {
		if got := Neighbours(m, tt.r, tt.c); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("Neighbours(%d,%d) = %v, want %v", tt.r, tt.c, got, tt.want)
		}
	}
}

func TestFlatUnflat(t *testing.T) {
	m := [][]int{{1, 2, 3}, {4, 5, 6}}
	flat, stride := Flat(m)
	if !reflect.DeepEqual(flat, []int{1, 2, 3, 4, 5, 6}) || stride != 3 {
		t.Fatalf("Flat = %v, %d", flat, stride)
	}
	if got := Unflat(flat, stride); !reflect.DeepEqual(got, m) {
		t.Errorf("Unflat = %v, want %v", got, m)
	}
	defer func() {
		if recover() == nil {
			t.Error("Unflat with a bad stride did not panic")
		}
	}()
	Unflat([]int{1, 2, 3}, 2)
}
