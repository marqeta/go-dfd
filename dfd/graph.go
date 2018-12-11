package dfd

import (
	"strconv"
	"sync"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/simple"
)

type dfdGraph struct {
	id                string
	graph, node, edge attributes
	nodes             map[int64]graph.Node
	from              map[int64]map[int64]graph.Edge
	to                map[int64]map[int64]graph.Edge

	nodeIDs Set

	mtx sync.RWMutex
}

func NewDfdGraph() *dfdGraph {
	return &dfdGraph{
		id:    genID(),
		nodes: make(map[int64]graph.Node),
		from:  make(map[int64]map[int64]graph.Edge),
		to:    make(map[int64]map[int64]graph.Edge),

		nodeIDs: NewSet(),
	}
}

// SetDOTID sets the DOT ID of the graph.
func (g *dfdGraph) SetDOTID(id string) {
	g.id = id
}

func (g *dfdGraph) NewNode() graph.Node {
	//xid, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	xid := genID()
	xid64, _ := strconv.ParseInt(xid, 10, 64)
	return &dotNode{Node: simple.Node(xid64)}
}

// DOTAttributers implements the dot.Attributers interface.
func (g *dfdGraph) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return g.graph, g.node, g.edge
}

// DOTAttributeSetters implements the dot.AttributeSetters interface.
func (g *dfdGraph) DOTAttributeSetters() (graph, node, edge encoding.AttributeSetter) {
	return &g.graph, &g.node, &g.edge
}
