package syncprim

import (
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestOnceRunsOnce(t *testing.T) {
	var o Once
	var n atomic.Int64
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			o.Do(func() { n.Add(1) })
		}()
	}
	wg.Wait()
	if n.Load() != 1 {
		t.Errorf("f ran %d times, want 1", n.Load())
	}
	if !o.Done() {
		t.Error("Done() = false after Do")
	}
	var fresh Once
	if fresh.Done() {
		t.Error("a fresh Once reports Done")
	}
}

// TestOnceBlocksUntilComplete is the guarantee people forget: the second caller
// must not proceed while the first is still initialising.
func TestOnceBlocksUntilComplete(t *testing.T) {
	var o Once
	var ready atomic.Bool
	started := make(chan struct{})
	go func() {
		o.Do(func() {
			close(started)
			time.Sleep(50 * time.Millisecond)
			ready.Store(true)
		})
	}()
	<-started
	o.Do(func() { t.Error("f ran twice") })
	if !ready.Load() {
		t.Error("Do returned before the first call finished initialising")
	}
}

func TestOncePanicCountsAsDone(t *testing.T) {
	var o Once
	func() {
		defer func() { recover() }()
		o.Do(func() { panic("boom") })
	}()
	ran := false
	o.Do(func() { ran = true })
	if ran {
		t.Error("f ran again after a panicking first call; it must still count as done")
	}
}

func TestMap(t *testing.T) {
	var m Map[string, int]
	if _, ok := m.Get("x"); ok {
		t.Error("empty map returned a value")
	}
	m.Set("x", 1)
	m.Set("y", 2)
	if v, ok := m.Get("x"); !ok || v != 1 {
		t.Errorf("Get = %v, %v", v, ok)
	}
	if m.Len() != 2 {
		t.Errorf("Len = %d", m.Len())
	}
	m.Delete("x")
	if _, ok := m.Get("x"); ok || m.Len() != 1 {
		t.Error("Delete failed")
	}
	snap := m.Snapshot()
	snap["y"] = 99
	if v, _ := m.Get("y"); v != 2 {
		t.Error("Snapshot must be a copy")
	}
}

func TestMapConcurrent(t *testing.T) {
	var m Map[int, int]
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Add(2)
		go func() { defer wg.Done(); m.Set(i, i) }()
		go func() { defer wg.Done(); m.Get(i); m.Len(); m.Snapshot() }()
	}
	wg.Wait()
	if m.Len() != 50 {
		t.Errorf("Len = %d, want 50", m.Len())
	}
}

func TestGetOrCompute(t *testing.T) {
	var m Map[string, int]
	var calls atomic.Int64
	v, loaded := m.GetOrCompute("k", func() int { calls.Add(1); return 42 })
	if v != 42 || loaded {
		t.Errorf("= %v, %v; want 42, false", v, loaded)
	}
	v, loaded = m.GetOrCompute("k", func() int { calls.Add(1); return 99 })
	if v != 42 || !loaded {
		t.Errorf("= %v, %v; want 42, true", v, loaded)
	}
	if calls.Load() != 1 {
		t.Errorf("compute ran %d times, want 1", calls.Load())
	}

	var m2 Map[string, int]
	var n atomic.Int64
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m2.GetOrCompute("same", func() int { n.Add(1); return 1 })
		}()
	}
	wg.Wait()
	if n.Load() != 1 {
		t.Errorf("compute ran %d times under concurrency, want 1", n.Load())
	}
}

func TestCounter(t *testing.T) {
	var c Counter
	var wg sync.WaitGroup
	for range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Inc("hits")
			c.Add("bytes", 10)
		}()
	}
	wg.Wait()
	if c.Get("hits") != 100 {
		t.Errorf("hits = %d, want 100", c.Get("hits"))
	}
	if c.Get("bytes") != 1000 {
		t.Errorf("bytes = %d, want 1000", c.Get("bytes"))
	}
	if c.Get("never") != 0 {
		t.Errorf("unknown counter = %d, want 0", c.Get("never"))
	}
}

func TestCounterMax(t *testing.T) {
	var c Counter
	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.Max("peak", int64(i))
		}()
	}
	wg.Wait()
	if got := c.Get("peak"); got != 99 {
		t.Errorf("peak = %d, want 99", got)
	}
	if got := c.Max("peak", 5); got != 99 {
		t.Errorf("Max with a smaller value returned %d, want 99", got)
	}
	if got := c.Max("peak", 500); got != 500 {
		t.Errorf("Max = %d, want 500", got)
	}
}

func TestSingleflight(t *testing.T) {
	var g Group[string, int]
	var calls atomic.Int64
	release := make(chan struct{})
	var wg sync.WaitGroup
	results := make([]int, 20)
	sharedCount := atomic.Int64{}

	for i := range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			v, shared, err := g.Do("key", func() (int, error) {
				calls.Add(1)
				<-release
				return 42, nil
			})
			if err != nil {
				t.Error(err)
			}
			results[i] = v
			if shared {
				sharedCount.Add(1)
			}
		}()
	}
	time.Sleep(20 * time.Millisecond)
	close(release)
	wg.Wait()

	if calls.Load() != 1 {
		t.Errorf("fn ran %d times, want 1", calls.Load())
	}
	for i, v := range results {
		if v != 42 {
			t.Fatalf("goroutine %d got %d, want 42", i, v)
		}
	}
	if sharedCount.Load() == 0 {
		t.Error("no caller was told the result was shared")
	}

	// After the flight completes, the key is forgotten.
	g.Do("key", func() (int, error) { calls.Add(1); return 1, nil })
	if calls.Load() != 2 {
		t.Errorf("a later call reused a completed result (%d calls); Do is not a cache", calls.Load())
	}
}

func TestSingleflightErrors(t *testing.T) {
	var g Group[string, int]
	boom := errors.New("boom")
	var wg sync.WaitGroup
	for range 5 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, err := g.Do("k", func() (int, error) {
				time.Sleep(5 * time.Millisecond)
				return 0, boom
			})
			if !errors.Is(err, boom) {
				t.Errorf("err = %v, want boom", err)
			}
		}()
	}
	wg.Wait()
}

func TestSingleflightForget(t *testing.T) {
	var g Group[string, int]
	var calls atomic.Int64
	g.Do("k", func() (int, error) { calls.Add(1); return 1, nil })
	g.Forget("k")
	g.Do("k", func() (int, error) { calls.Add(1); return 1, nil })
	if calls.Load() != 2 {
		t.Errorf("calls = %d, want 2", calls.Load())
	}
	g.Forget("never seen") // must not panic
}
