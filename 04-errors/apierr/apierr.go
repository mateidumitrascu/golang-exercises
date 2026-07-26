// Package apierr designs an error type for a service: something you can attach
// an operation, a category and a field to, that composes across layers, and
// that can be safely shown to a user or turned into a status code.
//
// The variadic constructor is Rob Pike's design from Upspin. It reads well at
// call sites and it is a genuinely good use of `any`:
//
//	return E(Op("user.Create"), KindInvalid, Field("email"), "must not be empty")
package apierr

// Kind is a coarse category. Callers switch on this instead of on string
// matching, and transports (HTTP, gRPC) map it to their own codes.
type Kind uint8

const (
	KindOther Kind = iota
	KindInvalid
	KindNotFound
	KindConflict
	KindPermission
	KindInternal
	KindUnavailable
)

// String returns the lowercase name: "other", "invalid", "not found",
// "conflict", "permission", "internal", "unavailable".
func (k Kind) String() string { panic("TODO: implement Kind.String") }

// Op is the logical operation that failed, e.g. "user.Create".
type Op string

// Field names the input that was wrong, for validation errors.
type Field string

// Error is the package's error type. Every field is optional.
type Error struct {
	Op    Op
	Kind  Kind
	Field Field
	Msg   string
	Err   error
}

// E builds an *Error from whatever you pass it, in any order:
//
//	Op     -> Op          Kind  -> Kind
//	Field  -> Field       string -> Msg
//	error  -> Err
//
// Later arguments of the same type overwrite earlier ones. Any other argument
// type panics - that is a bug in the caller, caught at the first test run.
// E() with no arguments returns an *Error with all fields zero.
func E(args ...any) error { panic("TODO: implement E") }

// Error joins the non-empty parts with ": ", in this order:
//
//	Op, Field, Kind (only when it is not KindOther), Msg, Err
//
// So E(Op("user.Create"), KindInvalid, Field("email"), "must not be empty")
// prints as:
//
//	user.Create: email: invalid: must not be empty
//
// and wrapping that in E(Op("api.Signup"), inner) prints as:
//
//	api.Signup: user.Create: email: invalid: must not be empty
func (e *Error) Error() string { panic("TODO: implement Error.Error") }

func (e *Error) Unwrap() error { panic("TODO: implement Error.Unwrap") }

// Is lets errors.Is(err, E(KindNotFound)) work: an *Error target matches when
// every non-zero field of the target equals ours. (Compare Kind only when the
// target's Kind is not KindOther, Op only when non-empty, and so on.)
func (e *Error) Is(target error) bool { panic("TODO: implement Error.Is") }

// KindOf walks the error tree and returns the Kind of the outermost *Error that
// has one. It returns KindOther for nil, and KindInternal for a non-nil error
// with no Kind anywhere - an unclassified failure is an internal failure.
func KindOf(err error) Kind { panic("TODO: implement KindOf") }

// Ops returns the chain of operations, outermost first, skipping empty ones:
// exactly the breadcrumb you want in a log line.
func Ops(err error) []Op { panic("TODO: implement Ops") }

// Public renders a message that is safe to show to a user:
//
//	nil                                   -> ""
//	KindInternal / KindUnavailable / other unclassified -> "internal error"
//	otherwise -> the outermost *Error's Msg, prefixed with "<field>: " when the
//	             field is set; if Msg is empty, use the Kind's String()
func Public(err error) string { panic("TODO: implement Public") }

// HTTPStatus maps an error to a status code:
//
//	nil 200, invalid 400, permission 403, not found 404, conflict 409,
//	unavailable 503, everything else 500
func HTTPStatus(err error) int { panic("TODO: implement HTTPStatus") }
