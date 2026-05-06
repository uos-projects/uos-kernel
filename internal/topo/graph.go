package topo

import (
	"sync"

	"github.com/uos-projects/uos-kernel/kernel"
)

type MemoryGraph struct {
	mu       sync.RWMutex
	outbound map[kernel.ResourceID][]kernel.Edge
	inbound  map[kernel.ResourceID][]kernel.Edge
}

func NewMemoryGraph() *MemoryGraph {
	return &MemoryGraph{
		outbound: make(map[kernel.ResourceID][]kernel.Edge),
		inbound:  make(map[kernel.ResourceID][]kernel.Edge),
	}
}

func (g *MemoryGraph) AddEdge(edge kernel.Edge) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.outbound[edge.From] = append(g.outbound[edge.From], edge)
	g.inbound[edge.To] = append(g.inbound[edge.To], edge)
}

func (g *MemoryGraph) RemoveEdge(from, to kernel.ResourceID, rel kernel.RelationType) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.outbound[from] = removeMatching(g.outbound[from], from, to, rel)
	g.inbound[to] = removeMatching(g.inbound[to], from, to, rel)
}

func removeMatching(edges []kernel.Edge, from, to kernel.ResourceID, rel kernel.RelationType) []kernel.Edge {
	n := 0
	for _, e := range edges {
		if e.From == from && e.To == to && e.Relation == rel {
			continue
		}
		edges[n] = e
		n++
	}
	return edges[:n]
}

func (g *MemoryGraph) Neighbors(id kernel.ResourceID, dir kernel.Direction, rel kernel.RelationType) []kernel.ResourceID {
	g.mu.RLock()
	defer g.mu.RUnlock()
	seen := make(map[kernel.ResourceID]struct{})
	var result []kernel.ResourceID

	if dir == kernel.Outbound || dir == kernel.Both {
		for _, e := range g.outbound[id] {
			if rel != "" && e.Relation != rel {
				continue
			}
			if _, ok := seen[e.To]; !ok {
				seen[e.To] = struct{}{}
				result = append(result, e.To)
			}
		}
	}
	if dir == kernel.Inbound || dir == kernel.Both {
		for _, e := range g.inbound[id] {
			if rel != "" && e.Relation != rel {
				continue
			}
			if _, ok := seen[e.From]; !ok {
				seen[e.From] = struct{}{}
				result = append(result, e.From)
			}
		}
	}
	return result
}

func (g *MemoryGraph) Traverse(start kernel.ResourceID, dir kernel.Direction, rel kernel.RelationType, maxDepth int) []kernel.ResourceID {
	g.mu.RLock()
	defer g.mu.RUnlock()

	visited := make(map[kernel.ResourceID]struct{})
	visited[start] = struct{}{}
	var result []kernel.ResourceID

	type item struct {
		id    kernel.ResourceID
		depth int
	}
	queue := []item{{start, 0}}

	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		if maxDepth > 0 && cur.depth >= maxDepth {
			continue
		}

		var neighbors []kernel.ResourceID
		if dir == kernel.Outbound || dir == kernel.Both {
			for _, e := range g.outbound[cur.id] {
				if rel != "" && e.Relation != rel {
					continue
				}
				neighbors = append(neighbors, e.To)
			}
		}
		if dir == kernel.Inbound || dir == kernel.Both {
			for _, e := range g.inbound[cur.id] {
				if rel != "" && e.Relation != rel {
					continue
				}
				neighbors = append(neighbors, e.From)
			}
		}

		for _, n := range neighbors {
			if _, ok := visited[n]; ok {
				continue
			}
			visited[n] = struct{}{}
			result = append(result, n)
			queue = append(queue, item{n, cur.depth + 1})
		}
	}
	return result
}

func (g *MemoryGraph) EdgesOf(id kernel.ResourceID) []kernel.Edge {
	g.mu.RLock()
	defer g.mu.RUnlock()
	var result []kernel.Edge
	result = append(result, g.outbound[id]...)
	result = append(result, g.inbound[id]...)
	return result
}
