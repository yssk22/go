package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const defaultOutput = "dispatch.yaml"

var (
	appName       = flag.String("application", "", "gae application name")
	deploymentDir = flag.String("deployment", "./deployment", "deployment directory")
	outputPath    = flag.String("output", "./deployment/default/dispatch.yaml", "output file path of dispatch.yaml")
)

func main() {
	log.SetPrefix("[gendispatch] ")
	log.SetFlags(0)
	flag.Parse()
	if len(*appName) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	genDispatch(*appName, *deploymentDir, *outputPath)
}

func genDispatch(appName, deploymentDir, outputPath string) {
	files, err := ioutil.ReadDir(deploymentDir)
	if err != nil {
		panic(err)
	}
	var output io.Writer
	if outputPath == "-" {
		output = os.Stdout
	} else {
		outputFile, err := os.OpenFile(outputPath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()
		output = outputFile
	}
	output.Write([]byte(fmt.Sprintf("application: %s\n", appName)))
	output.Write([]byte(fmt.Sprintf("\n")))
	output.Write([]byte(fmt.Sprintf("dispatch:\n")))
	for _, f := range files {
		if f.IsDir() {
			moduleName, err := extractModuleName(filepath.Join(deploymentDir, f.Name(), "app.yaml"))
			if err != nil {
				panic(err)
			}
			if moduleName != "default" {
				modulePath := strings.Replace(moduleName, "-", "/", -1)
				output.Write([]byte(fmt.Sprintf("- url: \"*/%s/*\"\n", modulePath)))
				output.Write([]byte(fmt.Sprintf("  module: %s\n", moduleName)))
			}
		}
	}
	output.Write([]byte(fmt.Sprintf("- url: \"*/*\"\n")))
	output.Write([]byte(fmt.Sprintf(" module: default\n")))
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
