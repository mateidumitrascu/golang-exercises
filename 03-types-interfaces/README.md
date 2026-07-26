# 03 — Types, methods, interfaces

Go has no inheritance, so the whole design language is: small interfaces,
concrete types with methods, and composition. This module covers the parts that
bite people who arrive from an OO language.

| # | Exercise | Difficulty | What it drills |
|---|----------|-----------|----------------|
| 1 | `shapes` | ★★☆ | Method sets, `sort.Interface`, type switches over value *and* pointer |
| 2 | `nilcheck` | ★★★ | Typed nil, interface identity, uncomparable dynamic types |
| 3 | `options` | ★★☆ | Functional options with validation and grouping |
| 4 | `expr` | ★★★ | Sealed interface as a sum type, tree walking, minimal-paren printing |
| 5 | `structtags` | ★★★ | Reflection: tags, kinds, recursion, unexported fields |

## Ideas worth internalising here

**Method sets.** With `func (t T) M()`, both `T` and `*T` have `M`. With
`func (t *T) M()`, only `*T` does — so `var s Shape = MyRect{}` fails to
compile if `Area` has a pointer receiver. Pick pointer receivers when the method
mutates or the struct is big, and then be consistent across the whole type.

**An interface value is two words: (type, value).** It is `nil` only if both are.
`var p *T = nil; var i any = p; i != nil` is **true**. This is the single most
common Go bug and exercise 2 makes you feel it before it costs you a production
incident.

**Accept interfaces, return structs.** `Scaled` returning `Shape` is the
exception, deliberately, so you can see what you lose: the caller has to type
assert to get back to the concrete type.

**A type switch on an unexported-method interface is exhaustive** — no other
package can add a case. That's how you get sum types in Go.

**Reflection is a last resort, but it is not magic**: `Kind` (the layout) is not
`Type` (the named type), `CanSet` is false for anything you reached by value or
that is unexported, and struct tags are just strings with a conventional format.

## Stretch goals

- Add `func (s ByArea) Sort()` and compare `sort.Sort` with `slices.SortFunc`
  in a benchmark. Why is the generic one usually faster?
- Give `expr` a `Derivative(e Expr, wrt string) Expr` and check it against a
  numeric difference quotient in a test.
- Make `structtags` support a `validate:"custom=positiveEven"` rule that looks
  up a registered `func(reflect.Value) error`.
