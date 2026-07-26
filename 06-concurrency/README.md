# 06 — Concurrency

Goroutines are cheap; correctness is not. This module covers the patterns you
actually need (bounded parallelism, pipelines, cancellation, rate limiting) and
the failure modes that come with them (leaks, races, deadlocks, lost errors).

**Always run these with the race detector:**

```
go test -race ./06-concurrency/...
```

| # | Exercise | Difficulty | What it drills |
|---|----------|-----------|----------------|
| 1 | `pool` | ★★★ | Bounded parallelism, ordered results, first-error cancellation |
| 2 | `pipeline` | ★★★ | Stage ownership, `select` on `ctx.Done()`, leak-free composition |
| 3 | `syncprim` | ★★★ | Writing your own `Once`, RWMutex map, CAS loops, singleflight |
| 4 | `ratelimit` | ★★★ | Semaphores, token buckets, `testing/synctest` |
| 5 | `shutdown` | ★★★ | `errgroup` from scratch, reverse-order time-boxed shutdown |
| 6 | `races` | ★★★ | Six real bugs, already written, for you to find and fix |

## Ideas worth internalising here

**The sender closes the channel.** Receivers never close. If two goroutines
might close, you have designed it wrong — or you need a `sync.Once`.

**Every goroutine needs an exit story.** "It sends on a channel nobody is
reading any more" is the most common leak in Go. `select { case out <- v: case
<-ctx.Done(): return }` is the shape that fixes it. Leaks are silent: nothing
crashes, memory just grows. That is why several tests here count goroutines.

**Don't communicate by sharing memory — but do share memory when it's simpler.**
A mutex around a map is fine and often faster than a channel-based actor. The
rule is: pick one owner per piece of state and be able to say who it is.

**`sync.WaitGroup.Add` goes before `go`**, always. Inside the goroutine it races
with `Wait`.

**Cancellation is cooperative.** A context can only stop code that checks it.
Anything that blocks — a channel op, a sleep, a syscall — needs a `select` with
`ctx.Done()` or a deadline.

**The race detector finds data races, not logic races.** `-race` will never
notice that your check-then-act cache computed the same value twice. Bug 6 in
`races` is that kind, on purpose.

**`testing/synctest` makes concurrency deterministic.** Inside a bubble, time is
virtual and `synctest.Wait()` blocks until every other goroutine is idle. Timing
tests become exact instead of flaky, and a one-hour timeout test runs instantly.

## Stretch goals

- Add `ParallelMap` support for *streaming* input (`<-chan T` in, `<-chan U`
  out) while keeping order. Notice how much harder ordering gets.
- Extend `Bucket` with `Reserve(n)` returning a cancellable reservation, the way
  `golang.org/x/time/rate` does.
- Benchmark `syncprim.Map` against `sync.Map` for read-heavy and write-heavy
  workloads. Explain the crossover point.
