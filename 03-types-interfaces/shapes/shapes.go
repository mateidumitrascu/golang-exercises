// Package shapes is the classic interface exercise, with the parts that
// actually teach you something: method sets, sort.Interface, type switches,
// and returning a concrete type through an interface.
package shapes

import "sort"

// Shape is the core abstraction. Keep interfaces small - that is the Go way.
type Shape interface {
	Area() float64
	Perimeter() float64
}

// Scaler is a shape that can produce a resized copy of itself. It returns
// Shape, so each implementation decides its own concrete result type.
type Scaler interface {
	Shape
	Scaled(f float64) Shape
}

type Point struct{ X, Y float64 }

type Rect struct{ W, H float64 }

type Circle struct{ R float64 }

// Polygon is a closed polygon given by its vertices in order.
type Polygon struct{ Points []Point }

// All the methods below use VALUE receivers, so that Rect satisfies Shape and
// so does *Rect, and a []Shape can hold either. With a pointer receiver only
// *Rect would satisfy it. Make sure you can explain why before moving on.
//
// Scaled multiplies every linear dimension by f, so the area scales by f*f, and
// the concrete type is preserved (Rect.Scaled returns a Rect).
//
// Polygon: area is the shoelace formula,
//
//	|sum over i of (x_i*y_{i+1} - x_{i+1}*y_i)| / 2
//
// with the vertices wrapping around. Fewer than 3 points means zero area. The
// perimeter is the sum of all edges including the closing one (0 for fewer than
// 2 points). Polygon.Scaled must not modify the receiver's slice.

func (r Rect) Area() float64          { panic("TODO: Rect.Area") }
func (r Rect) Perimeter() float64     { panic("TODO: Rect.Perimeter") }
func (r Rect) Scaled(f float64) Shape { panic("TODO: Rect.Scaled") }

func (c Circle) Area() float64          { panic("TODO: Circle.Area") }
func (c Circle) Perimeter() float64     { panic("TODO: Circle.Perimeter") }
func (c Circle) Scaled(f float64) Shape { panic("TODO: Circle.Scaled") }

func (p Polygon) Area() float64          { panic("TODO: Polygon.Area") }
func (p Polygon) Perimeter() float64     { panic("TODO: Polygon.Perimeter") }
func (p Polygon) Scaled(f float64) Shape { panic("TODO: Polygon.Scaled") }

// TotalArea sums the areas. A nil slice totals 0.
func TotalArea(shapes []Shape) float64 { panic("TODO: implement TotalArea") }

// Largest returns the shape with the greatest area, and false if the slice is
// empty. Ties go to the earlier element.
func Largest(shapes []Shape) (Shape, bool) { panic("TODO: implement Largest") }

// ByArea orders shapes by area ascending, breaking ties by perimeter ascending.
type ByArea []Shape

func (a ByArea) Len() int           { panic("TODO: ByArea.Len") }
func (a ByArea) Less(i, j int) bool { panic("TODO: ByArea.Less") }
func (a ByArea) Swap(i, j int)      { panic("TODO: ByArea.Swap") }

// This line is a compile-time assertion that ByArea implements the interface.
// It is a very common idiom - use it in your own code.
var _ sort.Interface = ByArea(nil)

// SortByArea sorts in place using ByArea. It must be a STABLE sort.
func SortByArea(shapes []Shape) { panic("TODO: implement SortByArea") }

// Describe returns a human description using a type switch:
//
//	nil interface            "nothing"
//	Rect / *Rect             "rect 3x4"          (%g formatting)
//	Circle / *Circle         "circle r=2"
//	Polygon / *Polygon       "polygon with 5 sides"
//	any other Shape          "shape with area 12.5"
//	anything else            "not a shape"
//
// It takes `any`, not Shape: handling both the value and the pointer form of
// each type in one switch is the point.
func Describe(v any) string { panic("TODO: implement Describe") }

// ScaleAll returns a new slice in which every element that also implements
// Scaler is scaled by f. Elements that are not Scalers pass through unchanged.
func ScaleAll(shapes []Shape, f float64) []Shape { panic("TODO: implement ScaleAll") }
