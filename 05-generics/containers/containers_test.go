package containers

import (
	"reflect"
	"sort"
	"testing"
)

func sorted(s *Set[int]) []int {
	var out []int
	for v := range s.All() {
		out = append(out, v)
	}
	sort.Ints(out)
	return out
}

func TestZeroValueSet(t *testing.T) {
	var s Set[int]
	if s.Len() != 0 || s.Has(1) {
		t.Error("the zero Set must read as empty")
	}
	for range s.All() {
		t.Fatal("the zero Set must iterate nothing")
	}
	s.Add(1)
	if !s.Has(1) {
		t.Error("Add on a zero Set must work")
	}
}

func TestSetBasics(t *testing.T) {
	s := NewSet(1, 2, 2, 3)
	if s.Len() != 3 {
		t.Errorf("Len = %d, want 3", s.Len())
	}
	if !s.Has(2) || s.Has(99) {
		t.Error("Has is wrong")
	}
	s.Delete(2, 99)
	if s.Has(2) || s.Len() != 2 {
		t.Errorf("after Delete: %v", sorted(s))
	}
	s.Add(4).Add(5)
	if got := sorted(s); !reflect.DeepEqual(got, []int{1, 3, 4, 5}) {
		t.Errorf("= %v", got)
	}
	c := s.Clone()
	c.Add(99)
	if s.Has(99) {
		t.Error("Clone must be independent")
	}
	s.Clear()
	if s.Len() != 0 {
		t.Error("Clear did not empty the set")
	}
}

func TestSetIterationStopsEarly(t *testing.T) {
	s := NewSet(1, 2, 3, 4, 5)
	n := 0
	for range s.All() {
		n++
		break
	}
	if n != 1 {
		t.Errorf("visited %d elements after break, want 1", n)
	}
}

func TestSetAlgebra(t *testing.T) {
	a := NewSet(1, 2, 3)
	b := NewSet(3, 4)

	if got := sorted(a.Union(b)); !reflect.DeepEqual(got, []int{1, 2, 3, 4}) {
		t.Errorf("Union = %v", got)
	}
	if got := sorted(a.Intersect(b)); !reflect.DeepEqual(got, []int{3}) {
		t.Errorf("Intersect = %v", got)
	}
	if got := sorted(a.Difference(b)); !reflect.DeepEqual(got, []int{1, 2}) {
		t.Errorf("Difference = %v", got)
	}
	if got := sorted(a.SymmetricDiff(b)); !reflect.DeepEqual(got, []int{1, 2, 4}) {
		t.Errorf("SymmetricDiff = %v", got)
	}
	if a.Len() != 3 || b.Len() != 2 {
		t.Error("set operations must not modify their operands")
	}

	if !NewSet(1, 2).IsSubsetOf(a) || a.IsSubsetOf(NewSet(1, 2)) {
		t.Error("IsSubsetOf is wrong")
	}
	if !NewSet[int]().IsSubsetOf(a) {
		t.Error("the empty set is a subset of everything")
	}
	if !NewSet(2, 1).Equal(NewSet(1, 2)) || a.Equal(b) {
		t.Error("Equal is wrong")
	}
}

func TestSetNilSafety(t *testing.T) {
	var nilSet *Set[int]
	s := NewSet(1)
	if got := s.Union(nilSet).Len(); got != 1 {
		t.Errorf("Union with nil = %d elements, want 1", got)
	}
	if got := s.Intersect(nilSet).Len(); got != 0 {
		t.Errorf("Intersect with nil = %d, want 0", got)
	}
}

func TestMapSetAndSorted(t *testing.T) {
	s := NewSet(1, 2, 3)
	doubled := MapSet(s, func(n int) int { return n % 2 })
	if got := sorted(doubled); !reflect.DeepEqual(got, []int{0, 1}) {
		t.Errorf("MapSet = %v, want [0 1] (collisions collapse)", got)
	}
	names := MapSet(s, func(n int) string { return string(rune('a' + n)) })
	if got := SortedSlice(names); !reflect.DeepEqual(got, []string{"b", "c", "d"}) {
		t.Errorf("SortedSlice = %v", got)
	}
	if got := SortedSlice(NewSet(3, 1, 2)); !reflect.DeepEqual(got, []int{1, 2, 3}) {
		t.Errorf("SortedSlice = %v", got)
	}
}

func TestOption(t *testing.T) {
	s := Some(42)
	if v, ok := s.Get(); !ok || v != 42 {
		t.Errorf("Get = %v, %v", v, ok)
	}
	if !s.IsSome() {
		t.Error("IsSome")
	}
	if got := s.OrElse(0); got != 42 {
		t.Errorf("OrElse = %d", got)
	}
	if got := s.MustGet(); got != 42 {
		t.Errorf("MustGet = %d", got)
	}
	if got := s.String(); got != "Some(42)" {
		t.Errorf("String = %q", got)
	}

	n := None[int]()
	if v, ok := n.Get(); ok || v != 0 {
		t.Errorf("None.Get = %v, %v", v, ok)
	}
	if n.IsSome() {
		t.Error("None.IsSome must be false")
	}
	if got := n.OrElse(7); got != 7 {
		t.Errorf("None.OrElse = %d, want 7", got)
	}
	if got := n.String(); got != "None" {
		t.Errorf("String = %q", got)
	}

	// Some(zero value) is not None. This is the whole point of the type.
	z := Some(0)
	if !z.IsSome() {
		t.Error("Some(0) must be present")
	}
	if got := z.OrElse(99); got != 0 {
		t.Errorf("Some(0).OrElse = %d, want 0", got)
	}
	if got := Some("").String(); got != "Some()" {
		t.Errorf("String = %q", got)
	}
}

func TestOptionMustGetPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("None.MustGet did not panic")
		}
	}()
	None[string]().MustGet()
}

func TestOptionCombinators(t *testing.T) {
	if got := MapOption(Some(3), func(n int) string { return string(rune('a' + n)) }); got.OrElse("") != "d" {
		t.Errorf("MapOption = %v", got)
	}
	called := false
	got := MapOption(None[int](), func(n int) string { called = true; return "x" })
	if got.IsSome() || called {
		t.Error("MapOption over None must not call f")
	}

	half := func(n int) Option[int] {
		if n%2 == 0 {
			return Some(n / 2)
		}
		return None[int]()
	}
	if v := FlatMapOption(Some(8), half); v.OrElse(-1) != 4 {
		t.Errorf("FlatMapOption = %v", v)
	}
	if v := FlatMapOption(Some(7), half); v.IsSome() {
		t.Errorf("FlatMapOption = %v, want None", v)
	}

	if v := Some(4).Filter(func(n int) bool { return n > 10 }); v.IsSome() {
		t.Error("Filter should have removed the value")
	}
	if v := Some(40).Filter(func(n int) bool { return n > 10 }); !v.IsSome() {
		t.Error("Filter removed a matching value")
	}
}
