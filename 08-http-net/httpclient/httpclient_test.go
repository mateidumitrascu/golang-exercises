package httpclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestHeaderTransport(t *testing.T) {
	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Header.Clone()
	}))
	defer srv.Close()

	c := &http.Client{Transport: &HeaderTransport{
		Headers: map[string]string{"Authorization": "Bearer t0ken", "X-Client": "goex"},
	}}
	req, _ := http.NewRequest("GET", srv.URL, nil)
	req.Header.Set("X-Original", "keep me")

	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if got.Get("Authorization") != "Bearer t0ken" || got.Get("X-Client") != "goex" {
		t.Errorf("headers not added: %v", got)
	}
	if got.Get("X-Original") != "keep me" {
		t.Error("existing headers must be preserved")
	}
	if req.Header.Get("Authorization") != "" {
		t.Error("RoundTrip must not modify the caller's request - clone it")
	}
}

func TestRetryOnStatus(t *testing.T) {
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits.Add(1) < 3 {
			w.WriteHeader(503)
			io.WriteString(w, "not yet")
			return
		}
		io.WriteString(w, "finally")
	}))
	defer srv.Close()

	rt := &RetryTransport{MaxAttempts: 5, Backoff: time.Millisecond}
	c := &http.Client{Transport: rt}
	resp, err := c.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 || string(body) != "finally" {
		t.Errorf("= %d %q", resp.StatusCode, body)
	}
	if hits.Load() != 3 {
		t.Errorf("server saw %d requests, want 3", hits.Load())
	}
	if rt.Attempts != 3 {
		t.Errorf("Attempts = %d, want 3", rt.Attempts)
	}
}

func TestRetryGivesUp(t *testing.T) {
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits.Add(1)
		w.WriteHeader(503)
	}))
	defer srv.Close()

	c := &http.Client{Transport: &RetryTransport{MaxAttempts: 3, Backoff: time.Millisecond}}
	resp, err := c.Get(srv.URL)
	if err != nil {
		t.Fatalf("after exhausting retries the last response should be returned, got err %v", err)
	}
	resp.Body.Close()
	if resp.StatusCode != 503 {
		t.Errorf("= %d", resp.StatusCode)
	}
	if hits.Load() != 3 {
		t.Errorf("server saw %d requests, want 3", hits.Load())
	}
}

func TestNoRetryOnSuccessOr4xx(t *testing.T) {
	for _, code := range []int{200, 400, 404, 500} {
		var hits atomic.Int64
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			hits.Add(1)
			w.WriteHeader(code)
		}))
		c := &http.Client{Transport: &RetryTransport{MaxAttempts: 3, Backoff: time.Millisecond}}
		resp, err := c.Get(srv.URL)
		if err == nil {
			resp.Body.Close()
		}
		if hits.Load() != 1 {
			t.Errorf("status %d caused %d attempts, want 1", code, hits.Load())
		}
		srv.Close()
	}
}

func TestRetryHonoursRetryAfter(t *testing.T) {
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hits.Add(1) == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(429)
			return
		}
		w.WriteHeader(200)
	}))
	defer srv.Close()

	c := &http.Client{Transport: &RetryTransport{MaxAttempts: 2, Backoff: time.Millisecond}}
	start := time.Now()
	resp, err := c.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if d := time.Since(start); d < 900*time.Millisecond {
		t.Errorf("waited %v, want at least the 1s from Retry-After", d)
	}
}

func TestRetryReplaysBody(t *testing.T) {
	var bodies []string
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		bodies = append(bodies, string(b))
		if hits.Add(1) < 2 {
			w.WriteHeader(503)
			return
		}
	}))
	defer srv.Close()

	c := &http.Client{Transport: &RetryTransport{MaxAttempts: 3, Backoff: time.Millisecond}}
	// http.NewRequest with a *strings.Reader sets GetBody for you.
	req, _ := http.NewRequest("POST", srv.URL, strings.NewReader("payload"))
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if len(bodies) != 2 || bodies[0] != "payload" || bodies[1] != "payload" {
		t.Errorf("server saw bodies %q, want the payload twice", bodies)
	}
}

