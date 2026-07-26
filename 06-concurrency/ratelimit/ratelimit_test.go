package ratelimit

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"testing/synctest"
	"time"
)

func TestSemaphoreBasics(t *testing.T) {
	s := NewSemaphore(2)
	if s.Available() != 2 {
		t.Errorf("Available = %d, want 2", s.Available())
	}
	if !s.TryAcquire() || !s.TryAcquire() {
		t.Fatal("the first two TryAcquire calls should succeed")
	}
	if s.TryAcquire() {
		t.Error("the third TryAcquire should fail")
	}
	if s.Available() != 0 {
		t.Errorf("Available = %d, want 0", s.Available())
	}
	s.Release()
	if !s.TryAcquire() {
		t.Error("TryAcquire after Release should succeed")
	}
}

func TestSemaphoreBlocksAndUnblocks(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := NewSemaphore(1)
		if err := s.Acquire(t.Context()); err != nil {
			t.Fatal(err)
		}
		acquired := make(chan struct{})
		go func() {
			s.Acquire(context.Background())
			close(acquired)
		}()
		synctest.Wait()
		select {
		case <-acquired:
			t.Fatal("Acquire returned while the semaphore was full")
		default:
		}
		s.Release()
		synctest.Wait()
		select {
		case <-acquired:
		default:
			t.Fatal("Acquire did not return after Release")
		}
	})
}

func TestSemaphoreContext(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		s := NewSemaphore(1)
		s.Acquire(t.Context())
		ctx, cancel := context.WithTimeout(t.Context(), time.Second)
		defer cancel()
		start := time.Now()
		err := s.Acquire(ctx)
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("err = %v, want DeadlineExceeded", err)
		}
		if d := time.Since(start); d != time.Second {
			t.Errorf("returned after %v, want 1s", d)
		}
		if s.Available() != 0 {
			t.Error("a failed Acquire must not take a slot")
		}
	})
}

func TestSemaphoreOverRelease(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("releasing more than acquired did not panic")
		}
	}()
	NewSemaphore(1).Release()
}

func TestSemaphoreLimitsConcurrency(t *testing.T) {
	s := NewSemaphore(3)
	var cur, peak atomic.Int64
	var wg sync.WaitGroup
	for range 50 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := s.Acquire(context.Background()); err != nil {
				return
			}
			defer s.Release()
			c := cur.Add(1)
			for {
				p := peak.Load()
				if c <= p || peak.CompareAndSwap(p, c) {
					break
				}
			}
			time.Sleep(time.Millisecond)
			cur.Add(-1)
		}()
	}
	wg.Wait()
	if peak.Load() > 3 {
		t.Errorf("peak = %d, want at most 3", peak.Load())
	}
}

func TestBucketBurstThenRefill(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := NewBucket(10, 3) // 10/s, burst 3
		for i := range 3 {
			if !b.Allow() {
				t.Fatalf("Allow %d denied; the bucket starts full", i)
			}
		}
		if b.Allow() {
			t.Error("the fourth Allow should be denied")
		}
		time.Sleep(100 * time.Millisecond) // exactly one token
		if !b.Allow() {
			t.Error("a token should have refilled after 100ms at 10/s")
		}
		if b.Allow() {
			t.Error("only one token should have refilled")
		}
		// Idle for a long time: the bucket fills to burst, not beyond.
		time.Sleep(time.Hour)
		if got := b.Tokens(); got != 3 {
			t.Errorf("Tokens after an hour = %v, want 3 (capped at burst)", got)
		}
	})
}

func TestBucketAllowN(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := NewBucket(1, 5)
		if !b.AllowN(5) {
			t.Fatal("AllowN(5) on a full bucket of 5 should succeed")
		}
		if b.AllowN(1) {
			t.Error("the bucket should be empty")
		}
		time.Sleep(2 * time.Second)
		if b.AllowN(3) {
			t.Error("AllowN(3) with only 2 tokens should fail")
		}
		if !b.AllowN(2) {
			t.Error("AllowN(2) with 2 tokens should succeed")
		}
	})
}

func TestBucketWait(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := NewBucket(2, 1) // 2/s
		if err := b.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
		start := time.Now()
		if err := b.Wait(t.Context()); err != nil {
			t.Fatal(err)
		}
		if d := time.Since(start); d != 500*time.Millisecond {
			t.Errorf("Wait blocked for %v, want exactly 500ms", d)
		}
	})
}

func TestBucketWaitContext(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := NewBucket(1, 1)
		b.Wait(t.Context()) // drain
		ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()
		if err := b.Wait(ctx); !errors.Is(err, context.DeadlineExceeded) {
			t.Errorf("err = %v, want DeadlineExceeded", err)
		}
		time.Sleep(time.Second)
		if !b.Allow() {
			t.Error("the cancelled Wait consumed a token")
		}
	})
}

func TestBucketRate(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		b := NewBucket(100, 1)
		start := time.Now()
		n := 0
		for time.Since(start) < time.Second {
			if b.Allow() {
				n++
			}
			time.Sleep(time.Millisecond)
		}
		if n < 95 || n > 105 {
			t.Errorf("allowed %d requests in one second at 100/s", n)
		}
	})
}

func TestLimiter(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := NewLimiter(1000, 100, 2)
		var cur, peak atomic.Int64
		var wg sync.WaitGroup
		for range 20 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				l.Do(t.Context(), func() {
					c := cur.Add(1)
					if c > peak.Load() {
						peak.Store(c)
					}
					time.Sleep(10 * time.Millisecond)
					cur.Add(-1)
				})
			}()
		}
		wg.Wait()
		if peak.Load() > 2 {
			t.Errorf("peak concurrency = %d, want at most 2", peak.Load())
		}
	})
}

func TestLimiterReleasesOnPanic(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		l := NewLimiter(1000, 100, 1)
		func() {
			defer func() { recover() }()
			l.Do(t.Context(), func() { panic("boom") })
		}()
		done := make(chan struct{})
		go func() {
			l.Do(t.Context(), func() {})
			close(done)
		}()
		synctest.Wait()
		select {
		case <-done:
		default:
			t.Error("the concurrency slot was not released after a panic")
		}
	})
}

func TestDebounce(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		var n atomic.Int64
		call, stop := Debounce(100*time.Millisecond, func() { n.Add(1) })
		defer stop()

		for range 5 {
			call()
			time.Sleep(20 * time.Millisecond)
		}
		if n.Load() != 0 {
			t.Errorf("fn ran %d times during the burst, want 0", n.Load())
		}
		time.Sleep(200 * time.Millisecond)
		if n.Load() != 1 {
			t.Errorf("fn ran %d times, want 1", n.Load())
		}

		call()
		stop()
		time.Sleep(200 * time.Millisecond)
		if n.Load() != 1 {
			t.Errorf("stop() did not cancel the pending call (n = %d)", n.Load())
		}
	})
}
