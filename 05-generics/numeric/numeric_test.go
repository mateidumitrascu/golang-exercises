package numeric

import (
	"math"
	"testing"
)

// Named types with the underlying types your constraints must accept.
type Celsius float64
type Count int
type Byte uint8
type Name string

func TestSum(t *testing.T) {
	if got := Sum([]int{1, 2, 3}); got != 6 {
		t.Errorf("Sum(ints) = %d, want 6", got)
	}
	if got := Sum([]float64{0.5, 0.25}); got != 0.75 {
		t.Errorf("Sum(floats) = %v", got)
	}
	if got := Sum([]int{}); got != 0 {
		t.Errorf("Sum(empty) = %d", got)
	}
	// Named types: this is what the tilde in the constraint buys you.
	if got := Sum([]Celsius{1.5, 2.5}); got != 4 {
		t.Errorf("Sum(Celsius) = %v, want 4", got)
	}
	if got := Sum([]Count{1, 2}); got != 3 {
		t.Errorf("Sum(Count) = %v", got)
	}
	// Arithmetic happens in T, so it wraps.
	if got := Sum([]uint8{200, 100}); got != 44 {
		t.Errorf("Sum(uint8) = %d, want 44 (wrapping is correct here)", got)
	}
}

func TestMean(t *testing.T) {
	got, ok := Mean([]int{1, 2, 3, 4})
	if !ok || got != 2.5 {
		t.Errorf("Mean = %v, %v; want 2.5, true", got, ok)
	}
	if _, ok := Mean([]float64{}); ok {
		t.Error("Mean of an empty slice must report false")
	}
}

func TestMinMax(t *testing.T) {
	if got, ok := MinOf(3, 1, 2); !ok || got != 1 {
		t.Errorf("MinOf = %v, %v", got, ok)
	}
	if got, ok := MaxOf("a", "c", "b"); !ok || got != "c" {
		t.Errorf("MaxOf(strings) = %v, %v", got, ok)
	}
	if _, ok := MinOf[int](); ok {
		t.Error("MinOf() must report false")
	}
	lo, hi, ok := MinMax([]Celsius{3, -1, 7})
	if !ok || lo != -1 || hi != 7 {
		t.Errorf("MinMax = %v, %v, %v", lo, hi, ok)
	}
	if _, _, ok := MinMax([]Name{}); ok {
		t.Error("MinMax of an empty slice must report false")
	}
	if got, _ := MaxOf(Name("a"), Name("b")); got != "b" {
		t.Errorf("MaxOf(Name) = %v", got)
	}
}

func TestClamp(t *testing.T) {
	tests := []struct{ v, lo, hi, want int }{
		{5, 1, 10, 5},
		{0, 1, 10, 1},
		{99, 1, 10, 10},
		{1, 1, 1, 1},
	}
	for _, tt := range tests {
		if got := Clamp(tt.v, tt.lo, tt.hi); got != tt.want {
			t.Errorf("Clamp(%d, %d, %d) = %d, want %d", tt.v, tt.lo, tt.hi, got, tt.want)
		}
	}
	if got := Clamp(2.5, 0.0, 1.0); got != 1.0 {
		t.Errorf("Clamp(float) = %v", got)
	}
	defer func() {
		if recover() == nil {
			t.Error("Clamp with lo > hi did not panic")
		}
	}()
	Clamp(1, 10, 0)
}

func TestAbs(t *testing.T) {
	if got := Abs(-3); got != 3 {
		t.Errorf("Abs(-3) = %d", got)
	}
	if got := Abs(-2.5); got != 2.5 {
		t.Errorf("Abs(-2.5) = %v", got)
	}
	if got := Abs(Celsius(-1)); got != 1 {
		t.Errorf("Abs(Celsius) = %v", got)
	}
	// The one case that cannot be right: the most negative integer has no
	// positive counterpart in its own type.
	if got := Abs(int8(math.MinInt8)); got != math.MinInt8 {
		t.Errorf("Abs(MinInt8) = %d, want %d unchanged (it has no positive counterpart)", got, math.MinInt8)
	}
}

func TestSumBy(t *testing.T) {
	type item struct {
		name  string
		price float64
		qty   int
	}
	items := []item{{"a", 1.5, 2}, {"b", 3, 1}}
	if got := SumBy(items, func(i item) float64 { return i.price * float64(i.qty) }); got != 6 {
		t.Errorf("SumBy = %v, want 6", got)
	}
	if got := SumBy(items, func(i item) int { return i.qty }); got != 3 {
		t.Errorf("SumBy(int) = %v, want 3", got)
	}
	if got := SumBy([]item{}, func(i item) int { return 1 }); got != 0 {
		t.Errorf("SumBy(empty) = %v", got)
	}
}

func TestCompare(t *testing.T) {
	if Compare(1, 2) != -1 || Compare(2, 1) != 1 || Compare(2, 2) != 0 {
		t.Error("Compare(ints) is wrong")
	}
	if Compare("a", "b") != -1 {
		t.Error("Compare(strings) is wrong")
	}
	if Compare(Byte(1), Byte(1)) != 0 {
		t.Error("Compare(named) is wrong")
	}
}

func TestInDelta(t *testing.T) {
	if !InDelta(0.1+0.2, 0.3, 1e-9) {
		t.Error("InDelta should smooth over float noise")
	}
	if InDelta(1.0, 2.0, 0.5) {
		t.Error("InDelta(1, 2, 0.5) should be false")
	}
	if !InDelta(float32(1), float32(1.0001), float32(0.001)) {
		t.Error("InDelta must work for float32 too")
	}
}
