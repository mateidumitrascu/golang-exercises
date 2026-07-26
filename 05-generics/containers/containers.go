// Package containers: generic types with methods, and the one rule that shapes
// all generic API design in Go.
//
// THE RULE: methods cannot introduce new type parameters. `func (s *Set[T])
// Map[U any](f func(T) U) *Set[U]` does not compile. Anything that changes the
// element type has to be a free function. Write the illegal version once and
// read the compiler error - it will save you an hour some day.
package containers

import "iter"

// Set is an unordered collection of distinct values.
//
// The zero value is usable for reads: Has, Len and All must work on a
// `var s Set[int]` without panicking. Only the mutating methods create the map.
type Set[T comparable] struct {
	// TODO: your field(s)
}

// NewSet builds a set from values.
func NewSet[T comparable](vals ...T) *Set[T] { panic("TODO: implement NewSet") }

func (s *Set[T]) Add(vals ...T) *Set[T] { panic("TODO: implement Add") } // returns s, for chaining
func (s *Set[T]) Delete(vals ...T)      { panic("TODO: implement Delete") }
func (s *Set[T]) Has(v T) bool          { panic("TODO: implement Has") }
func (s *Set[T]) Len() int              { panic("TODO: implement Len") }
func (s *Set[T]) Clear()                { panic("TODO: implement Clear") }
func (s *Set[T]) Clone() *Set[T]        { panic("TODO: implement Clone") }

// All iterates the elements in unspecified order. It must honour early
// termination, and iterating the zero value must yield nothing.
func (s *Set[T]) All() iter.Seq[T] { panic("TODO: implement All") }

// The set algebra. None of these modify their receiver or argument, and all of
// them handle nil receivers/arguments as empty sets.
func (s *Set[T]) Union(o *Set[T]) *Set[T]         { panic("TODO: implement Union") }
func (s *Set[T]) Intersect(o *Set[T]) *Set[T]     { panic("TODO: implement Intersect") }
func (s *Set[T]) Difference(o *Set[T]) *Set[T]    { panic("TODO: implement Difference") }
func (s *Set[T]) SymmetricDiff(o *Set[T]) *Set[T] { panic("TODO: implement SymmetricDiff") }

func (s *Set[T]) IsSubsetOf(o *Set[T]) bool { panic("TODO: implement IsSubsetOf") }
func (s *Set[T]) Equal(o *Set[T]) bool      { panic("TODO: implement Equal") }

// MapSet is Set.Map, except it cannot be a method, because it changes the
// element type. This is the workaround, and it is why the standard library is
// full of free functions like slices.SortFunc.
func MapSet[T, U comparable](s *Set[T], f func(T) U) *Set[U] { panic("TODO: implement MapSet") }

// SortedSlice returns the elements in ascending order. It also cannot be a
// method: Set's T is only `comparable`, and sorting needs `cmp.Ordered`, which
// a method has no way to add.
func SortedSlice[T interface{ ~int | ~string }](s *Set[T]) []T {
	panic("TODO: implement SortedSlice")
}

// Option is a value that may be absent. It is not idiomatic Go for error
// handling - use (T, error) for that - but it is an excellent exercise in
// generic types, and it is genuinely useful for "field not set" cases.
type Option[T any] struct {
	// TODO: your fields. Note you need to distinguish "set to the zero value"
	// from "not set".
}

func Some[T any](v T) Option[T] { panic("TODO: implement Some") }
func None[T any]() Option[T]    { panic("TODO: implement None") }

func (o Option[T]) Get() (T, bool) { panic("TODO: implement Get") }
func (o Option[T]) IsSome() bool   { panic("TODO: implement IsSome") }
func (o Option[T]) OrElse(def T) T { panic("TODO: implement OrElse") }

// MustGet panics if there is no value.
func (o Option[T]) MustGet() T { panic("TODO: implement MustGet") }

// Filter keeps the value only if it matches. This one CAN be a method: it does
// not introduce a new type parameter.
func (o Option[T]) Filter(pred func(T) bool) Option[T] { panic("TODO: implement Filter") }

// String is "Some(42)" or "None", using %v for the value.
func (o Option[T]) String() string { panic("TODO: implement String") }

// MapOption transforms the value if present. A free function, for the same
// reason as MapSet.
func MapOption[A, B any](o Option[A], f func(A) B) Option[B] { panic("TODO: implement MapOption") }

// FlatMapOption chains operations that may themselves be absent.
func FlatMapOption[A, B any](o Option[A], f func(A) Option[B]) Option[B] {
	panic("TODO: implement FlatMapOption")
}
