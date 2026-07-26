package rotate

import (
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	tests := []struct{ in, want []int }{
		{[]int{1, 2, 3}, []int{3, 2, 1}},
		{[]int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
		{[]int{1}, []int{1}},
		{[]int{}, []int{}},
	}
	for _, tt := range tests {
		in := append([]int(nil), tt.in...)
		Reverse(in)
		if !reflect.DeepEqual(in, tt.want) && len(in) > 0 {
			t.Errorf("Reverse(%v) -> %v, want %v", tt.in, in, tt.want)
		}
	}
	Reverse([]int(nil)) // must not panic
}

func TestRotateLeft(t *testing.T) {
	tests := []struct {
		in   []int
		k    int
		want []int
	}{
		{[]int{1, 2, 3, 4, 5}, 2, []int{3, 4, 5, 1, 2}},
		{[]int{1, 2, 3, 4, 5}, 0, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 5, []int{1, 2, 3, 4, 5}},
		{[]int{1, 2, 3, 4, 5}, 7, []int{3, 4, 5, 1, 2}},
		{[]int{1, 2, 3, 4, 5}, -1, []int{5, 1, 2, 3, 4}},
		{[]int{1, 2, 3, 4, 5}, -7, []int{4, 5, 1, 2, 3}},
		{[]int{1, 2, 3, 4, 5, 6}, 3, []int{4, 5, 6, 1, 2, 3}},
		{[]int{1, 2, 3, 4, 5, 6}, 2, []int{3, 4, 5, 6, 1, 2}},
		{[]int{1}, 3, []int{1}},
		{[]int{}, 3, []int{}},
	}
	for _, tt := range tests {
		in := append([]int(nil), tt.in...)
		RotateLeft(in, tt.k)
		if len(in) > 0 && !reflect.DeepEqual(in, tt.want) {
			t.Errorf("RotateLeft(%v, %d) -> %v, want %v", tt.in, tt.k, in, tt.want)
		}
	}
}

func TestRotateRight(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	RotateRight(in, 2)
	want := []int{4, 5, 1, 2, 3}
	if !reflect.DeepEqual(in, want) {
		t.Errorf("RotateRight -> %v, want %v", in, want)
	}
}

func TestRotateNoAllocations(t *testing.T) {
	s := make([]int, 1024)
	got := testing.AllocsPerRun(100, func() { RotateLeft(s, 37) })
	if got > 0 {
		t.Errorf("RotateLeft allocated %.0f times per call; do it in place", got)
	}
}

func TestIsRotation(t *testing.T) {
	tests := []struct {
		a, b []int
		want bool
	}{
		{[]int{1, 2, 3, 4}, []int{3, 4, 1, 2}, true},
		{[]int{1, 2, 3, 4}, []int{1, 2, 3, 4}, true},
		{[]int{1, 2, 3, 4}, []int{4, 1, 2, 3}, true},
		{[]int{1, 2, 3, 4}, []int{1, 2, 4, 3}, false},
		{[]int{1, 2, 3}, []int{1, 2, 3, 4}, false},
		{[]int{}, []int{}, true},
		{[]int{1, 1, 2}, []int{1, 2, 1}, true},
		{[]int{1, 1, 2}, []int{2, 1, 1}, true},
		{[]int{1, 2, 1}, []int{1, 1, 1}, false},
	}
	for _, tt := range tests {
		a := append([]int(nil), tt.a...)
		b := append([]int(nil), tt.b...)
		if got := IsRotation(a, b); got != tt.want {
			t.Errorf("IsRotation(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
		}
		if !reflect.DeepEqual(a, tt.a) || !reflect.DeepEqual(b, tt.b) {
			t.Fatalf("IsRotation modified its inputs: %v %v", a, b)
		}
	}
}

func BenchmarkRotateLeft(b *testing.B) {
	s := make([]int, 1<<16)
	b.ReportAllocs()
	for b.Loop() {
		RotateLeft(s, 1234)
	}
}
