package middleware

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

var okHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
})

func TestChainOrder(t *testing.T) {
	var order []string
	mk := func(name string) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				order = append(order, name)
				next.ServeHTTP(w, r)
			})
		}
	}
	h := Chain(mk("a"), mk("b"), mk("c"))(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		order = append(order, "handler")
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	want := "a b c handler"
	if got := strings.Join(order, " "); got != want {
		t.Errorf("order = %q, want %q", got, want)
	}
	// Chain with nothing must be a no-op, not a nil panic.
	Chain()(okHandler).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
}

func TestRequestID(t *testing.T) {
	var seen string
	h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		seen = RequestIDFrom(r.Context())
	}))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if seen == "" {
		t.Error("no request ID in the context")
	}
	if w.Header().Get("X-Request-ID") != seen {
		t.Errorf("header = %q, context = %q; they must match", w.Header().Get("X-Request-ID"), seen)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Set("X-Request-ID", "given-by-caller")
	h.ServeHTTP(httptest.NewRecorder(), r)
	if seen != "given-by-caller" {
		t.Errorf("= %q, want the incoming header to be reused", seen)
	}

	if got := RequestIDFrom(context.Background()); got != "" {
		t.Errorf("RequestIDFrom(empty) = %q", got)
	}

	// IDs must differ between requests.
	ids := map[string]bool{}
	for range 100 {
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
		if ids[seen] {
			t.Fatal("duplicate request ID")
		}
		ids[seen] = true
	}
}

func TestStatusRecorder(t *testing.T) {
	w := httptest.NewRecorder()
	rec := NewStatusRecorder(w)
	io.WriteString(rec, "hello")
	if rec.Status() != 200 {
		t.Errorf("Status = %d, want 200 for a handler that never calls WriteHeader", rec.Status())
	}
	if rec.Bytes() != 5 {
		t.Errorf("Bytes = %d, want 5", rec.Bytes())
	}

	w2 := httptest.NewRecorder()
	rec2 := NewStatusRecorder(w2)
	rec2.WriteHeader(404)
	io.WriteString(rec2, "nope")
	if rec2.Status() != 404 || rec2.Bytes() != 4 {
		t.Errorf("= %d, %d", rec2.Status(), rec2.Bytes())
	}
	if w2.Code != 404 || w2.Body.String() != "nope" {
		t.Error("the recorder must pass everything through")
	}
}

// TestStatusRecorderKeepsFlusher checks that wrapping does not break streaming.
func TestStatusRecorderKeepsFlusher(t *testing.T) {
	h := Logger(io.Discard)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "chunk")
		if err := http.NewResponseController(w).Flush(); err != nil {
			t.Errorf("Flush through the wrapper failed: %v (implement Unwrap)", err)
		}
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
}

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	h := Chain(RequestID, Logger(&buf))(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		io.WriteString(w, "hello")
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/tasks?x=1", nil))
	line := buf.String()
	for _, want := range []string{"GET", "/tasks", "201", "5b", "id="} {
		if !strings.Contains(line, want) {
			t.Errorf("log line %q is missing %q", strings.TrimSpace(line), want)
		}
	}

	buf.Reset()
	Logger(&buf)(okHandler).ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if strings.Contains(buf.String(), "id=") {
		t.Errorf("with no request ID the id field must be omitted: %q", buf.String())
	}
}

func TestRecoverer(t *testing.T) {
	var buf bytes.Buffer
	h := Recoverer(&buf)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("handler exploded")
	}))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 500 {
		t.Errorf("status = %d, want 500", w.Code)
	}
	if !strings.Contains(buf.String(), "exploded") {
		t.Errorf("log = %q, should contain the panic value", buf.String())
	}
}

func TestRecovererRepanicsOnAbort(t *testing.T) {
	h := Recoverer(io.Discard)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(http.ErrAbortHandler)
	}))
	defer func() {
		if r := recover(); r != http.ErrAbortHandler {
			t.Errorf("recovered %v; http.ErrAbortHandler must be re-panicked", r)
		}
	}()
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
}

