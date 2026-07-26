package wrapping

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func TestPathError(t *testing.T) {
	e := &PathError{Op: "open", Path: "/etc/shadow", Err: ErrPermission}
	want := "open /etc/shadow: permission denied"
	if got := e.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
	if !errors.Is(e, ErrPermission) {
		t.Error("errors.Is must see through PathError; did you implement Unwrap?")
	}
}

func TestOpen(t *testing.T) {
	if err := Open("fine.txt"); err != nil {
		t.Errorf("Open(fine.txt) = %v, want nil", err)
	}
	tests := []struct {
		path     string
		sentinel error
	}{
		{"missing", ErrNotFound},
		{"secret", ErrPermission},
		{"corrupt", ErrCorrupt},
	}
	for _, tt := range tests {
		err := Open(tt.path)
		if !errors.Is(err, tt.sentinel) {
			t.Errorf("Open(%q) = %v, want it to wrap %v", tt.path, err, tt.sentinel)
		}
		var pe *PathError
		if !errors.As(err, &pe) {
			t.Errorf("Open(%q) = %v, want a *PathError in the chain", tt.path, err)
			continue
		}
		if pe.Path != tt.path {
			t.Errorf("PathError.Path = %q, want %q", pe.Path, tt.path)
		}
	}
	if err := Open("corrupt"); !strings.Contains(err.Error(), "checksum") {
		t.Errorf("Open(corrupt) = %q, want the extra wrapping layer", err)
	}
}

func TestHTTPErrorCustomIs(t *testing.T) {
	e := HTTPError{Code: 404, Msg: "not found"}
	if got := e.Error(); got != "404: not found" {
		t.Errorf("Error() = %q, want %q", got, "404: not found")
	}
	if !errors.Is(e, HTTPError{Code: 404}) {
		t.Error("errors.Is should match on Code alone")
	}
	if errors.Is(e, HTTPError{Code: 500}) {
		t.Error("errors.Is matched a different code")
	}
	wrapped := fmt.Errorf("fetching user: %w", e)
	if !errors.Is(wrapped, HTTPError{Code: 404}) {
		t.Error("custom Is must still work through wrapping")
	}
	if errors.Is(wrapped, ErrNotFound) {
		t.Error("HTTPError.Is is matching things it should not")
	}
}

func TestChain(t *testing.T) {
	if got := Chain(); got != nil {
		t.Errorf("Chain() = %v, want nil", got)
	}
	if got := Chain(nil, nil); got != nil {
		t.Errorf("Chain(nil, nil) = %v, want nil", got)
	}
	single := errors.New("one")
	if got := Chain(nil, single, nil); got != single {
		t.Errorf("Chain with one error = %v, want the error itself", got)
	}
	a, b := errors.New("a"), errors.New("b")
	both := Chain(a, b)
	if !errors.Is(both, a) || !errors.Is(both, b) {
		t.Errorf("Chain(a, b) = %v, want both findable", both)
	}
}

func TestRootCause(t *testing.T) {
	if got := RootCause(nil); got != nil {
		t.Errorf("RootCause(nil) = %v", got)
	}
	deep := fmt.Errorf("a: %w", fmt.Errorf("b: %w", ErrNotFound))
	if got := RootCause(deep); got != ErrNotFound {
		t.Errorf("RootCause = %v, want ErrNotFound", got)
	}
	joined := errors.Join(ErrNotFound, ErrCorrupt)
	if got := RootCause(fmt.Errorf("x: %w", joined)); got != joined {
		t.Errorf("RootCause stopped at %v, want the join node itself", got)
	}
}

func TestFlattenAndCount(t *testing.T) {
	if len(Flatten(nil)) != 0 {
		t.Error("Flatten(nil) must be empty")
	}
	leaf1, leaf2 := errors.New("l1"), errors.New("l2")
	tree := fmt.Errorf("top: %w", errors.Join(leaf1, fmt.Errorf("mid: %w", leaf2)))
	all := Flatten(tree)
	if len(all) != 5 {
		t.Errorf("Flatten found %d nodes (%v), want 5: top, join, l1, mid, l2", len(all), all)
	}
	if all[0] != tree {
		t.Error("Flatten must start with the error itself")
	}
	if got := CountLeaves(tree); got != 2 {
		t.Errorf("CountLeaves = %d, want 2", got)
	}
	if got := CountLeaves(leaf1); got != 1 {
		t.Errorf("CountLeaves(leaf) = %d, want 1", got)
	}
	if got := CountLeaves(nil); got != 0 {
		t.Errorf("CountLeaves(nil) = %d, want 0", got)
	}
}

func TestFirstPathError(t *testing.T) {
	inner := &PathError{Op: "read", Path: "b", Err: ErrCorrupt}
	outer := &PathError{Op: "open", Path: "a", Err: inner}
	got, ok := FirstPathError(fmt.Errorf("x: %w", outer))
	if !ok || got != outer {
		t.Errorf("FirstPathError = %v, %v; want the outermost one", got, ok)
	}
	if _, ok := FirstPathError(ErrNotFound); ok {
		t.Error("FirstPathError found one where there is none")
	}
}

func TestSummarise(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{nil, "ok"},
		{Open("missing"), "missing"},
		{Open("secret"), "denied"},
		{Open("corrupt"), "corrupt"},
		{errors.New("who knows"), "unknown"},
		{errors.Join(errors.New("x"), ErrPermission), "denied"},
		{errors.Join(ErrNotFound, ErrPermission), "missing"},
	}
	for _, tt := range tests {
		if got := Summarise(tt.err); got != tt.want {
			t.Errorf("Summarise(%v) = %q, want %q", tt.err, got, tt.want)
		}
	}
}
