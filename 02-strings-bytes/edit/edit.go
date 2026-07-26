// Package edit: dynamic programming over runes, with a memory budget.
//
// The textbook Levenshtein builds an (n+1)x(m+1) table. You are going to build
// the two-row version, because that is the one people actually ship, and one of
// the tests measures your allocations.
package edit

// Distance returns the Levenshtein edit distance between a and b: the minimum
// number of single-rune insertions, deletions or substitutions to turn a into b.
//
// It must compare RUNES, not bytes ("é" is one edit away from "e", not two),
// and use O(min(len(a), len(b))) extra space.
func Distance(a, b string) int {
	panic("TODO: implement Distance")
}

// DistanceAtMost is Distance with an early exit: if the true distance exceeds
// max, it may stop and return max+1. Useful when you only care whether two
// strings are "close enough". Must never be slower than Distance.
func DistanceAtMost(a, b string, max int) int {
	panic("TODO: implement DistanceAtMost")
}

// Similarity is 1 - Distance(a,b)/max(runeLen(a), runeLen(b)), so identical
// strings score 1 and completely different ones score 0. Two empty strings
// score 1.
func Similarity(a, b string) float64 {
	panic("TODO: implement Similarity")
}

// LCS returns a longest common subsequence of a and b (not substring:
// the runes need not be adjacent). If several are longest, return any one.
//
// LCS("ABCBDAB", "BDCABA") has length 4.
func LCS(a, b string) string {
	panic("TODO: implement LCS")
}

// Op is one step of an edit script.
type Op struct {
	Kind byte // '=' keep, '+' insert, '-' delete, '~' substitute
	Rune rune // the rune involved: the new one for '+' and '~', the old one otherwise
}

// Script returns an edit script that transforms a into b, with exactly
// Distance(a,b) non-'=' operations. This one needs the full table plus a
// backtrace, so the memory rule does not apply here.
func Script(a, b string) []Op {
	panic("TODO: implement Script")
}

// Apply runs an edit script against a and returns the result. Apply(a,
// Script(a, b)) must equal b - that is how Script gets checked.
func Apply(a string, ops []Op) string {
	panic("TODO: implement Apply")
}
