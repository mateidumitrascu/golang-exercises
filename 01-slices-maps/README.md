# 01 — Slices, maps, and memory

The container types you already use every day, from the angle you probably
haven't: aliasing, capacity, growth, tail-zeroing, and the fact that map
iteration order is deliberately random.

Run one exercise: `go test ./01-slices-maps/chunk/`
Run the module: `./check.sh 01`

| # | Exercise | Difficulty | What it drills |
|---|----------|-----------|----------------|
| 1 | `chunk` | ★★☆ | Three-index slice expressions, sharing a backing array safely |
| 2 | `dedup` | ★★☆ | In-place filtering, the read/write index idiom, zeroing the tail |
| 3 | `rotate` | ★★★ | O(1)-space rotation, `testing.AllocsPerRun` as a spec |
| 4 | `sliceops` | ★★☆ | Re-implementing `slices`: insert, delete, clone, grow |
| 5 | `matrix` | ★★★ | 2D layout, one flat array vs many, spiral traversal |
| 6 | `multiset` | ★★☆ | Map-backed container, deterministic ordering, zero-value usability |
| 7 | `groupby` | ★★☆ | Generic collection helpers with two type parameters |
| 8 | `orderedmap` | ★★★ | O(1) ordered map, `iter.Seq2`, early-termination iterators |

## Ideas worth internalising here

**A slice is a struct.** `{ptr, len, cap}`, passed by value. Passing a slice to a
function copies the header, not the data — which is why a function can change
your elements but not your length.

**`append` is not pure.** If `cap > len`, it writes into memory you may still be
sharing with someone else. `s[i:j:j]` caps the capacity and turns that hazard
into a guaranteed copy. Exercise 1 has a test that only passes if you know this.

**Zero the tail.** `s = s[:len(s)-1]` leaves the removed element reachable from
the backing array. For `[]*T` or `[]string` that is a genuine leak: the GC can't
collect what the array still points at.

**Growing is amortised, not free.** `make([]T, 0, n)` when you know `n`.
Two of the tests here assert an allocation count, not just a result.

**Map iteration order is randomised on purpose.** Any function that returns a
"top N" or an ordered view has to sort, and needs a documented tie-break rule,
or its output is non-deterministic and its tests flake.

## Stretch goals

- Make `Chunk` return an `iter.Seq[[]T]` instead of a `[][]T`. Which one
  allocates less? Prove it with a benchmark.
- Add a `Bag.Filter(func(T, int) bool) *Bag[T]` that returns a new bag.
- Give `orderedmap` a `MarshalJSON` that preserves the insertion order —
  `encoding/json` sorts map keys, so you need to build the JSON yourself.
