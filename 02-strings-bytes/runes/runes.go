// Package runes: strings are bytes, ranging over them yields runes, and
// indexing them yields neither.
//
// Rule of the module: you may use the unicode package, but NOT unicode/utf8.
// You are going to write the decoder yourself.
package runes

// FirstRune decodes the first UTF-8 encoded rune in s and returns it with the
// number of bytes it occupied. It is your own utf8.DecodeRuneInString.
//
// For an empty string it returns (RuneError, 0). For any invalid encoding it
// returns (RuneError, 1). "Invalid" means all of:
//
//	a leading byte that is a continuation byte (10xxxxxx) or 0xFE/0xFF
//	a truncated sequence (not enough continuation bytes present)
//	an overlong encoding (0xC0 0x80 encodes NUL in two bytes - forbidden)
//	a surrogate half, U+D800..U+DFFF
//	a value above U+10FFFF
//
// The encoding, for reference:
//
//	0xxxxxxx                                        U+0000..U+007F
//	110xxxxx 10xxxxxx                               U+0080..U+07FF
//	1110xxxx 10xxxxxx 10xxxxxx                      U+0800..U+FFFF
//	11110xxx 10xxxxxx 10xxxxxx 10xxxxxx             U+10000..U+10FFFF
const RuneError = '�'

func FirstRune(s string) (r rune, size int) {
	panic("TODO: implement FirstRune")
}

// EncodeRune writes the UTF-8 encoding of r into buf and returns the number of
// bytes written. Invalid runes (negative, surrogates, > U+10FFFF) are encoded
// as RuneError. EncodeRune panics if buf is too small.
func EncodeRune(buf []byte, r rune) int {
	panic("TODO: implement EncodeRune")
}

// CountRunes returns the number of runes in s, counting each invalid byte as
// one rune. Do it in one pass with no allocations.
func CountRunes(s string) int {
	panic("TODO: implement CountRunes")
}

// Reverse reverses s rune by rune. Invalid bytes are preserved as-is in the
// reversed output.
func Reverse(s string) string {
	panic("TODO: implement Reverse")
}

// ReverseClusters reverses s but keeps each base rune together with the
// combining marks that follow it (unicode.Mn), so that "café" reverses to
// "éfac" and not to "́efac".
//
// This is a poor man's grapheme clustering - real clustering also has to deal
// with emoji ZWJ sequences and regional indicators. Handling Mn is enough here.
func ReverseClusters(s string) string {
	panic("TODO: implement ReverseClusters")
}

// TruncateRunes shortens s to at most n runes. If it had to cut anything, the
// result ends with the single rune U+2026 (…) and is still at most n runes long.
// TruncateRunes(s, 0) is "". Negative n panics.
func TruncateRunes(s string, n int) string {
	panic("TODO: implement TruncateRunes")
}

// ByteToRuneIndex converts a byte offset into the index of the rune that
// contains it. ByteToRuneIndex("héllo", 2) is 1, because byte 2 is the second
// byte of 'é'. It returns -1 if b is out of range.
func ByteToRuneIndex(s string, b int) int {
	panic("TODO: implement ByteToRuneIndex")
}
