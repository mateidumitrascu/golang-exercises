package lru

import (
	"reflect"
	"testing"
	"time"
)

func TestBasics(t *testing.T) {
	c := New[string, int](2)
	if c.Cap() != 2 || c.Len() != 0 {
		t.Fatalf("cap %d len %d", c.Cap(), c.Len())
	}
	if _, ok := c.Get("nope"); ok {
		t.Error("empty cache returned a value")
	}
	c.Put("a", 1)
	c.Put("b", 2)
	if v, ok := c.Get("a"); !ok || v != 1 {
		t.Errorf("Get(a) = %v, %v", v, ok)
	}
	if c.Len() != 2 {
		t.Errorf("Len = %d", c.Len())
	}
	c.Put("a", 10) // update, not insert
	if v, _ := c.Get("a"); v != 10 || c.Len() != 2 {
		t.Errorf("after update: %v, len %d", v, c.Len())
	}
	if !c.Delete("a") || c.Delete("a") {
		t.Error("Delete is wrong")
	}
	if c.Len() != 1 {
		t.Errorf("Len = %d after delete", c.Len())
	}
}

func TestEviction(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Get("a")    // a is now the most recent
	c.Put("c", 3) // b should go
	if _, ok := c.Peek("b"); ok {
		t.Error("b should have been evicted (it was least recently used)")
	}
	if _, ok := c.Peek("a"); !ok {
		t.Error("a was used recently and should have survived")
	}
	if got := c.Keys(); !reflect.DeepEqual(got, []string{"c", "a"}) {
		t.Errorf("Keys = %v, want [c a] (most recent first)", got)
	}
}

func TestPeekDoesNotRefresh(t *testing.T) {
	c := New[string, int](2)
	c.Put("a", 1)
	c.Put("b", 2)
	c.Peek("a")
	c.Put("c", 3)
	if _, ok := c.Peek("a"); ok {
		t.Error("Peek must not count as a use, so a should have been evicted")
	}
}

func TestOnEvict(t *testing.T) {
	var evicted []string
	c := New(2, WithOnEvict[string, int](func(k string, v int) {
		evicted = append(evicted, k)
	}))
	c.Put("a", 1)
	c.Put("b", 2)
	c.Put("c", 3)
	if !reflect.DeepEqual(evicted, []string{"a"}) {
		t.Errorf("evicted = %v, want [a]", evicted)
	}
	c.Put("c", 30) // replacement, not eviction
	c.Delete("b")  // explicit, not eviction
	if len(evicted) != 1 {
		t.Errorf("evicted = %v; replacing and deleting must not fire the callback", evicted)
	}
	c.Purge()
	if len(evicted) != 2 {
		t.Errorf("evicted = %v; Purge must fire it for everything left", evicted)
	}
}

func TestTTL(t *testing.T) {
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	c := New(10, WithClock[string, int](func() time.Time { return now }))
	c.PutTTL("short", 1, time.Second)
	c.PutTTL("long", 2, time.Hour)
	c.Put("forever", 3)

	if _, ok := c.Get("short"); !ok {
		t.Error("not expired yet")
	}
	now = now.Add(2 * time.Second)
	if _, ok := c.Get("short"); ok {
		t.Error("expired entries must read as missing")
	}
	if c.Len() != 2 {
		t.Errorf("Len = %d, want 2 (the expired entry should be gone)", c.Len())
	}
	if _, ok := c.Get("long"); !ok {
		t.Error("long should still be alive")
	}
	now = now.Add(2 * time.Hour)
	if _, ok := c.Get("long"); ok {
		t.Error("long should have expired")
	}
	if _, ok := c.Get("forever"); !ok {
		t.Error("a ttl of 0 means no expiry")
	}
	if got := c.Stats().Evictions; got != 2 {
		t.Errorf("Evictions = %d, want 2", got)
	}
}

func TestResize(t *testing.T) {
	c := New[int, int](5)
	for i := range 5 {
		c.Put(i, i)
	}
	c.Resize(2)
	if c.Cap() != 2 || c.Len() != 2 {
		t.Errorf("cap %d len %d, want 2 2", c.Cap(), c.Len())
	}
	if got := c.Keys(); !reflect.DeepEqual(got, []int{4, 3}) {
		t.Errorf("Keys = %v, want [4 3]", got)
	}
	c.Resize(4)
	c.Put(9, 9)
	if c.Len() != 3 {
		t.Errorf("Len = %d after growing", c.Len())
	}
	defer func() {
		if recover() == nil {
			t.Error("Resize(0) did not panic")
		}
	}()
	c.Resize(0)
}

func TestNewPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("New(0) did not panic")
		}
	}()
	New[string, int](0)
}

func TestStats(t *testing.T) {
	c := New[string, int](1)
	c.Put("a", 1)
	c.Get("a")
	c.Get("a")
	c.Get("zzz")
	c.Put("b", 2) // evicts a
	got := c.Stats()
	want := Stats{Hits: 2, Misses: 1, Evictions: 1}
	if got != want {
		t.Errorf("Stats = %+v, want %+v", got, want)
	}
}

// TestConstantTime hammers the cache; an O(n) eviction scan turns this into
// O(n^2) and it will not finish in a sensible time.
func TestConstantTime(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	const n = 300000
	c := New[int, int](1000)
	for i := range n {
		c.Put(i, i)
		if i%3 == 0 {
			c.Get(i - 500)
		}
	}
	if c.Len() != 1000 {
		t.Errorf("Len = %d, want 1000", c.Len())
	}
}

func BenchmarkPutGet(b *testing.B) {
	c := New[int, int](1024)
	b.ReportAllocs()
	i := 0
	for b.Loop() {
		c.Put(i%2048, i)
		c.Get(i % 2048)
		i++
	}
}
