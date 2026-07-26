package heapq

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func lessInt(a, b int) bool { return a < b }

func checkInvariant(t *testing.T, h *Heap[int]) {
	t.Helper()
	s := h.Slice()
	for i := range s {
		for _, c := range []int{2*i + 1, 2*i + 2} {
			if c < len(s) && s[c] < s[i] {
				t.Fatalf("heap invariant broken at %d: %v", i, s)
			}
		}
	}
}

func TestPushPopOrder(t *testing.T) {
	h := New(lessInt)
	if _, ok := h.Pop(); ok {
		t.Error("Pop on an empty heap returned ok")
	}
	if _, ok := h.Peek(); ok {
		t.Error("Peek on an empty heap returned ok")
	}
	input := []int{5, 3, 8, 1, 9, 2, 7}
	for _, v := range input {
		h.Push(v)
		checkInvariant(t, h)
	}
	if h.Len() != len(input) {
		t.Errorf("Len = %d", h.Len())
	}
	if v, _ := h.Peek(); v != 1 {
		t.Errorf("Peek = %d, want 1", v)
	}
	var got []int
	for h.Len() > 0 {
		v, _ := h.Pop()
		got = append(got, v)
		checkInvariant(t, h)
	}
	want := []int{1, 2, 3, 5, 7, 8, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("= %v, want %v", got, want)
	}
}

func TestMaxHeap(t *testing.T) {
	h := New(func(a, b int) bool { return a > b })
	for _, v := range []int{1, 5, 3} {
		h.Push(v)
	}
	if v, _ := h.Pop(); v != 5 {
		t.Errorf("= %d, want 5", v)
	}
}

func TestFromIsLinearAndCorrect(t *testing.T) {
	s := rand.Perm(1000)
	h := From(s, lessInt)
	checkInvariant(t, h)
	prev := -1
	for h.Len() > 0 {
		v, _ := h.Pop()
		if v < prev {
			t.Fatalf("out of order: %d after %d", v, prev)
		}
		prev = v
	}
}

func TestPushPopCombined(t *testing.T) {
	h := From([]int{2, 4, 6}, lessInt)
	v, ok := h.PushPop(5)
	if !ok || v != 2 {
		t.Errorf("PushPop = %d, %v; want 2", v, ok)
	}
	checkInvariant(t, h)
	if h.Len() != 3 {
		t.Errorf("Len = %d, want 3", h.Len())
	}
	// Pushing something smaller than the head returns it straight back.
	if v, _ := h.PushPop(0); v != 0 {
		t.Errorf("PushPop(0) = %d, want 0", v)
	}
	empty := New(lessInt)
	if v, ok := empty.PushPop(7); !ok || v != 7 {
		t.Errorf("PushPop on an empty heap = %d, %v", v, ok)
	}
}

func TestPQ(t *testing.T) {
	q := NewPQ[string]()
	if _, ok := q.Pop(); ok {
		t.Error("Pop on an empty queue")
	}
	q.Push("low", 5)
	urgent := q.Push("urgent", 1)
	q.Push("mid", 3)

	if q.Len() != 3 {
		t.Errorf("Len = %d", q.Len())
	}
	item, _ := q.Pop()
	if item.Value != "urgent" {
		t.Errorf("= %q, want urgent", item.Value)
	}
	_ = urgent

	// Reprioritise: "low" jumps the queue.
	low := q.Push("low2", 9)
	q.Update(low, 0)
	item, _ = q.Pop()
	if item.Value != "low2" {
		t.Errorf("after Update = %q, want low2", item.Value)
	}
	item, _ = q.Pop()
	if item.Value != "mid" {
		t.Errorf("= %q, want mid", item.Value)
	}
}

func TestPQUpdateAfterPop(t *testing.T) {
	q := NewPQ[int]()
	it := q.Push(1, 1)
	q.Pop()
	q.Update(it, 99) // must not panic or corrupt anything
	if q.Len() != 0 {
		t.Errorf("Len = %d", q.Len())
	}
}

