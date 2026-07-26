# 04 — Errors

Go has no exceptions, so error handling is ordinary code — which means you can
design it. This module is about designing it well: classification, wrapping,
context, and the small number of places where `panic` is actually right.

| # | Exercise | Difficulty | What it drills |
|---|----------|-----------|----------------|
| 1 | `wrapping` | ★★☆ | `%w`, `Is`/`As`, custom `Is`, `Join`, walking the error tree |
| 2 | `recovery` | ★★★ | `defer`/`recover`, named returns, per-goroutine recovery |
| 3 | `retry` | ★★★ | Backoff, error classification, `context`, `testing/synctest` |
| 4 | `apierr` | ★★★ | Designing one error type for a whole service (the Upspin pattern) |

## Ideas worth internalising here

**Errors are values, and wrapping builds a tree.** `fmt.Errorf("...: %w", err)`
adds a link; `errors.Join(a, b)` adds a branch. `errors.Is` and `errors.As`
search that tree — so the caller never needs to know how deep the failure was.

**Sentinel or type?** A sentinel (`var ErrNotFound = errors.New(...)`) is for
"which category", a custom type is for "what were the details". `errors.Is` for
the first, `errors.As` for the second. Implementing `Is(target error) bool`
yourself gives you category matching on a structured type — see `HTTPError`.

**`defer` + named return is the only way to convert a panic into an error**,
and `recover()` must be called *directly* by the deferred function. A panic in
another goroutine cannot be recovered by its parent: it takes the process down.
Every long-lived goroutine you spawn needs its own `defer recover` if it can
panic.

**Retry logic is where errors get classified.** "Is this worth doing again?" is
a property of the error, not of the call site — which is why `Retryable() bool`
and `Permanent{}` belong in the error, and why `context.Canceled` must never be
retried.

**`testing/synctest`** makes time-dependent code testable: inside
`synctest.Test`, the clock is fake and advances only when every goroutine is
blocked. A 24-hour backoff test runs in microseconds and asserts exact
durations. This is a Go 1.25+ superpower most people don't know about yet.

## Stretch goals

- Add `retry.Do` support for a circuit breaker: after N consecutive failures,
  fail fast for a cooldown period without calling `fn` at all.
- Give `apierr` a `Log(err)` that prints the internal detail and a request ID
  while `Public(err)` stays safe — the split that keeps you out of trouble.
- Write a `errors.Is`-compatible `Timeout()` interface and see how `net.Error`
  does the same thing in the standard library.
