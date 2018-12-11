package dfd

import (
	"fmt"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
)

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
