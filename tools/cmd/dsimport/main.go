package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"golang.org/x/oauth2"

	"github.com/yssk22/go/tools/dsutil"
)

var (
	key       = flag.String("key", "", "service account key file")
	kind      = flag.String("kind", "", "kind to export")
	namespace = flag.String("namespace", "", "namespace on the kind")
	host      = flag.String("host", "", "appengine host name")
	output    = flag.String("input", "", "input file")
	appID     = flag.String("appid", "", "appid (needed if you import the data exported by another app)")
	skip      = flag.Int("skip", 0, "skip N rows (useful to resume import operation to recover errors)")
)

func main() {
	log.SetPrefix("[dsimport] ")
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
			*output = fmt.Sprintf("%s.%s.bk", *host, *kind)
		} else {
			*output = fmt.Sprintf("%s.%s.%s.bk", *host, *namespace, *kind)
		}
	}
	ctx, err := dsutil.GetRemoteContext(oauth2.NoContext, *host, *namespace, *key)
	if err != nil {
		log.Fatal(err)
	}
	f, err := os.Open(*output)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	option := &dsutil.ImportOption{}
	option.AppID = *appID
	option.Skip = *skip
	if n, err := dsutil.Import(ctx, *kind, f, option); err != nil {
		log.Fatalf("%v, you can resume using -skip=%d", err, *skip+n)
	}
}
