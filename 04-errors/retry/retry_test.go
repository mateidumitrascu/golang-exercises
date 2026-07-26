package retry

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

func TestDelay(t *testing.T) {
	p := Policy{Base: time.Second, Multiplier: 2, Max: 10 * time.Second}
	want := []time.Duration{0, time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 10 * time.Second, 10 * time.Second}
	for i, w := range want {
		if got := p.Delay(i); got != w {
			t.Errorf("Delay(%d) = %v, want %v", i, got, w)
		}
	}
	if got := p.Delay(-1); got != 0 {
		t.Errorf("Delay(-1) = %v, want 0", got)
	}

	noCap := Policy{Base: time.Second, Multiplier: 3}
	if got := noCap.Delay(3); got != 9*time.Second {
		t.Errorf("Delay(3) = %v, want 9s", got)
	}
	// Multiplier <= 1 must behave as 2.
	def := Policy{Base: time.Second}
	if got := def.Delay(2); got != 2*time.Second {
		t.Errorf("default multiplier: Delay(2) = %v, want 2s", got)
	}
	// No overflow, no negative durations.
	huge := Policy{Base: time.Hour, Multiplier: 10, Max: 24 * time.Hour}
	if got := huge.Delay(40); got != 24*time.Hour {
		t.Errorf("Delay(40) = %v, want 24h (watch for int64 overflow)", got)
	}
	jittered := Policy{Base: time.Second, Jitter: func(d time.Duration) time.Duration { return d / 2 }}
	if got := jittered.Delay(1); got != 500*time.Millisecond {
		t.Errorf("jittered Delay(1) = %v, want 500ms", got)
	}
}

func TestIsRetryable(t *testing.T) {
	plain := errors.New("plain")
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"nil", nil, false},
		{"plain", plain, true},
		{"permanent", Stop(plain), false},
		{"wrapped permanent", fmt.Errorf("ctx: %w", Stop(plain)), false},
		{"canceled", context.Canceled, false},
		{"deadline", context.DeadlineExceeded, false},
		{"wrapped deadline", fmt.Errorf("x: %w", context.DeadlineExceeded), false},
		{"custom yes", customErr{true}, true},
		{"custom no", customErr{false}, false},
		{"wrapped custom no", fmt.Errorf("x: %w", customErr{false}), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsRetryable(tt.err); got != tt.want {
				t.Errorf("IsRetryable(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
	if Stop(nil) != nil {
		t.Error("Stop(nil) must be nil")
	}
}

type customErr struct{ ok bool }

func (c customErr) Error() string   { return "custom" }
func (c customErr) Retryable() bool { return c.ok }

var testPolicy = Policy{MaxAttempts: 4, Base: time.Second, Multiplier: 2, Max: time.Minute}

func TestDoSucceedsFirstTry(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := 0
		start := time.Now()
		err := Do(t.Context(), testPolicy, func(context.Context) error {
			calls++
			return nil
		})
		if err != nil || calls != 1 {
			t.Fatalf("err = %v, calls = %d", err, calls)
		}
		if d := time.Since(start); d != 0 {
			t.Errorf("slept %v before the first attempt, want 0", d)
		}
	})
}

func TestDoBacksOff(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := 0
		start := time.Now()
		err := Do(t.Context(), testPolicy, func(context.Context) error {
			calls++
			if calls < 3 {
				return errors.New("temporary")
			}
			return nil
		})
		if err != nil {
			t.Fatalf("err = %v", err)
		}
		if calls != 3 {
			t.Errorf("calls = %d, want 3", calls)
		}
		// 1s after attempt 1, 2s after attempt 2.
		if got := time.Since(start); got != 3*time.Second {
			t.Errorf("elapsed = %v, want 3s", got)
		}
	})
}

func TestDoExhausts(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		last := errors.New("still broken")
		calls := 0
		start := time.Now()
		err := Do(t.Context(), testPolicy, func(context.Context) error {
			calls++
			return last
		})
		if calls != 4 {
			t.Errorf("calls = %d, want MaxAttempts (4)", calls)
		}
		if !errors.Is(err, last) {
			t.Errorf("err = %v, want it to wrap the last error", err)
		}
		if !errors.Is(err, ErrExhausted) {
			t.Errorf("err = %v, want it to wrap ErrExhausted", err)
		}
		if !strings.Contains(err.Error(), "4") {
			t.Errorf("err = %q, should mention the attempt count", err)
		}
		if got := time.Since(start); got != 7*time.Second {
			t.Errorf("elapsed = %v, want 7s (1+2+4)", got)
		}
	})
}

func TestDoStopsOnPermanent(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		fatal := errors.New("bad request")
		calls := 0
		err := Do(t.Context(), testPolicy, func(context.Context) error {
			calls++
			return Stop(fatal)
		})
		if calls != 1 {
			t.Errorf("calls = %d, want 1", calls)
		}
		if !errors.Is(err, fatal) {
			t.Errorf("err = %v, want the wrapped error", err)
		}
		var p *Permanent
		if errors.As(err, &p) {
			t.Errorf("err = %#v: unwrap the *Permanent before returning it", err)
		}
	})
}

func TestDoRespectsContext(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 2500*time.Millisecond)
		defer cancel()
		boom := errors.New("boom")
		calls := 0
		start := time.Now()
		err := Do(ctx, testPolicy, func(context.Context) error {
			calls++
			return boom
		})
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("err = %v, want it to wrap context.DeadlineExceeded", err)
		}
		if !errors.Is(err, boom) {
			t.Errorf("err = %v, want it to also wrap the last attempt's error", err)
		}
		if calls != 2 {
			t.Errorf("calls = %d, want 2 (the third sleep is cut short)", calls)
		}
		if got := time.Since(start); got != 2500*time.Millisecond {
			t.Errorf("returned after %v, want as soon as the context expired", got)
		}
	})
}

func TestDoValue(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		calls := 0
		v, err := DoValue(t.Context(), testPolicy, func(context.Context) (string, error) {
			calls++
			if calls < 2 {
				return "", errors.New("nope")
			}
			return "ok", nil
		})
		if v != "ok" || err != nil {
			t.Fatalf("= %q, %v", v, err)
		}
		v, err = DoValue(t.Context(), Policy{MaxAttempts: 1}, func(context.Context) (string, error) {
			return "", errors.New("nope")
		})
		if v != "" || err == nil {
			t.Errorf("= %q, %v; want the zero value and an error", v, err)
		}
	})
}
