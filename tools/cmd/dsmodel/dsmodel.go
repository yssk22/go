package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/speedland/go/tools/generator"
	"github.com/speedland/go/x/xstrings"
)

const defaultOutput = "datastore_helper.go"

var (
	typeName = flag.String("type", "", "type name must be set")
	output   = flag.String("output", "", "output file name; default srcdir/<type>_datastore.go.go")
)

func main() {
	log.SetPrefix("[enum] ")
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
	for _, directory := range args {
		g := &Generator{
			Type: *typeName,
		}
		src, err := generator.Run(directory, g)
		if err != nil {
			log.Printf("error: %v", err)
			log.Printf("generated code:\n%s\n", src)
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
