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
		{Title: "broken", Value: r.Broken()},
		{Title: "bytesize", Value: r.Bytesize},
		{Title: "capacity", Value: r.Capacity},
		{Title: "default", Value: r.Default},
		{Title: "embedded", Value: r.Embedded()},
		{Title: "fd", Value: r.Fd},
		{Title: "frozen", Value: r.Frozen()},
		{Title: "fstring", Value: r.Fstring()},
		{Title: "generation", Value: r.Generation},
		{Title: "ivars", Value: r.Ivars},
		{Title: "length", Value: r.Length},
		{Title: "line", Value: r.Line},
		{Title: "marked", Value: r.GcMarked()},
		{Title: "memsize", Value: r.Memsize},
		{Title: "old", Value: r.GcOld()},
		{Title: "shared", Value: r.Shared()},
		{Title: "size", Value: r.Size},
		{Title: "wbProtected", Value: r.GcWbProtected()},
	}

	accString := func(title string, val string) {
		if val != "" {
			attr = append(attr, gexf.AttrValue{Title: title, Value: val})
		}
	}

	accInterface := func(title string, val interface{}) {
		if val != nil {
			attr = append(attr, gexf.AttrValue{Title: title, Value: val})
		}
	}

	accString("address", addr)
	accString("class", strconv.FormatUint(r.Class, 16))
	accString("encoding", r.Encoding)
	accString("file", r.File)
	accString("method", r.Method)
	accString("name", r.Name)
	accString("nodeType", r.NodeType)
	accString("struct", r.Struct)
	accString("type", label)
	accInterface("value", r.Value)

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
