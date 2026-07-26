// Package tcpproto drops below HTTP: a line-oriented TCP protocol with a
// server, a client, and a graceful shutdown that actually works.
//
// The protocol, one command per line, terminated by "\n":
//
//	SET <key> <value>   -> "+OK"
//	GET <key>           -> "$<value>" or "-ERR not found"
//	DEL <key>           -> "+OK" or "-ERR not found"
//	KEYS                -> "*<n>" then n lines of "$<key>", keys sorted
//	PING                -> "+PONG"
//	QUIT                -> "+BYE" and the server closes the connection
//	anything else       -> "-ERR unknown command"
//
// Keys and values may not contain spaces or newlines; a malformed command is
// "-ERR bad syntax". Lines longer than MaxLine are "-ERR line too long"
// followed by closing the connection - an unbounded line is a memory attack.
package tcpproto

import (
	"context"
	"net"
	"time"
)

// MaxLine is the largest command line accepted.
const MaxLine = 4096

// Server serves the protocol on a listener.
type Server struct {
	// IdleTimeout closes connections that send nothing for this long.
	// Zero means no timeout.
	IdleTimeout time.Duration
	// TODO: your fields: the store and its mutex, the listener, a WaitGroup
	// for live connections, a shutdown signal, and a set of open connections
	// so Shutdown can force them closed.
}

func NewServer() *Server { panic("TODO: implement NewServer") }

// Serve accepts connections until the listener is closed, handling each in its
// own goroutine. It returns nil after a Shutdown, and the accept error
// otherwise. A failure on one connection must never take the server down.
func (s *Server) Serve(ln net.Listener) error { panic("TODO: implement Serve") }

// ListenAndServe listens on addr (use "127.0.0.1:0" in tests to get a free
// port) and serves. Addr reports the actual address afterwards.
func (s *Server) ListenAndServe(addr string) error { panic("TODO: implement ListenAndServe") }

// Addr is the address the server is listening on, or "" if it is not.
func (s *Server) Addr() string { panic("TODO: implement Addr") }

// Shutdown stops accepting new connections and waits for the in-flight ones to
// finish, up to the context's deadline. On timeout it closes them by force and
// returns ctx.Err(). It must be safe to call twice.
func (s *Server) Shutdown(ctx context.Context) error { panic("TODO: implement Shutdown") }

// Conns reports how many connections are currently open.
func (s *Server) Conns() int { panic("TODO: implement Conns") }

// Client is a connection to a Server. It is NOT safe for concurrent use by
// several goroutines - document that, and make the tests hold you to it by
// giving each goroutine its own client.
type Client struct {
	// TODO
}

// Dial connects to addr.
func Dial(addr string) (*Client, error) { panic("TODO: implement Dial") }

// DialContext connects with a timeout.
func DialContext(ctx context.Context, addr string) (*Client, error) {
	panic("TODO: implement DialContext")
}

func (c *Client) Set(key, value string) error    { panic("TODO: implement Set") }
func (c *Client) Get(key string) (string, error) { panic("TODO: implement Get") }
func (c *Client) Del(key string) error           { panic("TODO: implement Del") }
func (c *Client) Keys() ([]string, error)        { panic("TODO: implement Keys") }
func (c *Client) Ping() error                    { panic("TODO: implement Ping") }

// Close sends QUIT and closes the connection.
func (c *Client) Close() error { panic("TODO: implement Close") }

// ErrNotFound is what Get and Del return for a missing key. It must be
// distinguishable from a protocol or network error.
var ErrNotFound = errNotFound{}

type errNotFound struct{}

func (errNotFound) Error() string { return "not found" }
