package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	ent "github.com/speedland/go/ent/generator"
	"github.com/speedland/go/tools/generator"
	"github.com/speedland/go/x/xstrings"
)

const defaultOutput = "datastore_helper.go"

var (
	typeName = flag.String("type", "", "type name must be set")
	kindName = flag.String("kind", "", "kind name (same value of `type` by default)")
	output   = flag.String("output", "", "output file name; default srcdir/<type>_datastore.go.go")
)

func main() {
	log.SetPrefix("[ent] ")
	log.SetFlags(0)
	flag.Parse()
	if len(*typeName) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	if *kindName == "" {
		*kindName = *typeName
	}
	for _, directory := range args {
		var g *ent.Struct
		if *kindName == "" {
			g = ent.NewStruct(*typeName, *typeName)
		} else {
			g = ent.NewStruct(*typeName, *kindName)
		}
		src, err := generator.Run(directory, g)
		if err != nil {
			log.Printf("error: %v", err)
			if e, ok := err.(*generator.InvalidSourceError); ok {
				log.Printf("generated code:\n%s\n", e.SourceWithLine(false))
			}
			log.Fatalf("Exiting")
		}
		output := filepath.Join(
			directory,
			fmt.Sprintf("%s_datastore.go", xstrings.ToSnakeCase(*typeName)),
		)
		if err := ioutil.WriteFile(output, src, 0644); err != nil {
			log.Fatalf("writing output: %s", err)
		}
		log.Printf("Generated %s - %s", g.Type, output)
	}
}
