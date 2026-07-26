// Package funcops: higher-order functions over slices, and the places where
// Go's type inference gives up.
//
// Two things to notice as you work through this:
//   - Methods cannot have their own type parameters. That is why everything
//     here is a function and not a method on some Stream[T] type. You cannot
//     write `func (s Stream[T]) Map[U any](...)`. Try it once so the compiler
//     error is familiar.
//   - Inference flows from arguments, not from return values. Compose needs its
//     middle type parameter written out at the call site in some situations.
package funcops

// Map applies f to every element, returning a new slice. Never returns nil for
// a non-nil input.
func Map[T, U any](s []T, f func(T) U) []U { panic("TODO: implement Map") }

// MapIndex is Map with the index passed along.
func MapIndex[T, U any](s []T, f func(int, T) U) []U { panic("TODO: implement MapIndex") }

// Filter keeps the elements for which keep returns true, in order.
func Filter[T any](s []T, keep func(T) bool) []T { panic("TODO: implement Filter") }

// Reduce folds left: Reduce([1,2,3], 0, add) is ((0+1)+2)+3.
func Reduce[T, A any](s []T, init A, f func(acc A, v T) A) A { panic("TODO: implement Reduce") }

// ReduceRight folds from the right: ReduceRight([1,2,3], "", concat) visits
// 3, then 2, then 1.
func ReduceRight[T, A any](s []T, init A, f func(acc A, v T) A) A {
	panic("TODO: implement ReduceRight")
}

// FlatMap maps each element to a slice and concatenates the results.
func FlatMap[T, U any](s []T, f func(T) []U) []U { panic("TODO: implement FlatMap") }

// Any, All and None short-circuit. All of an empty slice is true; Any is false.
func Any[T any](s []T, pred func(T) bool) bool  { panic("TODO: implement Any") }
func All[T any](s []T, pred func(T) bool) bool  { panic("TODO: implement All") }
func None[T any](s []T, pred func(T) bool) bool { panic("TODO: implement None") }

// Find returns the first match and whether there was one.
func Find[T any](s []T, pred func(T) bool) (T, bool) { panic("TODO: implement Find") }

// IndexFunc returns the index of the first match, or -1.
func IndexFunc[T any](s []T, pred func(T) bool) int { panic("TODO: implement IndexFunc") }

// Uniq removes duplicates by a key function, keeping the first of each.
func Uniq[T any, K comparable](s []T, key func(T) K) []T { panic("TODO: implement Uniq") }

// Zip pairs up two slices, stopping at the shorter one.
type Pair[A, B any] struct {
	First  A
	Second B
}

func Zip[A, B any](as []A, bs []B) []Pair[A, B] { panic("TODO: implement Zip") }

// Unzip is the inverse.
func Unzip[A, B any](ps []Pair[A, B]) ([]A, []B) { panic("TODO: implement Unzip") }

// Compose returns a function that applies f then g. Note the order of the type
// parameters and where inference does and does not work.
func Compose[A, B, C any](f func(A) B, g func(B) C) func(A) C {
	panic("TODO: implement Compose")
}

// Memoize wraps f so that each distinct argument is computed once. The returned
// function is NOT safe for concurrent use - that is module 06's problem.
func Memoize[K comparable, V any](f func(K) V) func(K) V { panic("TODO: implement Memoize") }

// Partial binds the first argument of a two-argument function.
func Partial[A, B, C any](f func(A, B) C, a A) func(B) C { panic("TODO: implement Partial") }
