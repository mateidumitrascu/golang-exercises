// Package chunk is about carving slices up without copying the elements.
//
// The whole point of this exercise is the three-index slice expression
// s[low:high:max]. If you don't use it, one of the tests will catch you.
package chunk

// Chunk splits s into consecutive groups of at most size elements. The last
// group may be shorter. Chunk panics if size <= 0.
//
// Rules:
//   - The returned groups must SHARE memory with s. Writing g[0] = x through a
//     returned group must be visible in s. (So: no copying.)
//   - Appending to a returned group must never overwrite an element of s.
//     Think about what append does when there is spare capacity.
//   - Chunk(nil, n) and Chunk([]T{}, n) return an empty result (len 0).
func Chunk[T any](s []T, size int) [][]T {
	panic("TODO: implement Chunk")
}

// Windows returns every contiguous window of exactly size elements, in order.
// Windows([1,2,3,4], 2) == [[1,2],[2,3],[3,4]].
//
// If size > len(s) the result is empty. Windows panics if size <= 0.
// The same memory-sharing and append-safety rules as Chunk apply.
func Windows[T any](s []T, size int) [][]T {
	panic("TODO: implement Windows")
}

// Split divides s into exactly n groups whose lengths differ by at most one.
// Earlier groups get the extra elements: Split([1..10], 3) has lengths 4,3,3.
//
// If n > len(s), the trailing groups are empty (but the result still has
// exactly n entries). Split panics if n <= 0. Memory is shared with s.
func Split[T any](s []T, n int) [][]T {
	panic("TODO: implement Split")
}
