package apierr

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestKindString(t *testing.T) {
	want := []string{"other", "invalid", "not found", "conflict", "permission", "internal", "unavailable"}
	for i, w := range want {
		if got := Kind(i).String(); got != w {
			t.Errorf("Kind(%d) = %q, want %q", i, got, w)
		}
	}
}

func TestEAndError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			"full",
			E(Op("user.Create"), KindInvalid, Field("email"), "must not be empty"),
			"user.Create: email: invalid: must not be empty",
		},
		{"just a message", E("boom"), "boom"},
		{"op and message", E(Op("db.Get"), "timed out"), "db.Get: timed out"},
		{"kind only", E(KindNotFound), "not found"},
		{"other kind is hidden", E(KindOther, "x"), "x"},
		{
			"wrapped",
			E(Op("api.Signup"), E(Op("user.Create"), KindInvalid, Field("email"), "must not be empty")),
			"api.Signup: user.Create: email: invalid: must not be empty",
		},
		{"wrapping a plain error", E(Op("db.Get"), errors.New("connection refused")), "db.Get: connection refused"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEArgumentTypes(t *testing.T) {
	err := E(Op("a"), KindConflict, Field("f"), "msg", errors.New("inner"))
	var e *Error
	if !errors.As(err, &e) {
		t.Fatalf("E returned %T", err)
	}
	if e.Op != "a" || e.Kind != KindConflict || e.Field != "f" || e.Msg != "msg" || e.Err == nil {
		t.Errorf("E filled in %+v", e)
	}
	// Last of a type wins.
	err = E(Op("a"), Op("b"))
	e = err.(*Error)
	if e.Op != "b" {
		t.Errorf("Op = %q, want b", e.Op)
	}
	// Unknown argument types are a programming error.
	defer func() {
		if recover() == nil {
			t.Error("E with an unsupported argument type must panic")
		}
	}()
	E(42)
}

func TestUnwrapAndIs(t *testing.T) {
	root := errors.New("root")
	err := E(Op("outer"), E(Op("inner"), KindNotFound, root))

	if !errors.Is(err, root) {
		t.Error("errors.Is must reach the wrapped root error")
	}
	if !errors.Is(err, E(KindNotFound)) {
		t.Error("errors.Is(err, E(KindNotFound)) should match on Kind alone")
	}
	if errors.Is(err, E(KindConflict)) {
		t.Error("matched the wrong Kind")
	}
	if !errors.Is(err, E(Op("inner"))) {
		t.Error("should match on Op alone")
	}
	if errors.Is(err, E(Op("nope"))) {
		t.Error("matched the wrong Op")
	}
	if !errors.Is(err, E(Op("inner"), KindNotFound)) {
		t.Error("should match when every set field matches")
	}
	if errors.Is(err, E(Op("inner"), KindConflict)) {
		t.Error("must not match when one set field differs")
	}
}

func TestKindOf(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want Kind
	}{
		{"nil", nil, KindOther},
		{"plain", errors.New("x"), KindInternal},
		{"direct", E(KindNotFound), KindNotFound},
		{"nested", E(Op("a"), E(Op("b"), KindConflict)), KindConflict},
		{"outermost wins", E(KindInvalid, E(KindConflict)), KindInvalid},
		{"through fmt.Errorf", fmt.Errorf("wrap: %w", E(KindPermission)), KindPermission},
		{"no kind anywhere", E(Op("a"), E(Op("b"))), KindInternal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KindOf(tt.err); got != tt.want {
				t.Errorf("KindOf = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOps(t *testing.T) {
	err := E(Op("api.Signup"), E(Op("user.Create"), E(KindInvalid, "bad")))
	if got := Ops(err); !reflect.DeepEqual(got, []Op{"api.Signup", "user.Create"}) {
		t.Errorf("Ops = %v, want [api.Signup user.Create]", got)
	}
	if got := Ops(nil); len(got) != 0 {
		t.Errorf("Ops(nil) = %v", got)
	}
	if got := Ops(errors.New("x")); len(got) != 0 {
		t.Errorf("Ops(plain) = %v", got)
	}
}

func TestPublic(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"nil", nil, ""},
		{"internal", E(KindInternal, "database on fire at 10.0.0.4"), "internal error"},
		{"unclassified", errors.New("goroutine leak in shard 7"), "internal error"},
		{"validation", E(KindInvalid, Field("email"), "must not be empty"), "email: must not be empty"},
		{"no field", E(KindNotFound, "no such user"), "no such user"},
		{"no message", E(KindConflict), "conflict"},
		{"unavailable", E(KindUnavailable, "upstream down"), "internal error"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Public(tt.err); got != tt.want {
				t.Errorf("Public = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		err  error
		want int
	}{
		{nil, 200},
		{E(KindInvalid), 400},
		{E(KindPermission), 403},
		{E(KindNotFound), 404},
		{E(KindConflict), 409},
		{E(KindInternal), 500},
		{E(KindUnavailable), 503},
		{errors.New("x"), 500},
		{fmt.Errorf("w: %w", E(Op("a"), E(KindNotFound))), 404},
	}
	for _, tt := range tests {
		if got := HTTPStatus(tt.err); got != tt.want {
			t.Errorf("HTTPStatus(%v) = %d, want %d", tt.err, got, tt.want)
		}
	}
}
