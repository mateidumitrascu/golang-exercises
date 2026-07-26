# Go exercises

Forty-five test-driven exercises and six capstone projects, built to take you from
"I can write Go" to "I can build a Go system without help".

There are no solutions in this repo, on purpose. **The tests are the
specification.** Every function has a doc comment telling you what it must do
and a test that checks it — including the edge cases you'd forget, the
allocation budgets, and the concurrency guarantees.

```
go test ./01-slices-maps/chunk/     # one exercise
./check.sh                          # everything, with a summary
./check.sh 06                       # one module
./check.sh -v 06                    # ...with the full output
```

Everything is standard library only. No dependencies, works offline. Go 1.26.

## The modules

| # | Module | Exercises | What it's for |
|---|--------|-----------|---------------|
| 01 | [slices-maps](01-slices-maps/) | 8 | Aliasing, capacity, tail-zeroing, map iteration order |
| 02 | [strings-bytes](02-strings-bytes/) | 5 | UTF-8 by hand, `bufio.SplitFunc`, state machines |
| 03 | [types-interfaces](03-types-interfaces/) | 5 | Method sets, typed nil, sum types, reflection |
| 04 | [errors](04-errors/) | 4 | Wrapping, `recover`, retries, designing an error type |
| 05 | [generics](05-generics/) | 4 | Constraints, generic containers, `iter.Seq` pipelines |
| 06 | [concurrency](06-concurrency/) | 6 | Pools, pipelines, `sync` from scratch, rate limits, races |
| 07 | [io-encoding](07-io-encoding/) | 4 | Reader/Writer contracts, wire formats, JSON, `io/fs` |
| 08 | [http-net](08-http-net/) | 4 | Services, middleware, clients, raw TCP |
| 09 | [datastructures](09-datastructures/) | 5 | LRU, heaps, tries, graphs, the classic algorithms |
| 10 | [projects](10-projects/) | 6 | Specs only. No stubs, no tests, no hand-holding. |

Each module has a `README.md` with a table of its exercises, the ideas it's
trying to teach, and stretch goals for when you want more.

Difficulty is marked ★★☆ (a solid hour) to ★★★ (an evening, and you'll learn
something). There is nothing here below ★★☆.

## How to work through it

**Pick the module you're weakest at, not module 01.** They're independent. If
you already know slices cold, start at 04 or 06. The only soft dependency is
that module 10's projects lean on things you built earlier.

**Read the doc comment, then the test, then write the code.** The doc comment
tells you the contract; the test tells you what the contract means in practice.
When they disagree, the test wins — and you should ask why.

**Run the tests before you write anything.** Every stub panics with `TODO`, so
you'll see the whole failing surface at once. That's the map.

**Use `-race` in modules 06 and 08.** Not optional. Several bugs there are
invisible without it:

```
go test -race ./06-concurrency/...
```

**Don't look things up until you've been stuck for twenty minutes.** Then look
up the *standard library source*, not a tutorial. Most of these exercises are
re-implementations of something in the stdlib, and reading `src/sync/once.go`
after you've written your own `Once` is worth ten blog posts.

**When an exercise is done, read the module README's "ideas" section.** It names
the thing you just learned, which is how it sticks.

## The rules that make this work

1. **No AI, no Stack Overflow, for the first attempt.** The whole point is to
   build the muscle of getting unstuck yourself. Once a test passes, comparing
   notes with anything you like is fair game — that's where the second half of
   the learning is.
2. **No `golang.org/x/...` and no third-party packages.** Where an exercise
   re-implements one (`errgroup`, `singleflight`, `rate.Limiter`, `slices`),
   that's the exercise.
3. **`gofmt -l .` prints nothing** before you call something finished.
4. **`go vet ./...` is clean** — with one deliberate exception in
   `06-concurrency/races`, which exists so you see vet catch a real bug.
5. **Make the test pass for the right reason.** Special-casing the test input is
   always possible and always a waste of your evening.

## Things that will trip you up (and are meant to)

- Tests that assert an **allocation count** (`testing.AllocsPerRun`) — a correct
  answer computed wastefully still fails.
- Tests that assert a **time complexity** by making n large — the O(n²) version
  simply won't finish.
- Tests that count **goroutines before and after** — a leak is a failure, not a
  smell.
- `testing/synctest` tests where time is **virtual** and durations are asserted
  **exactly**. A one-hour backoff test runs in microseconds.
- **Fuzz targets.** Run them when the unit tests pass, and expect to be
  embarrassed:
  ```
  go test -run xxx -fuzz FuzzFirstRune ./02-strings-bytes/runes/
  ```

## Tracking progress

`PROGRESS.md` has a checklist of every exercise. Tick them off as you go;
regenerate it with `./gen_progress.sh` if you rearrange things.

```
./check.sh

  passing packages:  12
  failing packages:  37
  (203 tests still hitting a TODO stub)
```

## If you want a suggested order

Slower, thorough route: 01 → 02 → 03 → 04 → 05 → 06 → 07 → 08 → 09 → 10.

Faster route if you're already comfortable and want the hard parts:
**06** (concurrency) → **04** (errors) → **07** (I/O) → **08** (HTTP) →
**09** (data structures) → **10** (a project). Then backfill 01–03 and 05 as
reference when something surprises you.

Either way: finish a project. The exercises make you fluent; the project is
where you find out what fluency was for.
