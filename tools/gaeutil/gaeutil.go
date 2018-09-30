package gaeutil

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

// CollectModules returns a list of module names from diretroies under the deployment directory.
func CollectModules(deploymentDir string) ([]string, error) {
	files, err := ioutil.ReadDir(deploymentDir)
	if err != nil {
		return nil, fmt.Errorf("could not read the diretory %s: %v", err)
	}
	var list []string
	for _, f := range files {
		if f.IsDir() {
			yamlPath := filepath.Join(deploymentDir, f.Name(), "app.yaml")
			moduleName, err := extractModuleName(yamlPath)
			if err != nil {
				return nil, fmt.Errorf("could not extract module name in %s: %v", yamlPath, err)
			}
			if moduleName != "default" {
				list = append(list, moduleName)
			}
		}
	}
	return list, nil
}

var reModuleName = regexp.MustCompile("\\s*service:\\s*([^\\s#]+)")

func extractModuleName(yamlPath string) (string, error) {
	buff, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		return "", err
	}
	matches := reModuleName.Copy().FindSubmatch(buff)
	if matches != nil {
		return string(matches[1]), nil
	}
	return "", fmt.Errorf("no module name is defined in %s", yamlPath)
}
