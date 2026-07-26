// Package httpclient is the other half of net/http: being a good client.
//
// The key extension point is http.RoundTripper - one method, Do-like, that sits
// underneath http.Client. Retries, authentication, instrumentation and caching
// are all transports wrapping other transports. The rules for writing one are
// strict, and worth knowing:
//
//   - Do not modify the request. Clone it if you must change headers.
//   - Do not touch resp.Body except to read/close it in your own logic.
//   - Return either a response or an error, never both.
//   - A body that you do not fully drain before discarding costs you the
//     connection: the transport can only reuse a connection whose response has
//     been read to EOF and closed.
package httpclient

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// HeaderTransport adds fixed headers to every request.
type HeaderTransport struct {
	Base    http.RoundTripper // nil means http.DefaultTransport
	Headers map[string]string
}

// RoundTrip must not mutate req - clone it, because the caller may reuse the
// request and because a RoundTripper is called from many goroutines.
func (t *HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	panic("TODO: implement HeaderTransport.RoundTrip")
}

// RetryTransport retries failed requests.
//
// Retry when:
//   - the underlying transport returns an error that is not a context
//     cancellation, or
//   - the status is 429, 502, 503 or 504.
//
// Never retry when:
//   - the context is done,
//   - the request has a body that cannot be replayed (req.GetBody == nil and
//     req.Body != nil),
//   - attempts are exhausted.
//
// Between attempts, wait Backoff * 2^(attempt-1), unless the response carried a
// Retry-After header with a number of seconds, in which case honour that.
// Always drain and close the body of a response you are discarding.
type RetryTransport struct {
	Base        http.RoundTripper
	MaxAttempts int           // total attempts; <= 0 means 1
	Backoff     time.Duration // 0 means no delay
	Attempts    int           // set by RoundTrip: how many attempts the last call took
}

func (t *RetryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	panic("TODO: implement RetryTransport.RoundTrip")
}

// New builds a client with sane defaults: an overall timeout, retries, and any
// fixed headers. Note that http.Client.Timeout covers the whole exchange
// including reading the body - so a retrying transport lives INSIDE it.
func New(timeout time.Duration, attempts int, headers map[string]string) *http.Client {
	panic("TODO: implement New")
}

// APIError is returned by DoJSON for a non-2xx response.
type APIError struct {
	StatusCode int
	Status     string
	Body       string // at most the first 512 bytes, for the error message
}

func (e *APIError) Error() string { panic("TODO: implement APIError.Error") }

// DoJSON performs a JSON request/response round trip:
//
//	marshal in (skip the body entirely if in is nil)
//	set Content-Type and Accept
//	attach ctx
//	on 2xx, decode into out (skip if out is nil) and return nil
//	on anything else, return an *APIError
//	always close the body
//
// It must return an error, not panic, for a nil client (use http.DefaultClient).
func DoJSON(ctx context.Context, c *http.Client, method, url string, in, out any) error {
	panic("TODO: implement DoJSON")
}

// ErrTooManyRedirects is returned by the client built by NewNoRedirect.
var ErrTooManyRedirects = errors.New("too many redirects")

// NewNoRedirect returns a client that does not follow redirects at all: the
// 3xx response is returned to the caller as-is. (http.Client.CheckRedirect
// returning http.ErrUseLastResponse.)
func NewNoRedirect() *http.Client { panic("TODO: implement NewNoRedirect") }

// Fetch downloads a URL and returns the body, with a hard cap on how much it
// will read - because "the server is honest about Content-Length" is not a
// safe assumption.
func Fetch(ctx context.Context, c *http.Client, url string, maxBytes int64) ([]byte, error) {
	panic("TODO: implement Fetch")
}
