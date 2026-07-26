// Package rotate: rearranging a slice with O(1) extra memory.
//
// The easy version allocates a temporary slice. Don't do the easy version -
// TestNoAllocations will fail you. Two classic techniques work:
//
//	the three-reverse trick   reverse(a) reverse(b) reverse(whole)
//	the juggling algorithm    gcd(n, k) cycles of element shuffling
package rotate

// Reverse reverses s in place.
func Reverse[T any](s []T) {
	panic("TODO: implement Reverse")
}

// RotateLeft shifts every element k positions towards index 0, wrapping around.
// RotateLeft([1,2,3,4,5], 2) leaves s as [3,4,5,1,2].
//
// k may be negative (rotate right instead) or larger than len(s). Rotating an
// empty or one-element slice does nothing. Must not allocate.
func RotateLeft[T any](s []T, k int) {
	panic("TODO: implement RotateLeft")
}

// RotateRight is RotateLeft in the other direction.
func RotateRight[T any](s []T, k int) {
	panic("TODO: implement RotateRight")
}

// IsRotation reports whether b is some rotation of a (including zero).
// [3,4,1,2] is a rotation of [1,2,3,4]; [1,2,4,3] is not.
// Must run in O(n) time on average and must not modify either input.
//
// Hint: there is a one-line trick involving concatenation and substring search,
// and a from-scratch version that does not build the doubled slice.
func IsRotation[T comparable](a, b []T) bool {
	panic("TODO: implement IsRotation")
}
