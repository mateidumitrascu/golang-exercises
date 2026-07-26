// Package httpserver builds a small JSON API with nothing but net/http.
//
// Since Go 1.22 the standard ServeMux does method matching and path wildcards,
// so you do not need a router library for this:
//
//	mux.HandleFunc("GET /tasks/{id}", h)   // r.PathValue("id")
//	mux.HandleFunc("POST /tasks", h)
//
// A pattern that matches the path but not the method produces 405 with an
// Allow header automatically. "GET" also matches HEAD.
package httpserver

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

// Task is the resource.
type Task struct {
	ID      int       `json:"id"`
	Title   string    `json:"title"`
	Done    bool      `json:"done"`
	Created time.Time `json:"created"`
}

var (
	ErrNotFound = errors.New("task not found")
	ErrInvalid  = errors.New("invalid task")
)

// Store is an in-memory repository. It must be safe for concurrent use - an
// http.Handler is called from many goroutines at once, always.
type Store struct {
	mu sync.RWMutex
	// TODO: your fields. IDs start at 1 and increase.
}

func NewStore() *Store { panic("TODO: implement NewStore") }

// Add validates (a title is required, at most 200 runes), assigns an ID and a
// Created timestamp, and stores a copy. It returns ErrInvalid for a bad task.
func (s *Store) Add(t Task) (Task, error) { panic("TODO: implement Add") }

// Get returns a copy, or ErrNotFound.
func (s *Store) Get(id int) (Task, error) { panic("TODO: implement Get") }

// List returns tasks ordered by ID. done == nil means "any".
func (s *Store) List(done *bool, limit int) []Task { panic("TODO: implement List") }

// Update replaces Title and Done for an existing task, keeping its ID and
// Created. Same validation as Add.
func (s *Store) Update(id int, t Task) (Task, error) { panic("TODO: implement Update") }

// Delete returns ErrNotFound if the id is unknown.
func (s *Store) Delete(id int) error { panic("TODO: implement Delete") }

// NewServer wires up the routes and returns the handler:
//
//	GET    /tasks          list; optional ?done=true|false and ?limit=N
//	POST   /tasks          create; 201 with a Location header
//	GET    /tasks/{id}     fetch
//	PUT    /tasks/{id}     replace
//	DELETE /tasks/{id}     remove; 204 with an empty body
//	GET    /healthz        200 with the body "ok"
//
// Rules for every JSON endpoint:
//   - Content-Type: application/json on every response with a body.
//   - A request body that is not valid JSON is 400, not 500.
//   - A POST or PUT without Content-Type: application/json is 415.
//   - An unknown /tasks/{id} is 404. A non-numeric id is 400.
//   - Errors are JSON: {"error": "some message"}.
//   - A bad ?done= or ?limit= value is 400.
//   - Bodies are limited to 1 MiB (http.MaxBytesReader).
func NewServer(store *Store) http.Handler { panic("TODO: implement NewServer") }

// WriteJSON is the helper every handler should go through: set the header,
// write the status, encode. Order matters - WriteHeader after the first Write
// is ignored and logs a warning.
func WriteJSON(w http.ResponseWriter, status int, v any) { panic("TODO: implement WriteJSON") }

// WriteError writes {"error": msg} with the given status.
func WriteError(w http.ResponseWriter, status int, msg string) { panic("TODO: implement WriteError") }
