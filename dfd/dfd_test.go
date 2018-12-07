package dfd

import (
	"fmt"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/graph"
)

type TestNode struct {
	*dotNode
}

func TestInitializeDFD(t *testing.T) {
	g_name := "my dfd"
	g := InitializeDFD(g_name)

	if g.Name != g_name {
		t.Errorf("Expected a DFD name of %s, but got %s", g_name, g.Name)
	}

	if g.graph == nil {
		t.Error("Graph attribute setter should not be nil")
	}

	if g.node == nil {
		t.Error("Node attribute setter should not be nil")
	}

	if g.edge == nil {
		t.Error("Edge attribute setter should not be nil")
	}
}

func TestInitializeTrustBoundary(t *testing.T) {
	g_name := "my dfd"
	g := InitializeTrustBoundary(g_name)

	if g.Name != g_name {
		t.Errorf("Expected a TrustBoundary name of %s, but got %s", g_name, g.Name)
	}

	if g.graph == nil {
		t.Error("Graph attribute setter should not be nil")
	}
}

func TestDeserializeDFD(t *testing.T) {
	cases := []struct {
		id string
	}{
		{"1234"},
		{"4321"},
		{"0"},
		{"-1"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Deserializing DFD with an id of %s", c.id), func(t *testing.T) {
			g := DeserializeDFD(c.id)
			if g.id != c.id {
				t.Errorf("Expected a DFD id of %s, but got %s", c.id, g.id)
			}
		})
	}
}

func TestDeserializeTrustBoundary(t *testing.T) {
	cases := []struct {
		id string
	}{
		{"1234"},
		{"4321"},
		{"0"},
		{"-1"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Deserializing TrustBoundary with an id of %s", c.id), func(t *testing.T) {
			g := DeserializeTrustBoundary(c.id)
			if g.id != c.id {
				t.Errorf("Expected a TrustBoundary id of %s, but got %s", c.id, g.id)
			}
		})
	}
}

func TestDFDAddNodeElem(t *testing.T) {
	cases := []struct {
		node         graph.Node
		expect_panic bool
	}{
		{node: NewProcess("")},
		{node: NewExternalService("")},
		{node: NewDataStore("")},
		{&TestNode{}, true},
	}

	g := InitializeDFD("")

	for _, c := range cases {
		t.Run(fmt.Sprintf("Adding node element of type %s to DFD", reflect.TypeOf(c.node)), func(t *testing.T) {
			if c.expect_panic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("The code did not panic when adding unknown node type")
					}
				}()
			}

			g.AddNodeElem(c.node)
			switch el := c.node.(type) {
			case *Process:
				if g.Processes[el.ExternalID()] == nil {
					t.Errorf("Expected to find a process with ID %d, but none exist", el.ID())
				}
			case *ExternalService:
				if g.ExternalServices[el.ExternalID()] == nil {
					t.Errorf("Expected to find an external service with ID %d, but none exist", el.ID())
				}
			case *DataStore:
				if g.DataStores[el.ExternalID()] == nil {
					t.Errorf("Expected to find a data store with ID %d, but none exist", el.ID())
				}
			}
		})
	}
}

func TestTrustBoundaryAddNodeElem(t *testing.T) {
	cases := []struct {
		node         graph.Node
		expect_panic bool
	}{
		{node: NewProcess("")},
		{node: NewExternalService("")},
		{node: NewDataStore("")},
		{&TestNode{}, true},
	}

	g := InitializeTrustBoundary("")

	for _, c := range cases {
		t.Run(fmt.Sprintf("Adding node element of type %s to TrustBoundary", reflect.TypeOf(c.node)), func(t *testing.T) {
			if c.expect_panic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("The code did not panic when adding unknown node type")
					}
				}()
			}

			g.AddNodeElem(c.node)
			switch el := c.node.(type) {
			case *Process:
				if g.Processes[el.ExternalID()] == nil {
					t.Errorf("Expected to find a process with ID %d, but none exist", el.ID())
				}
			case *ExternalService:
				if g.ExternalServices[el.ExternalID()] == nil {
					t.Errorf("Expected to find an external service with ID %d, but none exist", el.ID())
				}
			case *DataStore:
				if g.DataStores[el.ExternalID()] == nil {
					t.Errorf("Expected to find a data store with ID %d, but none exist", el.ID())
				}
			}
		})
	}
}

func TestDFDAddRemoveFlowRace(t *testing.T) {
	g := InitializeDFD("test")
	for i := 0; i < 5; i++ {
		go func() {
			p1 := NewProcess("p1")
			p2 := NewProcess("p2")
			g.AddFlow(p1, p2, "foo")
			g.RemoveFlow(p1.DOTID(), p2.DOTID())
		}()
	}
}

func TestDFDAddNodeRace(t *testing.T) {
	g := InitializeDFD("test")
	for i := 0; i < 5; i++ {
		go func() {
			p1 := NewProcess("p1")
			g.AddNodeElem(p1)
		}()
	}
}
