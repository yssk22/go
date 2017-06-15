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
	files, err := ioutil.ReadDir(deploymentDir)
	if err != nil {
		log.Fatal(err)
	}
	dir, err := ioutil.TempDir("", "genyamltmp")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	bindings := &TemplateVar{
		AppName: appName,
	}
	bindings.CronYamlPath = filepath.Join(dir, "cron.yaml")
	bindings.QueueYamlPath = filepath.Join(dir, "queue.yaml")
	for _, f := range files {
		if f.IsDir() {
			moduleName, err := extractModuleName(filepath.Join(deploymentDir, f.Name(), "app.yaml"))
			if err != nil {
				log.Fatal(err)
			}
			if moduleName != "default" {
				bindings.Modules = append(bindings.Modules, &Module{
					Name:        moduleName,
					URL:         strings.Replace(moduleName, "-", "/", -1),
					Package:     moduleName,
					PackagePath: filepath.Join(packageSuffix, moduleName),
				})
			}
		}
	}
	var goFilePath = filepath.Join(dir, "main.go")
	var dispatchFilePath = filepath.Join(outputDir, "dispatch.yaml")
	var cronFilePath = filepath.Join(outputDir, "cron.yaml")
	var queueFilePath = filepath.Join(outputDir, "queue.yaml")
	var goFile, dispatchFile *os.File
	log.Println("generaing yaml files...")
	dispatchFile, err = os.Create(dispatchFilePath)
	must(err)
	defer dispatchFile.Close()
	must(dispatchFileTemplate.Execute(dispatchFile, bindings))
	log.Println("\t", dispatchFilePath)

	goFile, err = os.Create(goFilePath)
	must(err)
	defer goFile.Close()
	must(goFileTemplate.Execute(goFile, bindings))

	// generate other yaml files by go command
	cmd := exec.Command("go", "run", goFilePath)
	// cmd.Stderr = os.Stderr
	must(cmd.Run())
	must(cp(cronFilePath, bindings.CronYamlPath))
	log.Println("\t", cronFilePath)
	must(cp(queueFilePath, bindings.QueueYamlPath))
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

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
