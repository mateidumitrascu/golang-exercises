package multiset

import (
	"reflect"
	"sort"
	"testing"
)

func build(words ...string) *Bag[string] {
	b := &Bag[string]{}
	for _, w := range words {
		b.Add(w, 1)
	}
	return b
}

func sortedItems(b *Bag[string]) []string {
	items := b.Items()
	sort.Strings(items)
	return items
}

func TestZeroValueIsReadable(t *testing.T) {
	var b Bag[string]
	if b.Count("x") != 0 || b.Len() != 0 || b.Size() != 0 || len(b.Items()) != 0 {
		t.Error("the zero Bag must behave like an empty bag")
	}
	b.Add("x", 1)
	if b.Count("x") != 1 {
		t.Error("Add on a zero Bag must work")
	}
}

func TestAddAndCount(t *testing.T) {
	b := build("a", "b", "a", "c", "a")
	if got := b.Count("a"); got != 3 {
		t.Errorf(`Count("a") = %d, want 3`, got)
	}
	if got := b.Count("zzz"); got != 0 {
		t.Errorf("missing element gave %d, want 0", got)
	}
	if b.Len() != 3 {
		t.Errorf("Len = %d, want 3", b.Len())
	}
	if b.Size() != 5 {
		t.Errorf("Size = %d, want 5", b.Size())
	}
}

func TestRemoveDeletesKey(t *testing.T) {
	b := build("a", "a", "b")
	b.Add("a", -2)
	if b.Count("a") != 0 {
		t.Errorf(`Count("a") = %d after removing both, want 0`, b.Count("a"))
	}
	if got := sortedItems(b); !reflect.DeepEqual(got, []string{"b"}) {
		t.Errorf("Items = %v, want [b]; zero-count elements must be deleted", got)
	}
	b.Add("b", -100)
	if b.Len() != 0 || b.Size() != 0 {
		t.Errorf("over-removing left Len=%d Size=%d, want 0 0", b.Len(), b.Size())
	}
}

func TestItemsIsACopy(t *testing.T) {
	b := build("a", "b")
	items := b.Items()
	if len(items) != 2 {
		t.Fatalf("Items = %v", items)
	}
	items[0] = "mutated"
	if got := sortedItems(b); !reflect.DeepEqual(got, []string{"a", "b"}) {
		t.Errorf("mutating the Items slice changed the bag: %v", got)
	}
}

func TestMostCommon(t *testing.T) {
	b := build("go", "rust", "go", "zig", "go", "rust", "c")
	got := b.MostCommon(2)
	want := []Entry[string]{{"go", 3}, {"rust", 2}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MostCommon(2) = %v, want %v", got, want)
	}

	// zig and c are tied at 1; zig was added first, so it comes first.
	all := b.MostCommon(0)
	wantAll := []Entry[string]{{"go", 3}, {"rust", 2}, {"zig", 1}, {"c", 1}}
	if !reflect.DeepEqual(all, wantAll) {
		t.Errorf("MostCommon(0) = %v, want %v (ties break by insertion order)", all, wantAll)
	}
	if got := b.MostCommon(100); len(got) != 4 {
		t.Errorf("MostCommon(100) returned %d entries, want 4", len(got))
	}
}

func TestMostCommonIsDeterministic(t *testing.T) {
	b := build("a", "b", "c", "d", "e", "f", "g", "h")
	first := b.MostCommon(0)
	for range 50 {
		if got := b.MostCommon(0); !reflect.DeepEqual(got, first) {
			t.Fatalf("MostCommon is not deterministic:\n%v\n%v\n(map iteration order is random - you must sort)", got, first)
		}
	}
}

func TestSetOps(t *testing.T) {
	a := build("x", "x", "y")
	c := build("x", "z", "z", "z")

	check := func(name string, got *Bag[string], want map[string]int) {
		t.Helper()
		if got.Len() != len(want) {
			t.Errorf("%s: Len = %d, want %d", name, got.Len(), len(want))
		}
		for k, v := range want {
			if got.Count(k) != v {
				t.Errorf("%s: Count(%q) = %d, want %d", name, k, got.Count(k), v)
			}
		}
	}
	check("Union", a.Union(c), map[string]int{"x": 2, "y": 1, "z": 3})
	check("Intersect", a.Intersect(c), map[string]int{"x": 1})
	check("Sum", a.Sum(c), map[string]int{"x": 3, "y": 1, "z": 3})

	if a.Count("z") != 0 || c.Count("y") != 0 {
		t.Error("set operations must not modify their operands")
	}
}
