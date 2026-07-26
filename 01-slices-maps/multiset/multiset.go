// Package multiset builds a counting container on top of a map.
//
// Watch out for: the zero value of a map is nil and writing to it panics;
// map iteration order is deliberately randomised, so anything that returns an
// ordered result has to sort, and ties have to be broken deterministically or
// your tests will flake.
package multiset

// Bag is a multiset: a set that remembers how many times it saw each element.
// The zero Bag must be usable for reads (Count, Len, Items on a zero Bag return
// zero values and must not panic). Only Add needs to lazily create the map.
type Bag[T comparable] struct {
	// TODO: pick your fields. Keeping a running total is cheaper than summing
	// the map every time Size is called.
}

// Add records n more occurrences of v. A negative n removes occurrences. When
// an element's count reaches zero or below it must be deleted from the bag
// entirely, so that Len and Items don't report it.
func (b *Bag[T]) Add(v T, n int) {
	panic("TODO: implement Add")
}

// Count returns how many times v is in the bag (0 if absent).
func (b *Bag[T]) Count(v T) int {
	panic("TODO: implement Count")
}

// Len is the number of distinct elements. Size is the total of all counts.
func (b *Bag[T]) Len() int  { panic("TODO: implement Len") }
func (b *Bag[T]) Size() int { panic("TODO: implement Size") }

// Items returns every distinct element with a positive count, in unspecified
// order. It must return a fresh slice - callers must not be able to corrupt
// the bag through it.
func (b *Bag[T]) Items() []T {
	panic("TODO: implement Items")
}

// Entry pairs an element with its count.
type Entry[T comparable] struct {
	Value T
	Count int
}

// MostCommon returns the n most frequent entries, highest count first. Ties are
// broken by the order the tied elements were FIRST added to the bag (insertion
// order), so the result is deterministic. If n <= 0 or n > Len(), everything is
// returned.
func (b *Bag[T]) MostCommon(n int) []Entry[T] {
	panic("TODO: implement MostCommon")
}

// Union returns a new bag where each element's count is the max of the two
// bags; Intersect uses the min and drops zeroes; Sum adds the counts.
// None of these may modify their receivers or arguments.
func (b *Bag[T]) Union(other *Bag[T]) *Bag[T]     { panic("TODO: implement Union") }
func (b *Bag[T]) Intersect(other *Bag[T]) *Bag[T] { panic("TODO: implement Intersect") }
func (b *Bag[T]) Sum(other *Bag[T]) *Bag[T]       { panic("TODO: implement Sum") }
