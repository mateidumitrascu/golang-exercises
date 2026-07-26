// Package jsonx covers the encoding/json you need once the data stops being
// simple: custom marshalers, values that arrive in more than one shape,
// polymorphic documents, and streams too big to hold in memory.
package jsonx

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

// Temperature marshals as a JSON string with a unit suffix: 21.5 -> "21.5C".
// It unmarshals from that same form, and rejects anything else with a helpful
// error.
//
// Note the receivers: MarshalJSON on the value type, UnmarshalJSON on the
// pointer type (it has to mutate). Getting that wrong means your marshaler is
// silently ignored for values - a genuinely common bug.
type Temperature float64

func (t Temperature) MarshalJSON() ([]byte, error) { panic("TODO: implement Temperature.MarshalJSON") }
func (t *Temperature) UnmarshalJSON(b []byte) error {
	panic("TODO: implement Temperature.UnmarshalJSON")
}

// Date marshals as "2006-01-02" (no time, no zone).
type Date time.Time

func (d Date) MarshalJSON() ([]byte, error)  { panic("TODO: implement Date.MarshalJSON") }
func (d *Date) UnmarshalJSON(b []byte) error { panic("TODO: implement Date.UnmarshalJSON") }
func (d Date) String() string                { panic("TODO: implement Date.String") }

// Tags accepts either a single string or an array of strings, because that is
// what real-world APIs do:
//
//	{"tags": "go"}            -> Tags{"go"}
//	{"tags": ["go", "json"]}  -> Tags{"go", "json"}
//	{"tags": null}            -> nil
//
// It always marshals as an array.
type Tags []string

func (t *Tags) UnmarshalJSON(b []byte) error { panic("TODO: implement Tags.UnmarshalJSON") }

// Config is a settings struct. Fill in the tags so that:
//
//	Name     is "name" and is always present
//	Port     is "port" and is omitted when zero
//	Debug    is "debug"
//	Password is never marshalled at all
//	Timeout  is "timeout_seconds" and holds a duration in seconds
//	Extra    collects nothing - it is skipped entirely by the JSON package
type Config struct {
	Name     string
	Port     int
	Debug    bool
	Password string
	Timeout  time.Duration
	Extra    map[string]string
}

// ParseConfig decodes strict JSON: an unknown field is an error, not something
// to ignore silently. (json.Decoder.DisallowUnknownFields.)
func ParseConfig(data []byte) (*Config, error) { panic("TODO: implement ParseConfig") }

// Event is a polymorphic message. Documents look like:
//
//	{"type": "click", "x": 10, "y": 20}
//	{"type": "keypress", "key": "a"}
//	{"type": "scroll", "delta": -3}
type Event interface{ Kind() string }

type Click struct{ X, Y int }
type KeyPress struct{ Key string }
type Scroll struct{ Delta int }

func (Click) Kind() string    { return "click" }
func (KeyPress) Kind() string { return "keypress" }
func (Scroll) Kind() string   { return "scroll" }

var ErrUnknownEvent = errors.New("unknown event type")

// DecodeEvent decodes one event. The trick: unmarshal into a struct that has
// only the discriminator plus a json.RawMessage of the whole object (or decode
// twice), then dispatch on the type and unmarshal properly.
func DecodeEvent(data []byte) (Event, error) { panic("TODO: implement DecodeEvent") }

// EncodeEvent writes an event back out in the same shape, with the type field
// first... except JSON object key order does not matter, so just make sure the
// "type" field is present and the payload fields are at the top level, not
// nested.
func EncodeEvent(e Event) ([]byte, error) { panic("TODO: implement EncodeEvent") }

// DecodeEventStream reads a stream of newline-or-whitespace separated JSON
// objects (not an array) and returns them all. json.Decoder does this natively:
// call Decode in a loop until io.EOF.
func DecodeEventStream(r io.Reader) ([]Event, error) { panic("TODO: implement DecodeEventStream") }

// SumMeasurements streams a potentially enormous JSON array of objects like
//
//	[{"value": 1.5}, {"value": 2.5}, ...]
//
// and returns the sum of the "value" fields WITHOUT holding the whole array in
// memory. Use json.Decoder: Token() to read the opening '[', More() to loop,
// Decode() to read one element at a time, Token() again for the closing ']'.
//
// It returns an error if the document is not an array of objects.
func SumMeasurements(r io.Reader) (float64, error) { panic("TODO: implement SumMeasurements") }

// Flatten turns nested JSON into dotted keys:
//
//	{"a": {"b": 1}, "c": [10, 20]}  ->  {"a.b": 1, "c.0": 10, "c.1": 20}
//
// Values keep their decoded Go types (float64, string, bool, nil).
func Flatten(data []byte) (map[string]any, error) { panic("TODO: implement Flatten") }

// Compact rewrites JSON with no insignificant whitespace, preserving key order,
// using json.Decoder's token stream. (json.Compact exists; write your own.)
func Compact(w io.Writer, data []byte) error { panic("TODO: implement Compact") }

var _ = json.Marshal // remove once you use the package
