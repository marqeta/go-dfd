package dfd

import (
	"fmt"
	"strings"
	"unsafe"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/formats/dot/ast"
)

func initGenerator(dst encoding.Builder) *generator {
	gen := &generator{directed: true, ids: make(map[string]graph.Node)}
	if a, ok := dst.(dot.AttributeSetters); ok {
		gen.graphAttr, gen.nodeAttr, gen.edgeAttr = a.DOTAttributeSetters()
	}
	return gen
}

// A generator keeps track of the information required for generating a gonum
// graph from a dot AST graph.
type generator struct {
	// Directed graph.
	directed bool
	// Map from dot AST node ID to gonum node.
	ids map[string]graph.Node
	// Nodes processed within the context of a subgraph, that is to be used as a
	// vertex of an edge.
	subNodes []graph.Node
	// Stack of start indices into the subgraph node slice. The top element
	// corresponds to the start index of the active (or inner-most) subgraph.
	subStart []int
	// graphAttr, nodeAttr and edgeAttr are global graph attributes.
	graphAttr, nodeAttr, edgeAttr encoding.AttributeSetter
}

// node returns the gonum node corresponding to the given dot AST node ID,
// generating a new such node if none exist.
func (gen *generator) node(dst encoding.Builder, id string) graph.Node {
	var ntype string
	var n graph.Node
	node_id_obj := strings.Split(id, "_")
	if len(node_id_obj) != 2 {
		panic(fmt.Errorf("Malformed Node ID: %s", id))
	}

	ntype = node_id_obj[0]
	id = node_id_obj[1]

	// bail if already visited
	// FIXME subgraphs need to be accounted for
	if n, ok := gen.ids[id]; ok {
		return n
	}

	switch ntype {
	case "process":
		n = DeserializeProcess(id)
	case "externalservice":
		n = DeserializeExternalService(id)
	case "datastore":
		n = DeserializeDataStore(id)
	default:
		panic(fmt.Sprintf("unknown node type %s", ntype))
	}

	if g, ok := dst.(DfdGraph); ok {
		g.AddNodeElem(n)
	} else {
		dst.AddNode(n)
	}
	gen.ids[id] = n
	return n
}

// addStmt adds the given statement to the graph.
func (gen *generator) addStmt(dst encoding.Builder, stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.NodeStmt:
		n, ok := gen.node(dst, stmt.Node.ID).(encoding.AttributeSetter)
		if !ok {
			return
		}
		for _, attr := range stmt.Attrs {
			a := encoding.Attribute{
				Key:   attr.Key,
				Value: attr.Val,
			}
			if err := n.SetAttribute(a); err != nil {
				panic(fmt.Errorf("unable to unmarshal node DOT attribute (%s=%s)", a.Key, a.Value))
			}
		}
	case *ast.EdgeStmt:
		gen.addEdgeStmt(dst, stmt)
	case *ast.AttrStmt:
		var n encoding.AttributeSetter
		var dst string
		switch stmt.Kind {
		case ast.GraphKind:
			if gen.graphAttr == nil {
				return
			}
			n = gen.graphAttr
			dst = "graph"
		case ast.NodeKind:
			if gen.nodeAttr == nil {
				return
			}
			n = gen.nodeAttr
			dst = "node"
		case ast.EdgeKind:
			if gen.edgeAttr == nil {
				return
			}
			n = gen.edgeAttr
			dst = "edge"
		default:
			panic("unreachable")
		}
		for _, attr := range stmt.Attrs {
			a := encoding.Attribute{
				Key:   attr.Key,
				Value: attr.Val,
			}
			if err := n.SetAttribute(a); err != nil {
				panic(fmt.Errorf("unable to unmarshal global %s DOT attribute (%s=%s)", dst, a.Key, a.Value))
			}
		}
	case *ast.Attr:
		// ignore.
	case *ast.Subgraph:
		tb_id := strings.Replace(stmt.ID, "cluster_", "", -1)
		sub := DeserializeTrustBoundary(tb_id)
		dst.(*DataFlowDiagram).TrustBoundaries[tb_id] = sub
		next_gen := initGenerator(sub)
		for _, stmt := range stmt.Stmts {
			next_gen.addStmt(sub, stmt)
		}
	default:
		panic(fmt.Sprintf("unknown statement type %T", stmt))
	}
}

// isInSubgraph reports whether the active context is within a subgraph, that is
// to be used as a vertex of an edge.
func (gen *generator) isInSubgraph() bool {
	return len(gen.subStart) > 0
}

// appendSubgraphNode appends the given node to the slice of nodes processed
// within the context of a subgraph.
func (gen *generator) appendSubgraphNode(n graph.Node) {
	gen.subNodes = append(gen.subNodes, n)
}

// addEdgeStmt adds the given edge statement to the graph.
func (gen *generator) addEdgeStmt(dst encoding.Builder, stmt *ast.EdgeStmt) {
	fs := gen.addVertex(dst, stmt.From)
	ts := gen.addEdge(dst, stmt.To, stmt.Attrs)
	for _, f := range fs {
		for _, t := range ts {
			edge := dst.(*DataFlowDiagram).AddFlow(f, t, "")
			applyPortsToEdge(stmt.From, stmt.To, edge)
			addEdgeAttrs(edge, stmt.Attrs)
		}
	}
}

