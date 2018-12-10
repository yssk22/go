package main

import "html/template"

type generatorTemplateVars struct {
	AppName       string
	Services      []*Service
	CronYamlPath  string
	QueueYamlPath string
	IndexYamlPath string
	AppYamlPath   string
}

var goGeneratorTemplate = template.Must(template.New("gendeployment.gogenerator").Parse(`package main

import (
	"fmt"
	"os"
	{{range .Services -}}
	{{if .PackageAlias}}
	{{.PackageAlias}} "{{.PackagePath}}"
	{{else}}
	"{{.PackagePath}}"
	{{end}}
	{{end}}
)

func main(){
	var fcron, fqueue, findex, fapp *os.File
	var err error
	if fcron, err = os.Create("{{.CronYamlPath}}"); err != nil {
		panic(err)
	}
	defer fcron.Close()
	if fqueue, err = os.Create("{{.QueueYamlPath}}"); err != nil {
		panic(err)
	}
	defer fqueue.Close()
	if findex, err = os.Create("{{.IndexYamlPath}}"); err != nil {
		panic(err)
	}
	defer findex.Close()
	if fapp, err = os.Create("{{.AppYamlPath}}"); err != nil {
		panic(err)
	}
	defer fapp.Close()

	fmt.Fprintf(fcron, "cron:\n")
	fmt.Fprintf(fqueue, "queue:\n")
	fmt.Fprintf(findex, "indexes:\n")
	{{with index .Services 0}}
	{{if .PackageAlias}}
	{{.PackageAlias}}.NewService().GenAppYAML(fapp)
	{{else}}
	{{.Package}}.NewService().GenAppYAML(fapp)
	{{end}}
	{{end}}

	{{range $index, $element := .Services -}}
	func(){
		{{with $element}}
		{{if .PackageAlias}}
		s := {{.PackageAlias}}.NewService()
		{{else}}
		s := {{.Package}}.NewService()
		{{end}}
		s.GenCronYAML(fcron)
		s.GenQueueYAML(fqueue)
		s.GenIndexYAML(findex)
		{{if $index }}
		s.GenHandlersYAML(fapp)
		{{end}}
		{{end}}
	}()
	{{end}}
	// fallback
	fmt.Fprintf(fapp, "- url: /.*\n")
	fmt.Fprintf(fapp, "  secure: always\n")
	fmt.Fprintf(fapp, "  script: _go_app\n")
}
`))

type appTemplateVars struct {
	PackageName string
	Services    []*Service
}

var goAppTemplate = template.Must(template.New("gendeployment.goapp").Parse(`package main

import (
	"github.com/yssk22/go/gae/service"
	{{range .Services -}}
	{{if .PackageAlias -}}
	{{.PackageAlias}} "{{.PackagePath}}"
	{{else -}}
	"{{.PackagePath}}"
	{{end -}}
	{{end}}
	"google.golang.org/appengine"
)

func init(){
	service.NewDispatcher(
	{{range .Services -}}
	{{if .PackageAlias -}}
	{{.PackageAlias}}.NewService(),
	{{else -}}
	{{.Package}}.NewService(),
	{{end -}}
	{{end }}
	).Run()
	appengine.Main()
}
`))
