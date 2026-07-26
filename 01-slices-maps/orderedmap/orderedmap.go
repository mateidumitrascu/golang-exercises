// Package orderedmap: a map that remembers insertion order, with O(1) Get, Set
// and Delete.
//
// The naive version (a map plus a []K that you linear-scan on delete) is O(n)
// per delete and will fail the benchmark-based test. Two designs work:
//
//	map[K]*node + doubly linked list of nodes
//	map[K]int index into a slice of entries + tombstones you compact amortised
//
// Pick one and defend it. The iterator must be lazy: it yields entries in
// insertion order and stops early if the caller stops ranging.
package orderedmap

import "iter"

// Map is an insertion-ordered map. The zero value is not usable; call New.
type Map[K comparable, V any] struct {
	// TODO: your fields
}

func New[K comparable, V any]() *Map[K, V] {
	panic("TODO: implement New")
}

// Set inserts or updates k. Updating an existing key must NOT change its
// position in the ordering.
func (m *Map[K, V]) Set(k K, v V) {
	panic("TODO: implement Set")
}

// Get returns the value and whether it was present.
func (m *Map[K, V]) Get(k K) (V, bool) {
	panic("TODO: implement Get")
}

// Delete removes k and reports whether it was there. Must be O(1).
func (m *Map[K, V]) Delete(k K) bool {
	panic("TODO: implement Delete")
}

// Len is the number of live entries.
func (m *Map[K, V]) Len() int {
	panic("TODO: implement Len")
}

// Keys returns the keys in insertion order.
func (m *Map[K, V]) Keys() []K {
	panic("TODO: implement Keys")
}

// All is a range-over-func iterator over the entries in insertion order:
//
//	for k, v := range m.All() { ... }
//
// It must honour early termination (when yield returns false, stop).
func (m *Map[K, V]) All() iter.Seq2[K, V] {
	panic("TODO: implement All")
}

// MoveToBack moves an existing key to the end of the ordering and reports
// whether it existed. (This is the operation an LRU cache needs.)
func (m *Map[K, V]) MoveToBack(k K) bool {
	panic("TODO: implement MoveToBack")
}

// Oldest returns the first-inserted live entry.
func (m *Map[K, V]) Oldest() (K, V, bool) {
	panic("TODO: implement Oldest")
}
