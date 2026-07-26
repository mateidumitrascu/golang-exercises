// Package matrix: two-dimensional slices, index arithmetic, and the difference
// between a [][]int and a flat []int with a stride.
package matrix

// New allocates a rows x cols matrix backed by ONE flat []int, with each row a
// window into it. Allocating rows+1 times instead of 2 is the mistake this
// exercise exists to break you of.
//
// New panics if rows or cols is negative.
func New(rows, cols int) [][]int {
	panic("TODO: implement New")
}

// Transpose returns a new matrix where out[c][r] == m[r][c]. It works for
// non-square input. Transpose panics if m is ragged (rows of differing length).
func Transpose(m [][]int) [][]int {
	panic("TODO: implement Transpose")
}

// RotateCW rotates a SQUARE matrix 90 degrees clockwise, in place, without
// allocating. The classic approach is transpose-then-reverse-each-row; the
// alternative is a four-way element cycle over concentric rings.
//
// RotateCW panics if m is not square.
func RotateCW(m [][]int) {
	panic("TODO: implement RotateCW")
}

// Spiral returns the elements of m in clockwise spiral order, starting at
// m[0][0] and moving right.
//
//	1 2 3
//	4 5 6   ->  1 2 3 6 9 8 7 4 5
//	7 8 9
//
// Works for any rectangular shape, including a single row or column.
func Spiral(m [][]int) []int {
	panic("TODO: implement Spiral")
}

// Neighbours returns the values of the up-to-8 cells surrounding (r, c),
// in reading order (top-left, top, top-right, left, right, ...), skipping
// positions that fall outside the matrix.
func Neighbours(m [][]int, r, c int) []int {
	panic("TODO: implement Neighbours")
}

// Flat converts a matrix into a flat slice plus its stride, and Unflat converts
// back. Unflat panics if len(flat) is not a multiple of stride, or stride <= 0.
func Flat(m [][]int) (flat []int, stride int) {
	panic("TODO: implement Flat")
}

func Unflat(flat []int, stride int) [][]int {
	panic("TODO: implement Unflat")
}
