// Package numeric is about writing constraints, not just using them.
//
// A constraint is an interface whose type set you define. `~int` means "any
// type whose underlying type is int", which is what makes your code work for
// `type Celsius float64` and not just for float64 itself. Forget the tilde and
// your library is useless to anyone with a named type - a very common bug.
//
// Do not import golang.org/x/exp/constraints. Write them.
package numeric

// Define these constraints yourself.
//
//	Signed    ~int ~int8 ~int16 ~int32 ~int64
//	Unsigned  ~uint ~uint8 ~uint16 ~uint32 ~uint64 ~uintptr
//	Integer   Signed | Unsigned
//	Float     ~float32 ~float64
//	Number    Integer | Float
//	Ordered   Integer | Float | ~string
type Signed interface{ TODOReplaceMe }
type Unsigned interface{ TODOReplaceMe }
type Integer interface{ TODOReplaceMe }
type Float interface{ TODOReplaceMe }
type Number interface{ TODOReplaceMe }
type Ordered interface{ TODOReplaceMe }

// TODOReplaceMe exists only so this file compiles before you start. Delete it
// once you have written the constraints above.
type TODOReplaceMe interface{}

// Sum adds everything up, in the element type - so summing []uint8 wraps around
// at 256, exactly as Go's arithmetic does. Sum of an empty slice is 0.
func Sum[T Number](s []T) T { panic("TODO: implement Sum") }

// Mean returns the arithmetic mean as a float64, and false for an empty slice.
func Mean[T Number](s []T) (float64, bool) { panic("TODO: implement Mean") }

// MinOf and MaxOf are variadic and report false when given nothing.
func MinOf[T Ordered](vals ...T) (T, bool) { panic("TODO: implement MinOf") }
func MaxOf[T Ordered](vals ...T) (T, bool) { panic("TODO: implement MaxOf") }

// MinMax returns both bounds in a single pass.
func MinMax[T Ordered](s []T) (lo, hi T, ok bool) { panic("TODO: implement MinMax") }

// Clamp constrains v to [lo, hi]. It panics if lo > hi.
func Clamp[T Ordered](v, lo, hi T) T { panic("TODO: implement Clamp") }

// Abs works for signed integers and floats. Note that Abs of the most negative
// integer cannot be represented - return it unchanged rather than pretending.
func Abs[T Signed | Float](v T) T { panic("TODO: implement Abs") }

// SumBy sums a projection of the elements. Two type parameters, one of which
// Go can only infer from the function argument - a good demonstration of what
// type inference can and cannot do.
func SumBy[T any, N Number](s []T, f func(T) N) N { panic("TODO: implement SumBy") }

// Compare returns -1, 0 or +1. It is what sort and slices.SortFunc want.
func Compare[T Ordered](a, b T) int { panic("TODO: implement Compare") }

// InDelta reports whether two floats are within tolerance of each other -
// the only sane way to compare floats.
func InDelta[T Float](a, b, tolerance T) bool { panic("TODO: implement InDelta") }
