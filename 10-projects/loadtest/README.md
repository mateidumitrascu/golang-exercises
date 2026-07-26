# Project: HTTP load generator

**Modules it exercises:** 06 (concurrency), 08 (HTTP client), 09 (heaps).
**Rough size:** 400–600 lines. **Difficulty:** ★★★

A `hey`/`wrk`-style tool: hammer an endpoint at a fixed rate, then report
latency percentiles honestly.

## The spec

```
loadtest -url http://localhost:8080/api -c 50 -rate 500 -d 30s -m POST -body @payload.json
```

```
Summary:
  requests    15000
  duration    30.0s
  throughput  499.8 req/s
  errors      12 (0.08%)

Latency:
  p50   4.2ms
  p90  11.7ms
  p95  18.3ms
  p99  84.1ms
  max 302.0ms

Status codes:
  200  14891
  429     97
  503     12

Errors:
  context deadline exceeded  12
```

## Requirements

1. **Open-model load.** Requests are issued on a schedule (`-rate`), not "as
   fast as the last one finished". If the server slows down, the request rate
   does *not* drop — that is the difference between measuring the server and
   measuring your client.
2. **Coordinated omission.** Latency is measured from when the request *should
   have been sent*, not from when it actually was. Getting this wrong is the
   single most common flaw in load-testing tools; write a comment explaining
   your choice.
3. **Bounded concurrency** (`-c`) with a clear behaviour when the limit is hit:
   either queue or drop, but say which and count it.
4. **Connection reuse.** Configure `http.Transport` properly
   (`MaxIdleConnsPerHost` — the default of 2 will silently cap your throughput)
   and drain every response body.
5. **Percentiles** without keeping every sample: a fixed-size reservoir sample,
   or an HDR-style bucketed histogram. Then prove your percentiles are within 1%
   of the exact answer on a synthetic distribution.
6. **Graceful stop.** Ctrl-C prints the report for what has run so far.

## Milestones

1. Fire N requests sequentially, time them, print the mean. Notice the mean is
   useless.
2. Concurrency with a worker pool; collect durations in a slice; sort for exact
   percentiles.
3. Switch to rate-based scheduling with a `time.Ticker`, and handle the case
   where a tick arrives while every worker is busy.
4. Replace the slice with a histogram; compare against the exact percentiles in
   a test.
5. Status/error breakdown, `-body @file`, custom headers, keep-alive control.

## The interesting problem

Testing a load tester. Point it at an `httptest.Server` whose handler sleeps for
a known, controlled distribution (say, 90% at 1 ms and 10% at 100 ms) and assert
that the reported p50 and p99 land where they should. Use `testing/synctest` if
you want that test to be instant and exact.

## Stretch

- Live output: a progress line updated in place with the current rate and p99.
- Multiple targets from a file, weighted.
- Export the histogram as JSON so two runs can be diffed.
