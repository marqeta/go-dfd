package dfd

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

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
