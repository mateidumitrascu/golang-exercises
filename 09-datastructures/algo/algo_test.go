package algo

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func TestSearch(t *testing.T) {
	s := []int{1, 3, 5, 7, 9}
	for i, v := range s {
		if got := Search(s, v); got != i {
			t.Errorf("Search(%d) = %d, want %d", v, got, i)
		}
	}
	for _, v := range []int{0, 2, 10} {
		if got := Search(s, v); got != -1 {
			t.Errorf("Search(%d) = %d, want -1", v, got)
		}
	}
	if got := Search([]int{}, 1); got != -1 {
		t.Errorf("= %d", got)
	}
	if got := Search([]string{"a", "b", "c"}, "b"); got != 1 {
		t.Errorf("strings: = %d", got)
	}
	// Randomised cross-check against a linear scan.
	for range 200 {
		n := rand.Intn(50)
		s := make([]int, n)
		for i := range s {
			s[i] = rand.Intn(100)
		}
		sort.Ints(s)
		target := rand.Intn(100)
		got := Search(s, target)
		if got == -1 {
			for _, v := range s {
				if v == target {
					t.Fatalf("Search missed %d in %v", target, s)
				}
			}
		} else if s[got] != target {
			t.Fatalf("Search(%v, %d) = %d pointing at %d", s, target, got, s[got])
		}
	}
}

func TestBounds(t *testing.T) {
	s := []int{1, 2, 2, 2, 3, 5}
	tests := []struct {
		target       int
		lower, upper int
	}{
		{0, 0, 0},
		{1, 0, 1},
		{2, 1, 4},
		{3, 4, 5},
		{4, 5, 5},
		{5, 5, 6},
		{9, 6, 6},
	}
	for _, tt := range tests {
		if got := LowerBound(s, tt.target); got != tt.lower {
			t.Errorf("LowerBound(%d) = %d, want %d", tt.target, got, tt.lower)
		}
		if got := UpperBound(s, tt.target); got != tt.upper {
			t.Errorf("UpperBound(%d) = %d, want %d", tt.target, got, tt.upper)
		}
		if got := CountEqual(s, tt.target); got != tt.upper-tt.lower {
			t.Errorf("CountEqual(%d) = %d, want %d", tt.target, got, tt.upper-tt.lower)
		}
	}
}

func TestSearchRotated(t *testing.T) {
	base := []int{0, 1, 2, 4, 5, 6, 7}
	for r := range len(base) {
		s := append(append([]int{}, base[r:]...), base[:r]...)
		for i, v := range s {
			if got := SearchRotated(s, v); got != i {
				t.Errorf("SearchRotated(%v, %d) = %d, want %d", s, v, got, i)
			}
		}
		if got := SearchRotated(s, 99); got != -1 {
			t.Errorf("SearchRotated(%v, 99) = %d", s, got)
		}
	}
	if got := SearchRotated(nil, 1); got != -1 {
		t.Errorf("= %d", got)
	}
}

func TestSearchAnswer(t *testing.T) {
	// The smallest x with x*x >= 100.
	if got := SearchAnswer(0, 1000, func(x int) bool { return x*x >= 100 }); got != 10 {
		t.Errorf("= %d, want 10", got)
	}
	if got := SearchAnswer(0, 10, func(x int) bool { return false }); got != 11 {
		t.Errorf("never true = %d, want hi+1", got)
	}
	if got := SearchAnswer(5, 10, func(x int) bool { return true }); got != 5 {
		t.Errorf("always true = %d, want lo", got)
	}
	calls := 0
	SearchAnswer(0, 1000000, func(x int) bool { calls++; return x >= 12345 })
	if calls > 25 {
		t.Errorf("%d predicate calls over a million values; that is not binary search", calls)
	}
}

