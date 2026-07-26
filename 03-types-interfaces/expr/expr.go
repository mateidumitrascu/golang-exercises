// Package expr builds an expression tree and walks it. This is the shape of
// every interpreter, query planner and rules engine you will ever write, and it
// is the best possible drill for type switches and sealed interfaces.
//
// Note the unexported exprNode method: it means no type outside this package
// can implement Expr, so a type switch over the known node types is exhaustive.
// That is Go's version of a sum type.
package expr

import "errors"

type Expr interface{ exprNode() }

type Num struct{ Value float64 }
type Var struct{ Name string }
type Unary struct { // Op is always '-'
	Op byte
	X  Expr
}
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
	panic("TODO: implement Eval")
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
