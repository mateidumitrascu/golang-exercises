// Package options builds the functional options pattern properly: defaults,
// validation, composability, and errors that name the thing that was wrong.
//
// This is how most Go libraries let you configure a constructor without a
// 12-field struct literal or a builder chain.
package options

import (
	"errors"
	"io"
	"time"
)

// Server is the thing being configured. Fields are unexported on purpose:
// the only way in is New plus options, so an invalid Server cannot exist.
type Server struct {
	host     string
	port     int
	timeout  time.Duration
	maxConns int
	tls      bool
	logger   io.Writer
	tags     []string
}

// Accessors, so tests (and users) can read the result.
func (s *Server) Host() string           { panic("TODO") }
func (s *Server) Port() int              { panic("TODO") }
func (s *Server) Timeout() time.Duration { panic("TODO") }
func (s *Server) MaxConns() int          { panic("TODO") }
func (s *Server) TLS() bool              { panic("TODO") }
func (s *Server) Logger() io.Writer      { panic("TODO") }
func (s *Server) Tags() []string         { panic("TODO") } // must return a copy

// Option mutates a Server under construction and reports what it did not like.
// Returning an error (rather than panicking, or silently clamping) is what
// makes this pattern usable in production.
type Option func(*Server) error

var ErrInvalidOption = errors.New("invalid option")

// The options. Each validates its own argument and returns an error wrapping
// ErrInvalidOption, with a message naming the option and the bad value.
//
//	WithPort(-1) -> `invalid option: port -1 out of range 1..65535`
func WithPort(port int) Option           { panic("TODO: implement WithPort") }
func WithTimeout(d time.Duration) Option { panic("TODO: implement WithTimeout") }
func WithMaxConns(n int) Option          { panic("TODO: implement WithMaxConns") }
func WithTLS() Option                    { panic("TODO: implement WithTLS") }
func WithLogger(w io.Writer) Option      { panic("TODO: implement WithLogger") }
func WithTags(tags ...string) Option     { panic("TODO: implement WithTags") }

// Rules:
//
//	port      must be 1..65535
//	timeout   must be > 0
//	maxConns  must be > 0
//	logger    must not be nil
//	tags      appended, not replaced, across multiple WithTags calls; no empty tags
//
// Defaults when an option is not given: port 8080, timeout 30s, maxConns 100,
// tls false, logger io.Discard, tags empty (non-nil).

// Group bundles several options into one. This is what lets users share
// "profiles" of configuration.
func Group(opts ...Option) Option { panic("TODO: implement Group") }

// New applies the defaults, then the options in order (so a later option beats
// an earlier one), and returns the first error it hits, or a *Server.
// host must be non-empty.
func New(host string, opts ...Option) (*Server, error) { panic("TODO: implement New") }

// MustNew is New but panics on error, for package-level initialisation.
func MustNew(host string, opts ...Option) *Server { panic("TODO: implement MustNew") }
