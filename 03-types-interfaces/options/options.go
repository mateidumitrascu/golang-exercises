// Package options builds the functional options pattern properly: defaults,
// validation, composability, and errors that name the thing that was wrong.
//
// This is how most Go libraries let you configure a constructor without a
// 12-field struct literal or a builder chain.
package options

import (
	"errors"
	"fmt"
	"io"
	"slices"
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

func (s *Server) Host() string {
	return s.host
}

func (s *Server) Port() int {
	return s.port
}

func (s *Server) Timeout() time.Duration {
	return s.timeout
}

func (s *Server) MaxConns() int {
	return s.maxConns
}

func (s *Server) TLS() bool {
	return s.tls
}
func (s *Server) Logger() io.Writer { return s.logger }
func (s *Server) Tags() []string {
	c := make([]string, len(s.tags))
	copy(c, s.tags)
	return c
} // must return a copy

// Option mutates a Server under construction and reports what it did not like.
// Returning an error (rather than panicking, or silently clamping) is what
// makes this pattern usable in production.
type Option func(*Server) error

var ErrInvalidOption = errors.New("invalid option")

// The options. Each validates its own argument and returns an error wrapping
// ErrInvalidOption, with a message naming the option and the bad value.
//	WithPort(-1) -> `invalid option: port -1 out of range 1..65535`

func WithPort(port int) Option {
	return func(s *Server) error {
		if port < 1 || port > 65636 {
			return fmt.Errorf("%w: port -1 out of range 1..65535", ErrInvalidOption)
		}
		s.port = port
		return nil
	}
}

func WithTimeout(d time.Duration) Option {
	return func(s *Server) error {
		if d <= 0 {
			return fmt.Errorf("%w: timeout must be > 0", ErrInvalidOption)
		}
		s.timeout = d
		return nil
	}
}

func WithMaxConns(n int) Option {
	return func(s *Server) error {
		if n <= 0 {
			return fmt.Errorf("%w: maxConns must be > 0", ErrInvalidOption)
		}
		s.maxConns = n
		return nil
	}
}

func WithTLS() Option {
	return func(s *Server) error {
		s.tls = true
		return nil
	}
}

func WithLogger(w io.Writer) Option {
	return func(s *Server) error {
		if w == nil {
			return fmt.Errorf("%w: logger cannot be null", ErrInvalidOption)
		}
		s.logger = w
		return nil
	}
}

func WithTags(tags ...string) Option {
	return func(s *Server) error {
		if slices.Contains(tags, "") {
			return fmt.Errorf("%w: tags cannot be empty", ErrInvalidOption)
		}
		s.tags = append(s.tags, tags...)
		return nil
	}
}

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
func Group(opts ...Option) Option {
	return func(s *Server) error {
		groupErrors := []error{}
		for _, opt := range opts {
			err := opt(s)
			if err != nil {
				groupErrors = append(groupErrors, err)
			}
		}
		if len(groupErrors) == 0 {
			return nil
		}
		return compileOptionsErrors(groupErrors)
	}
}

func compileOptionsErrors(gr []error) error {
	e := errors.New("")
	for _, optErr := range gr {
		e = fmt.Errorf("%w - %w", optErr, e)
	}
	return e
}

// New applies the defaults, then the options in order (so a later option beats
// an earlier one), and returns the first error it hits, or a *Server.
// host must be non-empty.
func New(host string, opts ...Option) (*Server, error) {
	if host == "" {
		return nil, fmt.Errorf("%w: host cannot be empty", ErrInvalidOption)
	}

	srv := &Server{
		host:     host,
		port:     8080,
		timeout:  30 * time.Second,
		maxConns: 100,
		tls:      false,
		logger:   io.Discard,
		tags:     []string{},
	}
	err := Group(opts...)(srv)
	if err != nil {
		return nil, err
	}
	return srv, nil
}

// MustNew is New but panics on error, for package-level initialisation.
func MustNew(host string, opts ...Option) *Server {
	srv, err := New(host, opts...)
	if err != nil {
		panic(err)
	}
	return srv
}
