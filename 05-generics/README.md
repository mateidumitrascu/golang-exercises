# 05 — Generics and iterators

Generics arrived in 1.18, range-over-func in 1.23. Together they change how you
write reusable code in Go. This module is about using them the way the standard
library does — sparingly, with small constraints, and without pretending Go is
Haskell.

| # | Exercise | Difficulty | What it drills |
|---|----------|-----------|----------------|
| 1 | `numeric` | ★★☆ | Writing constraints, `~` approximation, inference across two params |
| 2 | `funcops` | ★★☆ | Map/Filter/Reduce, short-circuiting, closures, memoisation |
| 3 | `containers` | ★★★ | Generic types with methods, and the no-type-params-on-methods rule |
| 4 | `iterseq` | ★★★ | `iter.Seq`, laziness, propagating `break`, `iter.Pull` |

## Ideas worth internalising here

**`~int` vs `int`.** A constraint of `int` accepts only `int`. `~int` accepts
every type whose *underlying* type is `int`, which is what makes your function
usable with `type UserID int`. Almost always you want the tilde.

**Methods cannot have type parameters.** `func (s *Set[T]) Map[U any](...)` is a
compile error. Anything that changes the element type is a free function. This
single rule explains why Go's generic APIs look like `slices.SortFunc(s, cmp)`
instead of `s.SortFunc(cmp)`.

**Inference flows from arguments.** `Map(s, f)` infers both parameters from `s`
and `f`. But a function whose type parameter appears only in the *return* type
must be instantiated explicitly: `None[int]()`.

**An iterator is just a function.** `iter.Seq[T]` is `func(yield func(T) bool)`.
Everything else — laziness, composition, early exit — falls out of that. The one
obligation is to stop calling `yield` the moment it returns `false`, and to
propagate that up the chain.

**`iter.Pull` runs the producer on a goroutine.** That's why it returns a `stop`
function and why you always `defer stop()` — otherwise a consumer that stops
early leaks a goroutine. Prove it to yourself with `-race` and a leak check.

**Don't reach for generics first.** If a concrete type would do, use it. The
best generic code in Go is boring: containers, and functions over slices/maps.

## Stretch goals

- Give `Set` a `MarshalJSON` that emits a sorted array — and notice you can't
  constrain `T` to `Ordered` there, so you need `cmp.Compare` with a type switch
  or a `less` function stored in the set.
- Write `Parallel(seq iter.Seq[T], n int, f func(T) U) iter.Seq[U]` after
  module 06 and decide what "order" means for it.
- Reimplement `Chunk` so it reuses one buffer, add a doc comment warning about
  it, and write the test that catches a caller who retains a chunk.
