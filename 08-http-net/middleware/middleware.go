// Package middleware: the http.Handler decorator pattern.
//
//	type Middleware func(http.Handler) http.Handler
//
// That single type is the whole framework. Everything below is an exercise in
// wrapping a handler without breaking the things underneath it - which is
// harder than it looks, because wrapping http.ResponseWriter hides the optional
// interfaces (Flusher, Hijacker) that the real one implements. Since Go 1.20
// the answer is http.NewResponseController.
package middleware

import (
	"context"
	"io"
	"net/http"
	"time"
)

// Middleware wraps a handler.
type Middleware func(http.Handler) http.Handler

// Chain composes middlewares so that Chain(a, b, c)(h) runs a, then b, then c,
// then h - reading order, which is not the order you get from the naive fold.
func Chain(ms ...Middleware) Middleware { panic("TODO: implement Chain") }

// requestIDKey is the context key. It must be an unexported type so that no
// other package can collide with it - a plain string key is a bug.
type contextKey struct{ name string }

var requestIDKey = &contextKey{"request-id"}

// RequestID puts a unique ID on every request: it uses the incoming
// X-Request-ID header if present, otherwise generates one (crypto/rand hex is
// fine). The ID goes into the request context and into the X-Request-ID
// response header.
func RequestID(next http.Handler) http.Handler { panic("TODO: implement RequestID") }

// RequestIDFrom reads the ID back out, returning "" if there is none.
func RequestIDFrom(ctx context.Context) string { panic("TODO: implement RequestIDFrom") }

// StatusRecorder wraps a ResponseWriter to remember the status code and how
// many bytes were written. A handler that never calls WriteHeader has status
// 200 - your recorder must report that, not 0.
type StatusRecorder struct {
	http.ResponseWriter
	// TODO: your fields
}

func NewStatusRecorder(w http.ResponseWriter) *StatusRecorder {
	panic("TODO: implement NewStatusRecorder")
}
func (r *StatusRecorder) WriteHeader(code int)        { panic("TODO: implement WriteHeader") }
func (r *StatusRecorder) Write(b []byte) (int, error) { panic("TODO: implement Write") }
func (r *StatusRecorder) Status() int                 { panic("TODO: implement Status") }
func (r *StatusRecorder) Bytes() int                  { panic("TODO: implement Bytes") }

// Unwrap lets http.ResponseController find the original writer, so that a
// handler calling Flush still works through the wrapper. One method, and it
// fixes a whole category of subtle breakage.
func (r *StatusRecorder) Unwrap() http.ResponseWriter { panic("TODO: implement Unwrap") }

// Logger writes one line per request to out, after the handler returns:
//
//	GET /tasks 200 137b 1.2ms id=abc123
//
// Fields in that order, space separated: method, path, status, bytes with a "b"
// suffix, duration, and id=<request id> (omit the id field entirely if there
// is none). Use the real elapsed time; the tests only check the shape.
func Logger(out io.Writer) Middleware { panic("TODO: implement Logger") }

// Recoverer turns a panic in a handler into a 500 and logs it to out. Two
// details that matter:
//   - If the response has already been started, you cannot change the status;
//     just stop.
//   - http.ErrAbortHandler is a deliberate signal from the handler that the
//     connection should be dropped. Re-panic on it rather than swallowing it.
func Recoverer(out io.Writer) Middleware { panic("TODO: implement Recoverer") }

// Timeout gives the handler a context with a deadline. If the handler has not
// responded by then, reply 503 with the body "timeout". If it has already
// started writing, do nothing (you cannot un-send bytes).
//
// The handler runs in its own goroutine, so it must not touch the
// ResponseWriter after Timeout has given up - guard with a mutex.
func Timeout(d time.Duration) Middleware { panic("TODO: implement Timeout") }

// BasicAuth checks HTTP basic auth against users (username -> password). On
// failure it replies 401 with a WWW-Authenticate header naming realm.
// Compare passwords with crypto/subtle.ConstantTimeCompare, not ==.
func BasicAuth(realm string, users map[string]string) Middleware { panic("TODO: implement BasicAuth") }

// MaxBody rejects bodies larger than n bytes with 413, and caps r.Body so a
// handler that reads it cannot be made to allocate more than n.
func MaxBody(n int64) Middleware { panic("TODO: implement MaxBody") }

// CORS adds the headers for a permissive cross-origin policy and answers
// OPTIONS preflight requests with 204 without calling the next handler.
func CORS(allowedOrigin string) Middleware { panic("TODO: implement CORS") }
