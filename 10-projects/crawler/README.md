# Project: concurrent web crawler

**Modules it exercises:** 05 (iterators), 06 (concurrency), 07 (I/O), 08 (HTTP).
**Rough size:** 400–700 lines. **Difficulty:** ★★★

Crawl a site starting from one URL, follow links, and report what you found —
concurrently, politely, and cancellably.

## The spec

```
crawler -url https://example.com -depth 3 -workers 8 -rate 5 -timeout 30s
```

Output, sorted by URL:

```
200  1.2 KB  https://example.com/
200  4.5 KB  https://example.com/about
404  0.0 KB  https://example.com/missing
ERR  dial tcp: i/o timeout  https://example.com/slow
```

Then a summary: pages crawled, total bytes, errors, wall time, and the deepest
path found.

## Requirements

1. **Bounded concurrency.** Exactly `-workers` fetches in flight, no more, no
   matter how many links are found. No goroutine-per-URL.
2. **No duplicate work.** Each URL is fetched once, even when three workers
   discover it simultaneously. Normalise first: strip the fragment, sort query
   parameters, resolve relative links against the page they came from.
3. **Same-host only** by default, with a `-external` flag to allow leaving.
4. **Rate limiting.** At most `-rate` requests per second overall (reuse your
   token bucket from 06).
5. **Depth limiting.** The seed is depth 0.
6. **Cancellation.** Ctrl-C stops cleanly: in-flight fetches are cancelled, the
   summary still prints. `signal.NotifyContext` is two lines.
7. **No leaked goroutines.** After the crawl returns, the goroutine count is
   back to its starting value. Prove it in a test.
8. **robots.txt**, at least the `Disallow:` lines for `User-agent: *`.

## Milestones

1. Fetch one URL and extract its links. Use `golang.org/x/net/html`… no — no
   dependencies. Write a small link extractor over `href="..."` with
   `strings.Index`, or a tiny state machine. It only has to be good enough.
2. Sequential recursive crawl with a `map[string]bool` of visited URLs.
3. Worker pool + results channel. Now find the deadlock you just wrote: the
   workers block sending results while the collector blocks sending work.
4. Fix it with a queue owned by one coordinator goroutine, or a `sync.WaitGroup`
   that counts outstanding URLs rather than a fixed set.
5. Add the rate limiter, depth, robots, cancellation, summary.

## The interesting problem

The crawl is done when the queue is empty *and* no worker is busy — not when the
queue is empty. Getting that termination condition right, without a race and
without a sleep-poll loop, is the whole project. Two good answers: a counter of
outstanding items guarded by a mutex plus a condition, or a single coordinator
goroutine that owns the queue and the in-flight count and communicates only over
channels.

## Testing it

Do not crawl the real internet in tests. `httptest.NewServer` with a handler
that serves a small in-memory site (a `map[string]string` of path to HTML) gives
you a deterministic, offline, instant test — including 404s, redirects, slow
pages (`time.Sleep` in the handler) and a link cycle.

## Stretch

- Respect `Retry-After` and back off on 429/503.
- Write the result as a site graph and reuse `09-datastructures/graph` to find
  orphan pages and the shortest click-path to each page.
- Stream results as they arrive via `iter.Seq[Result]` instead of a slice.
