package httpserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
)

func newTestServer(t *testing.T) (http.Handler, *Store) {
	t.Helper()
	s := NewStore()
	return NewServer(s), s
}

func do(t *testing.T, h http.Handler, method, target string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == nil {
		r = httptest.NewRequest(method, target, nil)
	} else {
		var buf bytes.Buffer
		switch b := body.(type) {
		case string:
			buf.WriteString(b)
		default:
			json.NewEncoder(&buf).Encode(b)
		}
		r = httptest.NewRequest(method, target, &buf)
		r.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestStore(t *testing.T) {
	s := NewStore()
	task, err := s.Add(Task{Title: "write go"})
	if err != nil {
		t.Fatal(err)
	}
	if task.ID != 1 {
		t.Errorf("first ID = %d, want 1", task.ID)
	}
	if task.Created.IsZero() {
		t.Error("Created was not set")
	}
	second, _ := s.Add(Task{Title: "write more go"})
	if second.ID != 2 {
		t.Errorf("second ID = %d, want 2", second.ID)
	}

	if _, err := s.Add(Task{Title: ""}); !errors.Is(err, ErrInvalid) {
		t.Errorf("empty title = %v, want ErrInvalid", err)
	}
	if _, err := s.Add(Task{Title: strings.Repeat("x", 201)}); !errors.Is(err, ErrInvalid) {
		t.Errorf("long title = %v, want ErrInvalid", err)
	}

	if _, err := s.Get(99); !errors.Is(err, ErrNotFound) {
		t.Errorf("Get(99) = %v, want ErrNotFound", err)
	}
	if err := s.Delete(1); err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(1); !errors.Is(err, ErrNotFound) {
		t.Errorf("double delete = %v", err)
	}
}

func TestStoreReturnsCopies(t *testing.T) {
	s := NewStore()
	added, _ := s.Add(Task{Title: "original"})
	added.Title = "mutated"
	got, _ := s.Get(1)
	if got.Title != "original" {
		t.Error("Add returned a value that aliases the stored task")
	}
}

func TestStoreConcurrent(t *testing.T) {
	s := NewStore()
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Add(Task{Title: fmt.Sprintf("task %d", i)})
			s.List(nil, 0)
		}()
	}
	wg.Wait()
	if got := len(s.List(nil, 0)); got != 100 {
		t.Errorf("stored %d tasks, want 100", got)
	}
	ids := map[int]bool{}
	for _, task := range s.List(nil, 0) {
		if ids[task.ID] {
			t.Fatalf("duplicate ID %d", task.ID)
		}
		ids[task.ID] = true
	}
}

func TestListFiltering(t *testing.T) {
	s := NewStore()
	for i := range 5 {
		task, _ := s.Add(Task{Title: fmt.Sprint(i)})
		if i%2 == 0 {
			s.Update(task.ID, Task{Title: task.Title, Done: true})
		}
	}
	yes, no := true, false
	if got := len(s.List(&yes, 0)); got != 3 {
		t.Errorf("done=true gave %d, want 3", got)
	}
	if got := len(s.List(&no, 0)); got != 2 {
		t.Errorf("done=false gave %d, want 2", got)
	}
	if got := len(s.List(nil, 2)); got != 2 {
		t.Errorf("limit=2 gave %d", got)
	}
	list := s.List(nil, 0)
	for i := 1; i < len(list); i++ {
		if list[i].ID < list[i-1].ID {
			t.Fatal("List must be ordered by ID")
		}
	}
}

func TestHealthz(t *testing.T) {
	h, _ := newTestServer(t)
	w := do(t, h, "GET", "/healthz", nil)
	if w.Code != 200 || strings.TrimSpace(w.Body.String()) != "ok" {
		t.Errorf("= %d %q", w.Code, w.Body.String())
	}
}

func TestCreateTask(t *testing.T) {
	h, _ := newTestServer(t)
	w := do(t, h, "POST", "/tasks", map[string]any{"title": "buy milk"})
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201; body = %s", w.Code, w.Body)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Errorf("Content-Type = %q", ct)
	}
	if loc := w.Header().Get("Location"); loc != "/tasks/1" {
		t.Errorf("Location = %q, want /tasks/1", loc)
	}
	var got Task
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("body is not a task: %s", w.Body)
	}
	if got.ID != 1 || got.Title != "buy milk" {
		t.Errorf("= %+v", got)
	}
}