func TestTimeout(t *testing.T) {
	h := Timeout(20 * time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-r.Context().Done():
		case <-time.After(2 * time.Second):
			t.Error("the handler's context was not cancelled")
		}
	}))
	w := httptest.NewRecorder()
	start := time.Now()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if d := time.Since(start); d > time.Second {
		t.Errorf("took %v; Timeout must give up", d)
	}
	if w.Code != http.StatusServiceUnavailable {
		t.Errorf("status = %d, want 503", w.Code)
	}
	if !strings.Contains(w.Body.String(), "timeout") {
		t.Errorf("body = %q", w.Body)
	}
}

func TestTimeoutFastHandler(t *testing.T) {
	h := Timeout(time.Second)(okHandler)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 200 || w.Body.String() != "hello" {
		t.Errorf("= %d %q", w.Code, w.Body)
	}
}

func TestTimeoutNoDataRace(t *testing.T) {
	// The handler keeps writing after the timeout fires; with -race this
	// catches an unguarded ResponseWriter.
	release := make(chan struct{})
	h := Timeout(10 * time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
		for range 100 {
			w.Write([]byte("late"))
		}
		close(release)
	}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	<-release
}

func TestBasicAuth(t *testing.T) {
	h := BasicAuth("test", map[string]string{"ana": "secret"})(okHandler)

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if w.Code != 401 {
		t.Errorf("no credentials = %d, want 401", w.Code)
	}
	if a := w.Header().Get("WWW-Authenticate"); !strings.Contains(a, "test") {
		t.Errorf("WWW-Authenticate = %q, should name the realm", a)
	}

	r := httptest.NewRequest("GET", "/", nil)
	r.SetBasicAuth("ana", "secret")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 {
		t.Errorf("valid credentials = %d", w.Code)
	}

	for _, creds := range [][2]string{{"ana", "wrong"}, {"nobody", "secret"}, {"", ""}} {
		r := httptest.NewRequest("GET", "/", nil)
		r.SetBasicAuth(creds[0], creds[1])
		w := httptest.NewRecorder()
		h.ServeHTTP(w, r)
		if w.Code != 401 {
			t.Errorf("%v = %d, want 401", creds, w.Code)
		}
	}
}

func TestMaxBody(t *testing.T) {
	var read int
	h := MaxBody(10)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		read = len(b)
	}))
	r := httptest.NewRequest("POST", "/", strings.NewReader("0123456789"))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != 200 || read != 10 {
		t.Errorf("exactly at the limit: %d, read %d", w.Code, read)
	}

	r = httptest.NewRequest("POST", "/", strings.NewReader(strings.Repeat("x", 100)))
	r.ContentLength = 100
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("oversized = %d, want 413", w.Code)
	}
}

func TestCORS(t *testing.T) {
	called := false
	h := CORS("https://example.com")(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		called = true
	}))
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("OPTIONS", "/", nil))
	if w.Code != http.StatusNoContent {
		t.Errorf("preflight = %d, want 204", w.Code)
	}
	if called {
		t.Error("a preflight request must not reach the handler")
	}
	if w.Header().Get("Access-Control-Allow-Origin") != "https://example.com" {
		t.Errorf("headers = %v", w.Header())
	}

	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/", nil))
	if !called || w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("normal requests must pass through with the header set")
	}
}

func TestFullStack(t *testing.T) {
	var logs bytes.Buffer
	var mu sync.Mutex
	safeLog := writerFunc(func(p []byte) (int, error) {
		mu.Lock()
		defer mu.Unlock()
		return logs.Write(p)
	})
	h := Chain(
		RequestID,
		Logger(safeLog),
		Recoverer(safeLog),
		Timeout(time.Second),
		MaxBody(1<<20),
	)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/panic" {
			panic("boom")
		}
		fmt.Fprintf(w, "id=%s", RequestIDFrom(r.Context()))
	}))

	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/ok", nil))
	if w.Code != 200 || !strings.HasPrefix(w.Body.String(), "id=") {
		t.Errorf("= %d %q", w.Code, w.Body)
	}
	w = httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", "/panic", nil))
	if w.Code != 500 {
		t.Errorf("panic through the stack = %d, want 500", w.Code)
	}
	mu.Lock()
	defer mu.Unlock()
	if strings.Count(logs.String(), "\n") < 2 {
		t.Errorf("expected a log line per request, got:\n%s", logs.String())
	}
}

type writerFunc func([]byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) { return f(p) }
