package dfd

import (
	"fmt"
	"strconv"

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
}

// Node circle
type Process struct {
	*dotNode
	Name string
}

// Node diamond
type ExternalService struct {
	*dotNode
	Name string
}

// Node cylinder
type DataStore struct {
	*dotNode
	Name string
}

// Subgraph
type TrustBoundary struct {
	*subGraph

	Name string

	Processes        map[string]*Process
	ExternalServices map[string]*ExternalService
	DataStores       map[string]*DataStore
}

// Edge
type Flow struct {
	*dotEdge
}

type DfdGraph interface {
	AddNodeElem(graph.Node)
	AddProcess(*Process) error
	AddExternalService(*ExternalService) error
	AddDataStore(*DataStore) error

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

// This method is used when loading a DFD from a DOT file, where an ID is already given
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

// This method is used when loading a DFD from a DOT file, where an ID is already given
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

// FIXME this should only conditionally generate a new ID
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

func DeserializeProcess(id string) *Process {
	xid64, _ := strconv.ParseInt(id, 10, 64)
	n := &Process{dotNode: &dotNode{Node: simple.Node(xid64)}}
	n.SetDOTID(id)
	return n
}

func NewProcess(name string) *Process {
	xid := genID()
	xid64, _ := strconv.ParseInt(xid, 10, 64)
	p := &Process{Name: name, dotNode: &dotNode{Node: simple.Node(xid64)}}
	p.SetDOTID(xid)
	p.Label = strconv.Quote(name)
	p.Shape = "circle"
	p.Name = name
	return p
}

func NewExternalService(name string) *ExternalService {
	xid := genID()
	xid64, _ := strconv.ParseInt(xid, 10, 64)
	es := &ExternalService{Name: name, dotNode: &dotNode{Node: simple.Node(xid64)}}
	es.SetDOTID(xid)
	es.Label = strconv.Quote(name)
	es.Shape = "diamond"
	es.Name = name
	return es
}

func DeserializeExternalService(id string) *ExternalService {
	xid64, _ := strconv.ParseInt(id, 10, 64)
	n := &ExternalService{dotNode: &dotNode{Node: simple.Node(xid64)}}
	n.SetDOTID(id)
	return n
}

func NewDataStore(name string) *DataStore {
	xid := genID()
	xid64, _ := strconv.ParseInt(xid, 10, 64)
	ds := &DataStore{Name: name, dotNode: &dotNode{Node: simple.Node(xid64)}}
	ds.SetDOTID(xid)
	ds.Label = strconv.Quote(name)
	ds.Shape = "cylinder"
	ds.Name = name
	return ds
}

func DeserializeDataStore(id string) *DataStore {
	xid64, _ := strconv.ParseInt(id, 10, 64)
	n := &DataStore{dotNode: &dotNode{Node: simple.Node(xid64)}}
	n.SetDOTID(id)
	return n
}

// FIXME This method assumes that all Nodes have unique IDs. This *should* be
// true, but further testing is required to really make this claim. It takes a
// brute force approach and runs in O(n). In the worst case, it will have to
// visit every node in the DFD
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

func (dfd *DataFlowDiagram) ExternalID() string {
	return dfd.id
}

func (dfd *DataFlowDiagram) UpdateName(new_name string) {
	dfd.Name = new_name
	dfd.setAttributes()
	return
}

func (g *DataFlowDiagram) AddNodeElem(n graph.Node) {
	switch el := n.(type) {
	case *Process:
		g.AddProcess(el)
	case *ExternalService:
		g.AddExternalService(el)
	case *DataStore:
		g.AddDataStore(el)
	default:
		panic(fmt.Sprintf("Unknown node type %T", g))
	}
	g.AddNode(n)
	return
}

func (g *DataFlowDiagram) AddProcess(p *Process) error {
	g.Processes[p.ExternalID()] = p
	return nil
}

func (g *DataFlowDiagram) RemoveProcess(id string) {
	delete(g.Processes, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *DataFlowDiagram) AddExternalService(es *ExternalService) error {
	g.ExternalServices[es.ExternalID()] = es
	return nil
}

func (g *DataFlowDiagram) RemoveExternalService(id string) {
	delete(g.ExternalServices, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *DataFlowDiagram) AddDataStore(es *DataStore) error {
	g.DataStores[es.ExternalID()] = es
	return nil
}

func (g *DataFlowDiagram) RemoveDataStore(id string) {
	delete(g.DataStores, id)
	g.RemoveNode(idToID64(id))
	return
}

func (dfd *DataFlowDiagram) AddTrustBoundary(name string) (*TrustBoundary, error) {
	tb := InitializeTrustBoundary(name)
	dfd.TrustBoundaries[tb.ExternalID()] = tb
	return dfd.TrustBoundaries[tb.ExternalID()], nil
}

func (g *DataFlowDiagram) RemoveTrustBoundary(id string) {
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
	switch el := n.(type) {
	case *Process:
		g.AddProcess(el)
	case *ExternalService:
		g.AddExternalService(el)
	case *DataStore:
		g.AddDataStore(el)
	default:
		panic(fmt.Sprintf("Unknown node type %T", g))
	}
	g.AddNode(n)
	return
}

func (g *TrustBoundary) AddProcess(p *Process) error {
	g.Processes[p.ExternalID()] = p
	return nil
}

func (g *TrustBoundary) RemoveProcess(id string) {
	delete(g.Processes, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *TrustBoundary) AddExternalService(es *ExternalService) error {
	g.ExternalServices[es.ExternalID()] = es
	return nil
}

func (g *TrustBoundary) RemoveExternalService(id string) {
	delete(g.ExternalServices, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *TrustBoundary) AddDataStore(es *DataStore) error {
	g.DataStores[es.ExternalID()] = es
	return nil
}

func (g *TrustBoundary) RemoveDataStore(id string) {
	delete(g.DataStores, id)
	g.RemoveNode(idToID64(id))
	return
}

func (g *DataFlowDiagram) AddFlow(f graph.Node, t graph.Node, name string) *Flow {
	flow := &Flow{dotEdge: &dotEdge{Label: formatFlowLabel(name), Edge: g.DirectedGraph.NewEdge(f, t)}}
	flow_id := genFlowID(f, t)
	g.SetEdge(flow)
	g.Flows[flow_id] = flow
	return flow
}

func (g *DataFlowDiagram) RemoveFlow(src_id, dest_id string) {
	delete(g.Flows, fmt.Sprintf("%s%s", src_id, dest_id))
	g.RemoveEdge(idToID64(src_id), idToID64(dest_id))
	return
}

func (p *Process) DOTID() string {
	return fmt.Sprintf("process_%s", p.dotID)
}

func (p *Process) ExternalID() string {
	return p.dotID
}

func (n *Process) UpdateName(new_name string) {
	n.Name = new_name
	n.Label = strconv.Quote(new_name)
	return
}

func (es *ExternalService) DOTID() string {
	return fmt.Sprintf("externalservice_%s", es.dotID)
}

func (n *ExternalService) UpdateName(new_name string) {
	n.Name = new_name
	n.Label = strconv.Quote(new_name)
	return
}

func (n *DataStore) UpdateName(new_name string) {
	n.Name = new_name
	n.Label = strconv.Quote(new_name)
	return
}

func (n *DataStore) ExternalID() string {
	return n.dotID
}

func (n *DataStore) DOTID() string {
	return fmt.Sprintf("datastore_%s", n.dotID)
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
