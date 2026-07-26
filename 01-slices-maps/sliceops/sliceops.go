// Package sliceops re-implements the parts of the standard "slices" package
// that people most often get wrong.
//
// Do NOT import "slices" here. The point is to write the pointer arithmetic
// yourself so that append's growth behaviour stops being magic.
package sliceops

// InsertAt returns s with the values v inserted at index i, shifting the rest
// right. It may reuse s's backing array. InsertAt panics if i is out of
// [0, len(s)].
//
//	InsertAt([1,2,5], 2, 3, 4) == [1,2,3,4,5]
func InsertAt[T any](s []T, i int, v ...T) []T {
	panic("TODO: implement InsertAt")
}

// DeleteRange returns s with elements [i, j) removed, preserving order and
// zeroing the freed tail. It panics if the range is invalid.
func DeleteRange[T any](s []T, i, j int) []T {
	panic("TODO: implement DeleteRange")
}

// DeleteUnordered removes the element at index i in O(1) by moving the last
// element into the hole. Order is not preserved. The tail must be zeroed.
func DeleteUnordered[T any](s []T, i int) []T {
	panic("TODO: implement DeleteUnordered")
}

// Clone returns a copy of s with its own backing array. Clone(nil) is nil, and
// Clone of an empty-but-non-nil slice is empty and non-nil.
func Clone[T any](s []T) []T {
	panic("TODO: implement Clone")
}

// Grow returns a slice with the same contents as s and capacity for at least n
// more elements without another allocation. If s already has room, s is
// returned untouched (same backing array).
func Grow[T any](s []T, n int) []T {
	panic("TODO: implement Grow")
}

// Equal reports whether two slices have the same length and elements.
// A nil slice equals an empty slice.
func Equal[T comparable](a, b []T) bool {
	panic("TODO: implement Equal")
}

// Reslice is a teaching function: return a slice sharing s's memory that
// contains s[i:j] but whose capacity is exactly j-i, so that appending to the
// result can never be seen through s.
func Reslice[T any](s []T, i, j int) []T {
	panic("TODO: implement Reslice")
}
