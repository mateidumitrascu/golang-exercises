package pipeline

import (
	"context"
	"runtime"
	"sort"
	"testing"
	"time"
)

// assertNoLeaks gives goroutines a moment to exit and then fails if the count
// has grown. Call it with defer at the top of a test.
func assertNoLeaks(t *testing.T) {
	t.Helper()
	before := runtime.NumGoroutine()
	t.Cleanup(func() {
		for range 100 {
			runtime.Gosched()
			time.Sleep(time.Millisecond)
			if runtime.NumGoroutine() <= before {
				return
			}
		}
		buf := make([]byte, 1<<16)
		n := runtime.Stack(buf, true)
		t.Errorf("goroutine leak: %d before, %d after\n%s",
			before, runtime.NumGoroutine(), buf[:n])
	})
}

func collect[T any](c <-chan T) []T {
	var out []T
	for v := range c {
		out = append(out, v)
	}
	return out
}

func TestGenerate(t *testing.T) {
	assertNoLeaks(t)
	got := collect(Generate(context.Background(), 1, 2, 3))
	if len(got) != 3 || got[0] != 1 || got[2] != 3 {
		t.Errorf("= %v", got)
	}
	if got := collect(Generate[int](context.Background())); len(got) != 0 {
		t.Errorf("empty Generate = %v", got)
	}
}

func TestGenerateStopsOnCancel(t *testing.T) {
	assertNoLeaks(t)
	ctx, cancel := context.WithCancel(context.Background())
	c := Generate(ctx, 1, 2, 3, 4, 5)
	<-c // take one, abandon the rest
	cancel()
}

func TestTransformAndFilter(t *testing.T) {
	assertNoLeaks(t)
	ctx := context.Background()
	src := Generate(ctx, 1, 2, 3, 4, 5)
	odd := FilterChan(ctx, src, func(n int) bool { return n%2 == 1 })
	doubled := Transform(ctx, odd, func(n int) int { return n * 2 })
	got := collect(doubled)
	want := []int{2, 6, 10}
	for i := range want {
		if i >= len(got) || got[i] != want[i] {
			t.Fatalf("= %v, want %v", got, want)
		}
	}
}

func TestStagesStopOnCancel(t *testing.T) {
	assertNoLeaks(t)
	ctx, cancel := context.WithCancel(context.Background())
	src := Generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	out := Transform(ctx, src, func(n int) int { return n })
	<-out
	cancel()
	// Draining is allowed but not required; the goroutines must exit either way.
	for range out {
	}
}

func TestFanOutFanIn(t *testing.T) {
	assertNoLeaks(t)
	ctx := context.Background()
	src := Generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8)
	outs := FanOut(ctx, src, 4, func(n int) int {
		time.Sleep(time.Millisecond)
		return n * n
	})
	if len(outs) != 4 {
		t.Fatalf("FanOut returned %d channels, want 4", len(outs))
	}
	got := collect(FanIn(ctx, outs...))
	sort.Ints(got)
	want := []int{1, 4, 9, 16, 25, 36, 49, 64}
	if len(got) != len(want) {
		t.Fatalf("= %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("= %v, want %v", got, want)
		}
	}
}

func TestFanInClosesOnce(t *testing.T) {
	assertNoLeaks(t)
	ctx := context.Background()
	a := Generate(ctx, 1)
	b := Generate(ctx, 2)
	out := FanIn(ctx, a, b)
	collect(out)
	if _, ok := <-out; ok {
		t.Error("FanIn's channel must be closed after all inputs are drained")
	}
	if got := collect(FanIn[int](ctx)); len(got) != 0 {
		t.Error("FanIn with no inputs must close immediately")
	}
}

func TestOrDone(t *testing.T) {
	assertNoLeaks(t)
	ctx, cancel := context.WithCancel(context.Background())
	// An unbuffered channel nobody will ever write to again.
	src := make(chan int)
	out := OrDone(ctx, src)
	go func() {
		src <- 1
		cancel()
	}()
	if v := <-out; v != 1 {
		t.Errorf("= %d", v)
	}
	if _, ok := <-out; ok {
		t.Error("OrDone must close its output when the context is cancelled")
	}
}

func TestTee(t *testing.T) {
	assertNoLeaks(t)
	ctx := context.Background()
	a, b := Tee(ctx, Generate(ctx, 1, 2, 3))
	done := make(chan []int, 1)
	go func() { done <- collect(b) }()
	gotA := collect(a)
	gotB := <-done
	if len(gotA) != 3 || len(gotB) != 3 {
		t.Fatalf("a = %v, b = %v; both should get everything", gotA, gotB)
	}
	for i := range gotA {
		if gotA[i] != gotB[i] {
			t.Errorf("streams differ at %d: %v vs %v", i, gotA, gotB)
		}
	}
}

func TestBridge(t *testing.T) {
	assertNoLeaks(t)
	ctx := context.Background()
	chans := make(chan (<-chan int), 3)
	for i := range 3 {
		chans <- Generate(ctx, i*10, i*10+1)
	}
	close(chans)
	got := collect(Bridge(ctx, chans))
	want := []int{0, 1, 10, 11, 20, 21}
	if len(got) != len(want) {
		t.Fatalf("= %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("= %v, want %v (order must be preserved)", got, want)
		}
	}
}

func TestTake(t *testing.T) {
	assertNoLeaks(t)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// An infinite source: Take must not try to drain it.
	src := make(chan int)
	go func() {
		defer close(src)
		for i := 0; ; i++ {
			select {
			case src <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	got := collect(Take(ctx, src, 3))
	if len(got) != 3 || got[0] != 0 || got[2] != 2 {
		t.Errorf("= %v, want [0 1 2]", got)
	}
	if got := collect(Take(ctx, Generate(ctx, 1, 2), 0)); len(got) != 0 {
		t.Errorf("Take(0) = %v", got)
	}
	if got := collect(Take(ctx, Generate(ctx, 1), 5)); len(got) != 1 {
		t.Errorf("Take(more than available) = %v", got)
	}
}
