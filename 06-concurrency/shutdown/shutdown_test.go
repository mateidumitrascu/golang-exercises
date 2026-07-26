package shutdown

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGroupSuccess(t *testing.T) {
	var g Group
	var n atomic.Int64
	for range 10 {
		g.Go(func() error { n.Add(1); return nil })
	}
	if err := g.Wait(); err != nil {
		t.Errorf("Wait = %v, want nil", err)
	}
	if n.Load() != 10 {
		t.Errorf("ran %d, want 10", n.Load())
	}
}

func TestGroupFirstError(t *testing.T) {
	var g Group
	first := errors.New("first")
	second := errors.New("second")
	g.Go(func() error { return first })
	time.Sleep(10 * time.Millisecond)
	g.Go(func() error { return second })
	if err := g.Wait(); !errors.Is(err, first) {
		t.Errorf("Wait = %v, want the first error", err)
	}
}

func TestGroupContextCancelledOnError(t *testing.T) {
	g, ctx := WithContext(context.Background())
	boom := errors.New("boom")
	g.Go(func() error { return boom })
	g.Go(func() error {
		<-ctx.Done()
		return nil
	})
	if err := g.Wait(); !errors.Is(err, boom) {
		t.Fatalf("Wait = %v", err)
	}
	if ctx.Err() == nil {
		t.Error("the group's context must be cancelled after Wait")
	}
}

func TestGroupCancelsOnSuccessToo(t *testing.T) {
	g, ctx := WithContext(context.Background())
	g.Go(func() error { return nil })
	g.Wait()
	if ctx.Err() == nil {
		t.Error("Wait must cancel the context even when nothing failed")
	}
}

func TestGroupRecoversPanics(t *testing.T) {
	var g Group
	g.Go(func() error { panic("task exploded") })
	err := g.Wait()
	if err == nil {
		t.Fatal("a panicking task must surface as an error")
	}
	if !strings.Contains(err.Error(), "exploded") {
		t.Errorf("err = %q, should mention the panic value", err)
	}
}

func TestGroupLimit(t *testing.T) {
	var g Group
	g.SetLimit(3)
	var cur, peak atomic.Int64
	for range 20 {
		g.Go(func() error {
			c := cur.Add(1)
			for {
				p := peak.Load()
				if c <= p || peak.CompareAndSwap(p, c) {
					break
				}
			}
			time.Sleep(2 * time.Millisecond)
			cur.Add(-1)
			return nil
		})
	}
	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}
	if peak.Load() > 3 {
		t.Errorf("peak = %d, want at most 3", peak.Load())
	}
}

func TestGroupTryGo(t *testing.T) {
	var g Group
	g.SetLimit(1)
	block := make(chan struct{})
	if !g.TryGo(func() error { <-block; return nil }) {
		t.Fatal("the first TryGo should succeed")
	}
	if g.TryGo(func() error { return nil }) {
		t.Error("TryGo should fail when the limit is reached")
	}
	close(block)
	g.Wait()
	if !g.TryGo(func() error { return nil }) {
		t.Error("TryGo should succeed again once a slot is free")
	}
	g.Wait()
}

func TestGroupSetLimitPanics(t *testing.T) {
	var g Group
	block := make(chan struct{})
	g.Go(func() error { <-block; return nil })
	func() {
		defer func() {
			if recover() == nil {
				t.Error("SetLimit with active goroutines did not panic")
			}
			close(block)
		}()
		g.SetLimit(5)
	}()
	g.Wait()
}

func TestShutdownReverseOrder(t *testing.T) {
	var mu sync.Mutex
	var order []string
	mk := func(name string) Closer {
		return Closer{Name: name, Close: func(context.Context) error {
			mu.Lock()
			defer mu.Unlock()
			order = append(order, name)
			return nil
		}}
	}
	err := Shutdown(context.Background(), time.Second, mk("db"), mk("cache"), mk("server"))
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"server", "cache", "db"}; !reflect.DeepEqual(order, want) {
		t.Errorf("order = %v, want %v (reverse, like defer)", order, want)
	}
}

func TestShutdownCollectsErrors(t *testing.T) {
	e1, e2 := errors.New("db is angry"), errors.New("cache is angry")
	err := Shutdown(context.Background(), time.Second,
		Closer{Name: "db", Close: func(context.Context) error { return e1 }},
		Closer{Name: "ok", Close: func(context.Context) error { return nil }},
		Closer{Name: "cache", Close: func(context.Context) error { return e2 }},
	)
	if !errors.Is(err, e1) || !errors.Is(err, e2) {
		t.Fatalf("err = %v, want both failures", err)
	}
	for _, name := range []string{"db", "cache"} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("err = %q, should name %q", err, name)
		}
	}
}

func TestShutdownTimeout(t *testing.T) {
	slowDone := make(chan struct{})
	defer close(slowDone)
	var lastRan atomic.Bool
	start := time.Now()
	err := Shutdown(context.Background(), 100*time.Millisecond,
		Closer{Name: "first", Close: func(ctx context.Context) error {
			lastRan.Store(true)
			return nil
		}},
		Closer{Name: "slow", Close: func(ctx context.Context) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-slowDone:
				return nil
			}
		}},
	)
	if elapsed := time.Since(start); elapsed > 500*time.Millisecond {
		t.Errorf("Shutdown took %v, want it bounded by the timeout", elapsed)
	}
	if err == nil || !strings.Contains(err.Error(), "slow") {
		t.Errorf("err = %v, want it to name the closer that timed out", err)
	}
	if !lastRan.Load() {
		t.Error("the remaining closers must still be attempted after one times out")
	}
}

func TestWorkerDrains(t *testing.T) {
	in := make(chan int, 10)
	for i := range 10 {
		in <- i
	}
	var sum atomic.Int64
	w := NewWorker(in, func(n int) {
		sum.Add(int64(n))
		time.Sleep(time.Millisecond)
	})
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(3 * time.Millisecond)
		cancel()
	}()
	err := w.Run(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Run = %v, want context.Canceled", err)
	}
	if w.Processed() != 10 {
		t.Errorf("processed %d jobs, want all 10 (drain what is already queued)", w.Processed())
	}
	if sum.Load() != 45 {
		t.Errorf("sum = %d, want 45", sum.Load())
	}
}

func TestWorkerEndsWithInput(t *testing.T) {
	in := make(chan int, 3)
	for i := range 3 {
		in <- i
	}
	close(in)
	w := NewWorker(in, func(int) {})
	if err := w.Run(context.Background()); err != nil {
		t.Errorf("Run = %v, want nil when the input closes normally", err)
	}
	if w.Processed() != 3 {
		t.Errorf("processed %d, want 3", w.Processed())
	}
}

func ExampleGroup() {
	g, ctx := WithContext(context.Background())
	g.SetLimit(2)
	for i := range 3 {
		g.Go(func() error {
			if i == 1 {
				return fmt.Errorf("task %d failed", i)
			}
			<-ctx.Done()
			return nil
		})
	}
	fmt.Println(g.Wait())
	// Output: task 1 failed
}
