// Package lru is the interview classic that is also genuinely useful: a cache
// with a size limit that evicts the least recently used entry.
//
// Every operation must be O(1). That rules out scanning a slice to find the
// oldest entry. The standard answer is a hash map from key to list node, plus a
// doubly linked list ordered by recency: the map gives you O(1) lookup, the
// list gives you O(1) reordering, and each node knows its key so eviction can
// remove it from the map.
//
// You may use container/list, or write the four pointer updates yourself.
// Writing them yourself is the better exercise, and it is about 30 lines.
package lru

import "time"

// Cache is an LRU cache with optional per-entry expiry.
type Cache[K comparable, V any] struct {
	// TODO: your fields.
	//
	// For testability, keep a `now func() time.Time` field that New sets to
	// time.Now. The tests replace it to control expiry without sleeping.
}

// Option configures a Cache.
type Option[K comparable, V any] func(*Cache[K, V])

// WithOnEvict registers a callback invoked for every entry that leaves the
// cache because of capacity or expiry - but NOT when it is replaced by Put or
// removed by Delete.
func WithOnEvict[K comparable, V any](fn func(K, V)) Option[K, V] {
	panic("TODO: implement WithOnEvict")
}

// WithClock replaces the time source.
func WithClock[K comparable, V any](now func() time.Time) Option[K, V] {
	panic("TODO: implement WithClock")
}

// New panics if capacity <= 0.
func New[K comparable, V any](capacity int, opts ...Option[K, V]) *Cache[K, V] {
	panic("TODO: implement New")
}

// Get returns the value and marks the entry as most recently used. An expired
// entry counts as a miss and is evicted on the spot.
func (c *Cache[K, V]) Get(k K) (V, bool) { panic("TODO: implement Get") }

// Peek returns the value WITHOUT changing recency. Useful for metrics and for
// tests; also the thing people forget to provide.
func (c *Cache[K, V]) Peek(k K) (V, bool) { panic("TODO: implement Peek") }

// Put inserts or updates. An update refreshes recency but does not change the
// insertion-order position of anything else. When the cache is over capacity,
// evict the least recently used entry.
func (c *Cache[K, V]) Put(k K, v V) { panic("TODO: implement Put") }

// PutTTL is Put with an expiry. A ttl <= 0 means "never expires".
func (c *Cache[K, V]) PutTTL(k K, v V, ttl time.Duration) { panic("TODO: implement PutTTL") }

// Delete removes an entry and reports whether it was there. No eviction
// callback: the caller asked for this.
func (c *Cache[K, V]) Delete(k K) bool { panic("TODO: implement Delete") }

// Len counts live (unexpired) entries. Cap is the configured capacity.
func (c *Cache[K, V]) Len() int { panic("TODO: implement Len") }
func (c *Cache[K, V]) Cap() int { panic("TODO: implement Cap") }

// Keys returns the keys from most to least recently used.
func (c *Cache[K, V]) Keys() []K { panic("TODO: implement Keys") }

// Resize changes the capacity, evicting from the least recently used end if it
// shrinks. It panics if n <= 0.
func (c *Cache[K, V]) Resize(n int) { panic("TODO: implement Resize") }

// Purge empties the cache, firing the eviction callback for everything.
func (c *Cache[K, V]) Purge() { panic("TODO: implement Purge") }

// Stats are cumulative counters. Hits and Misses come from Get only;
// Evictions counts capacity and expiry evictions.
type Stats struct {
	Hits, Misses, Evictions int
}

func (c *Cache[K, V]) Stats() Stats { panic("TODO: implement Stats") }
