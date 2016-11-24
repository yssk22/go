package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/speedland/go/tools/dsutil"
)

var (
	kind      = flag.String("kind", "", "kind to export")
	namespace = flag.String("namespace", "", "namespace on the kind")
	host      = flag.String("host", "", "appengine host name")
	output    = flag.String("output", "", "output file")
)

func main() {
	log.SetPrefix("[dsexport] ")
	log.SetFlags(0)
	flag.Parse()
	if len(*kind) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if len(*host) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	if len(*output) == 0 {
		if len(*namespace) == 0 {
			*output = fmt.Sprintf("%s.%s.bk", *host, *namespace, *kind)
		} else {
			*output = fmt.Sprintf("%s.%s.%s.bk", *host, *namespace, *kind)
		}
	}
	ctx, err := dsutil.GetRemoteContext(*host, *namespace)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Create(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	if err := dsutil.Export(ctx, *kind, f); err != nil {
		log.Fatal(err)
	}
}