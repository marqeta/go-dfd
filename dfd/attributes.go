package dfd

import (
	"gonum.org/v1/gonum/graph/encoding"
)

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
