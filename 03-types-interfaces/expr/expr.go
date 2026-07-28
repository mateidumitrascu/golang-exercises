// Package expr builds an expression tree and walks it. This is the shape of
// every interpreter, query planner and rules engine you will ever write, and it
// is the best possible drill for type switches and sealed interfaces.
//
// Note the unexported exprNode method: it means no type outside this package
// can implement Expr, so a type switch over the known node types is exhaustive.
// That is Go's version of a sum type.
package expr

import (
	"errors"
	"fmt"
	"log"
	"math"
)

type Expr interface{ exprNode() }

type (
	Num   struct{ Value float64 }
	Var   struct{ Name string }
	Unary struct { // Op is always '-'
		Op byte
		X  Expr
	}
)

type Binary struct { // Op is one of + - * / ^
	Op   byte
	L, R Expr
}
type Call struct {
	Fn   string
	Args []Expr
}

// These are given: the sealing methods.
func (Num) exprNode()    {}
func (Var) exprNode()    {}
func (Unary) exprNode()  {}
func (Binary) exprNode() {}
func (Call) exprNode()   {}

var (
	ErrUnknownVar  = errors.New("unknown variable")
	ErrUnknownFunc = errors.New("unknown function")
	ErrDivByZero   = errors.New("division by zero")
	ErrBadArity    = errors.New("wrong number of arguments")
	ErrBadNode     = errors.New("malformed expression")
)

// Eval computes the value of e. Supported functions, with their arities:
//
//	abs(x) sqrt(x)      1 argument
//	min(a,b) max(a,b) pow(a,b)   2 arguments
//
// Errors wrap the sentinels above and must include the offending name:
// an unknown variable "z" produces an error whose message contains "z" and for
// which errors.Is(err, ErrUnknownVar) is true. A nil child node, or an unknown
// operator, is ErrBadNode.
func Eval(e Expr, env map[string]float64) (float64, error) {
	switch v := e.(type) {
	case nil:
		return 0, fmt.Errorf("%w: null node", ErrBadNode)
	case Num:
		return v.Value, nil

	case Var:
		return env[v.Name], nil

	case Unary:
		f, err := Eval(v.X, env)
		return -1 * f, err

	case Binary:
		l, err := Eval(v.L, env)
		if err != nil {
			return 0, err
		}
		r, err := Eval(v.R, env)
		if err != nil {
			return 0, err
		}
		log.Printf("PERFORMING OPERATION %c WITH %f AND %f \n\n", v.Op, l, r)
		return ops[v.Op](l, r)

	case Call:
		results := []float64{}
		for _, exp := range v.Args {
			r, err := Eval(exp, env)
			if err != nil {
				return 0, err
			}
			results = append(results, r)
		}
		return calls[v.Fn](results...)
	}
	return 0, nil
}

// String renders e with the FEWEST parentheses that still parse back to the
// same tree. Precedence: ^ (3) binds tighter than * / (2), which bind tighter
// than + - (1). Unary minus binds tighter than any binary operator. + - * /
// are left-associative; ^ is right-associative.
//
//	Binary{'+', 1, Binary{'+', 2, 3}}  ->  "1 + (2 + 3)"
//	Binary{'+', Binary{'+', 1, 2}, 3}  ->  "1 + 2 + 3"
//	Binary{'*', Binary{'+', 1, 2}, 3}  ->  "(1 + 2) * 3"
//	Binary{'^', 2, Binary{'^', 3, 4}}  ->  "2 ^ 3 ^ 4"
//	Binary{'^', Binary{'^', 2, 3}, 4}  ->  "(2 ^ 3) ^ 4"
//	Unary{'-', Binary{'*', x, y}}      ->  "-(x * y)"
//
// Binary operators are surrounded by single spaces; call arguments are joined
// with ", "; numbers use %g.
func String(e Expr) string {
	panic("TODO: implement String")
}

// Simplify returns an equivalent, smaller tree. It must apply, repeatedly,
// until nothing changes:
//
//	constant folding, when a subtree has no variables and no error
//	x + 0, 0 + x, x - 0, x * 1, 1 * x, x / 1, x ^ 1   ->  x
//	x * 0, 0 * x                                     ->  0
//	x ^ 0                                            ->  1
//
// Do not fold anything that would error (division by zero, unknown function):
// leave those subtrees alone. Simplify must not modify the input tree.
func Simplify(e Expr) Expr {
	panic("TODO: implement Simplify")
}

// Vars returns the distinct variable names in e, sorted.
func Vars(e Expr) []string {
	panic("TODO: implement Vars")
}

// Depth is the height of the tree: a leaf is 1.
func Depth(e Expr) int {
	panic("TODO: implement Depth")
}

// Walk calls fn for every node, parents before children, and stops early if fn
// returns false for a node (that node's children are skipped, siblings are not).
func Walk(e Expr, fn func(Expr) bool) {
	panic("TODO: implement Walk")
}

var ops = map[byte]func(a, b float64) (float64, error){
	'+': func(a, b float64) (float64, error) {
		return a + b, nil
	},
	'-': func(a, b float64) (float64, error) {
		return a - b, nil
	},
	'*': func(a, b float64) (float64, error) {
		return a * b, nil
	},
	'/': func(a, b float64) (float64, error) {
		if b == 0 {
			return 0, ErrDivByZero
		}
		return a / b, nil
	},
	'^': func(a, b float64) (float64, error) {
		if a == 0 && b == 0 {
			return 0, ErrDivByZero
		}
		return math.Pow(a, b), nil
	},
}

var calls = map[string]func(args ...float64) (float64, error){
	"abs": func(args ...float64) (float64, error) {
		if len(args) != 1 {
			return 0, fmt.Errorf("%w: provided %d arguments instead of 1", ErrBadArity, len(args))
		}
		return math.Abs(args[0]), nil
	},

	"sqrt": func(args ...float64) (float64, error) {
		if len(args) != 1 {
			return 0, fmt.Errorf("%w: provided %d arguments instead of 1", ErrBadArity, len(args))
		}
		return math.Sqrt(args[0]), nil
	},
	"min": func(args ...float64) (float64, error) {
		if len(args) != 2 {
			return 0, fmt.Errorf("%w: provided %d arguments instead of 2", ErrBadArity, len(args))
		}
		return math.Min(args[0], args[1]), nil
	},
	"max": func(args ...float64) (float64, error) {
		if len(args) != 2 {
			return 0, fmt.Errorf("%w: provided %d arguments instead of 2", ErrBadArity, len(args))
		}
		return math.Max(args[0], args[1]), nil
	},
	"pow": func(args ...float64) (float64, error) {
		if len(args) != 2 {
			return 0, fmt.Errorf("%w: provided %d arguments instead of 2", ErrBadArity, len(args))
		}
		return math.Pow(args[0], args[1]), nil
	},
}
