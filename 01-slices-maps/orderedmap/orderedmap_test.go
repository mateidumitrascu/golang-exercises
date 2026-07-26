package orderedmap

import (
	"reflect"
	"testing"
)

func TestInsertionOrder(t *testing.T) {
	m := New[string, int]()
	for i, k := range []string{"c", "a", "b"} {
		m.Set(k, i)
	}
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"c", "a", "b"}) {
		t.Errorf("Keys = %v, want [c a b]", got)
	}
	m.Set("c", 99) // update must not reorder
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"c", "a", "b"}) {
		t.Errorf("after update Keys = %v, want [c a b]", got)
	}
	if v, ok := m.Get("c"); !ok || v != 99 {
		t.Errorf("Get(c) = %v, %v; want 99, true", v, ok)
	}
}

func TestGetMissing(t *testing.T) {
	m := New[string, int]()
	if v, ok := m.Get("nope"); ok || v != 0 {
		t.Errorf("Get on empty map = %v, %v", v, ok)
	}
	if m.Len() != 0 {
		t.Error("empty map has non-zero Len")
	}
}

func TestDelete(t *testing.T) {
	m := New[string, int]()
	for i, k := range []string{"a", "b", "c", "d"} {
		m.Set(k, i)
	}
	if !m.Delete("b") {
		t.Error("Delete of an existing key returned false")
	}
	if m.Delete("b") {
		t.Error("Delete of an absent key returned true")
	}
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"a", "c", "d"}) {
		t.Errorf("Keys after delete = %v, want [a c d]", got)
	}
	if m.Len() != 3 {
		t.Errorf("Len = %d, want 3", m.Len())
	}
	// Re-inserting a deleted key puts it at the back.
	m.Set("b", 9)
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"a", "c", "d", "b"}) {
		t.Errorf("Keys after re-insert = %v, want [a c d b]", got)
	}
}

func TestAllIterator(t *testing.T) {
	m := New[string, int]()
	for i, k := range []string{"x", "y", "z"} {
		m.Set(k, i)
	}
	var keys []string
	var vals []int
	for k, v := range m.All() {
		keys = append(keys, k)
		vals = append(vals, v)
	}
	if !reflect.DeepEqual(keys, []string{"x", "y", "z"}) || !reflect.DeepEqual(vals, []int{0, 1, 2}) {
		t.Errorf("iteration = %v %v", keys, vals)
	}
}

func TestAllStopsEarly(t *testing.T) {
	m := New[int, int]()
	for i := range 100 {
		m.Set(i, i)
	}
	seen := 0
	for range m.All() {
		seen++
		if seen == 3 {
			break
		}
	}
	if seen != 3 {
		t.Errorf("iterator visited %d entries after break, want 3", seen)
	}
}

func TestMoveToBackAndOldest(t *testing.T) {
	m := New[string, int]()
	for i, k := range []string{"a", "b", "c"} {
		m.Set(k, i)
	}
	if !m.MoveToBack("a") {
		t.Error("MoveToBack returned false for an existing key")
	}
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"b", "c", "a"}) {
		t.Errorf("Keys = %v, want [b c a]", got)
	}
	k, v, ok := m.Oldest()
	if !ok || k != "b" || v != 1 {
		t.Errorf("Oldest = %v, %v, %v; want b, 1, true", k, v, ok)
	}
	if m.MoveToBack("nope") {
		t.Error("MoveToBack on a missing key returned true")
	}
	empty := New[string, int]()
	if _, _, ok := empty.Oldest(); ok {
		t.Error("Oldest on an empty map returned ok")
	}
}

// TestDeleteIsConstantTime deletes from the front of a large map many times.
// A linear scan per delete turns this into O(n^2) and it will take far too long.
func TestDeleteIsConstantTime(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	const n = 200000
	m := New[int, int]()
	for i := range n {
		m.Set(i, i)
	}
	for i := range n {
		if !m.Delete(i) {
			t.Fatalf("Delete(%d) returned false", i)
		}
	}
	if m.Len() != 0 {
		t.Errorf("Len = %d after deleting everything", m.Len())
	}
}

func BenchmarkSetGet(b *testing.B) {
	m := New[int, int]()
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		m.Set(i, i)
		_, _ = m.Get(i)
		i++
	}
}
