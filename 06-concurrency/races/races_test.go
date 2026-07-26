package races

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestSumConcurrent(t *testing.T) {
	nums := make([]int, 10000)
	want := 0
	for i := range nums {
		nums[i] = i
		want += i
	}
	for _, workers := range []int{1, 4, 8, 64} {
		for range 20 {
			if got := SumConcurrent(nums, workers); got != want {
				t.Fatalf("SumConcurrent(workers=%d) = %d, want %d", workers, got, want)
			}
		}
	}
	if got := SumConcurrent(nil, 4); got != 0 {
		t.Errorf("SumConcurrent(nil) = %d", got)
	}
}

func TestWordCount(t *testing.T) {
	texts := make([]string, 200)
	for i := range texts {
		texts[i] = "go is fun go go"
	}
	got := WordCount(texts)
	want := map[string]int{"go": 600, "is": 200, "fun": 200}
	for k, v := range want {
		if got[k] != v {
			t.Errorf("count[%q] = %d, want %d", k, got[k], v)
		}
	}
	if len(got) != len(want) {
		t.Errorf("got %d distinct words, want %d", len(got), len(want))
	}
}

func TestFirstErrorNoLeak(t *testing.T) {
	before := runtime.NumGoroutine()
	sentinel := errors.New("first")
	fns := make([]func() error, 100)
	for i := range fns {
		fns[i] = func() error {
			time.Sleep(time.Millisecond)
			return sentinel
		}
	}
	if err := FirstError(fns); !errors.Is(err, sentinel) {
		t.Fatalf("err = %v", err)
	}
	for range 200 {
		time.Sleep(time.Millisecond)
		if runtime.NumGoroutine() <= before+2 {
			return
		}
	}
	t.Errorf("goroutines leaked: %d before, %d after - the senders are still blocked",
		before, runtime.NumGoroutine())
}

func TestFirstErrorEmpty(t *testing.T) {
	if err := FirstError(nil); err != nil {
		t.Errorf("= %v, want nil", err)
	}
}

func TestCountUp(t *testing.T) {
	for range 50 {
		if got := CountUp(100); got != 100 {
			t.Fatalf("CountUp(100) = %d, want 100", got)
		}
	}
}

func TestBroadcaster(t *testing.T) {
	b := &Broadcaster{}
	c1 := b.Subscribe()
	c2 := b.Subscribe()
	if !b.Send("hello") {
		t.Error("Send on an open broadcaster returned false")
	}
	if v := <-c1; v != "hello" {
		t.Errorf("c1 = %q", v)
	}
	if v := <-c2; v != "hello" {
		t.Errorf("c2 = %q", v)
	}

	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b.Close() // must be idempotent and safe from many goroutines
		}()
	}
	wg.Wait()

	if b.Send("after close") {
		t.Error("Send after Close must return false, not panic and not send")
	}
	if _, ok := <-c1; ok {
		t.Error("subscriber channels must be closed")
	}
}

func TestBroadcasterConcurrentSendAndClose(t *testing.T) {
	for range 20 {
		b := &Broadcaster{}
		ch := b.Subscribe()
		go func() {
			for range ch {
			}
		}()
		var wg sync.WaitGroup
		for range 5 {
			wg.Add(1)
			go func() { defer wg.Done(); b.Send("x") }()
		}
		wg.Add(1)
		go func() { defer wg.Done(); b.Close() }()
		wg.Wait()
	}
}

func TestCacheComputesOncePerKey(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup
	var calls atomic.Int64
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v := c.Get("k", func() int {
				calls.Add(1)
				time.Sleep(5 * time.Millisecond)
				return 42
			})
			if v != 42 {
				t.Errorf("Get = %d", v)
			}
		}()
	}
	wg.Wait()
	if calls.Load() != 1 {
		t.Errorf("compute ran %d times for one key, want 1", calls.Load())
	}
	if c.Computed != 1 {
		t.Errorf("Computed = %d, want 1", c.Computed)
	}
}

func TestCacheDifferentKeys(t *testing.T) {
	c := NewCache()
	var wg sync.WaitGroup
	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			key := string(rune('a' + i%5))
			c.Get(key, func() int { return i })
		}()
	}
	wg.Wait()
	if c.Computed != 5 {
		t.Errorf("Computed = %d, want 5 (one per distinct key)", c.Computed)
	}
}

func TestNoLeftoverSentinel(t *testing.T) {
	// Delete ErrNotFixed from races.go when you are done. This test is a
	// checklist item, not a bug.
	if ErrNotFixed != nil {
		t.Skip("still working through the file")
	}
}
