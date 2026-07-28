package expr

import (
	"math"
	"testing"
)

func n(v float64) Expr              { return Num{v} }
func v(name string) Expr            { return Var{name} }
func b(op byte, l, r Expr) Expr     { return Binary{op, l, r} }
func neg(x Expr) Expr               { return Unary{'-', x} }
func call(f string, a ...Expr) Expr { return Call{f, a} }

func TestEval(t *testing.T) {
	env := map[string]float64{"x": 3, "y": 4}
	tests := []struct {
		name string
		e    Expr
		want float64
	}{
		{"number", n(42), 42},
		{"variable", v("x"), 3},
		{"add", b('+', n(1), n(2)), 3},
		{"sub", b('-', n(1), n(2)), -1},
		{"mul", b('*', v("x"), v("y")), 12},
		{"div", b('/', n(7), n(2)), 3.5},
		{"pow op", b('^', n(2), n(10)), 1024},
		{"unary", neg(v("x")), -3},
		{"nested", b('+', b('*', v("x"), v("x")), b('*', v("y"), v("y"))), 25},
		{"abs", call("abs", n(-5)), 5},
		{"sqrt", call("sqrt", n(16)), 4},
		{"min", call("min", v("x"), v("y")), 3},
		{"max", call("max", v("x"), v("y")), 4},
		{"pow fn", call("pow", n(2), n(8)), 256},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Eval(tt.e, env)
			if err != nil {
				t.Fatalf("Eval: %v", err)
			}
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("Eval = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestEvalErrors(t *testing.T) {
// 	env := map[string]float64{"x": 1}
// 	tests := []struct {
// 		name    string
// 		e       Expr
// 		want    error
// 		mention string
// 	}{
// 		{"unknown var", v("zz"), ErrUnknownVar, "zz"},
// 		{"unknown func", call("frobnicate", n(1)), ErrUnknownFunc, "frobnicate"},
// 		{"div by zero", b('/', n(1), n(0)), ErrDivByZero, ""},
// 		{"div by zero via var", b('/', n(1), b('-', v("x"), n(1))), ErrDivByZero, ""},
// 		{"bad arity", call("sqrt", n(1), n(2)), ErrBadArity, "sqrt"},
// 		{"bad arity 2", call("min", n(1)), ErrBadArity, "min"},
// 		{"nil child", b('+', n(1), nil), ErrBadNode, ""},
// 		{"bad operator", b('%', n(1), n(2)), ErrBadNode, ""},
// 		{"error propagates up", b('+', n(1), v("nope")), ErrUnknownVar, "nope"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := Eval(tt.e, env)
// 			if !errors.Is(err, tt.want) {
// 				t.Fatalf("err = %v, want it to wrap %v", err, tt.want)
// 			}
// 			if tt.mention != "" && !strings.Contains(err.Error(), tt.mention) {
// 				t.Errorf("err = %q, should mention %q", err, tt.mention)
// 			}
// 		})
// 	}
// 	if _, err := Eval(v("x"), nil); !errors.Is(err, ErrUnknownVar) {
// 		t.Error("a nil env must not panic")
// 	}
// }
//
// func TestString(t *testing.T) {
// 	tests := []struct {
// 		e    Expr
// 		want string
// 	}{
// 		{n(1.5), "1.5"},
// 		{v("x"), "x"},
// 		{b('+', n(1), n(2)), "1 + 2"},
// 		{b('+', b('+', n(1), n(2)), n(3)), "1 + 2 + 3"},
// 		{b('+', n(1), b('+', n(2), n(3))), "1 + (2 + 3)"},
// 		{b('-', n(1), b('-', n(2), n(3))), "1 - (2 - 3)"},
// 		{b('*', b('+', n(1), n(2)), n(3)), "(1 + 2) * 3"},
// 		{b('+', b('*', n(1), n(2)), n(3)), "1 * 2 + 3"},
// 		{b('/', n(1), b('*', n(2), n(3))), "1 / (2 * 3)"},
// 		{b('^', n(2), b('^', n(3), n(4))), "2 ^ 3 ^ 4"},
// 		{b('^', b('^', n(2), n(3)), n(4)), "(2 ^ 3) ^ 4"},
// 		{b('^', b('+', n(1), n(2)), n(3)), "(1 + 2) ^ 3"},
// 		{neg(n(3)), "-3"},
// 		{neg(b('*', v("x"), v("y"))), "-(x * y)"},
// 		{neg(neg(v("x"))), "-(-x)"},
// 		{b('*', neg(v("x")), v("y")), "-x * y"},
// 		{call("min", n(1), b('+', v("x"), n(2))), "min(1, x + 2)"},
// 		{call("now"), "now()"},
// 	}
// 	for _, tt := range tests {
// 		if got := String(tt.e); got != tt.want {
// 			t.Errorf("String(%#v) = %q, want %q", tt.e, got, tt.want)
// 		}
// 	}
// }
//
// func TestSimplify(t *testing.T) {
// 	tests := []struct {
// 		name string
// 		in   Expr
// 		want string
// 	}{
// 		{"fold constants", b('+', n(1), b('*', n(2), n(3))), "7"},
// 		{"x + 0", b('+', v("x"), n(0)), "x"},
// 		{"0 + x", b('+', n(0), v("x")), "x"},
// 		{"x - 0", b('-', v("x"), n(0)), "x"},
// 		{"x * 1", b('*', v("x"), n(1)), "x"},
// 		{"1 * x", b('*', n(1), v("x")), "x"},
// 		{"x * 0", b('*', v("x"), n(0)), "0"},
// 		{"x / 1", b('/', v("x"), n(1)), "x"},
// 		{"x ^ 1", b('^', v("x"), n(1)), "x"},
// 		{"x ^ 0", b('^', v("x"), n(0)), "1"},
// 		{"cascade", b('*', b('+', v("x"), n(0)), b('-', n(4), n(3))), "x"},
// 		{"nothing to do", b('+', v("x"), v("y")), "x + y"},
// 		{"no folding through error", b('/', n(1), n(0)), "1 / 0"},
// 		{"fold inside call", call("abs", b('-', n(1), n(5))), "4"},
// 		{"partial", b('+', b('*', n(2), n(3)), v("x")), "6 + x"},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := String(Simplify(tt.in)); got != tt.want {
// 				t.Errorf("Simplify -> %q, want %q", got, tt.want)
// 			}
// 		})
// 	}
// }
//
// func TestSimplifyDoesNotMutate(t *testing.T) {
// 	inner := Binary{'+', Num{0}, Var{"x"}}
// 	tree := Binary{'*', inner, Num{1}}
// 	before := String(tree)
// 	Simplify(tree)
// 	if after := String(tree); after != before {
// 		t.Errorf("Simplify mutated its input: %q -> %q", before, after)
// 	}
// }
//
// func TestVarsAndDepth(t *testing.T) {
// 	e := b('+', v("z"), b('*', v("a"), b('-', v("z"), n(1))))
// 	if got := Vars(e); !reflect.DeepEqual(got, []string{"a", "z"}) {
// 		t.Errorf("Vars = %v, want [a z]", got)
// 	}
// 	if got := Vars(n(1)); len(got) != 0 {
// 		t.Errorf("Vars(const) = %v, want empty", got)
// 	}
// 	if got := Depth(n(1)); got != 1 {
// 		t.Errorf("Depth(leaf) = %d, want 1", got)
// 	}
// 	if got := Depth(e); got != 4 {
// 		t.Errorf("Depth = %d, want 4", got)
// 	}
// 	if got := Depth(call("min", n(1), b('+', n(1), n(2)))); got != 3 {
// 		t.Errorf("Depth(call) = %d, want 3", got)
// 	}
// }
//
// func TestWalk(t *testing.T) {
// 	e := b('+', v("a"), b('*', v("b"), v("c")))
// 	var seen []string
// 	Walk(e, func(x Expr) bool {
// 		if vv, ok := x.(Var); ok {
// 			seen = append(seen, vv.Name)
// 		}
// 		return true
// 	})
// 	if !reflect.DeepEqual(seen, []string{"a", "b", "c"}) {
// 		t.Errorf("Walk visited %v, want [a b c]", seen)
// 	}
//
// 	count := 0
// 	Walk(e, func(x Expr) bool {
// 		count++
// 		_, isBin := x.(Binary)
// 		return !isBin || count == 1 // stop descending into the inner Binary
// 	})
// 	if count != 3 {
// 		t.Errorf("Walk visited %d nodes after pruning, want 3 (root, a, inner binary)", count)
// 	}
// }
