package pool

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestParallelMapOrder(t *testing.T) {
	in := make([]int, 100)
	for i := range in {
		in[i] = i
	}
	// Reverse the durations so that completion order is the opposite of input order.
	got, err := ParallelMap(context.Background(), in, 8, func(ctx context.Context, n int) (string, error) {
		time.Sleep(time.Duration(100-n) * time.Microsecond)
		return fmt.Sprint(n), nil
	})
	if err != nil {
		t.Fatal(err)
	}
	for i := range in {
		if got[i] != fmt.Sprint(i) {
			t.Fatalf("results are not in input order: got[%d] = %q", i, got[i])
		}
	}
}

func TestParallelMapBoundsConcurrency(t *testing.T) {
	var current, peak atomic.Int64
	in := make([]int, 200)
	const workers = 5
	_, err := ParallelMap(context.Background(), in, workers, func(ctx context.Context, n int) (int, error) {
		cur := current.Add(1)
		for {
			old := peak.Load()
			if cur <= old || peak.CompareAndSwap(old, cur) {
				break
			}
		}
		time.Sleep(time.Millisecond)
		current.Add(-1)
		return n, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if peak.Load() > workers {
		t.Errorf("peak concurrency was %d, want at most %d", peak.Load(), workers)
	}
	if peak.Load() < 2 {
		t.Error("no parallelism at all - are you running the work sequentially?")
	}
}

func TestParallelMapEmpty(t *testing.T) {
	before := runtime.NumGoroutine()
	got, err := ParallelMap(context.Background(), []int{}, 4, func(context.Context, int) (int, error) {
		t.Error("f must not be called for an empty input")
		return 0, nil
	})
	if err != nil || len(got) != 0 {
		t.Errorf("= %v, %v", got, err)
	}
	time.Sleep(10 * time.Millisecond)
	if after := runtime.NumGoroutine(); after > before {
		t.Errorf("started %d goroutines for an empty input", after-before)
	}
}

func TestParallelMapFirstErrorWins(t *testing.T) {
	boom := errors.New("boom")
	var calls atomic.Int64
	in := make([]int, 500)
	for i := range in {
		in[i] = i
	}
	got, err := ParallelMap(context.Background(), in, 4, func(ctx context.Context, n int) (int, error) {
		calls.Add(1)
		if n == 10 {
			return 0, boom
		}
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		case <-time.After(time.Millisecond):
		}
		return n, nil
	})
	if !errors.Is(err, boom) {
		t.Fatalf("err = %v, want boom", err)
	}
	if got != nil {
		t.Errorf("results = %v, want nil on error", got)
	}
	if calls.Load() >= int64(len(in)) {
		t.Errorf("f was called %d times out of %d; the rest should have been skipped", calls.Load(), len(in))
	}
}

func TestParallelMapContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := ParallelMap(ctx, []int{1, 2, 3}, 2, func(context.Context, int) (int, error) {
		t.Error("f must not be called with an already-cancelled context")
		return 0, nil
	})
	if !errors.Is(err, context.Canceled) {
		t.Errorf("err = %v, want context.Canceled", err)
	}
}

func TestParallelMapPanicsOnBadWorkers(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("workers = 0 did not panic")
		}
	}()
	ParallelMap(context.Background(), []int{1}, 0, func(context.Context, int) (int, error) { return 0, nil })
}

func TestForEachCollectsAllErrors(t *testing.T) {
	in := []int{1, 2, 3, 4, 5}
	err := ForEach(context.Background(), in, 3, func(_ context.Context, n int) error {
		if n%2 == 1 {
			return fmt.Errorf("odd: %d", n)
		}
		return nil
	})
	if err == nil {
		t.Fatal("want an error")
	}
	for _, want := range []string{"1", "3", "5"} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("err = %q, should mention every failure (%s missing)", err, want)
		}
	}
	if err := ForEach(context.Background(), in, 3, func(context.Context, int) error { return nil }); err != nil {
		t.Errorf("err = %v, want nil", err)
	}
}

func TestPoolRunsEverything(t *testing.T) {
	p := New(4, 10)
	var count atomic.Int64
	for range 100 {
		if err := p.Submit(func() { count.Add(1) }); err != nil {
			t.Fatal(err)
		}
	}
	p.Close()
	p.Wait()
	if count.Load() != 100 {
		t.Errorf("ran %d tasks, want 100", count.Load())
	}
}

func TestPoolWorkerCount(t *testing.T) {
	before := runtime.NumGoroutine()
	p := New(3, 1)
	time.Sleep(20 * time.Millisecond)
	if got := runtime.NumGoroutine() - before; got != 3 {
		t.Errorf("New(3, 1) started %d goroutines, want exactly 3", got)
	}
	var running, peak atomic.Int64
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		p.Submit(func() {
			defer wg.Done()
			c := running.Add(1)
			if c > peak.Load() {
				peak.Store(c)
			}
			time.Sleep(2 * time.Millisecond)
			running.Add(-1)
		})
	}
	wg.Wait()
	if peak.Load() > 3 {
		t.Errorf("peak concurrency %d, want at most 3", peak.Load())
	}
	p.Close()
	p.Wait()
	time.Sleep(20 * time.Millisecond)
	if leaked := runtime.NumGoroutine() - before; leaked > 0 {
		t.Errorf("%d goroutines still running after Close+Wait", leaked)
	}
}

func TestPoolCloseAndSubmit(t *testing.T) {
	p := New(2, 2)
	p.Close()
	p.Close() // must be idempotent
	if err := p.Submit(func() {}); !errors.Is(err, ErrPoolClosed) {
		t.Errorf("Submit after Close = %v, want ErrPoolClosed", err)
	}
	if _, err := p.TrySubmit(func() {}); !errors.Is(err, ErrPoolClosed) {
		t.Errorf("TrySubmit after Close = %v, want ErrPoolClosed", err)
	}
	p.Wait()
}

func TestPoolTrySubmit(t *testing.T) {
	p := New(1, 1)
	block := make(chan struct{})
	p.Submit(func() { <-block })     // occupies the worker
	ok1, _ := p.TrySubmit(func() {}) // fills the queue
	full := false
	for range 10 {
		ok, err := p.TrySubmit(func() {})
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			full = true
			break
		}
	}
	if !ok1 {
		t.Error("the first TrySubmit should have fit in the queue")
	}
	if !full {
		t.Error("TrySubmit must return false when the queue is full instead of blocking")
	}
	close(block)
	p.Close()
	p.Wait()
}

func TestPoolSurvivesPanics(t *testing.T) {
	p := New(2, 4)
	var ok atomic.Int64
	for i := range 10 {
		p.Submit(func() {
			if i%2 == 0 {
				panic("task blew up")
			}
			ok.Add(1)
		})
	}
	p.Close()
	p.Wait()
	if ok.Load() != 5 {
		t.Errorf("%d good tasks completed, want 5 - a panicking task must not kill the worker", ok.Load())
	}
	if p.Panics() != 5 {
		t.Errorf("Panics = %d, want 5", p.Panics())
	}
}