func TestRetryRespectsContext(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(503)
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	req, _ := http.NewRequestWithContext(ctx, "GET", srv.URL, nil)
	c := &http.Client{Transport: &RetryTransport{MaxAttempts: 100, Backoff: 30 * time.Millisecond}}
	start := time.Now()
	resp, err := c.Do(req)
	if err == nil {
		resp.Body.Close()
	}
	if d := time.Since(start); d > time.Second {
		t.Errorf("took %v; a cancelled context must stop the retry loop", d)
	}
}

func TestDoJSON(t *testing.T) {
	type req struct {
		Name string `json:"name"`
	}
	type resp struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Content-Type = %q", r.Header.Get("Content-Type"))
		}
		var in req
		json.NewDecoder(r.Body).Decode(&in)
		json.NewEncoder(w).Encode(resp{ID: 7, Name: in.Name})
	}))
	defer srv.Close()

	var out resp
	if err := DoJSON(context.Background(), nil, "POST", srv.URL, req{Name: "ana"}, &out); err != nil {
		t.Fatal(err)
	}
	if out.ID != 7 || out.Name != "ana" {
		t.Errorf("= %+v", out)
	}
}

func TestDoJSONErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		io.WriteString(w, `{"error":"nope"}`)
	}))
	defer srv.Close()

	err := DoJSON(context.Background(), nil, "GET", srv.URL, nil, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %v (%T), want *APIError", err, err)
	}
	if apiErr.StatusCode != 422 {
		t.Errorf("StatusCode = %d", apiErr.StatusCode)
	}
	if !strings.Contains(apiErr.Body, "nope") {
		t.Errorf("Body = %q", apiErr.Body)
	}
	if !strings.Contains(apiErr.Error(), "422") {
		t.Errorf("Error() = %q, should mention the status", apiErr.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := DoJSON(ctx, nil, "GET", srv.URL, nil, nil); !errors.Is(err, context.Canceled) {
		t.Errorf("cancelled context = %v", err)
	}
}

func TestNewClient(t *testing.T) {
	var hits atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Api-Key") != "abc" {
			t.Errorf("header missing: %v", r.Header)
		}
		if hits.Add(1) < 2 {
			w.WriteHeader(502)
			return
		}
		io.WriteString(w, "ok")
	}))
	defer srv.Close()

	c := New(5*time.Second, 3, map[string]string{"X-Api-Key": "abc"})
	if c.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v", c.Timeout)
	}
	resp, err := c.Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if string(b) != "ok" {
		t.Errorf("= %q", b)
	}
}

func TestNoRedirect(t *testing.T) {
	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "final")
	}))
	defer target.Close()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, target.URL, http.StatusFound)
	}))
	defer srv.Close()

	resp, err := NewNoRedirect().Get(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusFound {
		t.Errorf("status = %d, want 302", resp.StatusCode)
	}
	if loc := resp.Header.Get("Location"); loc != target.URL {
		t.Errorf("Location = %q", loc)
	}
}

func TestFetchLimits(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// A server with far more to say than the client asked for.
		for range 100 {
			w.Write([]byte(strings.Repeat("x", 1000)))
		}
	}))
	defer srv.Close()

	got, err := Fetch(context.Background(), nil, srv.URL, 1000)
	if len(got) > 1000 {
		t.Errorf("read %d bytes despite a 1000-byte cap", len(got))
	}
	if err == nil && len(got) != 1000 {
		t.Errorf("read %d bytes, want either 1000 bytes or an error", len(got))
	}

	small := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello")
	}))
	defer small.Close()
	got, err = Fetch(context.Background(), nil, small.URL, 1000)
	if err != nil || string(got) != "hello" {
		t.Errorf("= %q, %v", got, err)
	}
}
