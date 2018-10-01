package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/yssk22/go/x/xerrors"
)

var (
	appName       = flag.String("application", "", "gae application name")
	packagePrefix = flag.String("package", "", "package prefix for each module")
	servicesDir   = flag.String("services", "./services", "the root directory of services")
	deploymentDir = flag.String("deployment", "./deployment", "deployment directory")
	outputDir     = flag.String("output", "./deployment/default/", "output file path of dispatch.yaml")
	fallback      = flag.String("fallback", "default", "fallback service name")
	singleService = flag.Bool("single", false, "true to have a single module rather than per-service modules")
)

func main() {
	log.SetPrefix("[gendeployment] ")
	log.SetFlags(0)
	flag.Parse()
	if len(*appName) == 0 || len(*packagePrefix) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	fallbackService, nonFallbackServices := collectServices(*servicesDir, *packagePrefix, *fallback)
	if fallbackService == nil {
		log.Fatalf("No fallback service (%q) is found", *fallback)
	}
	if !*singleService {
		for _, s := range nonFallbackServices {
			log.Printf("Creating %q deployment\n", s.Name)
			createDeployment(path.Join(*deploymentDir, s.Name), *appName, s)
		}
	}
	// create a fallback service deployment
	log.Printf("Creating %q deployment (fallback)\n", fallbackService.Name)
	createDeployment(path.Join(*deploymentDir, *fallback), *appName, fallbackService, nonFallbackServices...)
	if !*singleService {
		// dispatch.yaml
		log.Printf("Creating dispath.yaml\n")
		createDispatch(path.Join(*deploymentDir, *fallback), fallbackService, nonFallbackServices)
	}
}

func createDeployment(deploymentDir string, appName string, main *Service, services ...*Service) {
	os.MkdirAll(deploymentDir, 0755)

	tempdir, err := ioutil.TempDir("", "gendeployment")
	xerrors.MustNil(err)
	defer os.RemoveAll(tempdir)
	// create a generator go file
	var generatedFile *os.File
	var generatedFilePath = filepath.Join(tempdir, "main.go")
	generatedFile, err = os.Create(generatedFilePath)
	xerrors.MustNil(err)
	xerrors.MustNil(goGeneratorTemplate.Execute(generatedFile, &generatorTemplateVars{
		AppName:       appName,
		CronYamlPath:  filepath.Join(deploymentDir, "cron.yaml"),
		QueueYamlPath: filepath.Join(deploymentDir, "queue.yaml"),
		Services:      append([]*Service{main}, services...),
	}))
	xerrors.MustNil(generatedFile.Close())

	// execute generator
	var buff bytes.Buffer
	cmd := exec.Command("go", "run", generatedFilePath)
	cmd.Stderr = &buff
	if err := cmd.Run(); err != nil {
		log.Println("Failed to run the generator go file: ")
		fmt.Fprintln(os.Stderr, buff.String())
		log.Println("generator file source:")
		content, _ := ioutil.ReadFile(generatedFilePath)
		fmt.Fprintln(os.Stderr, string(content))
		log.Println("Please check your source is:")
		log.Println("  - package name must be the same name with the directory name")
		log.Println("  - package must export `func NewService() *service.Service` function")
		os.Exit(1)
	}

	// create app.go
	generatedFile, err = os.Create(filepath.Join(deploymentDir, "app.go"))
	xerrors.MustNil(err)
	xerrors.MustNil(goAppTemplate.Execute(generatedFile, &appTemplateVars{
		PackageName: main.PackageAlias,
		Services:    append([]*Service{main}, services...),
	}))
	xerrors.MustNil(generatedFile.Close())

	// create app.yaml
	generatedFile, err = os.Create(filepath.Join(deploymentDir, "app.yaml"))
	xerrors.MustNil(err)
	xerrors.MustNil(appYamlTemplate.Execute(generatedFile, &appYamlTemplateVars{
		ServiceName: main.Name,
		GoVersion:   "go1.8",
	}))
	xerrors.MustNil(generatedFile.Close())
}

func createDispatch(deploymentDir string, main *Service, services []*Service) {
	// dispatch.yaml
	var dispatchFilePath = filepath.Join(deploymentDir, "dispatch.yaml")
	dispatchFile, err := os.Create(dispatchFilePath)
	xerrors.MustNil(err)
	defer dispatchFile.Close()
	for _, s := range services {
		dispatchFile.WriteString(fmt.Sprintf("- url: \"*/%s/*\"\n", s.URL))
		dispatchFile.WriteString(fmt.Sprintf("  module: %s\n", s.Name))
	}
	dispatchFile.WriteString("- url: \"*/*\"\n")
	dispatchFile.WriteString(fmt.Sprintf("  module: %s\n", main.PackageAlias))
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
