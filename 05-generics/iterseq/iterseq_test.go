package iterseq

import (
	"reflect"
	"testing"
)

func TestRange(t *testing.T) {
	tests := []struct {
		start, stop, step int
		want              []int
	}{
		{0, 5, 1, []int{0, 1, 2, 3, 4}},
		{0, 10, 3, []int{0, 3, 6, 9}},
		{5, 0, -1, []int{5, 4, 3, 2, 1}},
		{0, 0, 1, nil},
		{0, 5, -1, nil},
		{5, 0, 1, nil},
	}
	for _, tt := range tests {
		got := Collect(Range(tt.start, tt.stop, tt.step))
		if len(got) != len(tt.want) || (len(got) > 0 && !reflect.DeepEqual(got, tt.want)) {
			t.Errorf("Range(%d,%d,%d) = %v, want %v", tt.start, tt.stop, tt.step, got, tt.want)
		}
	}
	defer func() {
		if recover() == nil {
			t.Error("Range with step 0 did not panic")
		}
	}()
	Collect(Range(0, 5, 0))
}

func TestMapFilter(t *testing.T) {
	seq := FromSlice([]int{1, 2, 3, 4, 5})
	got := Collect(Map(Filter(seq, func(n int) bool { return n%2 == 1 }), func(n int) int { return n * 10 }))
	if !reflect.DeepEqual(got, []int{10, 30, 50}) {
		t.Errorf("= %v", got)
	}
	strs := Collect(Map(FromSlice([]int{1, 2}), func(n int) string { return string(rune('a' + n)) }))
	if !reflect.DeepEqual(strs, []string{"b", "c"}) {
		t.Errorf("= %v", strs)
	}
}

// TestLaziness is the important one: an infinite source must be fine as long as
// the consumer stops.
func TestLaziness(t *testing.T) {
	produced := 0
	counter := func(yield func(int) bool) {
		for i := 0; ; i++ {
			produced++
			if !yield(i) {
				return
			}
		}
	}
	got := Collect(Take(Map(counter, func(n int) int { return n * 2 }), 3))
	if !reflect.DeepEqual(got, []int{0, 2, 4}) {
		t.Errorf("= %v", got)
	}
	if produced != 3 {
		t.Errorf("the source produced %d elements, want 3 - Take must stop pulling", produced)
	}
}

func TestTakeDropTakeWhile(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}
	if got := Collect(Take(FromSlice(s), 2)); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Errorf("Take = %v", got)
	}
	if got := Collect(Take(FromSlice(s), 99)); !reflect.DeepEqual(got, s) {
		t.Errorf("Take(99) = %v", got)
	}
	if got := Collect(Take(FromSlice(s), 0)); len(got) != 0 {
		t.Errorf("Take(0) = %v", got)
	}
	if got := Collect(Drop(FromSlice(s), 3)); !reflect.DeepEqual(got, []int{4, 5}) {
		t.Errorf("Drop = %v", got)
	}
	if got := Collect(Drop(FromSlice(s), 99)); len(got) != 0 {
		t.Errorf("Drop(99) = %v", got)
	}
	if got := Collect(TakeWhile(FromSlice(s), func(n int) bool { return n < 4 })); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("TakeWhile = %v", got)
	}
}

func TestBreakPropagates(t *testing.T) {
	produced := 0
	src := func(yield func(int) bool) {
		for i := range 100 {
			produced++
			if !yield(i) {
				return
			}
		}
	}
	n := 0
	for range Filter(Map(src, func(v int) int { return v }), func(int) bool { return true }) {
		n++
		if n == 2 {
			break
		}
	}
	if produced > 2 {
		t.Errorf("source produced %d elements after the consumer broke out at 2; "+
			"every combinator must return false up the chain", produced)
	}
}

