# 08 — HTTP and networking

`net/http` is a full HTTP/1.1 and HTTP/2 stack in the standard library, and
since 1.22 its router is good enough that most services need no dependencies at
all. This module builds a service, the middleware around it, a well-behaved
client, and then drops to raw TCP.

| # | Exercise | Difficulty | What it drills |
|---|----------|-----------|----------------|
| 1 | `httpserver` | ★★☆ | `ServeMux` patterns, status codes, `httptest`, concurrent state |
| 2 | `middleware` | ★★★ | The decorator pattern, wrapping `ResponseWriter`, context values |
| 3 | `httpclient` | ★★★ | `RoundTripper`, retries, body replay, hostile servers |
| 4 | `tcpproto` | ★★★ | `net.Listener`, per-connection goroutines, graceful shutdown |

## Ideas worth internalising here

**`http.Handler` is one method, and everything composes through it.** Routers,
middleware, timeouts and mocks are all just handlers wrapping handlers. There is
no framework to learn.

**Your handler runs concurrently.** Any state it touches needs a mutex or an
atomic. This is the single most common bug in a first Go service.

**Wrapping `ResponseWriter` hides its optional interfaces.** The real one also
implements `Flusher` and `Hijacker`; your wrapper doesn't, so streaming quietly
breaks. Implement `Unwrap() http.ResponseWriter` and use
`http.NewResponseController` — that is the modern fix.

**Context keys must be an unexported type.** `ctx.Value("user")` with a string
key can collide with any other package in the process. `type ctxKey struct{}` is
the idiom.

**A client must drain and close every response body**, even one it is throwing
away, or the connection cannot be reused and you will leak file descriptors
under load. And it must never trust `Content-Length` — cap the read.

**Retries need replayable bodies.** `req.Body` is a one-shot reader;
`req.GetBody` is how the transport gets a fresh one. If it is nil and there is a
body, you cannot safely retry.

**Graceful shutdown is: stop accepting, then wait, then force.** Close the
listener, wait on a `WaitGroup` of live connections, and hard-close whatever is
left when the deadline passes.

## Stretch goals

- Add Server-Sent Events to `httpserver`: `GET /tasks/stream` that pushes an
  event whenever a task changes, using `http.NewResponseController().Flush()`.
- Give `httpclient` a caching transport that honours `Cache-Control: max-age`
  and `ETag`/`If-None-Match`.
- Make `tcpproto` speak the real RESP protocol (Redis) well enough that
  `redis-cli PING` works against it. It is a surprisingly small step.
