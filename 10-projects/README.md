# 10 — Projects

The exercises taught you the parts. These are for putting them together: no
stub functions, no tests written for you, no signature telling you what to do.
Just a specification and a blank file — which is exactly the gap between
"I can do the exercise" and "I can build the thing".

| Project | Difficulty | What you'll fight with |
|---------|-----------|------------------------|
| `crawler` | ★★★ | Knowing when a concurrent job is *finished* |
| `kvstore` | ★★★ | Crash safety: a torn write must not lose the whole file |
| `sitegen` | ★★☆ | Deterministic output, golden-file testing, escaping |
| `loadtest` | ★★★ | Measuring latency honestly (coordinated omission) |
| `minigit` | ★★★ | Content addressing and recursive trees |
| `interp` | ★★★ | Pratt parsing and closures over environments |

Each has a `README.md` with the full spec and a `main.go` stub. Run one with:

```
go run ./10-projects/kvstore
```

## How to work on these

**Write the test first for anything you don't understand yet.** Not for
everything — for the part you're unsure about. That's where a test earns its
keep.

**Start with the smallest thing that runs end to end**, then grow it. A crawler
that fetches one page and prints its links is worth more than a beautiful
half-finished worker pool.

**Commit at each milestone.** They're written so that each one leaves you with
something that works.

**When you get stuck, write the type first.** In Go, most design problems are
solved by asking "what is the interface here, and who owns this state?"

**Then go back and read the standard library.** After you've written your own
`errgroup`, read `x/sync/errgroup`. After your own LRU, read
`groupcache/lru`. After your crawler, read the "Web Crawler" section of
*The Go Programming Language*. Reading source you now have opinions about is the
fastest way to improve.

## Definition of done

For each project, before you call it finished:

- `go vet ./...` is clean.
- `go test -race ./...` passes, and there *are* tests.
- `gofmt -l .` prints nothing.
- The tool handles Ctrl-C without leaving a mess.
- A stranger could run it from your README alone.
