# 02 — Strings, bytes, runes

A Go string is an immutable slice of bytes with no declared encoding. Indexing
gives bytes, ranging gives runes, and `len` gives neither of the things people
expect. This module makes that concrete.

| # | Exercise | Difficulty | What it drills |
|---|----------|-----------|----------------|
| 1 | `runes` | ★★★ | Writing a UTF-8 decoder/encoder by hand; differential fuzzing |
| 2 | `wordwrap` | ★★☆ | Rune-width layout, `strings.Builder`, whitespace edge cases |
| 3 | `tokenizer` | ★★★ | `bufio.SplitFunc`: partial buffers, EOF, framing |
| 4 | `csvlite` | ★★★ | A byte-level state machine with positional errors, round-trip fuzzing |
| 5 | `edit` | ★★★ | Two-row dynamic programming, backtracing an edit script |

## Ideas worth internalising here

**`for i, r := range s`** decodes UTF-8 as it goes: `i` jumps by the rune's byte
width, and invalid bytes yield `U+FFFD` with width 1. `s[i]` is a byte. `[]rune(s)`
allocates a whole new slice — in a hot path, don't.

**`strings.Builder` beats `+=`** because `+=` reallocates and copies on every
step. `Builder.Grow(n)` when you can estimate the result makes it one allocation.

**A `SplitFunc` is called repeatedly with more data each time.** It must be
prepared for a token — or a single rune — to be cut in half by the buffer
boundary, and must distinguish "need more" from "done".

**Fuzzing is cheap here.** Two of these exercises have differential fuzz targets:
compare your implementation against the standard library, or against your own
inverse function. `go test -fuzz FuzzName ./path/` and let it find the case you
didn't think of.

## Stretch goals

- Make `csvlite` streaming: a `Reader` with a `Read() ([]string, error)` method
  over an `io.Reader`, so a 10 GB file doesn't have to fit in memory.
- Add `Distance` with transpositions (Damerau–Levenshtein) and check that
  `Distance("form","from")` drops from 2 to 1.
- Write a `strings.Replacer`-style multi-pattern replacer using a trie
  (you'll build the trie for real in module 09).
