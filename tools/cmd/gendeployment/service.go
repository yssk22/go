package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/yssk22/go/x/xerrors"
)

// Service is a definition of service
type Service struct {
	Name         string
	URL          string
	Package      string
	PackageAlias string
	PackagePath  string
}

func collectServices(servicesDir string, packagePrefix string, fallback string) (*Service, []*Service) {
	dirs, err := ioutil.ReadDir(servicesDir)
	xerrors.MustNil(err)
	var list []*Service
	for _, d := range dirs {
		if !d.IsDir() {
			panic(fmt.Errorf("%s is not a directory", d.Name()))
		}
		name := d.Name()
		s := &Service{
			Name:        name,
			URL:         strings.Replace(name, "-", "/", -1),
			Package:     name,
			PackagePath: filepath.Join(packagePrefix, name),
		}
		if name == "default" {
			s.PackageAlias = "def"
		}
		list = append(list, s)
		log.Printf(
			"Service: Name=%s, URL=%s, Package=%s, PackagePath=%s\n",
			s.Name, s.URL, s.Package, s.PackagePath,
		)
	}
	var fallbackService *Service
	var nonFallbackServices []*Service
	for _, s := range list {
		if s.Name != fallback {
			nonFallbackServices = append(nonFallbackServices, s)
		} else {
			fallbackService = s
		}
	}
	return fallbackService, nonFallbackServices
}
