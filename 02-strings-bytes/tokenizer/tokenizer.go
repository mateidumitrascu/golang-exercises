// Package tokenizer is about bufio.SplitFunc - the small, awkward, extremely
// useful interface behind bufio.Scanner.
//
// A SplitFunc is called with whatever bytes the Scanner currently has buffered
// and whether that is everything there will ever be. It answers:
//
//	(0, nil, nil)         "I need more data" (only legal when !atEOF)
//	(n, token, nil)       "consume n bytes, here is a token"
//	(n, nil, nil)         "consume n bytes, no token this time"
//	(0, nil, err)         "stop, this is broken"
//
// Getting the boundary cases right - a multi-byte rune split across the buffer
// edge, a final token with no trailing delimiter, an empty token - is the whole
// exercise. Your functions must work when data arrives one byte at a time.
package tokenizer

import (
	"bufio"
	"io"
)

// ScanWords is a SplitFunc whose tokens are maximal runs of letters and digits
// (unicode-aware). Everything else - spaces, punctuation, symbols - is a
// separator and never appears in a token.
//
//	"Hello, wörld! 42x" -> "Hello" "wörld" "42x"
func ScanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	panic("TODO: implement ScanWords")
}

// ScanFixed returns a SplitFunc that emits chunks of exactly n bytes. The final
// chunk may be shorter but is never empty. ScanFixed panics if n <= 0.
func ScanFixed(n int) bufio.SplitFunc {
	panic("TODO: implement ScanFixed")
}

// ScanDelimited returns a SplitFunc that splits on a single byte, dropping the
// delimiter. Consecutive delimiters produce empty tokens. A trailing delimiter
// at EOF does NOT produce a final empty token, but "a,,b" produces three
// tokens: "a", "", "b".
func ScanDelimited(delim byte) bufio.SplitFunc {
	panic("TODO: implement ScanDelimited")
}

// ScanLengthPrefixed returns a SplitFunc for a binary framing where each frame
// is a 4-byte big-endian length followed by that many payload bytes. The token
// is the payload without the header.
//
// It returns io.ErrUnexpectedEOF if the stream ends mid-frame, and refuses
// frames larger than maxFrame with an error of your choosing.
func ScanLengthPrefixed(maxFrame int) bufio.SplitFunc {
	panic("TODO: implement ScanLengthPrefixed")
}

// WordFreq counts lowercased words from r using ScanWords. It returns any read
// error from the underlying reader.
func WordFreq(r io.Reader) (map[string]int, error) {
	panic("TODO: implement WordFreq")
}
