package graph

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"
)

func directed(edges ...[2]string) *Graph[string] {
	g := NewDirected[string]()
	for _, e := range edges {
		g.AddEdge(e[0], e[1], 1)
	}
	return g
}

func TestAddAndInspect(t *testing.T) {
	g := NewDirected[string]()
	g.AddNode("lonely")
	g.AddEdge("a", "b", 2)
	g.AddEdge("a", "c", 1)
	g.AddEdge("a", "b", 5) // replaces the weight

	if got := g.Nodes(); !reflect.DeepEqual(got, []string{"a", "b", "c", "lonely"}) {
		t.Errorf("Nodes = %v", got)
	}
	if g.Len() != 4 {
		t.Errorf("Len = %d", g.Len())
	}
	if g.Edges() != 2 {
		t.Errorf("Edges = %d, want 2 (the duplicate replaces)", g.Edges())
	}
	want := []Edge[string]{{"b", 5}, {"c", 1}}
	if got := g.Neighbours("a"); !reflect.DeepEqual(got, want) {
		t.Errorf("Neighbours = %v, want %v", got, want)
	}
	if got := g.Neighbours("b"); len(got) != 0 {
		t.Errorf("Neighbours(b) = %v, want none in a directed graph", got)
	}
	if got := g.Neighbours("unknown"); len(got) != 0 {
		t.Errorf("= %v", got)
	}
}

func TestUndirected(t *testing.T) {
	g := NewUndirected[string]()
	g.AddEdge("a", "b", 3)
	if got := g.Neighbours("b"); len(got) != 1 || got[0].To != "a" || got[0].Weight != 3 {
		t.Errorf("Neighbours(b) = %v, want the reverse edge", got)
	}
	if g.Edges() != 1 {
		t.Errorf("Edges = %d, want 1", g.Edges())
	}
}

func TestNegativeWeightPanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a negative weight did not panic")
		}
	}()
	NewDirected[string]().AddEdge("a", "b", -1)
}

func TestBFSDFS(t *testing.T) {
	//   a -> b -> d
	//   a -> c -> d -> e
	g := directed([2]string{"a", "b"}, [2]string{"a", "c"}, [2]string{"b", "d"},
		[2]string{"c", "d"}, [2]string{"d", "e"})

	if got := g.BFS("a"); !reflect.DeepEqual(got, []string{"a", "b", "c", "d", "e"}) {
		t.Errorf("BFS = %v", got)
	}
	if got := g.DFS("a"); !reflect.DeepEqual(got, []string{"a", "b", "d", "e", "c"}) {
		t.Errorf("DFS = %v, want [a b d e c] (neighbours in sorted order)", got)
	}
	if got := g.BFS("e"); !reflect.DeepEqual(got, []string{"e"}) {
		t.Errorf("BFS(e) = %v", got)
	}
	if got := g.BFS("nope"); got != nil {
		t.Errorf("BFS of an unknown node = %v, want nil", got)
	}
}

func TestTraversalHandlesCycles(t *testing.T) {
	g := directed([2]string{"a", "b"}, [2]string{"b", "c"}, [2]string{"c", "a"})
	if got := g.BFS("a"); len(got) != 3 {
		t.Errorf("BFS = %v, want 3 nodes and no infinite loop", got)
	}
	if got := g.DFS("a"); len(got) != 3 {
		t.Errorf("DFS = %v", got)
	}
	self := directed([2]string{"a", "a"})
	if got := self.DFS("a"); len(got) != 1 {
		t.Errorf("self-loop DFS = %v", got)
	}
}

func TestShortestPath(t *testing.T) {
	g := directed([2]string{"a", "b"}, [2]string{"b", "c"}, [2]string{"a", "z"},
		[2]string{"z", "c"}, [2]string{"c", "d"})
	if got := g.ShortestPath("a", "d"); !reflect.DeepEqual(got, []string{"a", "b", "c", "d"}) {
		t.Errorf("= %v, want the lexically smallest of the equal-length paths", got)
	}
	if got := g.ShortestPath("a", "a"); !reflect.DeepEqual(got, []string{"a"}) {
		t.Errorf("= %v", got)
	}
	if got := g.ShortestPath("d", "a"); got != nil {
		t.Errorf("= %v, want nil", got)
	}
	if got := g.ShortestPath("a", "nope"); got != nil {
		t.Errorf("= %v, want nil", got)
	}
}

func TestTopoSort(t *testing.T) {
	// A classic build-dependency graph.
	g := directed(
		[2]string{"config", "server"},
		[2]string{"logger", "server"},
		[2]string{"logger", "db"},
		[2]string{"db", "server"},
		[2]string{"server", "main"},
	)
	got, err := g.TopoSort()
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"config", "logger", "db", "server", "main"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("TopoSort = %v, want %v (smallest-first tie-break)", got, want)
	}

	pos := map[string]int{}
	for i, n := range got {
		pos[n] = i
	}
	for _, n := range g.Nodes() {
		for _, e := range g.Neighbours(n) {
			if pos[n] > pos[e.To] {
				t.Errorf("edge %s->%s points backwards", n, e.To)
			}
		}
	}

	cyclic := directed([2]string{"a", "b"}, [2]string{"b", "c"}, [2]string{"c", "a"})
	if _, err := cyclic.TopoSort(); !errors.Is(err, ErrCycle) {
		t.Errorf("err = %v, want ErrCycle", err)
	}
	empty, err := NewDirected[string]().TopoSort()
	if err != nil || len(empty) != 0 {
		t.Errorf("empty graph = %v, %v", empty, err)
	}
}