func TestCreateValidation(t *testing.T) {
	h, _ := newTestServer(t)
	tests := []struct {
		name string
		body any
		want int
	}{
		{"empty title", map[string]any{"title": ""}, 400},
		{"missing title", map[string]any{}, 400},
		{"malformed json", "{not json", 400},
		{"long title", map[string]any{"title": strings.Repeat("x", 300)}, 400},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := do(t, h, "POST", "/tasks", tt.body)
			if w.Code != tt.want {
				t.Errorf("status = %d, want %d (body: %s)", w.Code, tt.want, w.Body)
			}
			var e struct{ Error string }
			if err := json.Unmarshal(w.Body.Bytes(), &e); err != nil || e.Error == "" {
				t.Errorf("error body = %s, want {\"error\": \"...\"}", w.Body)
			}
		})
	}
}

func TestCreateRequiresJSONContentType(t *testing.T) {
	h, _ := newTestServer(t)
	r := httptest.NewRequest("POST", "/tasks", strings.NewReader(`{"title":"x"}`))
	r.Header.Set("Content-Type", "text/plain")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusUnsupportedMediaType {
		t.Errorf("status = %d, want 415", w.Code)
	}
}

func TestGetUpdateDelete(t *testing.T) {
	h, _ := newTestServer(t)
	do(t, h, "POST", "/tasks", map[string]any{"title": "one"})

	w := do(t, h, "GET", "/tasks/1", nil)
	if w.Code != 200 {
		t.Fatalf("get = %d", w.Code)
	}

	w = do(t, h, "PUT", "/tasks/1", map[string]any{"title": "one edited", "done": true})
	if w.Code != 200 {
		t.Fatalf("put = %d, body %s", w.Code, w.Body)
	}
	var task Task
	json.Unmarshal(w.Body.Bytes(), &task)
	if task.Title != "one edited" || !task.Done || task.ID != 1 {
		t.Errorf("after update: %+v", task)
	}

	w = do(t, h, "DELETE", "/tasks/1", nil)
	if w.Code != http.StatusNoContent {
		t.Errorf("delete = %d, want 204", w.Code)
	}
	if w.Body.Len() != 0 {
		t.Errorf("204 must have an empty body, got %q", w.Body)
	}
	if w := do(t, h, "GET", "/tasks/1", nil); w.Code != 404 {
		t.Errorf("after delete: %d, want 404", w.Code)
	}
}

func TestIDErrors(t *testing.T) {
	h, _ := newTestServer(t)
	if w := do(t, h, "GET", "/tasks/999", nil); w.Code != 404 {
		t.Errorf("unknown id = %d, want 404", w.Code)
	}
	if w := do(t, h, "GET", "/tasks/abc", nil); w.Code != 400 {
		t.Errorf("non-numeric id = %d, want 400", w.Code)
	}
	if w := do(t, h, "DELETE", "/tasks/999", nil); w.Code != 404 {
		t.Errorf("delete unknown = %d, want 404", w.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	h, _ := newTestServer(t)
	w := do(t, h, "PATCH", "/tasks/1", nil)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405 (the 1.22 mux does this for you)", w.Code)
	}
	if allow := w.Header().Get("Allow"); allow == "" {
		t.Error("405 should carry an Allow header")
	}
}

func TestListQueryParams(t *testing.T) {
	h, _ := newTestServer(t)
	for i := range 5 {
		do(t, h, "POST", "/tasks", map[string]any{"title": fmt.Sprint(i)})
	}
	do(t, h, "PUT", "/tasks/1", map[string]any{"title": "0", "done": true})

	var list []Task
	w := do(t, h, "GET", "/tasks?done=true", nil)
	json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 1 {
		t.Errorf("?done=true gave %d tasks", len(list))
	}
	w = do(t, h, "GET", "/tasks?limit=2", nil)
	json.Unmarshal(w.Body.Bytes(), &list)
	if len(list) != 2 {
		t.Errorf("?limit=2 gave %d tasks", len(list))
	}
	for _, bad := range []string{"?done=maybe", "?limit=-1", "?limit=abc"} {
		if w := do(t, h, "GET", "/tasks"+bad, nil); w.Code != 400 {
			t.Errorf("%s = %d, want 400", bad, w.Code)
		}
	}
}

func TestEmptyListIsArray(t *testing.T) {
	h, _ := newTestServer(t)
	w := do(t, h, "GET", "/tasks", nil)
	if strings.TrimSpace(w.Body.String()) != "[]" {
		t.Errorf("empty list = %s, want [] (a nil slice marshals to null - fix that)", w.Body)
	}
}

func TestBodyLimit(t *testing.T) {
	h, _ := newTestServer(t)
	huge := map[string]any{"title": strings.Repeat("x", 2<<20)}
	w := do(t, h, "POST", "/tasks", huge)
	if w.Code != 400 && w.Code != http.StatusRequestEntityTooLarge {
		t.Errorf("oversized body = %d, want 400 or 413", w.Code)
	}
}
