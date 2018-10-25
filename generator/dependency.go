package generator

import (
	"bytes"
	"fmt"
	"path"
)

// Dependency object provides a utility to manage dependencies for the generated source file.
type Dependency struct {
	packages       []string
	packageToAlias map[string]string
	aliasToPackage map[string]string
}

// NewDependency returns a new *Dependency object
func NewDependency() *Dependency {
	return &Dependency{
		packages:       []string{},
		packageToAlias: map[string]string{},
		aliasToPackage: map[string]string{},
	}
}

// Add a `pkg` into the dependency
func (d *Dependency) Add(pkg string) string {
	alias, ok := d.packageToAlias[pkg]
	if ok {
		return alias
	}
	alias = path.Base(pkg)
	defer func() {
		d.packages = append(d.packages, pkg)
		d.packageToAlias[pkg] = alias
		d.aliasToPackage[alias] = pkg
	}()

	base := alias
	c := 0
	for {
		_, ok := d.aliasToPackage[alias]
		if !ok {
			return alias
		}
		c++
		alias = fmt.Sprintf("%s%d", base, c)
	}
}

// GenImport generates import statements
func (d *Dependency) GenImport() string {
	var buff bytes.Buffer
	buff.WriteString("import (\n")
	for _, pkg := range d.packages {
		base := path.Base(pkg)
		alias := d.packageToAlias[pkg]
		if alias == base {
			buff.WriteString(fmt.Sprintf("\t%q\n", pkg))
		} else {
			buff.WriteString(fmt.Sprintf("\t%s %q\n", alias, pkg))
		}
	}
	buff.WriteString(")\n")
	return buff.String()
}
