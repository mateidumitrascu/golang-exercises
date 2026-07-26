// Package graph implements the traversals and shortest-path algorithms every
// backend engineer eventually needs: dependency resolution, deadlock detection,
// routing, scheduling.
//
// One design decision worth noticing: node IDs are constrained to cmp.Ordered
// so that every result can be returned in a deterministic order. Graph
// algorithms over a Go map are otherwise randomly ordered, and randomly ordered
// output is untestable.
package graph

import (
	"cmp"
	"errors"
)

// Edge is a weighted connection.
type Edge[T cmp.Ordered] struct {
	To     T
	Weight float64
}

// Graph is an adjacency-list graph, directed or not.
type Graph[T cmp.Ordered] struct {
	// TODO: your fields.
	// In an undirected graph, AddEdge must record the edge in both directions.
}

func NewDirected[T cmp.Ordered]() *Graph[T]   { panic("TODO: implement NewDirected") }
func NewUndirected[T cmp.Ordered]() *Graph[T] { panic("TODO: implement NewUndirected") }

// AddNode adds an isolated node. Adding one twice is a no-op.
func (g *Graph[T]) AddNode(n T) { panic("TODO: implement AddNode") }

// AddEdge adds an edge, creating either endpoint if it is unknown. Adding the
// same edge twice replaces the weight rather than duplicating it. A self-loop
// is allowed. It panics on a negative weight - Dijkstra silently gives wrong
// answers for those, so refuse them at the door.
func (g *Graph[T]) AddEdge(from, to T, weight float64) { panic("TODO: implement AddEdge") }

// Nodes returns every node, sorted.
func (g *Graph[T]) Nodes() []T { panic("TODO: implement Nodes") }

// Neighbours returns the outgoing edges of n, sorted by destination.
func (g *Graph[T]) Neighbours(n T) []Edge[T] { panic("TODO: implement Neighbours") }

// Len is the number of nodes; Edges is the number of edges (an undirected edge
// counts once).
func (g *Graph[T]) Len() int   { panic("TODO: implement Len") }
func (g *Graph[T]) Edges() int { panic("TODO: implement Edges") }

// BFS returns the nodes reachable from start in breadth-first order, start
// included. Neighbours are visited in sorted order so the result is
// deterministic. An unknown start returns nil.
func (g *Graph[T]) BFS(start T) []T { panic("TODO: implement BFS") }

// DFS returns the nodes in depth-first pre-order, again visiting neighbours in
// sorted order. Write it ITERATIVELY with an explicit stack: recursion blows up
// on a deep graph, and the iterative version is the one that teaches you what
// the recursion was doing.
func (g *Graph[T]) DFS(start T) []T { panic("TODO: implement DFS") }

// ShortestPath returns the fewest-hops path from a to b (ignoring weights),
// inclusive of both ends, or nil if there is none. Among equal-length paths,
// return the lexically smallest.
func (g *Graph[T]) ShortestPath(a, b T) []T { panic("TODO: implement ShortestPath") }

var ErrCycle = errors.New("graph: cycle detected")

// TopoSort returns a topological ordering of a directed graph: every edge
// points forwards in the result. Ties are broken by choosing the smallest
// available node, which makes the answer unique (that is Kahn's algorithm with
// a priority queue). It returns ErrCycle if the graph has one.
func (g *Graph[T]) TopoSort() ([]T, error) { panic("TODO: implement TopoSort") }

// HasCycle reports whether the graph contains a cycle. For a directed graph
// that is the classic three-colour DFS; for an undirected one it is a DFS that
// ignores the edge it arrived on.
func (g *Graph[T]) HasCycle() bool { panic("TODO: implement HasCycle") }

// Dijkstra returns the cost of the cheapest path from start to every reachable
// node, and the predecessor map needed to rebuild the paths. start maps to 0.
// Unreachable nodes are absent from both maps.
//
// Use a priority queue. The O(V^2) version passes the correctness tests but
// fails the one with 50,000 nodes.
func (g *Graph[T]) Dijkstra(start T) (dist map[T]float64, prev map[T]T) {
	panic("TODO: implement Dijkstra")
}

// Path rebuilds the cheapest route from start to end using Dijkstra's output.
// It returns nil if end is unreachable.
func (g *Graph[T]) Path(start, end T) ([]T, float64) { panic("TODO: implement Path") }

// Components returns the connected components of an undirected graph, each
// sorted, and the components themselves sorted by their first element.
func (g *Graph[T]) Components() [][]T { panic("TODO: implement Components") }
