[![Build Status](https://cloud.drone.io/api/badges/marqeta/go-dfd/status.svg)](https://cloud.drone.io/marqeta/go-dfd)

# go-dfd

A utility written in Go for generating Data Flow Diagrams in DOT (Graphviz) format.

## Installation

```sh
$> go get github.com/marqeta/go-dfd/dfd
```
## Usage

```go
package main

import (
	dfd "github.com/marqeta/go-dfd/dfd"
)

func main() {
	client := dfd.NewClient("/path/to/dfd.dot")
	toDOT(client)
	fromDOT(client)
}

// You can write out Data Flow Diagram objects to DOT files
func toDOT(client *dfd.Client) {
	graph := dfd.InitializeDFD("WebApp Thing")
	p := dfd.NewProcess("Google Analytics")
	graph.AddNodeElem(p)

	external_tb, _ := graph.AddTrustBoundary("Browser")
	pclient := dfd.NewProcess("Client")
	external_tb.AddNodeElem(pclient)
	db := dfd.NewDataStore("sqlite")
	external_tb.AddNodeElem(db)
	graph.AddFlow(pclient, db, "HTTP")

	aws_tb, _ := graph.AddTrustBoundary("AWS")
	ws := dfd.NewProcess("Web Server")
	aws_tb.AddNodeElem(ws)
	logs := dfd.NewExternalService("Logs")
	aws_tb.AddNodeElem(logs)
	graph.AddFlow(ws, logs, "TCP")
	graph.AddFlow(pclient, ws, "HTTPS")
	graph.AddFlow(pclient, logs, "HTTPS")

	client.DFDToDOT(graph)
}

// You can read in DOT files as long as they follow the expected format
func fromDOT(client *dfd.Client) {
	client.DFDFromDOT()
}
```

The above code will generate a file at `/path/to/dfd.dot` which, when rendered with GraphViz, looks like the example provided below.

![scratch](https://user-images.githubusercontent.com/647423/49312152-dc269800-f4a8-11e8-825d-58c05144945f.png)
