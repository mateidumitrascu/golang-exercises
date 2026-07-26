// Package csvlite is a hand-written CSV parser: a state machine over bytes.
//
// Do not use encoding/csv. The point is to write the state machine, get the
// quoting rules right, and produce errors that tell the user where the problem
// is - the difference between a toy parser and a usable one.
//
// The dialect (a subset of RFC 4180):
//
//	fields are separated by commas, records by "\n" (a trailing "\r" is dropped)
//	a field may be wrapped in double quotes
//	inside quotes, "" means one literal quote, and commas and newlines are data
//	a quote may only appear in an unquoted field if... it may not. That is an error.
//	leading and trailing spaces are significant and are never trimmed
package csvlite

import "fmt"

// ParseError says what went wrong and where.
type ParseError struct {
	Line   int // 1-based
	Column int // 1-based, in bytes
	Msg    string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("csvlite: line %d, column %d: %s", e.Line, e.Column, e.Msg)
}

// ParseRecord parses exactly one record from a string that contains no
// unquoted newline. It returns a *ParseError with Line == 1 on failure.
//
//	`a,"b,c",""" d"` -> ["a", "b,c", `" d`]
//
// An empty input is one record containing one empty field.
func ParseRecord(line string) ([]string, error) {
	panic("TODO: implement ParseRecord")
}

// ParseAll parses a whole document. Quoted fields may contain newlines, which
// do not end the record. A trailing newline at the end of the document does not
// create a final empty record.
//
// If records have differing field counts, that is NOT an error here - report it
// with the separate Validate function.
func ParseAll(s string) ([][]string, error) {
	panic("TODO: implement ParseAll")
}

// Validate reports the first record whose field count differs from the first
// record's, as a *ParseError. It returns nil if the table is rectangular.
func Validate(records [][]string) error {
	panic("TODO: implement Validate")
}

// FormatField quotes a field if and only if it has to be: if it contains a
// comma, a quote, a newline, or a carriage return. Quotes inside are doubled.
func FormatField(s string) string {
	panic("TODO: implement FormatField")
}

// FormatAll renders records back to CSV text, one record per line, each line
// ending in "\n".
func FormatAll(records [][]string) string {
	panic("TODO: implement FormatAll")
}
