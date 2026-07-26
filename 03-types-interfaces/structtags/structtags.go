// Package structtags is reflection with a purpose: a struct validator driven by
// tags, and a struct-to-map converter. Every config loader, ORM and JSON codec
// in Go is a variation on this code.
//
// Reflection rules you will meet here: you cannot read unexported fields, Kind
// is not Type, and reflect.Value.IsZero is your friend.
package structtags

import (
	"errors"
	"fmt"
)

// FieldError is one failed rule on one field.
type FieldError struct {
	Path string // "Name", or "Address.City", or "Items[2].SKU"
	Rule string // "required", "min", "max", "oneof", "email"
	Msg  string
}

func (e FieldError) Error() string { return fmt.Sprintf("%s: %s", e.Path, e.Msg) }

// FieldErrors is the collection of everything wrong with a value, in field
// declaration order (depth first).
type FieldErrors []FieldError

// Error joins the field errors with "; ". An empty FieldErrors must never be
// returned as a non-nil error - see the note on Validate.
func (fe FieldErrors) Error() string { panic("TODO: implement FieldErrors.Error") }

// Fields returns the paths that failed, in order.
func (fe FieldErrors) Fields() []string { panic("TODO: implement FieldErrors.Fields") }

var (
	ErrNotStruct = errors.New("not a struct")
	ErrBadRule   = errors.New("bad validation rule")
)

// Validate checks v against its `validate` struct tags and returns FieldErrors
// (as an error) or nil. It must return a nil error - not an empty non-nil
// FieldErrors - when everything is fine. (Remember module 03's nilcheck: this
// is exactly that trap.)
//
// v may be a struct or a non-nil pointer to one; anything else is ErrNotStruct.
//
// Supported rules, comma separated in the tag:
//
//	required    the field must not be its zero value; for slices and maps,
//	            len > 0; for pointers, non-nil
//	min=N       strings: at least N runes; numbers: value >= N;
//	            slices/maps: at least N elements
//	max=N       the same, upper bound
//	oneof=a b c the field's string form must be one of the space-separated words
//	email       must contain exactly one '@', a non-empty part before it, and a
//	            part after it containing a '.' that is neither first nor last
//
// A tag of "-" skips the field. Unexported fields are always skipped. An
// unknown rule, or a malformed one, is returned as an error wrapping ErrBadRule
// (not as a FieldError - it is a bug in the program, not in the data).
//
// Nesting: a struct field, or a non-nil pointer to a struct, is validated
// recursively with its path prefixed ("Address.City"). A slice of structs is
// validated element by element ("Items[0].SKU"). A nil pointer to a struct is
// not descended into (but `required` still applies to the pointer itself).
func Validate(v any) error {
	panic("TODO: implement Validate")
}

// ToMap converts a struct into a map using the given tag name for keys
// (e.g. "json"). Rules:
//
//	no tag             -> the field name is the key
//	tag "-"            -> skipped
//	tag "name,omitempty" -> key is "name", and the field is skipped if it is
//	                        the zero value
//	nested struct      -> a nested map[string]any
//	pointer to struct  -> nested map, or nil if the pointer is nil
//	everything else    -> the value as-is
//
// Unexported fields are skipped. Non-struct input is ErrNotStruct.
func ToMap(v any, tagName string) (map[string]any, error) {
	panic("TODO: implement ToMap")
}
