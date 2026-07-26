package csvlite

import (
	"errors"
	"reflect"
	"testing"
)

func TestParseRecord(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []string
	}{
		{"plain", "a,b,c", []string{"a", "b", "c"}},
		{"empty fields", "a,,c", []string{"a", "", "c"}},
		{"empty input", "", []string{""}},
		{"trailing comma", "a,", []string{"a", ""}},
		{"quoted", `a,"b,c",d`, []string{"a", "b,c", "d"}},
		{"escaped quotes", `"""quoted"""`, []string{`"quoted"`}},
		{"quote in middle", `a,""" d"`, []string{"a", `" d`}},
		{"spaces are data", ` a , b `, []string{" a ", " b "}},
		{"empty quoted", `a,"",b`, []string{"a", "", "b"}},
		{"carriage return dropped", "a,b\r", []string{"a", "b"}},
		{"quoted newline", "\"a\nb\"", []string{"a\nb"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRecord(tt.in)
			if err != nil {
				t.Fatalf("ParseRecord(%q) failed: %v", tt.in, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRecord(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestParseRecordErrors(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		column int
	}{
		{"unterminated quote", `"abc`, 5},
		{"quote in unquoted field", `ab"c`, 3},
		{"junk after closing quote", `"ab"c`, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseRecord(tt.in)
			if err == nil {
				t.Fatalf("ParseRecord(%q) should have failed", tt.in)
			}
			var pe *ParseError
			if !errors.As(err, &pe) {
				t.Fatalf("err is %T, want *ParseError", err)
			}
			if pe.Line != 1 {
				t.Errorf("Line = %d, want 1", pe.Line)
			}
			if pe.Column != tt.column {
				t.Errorf("Column = %d, want %d (err: %v)", pe.Column, tt.column, err)
			}
			if pe.Msg == "" {
				t.Error("Msg is empty; say what went wrong")
			}
		})
	}
}

func TestParseAll(t *testing.T) {
	in := "name,age\n\"Doe, John\",42\n\"multi\nline\",1\n"
	want := [][]string{
		{"name", "age"},
		{"Doe, John", "42"},
		{"multi\nline", "1"},
	}
	got, err := ParseAll(in)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ParseAll =\n%q\nwant\n%q", got, want)
	}
}

func TestParseAllEdgeCases(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want [][]string
	}{
		{"crlf", "a,b\r\nc,d\r\n", [][]string{{"a", "b"}, {"c", "d"}}},
		{"no trailing newline", "a\nb", [][]string{{"a"}, {"b"}}},
		{"empty document", "", [][]string{}},
		{"blank line is a record", "a\n\nb", [][]string{{"a"}, {""}, {"b"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseAll(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("ParseAll(%q) = %q, want %q", tt.in, got, tt.want)
			}
			for i := range got {
				if !reflect.DeepEqual(got[i], tt.want[i]) {
					t.Errorf("record %d = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestParseAllReportsLineNumbers(t *testing.T) {
	_, err := ParseAll("a,b\nc,\"d\ne,f")
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("err = %v (%T), want *ParseError", err, err)
	}
	if pe.Line != 2 {
		t.Errorf("Line = %d, want 2 (the quote opens on line 2)", pe.Line)
	}
}

func TestValidate(t *testing.T) {
	ok := [][]string{{"a", "b"}, {"c", "d"}}
	if err := Validate(ok); err != nil {
		t.Errorf("rectangular table reported an error: %v", err)
	}
	bad := [][]string{{"a", "b"}, {"c"}, {"d", "e"}}
	err := Validate(bad)
	var pe *ParseError
	if !errors.As(err, &pe) {
		t.Fatalf("err = %v, want *ParseError", err)
	}
	if pe.Line != 2 {
		t.Errorf("Line = %d, want 2", pe.Line)
	}
	if err := Validate(nil); err != nil {
		t.Errorf("empty table: %v", err)
	}
}

func TestFormatField(t *testing.T) {
	tests := []struct{ in, want string }{
		{"plain", "plain"},
		{"", ""},
		{"has,comma", `"has,comma"`},
		{`has"quote`, `"has""quote"`},
		{"has\nnewline", "\"has\nnewline\""},
		{" spaces ", " spaces "},
	}
	for _, tt := range tests {
		if got := FormatField(tt.in); got != tt.want {
			t.Errorf("FormatField(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	records := [][]string{
		{"name", "quote", "note"},
		{"Doe, John", `he said "hi"`, "line1\nline2"},
		{"", " ", "plain"},
	}
	text := FormatAll(records)
	got, err := ParseAll(text)
	if err != nil {
		t.Fatalf("re-parsing our own output failed: %v\noutput was:\n%s", err, text)
	}
	if !reflect.DeepEqual(got, records) {
		t.Errorf("round trip changed the data:\ngot  %q\nwant %q", got, records)
	}
}

func FuzzRoundTrip(f *testing.F) {
	f.Add("a,b", "c")
	f.Add(`"x"`, "\n")
	f.Fuzz(func(t *testing.T, a, b string) {
		records := [][]string{{a, b}}
		got, err := ParseAll(FormatAll(records))
		if err != nil {
			t.Fatalf("formatting then parsing %q failed: %v", records, err)
		}
		if len(got) != 1 || !reflect.DeepEqual(got[0], records[0]) {
			t.Fatalf("round trip: got %q, want %q", got, records)
		}
	})
}
