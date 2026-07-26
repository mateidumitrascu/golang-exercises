// Package iterseq builds a lazy pipeline library on Go's range-over-func
// iterators.
//
//	iter.Seq[T]     = func(yield func(T) bool)
//	iter.Seq2[K, V] = func(yield func(K, V) bool)
//
// Two rules that everything here depends on:
//   - When yield returns false, the consumer wants to stop. You MUST return
//     immediately and not call yield again. Every combinator has to propagate
//     that signal down the chain.
//   - Nothing runs until someone ranges over the sequence. A pipeline of ten
//     Maps over a million elements does no work at all until it is consumed,
//     and does only 3 elements' worth if the consumer takes 3.
//
// iter.Pull turns a push iterator into a pull one, for when you need to advance
// two sequences in lockstep. Always defer its stop function.
package iterseq

import "iter"

// Range yields start, start+step, ... while it is on the correct side of stop
// (exclusive). It panics if step is 0, and yields nothing if step points the
// wrong way.
func Range(start, stop, step int) iter.Seq[int] { panic("TODO: implement Range") }

// FromSlice yields the elements of s. Collect is the inverse.
func FromSlice[T any](s []T) iter.Seq[T] { panic("TODO: implement FromSlice") }
func Collect[T any](seq iter.Seq[T]) []T { panic("TODO: implement Collect") }

// Map, Filter, Take, Drop and TakeWhile are lazy: they return immediately and
// do work only as the result is consumed.
func Map[A, B any](seq iter.Seq[A], f func(A) B) iter.Seq[B]       { panic("TODO: implement Map") }
func Filter[T any](seq iter.Seq[T], keep func(T) bool) iter.Seq[T] { panic("TODO: implement Filter") }

// Take yields at most n elements and then stops pulling from the source.
// n <= 0 yields nothing.
func Take[T any](seq iter.Seq[T], n int) iter.Seq[T] { panic("TODO: implement Take") }

// Drop skips the first n elements.
func Drop[T any](seq iter.Seq[T], n int) iter.Seq[T] { panic("TODO: implement Drop") }

// TakeWhile yields elements until pred fails, then stops.
func TakeWhile[T any](seq iter.Seq[T], pred func(T) bool) iter.Seq[T] {
	panic("TODO: implement TakeWhile")
}

// Chain concatenates sequences.
func Chain[T any](seqs ...iter.Seq[T]) iter.Seq[T] { panic("TODO: implement Chain") }

// Repeat yields v forever. It is infinite on purpose: combined with Take it
// shows that laziness actually works.
func Repeat[T any](v T) iter.Seq[T] { panic("TODO: implement Repeat") }

// Enumerate pairs each element with its index.
func Enumerate[T any](seq iter.Seq[T]) iter.Seq2[int, T] { panic("TODO: implement Enumerate") }

// Zip walks two sequences in lockstep, stopping when either runs out.
// This is the one that needs iter.Pull - and the one where forgetting to call
// stop leaks a goroutine.
func Zip[A, B any](a iter.Seq[A], b iter.Seq[B]) iter.Seq2[A, B] { panic("TODO: implement Zip") }

// MergeSorted merges two ascending sequences into one ascending sequence,
// keeping duplicates. Also needs iter.Pull.
func MergeSorted[T interface{ ~int | ~string | ~float64 }](a, b iter.Seq[T]) iter.Seq[T] {
	panic("TODO: implement MergeSorted")
}

// Chunk batches elements into slices of at most n. The final chunk may be
// shorter. Panics if n <= 0. Reusing one buffer across chunks would be a bug -
// the consumer may keep them.
func Chunk[T any](seq iter.Seq[T], n int) iter.Seq[[]T] { panic("TODO: implement Chunk") }

// Reduce, Count and First consume the sequence.
func Reduce[T, A any](seq iter.Seq[T], init A, f func(A, T) A) A { panic("TODO: implement Reduce") }
func Count[T any](seq iter.Seq[T]) int                           { panic("TODO: implement Count") }
func First[T any](seq iter.Seq[T]) (T, bool)                     { panic("TODO: implement First") }
