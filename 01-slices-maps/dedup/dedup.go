// Package dedup is about editing a slice in place.
//
// Every function here uses the "filter in place" idiom: you walk the slice with
// a read index and a write index and return s[:w]. No new backing arrays.
package dedup

// Dedup removes duplicate values from s, keeping the FIRST occurrence of each
// value and preserving the order of what remains. It returns the shortened
// slice, which shares the backing array with s.
//
// Dedup([3,1,3,2,1]) == [3,1,2].
//
// Requirement: the elements past the returned length must be zeroed. A slice of
// pointers whose tail still points at dead objects keeps them alive - that is a
// real memory leak, and TestDedupZeroesTail checks for it.
func Dedup[T comparable](s []T) []T {
	panic("TODO: implement Dedup")
}

// DedupSorted does the same for an already-sorted slice, in O(n) time and O(1)
// space, comparing only neighbours. Behaviour is undefined if s is not sorted.
func DedupSorted[T comparable](s []T) []T {
	panic("TODO: implement DedupSorted")
}

// DeleteIf removes every element for which drop returns true, in place,
// preserving order. It returns the shortened slice and zeroes the tail.
func DeleteIf[T any](s []T, drop func(T) bool) []T {
	panic("TODO: implement DeleteIf")
}

// Compact collapses runs of equal neighbouring elements down to one, using eq
// to compare. Compact is the general form of DedupSorted:
//
//	Compact([1,1,2,2,2,1], ==) == [1,2,1]
func Compact[T any](s []T, eq func(a, b T) bool) []T {
	panic("TODO: implement Compact")
}
