package tcpproto

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"net"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

// startServer runs a server on a random port and shuts it down when the test ends.
func startServer(t *testing.T) (*Server, string) {
	t.Helper()
	s := NewServer()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	done := make(chan error, 1)
	go func() { done <- s.Serve(ln) }()
	t.Cleanup(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		s.Shutdown(ctx)
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("Serve returned %v, want nil after Shutdown", err)
			}
		case <-time.After(2 * time.Second):
			t.Error("Serve did not return after Shutdown")
		}
	})
	return s, ln.Addr().String()
}

// rawConn talks the protocol by hand, so the tests do not depend on Client
// being correct.
type rawConn struct {
	c net.Conn
	r *bufio.Reader
	t *testing.T
}

func dialRaw(t *testing.T, addr string) *rawConn {
	t.Helper()
	c, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { c.Close() })
	return &rawConn{c: c, r: bufio.NewReader(c), t: t}
}

func (rc *rawConn) send(line string) string {
	rc.t.Helper()
	rc.c.SetDeadline(time.Now().Add(2 * time.Second))
	if _, err := fmt.Fprintf(rc.c, "%s\n", line); err != nil {
		rc.t.Fatalf("write %q: %v", line, err)
	}
	resp, err := rc.r.ReadString('\n')
	if err != nil {
		rc.t.Fatalf("read after %q: %v", line, err)
	}
	return strings.TrimRight(resp, "\r\n")
}

func (rc *rawConn) readLine() string {
	rc.t.Helper()
	resp, err := rc.r.ReadString('\n')
	if err != nil {
		rc.t.Fatalf("read: %v", err)
	}
	return strings.TrimRight(resp, "\r\n")
}

func TestProtocol(t *testing.T) {
	_, addr := startServer(t)
	rc := dialRaw(t, addr)

	tests := []struct{ send, want string }{
		{"PING", "+PONG"},
		{"SET a 1", "+OK"},
		{"GET a", "$1"},
		{"GET nope", "-ERR not found"},
		{"SET a 2", "+OK"},
		{"GET a", "$2"},
		{"DEL a", "+OK"},
		{"DEL a", "-ERR not found"},
		{"GET a", "-ERR not found"},
		{"FROBNICATE", "-ERR unknown command"},
		{"SET", "-ERR bad syntax"},
		{"SET onlykey", "-ERR bad syntax"},
		{"GET", "-ERR bad syntax"},
		{"", "-ERR bad syntax"},
	}
	for _, tt := range tests {
		if got := rc.send(tt.send); got != tt.want {
			t.Errorf("%q -> %q, want %q", tt.send, got, tt.want)
		}
	}
}

func TestKeys(t *testing.T) {
	_, addr := startServer(t)
	rc := dialRaw(t, addr)
	if got := rc.send("KEYS"); got != "*0" {
		t.Errorf("empty KEYS = %q, want *0", got)
	}
	rc.send("SET zebra 1")
	rc.send("SET alpha 2")
	rc.send("SET mango 3")
	if got := rc.send("KEYS"); got != "*3" {
		t.Fatalf("KEYS = %q, want *3", got)
	}
	var keys []string
	for range 3 {
		keys = append(keys, rc.readLine())
	}
	want := []string{"$alpha", "$mango", "$zebra"}
	if !reflect.DeepEqual(keys, want) {
		t.Errorf("keys = %q, want %q (sorted)", keys, want)
	}
}

func TestQuit(t *testing.T) {
	_, addr := startServer(t)
	rc := dialRaw(t, addr)
	if got := rc.send("QUIT"); got != "+BYE" {
		t.Errorf("QUIT = %q", got)
	}
	rc.c.SetDeadline(time.Now().Add(time.Second))
	if _, err := rc.r.ReadString('\n'); err == nil {
		t.Error("the server must close the connection after QUIT")
	}
}

func TestLongLineRejected(t *testing.T) {
	_, addr := startServer(t)
	rc := dialRaw(t, addr)
	rc.c.SetDeadline(time.Now().Add(2 * time.Second))
	fmt.Fprintf(rc.c, "SET k %s\n", strings.Repeat("x", MaxLine*2))
	resp, err := rc.r.ReadString('\n')
	if err != nil {
		t.Fatalf("expected an error response, got %v", err)
	}
	if !strings.Contains(resp, "too long") {
		t.Errorf("= %q, want -ERR line too long", strings.TrimSpace(resp))
	}
}

