package jsonx

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestTemperature(t *testing.T) {
	b, err := json.Marshal(Temperature(21.5))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"21.5C"` {
		t.Errorf("= %s, want \"21.5C\"", b)
	}
	// Inside a struct, by value - this is where a pointer-receiver marshaler
	// would silently not fire.
	s := struct {
		T Temperature `json:"t"`
	}{-3}
	b, _ = json.Marshal(s)
	if string(b) != `{"t":"-3C"}` {
		t.Errorf("= %s", b)
	}

	var got Temperature
	if err := json.Unmarshal([]byte(`"36.6C"`), &got); err != nil {
		t.Fatal(err)
	}
	if got != 36.6 {
		t.Errorf("= %v, want 36.6", got)
	}
	for _, bad := range []string{`"36.6"`, `36.6`, `"abcC"`, `"C"`, `""`} {
		if err := json.Unmarshal([]byte(bad), &got); err == nil {
			t.Errorf("Unmarshal(%s) should have failed", bad)
		}
	}
}

func TestDate(t *testing.T) {
	d := Date(time.Date(2026, 7, 25, 15, 4, 5, 0, time.UTC))
	b, err := json.Marshal(d)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"2026-07-25"` {
		t.Errorf("= %s", b)
	}
	var got Date
	if err := json.Unmarshal([]byte(`"2001-02-03"`), &got); err != nil {
		t.Fatal(err)
	}
	if got.String() != "2001-02-03" {
		t.Errorf("String = %q", got.String())
	}
	if err := json.Unmarshal([]byte(`"03/02/2001"`), &got); err == nil {
		t.Error("a bad date format should fail")
	}
}

func TestTags(t *testing.T) {
	tests := []struct {
		in   string
		want Tags
	}{
		{`"go"`, Tags{"go"}},
		{`["go","json"]`, Tags{"go", "json"}},
		{`[]`, Tags{}},
		{`null`, nil},
	}
	for _, tt := range tests {
		var got Tags
		if err := json.Unmarshal([]byte(tt.in), &got); err != nil {
			t.Fatalf("Unmarshal(%s): %v", tt.in, err)
		}
		if len(got) != len(tt.want) {
			t.Errorf("Unmarshal(%s) = %v, want %v", tt.in, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("Unmarshal(%s) = %v, want %v", tt.in, got, tt.want)
			}
		}
	}
	if err := json.Unmarshal([]byte(`42`), new(Tags)); err == nil {
		t.Error("a number should not decode into Tags")
	}
	b, _ := json.Marshal(Tags{"a"})
	if string(b) != `["a"]` {
		t.Errorf("marshal = %s, want an array", b)
	}
}

func TestConfigTags(t *testing.T) {
	c := Config{Name: "svc", Debug: true, Password: "hunter2", Timeout: 30 * time.Second,
		Extra: map[string]string{"a": "b"}}
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	s := string(b)
	for _, want := range []string{`"name":"svc"`, `"debug":true`, `"timeout_seconds":30`} {
		if !strings.Contains(s, want) {
			t.Errorf("marshal = %s, missing %s", s, want)
		}
	}
	for _, unwanted := range []string{"hunter2", "Password", "port", "Extra"} {
		if strings.Contains(s, unwanted) {
			t.Errorf("marshal = %s, should not contain %q", s, unwanted)
		}
	}
	c.Port = 8080
	b, _ = json.Marshal(c)
	if !strings.Contains(string(b), `"port":8080`) {
		t.Errorf("marshal = %s, want the port when it is set", b)
	}
}

func TestParseConfig(t *testing.T) {
	got, err := ParseConfig([]byte(`{"name":"svc","port":80,"timeout_seconds":5}`))
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "svc" || got.Port != 80 || got.Timeout != 5*time.Second {
		t.Errorf("= %+v", got)
	}
	if _, err := ParseConfig([]byte(`{"name":"x","nope":1}`)); err == nil {
		t.Error("an unknown field must be rejected")
	}
	if _, err := ParseConfig([]byte(`{`)); err == nil {
		t.Error("malformed JSON must fail")
	}
}

func TestDecodeEvent(t *testing.T) {
	tests := []struct {
		in   string
		want Event
	}{
		{`{"type":"click","x":10,"y":20}`, Click{10, 20}},
		{`{"type":"keypress","key":"a"}`, KeyPress{"a"}},
		{`{"type":"scroll","delta":-3}`, Scroll{-3}},
	}
	for _, tt := range tests {
		got, err := DecodeEvent([]byte(tt.in))
		if err != nil {
			t.Fatalf("DecodeEvent(%s): %v", tt.in, err)
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("DecodeEvent(%s) = %#v, want %#v", tt.in, got, tt.want)
		}
	}
	if _, err := DecodeEvent([]byte(`{"type":"explode"}`)); !errors.Is(err, ErrUnknownEvent) {
		t.Errorf("err = %v, want ErrUnknownEvent", err)
	}
	if _, err := DecodeEvent([]byte(`{}`)); err == nil {
		t.Error("a missing type must be an error")
	}
}

