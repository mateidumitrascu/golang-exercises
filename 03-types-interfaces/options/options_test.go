package options

import (
	"bytes"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestDefaults(t *testing.T) {
	s, err := New("localhost")
	if err != nil {
		t.Fatal(err)
	}
	if s.Host() != "localhost" {
		t.Errorf("Host = %q", s.Host())
	}
	if s.Port() != 8080 {
		t.Errorf("Port = %d, want 8080", s.Port())
	}
	if s.Timeout() != 30*time.Second {
		t.Errorf("Timeout = %v, want 30s", s.Timeout())
	}
	if s.MaxConns() != 100 {
		t.Errorf("MaxConns = %d, want 100", s.MaxConns())
	}
	if s.TLS() {
		t.Error("TLS should default to false")
	}
	if s.Logger() != io.Discard {
		t.Error("Logger should default to io.Discard")
	}
	if tags := s.Tags(); tags == nil || len(tags) != 0 {
		t.Errorf("Tags = %v, want empty non-nil", tags)
	}
}

func TestOptionsApply(t *testing.T) {
	var buf bytes.Buffer
	s, err := New("example.com",
		WithPort(443),
		WithTLS(),
		WithTimeout(5*time.Second),
		WithMaxConns(7),
		WithLogger(&buf),
		WithTags("a", "b"),
		WithTags("c"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if s.Port() != 443 || !s.TLS() || s.Timeout() != 5*time.Second || s.MaxConns() != 7 {
		t.Errorf("options did not all apply: %+v", s)
	}
	if s.Logger() != io.Writer(&buf) {
		t.Error("WithLogger did not apply")
	}
	if got := s.Tags(); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("Tags = %v, want [a b c] (WithTags appends)", got)
	}
}

func TestLaterOptionWins(t *testing.T) {
	s, err := New("h", WithPort(1), WithPort(2))
	if err != nil {
		t.Fatal(err)
	}
	if s.Port() != 2 {
		t.Errorf("Port = %d, want 2", s.Port())
	}
}

func TestTagsAreCopied(t *testing.T) {
	s, _ := New("h", WithTags("a"))
	tags := s.Tags()
	tags[0] = "mutated"
	if s.Tags()[0] != "a" {
		t.Error("Tags() must return a copy")
	}
}

func TestValidation(t *testing.T) {
	tests := []struct {
		name    string
		opt     Option
		wantSub string
	}{
		{"port too low", WithPort(0), "port"},
		{"port too high", WithPort(70000), "port"},
		{"negative timeout", WithTimeout(-time.Second), "timeout"},
		{"zero timeout", WithTimeout(0), "timeout"},
		{"zero conns", WithMaxConns(0), "maxconns"},
		{"nil logger", WithLogger(nil), "logger"},
		{"empty tag", WithTags("ok", ""), "tag"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := New("h", tt.opt)
			if err == nil {
				t.Fatal("expected an error")
			}
			if !errors.Is(err, ErrInvalidOption) {
				t.Errorf("err = %v, should wrap ErrInvalidOption", err)
			}
			if !strings.Contains(strings.ToLower(err.Error()), tt.wantSub) {
				t.Errorf("err = %q, should mention %q", err, tt.wantSub)
			}
		})
	}
}

func TestEmptyHost(t *testing.T) {
	if _, err := New(""); err == nil {
		t.Error("New with an empty host must fail")
	}
}

func TestFirstErrorWins(t *testing.T) {
	_, err := New("h", WithPort(-1), WithMaxConns(-1))
	if err == nil || !strings.Contains(strings.ToLower(err.Error()), "port") {
		t.Errorf("err = %v, want the first failing option (port)", err)
	}
}

func TestGroup(t *testing.T) {
	prod := Group(WithPort(443), WithTLS(), WithMaxConns(1000))
	s, err := New("h", prod, WithMaxConns(5))
	if err != nil {
		t.Fatal(err)
	}
	if s.Port() != 443 || !s.TLS() || s.MaxConns() != 5 {
		t.Errorf("group + override gave %d %v %d", s.Port(), s.TLS(), s.MaxConns())
	}
	if _, err := New("h", Group(WithPort(0))); err == nil {
		t.Error("Group must propagate errors from the options inside it")
	}
}

func TestMustNew(t *testing.T) {
	s := MustNew("h", WithPort(99))
	if s.Port() != 99 {
		t.Errorf("Port = %d", s.Port())
	}
	defer func() {
		if recover() == nil {
			t.Error("MustNew with a bad option did not panic")
		}
	}()
	MustNew("h", WithPort(-5))
}
