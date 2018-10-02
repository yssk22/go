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
	"path/filepath"

	"github.com/yssk22/go/x/xerrors"
)

var (
	appName       = flag.String("application", "", "gae application name")
	packagePrefix = flag.String("package", "", "package prefix for each module")
	servicesDir   = flag.String("services", "./services", "the root directory of services")
	deploymentDir = flag.String("deployment", "./deployment", "deployment directory")
	fallback      = flag.String("fallback", "default", "fallback service name")
	singleService = flag.Bool("single", false, "set to have a single module rather than per-service modules")
	startLocal    = flag.Bool("start-local", false, "set to start the local server by dev_appserver.py")
	cleanup       = flag.Bool("cleanup", true, "set to recreate deployment directory from scratch")
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
	if *cleanup {
		log.Printf("Cleanup %q\n", *deploymentDir)
		os.RemoveAll(*deploymentDir)
	}
	if !*singleService {
		for _, s := range nonFallbackServices {
			log.Printf("Creating %q deployment\n", s.Name)
			createDeployment(*deploymentDir, *appName, s)
		}
	}
	// create a fallback service deployment
	log.Printf("Creating %q deployment (fallback)\n", fallbackService.Name)
	createDeployment(*deploymentDir, *appName, fallbackService, nonFallbackServices...)
	if !*singleService {
		// dispatch.yaml
		log.Printf("Creating dispath.yaml\n")
		createDispatch(*deploymentDir, fallbackService, nonFallbackServices)
	}

	if *startLocal {
		log.Printf("Start dev_appserver.py\n")
		if *singleService {
			startLocalServer(*appName, *deploymentDir, fallbackService)
		} else {
			startLocalServer(*appName, *deploymentDir, fallbackService, nonFallbackServices...)
		}
	}
}

func createDeployment(deploymentDir string, appName string, main *Service, services ...*Service) {
	os.MkdirAll(filepath.Join(deploymentDir, main.Name), 0755)

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
		CronYamlPath:  filepath.Join(deploymentDir, main.Name, "cron.yaml"),
		QueueYamlPath: filepath.Join(deploymentDir, main.Name, "queue.yaml"),
		IndexYamlPath: filepath.Join(deploymentDir, main.Name, "index.yaml"),
		AppYamlPath:   filepath.Join(deploymentDir, main.Name, "app.yaml"),
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
	generatedFile, err = os.Create(filepath.Join(deploymentDir, main.Name, "app.go"))
	xerrors.MustNil(err)
	if main.PackageAlias != "" {
		xerrors.MustNil(goAppTemplate.Execute(generatedFile, &appTemplateVars{
			PackageName: main.PackageAlias,
			Services:    append([]*Service{main}, services...),
		}))
	} else {
		xerrors.MustNil(goAppTemplate.Execute(generatedFile, &appTemplateVars{
			PackageName: main.Package,
			Services:    append([]*Service{main}, services...),
		}))
	}
	xerrors.MustNil(generatedFile.Close())

	cmd = exec.Command("go", "fmt", "./"+filepath.Join(deploymentDir, main.Name))
	cmd.Stderr = &buff
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to run go fmt on %s\n", deploymentDir)
		fmt.Fprintln(os.Stderr, buff.String())
	}
}

func createDispatch(deploymentDir string, main *Service, services []*Service) {
	// dispatch.yaml
	var dispatchFilePath = filepath.Join(deploymentDir, main.Name, "dispatch.yaml")
	dispatchFile, err := os.Create(dispatchFilePath)
	xerrors.MustNil(err)
	defer dispatchFile.Close()
	dispatchFile.WriteString("dispatch:\n")
	for _, s := range services {
		dispatchFile.WriteString(fmt.Sprintf("  - url: \"*/%s/*\"\n", s.URL))
		dispatchFile.WriteString(fmt.Sprintf("    module: %s\n", s.Name))
	}
	dispatchFile.WriteString("  - url: \"*/*\"\n")
	dispatchFile.WriteString(fmt.Sprintf("    module: %s\n", main.Name))
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
