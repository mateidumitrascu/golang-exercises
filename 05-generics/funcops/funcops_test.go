package funcops

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestMap(t *testing.T) {
	got := Map([]int{1, 2, 3}, func(n int) string { return strconv.Itoa(n * 2) })
	if !reflect.DeepEqual(got, []string{"2", "4", "6"}) {
		t.Errorf("Map = %v", got)
	}
	if got := Map([]int{}, strconv.Itoa); got == nil || len(got) != 0 {
		t.Errorf("Map of an empty slice = %v, want empty non-nil", got)
	}
	if got := MapIndex([]string{"a", "b"}, func(i int, s string) string {
		return strconv.Itoa(i) + s
	}); !reflect.DeepEqual(got, []string{"0a", "1b"}) {
		t.Errorf("MapIndex = %v", got)
	}
}

func TestFilter(t *testing.T) {
	got := Filter([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 1 })
	if !reflect.DeepEqual(got, []int{1, 3, 5}) {
		t.Errorf("Filter = %v", got)
	}
	if got := Filter([]int{1, 2}, func(int) bool { return false }); len(got) != 0 {
		t.Errorf("Filter = %v, want empty", got)
	}
}

func TestReduce(t *testing.T) {
	if got := Reduce([]int{1, 2, 3}, 0, func(a, v int) int { return a + v }); got != 6 {
		t.Errorf("Reduce = %d", got)
	}
	// Accumulator of a different type.
	got := Reduce([]int{1, 2, 3}, "", func(a string, v int) string { return a + strconv.Itoa(v) })
	if got != "123" {
		t.Errorf("Reduce = %q, want 123", got)
	}
	// Order matters: left fold vs right fold.
	sub := func(a, v int) int { return a - v }
	if l, r := Reduce([]int{1, 2, 3}, 0, sub), ReduceRight([]int{1, 2, 3}, 0, sub); l != -6 || r != -6 {
		t.Errorf("folds = %d, %d", l, r)
	}
	rs := ReduceRight([]string{"a", "b", "c"}, "", func(acc, v string) string { return acc + v })
	if rs != "cba" {
		t.Errorf("ReduceRight = %q, want cba", rs)
	}
}

func TestFlatMap(t *testing.T) {
	got := FlatMap([]string{"a b", "c"}, func(s string) []string { return strings.Fields(s) })
	if !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("FlatMap = %v", got)
	}
}

func TestPredicates(t *testing.T) {
	even := func(n int) bool { return n%2 == 0 }
	nums := []int{2, 4, 6}
	if !All(nums, even) || !Any(nums, even) || None(nums, even) {
		t.Error("predicates on all-even")
	}
	if All([]int{1, 2}, even) || !Any([]int{1, 2}, even) {
		t.Error("predicates on mixed")
	}
	if !All([]int{}, even) || Any([]int{}, even) || !None([]int{}, even) {
		t.Error("empty slice: All is true, Any is false, None is true")
	}

	// Short-circuiting: Any must stop at the first match.
	calls := 0
	Any([]int{2, 4, 6}, func(n int) bool { calls++; return true })
	if calls != 1 {
		t.Errorf("Any made %d calls, want 1", calls)
	}
	calls = 0
	All([]int{1, 3, 5}, func(n int) bool { calls++; return false })
	if calls != 1 {
		t.Errorf("All made %d calls, want 1", calls)
	}
}

func TestFindAndIndex(t *testing.T) {
	v, ok := Find([]string{"a", "bb", "ccc"}, func(s string) bool { return len(s) == 2 })
	if !ok || v != "bb" {
		t.Errorf("Find = %q, %v", v, ok)
	}
	if v, ok := Find([]string{"a"}, func(string) bool { return false }); ok || v != "" {
		t.Errorf("Find = %q, %v; want the zero value and false", v, ok)
	}
	if got := IndexFunc([]int{1, 2, 3}, func(n int) bool { return n == 3 }); got != 2 {
		t.Errorf("IndexFunc = %d", got)
	}
	if got := IndexFunc([]int{1}, func(int) bool { return false }); got != -1 {
		t.Errorf("IndexFunc = %d, want -1", got)
	}
}

func TestUniq(t *testing.T) {
	type user struct{ id, age int }
	users := []user{{1, 20}, {2, 30}, {1, 99}}
	got := Uniq(users, func(u user) int { return u.id })
	want := []user{{1, 20}, {2, 30}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Uniq = %v, want %v (keep the first of each key)", got, want)
	}
	if got := Uniq([]string{"a", "A", "b"}, strings.ToLower); len(got) != 2 {
		t.Errorf("Uniq = %v, want 2 elements", got)
	}
}

func TestZipUnzip(t *testing.T) {
	ps := Zip([]int{1, 2, 3}, []string{"a", "b"})
	want := []Pair[int, string]{{1, "a"}, {2, "b"}}
	if !reflect.DeepEqual(ps, want) {
		t.Errorf("Zip = %v, want %v (stop at the shorter slice)", ps, want)
	}
	as, bs := Unzip(want)
	if !reflect.DeepEqual(as, []int{1, 2}) || !reflect.DeepEqual(bs, []string{"a", "b"}) {
		t.Errorf("Unzip = %v, %v", as, bs)
	}
	if got := Zip([]int{}, []string{"a"}); len(got) != 0 {
		t.Errorf("Zip = %v, want empty", got)
	}
}

func TestCompose(t *testing.T) {
	double := func(n int) int { return n * 2 }
	str := func(n int) string { return strconv.Itoa(n) }
	f := Compose(double, str)
	if got := f(21); got != "42" {
		t.Errorf("Compose = %q, want 42", got)
	}
	g := Compose(str, func(s string) int { return len(s) })
	if got := g(1000); got != 4 {
		t.Errorf("Compose = %d, want 4", got)
	}
}

func TestMemoize(t *testing.T) {
	calls := 0
	slow := func(n int) int {
		calls++
		return n * n
	}
	fast := Memoize(slow)
	if fast(4) != 16 || fast(4) != 16 || fast(5) != 25 {
		t.Error("memoized function returned wrong values")
	}
	if calls != 2 {
		t.Errorf("underlying function called %d times, want 2", calls)
	}
	// The zero value of V must be cached too, not recomputed.
	calls = 0
	zero := Memoize(func(string) int { calls++; return 0 })
	zero("x")
	zero("x")
	if calls != 1 {
		t.Errorf("a cached zero value was recomputed (%d calls); use the comma-ok form", calls)
	}
}

func TestPartial(t *testing.T) {
	add := func(a, b int) int { return a + b }
	add5 := Partial(add, 5)
	if got := add5(3); got != 8 {
		t.Errorf("Partial = %d", got)
	}
	prefix := Partial(func(p, s string) string { return p + s }, ">> ")
	if got := prefix("hi"); got != ">> hi" {
		t.Errorf("Partial = %q", got)
	}
}