func TestHasCycle(t *testing.T) {
	if directed([2]string{"a", "b"}, [2]string{"b", "c"}).HasCycle() {
		t.Error("a DAG has no cycle")
	}
	if !directed([2]string{"a", "b"}, [2]string{"b", "a"}).HasCycle() {
		t.Error("missed a two-node cycle")
	}
	if !directed([2]string{"a", "a"}).HasCycle() {
		t.Error("missed a self-loop")
	}
	// A diamond is not a cycle, even though a node is reachable twice.
	if directed([2]string{"a", "b"}, [2]string{"a", "c"}, [2]string{"b", "d"}, [2]string{"c", "d"}).HasCycle() {
		t.Error("a diamond DAG was reported as cyclic - are you marking nodes as " +
			"'in progress' and 'done' separately?")
	}
	u := NewUndirected[string]()
	u.AddEdge("a", "b", 1)
	if u.HasCycle() {
		t.Error("a single undirected edge is not a cycle (do not count the way back)")
	}
	u.AddEdge("b", "c", 1)
	u.AddEdge("c", "a", 1)
	if !u.HasCycle() {
		t.Error("missed an undirected triangle")
	}
}

func TestDijkstra(t *testing.T) {
	g := NewDirected[string]()
	g.AddEdge("a", "b", 1)
	g.AddEdge("b", "c", 2)
	g.AddEdge("a", "c", 10)
	g.AddEdge("c", "d", 1)
	g.AddNode("island")

	dist, _ := g.Dijkstra("a")
	want := map[string]float64{"a": 0, "b": 1, "c": 3, "d": 4}
	if !reflect.DeepEqual(dist, want) {
		t.Errorf("dist = %v, want %v", dist, want)
	}
	if _, ok := dist["island"]; ok {
		t.Error("unreachable nodes must be absent, not infinite")
	}

	path, cost := g.Path("a", "d")
	if !reflect.DeepEqual(path, []string{"a", "b", "c", "d"}) || cost != 4 {
		t.Errorf("Path = %v, %v", path, cost)
	}
	if p, _ := g.Path("a", "island"); p != nil {
		t.Errorf("Path to an unreachable node = %v", p)
	}
	if p, c := g.Path("a", "a"); !reflect.DeepEqual(p, []string{"a"}) || c != 0 {
		t.Errorf("Path to self = %v, %v", p, c)
	}
}

func TestDijkstraFloatWeights(t *testing.T) {
	g := NewUndirected[string]()
	g.AddEdge("a", "b", 0.5)
	g.AddEdge("b", "c", 0.25)
	dist, _ := g.Dijkstra("a")
	if math.Abs(dist["c"]-0.75) > 1e-9 {
		t.Errorf("dist[c] = %v, want 0.75", dist["c"])
	}
}

func TestComponents(t *testing.T) {
	g := NewUndirected[string]()
	g.AddEdge("a", "b", 1)
	g.AddEdge("b", "c", 1)
	g.AddEdge("x", "y", 1)
	g.AddNode("solo")

	got := g.Components()
	want := [][]string{{"a", "b", "c"}, {"solo"}, {"x", "y"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Components = %v, want %v", got, want)
	}
	if got := NewUndirected[string]().Components(); len(got) != 0 {
		t.Errorf("empty graph = %v", got)
	}
}

// TestDijkstraScales builds a 50k-node grid-ish graph. The O(V^2) version does
// 2.5 billion operations here; the priority-queue version does a few million.
func TestDijkstraScales(t *testing.T) {
	if testing.Short() {
		t.Skip("short mode")
	}
	const n = 50000
	g := NewDirected[int]()
	for i := range n {
		g.AddEdge(i, (i+1)%n, 1)
		g.AddEdge(i, (i+7)%n, 3)
	}
	dist, _ := g.Dijkstra(0)
	if len(dist) != n {
		t.Fatalf("reached %d nodes, want %d", len(dist), n)
	}
	if dist[1] != 1 || dist[7] != 3 {
		t.Errorf("dist[1] = %v, dist[7] = %v", dist[1], dist[7])
	}
}

func ExampleGraph_TopoSort() {
	g := NewDirected[string]()
	g.AddEdge("eggs", "cake", 1)
	g.AddEdge("flour", "cake", 1)
	g.AddEdge("cake", "party", 1)
	order, _ := g.TopoSort()
	fmt.Println(order)
	// Output: [eggs flour cake party]
}
