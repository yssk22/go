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
	"github.com/yssk22/go/generator/flow"
	"github.com/yssk22/go/generator/api"
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

	flowOptions := flow.NewOptions()
	flowOptions.RootDir, _ = os.Getwd()
	generators := []generator.Generator{
		api.NewGenerator(),
		enum.NewGenerator(),
		flow.NewGenerator(flowOptions),		
	}
	if *annotation != "" {
		anns := xstrings.SplitAndTrim(*annotation, ",")
		generators = slice.Filter(generators, func(i int, g interface{}) bool{
			for _, a := range anns {
				if g.(generator.Generator).GetAnnotation().Is(a) {
					return false
				}	
			}
			return true
		}).([]generator.Generator)
	}
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
