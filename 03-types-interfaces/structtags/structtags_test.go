package structtags

import (
	"errors"
	"reflect"
	"strings"
	"testing"
)

type Address struct {
	City string `validate:"required,min=2"`
	Zip  string `validate:"min=4,max=10"`
}

type Item struct {
	SKU string `validate:"required"`
	Qty int    `validate:"min=1,max=99"`
}

type User struct {
	Name    string   `validate:"required,min=2,max=20"`
	Email   string   `validate:"required,email"`
	Age     int      `validate:"min=0,max=150"`
	Role    string   `validate:"oneof=admin user guest"`
	Tags    []string `validate:"max=3"`
	Addr    Address
	Ptr     *Address
	Items   []Item
	Ignored string `validate:"-"`
	secret  string
}

func validUser() User {
	return User{
		Name:  "Ana",
		Email: "ana@example.com",
		Age:   31,
		Role:  "admin",
		Tags:  []string{"a"},
		Addr:  Address{City: "Cluj", Zip: "400000"},
		Items: []Item{{SKU: "x", Qty: 2}},
	}
}

func TestValidateOK(t *testing.T) {
	u := validUser()
	if err := Validate(u); err != nil {
		t.Errorf("Validate = %v, want nil", err)
	}
	if err := Validate(&u); err != nil {
		t.Errorf("Validate(pointer) = %v, want nil", err)
	}
	// The nil-interface trap: a nil error must compare == nil.
	err := Validate(u)
	if err != nil {
		t.Errorf("returned a typed nil: %#v", err)
	}
}

func TestValidateFailures(t *testing.T) {
	tests := []struct {
		name      string
		mutate    func(*User)
		wantPaths []string
	}{
		{"missing name", func(u *User) { u.Name = "" }, []string{"Name"}},
		{"short name", func(u *User) { u.Name = "A" }, []string{"Name"}},
		{"bad email", func(u *User) { u.Email = "not-an-email" }, []string{"Email"}},
		{"bad email 2", func(u *User) { u.Email = "a@b" }, []string{"Email"}},
		{"bad email 3", func(u *User) { u.Email = "@example.com" }, []string{"Email"}},
		{"age too big", func(u *User) { u.Age = 200 }, []string{"Age"}},
		{"bad role", func(u *User) { u.Role = "root" }, []string{"Role"}},
		{"too many tags", func(u *User) { u.Tags = []string{"a", "b", "c", "d"} }, []string{"Tags"}},
		{"nested", func(u *User) { u.Addr.City = "" }, []string{"Addr.City"}},
		{"nested zip", func(u *User) { u.Addr.Zip = "1" }, []string{"Addr.Zip"}},
		{"through pointer", func(u *User) { u.Ptr = &Address{City: "X", Zip: "1234"} }, []string{"Ptr.City"}},
		{"slice element", func(u *User) { u.Items = []Item{{SKU: "a", Qty: 1}, {SKU: "", Qty: 0}} },
			[]string{"Items[1].SKU", "Items[1].Qty"}},
		{"several at once", func(u *User) { u.Name = ""; u.Role = "x" }, []string{"Name", "Role"}},
		{"tag - is skipped", func(u *User) { u.Ignored = "" }, nil},
		{"unexported is skipped", func(u *User) { u.secret = "" }, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := validUser()
			tt.mutate(&u)
			err := Validate(u)
			if len(tt.wantPaths) == 0 {
				if err != nil {
					t.Fatalf("Validate = %v, want nil", err)
				}
				return
			}
			var fe FieldErrors
			if !errors.As(err, &fe) {
				t.Fatalf("err = %v (%T), want FieldErrors", err, err)
			}
			if !reflect.DeepEqual(fe.Fields(), tt.wantPaths) {
				t.Errorf("failed fields = %v, want %v", fe.Fields(), tt.wantPaths)
			}
			for _, f := range fe {
				if f.Rule == "" || f.Msg == "" {
					t.Errorf("FieldError %+v has an empty Rule or Msg", f)
				}
			}
			if !strings.Contains(err.Error(), tt.wantPaths[0]) {
				t.Errorf("Error() = %q, should mention %q", err, tt.wantPaths[0])
			}
		})
	}
}

func TestValidateNilPointerNotDescended(t *testing.T) {
	u := validUser()
	u.Ptr = nil // no `required` on Ptr, so this is fine
	if err := Validate(u); err != nil {
		t.Errorf("a nil pointer field must not be descended into: %v", err)
	}
}

func TestValidateBadInput(t *testing.T) {
	for _, v := range []any{42, "hello", nil, []int{1}, (*User)(nil)} {
		if err := Validate(v); !errors.Is(err, ErrNotStruct) {
			t.Errorf("Validate(%#v) = %v, want ErrNotStruct", v, err)
		}
	}
}

func TestValidateBadRule(t *testing.T) {
	type bad1 struct {
		A string `validate:"frobnicate"`
	}
	type bad2 struct {
		A string `validate:"min=abc"`
	}
	type bad3 struct {
		A string `validate:"min"`
	}
	for _, v := range []any{bad1{}, bad2{"x"}, bad3{"x"}} {
		if err := Validate(v); !errors.Is(err, ErrBadRule) {
			t.Errorf("Validate(%T) = %v, want ErrBadRule", v, err)
		}
	}
}

type mapped struct {
	ID      int    `json:"id"`
	Name    string `json:"name,omitempty"`
	Secret  string `json:"-"`
	NoTag   bool
	Nested  Address  `json:"nested"`
	Pointer *Address `json:"pointer"`
	hidden  int
}

func TestToMap(t *testing.T) {
	in := mapped{ID: 7, Name: "x", Secret: "s", NoTag: true,
		Nested: Address{City: "Cluj", Zip: "400"}}
	got, err := ToMap(in, "json")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"id":      7,
		"name":    "x",
		"NoTag":   true,
		"nested":  map[string]any{"City": "Cluj", "Zip": "400"},
		"pointer": nil,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ToMap =\n%#v\nwant\n%#v", got, want)
	}
}

func TestToMapOmitempty(t *testing.T) {
	got, err := ToMap(mapped{ID: 0, Name: ""}, "json")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := got["name"]; ok {
		t.Error("omitempty field with a zero value should be absent")
	}
	if _, ok := got["id"]; !ok {
		t.Error("a zero field without omitempty must still be present")
	}
}

func TestToMapBadInput(t *testing.T) {
	if _, err := ToMap(42, "json"); !errors.Is(err, ErrNotStruct) {
		t.Errorf("err = %v, want ErrNotStruct", err)
	}
}