func TestPQManyUpdates(t *testing.T) {
	q := NewPQ[int]()
	items := make([]*Item[int], 100)
	for i := range items {
		items[i] = q.Push(i, 1000-i)
	}
	for i := range items {
		q.Update(items[i], i)
	}
	prev := -1
	for q.Len() > 0 {
		it, _ := q.Pop()
		if it.Priority < prev {
			t.Fatalf("out of order: %d after %d", it.Priority, prev)
		}
		prev = it.Priority
	}
}

func TestTopK(t *testing.T) {
	s := []int{5, 1, 9, 3, 7, 2, 8}
	got := TopK(s, 3, lessInt)
	if !reflect.DeepEqual(got, []int{9, 8, 7}) {
		t.Errorf("TopK = %v, want [9 8 7]", got)
	}
	if got := TopK(s, 100, lessInt); len(got) != len(s) || got[0] != 9 {
		t.Errorf("k > len = %v", got)
	}
	if got := TopK(s, 0, lessInt); len(got) != 0 {
		t.Errorf("k = 0 gave %v", got)
	}
	if got := TopK([]int(nil), 3, lessInt); len(got) != 0 {
		t.Errorf("empty input = %v", got)
	}
	// Strings, by length.
	words := []string{"aaa", "b", "cccc", "dd"}
	if got := TopK(words, 2, func(a, b string) bool { return len(a) < len(b) }); !reflect.DeepEqual(got, []string{"cccc", "aaa"}) {
		t.Errorf("= %v", got)
	}
	if !reflect.DeepEqual(s, []int{5, 1, 9, 3, 7, 2, 8}) {
		t.Error("TopK must not modify its input")
	}
}

func TestTopKMemory(t *testing.T) {
	s := rand.Perm(200000)
	var got []int
	allocs := testing.AllocsPerRun(3, func() {
		got = TopK(s, 10, lessInt)
	})
	if len(got) != 10 || got[0] != 199999 {
		t.Fatalf("= %v", got[:min(3, len(got))])
	}
	if allocs > 20 {
		t.Errorf("%.0f allocations; keep a heap of k, do not copy or sort the input", allocs)
	}
}

func TestHeapSort(t *testing.T) {
	for _, n := range []int{0, 1, 2, 7, 100, 1001} {
		s := rand.Perm(n)
		HeapSort(s, lessInt)
		if !sort.IntsAreSorted(s) {
			t.Fatalf("n=%d not sorted", n)
		}
	}
	s := []int{3, 1, 2}
	if allocs := testing.AllocsPerRun(10, func() { HeapSort(s, lessInt) }); allocs > 0 {
		t.Errorf("HeapSort allocated %.0f times; sort in place", allocs)
	}
	desc := []int{1, 2, 3}
	HeapSort(desc, func(a, b int) bool { return a > b })
	if !reflect.DeepEqual(desc, []int{3, 2, 1}) {
		t.Errorf("= %v", desc)
	}
}

func TestMergeSorted(t *testing.T) {
	lists := [][]int{
		{1, 4, 7},
		{2, 5, 8},
		{3, 6, 9},
		{},
		{0},
	}
	got := MergeSorted(lists, lessInt)
	want := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("= %v, want %v", got, want)
	}
	if got := MergeSorted(nil, lessInt); len(got) != 0 {
		t.Errorf("= %v", got)
	}
	if got := MergeSorted([][]int{{1, 1, 2}, {1, 3}}, lessInt); !reflect.DeepEqual(got, []int{1, 1, 1, 2, 3}) {
		t.Errorf("duplicates = %v", got)
	}
}

func BenchmarkTopK(b *testing.B) {
	s := rand.Perm(100000)
	b.ReportAllocs()
	for b.Loop() {
		TopK(s, 20, lessInt)
	}
}
