// Package readers is about io.Reader and io.Writer - two of the best-designed
// interfaces in any standard library, and the source of a specific family of
// bugs.
//
// The Reader contract, in full, because almost everyone gets it wrong:
//
//	Read(p) (n int, err error)
//
//	- It may return 0 < n < len(p) and nil. A short read is NOT an error and
//	  NOT the end of the stream. Loop, or use io.ReadFull.
//	- It may return n > 0 AND err == io.EOF in the same call. You must process
//	  those n bytes before you stop.
//	- (0, nil) is allowed but discouraged; treat it as "try again", never as EOF.
//	- Callers must not use p[n:] - implementations may not have touched it.
//
// The Writer contract is stricter: Write must write all of p or return an
// error explaining why it did not.
package readers

import "io"

// CountingWriter counts the bytes written through it. The zero value is not
// usable; use NewCountingWriter. It must pass short writes and errors through
// unchanged, and count only what actually landed.
type CountingWriter struct {
	// TODO
}

func NewCountingWriter(w io.Writer) *CountingWriter   { panic("TODO: implement NewCountingWriter") }
func (c *CountingWriter) Write(p []byte) (int, error) { panic("TODO: implement Write") }
func (c *CountingWriter) N() int64                    { panic("TODO: implement N") }

// UpperReader wraps r, converting ASCII lowercase to uppercase as bytes flow
// through. It must not buffer: whatever the underlying reader returns, it
// transforms that much and returns immediately.
func UpperReader(r io.Reader) io.Reader { panic("TODO: implement UpperReader") }

// ExactReader reads exactly n bytes from r and then reports io.EOF. If the
// underlying reader runs out first, it returns io.ErrUnexpectedEOF - the
// distinction that tells "this stream is finished" apart from "this stream is
// truncated".
func ExactReader(r io.Reader, n int64) io.Reader { panic("TODO: implement ExactReader") }

// TeeReader returns a reader that writes to w everything it reads from r.
// Write errors from w are returned as read errors. Your own io.TeeReader.
func TeeReader(r io.Reader, w io.Writer) io.Reader { panic("TODO: implement TeeReader") }

// MultiWriter duplicates writes to every writer, in order. If one fails, it
// stops and returns that error. If one accepts fewer bytes than it was given,
// that is io.ErrShortWrite.
func MultiWriter(ws ...io.Writer) io.Writer { panic("TODO: implement MultiWriter") }

// ReadAll reads r to EOF. Your own io.ReadAll: start with a small buffer, grow
// it geometrically, and remember that the final Read can return data together
// with io.EOF.
func ReadAll(r io.Reader) ([]byte, error) { panic("TODO: implement ReadAll") }

// Copy copies src to dst using buf (allocating a 32 KiB one if buf is nil) and
// returns the number of bytes copied. Your own io.CopyBuffer:
//
//	read n bytes; if n > 0, write them and check for a short write;
//	only then look at the read error; io.EOF means success, not failure.
func Copy(dst io.Writer, src io.Reader, buf []byte) (int64, error) {
	panic("TODO: implement Copy")
}

// CopyProgress is Copy with a callback invoked after each successful write with
// the running total. It is how you draw a progress bar.
func CopyProgress(dst io.Writer, src io.Reader, progress func(total int64)) (int64, error) {
	panic("TODO: implement CopyProgress")
}

// LineWriter wraps w so that data is only passed through one whole line at a
// time (each ending in "\n"). Partial lines are held in an internal buffer
// until the newline arrives. Flush writes out whatever is left, adding no
// newline. This is what makes concurrent writers to one log file not interleave
// mid-line.
type LineWriter struct {
	// TODO
}

func NewLineWriter(w io.Writer) *LineWriter       { panic("TODO: implement NewLineWriter") }
func (l *LineWriter) Write(p []byte) (int, error) { panic("TODO: implement LineWriter.Write") }
func (l *LineWriter) Flush() error                { panic("TODO: implement Flush") }
