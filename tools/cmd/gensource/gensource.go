package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/yssk22/go/generator"
	"github.com/yssk22/go/generator/enum"
	api "github.com/yssk22/go/web/api/generator"
	"github.com/yssk22/go/gcp/datastore/typed"
	"github.com/yssk22/go/iterator/slice"
	"github.com/yssk22/go/x/xstrings"
)

var (
	annotation = flag.String("a", "", "annotation name to generate the sources")
)


func main() {
	log.SetPrefix("[gensource] ")
	log.SetFlags(0)
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		args = []string{"."}
	}
	log.Println("gensource", args)

	generators := []generator.Generator{
		api.NewGenerator(),
		enum.NewGenerator(),
		typed.NewGenerator(),
	}
	anns := xstrings.SplitAndTrim(*annotation, ",")
	generators = slice.Filter(generators, func(i int, g interface{}) bool{
		gena := g.(generator.Generator).GetAnnotationSymbol()
		if *annotation == "" {
			log.Printf("%s: yes\n", gena)
			return false
		}
		for _, a := range anns {
			if gena.Is(a) {
				log.Printf("%s: yes\n", gena)
				return false
			}	
		}
		log.Printf("%s: no\n", gena)
		return true
	}).([]generator.Generator)
	runner := generator.NewRunner(generators...)
	for _, dir := range args {
		runDirectory(runner, dir, false)
	}
}

func runDirectory(runner *generator.Runner, dir string, recursive bool) {
	filename := filepath.Base(dir)
	if filename == "..." {
		runDirectory(runner, filepath.Dir(dir), true)
		return
	}
	info, err := os.Stat(dir)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return
	}
	if !info.IsDir() {
		log.Printf("ERROR: %q is not a directory", dir)
		return
	}
	if recursive {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			log.Fatalf("FATAL: %s", err)
		}
		for _, file := range files {
			if file.IsDir() {
				runDirectory(runner, filepath.Join(dir, file.Name()), true)
			}
		}
	}
	err = runner.Run(dir)
	if err != nil {
		s := err.Error()
		if strings.Index(s, "no buildable Go source files") >= 0 {
			return
		}
		log.Printf("ERROR: %s", err)
	}
}