func TestStateIsSharedAcrossConnections(t *testing.T) {
	_, addr := startServer(t)
	a := dialRaw(t, addr)
	b := dialRaw(t, addr)
	a.send("SET shared yes")
	if got := b.send("GET shared"); got != "$yes" {
		t.Errorf("second connection sees %q", got)
	}
}

func TestConcurrentClients(t *testing.T) {
	s, addr := startServer(t)
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c, err := Dial(addr)
			if err != nil {
				t.Error(err)
				return
			}
			defer c.Close()
			key := fmt.Sprintf("k%d", i)
			for j := range 10 {
				if err := c.Set(key, fmt.Sprint(j)); err != nil {
					t.Error(err)
					return
				}
				v, err := c.Get(key)
				if err != nil {
					t.Error(err)
					return
				}
				if v != fmt.Sprint(j) {
					t.Errorf("got %q, want %d", v, j)
					return
				}
			}
		}()
	}
	wg.Wait()
	// Give the server a moment to notice the closed connections.
	for range 100 {
		if s.Conns() == 0 {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Errorf("Conns = %d after every client closed, want 0", s.Conns())
}

func TestClient(t *testing.T) {
	_, addr := startServer(t)
	c, err := Dial(addr)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	if err := c.Ping(); err != nil {
		t.Fatal(err)
	}
	if err := c.Set("name", "ana"); err != nil {
		t.Fatal(err)
	}
	v, err := c.Get("name")
	if err != nil || v != "ana" {
		t.Fatalf("Get = %q, %v", v, err)
	}
	if _, err := c.Get("missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get(missing) = %v, want ErrNotFound", err)
	}
	c.Set("b", "2")
	keys, err := c.Keys()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(keys, []string{"b", "name"}) {
		t.Errorf("Keys = %q", keys)
	}
	if err := c.Del("name"); err != nil {
		t.Fatal(err)
	}
	if err := c.Del("name"); !errors.Is(err, ErrNotFound) {
		t.Errorf("Del twice = %v, want ErrNotFound", err)
	}
}

func TestDialContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	// 203.0.113.0/24 is reserved for documentation and should not respond.
	_, err := DialContext(ctx, "203.0.113.1:9")
	if err == nil {
		t.Error("expected a dial failure")
	}
}

func TestGracefulShutdownWaitsForConnections(t *testing.T) {
	s := NewServer()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	served := make(chan error, 1)
	go func() { served <- s.Serve(ln) }()

	rc := dialRaw(t, ln.Addr().String())
	rc.send("SET k v")

	shutdownDone := make(chan error, 1)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		shutdownDone <- s.Shutdown(ctx)
	}()

	time.Sleep(50 * time.Millisecond)
	select {
	case <-shutdownDone:
		t.Fatal("Shutdown returned while a connection was still open")
	default:
	}

	// New connections must be refused once shutdown has started.
	if c, err := net.DialTimeout("tcp", ln.Addr().String(), 200*time.Millisecond); err == nil {
		c.Close()
		t.Error("the listener should be closed during shutdown")
	}

	// The live connection still works, then closes.
	if got := rc.send("GET k"); got != "$v" {
		t.Errorf("in-flight connection broke: %q", got)
	}
	rc.send("QUIT")

	select {
	case err := <-shutdownDone:
		if err != nil {
			t.Errorf("Shutdown = %v, want nil", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Shutdown did not return after the connection closed")
	}
	if err := <-served; err != nil {
		t.Errorf("Serve = %v, want nil", err)
	}
	if err := s.Shutdown(context.Background()); err != nil {
		t.Errorf("second Shutdown = %v, want nil", err)
	}
}

func TestShutdownTimeout(t *testing.T) {
	s := NewServer()
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	go s.Serve(ln)
	rc := dialRaw(t, ln.Addr().String())
	rc.send("PING") // hold the connection open and idle

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	start := time.Now()
	err := s.Shutdown(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Shutdown = %v, want DeadlineExceeded", err)
	}
	if d := time.Since(start); d > 2*time.Second {
		t.Errorf("Shutdown took %v", d)
	}
}

func TestIdleTimeout(t *testing.T) {
	s := NewServer()
	s.IdleTimeout = 100 * time.Millisecond
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	go s.Serve(ln)
	defer s.Shutdown(context.Background())

	rc := dialRaw(t, ln.Addr().String())
	rc.send("PING")
	rc.c.SetDeadline(time.Now().Add(2 * time.Second))
	if _, err := rc.r.ReadString('\n'); err == nil {
		t.Error("an idle connection should have been closed by the server")
	}
}
