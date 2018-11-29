package dfd

import (
	"strings"
	"testing"
)

// This can be thought of as an integration test
func TestClientDFDToDOT(t *testing.T) {
	client := &Client{
		Config: Config{
			DOTPath: "test.dot",
		},
	}
	dfd := DeserializeDFD("1552575689497326632")
	dfd.UpdateName("WebApp Thing")
	client.DFD = dfd

	google := DeserializeProcess("4404728580455388596")
	google.UpdateName("Google Analytics")
	dfd.AddNodeElem(google)

	external_tb := DeserializeTrustBoundary("2377452644169062617")
	external_tb.UpdateName("Browser")
	dfd.TrustBoundaries["2377452644169062617"] = external_tb
	pclient := DeserializeProcess("6384522904477046688")
	pclient.UpdateName("Client")
	external_tb.AddNodeElem(pclient)

	aws_tb := DeserializeTrustBoundary("7626181850182627084")
	aws_tb.UpdateName("AWS")
	dfd.TrustBoundaries["7626181850182627084"] = aws_tb
	ws := DeserializeProcess("6865082864924295608")
	ws.UpdateName("Web Server")
	aws_tb.AddNodeElem(ws)
	logs := DeserializeExternalService("4258120822598301454")
	logs.UpdateName("Logs")
	aws_tb.AddNodeElem(logs)

	dfd.AddFlow(aws_tb.Processes[ws.ExternalID()], aws_tb.ExternalServices[logs.ExternalID()], "TCP")
	dfd.AddFlow(pclient, ws, "HTTPS")
	dfd.AddFlow(pclient, logs, "HTTPS")

	dot, err := client.DFDToDOT(dfd)

	if err != nil {
		t.Error("An error occurred while generating the DOT file")
	}

	if strings.TrimSpace(dot) != strings.TrimSpace(testGraph) {
		t.Error("Returned DOT does not match the expected graph")
	}
}

const testGraph = `
strict digraph 1552575689497326632 {
	graph [
		label="WebApp Thing"
		fontname="Arial"
		fontsize="14"
		labelloc="t"
		fontsize="20"
		nodesep="1"
		rankdir="t"
	];
	node [
		fontname="Arial"
		fontsize="14"
	];
	edge [
		shape="none"
		fontname="Arial"
		fontsize="12"
	];

	subgraph cluster_2377452644169062617 {
		graph [
			label="Browser"
			fontsize="10"
			style="dashed"
			color="grey35"
			fontcolor="grey35"
		];

		// Node definitions.
		process_6384522904477046688 [
			label="Client"
			shape=circle
		];
	}
	subgraph cluster_7626181850182627084 {
		graph [
			label="AWS"
			fontsize="10"
			style="dashed"
			color="grey35"
			fontcolor="grey35"
		];

		// Node definitions.
		externalservice_4258120822598301454 [
			label="Logs"
			shape=diamond
		];
		process_6865082864924295608 [
			label="Web Server"
			shape=circle
		];
	}
	// Node definitions.
	externalservice_4258120822598301454 [
		label="Logs"
		shape=diamond
	];
	process_4404728580455388596 [
		label="Google Analytics"
		shape=circle
	];
	process_6384522904477046688 [
		label="Client"
		shape=circle
	];
	process_6865082864924295608 [
		label="Web Server"
		shape=circle
	];

	// Edge definitions.
	process_6384522904477046688 -> externalservice_4258120822598301454 [label=<<table border="0" cellborder="0" cellpadding="2"><tr><td><b>HTTPS</b></td></tr></table>>];
	process_6384522904477046688 -> process_6865082864924295608 [label=<<table border="0" cellborder="0" cellpadding="2"><tr><td><b>HTTPS</b></td></tr></table>>];
	process_6865082864924295608 -> externalservice_4258120822598301454 [label=<<table border="0" cellborder="0" cellpadding="2"><tr><td><b>TCP</b></td></tr></table>>];
}
`
