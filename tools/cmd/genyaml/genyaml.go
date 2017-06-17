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
	"regexp"
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
	log.SetPrefix("[gendispatch] ")
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
	var goFilePath = filepath.Join(dir, "main.go")
	var dispatchFilePath = filepath.Join(outputDir, "dispatch.yaml")
	var cronFilePath = filepath.Join(outputDir, "cron.yaml")
	var queueFilePath = filepath.Join(outputDir, "queue.yaml")
	var goFile, dispatchFile *os.File
	log.Println("generaing yaml files...")
	dispatchFile, err = os.Create(dispatchFilePath)
	xerrors.MustNil(err)
	defer dispatchFile.Close()
	xerrors.MustNil(dispatchFileTemplate.Execute(dispatchFile, bindings))
	log.Println("\t", dispatchFilePath)

	goFile, err = os.Create(goFilePath)
	xerrors.MustNil(err)
	defer goFile.Close()
	xerrors.MustNil(goFileTemplate.Execute(goFile, bindings))

	// generate other yaml files by go command
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
	xerrors.MustNil(cp(cronFilePath, bindings.CronYamlPath))
	log.Println("\t", cronFilePath)
	xerrors.MustNil(cp(queueFilePath, bindings.QueueYamlPath))
	log.Println("\t", queueFilePath)
}

var reModuleName = regexp.MustCompile("\\s+module:\\s*([^\\s#]+)")

func extractModuleName(yamlPath string) (string, error) {
	buff, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return "", err
	}
	matches := reModuleName.FindSubmatch(buff)
	if matches != nil {
		return string(matches[1]), nil
	}
	return "", fmt.Errorf("no module name is defined in %s", yamlPath)
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