// addVertex adds the given vertex to the graph, and returns its set of nodes.
func (gen *generator) addVertex(dst encoding.Builder, v ast.Vertex) []graph.Node {
	switch v := v.(type) {
	case *ast.Node:
		n := gen.node(dst, v.ID)
		return []graph.Node{n}
	case *ast.Subgraph:
		gen.pushSubgraph()
		for _, stmt := range v.Stmts {
			gen.addStmt(dst, stmt)
		}
		return gen.popSubgraph()
	default:
		panic(fmt.Sprintf("unknown vertex type %T", v))
	}
}

// addEdge adds the given edge to the graph, and returns its set of nodes.
func (gen *generator) addEdge(dst encoding.Builder, to *ast.Edge, attrs []*ast.Attr) []graph.Node {
	if !gen.directed && to.Directed {
		panic(fmt.Errorf("directed edge to %v in undirected graph", to.Vertex))
	}
	fs := gen.addVertex(dst, to.Vertex)
	if to.To != nil {
		ts := gen.addEdge(dst, to.To, attrs)
		for _, f := range fs {
			for _, t := range ts {
				edge := dst.NewEdge(f, t)
				dst.SetEdge(edge)
				applyPortsToEdge(to.Vertex, to.To, edge)
				addEdgeAttrs(edge, attrs)
			}
		}
	}
	return fs
}

// applyPortsToEdge applies the available port metadata from an ast.Edge
// to a graph.Edge
func applyPortsToEdge(from ast.Vertex, to *ast.Edge, edge graph.Edge) {
	if ps, isPortSetter := edge.(dot.PortSetter); isPortSetter {
		if n, vertexIsNode := from.(*ast.Node); vertexIsNode {
			if n.Port != nil {
				err := ps.SetFromPort(n.Port.ID, n.Port.CompassPoint.String())
				if err != nil {
					panic(fmt.Errorf("unable to unmarshal edge port (:%s:%s)", n.Port.ID, n.Port.CompassPoint.String()))
				}
			}
		}

		if n, vertexIsNode := to.Vertex.(*ast.Node); vertexIsNode {
			if n.Port != nil {
				err := ps.SetToPort(n.Port.ID, n.Port.CompassPoint.String())
				if err != nil {
					panic(fmt.Errorf("unable to unmarshal edge DOT port (:%s:%s)", n.Port.ID, n.Port.CompassPoint.String()))
				}
			}
		}
	}
}

// pushSubgraph pushes the node start index of the active subgraph onto the
// stack.
func (gen *generator) pushSubgraph() {
	gen.subStart = append(gen.subStart, len(gen.subNodes))
}

// popSubgraph pops the node start index of the active subgraph from the stack,
// and returns the nodes processed since.
func (gen *generator) popSubgraph() []graph.Node {
	// Get nodes processed since the subgraph became active.
	start := gen.subStart[len(gen.subStart)-1]
	// TODO: Figure out a better way to store subgraph nodes, so that duplicates
	// may not occur.
	nodes := unique(gen.subNodes[start:])
	// Remove subgraph from stack.
	gen.subStart = gen.subStart[:len(gen.subStart)-1]
	if len(gen.subStart) == 0 {
		// Remove subgraph nodes when the bottom-most subgraph has been processed.
		gen.subNodes = gen.subNodes[:0]
	}
	return nodes
}

// unique returns the set of unique nodes contained within ns.
func unique(ns []graph.Node) []graph.Node {
	var nodes []graph.Node
	seen := make(Int64s)
	for _, n := range ns {
		id := n.ID()
		if seen.Has(id) {
			// skip duplicate node
			continue
		}
		seen.Add(id)
		nodes = append(nodes, n)
	}
	return nodes
}

type Int64s map[int64]struct{}

// The simple accessor methods for Ints are provided to allow ease of
// implementation change should the need arise.

// Add inserts an element into the set.
func (s Int64s) Add(e int64) {
	s[e] = struct{}{}
}

// Has reports the existence of the element in the set.
func (s Int64s) Has(e int64) bool {
	_, ok := s[e]
	return ok
}

// Remove deletes the specified element from the set.
func (s Int64s) Remove(e int64) {
	delete(s, e)
}

// Count reports the number of elements stored in the set.
func (s Int64s) Count() int {
	return len(s)
}

// Int64sEqual reports set equality between the parameters. Sets are equal if
// and only if they have the same elements.
func Int64sEqual(a, b Int64s) bool {
	if int64sSame(a, b) {
		return true
	}

	if len(a) != len(b) {
		return false
	}

	for e := range a {
		if _, ok := b[e]; !ok {
			return false
		}
	}

	return true
}

// int64sSame determines whether two sets are backed by the same store. In the
// current implementation using hash maps it makes use of the fact that
// hash maps are passed as a pointer to a runtime Hmap struct. A map is
// not seen by the runtime as a pointer though, so we use unsafe to get
// the maps' pointer values to compare.
func int64sSame(a, b Int64s) bool {
	return *(*uintptr)(unsafe.Pointer(&a)) == *(*uintptr)(unsafe.Pointer(&b))
}

// addEdgeAttrs adds the attributes to the given edge.
func addEdgeAttrs(edge graph.Edge, attrs []*ast.Attr) {
	e, ok := edge.(encoding.AttributeSetter)
	if !ok {
		return
	}
	for _, attr := range attrs {
		a := encoding.Attribute{
			Key:   attr.Key,
			Value: attr.Val,
		}
		if err := e.SetAttribute(a); err != nil {
			panic(fmt.Errorf("unable to unmarshal edge DOT attribute (%s=%s)", a.Key, a.Value))
		}
	}
}
