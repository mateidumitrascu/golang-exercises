package groupby

import (
	"reflect"
	"sort"
	"testing"
)

type person struct {
	Name string
	Age  int
	City string
}

var people = []person{
	{"ana", 31, "cluj"},
	{"bogdan", 24, "iasi"},
	{"cristi", 31, "cluj"},
	{"dana", 24, "bucuresti"},
}

func TestGroupBy(t *testing.T) {
	got := GroupBy(people, func(p person) int { return p.Age })
	if len(got) != 2 {
		t.Fatalf("got %d groups, want 2", len(got))
	}
	want31 := []person{people[0], people[2]}
	if !reflect.DeepEqual(got[31], want31) {
		t.Errorf("group 31 = %v, want %v (order must be preserved)", got[31], want31)
	}
	if len(got[24]) != 2 || got[24][0].Name != "bogdan" {
		t.Errorf("group 24 = %v", got[24])
	}
	if g := GroupBy(nil, func(p person) int { return p.Age }); len(g) != 0 {
		t.Errorf("GroupBy(nil) = %v, want empty", g)
	}
}

func TestIndex(t *testing.T) {
	got := Index(people, func(p person) string { return p.City })
	if got["cluj"].Name != "cristi" {
		t.Errorf(`Index[cluj] = %v, want cristi (last wins)`, got["cluj"].Name)
	}
	if len(got) != 3 {
		t.Errorf("len = %d, want 3", len(got))
	}
}

func TestCountBy(t *testing.T) {
	got := CountBy(people, func(p person) string { return p.City })
	want := map[string]int{"cluj": 2, "iasi": 1, "bucuresti": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CountBy = %v, want %v", got, want)
	}
}

func TestPartition(t *testing.T) {
	yes, no := Partition([]int{1, 2, 3, 4, 5}, func(n int) bool { return n%2 == 1 })
	if !reflect.DeepEqual(yes, []int{1, 3, 5}) || !reflect.DeepEqual(no, []int{2, 4}) {
		t.Errorf("Partition = %v, %v", yes, no)
	}
	y, n := Partition([]int{}, func(int) bool { return true })
	if y == nil || n == nil {
		t.Error("Partition must return non-nil slices even when a side is empty")
	}
}

func TestKeysValues(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3}
	ks := Keys(m)
	sort.Strings(ks)
	if !reflect.DeepEqual(ks, []string{"a", "b", "c"}) {
		t.Errorf("Keys = %v", ks)
	}
	vs := Values(m)
	sort.Ints(vs)
	if !reflect.DeepEqual(vs, []int{1, 2, 3}) {
		t.Errorf("Values = %v", vs)
	}
	if got := Keys(map[string]int{}); len(got) != 0 {
		t.Errorf("Keys of empty map = %v", got)
	}
}

func TestKeysPreallocates(t *testing.T) {
	m := make(map[int]int, 1000)
	for i := range 1000 {
		m[i] = i
	}
	if got := testing.AllocsPerRun(20, func() { _ = Keys(m) }); got > 1 {
		t.Errorf("Keys allocated %.0f times; make the slice with the right capacity up front", got)
	}
}

func TestInvert(t *testing.T) {
	got := Invert(map[string]int{"a": 1, "b": 2})
	want := map[int]string{1: "a", 2: "b"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Invert = %v, want %v", got, want)
	}
	dup := Invert(map[string]int{"a": 1, "b": 1})
	if len(dup) != 1 {
		t.Errorf("Invert with duplicate values = %v, want one entry", dup)
	}
}

func TestMergeWith(t *testing.T) {
	a := map[string]int{"x": 1, "y": 2}
	b := map[string]int{"y": 10, "z": 3}
	got := MergeWith(func(p, q int) int { return p + q }, a, b)
	want := map[string]int{"x": 1, "y": 12, "z": 3}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("MergeWith = %v, want %v", got, want)
	}
	if a["y"] != 2 {
		t.Error("MergeWith must not modify its inputs")
	}
	if got := MergeWith[string, int](func(p, q int) int { return q }); len(got) != 0 {
		t.Errorf("MergeWith with no maps = %v, want empty non-nil map", got)
	}
}

func TestSetDefault(t *testing.T) {
	m := map[string][]int{}
	v, existed := SetDefault(m, "k", []int{})
	if existed {
		t.Error("first call should report the key as absent")
	}
	m["k"] = append(v, 1)
	v2, existed2 := SetDefault(m, "k", []int{99})
	if !existed2 || !reflect.DeepEqual(v2, []int{1}) {
		t.Errorf("second call = %v, %v; want [1], true", v2, existed2)
	}
}
