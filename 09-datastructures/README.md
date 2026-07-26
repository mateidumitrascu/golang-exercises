# 09 — Data structures and algorithms

Not for interview practice — for the day you need a cache that won't blow up
memory, a scheduler that pops the right job, or an autocomplete that doesn't scan
every key. Written in Go, with Go's constraints and Go's edge cases.

| # | Exercise | Difficulty | What it drills |
|---|----------|-----------|----------------|
| 1 | `lru` | ★★★ | Map + linked list for O(1), TTL, injectable clock, eviction hooks |
| 2 | `heapq` | ★★★ | Binary heap, priority queue with `Update`, top-k in O(n log k) |
| 3 | `trie` | ★★★ | Prefix tree, pruning on delete, wildcard DFS, lexical iteration |
| 4 | `graph` | ★★★ | BFS/DFS, Kahn's topological sort, cycle colours, Dijkstra with a PQ |
| 5 | `algo` | ★★★ | Binary search, two pointers, sliding window, dynamic programming |

## Ideas worth internalising here

**The complexity is the specification.** Several tests here don't check the
answer, they check the *cost*: `TestConstantTime` on the LRU, `TestTopKMemory`,
`TestDijkstraScales`, `TestLISIsFast`. A correct O(n²) solution fails. That is
deliberate — in real systems the wrong complexity is the bug.

**`lo + (hi-lo)/2`**, not `(lo+hi)/2`. The second one overflows, and it stayed
broken in the JDK for nine years.

**Deterministic output matters.** Anything that iterates a Go map produces a
different answer every run. Sorting neighbours, and breaking ties by the
smallest node, is what makes a graph algorithm testable at all.

**Injectable clocks beat `time.Sleep` in tests.** The LRU takes a
`now func() time.Time`; the test moves time forward by an hour instantly. Do
this in production code too — the alternative is a slow, flaky test suite.

**Cycle detection needs three colours, not two.** White/grey/black: a node you
are *currently* visiting is not the same as one you have *finished*. Two colours
report a diamond DAG as cyclic, and the test catches exactly that.

## Stretch goals

- Make `lru` concurrency-safe and benchmark it against a sharded version (16
  shards, each with its own mutex). Where is the crossover?
- Add `graph.BellmanFord` which *does* allow negative weights and detects
  negative cycles — then explain why Dijkstra can't.
- Extend `trie` into a radix tree (compress single-child chains) and measure the
  memory difference on a real word list from `/usr/share/dict/words`.