func TestTwoSumSorted(t *testing.T) {
	i, j := TwoSumSorted([]int{1, 2, 4, 7, 11}, 9)
	if i != 1 || j != 3 {
		t.Errorf("= %d, %d; want 1, 3", i, j)
	}
	if i, j := TwoSumSorted([]int{1, 2}, 100); i != -1 || j != -1 {
		t.Errorf("= %d, %d", i, j)
	}
	if i, j := TwoSumSorted([]int{-3, 0, 3}, 0); i != 0 || j != 2 {
		t.Errorf("= %d, %d", i, j)
	}
	if i, _ := TwoSumSorted(nil, 0); i != -1 {
		t.Error("empty input")
	}
}

func TestThreeSumZero(t *testing.T) {
	got := ThreeSumZero([]int{-1, 0, 1, 2, -1, -4})
	want := [][3]int{{-1, -1, 2}, {-1, 0, 1}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("= %v, want %v", got, want)
	}
	if got := ThreeSumZero([]int{0, 0, 0, 0}); !reflect.DeepEqual(got, [][3]int{{0, 0, 0}}) {
		t.Errorf("= %v, want one triple", got)
	}
	if got := ThreeSumZero([]int{1, 2, 3}); len(got) != 0 {
		t.Errorf("= %v", got)
	}
}

func TestMaxWater(t *testing.T) {
	if got := MaxWater([]int{1, 8, 6, 2, 5, 4, 8, 3, 7}); got != 49 {
		t.Errorf("= %d, want 49", got)
	}
	if got := MaxWater([]int{1, 1}); got != 1 {
		t.Errorf("= %d", got)
	}
	if got := MaxWater([]int{5}); got != 0 {
		t.Errorf("= %d", got)
	}
}

