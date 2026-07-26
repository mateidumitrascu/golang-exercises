// Package wordwrap does terminal-style text layout. Width is measured in RUNES,
// not bytes - an accented character is one column wide, not two.
package wordwrap

// Wrap greedily lays text out into lines of at most width runes.
//
//   - Input is split into words on any run of whitespace (space, tab, newline).
//   - Words are never broken, even if a single word exceeds width: it gets a
//     line to itself.
//   - Lines are joined with "\n". There is no trailing newline, no trailing
//     space, and no leading space.
//   - Wrap panics if width <= 0.
func Wrap(text string, width int) string {
	panic("TODO: implement Wrap")
}

// WrapHard is Wrap, except that a word longer than width is broken across as
// many lines as it needs (at rune boundaries).
func WrapHard(text string, width int) string {
	panic("TODO: implement WrapHard")
}

// WrapParagraphs preserves paragraph structure: a blank line (two or more
// consecutive newlines, possibly with spaces between them) separates
// paragraphs, and each paragraph is wrapped independently. Paragraphs are
// rejoined with exactly one blank line between them.
func WrapParagraphs(text string, width int) string {
	panic("TODO: implement WrapParagraphs")
}

// Indent prefixes every non-empty line of s with prefix. Empty lines are left
// completely empty (no trailing whitespace).
func Indent(s, prefix string) string {
	panic("TODO: implement Indent")
}

// Dedent removes the longest common leading whitespace prefix from every
// non-empty line. Lines that are entirely whitespace are ignored when computing
// the common prefix and are emitted as "".
//
// This is what you want for cleaning up a raw string literal in a test.
func Dedent(s string) string {
	panic("TODO: implement Dedent")
}
