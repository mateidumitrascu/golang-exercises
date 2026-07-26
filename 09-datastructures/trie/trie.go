// Package trie is the prefix tree: the structure behind autocomplete, routers,
// spell checkers and IP lookup tables.
//
// A node holds a map (or array) of children keyed by the next rune, and a flag
// for "a word ends here". Lookup is O(len(key)) and does not depend on how many
// words are stored - which is the whole point.
package trie

import "iter"

// Trie stores a set of strings. The zero value is not usable; call New.
type Trie struct {
	// TODO: your fields. A root node, and a count of words for O(1) Len.
	//
	// The node type is up to you. map[rune]*node is flexible;
	// [26]*node is faster but only works for a fixed alphabet.
}

func New() *Trie { panic("TODO: implement New") }

// Insert adds a word and reports whether it was new. The empty string is a
// valid word.
func (t *Trie) Insert(word string) bool { panic("TODO: implement Insert") }

// Contains reports whether the exact word is stored.
func (t *Trie) Contains(word string) bool { panic("TODO: implement Contains") }

// HasPrefix reports whether any stored word starts with prefix. Every word is
// a prefix of itself, and "" is a prefix of everything.
func (t *Trie) HasPrefix(prefix string) bool { panic("TODO: implement HasPrefix") }

// Delete removes a word and reports whether it was there. It must also prune
// nodes that are no longer part of any word - a trie that only ever grows is a
// memory leak.
func (t *Trie) Delete(word string) bool { panic("TODO: implement Delete") }

// Len is the number of words.
func (t *Trie) Len() int { panic("TODO: implement Len") }

// Nodes is the number of nodes currently allocated, root included. The tests
// use it to check that Delete really prunes.
func (t *Trie) Nodes() int { panic("TODO: implement Nodes") }

// WithPrefix returns the stored words starting with prefix, in lexical order,
// at most limit of them (limit <= 0 means no limit).
func (t *Trie) WithPrefix(prefix string, limit int) []string { panic("TODO: implement WithPrefix") }

// All iterates every word in lexical order. It must honour early termination.
func (t *Trie) All() iter.Seq[string] { panic("TODO: implement All") }

// LongestPrefixOf returns the longest stored word that is a prefix of s, and
// whether there was one. This is the operation an HTTP router or an IP routing
// table performs.
//
//	insert "go", "gopher"; LongestPrefixOf("gophers") -> "gopher", true
func (t *Trie) LongestPrefixOf(s string) (string, bool) { panic("TODO: implement LongestPrefixOf") }

// Match returns the words matching a pattern in which '.' stands for exactly
// one rune. Results are in lexical order.
//
//	insert "cat", "cot", "cart"; Match("c.t") -> ["cat", "cot"]
//
// This is a depth-first search that branches at every '.', and it is the reason
// a trie beats a plain map for this job.
func (t *Trie) Match(pattern string) []string { panic("TODO: implement Match") }

// CountPrefix returns how many stored words start with prefix.
func (t *Trie) CountPrefix(prefix string) int { panic("TODO: implement CountPrefix") }
