package main

import "html/template"

type TemplateVar struct {
	AppName       string
	Modules       []*Module
	CronYamlPath  string
	QueueYamlPath string
}

type Module struct {
	Name        string
	URL         string
	Package     string
	PackagePath string
}

var goFileTemplate = template.Must(template.New("genyaml").Parse(`package main

import (
	"fmt"
	"os"
	{{range .Modules -}}
	"{{.PackagePath}}"
	{{end}}
)

func main(){
	var fcron, fqueue *os.File
	var err error
	if fcron, err = os.Create("{{.CronYamlPath}}"); err != nil {
		panic(err)
	}
	defer fcron.Close()
	if fqueue, err = os.Create("{{.QueueYamlPath}}"); err != nil {
		panic(err)
	}
	defer fqueue.Close()
	fmt.Fprintf(fcron, "cron:\n")
	fmt.Fprintf(fqueue, "queue:\n")
	{{range .Modules -}}
	func(){
		s := {{.Package}}.NewService()
		s.GenCronYAML(fcron)
		s.GenQueueYAML(fqueue)
	}()
	{{end}}
}
`))

var dispatchFileTemplate = template.Must(template.New("dispatch").Parse(`dispatch:
{{range .Modules -}}
- url: "*/{{.URL}}/*"
  module: {{.Name}}
{{end -}}
- url: "*/*"
  module: default
`))
