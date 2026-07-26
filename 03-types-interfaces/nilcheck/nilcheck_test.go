package nilcheck

import (
	"errors"
	"strings"
	"testing"
)

func TestTheTrapExists(t *testing.T) {
	err := BuggyValidate(1)
	if err == nil {
		t.Fatal("BuggyValidate is supposed to demonstrate the bug: with `var e *MyError` " +
			"and `return e`, the returned interface is non-nil even on success")
	}
	var target *MyError
	if !errors.As(err, &target) || target != nil {
		t.Errorf("expected an error interface holding a nil *MyError, got %#v", err)
	}
}

func TestValidateIsFixed(t *testing.T) {
	if err := Validate(1); err != nil {
		t.Errorf("Validate(1) = %#v, want a truly nil error", err)
	}
	err := Validate(-1)
	if err == nil {
		t.Fatal("Validate(-1) must fail")
	}
	var me *MyError
	if !errors.As(err, &me) || me == nil {
		t.Fatalf("Validate(-1) = %#v, want a non-nil *MyError", err)
	}
	if me.Error() == "" {
		t.Error("MyError.Error returned an empty string")
	}
}

func TestIsNil(t *testing.T) {
	var nilPtr *MyError
	var nilMap map[string]int
	var nilSlice []int
	var nilChan chan int
	var nilFunc func()
	var nilIface error

	tests := []struct {
		name string
		v    any
		want bool
	}{
		{"untyped nil", nil, true},
		{"nil pointer", nilPtr, true},
		{"nil map", nilMap, true},
		{"nil slice", nilSlice, true},
		{"nil chan", nilChan, true},
		{"nil func", nilFunc, true},
		{"nil error interface", nilIface, true},
		{"non-nil pointer", &MyError{}, false},
		{"empty map", map[string]int{}, false},
		{"empty slice", []int{}, false},
		{"zero int", 0, false},
		{"empty string", "", false},
		{"zero struct", MyError{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNil(tt.v); got != tt.want {
				t.Errorf("IsNil(%#v) = %v, want %v", tt.v, got, tt.want)
			}
		})
	}
}

func TestSafeEqual(t *testing.T) {
	tests := []struct {
		name            string
		a, b            any
		wantEq, wantCmp bool
	}{
		{"equal ints", 1, 1, true, true},
		{"different ints", 1, 2, false, true},
		{"different types", 1, int64(1), false, true},
		{"strings", "a", "a", true, true},
		{"both nil", nil, nil, true, true},
		{"nil and value", nil, 1, false, true},
		{"slices", []int{1}, []int{1}, false, false},
		{"maps", map[string]int{}, map[string]int{}, false, false},
		{"funcs", func() {}, func() {}, false, false},
		{"struct with slice field", struct{ S []int }{}, struct{ S []int }{}, false, false},
		{"comparable structs", struct{ N int }{1}, struct{ N int }{1}, true, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eq, cmp := SafeEqual(tt.a, tt.b)
			if eq != tt.wantEq || cmp != tt.wantCmp {
				t.Errorf("SafeEqual = (%v, %v), want (%v, %v)", eq, cmp, tt.wantEq, tt.wantCmp)
			}
		})
	}
}

func TestFirstNonNil(t *testing.T) {
	a, b := 1, 2
	if got := FirstNonNil(nil, &a, &b); got != &a {
		t.Errorf("FirstNonNil = %v, want &a", got)
	}
	if got := FirstNonNil[int](nil, nil); got != nil {
		t.Errorf("FirstNonNil(all nil) = %v, want nil", got)
	}
	if got := FirstNonNil[int](); got != nil {
		t.Error("FirstNonNil() = non-nil")
	}
}

func TestCoalesce(t *testing.T) {
	if got := Coalesce("", "", "x", "y"); got != "x" {
		t.Errorf("Coalesce = %q, want x", got)
	}
	if got := Coalesce(0, 0, 3); got != 3 {
		t.Errorf("Coalesce = %d, want 3", got)
	}
	if got := Coalesce[string](); got != "" {
		t.Errorf("Coalesce() = %q", got)
	}
	if got := Coalesce(0, 0); got != 0 {
		t.Errorf("Coalesce(zeros) = %d", got)
	}
}

func TestLookup(t *testing.T) {
	m := map[string]int{"a": 1}
	v, err := Lookup(m, "a")
	if v != 1 || err != nil {
		t.Errorf("Lookup(a) = %v, %v", v, err)
	}
	_, err = Lookup(m, "zz")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want it to wrap ErrNotFound", err)
	}
	if !strings.Contains(err.Error(), "zz") {
		t.Errorf("error message %q should mention the key", err.Error())
	}
	if _, err := Lookup(nil, "a"); !errors.Is(err, ErrNotFound) {
		t.Errorf("lookup in a nil map = %v, want ErrNotFound", err)
	}
}