func TestIsPalindromeWords(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"A man, a plan, a canal: Panama", true},
		{"race a car", false},
		{"", true},
		{".,", true},
		{"Ana", true},
		{"éé", true},
		{"éa", false},
	}
	for _, tt := range tests {
		if got := IsPalindromeWords(tt.in); got != tt.want {
			t.Errorf("IsPalindromeWords(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestLongestUnique(t *testing.T) {
	tests := []struct {
		in   string
		want int
	}{
		{"abcabcbb", 3},
		{"bbbbb", 1},
		{"pwwkew", 3},
		{"", 0},
		{"日本語日本", 3},
		{"abcdef", 6},
	}
	for _, tt := range tests {
		if got := LongestUnique(tt.in); got != tt.want {
			t.Errorf("LongestUnique(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestMaxSubarraySum(t *testing.T) {
	got, ok := MaxSubarraySum([]int{2, 1, 5, 1, 3, 2}, 3)
	if !ok || got != 9 {
		t.Errorf("= %d, %v; want 9", got, ok)
	}
	if _, ok := MaxSubarraySum([]int{1, 2}, 5); ok {
		t.Error("k larger than the slice must report false")
	}
	if _, ok := MaxSubarraySum([]int{1, 2}, 0); ok {
		t.Error("k = 0 must report false")
	}
	if got, _ := MaxSubarraySum([]int{-1, -2, -3}, 2); got != -3 {
		t.Errorf("= %d, want -3", got)
	}
}

func TestMaxSubarray(t *testing.T) {
	sum, start, end := MaxSubarray([]int{-2, 1, -3, 4, -1, 2, 1, -5, 4})
	if sum != 6 || start != 3 || end != 7 {
		t.Errorf("= %d, [%d,%d); want 6, [3,7)", sum, start, end)
	}
	if sum, _, _ := MaxSubarray([]int{-5, -2, -8}); sum != -2 {
		t.Errorf("all negative = %d, want -2", sum)
	}
	if sum, s, e := MaxSubarray(nil); sum != 0 || s != 0 || e != 0 {
		t.Errorf("empty = %d, %d, %d", sum, s, e)
	}
	if sum, _, _ := MaxSubarray([]int{3}); sum != 3 {
		t.Errorf("= %d", sum)
	}
}

func TestMinWindow(t *testing.T) {
	tests := []struct {
		s, chars, want string
	}{
		{"ADOBECODEBANC", "ABC", "BANC"},
		{"a", "a", "a"},
		{"a", "aa", ""},
		{"", "a", ""},
		{"abc", "", ""},
		{"aabbcc", "abc", "abbc"},
	}
	for _, tt := range tests {
		if got := MinWindow(tt.s, tt.chars); got != tt.want {
			t.Errorf("MinWindow(%q, %q) = %q, want %q", tt.s, tt.chars, got, tt.want)
		}
	}
}

func TestLIS(t *testing.T) {
	tests := []struct {
		in   []int
		want int
	}{
		{[]int{10, 9, 2, 5, 3, 7, 101, 18}, 4},
		{[]int{0, 1, 0, 3, 2, 3}, 4},
		{[]int{7, 7, 7, 7}, 1},
		{nil, 0},
		{[]int{5}, 1},
	}
	for _, tt := range tests {
		if got := LIS(tt.in); got != tt.want {
			t.Errorf("LIS(%v) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestLISIsFast(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	// 300k elements: O(n^2) is 9e10 operations and will not finish.
	s := rand.Perm(300000)
	if got := LIS(s); got < 100 {
		t.Errorf("LIS = %d, suspiciously small", got)
	}
}

func TestCoinChange(t *testing.T) {
	tests := []struct {
		coins  []int
		amount int
		want   int
	}{
		{[]int{1, 5, 10, 25}, 30, 2},
		{[]int{2}, 3, -1},
		{[]int{1}, 0, 0},
		{[]int{}, 5, -1},
		{[]int{186, 419, 83, 408}, 6249, 20},
	}
	for _, tt := range tests {
		if got := CoinChange(tt.coins, tt.amount); got != tt.want {
			t.Errorf("CoinChange(%v, %d) = %d, want %d", tt.coins, tt.amount, got, tt.want)
		}
	}
}

func TestKnapsack(t *testing.T) {
	values := []int{60, 100, 120}
	weights := []int{10, 20, 30}
	if got := Knapsack(values, weights, 50); got != 220 {
		t.Errorf("= %d, want 220", got)
	}
	if got := Knapsack(values, weights, 0); got != 0 {
		t.Errorf("= %d", got)
	}
	if got := Knapsack(nil, nil, 10); got != 0 {
		t.Errorf("= %d", got)
	}
	// Each item once: with capacity 20 you cannot take the 10-weight item twice.
	if got := Knapsack([]int{60}, []int{10}, 20); got != 60 {
		t.Errorf("= %d, want 60 (0/1, not unbounded - iterate capacity downwards)", got)
	}
}

func TestWordBreak(t *testing.T) {
	tests := []struct {
		s    string
		dict []string
		want bool
	}{
		{"leetcode", []string{"leet", "code"}, true},
		{"applepenapple", []string{"apple", "pen"}, true},
		{"catsandog", []string{"cats", "dog", "sand", "and", "cat"}, false},
		{"", []string{"a"}, true},
		{"a", nil, false},
		{"aaaaaaaaaaaaaaaaaaaaaaaaaaaaab", []string{"a", "aa", "aaa"}, false},
	}
	for _, tt := range tests {
		if got := WordBreak(tt.s, tt.dict); got != tt.want {
			t.Errorf("WordBreak(%q, %v) = %v, want %v", tt.s, tt.dict, got, tt.want)
		}
	}
}

func TestGridPaths(t *testing.T) {
	tests := []struct {
		grid [][]int
		want int
	}{
		{[][]int{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}}, 6},
		{[][]int{{0, 0, 0}, {0, 1, 0}, {0, 0, 0}}, 2},
		{[][]int{{0}}, 1},
		{[][]int{{1}}, 0},
		{[][]int{{0, 1}, {1, 0}}, 0},
		{nil, 0},
	}
	for _, tt := range tests {
		if got := GridPaths(tt.grid); got != tt.want {
			t.Errorf("GridPaths(%v) = %d, want %d", tt.grid, got, tt.want)
		}
	}
}
