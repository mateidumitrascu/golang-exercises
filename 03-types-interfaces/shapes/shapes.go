// Package shapes is the classic interface exercise, with the parts that
// actually teach you something: method sets, sort.Interface, type switches,
// and returning a concrete type through an interface.
package shapes

import (
	"fmt"
	"math"
	"sort"
)

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

func (r Rect) Area() float64 {
	return r.H * r.W
}

func (r Rect) Perimeter() float64 {
	return 2*r.H + 2*r.W
}

func (r Rect) Scaled(f float64) Shape {
	return Rect{f * r.W, f * r.H}
}

func (c Circle) Area() float64 {
	return math.Pi * c.R * c.R
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.R
}

func (c Circle) Scaled(f float64) Shape {
	return Circle{f * c.R}
}

func (p Polygon) Area() float64 {
	if len(p.Points) < 3 {
		return 0
	}
	var sum float64 = 0
	for i := 0; i < len(p.Points)-1; i++ {
		sum += p.Points[i].X*p.Points[i+1].Y - p.Points[i+1].X*p.Points[i].Y
	}
	return sum / 2
}

func (p Polygon) Perimeter() float64 {
	if len(p.Points) < 2 {
		return 0
	}
	var sum float64 = 0
	for i := 0; i < len(p.Points)-1; i++ {
		sum += math.Sqrt(square(p.Points[i].X-p.Points[i+1].X) + square(p.Points[i].Y-p.Points[i+1].Y))
	}
	c := len(p.Points)
	sum += math.Sqrt(square(p.Points[0].X-p.Points[c-1].X) + square(p.Points[0].Y-p.Points[c-1].Y))
	return sum
}

func (p Polygon) Scaled(f float64) Shape {
	points := make([]Point, len(p.Points))
	for i, pt := range p.Points {
		points[i].X = pt.X * f
		points[i].Y = pt.Y * f
	}
	return Polygon{points}
}

// TotalArea sums the areas. A nil slice totals 0.
func TotalArea(shapes []Shape) float64 {
	if shapes == nil {
		return 0
	}
	var sum float64 = 0

	for _, shape := range shapes {
		sum += shape.Area()
	}
	return sum
}

// Largest returns the shape with the greatest area, and false if the slice is
// empty. Ties go to the earlier element.
func Largest(shapes []Shape) (Shape, bool) {
	if len(shapes) == 0 {
		return nil, false
	}
	var largest Shape
	var largestArea float64
	for _, shape := range shapes {
		if shape.Area() > largestArea {
			largestArea = shape.Area()
			largest = shape
		}
	}

	return largest, true
}

// ByArea orders shapes by area ascending, breaking ties by perimeter ascending.
type ByArea []Shape

func (a ByArea) Len() int {
	return len(a)
}

func (a ByArea) Less(i, j int) bool {
	return a[i].Area() < a[j].Area()
}

func (a ByArea) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// This line is a compile-time assertion that ByArea implements the interface.
// It is a very common idiom - use it in your own code.
var _ sort.Interface = ByArea(nil)

// SortByArea sorts in place using ByArea. It must be a STABLE sort.
func SortByArea(shapes []Shape) {
	sort.Stable(ByArea(shapes))
}

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
func Describe(v any) string {
	switch v.(type) {
	case nil:
		return "nothing"

	case Rect, *Rect:
		if r, ok := v.(Rect); ok {
			return fmt.Sprintf("rect %gx%g", r.W, r.H)
		}
		r := *v.(*Rect)
		return fmt.Sprintf("rect %gx%g", r.W, r.H)

	case Circle, *Circle:
		if c, ok := v.(Circle); ok {
			return fmt.Sprintf("circle r=%g", c.R)
		}
		c := *v.(*Circle)
		return fmt.Sprintf("circle r=%g", c.R)

	case Polygon, *Polygon:
		if p, ok := v.(Polygon); ok {
			return fmt.Sprintf("polygon with %d sides", len(p.Points))
		}
		p := *v.(*Polygon)
		return fmt.Sprintf("polygon with %d sides", len(p.Points))

	case Shape:
		return fmt.Sprintf("shape with area %g", v.(Shape).Area())

	default:
		return "not a shape"
	}
}

// ScaleAll returns a new slice in which every element that also implements
// Scaler is scaled by f. Elements that are not Scalers pass through unchanged.
func ScaleAll(shapes []Shape, f float64) []Shape {
	cp := make([]Shape, len(shapes))
	copy(cp, shapes)
	for i, s := range cp {
		if sc, ok := s.(Scaler); ok {
			cp[i] = sc.Scaled(f)
		}
	}
	return cp
}

func square(x float64) float64 {
	return x * x
}
