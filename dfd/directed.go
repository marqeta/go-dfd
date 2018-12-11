package dfd

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/iterator"
	"gonum.org/v1/gonum/graph/simple"
)

// AddNode adds n to the graph. It panics if the added node ID matches an existing node ID.
func (g *dfdGraph) AddNode(n graph.Node) {
	if _, exists := g.nodes[n.ID()]; exists {
		panic(fmt.Sprintf("simple: node ID collision: %d", n.ID()))
	}
	g.nodes[n.ID()] = n
	g.from[n.ID()] = make(map[int64]graph.Edge)
	g.to[n.ID()] = make(map[int64]graph.Edge)
	g.nodeIDs.Use(n.ID())
}

// Edge returns the edge from u to v if such an edge exists and nil otherwise.
// The node v must be directly reachable from u as defined by the From method.
func (g *dfdGraph) Edge(uid, vid int64) graph.Edge {
	g.mtx.RLock()
	defer g.mtx.RUnlock()
	edge, ok := g.from[uid][vid]
	if !ok {
		return nil
	}
	return edge
}

// Edges returns all the edges in the graph.
func (g *dfdGraph) Edges() graph.Edges {
	var edges []graph.Edge
	for _, u := range g.nodes {
		for _, e := range g.from[u.ID()] {
			edges = append(edges, e)
		}
	}
	return iterator.NewOrderedEdges(edges)
}

// From returns all nodes in g that can be reached directly from n.
func (g *dfdGraph) From(id int64) graph.Nodes {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if _, ok := g.from[id]; !ok {
		return nil
	}

	from := make([]graph.Node, len(g.from[id]))
	i := 0
	for vid := range g.from[id] {
		from[i] = g.nodes[vid]
		i++
	}
	return iterator.NewOrderedNodes(from)
}

// HasEdgeBetween returns whether an edge exists between nodes x and y without
// considering direction.
func (g *dfdGraph) HasEdgeBetween(xid, yid int64) bool {
	if _, ok := g.from[xid][yid]; ok {
		return true
	}
	_, ok := g.from[yid][xid]
	return ok
}

// HasEdgeFromTo returns whether an edge exists in the graph from u to v.
func (g *dfdGraph) HasEdgeFromTo(uid, vid int64) bool {
	if _, ok := g.from[uid][vid]; !ok {
		return false
	}
	return true
}

// NewEdge returns a new Edge from the source to the destination node.
func (g *dfdGraph) NewEdge(from, to graph.Node) graph.Edge {
	return &simple.Edge{F: from, T: to}
}

// Node returns the node with the given ID if it exists in the graph,
// and nil otherwise.
func (g *dfdGraph) Node(id int64) graph.Node {
	return g.nodes[id]
}

// Nodes returns all the nodes in the graph.
func (g *dfdGraph) Nodes() graph.Nodes {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	if len(g.nodes) == 0 {
		return nil
	}
	nodes := make([]graph.Node, len(g.nodes))
	i := 0
	for _, n := range g.nodes {
		nodes[i] = n
		i++
	}
	return iterator.NewOrderedNodes(nodes)
}

// RemoveEdge removes the edge with the given end point IDs from the graph, leaving the terminal
// nodes. If the edge does not exist it is a no-op.
func (g *dfdGraph) RemoveEdge(fid, tid int64) {
	if _, ok := g.nodes[fid]; !ok {
		return
	}
	if _, ok := g.nodes[tid]; !ok {
		return
	}

	delete(g.from[fid], tid)
	delete(g.to[tid], fid)
}

// RemoveNode removes the node with the given ID from the graph, as well as any edges attached
// to it. If the node is not in the graph it is a no-op.
func (g *dfdGraph) RemoveNode(id int64) {
	if _, ok := g.nodes[id]; !ok {
		return
	}
	delete(g.nodes, id)

	for from := range g.from[id] {
		delete(g.to[from], id)
	}
	delete(g.from, id)

	for to := range g.to[id] {
		delete(g.from[to], id)
	}
	delete(g.to, id)

	g.nodeIDs.Release(id)
}

// SetEdge adds e, an edge from one node to another. If the nodes do not exist, they are added
// and are set to the nodes of the edge otherwise.
// It will panic if the IDs of the e.From and e.To are equal.
func (g *dfdGraph) SetEdge(e graph.Edge) {
	var (
		from = e.From()
		fid  = from.ID()
		to   = e.To()
		tid  = to.ID()
	)

	if fid == tid {
		panic("simple: adding self edge")
	}

	if _, ok := g.nodes[fid]; !ok {
		g.AddNode(from)
	} else {
		g.nodes[fid] = from
	}
	if _, ok := g.nodes[tid]; !ok {
		g.AddNode(to)
	} else {
		g.nodes[tid] = to
	}

	g.from[fid][tid] = e
	g.to[tid][fid] = e
}

// To returns all nodes in g that can reach directly to n.
func (g *dfdGraph) To(id int64) graph.Nodes {
	if _, ok := g.from[id]; !ok {
		return nil
	}

	to := make([]graph.Node, len(g.to[id]))
	i := 0
	for uid := range g.to[id] {
		to[i] = g.nodes[uid]
		i++
	}
	return iterator.NewOrderedNodes(to)
}
