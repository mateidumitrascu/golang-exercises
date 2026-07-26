// Package heapq is the binary heap: the data structure behind priority queues,
// top-k queries, k-way merges, event simulations and schedulers.
//
// The whole thing is an array. The children of i are 2i+1 and 2i+2; the parent
// of i is (i-1)/2. Sifting up and sifting down are ten lines each. Once you
// have written them once you will never be afraid of container/heap again.
package heapq

// Heap is a generic binary heap ordered by a less function. less(a, b) reports
// whether a must come out before b, so a min-heap over ints uses a < b.
type Heap[T any] struct {
	// TODO: your fields: the slice and the less function.
}

// New builds an empty heap. It panics if less is nil.
func New[T any](less func(a, b T) bool) *Heap[T] { panic("TODO: implement New") }

// From builds a heap from an existing slice in O(n) - by sifting down from the
// last parent, NOT by pushing n times, which is O(n log n). It takes ownership
// of the slice.
func From[T any](s []T, less func(a, b T) bool) *Heap[T] { panic("TODO: implement From") }

func (h *Heap[T]) Len() int { panic("TODO: implement Len") }

// Push adds a value in O(log n).
func (h *Heap[T]) Push(v T) { panic("TODO: implement Push") }

// Pop removes and returns the smallest value (by less). It returns false when
// the heap is empty.
func (h *Heap[T]) Pop() (T, bool) { panic("TODO: implement Pop") }

// Peek returns the smallest value without removing it.
func (h *Heap[T]) Peek() (T, bool) { panic("TODO: implement Peek") }

// PushPop pushes v then pops, but does it in one sift instead of two - a
// worthwhile optimisation for the "keep the k largest" loop.
func (h *Heap[T]) PushPop(v T) (T, bool) { panic("TODO: implement PushPop") }

// Slice exposes the underlying array, for tests that want to check the heap
// invariant. Do not mutate it.
func (h *Heap[T]) Slice() []T { panic("TODO: implement Slice") }

// Item is an entry in a priority queue: a value with a priority, plus the index
// that makes Update possible in O(log n) instead of O(n).
type Item[T any] struct {
	Value    T
	Priority int
	index    int // maintained by the queue
}

// PQ is a priority queue where the LOWEST priority number comes out first, and
// where an item's priority can be changed after it has been queued (which is
// what Dijkstra and event schedulers need).
type PQ[T any] struct {
	// TODO
}

func NewPQ[T any]() *PQ[T] { panic("TODO: implement NewPQ") }

// Push queues value and returns the item handle, which the caller keeps if it
// wants to Update later.
func (q *PQ[T]) Push(value T, priority int) *Item[T] { panic("TODO: implement PQ.Push") }

// Pop returns the lowest-priority-number item. Ties may come out in any order.
func (q *PQ[T]) Pop() (*Item[T], bool) { panic("TODO: implement PQ.Pop") }

// Update changes an item's priority and restores the heap in O(log n). It is a
// no-op for an item that has already been popped.
func (q *PQ[T]) Update(item *Item[T], priority int) { panic("TODO: implement PQ.Update") }

func (q *PQ[T]) Len() int { panic("TODO: implement PQ.Len") }

// TopK returns the k largest elements of s (by less: the k for which less
// returns false against the rest), in descending order. It must run in
// O(n log k) time and O(k) space - keep a min-heap of size k, not a sorted
// copy of everything. If k >= len(s) it returns everything, sorted.
func TopK[T any](s []T, k int, less func(a, b T) bool) []T { panic("TODO: implement TopK") }

// HeapSort sorts s in place using a heap, in O(n log n) with no extra
// allocation. (It is not stable - say so in a comment, and know why.)
func HeapSort[T any](s []T, less func(a, b T) bool) { panic("TODO: implement HeapSort") }

// MergeSorted merges k sorted slices into one sorted slice in O(n log k), using
// a heap of the current head of each input. Concatenating and sorting is
// O(n log n) - do it the right way.
func MergeSorted[T any](lists [][]T, less func(a, b T) bool) []T {
	panic("TODO: implement MergeSorted")
}
