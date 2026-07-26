# Project: persistent key-value store

**Modules it exercises:** 04 (errors), 06 (concurrency), 07 (I/O), 08 (HTTP).
**Rough size:** 600–900 lines. **Difficulty:** ★★★

A durable key-value store with a write-ahead log, crash recovery, TTLs and an
HTTP API. This is the project that teaches you what a database actually does.

## The spec

```
kvstore -dir ./data -addr :8080 -sync=true
```

```
PUT    /v1/kv/{key}          body = value; optional ?ttl=60s   -> 204
GET    /v1/kv/{key}                                            -> 200 + value, or 404
DELETE /v1/kv/{key}                                            -> 204 or 404
GET    /v1/kv?prefix=user:   -> 200, JSON array of keys
GET    /v1/stats             -> 200, JSON: keys, bytes, log size, uptime
POST   /v1/compact           -> 204, rewrites the log
```

## Requirements

1. **Write-ahead log.** Every mutation is appended to `data/wal.log` *before*
   the in-memory map is updated. Records are length-prefixed with a CRC32
   (reuse `07-io-encoding/framing`).
2. **Recovery.** On startup, replay the log to rebuild the map. A record that
   fails its checksum, or a torn final record from a crash mid-write, is
   discarded — and everything before it is still recovered. **This must not be
   a fatal error**: a torn tail is the normal outcome of a power cut.
3. **Durability flag.** With `-sync=true`, `fsync` after every write and be able
   to say (with a benchmark) how much it costs. With `-sync=false`, batch.
4. **TTL.** Expired keys read as missing, and are purged in the background
   without stopping the world.
5. **Compaction.** The log grows forever; `POST /v1/compact` (and an automatic
   trigger at, say, 2× the live data size) writes a fresh log with only live
   keys, then atomically renames it into place. A crash *during* compaction must
   leave the store recoverable.
6. **Concurrency.** Many readers, one writer: `sync.RWMutex` or a sharded map.
   Correct under `-race` with 100 concurrent clients.
7. **Graceful shutdown.** SIGINT: stop accepting, finish in-flight requests,
   flush and close the log, exit 0.

## Milestones

1. In-memory map behind an interface, plus the HTTP handlers. Test with
   `httptest`.
2. Append-only log + replay on startup. Kill the process, restart, check the
   data is there.
3. Corrupt the log deliberately in a test (truncate it mid-record, flip a byte)
   and prove recovery does the right thing.
4. TTL and background expiry.
5. Compaction with an atomic rename.
6. Benchmarks: writes/sec with and without fsync, memory per key.

## The interesting problem

Crash safety. Write a test that builds a store, writes N records, then
*truncates the file at every possible byte offset* and re-opens it. For every
offset, recovery must succeed and must contain a prefix of the records — never a
corrupted value, never a panic. That single test is worth the whole project.

## Stretch

- Add a `WATCH` endpoint using Server-Sent Events that streams changes.
- Add an on-disk index so the map doesn't have to hold every value in memory
  (keys and file offsets in memory, values read on demand — that is Bitcask).
- Add snapshot + WAL truncation, and measure recovery time for 1M keys.
