package recovery

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

func TestSafelyNoPanic(t *testing.T) {
	called := false
	if err := Safely(func() { called = true }); err != nil {
		t.Errorf("Safely = %v, want nil", err)
	}
	if !called {
		t.Error("Safely did not call fn")
	}
}

func TestSafelyCatchesPanic(t *testing.T) {
	err := Safely(func() { panic("boom") })
	if err == nil {
		t.Fatal("Safely returned nil for a panicking function")
	}
	var pe *PanicError
	if !errors.As(err, &pe) {
		t.Fatalf("err = %#v, want *PanicError", err)
	}
	if pe.Value != "boom" {
		t.Errorf("Value = %v, want boom", pe.Value)
	}
	if !strings.Contains(err.Error(), "boom") {
		t.Errorf("Error() = %q, should contain the panic value", err)
	}
	if len(pe.Stack) == 0 {
		t.Error("Stack is empty; capture it inside the deferred function")
	}
	if !strings.Contains(string(pe.Stack), "recovery") {
		t.Errorf("Stack does not look like a stack trace:\n%s", pe.Stack)
	}
}

func TestSafelyPanicWithError(t *testing.T) {
	sentinel := errors.New("sentinel")
	err := Safely(func() { panic(sentinel) })
	if !errors.Is(err, sentinel) {
		t.Errorf("errors.Is could not find the panicked error: %v", err)
	}
}

func TestSafelyCatchesRuntimeErrors(t *testing.T) {
	err := Safely(func() {
		var m map[string]int
		m["boom"] = 1 // assignment to entry in nil map
	})
	if err == nil {
		t.Fatal("Safely missed a runtime panic")
	}
	var re error
	var pe *PanicError
	if errors.As(err, &pe) {
		if v, ok := pe.Value.(error); ok {
			re = v
		}
	}
	if re == nil {
		t.Errorf("expected the runtime error value, got %#v", err)
	}

	err = Safely(func() {
		s := []int{1, 2, 3}
		i := 5
		_ = s[i]
	})
	if err == nil {
		t.Error("Safely missed an index-out-of-range panic")
	}
}

func TestSafelyValue(t *testing.T) {
	v, err := SafelyValue(func() int { return 42 })
	if v != 42 || err != nil {
		t.Errorf("= %v, %v; want 42, nil", v, err)
	}
	s, err := SafelyValue(func() string { panic("nope") })
	if s != "" || err == nil {
		t.Errorf("= %q, %v; want zero value and an error", s, err)
	}
}

func TestWrap(t *testing.T) {
	sentinel := errors.New("normal failure")
	if err := Wrap(func() error { return sentinel }); !errors.Is(err, sentinel) {
		t.Errorf("Wrap = %v, want the returned error", err)
	}
	if err := Wrap(func() error { return nil }); err != nil {
		t.Errorf("Wrap = %v, want nil", err)
	}
	err := Wrap(func() error { panic(42) })
	var pe *PanicError
	if !errors.As(err, &pe) || pe.Value != 42 {
		t.Errorf("Wrap = %#v, want a *PanicError holding 42", err)
	}
}

func TestGo(t *testing.T) {
	errc := Go(func() error { return nil })
	if err, ok := <-errc; err != nil || !ok {
		t.Errorf("received %v, %v", err, ok)
	}
	if _, ok := <-errc; ok {
		t.Error("the channel must be closed after the result")
	}

	errc = Go(func() error { panic("in a goroutine") })
	err := <-errc
	var pe *PanicError
	if !errors.As(err, &pe) {
		t.Fatalf("err = %#v, want *PanicError (a panic in another goroutine cannot be "+
			"recovered from here - fn needs its own defer)", err)
	}

	sentinel := errors.New("x")
	if err := <-Go(func() error { return sentinel }); !errors.Is(err, sentinel) {
		t.Errorf("err = %v, want the sentinel", err)
	}
}

func TestCleanupOrder(t *testing.T) {
	got := CleanupOrder()
	want := []string{"c", "b", "a"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CleanupOrder = %v, want %v (defers run LIFO, and they still run "+
			"while panicking)", got, want)
	}
}

func TestMustPositive(t *testing.T) {
	if got := MustPositive(3); got != 3 {
		t.Errorf("= %d", got)
	}
	err := Safely(func() { MustPositive(0) })
	if err == nil {
		t.Error("MustPositive(0) must panic")
	}
	if err := Safely(func() { MustPositive(-1) }); err == nil {
		t.Error("MustPositive(-1) must panic")
	}
}

func TestRepanic(t *testing.T) {
	sentinel := errors.New("known")
	if err := Repanic(func() { panic(sentinel) }); !errors.Is(err, sentinel) {
		t.Errorf("Repanic = %v, want the error", err)
	}
	if err := Repanic(func() {}); err != nil {
		t.Errorf("Repanic = %v, want nil", err)
	}
	// A non-error panic must escape - so catching it needs an outer recover.
	outer := Safely(func() {
		_ = Repanic(func() { panic("not an error") })
	})
	var pe *PanicError
	if !errors.As(outer, &pe) || pe.Value != "not an error" {
		t.Errorf("outer = %#v, want the re-panicked value to reach the outer recover", outer)
	}
}
