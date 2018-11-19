package dfd

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gonum.org/v1/gonum/graph/encoding"
	"gonum.org/v1/gonum/graph/encoding/dot"
	fdot "gonum.org/v1/gonum/graph/formats/dot"
)

type Client struct {
	Config Config
	DFD    *DataFlowDiagram
}

func NewClient(dot_path string) *Client {
	client := &Client{
		Config: Config{
			DOTPath: dot_path,
		},
	}
	dfd, _ := client.DFDFromDOT()
	client.DFD = dfd.(*DataFlowDiagram)
	return client
}

func (client *Client) DFDFromDOT() (encoding.Builder, error) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	f, err := os.Open(client.Config.DOTPath)
	if os.IsNotExist(err) { // We'll initialize an empty DFD
		return InitializeDFD(""), nil
	} else if err != nil { // Bail if something else went wrong
		log.Fatal(err)
	}
	defer f.Close()

	fileinfo, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = f.Read(buffer)
	if err != nil {
		log.Fatal(err)
	}

	ast, err := fdot.ParseBytes(buffer)
	if err != nil { // Initialize an empty DFD if the file is malformed
		return InitializeDFD(""), nil
	}

	gast := ast.Graphs[0]
	dst := DeserializeDFD(gast.ID)
	gen := initGenerator(dst)
	for _, stmt := range gast.Stmts {
		gen.addStmt(dst, stmt)
	}
	return dst, nil
}

func (client *Client) DFDToDOT(dfd encoding.Builder) error {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	got, err := client.marshal(dfd)
	if err != nil {
		fmt.Println(err)
		return err
	}

	f, err := os.OpenFile(client.Config.DOTPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0660)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	nbytes, err := f.Write(got)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[INFO] %d bytes written\n", nbytes)

	fmt.Println(string(got))
	return nil
}

// Wrapper function for Marshal method in the dot package
func (client *Client) marshal(dfd encoding.Builder) ([]byte, error) {
	return dot.Marshal(dfd, "", "", "\t")
}