func TestChainRepeat(t *testing.T) {
	got := Collect(Chain(FromSlice([]int{1, 2}), FromSlice([]int{3}), FromSlice([]int{})))
	if !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("Chain = %v", got)
	}
	if got := Collect(Take(Repeat("x"), 3)); !reflect.DeepEqual(got, []string{"x", "x", "x"}) {
		t.Errorf("Repeat = %v", got)
	}
	if got := Collect(Chain[int]()); len(got) != 0 {
		t.Errorf("Chain() = %v", got)
	}
}

func TestEnumerate(t *testing.T) {
	var idx []int
	var vals []string
	for i, v := range Enumerate(FromSlice([]string{"a", "b", "c"})) {
		idx = append(idx, i)
		vals = append(vals, v)
	}
	if !reflect.DeepEqual(idx, []int{0, 1, 2}) || !reflect.DeepEqual(vals, []string{"a", "b", "c"}) {
		t.Errorf("Enumerate = %v %v", idx, vals)
	}
	n := 0
	for range Enumerate(Repeat(1)) {
		n++
		if n == 2 {
			break
		}
	}
}

func TestZip(t *testing.T) {
	var pairs []string
	for a, b := range Zip(FromSlice([]int{1, 2, 3}), FromSlice([]string{"x", "y"})) {
		pairs = append(pairs, string(rune('0'+a))+b)
	}
	if !reflect.DeepEqual(pairs, []string{"1x", "2y"}) {
		t.Errorf("Zip = %v, want [1x 2y] (stop at the shorter one)", pairs)
	}
	// Zipping an infinite sequence must terminate.
	n := 0
	for range Zip(Repeat(1), FromSlice([]int{1, 2})) {
		n++
	}
	if n != 2 {
		t.Errorf("Zip yielded %d pairs, want 2", n)
	}
}

func TestMergeSorted(t *testing.T) {
	a := FromSlice([]int{1, 3, 3, 7})
	b := FromSlice([]int{2, 3, 8})
	got := Collect(MergeSorted(a, b))
	want := []int{1, 2, 3, 3, 3, 7, 8}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeSorted = %v, want %v", got, want)
	}
	if got := Collect(MergeSorted(FromSlice([]string{}), FromSlice([]string{"a"}))); !reflect.DeepEqual(got, []string{"a"}) {
		t.Errorf("MergeSorted with an empty side = %v", got)
	}
	// Early exit must not deadlock or leak.
	n := 0
	for range MergeSorted(Repeat(1), Repeat(2)) {
		n++
		if n == 5 {
			break
		}
	}
}

func TestChunk(t *testing.T) {
	got := Collect(Chunk(Range(0, 7, 1), 3))
	want := [][]int{{0, 1, 2}, {3, 4, 5}, {6}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Chunk = %v, want %v", got, want)
	}
	// Chunks must not share a buffer.
	chunks := Collect(Chunk(Range(0, 4, 1), 2))
	chunks[0][0] = 99
	if chunks[1][0] == 99 {
		t.Error("chunks share memory; allocate a new slice per chunk")
	}
	if got := Collect(Chunk(Range(0, 0, 1), 3)); len(got) != 0 {
		t.Errorf("Chunk of an empty sequence = %v", got)
	}
}

func TestConsumers(t *testing.T) {
	if got := Reduce(Range(1, 5, 1), 0, func(a, v int) int { return a + v }); got != 10 {
		t.Errorf("Reduce = %d, want 10", got)
	}
	if got := Count(Range(0, 5, 1)); got != 5 {
		t.Errorf("Count = %d", got)
	}
	v, ok := First(Range(3, 9, 1))
	if !ok || v != 3 {
		t.Errorf("First = %v, %v", v, ok)
	}
	if _, ok := First(Range(0, 0, 1)); ok {
		t.Error("First of an empty sequence must report false")
	}
	// First must not consume the whole (infinite) sequence.
	if v, ok := First(Repeat(7)); !ok || v != 7 {
		t.Errorf("First(Repeat) = %v, %v", v, ok)
	}
}
