// Package framing puts a message boundary on a byte stream.
//
// TCP gives you bytes, not messages. Every network protocol therefore invents
// framing: fixed size, delimiters, or - the usual answer - a length prefix.
// You are building the length-prefixed kind, with a checksum, plus the varint
// encoding that protobuf and friends use.
//
// Wire format for one frame, all integers big-endian:
//
//	magic    2 bytes   0xC0 0xDE
//	type     1 byte
//	length   4 bytes   payload length
//	payload  <length> bytes
//	crc32    4 bytes   IEEE checksum of the payload only
package framing

import (
	"errors"
	"io"
)

const (
	Magic0    = 0xC0
	Magic1    = 0xDE
	HeaderLen = 7 // magic + type + length
)

var (
	ErrBadMagic    = errors.New("framing: bad magic")
	ErrChecksum    = errors.New("framing: checksum mismatch")
	ErrFrameTooBig = errors.New("framing: frame too big")
)

// Frame is one message.
type Frame struct {
	Type    uint8
	Payload []byte
}

// WriteFrame writes f to w. It must issue as few Write calls as it reasonably
// can (build the whole frame, write it once) - a syscall per field is how you
// make a protocol slow.
func WriteFrame(w io.Writer, f Frame) error { panic("TODO: implement WriteFrame") }

// ReadFrame reads one frame from r.
//
//	clean end of stream before any bytes  -> io.EOF
//	stream ends part way through a frame  -> io.ErrUnexpectedEOF
//	wrong magic                           -> ErrBadMagic
//	length > maxPayload                   -> ErrFrameTooBig (and do NOT allocate it)
//	crc mismatch                          -> ErrChecksum
//
// io.ReadFull is the tool for "give me exactly this many bytes or tell me why
// not" - note how it converts a partial read into ErrUnexpectedEOF for you.
func ReadFrame(r io.Reader, maxPayload int) (Frame, error) { panic("TODO: implement ReadFrame") }

// Encoder writes frames to a buffered stream. Flush is the caller's job -
// document that, because forgetting it is the classic bug with buffered output.
type Encoder struct {
	// TODO
}

func NewEncoder(w io.Writer) *Encoder   { panic("TODO: implement NewEncoder") }
func (e *Encoder) Encode(f Frame) error { panic("TODO: implement Encode") }
func (e *Encoder) Flush() error         { panic("TODO: implement Flush") }

// Decoder reads frames from a stream.
type Decoder struct {
	// TODO
}

func NewDecoder(r io.Reader, maxPayload int) *Decoder { panic("TODO: implement NewDecoder") }

// Decode returns the next frame, or io.EOF at a clean end of stream.
func (d *Decoder) Decode() (Frame, error) { panic("TODO: implement Decode") }

// All iterates the remaining frames. Iteration stops at the first error, which
// is then available from Err. (This is the bufio.Scanner design: iterate
// cleanly, check the error afterwards.)
func (d *Decoder) All(yield func(Frame) bool) { panic("TODO: implement All") }

// Err returns the error that stopped iteration, or nil at a clean EOF.
func (d *Decoder) Err() error { panic("TODO: implement Err") }

// AppendUvarint appends x to buf in the base-128 varint encoding: seven bits of
// payload per byte, little-endian, with the high bit set on every byte except
// the last. Write it by hand; do not call encoding/binary.
//
//	1     -> 0x01
//	127   -> 0x7F
//	128   -> 0x80 0x01
//	300   -> 0xAC 0x02
func AppendUvarint(buf []byte, x uint64) []byte { panic("TODO: implement AppendUvarint") }

// ConsumeUvarint decodes a varint from the front of buf and returns the value
// and the number of bytes used. It returns n == 0 if buf is too short, and
// n < 0 if the encoding overflows a uint64.
func ConsumeUvarint(buf []byte) (x uint64, n int) { panic("TODO: implement ConsumeUvarint") }

// AppendString appends a length-prefixed (varint) string.
func AppendString(buf []byte, s string) []byte { panic("TODO: implement AppendString") }

// ConsumeString is the inverse; ok is false if buf is truncated.
func ConsumeString(buf []byte) (s string, n int, ok bool) { panic("TODO: implement ConsumeString") }
