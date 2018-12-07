package dfd

import (
	"fmt"
	"strconv"
	"sync"

	"gonum.org/v1/gonum/graph"
	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	"gonum.org/v1/gonum/graph/simple"
)

// Graph
type DataFlowDiagram struct {
	*structuredGraph

	Name string

	Processes        map[string]*Process
	ExternalServices map[string]*ExternalService
	DataStores       map[string]*DataStore
	TrustBoundaries  map[string]*TrustBoundary
	Flows            map[string]*Flow

	mtx sync.Mutex
}

// Subgraph
type TrustBoundary struct {
	*subGraph

	Name string

	Processes        map[string]*Process
	ExternalServices map[string]*ExternalService
	DataStores       map[string]*DataStore

	mtx sync.Mutex
}

// Edge
type Flow struct {
	*dotEdge
}

type DfdGraph interface {
	AddNodeElem(graph.Node)
	addProcess(*Process) error
	addExternalService(*ExternalService) error
	addDataStore(*DataStore) error

	RemoveProcess(string)
	RemoveExternalService(string)
	RemoveDataStore(string)
}

func InitializeDFD(name string) *DataFlowDiagram {
	dfd := &DataFlowDiagram{
		Name:             name,
		Processes:        make(map[string]*Process),
		ExternalServices: make(map[string]*ExternalService),
		DataStores:       make(map[string]*DataStore),
		Flows:            make(map[string]*Flow),
		TrustBoundaries:  make(map[string]*TrustBoundary),
		structuredGraph: &structuredGraph{
			DirectedGraph: simple.NewDirectedGraph(), id: genID(),
		},
	}
	dfd.setAttributes()
	return dfd
}

// DeserializeDFD is used when loading a DFD from a DOT file, where an ID is already given
func DeserializeDFD(id string) *DataFlowDiagram {
	return &DataFlowDiagram{
		Processes:        make(map[string]*Process),
		ExternalServices: make(map[string]*ExternalService),
		DataStores:       make(map[string]*DataStore),
		Flows:            make(map[string]*Flow),
		TrustBoundaries:  make(map[string]*TrustBoundary),
		structuredGraph: &structuredGraph{
			DirectedGraph: simple.NewDirectedGraph(), id: id,
		},
	}
}

// DeserializeTrustBoundary is used when loading a DFD from a DOT file, where an ID is already given
func DeserializeTrustBoundary(id string) *TrustBoundary {
	return &TrustBoundary{
		Processes:        make(map[string]*Process),
		ExternalServices: make(map[string]*ExternalService),
		DataStores:       make(map[string]*DataStore),
		subGraph: &subGraph{
			DirectedGraph: simple.NewDirectedGraph(), id: id,
		},
	}
}

func InitializeTrustBoundary(name string) *TrustBoundary {
	tb := &TrustBoundary{
		Name:             name,
		Processes:        make(map[string]*Process),
		ExternalServices: make(map[string]*ExternalService),
		DataStores:       make(map[string]*DataStore),
		subGraph: &subGraph{
			DirectedGraph: simple.NewDirectedGraph(), id: genID(),
		},
	}
	tb.setAttributes()
	return tb
}

// FindNode looks for a node with a given id in either the top level graph or a
// subgraph. It assumes that all Nodes have unique IDs. This *should* be true,
// but further testing is required to really make this claim.
func (g *DataFlowDiagram) FindNode(id string) graph.Node {
	if n, ok := g.Processes[id]; ok {
		return n
	} else if n, ok := g.ExternalServices[id]; ok {
		return n
	} else if n, ok := g.DataStores[id]; ok {
		return n
	}
	for _, tb := range g.TrustBoundaries {
		if n := tb.FindNode(id); n != nil {
			return n
		}
	}
	return nil
}

// FindNode looks for a node with a given id It assumes that all Nodes have
// unique IDs. This *should* be true, but further testing is required to really
// make this claim.
func (g *TrustBoundary) FindNode(id string) graph.Node {
	if n, ok := g.Processes[id]; ok {
		return n
	} else if n, ok := g.ExternalServices[id]; ok {
		return n
	} else if n, ok := g.DataStores[id]; ok {
		return n
	}
	return nil
}

func (dfd *DataFlowDiagram) ToDOT() error {
	return nil
}

func (dfd *DataFlowDiagram) GetTrustBoundary(id string) *TrustBoundary {
	return dfd.TrustBoundaries[id]
}

func (dfd *DataFlowDiagram) ExternalID() string {
	return dfd.id
}

func (dfd *DataFlowDiagram) UpdateName(new_name string) {
	dfd.Name = new_name
	dfd.setAttributes()
	return
}

func (g *DataFlowDiagram) AddNodeElem(n graph.Node) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	switch el := n.(type) {
	case *Process:
		g.addProcess(el)
	case *ExternalService:
		g.addExternalService(el)
	case *DataStore:
		g.addDataStore(el)
	default:
		panic(fmt.Sprintf("Unknown node type %T", el))
	}
	g.AddNode(n)
	return
}

func (g *DataFlowDiagram) addProcess(p *Process) error {
	g.Processes[p.ExternalID()] = p
	return nil
}

