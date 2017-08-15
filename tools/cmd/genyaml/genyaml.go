package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"bytes"

	"github.com/speedland/go/tools/gaeutil"
	"github.com/speedland/go/x/xerrors"
)

var (
	appName       = flag.String("application", "", "gae application name")
	packageSuffix = flag.String("package", "", "package suffix for each module")
	deploymentDir = flag.String("deployment", "./deployment", "deployment directory")
	outputDir     = flag.String("output", "./deployment/default/", "output file path of dispatch.yaml")
)

func main() {
	log.SetPrefix("[genyaml] ")
	log.SetFlags(0)
	flag.Parse()
	if len(*appName) == 0 || len(*packageSuffix) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	genDispatch(*appName, *packageSuffix, *deploymentDir, *outputDir)
}

func genDispatch(appName, packageSuffix, deploymentDir, outputDir string) {
	modules, err := gaeutil.CollectModules(deploymentDir)
	xerrors.MustNil(err)
	dir, err := ioutil.TempDir("", "genyamltmp")
	xerrors.MustNil(err)
	defer os.RemoveAll(dir)
	bindings := &TemplateVar{
		AppName: appName,
	}
	bindings.CronYamlPath = filepath.Join(dir, "cron.yaml")
	bindings.QueueYamlPath = filepath.Join(dir, "queue.yaml")
	for _, m := range modules {
		module := &Module{
			Name:        m,
			URL:         strings.Replace(m, "-", "/", -1),
			Package:     m,
			PackagePath: filepath.Join(packageSuffix, m),
		}
		bindings.Modules = append(bindings.Modules, module)
		log.Printf(
			"Service: Name=%s, URL=%s, Package=%s, PackagePath=%s\n",
			module.Name, module.URL, module.Package, module.PackagePath,
		)
	}

	log.Println("generaing yaml files on each deployment directory...")
	for _, m := range bindings.Modules {
		generateConfigs(bindings.AppName, filepath.Join(deploymentDir, m.Name), m)
	}

	log.Println("generaing yaml files on defualt directory...")
	// dispatch.yaml
	var dispatchFilePath = filepath.Join(outputDir, "dispatch.yaml")
	var dispatchFile *os.File
	dispatchFile, err = os.Create(dispatchFilePath)
	xerrors.MustNil(err)
	defer dispatchFile.Close()
	xerrors.MustNil(dispatchFileTemplate.Execute(dispatchFile, bindings))
	log.Println("\t", dispatchFilePath)

	// configs for default
	generateConfigs(appName, outputDir, bindings.Modules...)
}

func generateConfigs(appName string, outputDir string, modules ...*Module) {
	// create a generator go file
	dir, err := ioutil.TempDir("", "genyamltmp")
	xerrors.MustNil(err)
	defer os.RemoveAll(dir)
	var goFile *os.File
	var goFilePath = filepath.Join(dir, "main.go")
	goFile, err = os.Create(goFilePath)
	xerrors.MustNil(err)
	defer goFile.Close()
	bindings := &TemplateVar{
		AppName: appName,
	}
	bindings.CronYamlPath = filepath.Join(dir, "cron.yaml")
	bindings.QueueYamlPath = filepath.Join(dir, "queue.yaml")
	bindings.Modules = modules
	xerrors.MustNil(goFileTemplate.Execute(goFile, bindings))

	// execute generator
	var buff bytes.Buffer
	cmd := exec.Command("go", "run", goFilePath)
	cmd.Stderr = &buff
	if err := cmd.Run(); err != nil {
		log.Println("Failed to run the generator go file: ")
		fmt.Fprintln(os.Stderr, buff.String())
		log.Println("generator file source:")
		content, _ := ioutil.ReadFile(goFilePath)
		fmt.Fprintln(os.Stderr, string(content))
		log.Println("Please check your source is:")
		log.Println("  - package name must be the same name with the directory name")
		log.Println("  - package must export `func NewService() *service.Service` function")
		os.Exit(1)
	}

	// copy the file to the final destination
	cronFilePath := filepath.Join(outputDir, "cron.yaml")
	queueFilePath := filepath.Join(outputDir, "queue.yaml")
	xerrors.MustNil(cp(cronFilePath, bindings.CronYamlPath))
	log.Println("\t", cronFilePath)
	xerrors.MustNil(cp(queueFilePath, bindings.QueueYamlPath))
	log.Println("\t", queueFilePath)
}

func cp(dst, src string) error {
	var err error
	var s, d *os.File
	if s, err = os.Open(src); err != nil {
		return err
	}
	defer s.Close()
	if d, err = os.Create(dst); err != nil {
		return err
	}
	defer d.Close()
	_, err = io.Copy(d, s)
	return err
}