func TestEncodeEventRoundTrip(t *testing.T) {
	for _, e := range []Event{Click{1, 2}, KeyPress{"z"}, Scroll{5}} {
		b, err := EncodeEvent(e)
		if err != nil {
			t.Fatal(err)
		}
		var probe map[string]any
		json.Unmarshal(b, &probe)
		if probe["type"] != e.Kind() {
			t.Errorf("encoded %s, missing or wrong type field: %s", e.Kind(), b)
		}
		got, err := DecodeEvent(b)
		if err != nil {
			t.Fatalf("re-decoding %s: %v", b, err)
		}
		if !reflect.DeepEqual(got, e) {
			t.Errorf("round trip: %#v -> %s -> %#v", e, b, got)
		}
	}
}

func TestDecodeEventStream(t *testing.T) {
	stream := `{"type":"click","x":1,"y":2}
	{"type":"scroll","delta":9}{"type":"keypress","key":"q"}`
	got, err := DecodeEventStream(strings.NewReader(stream))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 3 {
		t.Fatalf("got %d events, want 3", len(got))
	}
	if got[2].Kind() != "keypress" {
		t.Errorf("last event = %#v", got[2])
	}
	if got, err := DecodeEventStream(strings.NewReader("")); err != nil || len(got) != 0 {
		t.Errorf("empty stream = %v, %v", got, err)
	}
}

func TestSumMeasurements(t *testing.T) {
	doc := `[{"value": 1.5}, {"value": 2.5}, {"value": -1}]`
	got, err := SumMeasurements(strings.NewReader(doc))
	if err != nil {
		t.Fatal(err)
	}
	if got != 3 {
		t.Errorf("= %v, want 3", got)
	}
	if got, err := SumMeasurements(strings.NewReader("[]")); err != nil || got != 0 {
		t.Errorf("empty array = %v, %v", got, err)
	}
	for _, bad := range []string{`{"value":1}`, `[1,2]`, `[`, `"nope"`} {
		if _, err := SumMeasurements(strings.NewReader(bad)); err == nil {
			t.Errorf("SumMeasurements(%s) should have failed", bad)
		}
	}
}

// TestSumMeasurementsIsStreaming feeds in a ~2 MB document and measures how
// much memory the call allocates in total. Reading the whole thing and
// unmarshalling into a slice costs several megabytes; decoding one element at a
// time costs a fixed buffer.
func TestSumMeasurementsIsStreaming(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("[")
	for i := range 200000 {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(`{"value":1}`)
	}
	buf.WriteString("]")
	data := buf.Bytes()

	var before, after runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&before)
	got, err := SumMeasurements(bytes.NewReader(data))
	runtime.ReadMemStats(&after)
	if err != nil {
		t.Fatal(err)
	}
	if got != 200000 {
		t.Errorf("= %v, want 200000", got)
	}
	used := after.TotalAlloc - before.TotalAlloc
	if used > 1<<21 { // 2 MiB, against a 2.2 MB input
		t.Errorf("allocated %d bytes for a %d byte document; decode one element "+
			"at a time instead of unmarshalling the whole array", used, len(data))
	}
}

func TestFlatten(t *testing.T) {
	got, err := Flatten([]byte(`{"a":{"b":1,"c":{"d":"x"}},"e":[10,20],"f":null,"g":true}`))
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"a.b":   float64(1),
		"a.c.d": "x",
		"e.0":   float64(10),
		"e.1":   float64(20),
		"f":     nil,
		"g":     true,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Flatten =\n%#v\nwant\n%#v", got, want)
	}
	if got, err := Flatten([]byte(`{}`)); err != nil || len(got) != 0 {
		t.Errorf("= %v, %v", got, err)
	}
	if _, err := Flatten([]byte(`nope`)); err == nil {
		t.Error("invalid JSON must fail")
	}
}

func TestCompact(t *testing.T) {
	in := `{ "b" : 1 , "a" : [ 1 , 2 ] , "s" : "keep  spaces" }`
	var buf bytes.Buffer
	if err := Compact(&buf, []byte(in)); err != nil {
		t.Fatal(err)
	}
	want := `{"b":1,"a":[1,2],"s":"keep  spaces"}`
	if buf.String() != want {
		t.Errorf("= %s\nwant %s", buf.String(), want)
	}
	if err := Compact(&bytes.Buffer{}, []byte(`{`)); err == nil {
		t.Error("invalid JSON must fail")
	}
}