func (g *DataFlowDiagram) RemoveProcess(id string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	delete(g.Processes, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *DataFlowDiagram) addExternalService(es *ExternalService) error {
	g.ExternalServices[es.ExternalID()] = es
	return nil
}

func (g *DataFlowDiagram) RemoveExternalService(id string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	delete(g.ExternalServices, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *DataFlowDiagram) addDataStore(es *DataStore) error {
	g.DataStores[es.ExternalID()] = es
	return nil
}

func (g *DataFlowDiagram) RemoveDataStore(id string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	delete(g.DataStores, id)
	g.RemoveNode(idToID64(id))
	return
}

func (dfd *DataFlowDiagram) AddTrustBoundary(name string) (*TrustBoundary, error) {
	dfd.mtx.Lock()
	defer dfd.mtx.Unlock()
	tb := InitializeTrustBoundary(name)
	dfd.TrustBoundaries[tb.ExternalID()] = tb
	return dfd.TrustBoundaries[tb.ExternalID()], nil
}

func (g *DataFlowDiagram) RemoveTrustBoundary(id string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	delete(g.TrustBoundaries, id)
	return
}

func (g *DataFlowDiagram) Structure() []dot.Graph {
	graphs := []dot.Graph{}
	for _, tb := range g.TrustBoundaries {
		graphs = append(graphs, tb)
	}
	return graphs
}

func (dfd *DataFlowDiagram) setAttributes() {
	graph_setter, node_setter, edge_setter := dfd.DOTAttributeSetters()
	dfd.graph = nil
	dfd.node = nil
	dfd.edge = nil
	for _, attr := range dfd.graphAttributes() {
		graph_setter.SetAttribute(attr)
	}

	for _, attr := range dfd.nodeAttributes() {
		node_setter.SetAttribute(attr)
	}

	for _, attr := range dfd.edgeAttributes() {
		edge_setter.SetAttribute(attr)
	}
	return
}

func (dfd *DataFlowDiagram) graphAttributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, 7)
	attrs[0] = makeAttribute("label", dfd.Name)
	attrs[1] = makeAttribute("fontname", "Arial")
	attrs[2] = makeAttribute("fontsize", "14")
	attrs[3] = makeAttribute("labelloc", "t")
	attrs[4] = makeAttribute("fontsize", "20")
	attrs[5] = makeAttribute("nodesep", "1")
	attrs[6] = makeAttribute("rankdir", "t")
	return attrs
}

func (dfd *DataFlowDiagram) nodeAttributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, 2)
	attrs[0] = makeAttribute("fontname", "Arial")
	attrs[1] = makeAttribute("fontsize", "14")
	return attrs
}

func (dfd *DataFlowDiagram) edgeAttributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, 3)
	attrs[0] = makeAttribute("shape", "none")
	attrs[1] = makeAttribute("fontname", "Arial")
	attrs[2] = makeAttribute("fontsize", "12")
	return attrs
}
func (sg *TrustBoundary) ExternalID() string {
	return sg.id
}

func (g *TrustBoundary) Subgraph() graph.Graph {
	return g
}

func (sg *TrustBoundary) setAttributes() {
	sg.graph = nil
	sg.node = nil
	sg.edge = nil
	graph_setter, _, _ := sg.DOTAttributeSetters()
	for _, attr := range sg.graphAttributes() {
		graph_setter.SetAttribute(attr)
	}
	return
}

func (sg *TrustBoundary) graphAttributes() []encoding.Attribute {
	attrs := make([]encoding.Attribute, 5)
	attrs[0] = makeAttribute("label", sg.Name)
	attrs[1] = makeAttribute("fontsize", "10")
	attrs[2] = makeAttribute("style", "dashed")
	attrs[3] = makeAttribute("color", "grey35")
	attrs[4] = makeAttribute("fontcolor", "grey35")
	return attrs
}

func (sg *TrustBoundary) UpdateName(new_name string) {
	sg.Name = new_name
	sg.setAttributes()
	return
}

func (g *TrustBoundary) AddNodeElem(n graph.Node) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	switch el := n.(type) {
	case *Process:
		g.addProcess(el)
	case *ExternalService:
		g.addExternalService(el)
	case *DataStore:
		g.addDataStore(el)
	default:
		panic(fmt.Sprintf("Unknown node type %T", g))
	}
	g.AddNode(n)
	return
}

func (g *TrustBoundary) addProcess(p *Process) error {
	g.Processes[p.ExternalID()] = p
	return nil
}

func (g *TrustBoundary) RemoveProcess(id string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	delete(g.Processes, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *TrustBoundary) addExternalService(es *ExternalService) error {
	g.ExternalServices[es.ExternalID()] = es
	return nil
}

func (g *TrustBoundary) RemoveExternalService(id string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	delete(g.ExternalServices, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *TrustBoundary) addDataStore(es *DataStore) error {
	g.DataStores[es.ExternalID()] = es
	return nil
}

func (g *TrustBoundary) RemoveDataStore(id string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	delete(g.DataStores, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *DataFlowDiagram) AddFlow(f graph.Node, t graph.Node, name string) *Flow {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	flow := &Flow{dotEdge: &dotEdge{Label: formatFlowLabel(name), Edge: g.DirectedGraph.NewEdge(f, t)}}
	flow_id := genFlowID(f, t)
	g.SetEdge(flow)
	g.Flows[flow_id] = flow
	return flow
}

func (g *DataFlowDiagram) RemoveFlow(src_id, dest_id string) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	delete(g.Flows, fmt.Sprintf("%s%s", src_id, dest_id))
	g.RemoveEdge(idToID64(src_id), idToID64(dest_id))
	return
}

func (f *Flow) Attributes() []encoding.Attribute {
	if len(f.Label) == 0 {
		return nil
	}
	return []encoding.Attribute{{
		Key:   "label",
		Value: f.Label,
	}}
}

func makeAttribute(key, value string) encoding.Attribute {
	return encoding.Attribute{Key: key, Value: strconv.Quote(value)}
}

func genFlowID(f graph.Node, t graph.Node) string {
	return fmt.Sprintf("%d%d", f.ID(), t.ID())
}

func formatFlowLabel(name string) string {
	return fmt.Sprintf(`<<table border="0" cellborder="0" cellpadding="2"><tr><td><b>%s</b></td></tr></table>>`, name)
}
