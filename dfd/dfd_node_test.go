package dfd

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

func TestNewProcess(t *testing.T) {
	cases := []struct {
		name string
	}{
		{"node 1"},
		{"node1"},
		{"node_1"},
		{""},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Creating a Process with a name of %s", c.name), func(t *testing.T) {
			n := NewProcess(c.name)
			if n.Name != c.name {
				t.Errorf("Expected Name to be %s, but got %s", c.name, n.Name)
			}

			if n.Label != strconv.Quote(c.name) {
				t.Errorf("Expected Label to be %s, but got %s", strconv.Quote(c.name), n.Name)
			}

			if n.Shape != Circle {
				t.Errorf("Expected Shape to be %s, but got %s", Circle, n.Shape)
			}
		})
	}
}

func TestNewExternalService(t *testing.T) {
	cases := []struct {
		name string
	}{
		{"node 1"},
		{"node1"},
		{"node_1"},
		{""},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Creating a ExternalService with a name of %s", c.name), func(t *testing.T) {
			n := NewExternalService(c.name)
			if n.Name != c.name {
				t.Errorf("Expected Name to be %s, but got %s", c.name, n.Name)
			}

			if n.Label != strconv.Quote(c.name) {
				t.Errorf("Expected Label to be %s, but got %s", strconv.Quote(c.name), n.Name)
			}

			if n.Shape != Diamond {
				t.Errorf("Expected Shape to be %s, but got %s", Diamond, n.Shape)
			}
		})
	}
}

func TestNewDataStore(t *testing.T) {
	cases := []struct {
		name string
	}{
		{"node 1"},
		{"node1"},
		{"node_1"},
		{""},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Creating a DataStore with a name of %s", c.name), func(t *testing.T) {
			n := NewDataStore(c.name)
			if n.Name != c.name {
				t.Errorf("Expected Name to be %s, but got %s", c.name, n.Name)
			}

			if n.Label != strconv.Quote(c.name) {
				t.Errorf("Expected Label to be %s, but got %s", strconv.Quote(c.name), n.Name)
			}

			if n.Shape != Cylinder {
				t.Errorf("Expected Shape to be %s, but got %s", Cylinder, n.Shape)
			}
		})
	}
}

func TestDeserializeProcess(t *testing.T) {
	cases := []struct {
		id string
	}{
		{"1234"},
		{"4321"},
		{"0"},
		{"-1"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Deserializing Process with an id of %s", c.id), func(t *testing.T) {
			n := DeserializeProcess(c.id)
			if n.ID() != idToID64(c.id) {
				t.Errorf("Expected node to have an ID of %s, but got %d", c.id, n.ID())
			}
		})
	}
}

func TestDeserializeExternalService(t *testing.T) {
	cases := []struct {
		id string
	}{
		{"1234"},
		{"4321"},
		{"0"},
		{"-1"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Deserializing ExternalService with an id of %s", c.id), func(t *testing.T) {
			n := DeserializeExternalService(c.id)
			if n.ID() != idToID64(c.id) {
				t.Errorf("Expected node to have an ID of %s, but got %d", c.id, n.ID())
			}
		})
	}
}

func TestDeserializeDataStore(t *testing.T) {
	cases := []struct {
		id string
	}{
		{"1234"},
		{"4321"},
		{"0"},
		{"-1"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Deserializing DataStore with an id of %s", c.id), func(t *testing.T) {
			n := DeserializeDataStore(c.id)
			if n.ID() != idToID64(c.id) {
				t.Errorf("Expected node to have an ID of %s, but got %d", c.id, n.ID())
			}
		})
	}
}

func TestNodeDOTID(t *testing.T) {
	cases := []struct {
		node      DfdNode
		node_type string
	}{
		{NewProcess(""), "process"},
		{NewExternalService(""), "externalservice"},
		{NewDataStore(""), "datastore"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Testing DOTID prefix for %s", c.node_type), func(t *testing.T) {
			actual := c.node.DOTID()
			if len(actual) < 1 || strings.Count(actual, c.node_type) != 1 {
				t.Errorf("Expected %s to contain %s", actual, c.node_type)
			}
		})
	}
}

func TestNodeExternalID(t *testing.T) {
	cases := []struct {
		node      DfdNode
		node_type string
	}{
		{NewProcess(""), "process"},
		{NewExternalService(""), "externalservice"},
		{NewDataStore(""), "datastore"},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("Testing ExternalID for %s", c.node_type), func(t *testing.T) {
			expected := reflect.ValueOf(c.node).Elem().FieldByName("dotID").String()
			actual := c.node.ExternalID()
			if expected != actual {
				t.Errorf("Expected %s, got %s", expected, actual)
			}
		})
	}
}

func TestProcessUpdateName(t *testing.T) {
	test_name := "TestNodeName"
	n := NewProcess("")
	if n.Name == test_name {
		t.Error("Somehow the test name was set during initialization!")
	}

	n.UpdateName(test_name)
	if n.Name != test_name {
		t.Errorf("Expected Name to be %s, but got %s", test_name, n.Name)
	}

	if n.Label != strconv.Quote(test_name) {
		t.Errorf("Expected Name to be %s, but got %s", strconv.Quote(test_name), n.Name)
	}
}

func TestExternalServiceUpdateName(t *testing.T) {
	test_name := "TestNodeName"
	n := NewExternalService("")
	if n.Name == test_name {
		t.Error("Somehow the test name was set during initialization!")
	}

	n.UpdateName(test_name)
	if n.Name != test_name {
		t.Errorf("Expected Name to be %s, but got %s", test_name, n.Name)
	}

	if n.Label != strconv.Quote(test_name) {
		t.Errorf("Expected Name to be %s, but got %s", strconv.Quote(test_name), n.Name)
	}
}

func TestDataStoreUpdateName(t *testing.T) {
	test_name := "TestNodeName"
	n := NewDataStore("")
	if n.Name == test_name {
		t.Error("Somehow the test name was set during initialization!")
	}

	n.UpdateName(test_name)
	if n.Name != test_name {
		t.Errorf("Expected Name to be %s, but got %s", test_name, n.Name)
	}

	if n.Label != strconv.Quote(test_name) {
		t.Errorf("Expected Name to be %s, but got %s", strconv.Quote(test_name), n.Name)
	}
}
