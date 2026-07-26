package shapes

import (
	"math"
	"reflect"
	"sort"
	"testing"
)

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func isRect(s Shape) bool { _, ok := s.(Rect); return ok }

var (
	unitSquare = Polygon{Points: []Point{{0, 0}, {2, 0}, {2, 2}, {0, 2}}}
	tri        = Polygon{Points: []Point{{0, 0}, {3, 0}, {0, 4}}}
)

func TestAreaPerimeter(t *testing.T) {
	tests := []struct {
		name        string
		s           Shape
		area, perim float64
	}{
		{"rect", Rect{3, 4}, 12, 14},
		{"circle", Circle{2}, 4 * math.Pi, 4 * math.Pi},
		{"square polygon", unitSquare, 4, 8},
		{"3-4-5 triangle", tri, 6, 12},
		{"degenerate polygon", Polygon{Points: []Point{{0, 0}, {1, 1}}}, 0, 2 * math.Sqrt2},
		{"empty polygon", Polygon{}, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Area(); !near(got, tt.area) {
				t.Errorf("Area = %v, want %v", got, tt.area)
			}
			if got := tt.s.Perimeter(); !near(got, tt.perim) {
				t.Errorf("Perimeter = %v, want %v", got, tt.perim)
			}
		})
	}
}

// TestPointersAlsoSatisfyShape is the method-set check: with value receivers,
// both T and *T implement the interface.
func TestPointersAlsoSatisfyShape(t *testing.T) {
	var _ Shape = Rect{}
	var _ Shape = &Rect{}
	var _ Scaler = Circle{}
	var _ Scaler = &Circle{}
	list := []Shape{Rect{1, 1}, &Rect{2, 2}, Circle{1}, &Polygon{Points: tri.Points}}
	if got := len(list); got != 4 {
		t.Fatal(got)
	}
}

func TestScaled(t *testing.T) {
	if got := (Rect{3, 4}).Scaled(2); !near(got.Area(), 48) {
		t.Errorf("scaled rect area = %v, want 48", got.Area())
	}
	if got := (Circle{2}).Scaled(0.5); !near(got.Area(), math.Pi) {
		t.Errorf("scaled circle area = %v, want pi", got.Area())
	}
	// The concrete type must survive scaling.
	if scaled := (Rect{1, 1}).Scaled(2); !isRect(scaled) {
		t.Errorf("Rect.Scaled returned %T, want Rect", scaled)
	}
	p := Polygon{Points: []Point{{0, 0}, {1, 0}, {0, 1}}}
	scaled := p.Scaled(3)
	if !near(scaled.Area(), 4.5) {
		t.Errorf("scaled polygon area = %v, want 4.5", scaled.Area())
	}
	if p.Points[1].X != 1 {
		t.Error("Polygon.Scaled modified the receiver's points; copy the slice")
	}
}

func TestTotalAreaAndLargest(t *testing.T) {
	list := []Shape{Rect{1, 1}, Rect{2, 2}, Circle{1}}
	if got := TotalArea(list); !near(got, 1+4+math.Pi) {
		t.Errorf("TotalArea = %v", got)
	}
	if got := TotalArea(nil); got != 0 {
		t.Errorf("TotalArea(nil) = %v", got)
	}
	big, ok := Largest(list)
	if !ok || !near(big.Area(), 4) {
		t.Errorf("Largest = %v, %v", big, ok)
	}
	if _, ok := Largest(nil); ok {
		t.Error("Largest of an empty slice must report false")
	}
}

func TestSortByArea(t *testing.T) {
	list := []Shape{Rect{3, 3}, Circle{1}, Rect{1, 1}, Rect{2, 2}}
	SortByArea(list)
	areas := make([]float64, len(list))
	for i, s := range list {
		areas[i] = s.Area()
	}
	if !sort.Float64sAreSorted(areas) {
		t.Errorf("not sorted: %v", areas)
	}
	// Stability: two shapes with equal area and equal perimeter keep their order.
	a, b := Rect{2, 2}, Rect{2, 2}
	stable := []Shape{&a, &b}
	SortByArea(stable)
	if stable[0] != Shape(&a) {
		t.Error("SortByArea must be stable")
	}
}

func TestByAreaIsSortInterface(t *testing.T) {
	s := ByArea{Rect{2, 2}, Rect{1, 1}}
	if s.Len() != 2 {
		t.Errorf("Len = %d", s.Len())
	}
	if !s.Less(1, 0) || s.Less(0, 1) {
		t.Error("Less is wrong")
	}
	s.Swap(0, 1)
	if !near(s[0].Area(), 1) {
		t.Error("Swap is wrong")
	}
}

type weirdShape struct{}

func (weirdShape) Area() float64      { return 12.5 }
func (weirdShape) Perimeter() float64 { return 1 }

func TestDescribe(t *testing.T) {
	tests := []struct {
		in   any
		want string
	}{
		{nil, "nothing"},
		{Rect{3, 4}, "rect 3x4"},
		{&Rect{3, 4}, "rect 3x4"},
		{Circle{2}, "circle r=2"},
		{&Circle{2.5}, "circle r=2.5"},
		{Polygon{Points: make([]Point, 5)}, "polygon with 5 sides"},
		{weirdShape{}, "shape with area 12.5"},
		{"hello", "not a shape"},
		{42, "not a shape"},
	}
	for _, tt := range tests {
		if got := Describe(tt.in); got != tt.want {
			t.Errorf("Describe(%#v) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestScaleAll(t *testing.T) {
	in := []Shape{Rect{1, 1}, weirdShape{}}
	out := ScaleAll(in, 2)
	if !near(out[0].Area(), 4) {
		t.Errorf("scaled rect area = %v, want 4", out[0].Area())
	}
	if !near(out[1].Area(), 12.5) {
		t.Error("non-Scaler shapes must pass through unchanged")
	}
	if !reflect.DeepEqual(in[0], Shape(Rect{1, 1})) {
		t.Error("ScaleAll must not modify the input slice")
	}
}
