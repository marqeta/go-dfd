package dfd

import (
	//"crypto/rand"
	"fmt"
	"strconv"
	//"math"
	//"math/big"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/simple"
)

type structuredGraph struct {
	*simple.DirectedGraph
	id                string
	graph, node, edge attributes
}

// SetDOTID sets the DOT ID of the graph.
func (g *structuredGraph) SetDOTID(id string) {
	g.id = id
}

// DOTID returns the DOT ID of the graph.
func (g *structuredGraph) DOTID() string {
	return g.id
}

func (g *structuredGraph) NewNode() graph.Node {
	//xid, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
	xid := genID()
	xid64, _ := strconv.ParseInt(xid, 10, 64)
	return &dotNode{Node: simple.Node(xid64)}
}

// DOTAttributers implements the dot.Attributers interface.
func (g *structuredGraph) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return g.graph, g.node, g.edge
}

// DOTAttributeSetters implements the dot.AttributeSetters interface.
func (g *structuredGraph) DOTAttributeSetters() (graph, node, edge encoding.AttributeSetter) {
	return &g.graph, &g.node, &g.edge
}

type subGraph struct {
	*simple.DirectedGraph
	id                string
	graph, node, edge attributes
}

// SetDOTID sets the DOT ID of the graph.
func (g *subGraph) SetDOTID(id string) {
	g.id = id
}

func (g *subGraph) DOTID() string {
	return fmt.Sprintf("cluster_%s", g.id)
}

// NewNode returns a new node with a unique node ID for the graph.
func (g *subGraph) NewNode() graph.Node {
	xid := genID()
	xid64, _ := strconv.ParseInt(xid, 10, 64)
	return &dotNode{Node: simple.Node(xid64)}
}

// NewEdge returns a new Edge from the source to the destination node.
func (g *subGraph) NewEdge(from, to graph.Node) graph.Edge {
	return &dotEdge{Edge: g.DirectedGraph.NewEdge(from, to)}
}

// DOTAttributers implements the dot.Attributers interface.
func (g *subGraph) DOTAttributers() (graph, node, edge encoding.Attributer) {
	return g.graph, g.node, g.edge
}

// DOTAttributeSetters implements the dot.AttributeSetters interface.
func (g *subGraph) DOTAttributeSetters() (graph, node, edge encoding.AttributeSetter) {
	return &g.graph, &g.node, &g.edge
}

// dotNode extends simple.Node with a label field to test round-trip encoding
// and decoding of node DOT label attributes.
type dotNode struct {
	graph.Node
	dotID string
	// Node label.
	Label string
	Shape string
	Style string
	Dir   string
}

// DOTID returns the node's DOT ID.
//func (n *dotNode) DOTID() string {
//	return n.dotID
//}

func (n *dotNode) ExternalID() string {
	return n.dotID
}

// SetDOTID sets a DOT ID.
func (n *dotNode) SetDOTID(id string) {
	n.dotID = id
}

// SetAttribute sets a DOT attribute.
func (n *dotNode) SetAttribute(attr encoding.Attribute) error {
	switch attr.Key {
	case "label":
		n.Label = attr.Value
	case "shape":
		n.Shape = attr.Value
	case "style":
		n.Style = attr.Value
	case "dir":
		n.Dir = attr.Value
	default:
		return fmt.Errorf("unable to unmarshal node DOT attribute with key %q", attr.Key)
	}

	return nil
}

// Attributes returns the DOT attributes of the node.
func (n *dotNode) Attributes() []encoding.Attribute {
	if len(n.Label) == 0 && len(n.Shape) == 0 && len(n.Style) == 0 && len(n.Dir) == 0 {
		return nil
	}
	var attrs []encoding.Attribute
	if len(n.Label) != 0 {
		attrs = append(attrs, encoding.Attribute{Key: "label", Value: n.Label})
	}
	if len(n.Shape) != 0 {
		attrs = append(attrs, encoding.Attribute{Key: "shape", Value: n.Shape})
	}
	if len(n.Style) != 0 {
		attrs = append(attrs, encoding.Attribute{Key: "style", Value: n.Style})
	}
	if len(n.Dir) != 0 {
		attrs = append(attrs, encoding.Attribute{Key: "dir", Value: n.Dir})
	}
	return attrs
}

// dotEdge extends simple.Edge with a label field to test round-trip encoding and
// decoding of edge DOT label attributes.
type dotEdge struct {
	graph.Edge
	Dir            string
	Label          string
	FromPortLabels dotPortLabels
	ToPortLabels   dotPortLabels
}

// SetAttribute sets a DOT attribute.
func (e *dotEdge) SetAttribute(attr encoding.Attribute) error {
	switch attr.Key {
	case "label":
		e.Label = attr.Value
	case "dir":
		e.Dir = attr.Value
	default:
		return fmt.Errorf("unable to unmarshal edge DOT attribute with key %q", attr.Key)
	}
	return nil
}

func (e *dotEdge) SetFromPort(port, compass string) error {
	e.FromPortLabels.Port = port
	e.FromPortLabels.Compass = compass
	return nil
}

func (e *dotEdge) SetToPort(port, compass string) error {
	e.ToPortLabels.Port = port
	e.ToPortLabels.Compass = compass
	return nil
}

func (e *dotEdge) FromPort() (port, compass string) {
	return e.FromPortLabels.Port, e.FromPortLabels.Compass
}

func (e *dotEdge) ToPort() (port, compass string) {
	return e.ToPortLabels.Port, e.ToPortLabels.Compass
}

type attributes []encoding.Attribute

func (a attributes) Attributes() []encoding.Attribute {
	return []encoding.Attribute(a)
}
func (a *attributes) SetAttribute(attr encoding.Attribute) error {
	*a = append(*a, attr)
	return nil
}

type dotPortLabels struct {
	Port, Compass string
}
