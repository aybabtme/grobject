package main

import (
	"flag"
	"github.com/aybabtme/gexf"
	"github.com/aybabtme/rubyobj"
	"log"
	"os"
	"runtime"
	"strconv"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		objDump  string
		gexfDump string
	)
	flag.StringVar(&objDump, "src", "", "source of the object dump JSON")
	flag.StringVar(&gexfDump, "dst", "", "destination of the GEXF object dump")
	flag.Parse()

	switch {
	case objDump == "":
		flag.PrintDefaults()
		log.Fatalln("Need to provide source ObjectSpace file")
	case gexfDump == "":
		flag.PrintDefaults()
		log.Fatalln("Need to specify destination GEXF file")
	}

	jsonF, err := os.Open(objDump)
	if err != nil {
		log.Fatalf("Opening JSON source file: %v", err)
	}
	defer jsonF.Close()

	gexfF, err := os.Create(gexfDump)
	if err != nil {
		log.Fatalf("Opening GEXF destination file: %v", err)
	}
	defer gexfF.Close()

	decoded, errc := rubyobj.ParallelDecode(jsonF, uint(runtime.NumCPU()))

	go func() {
		for err := range errc {
			log.Printf("error: %v", err)
		}
	}()

	g := gexf.NewGraph()
	g.SetNodeAttrs(attrs)

	i := -1
	var id string

	for rObj := range decoded {
		addr, label, attr := extractNode(&rObj)

		if rObj.Address == 0 {
			id = g.GetID(i)
			i--
		} else {
			id = g.GetID(addr)
		}

		g.AddNode(id, label, attr)

		for _, ref := range rObj.References {
			g.AddEdge(id, g.GetID(ref))
		}
	}

	if err := gexf.Encode(gexfF, g); err != nil {
		log.Fatalf("Error encoding graph to GEXF: %v", err)
	}
}

func extractNode(r *rubyobj.RubyObject) (addr string, label string, attr []gexf.AttrValue) {
	addr = strconv.FormatUint(r.Address, 16)
	label = r.Type.Name()

	attr = []gexf.AttrValue{
		{Title: "type", Value: label},
		{Title: "value", Value: r.Value},
		{Title: "name", Value: r.Name},
		{Title: "nodeType", Value: r.NodeType},
		{Title: "address", Value: addr},
		{Title: "class", Value: r.Class},
		{Title: "default", Value: r.Default},
		{Title: "generation", Value: r.Generation},
		{Title: "bytesize", Value: r.Bytesize},
		{Title: "fd", Value: r.Fd},
		{Title: "file", Value: r.File},
		{Title: "encoding", Value: r.Encoding},
		{Title: "method", Value: r.Method},
		{Title: "ivars", Value: r.Ivars},
		{Title: "length", Value: r.Length},
		{Title: "line", Value: r.Line},
		{Title: "memsize", Value: r.Memsize},
		{Title: "capacity", Value: r.Capacity},
		{Title: "size", Value: r.Size},
		{Title: "struct", Value: r.Struct},
		{Title: "wbProtected", Value: r.GcWbProtected()},
		{Title: "old", Value: r.GcOld()},
		{Title: "marked", Value: r.GcMarked()},
		{Title: "broken", Value: r.Broken()},
		{Title: "frozen", Value: r.Frozen()},
		{Title: "fstring", Value: r.Fstring()},
		{Title: "shared", Value: r.Shared()},
		{Title: "embedded", Value: r.Embedded()},
	}
	return
}

var attrs = []gexf.Attr{
	{Title: "type", Type: gexf.String},
	{Title: "value", Type: gexf.String},
	{Title: "name", Type: gexf.String},
	{Title: "nodeType", Type: gexf.String},
	{Title: "address", Type: gexf.String},
	{Title: "class", Type: gexf.String},
	{Title: "default", Type: gexf.Long},
	{Title: "generation", Type: gexf.Long},
	{Title: "bytesize", Type: gexf.Long},
	{Title: "fd", Type: gexf.Long},
	{Title: "file", Type: gexf.String},
	{Title: "encoding", Type: gexf.String},
	{Title: "method", Type: gexf.String},
	{Title: "ivars", Type: gexf.Long},
	{Title: "length", Type: gexf.Long},
	{Title: "line", Type: gexf.Long},
	{Title: "memsize", Type: gexf.Long},
	{Title: "capacity", Type: gexf.Long},
	{Title: "size", Type: gexf.Long},
	{Title: "struct", Type: gexf.String},
	{Title: "wbProtected", Type: gexf.Boolean},
	{Title: "old", Type: gexf.Boolean},
	{Title: "marked", Type: gexf.Boolean},
	{Title: "broken", Type: gexf.Boolean},
	{Title: "frozen", Type: gexf.Boolean},
	{Title: "fstring", Type: gexf.Boolean},
	{Title: "shared", Type: gexf.Boolean},
	{Title: "embedded", Type: gexf.Boolean},
}
