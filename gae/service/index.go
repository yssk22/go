package service

import (
	"bytes"
	"fmt"
	"io"
)

// Index is a composit index definition for services to generate index.yaml
type Index struct {
	Kind       string
	Ancestor   bool
	Properties []*IndexProperty
}

// IndexProperty is a struct to represent name and direction
type IndexProperty struct {
	Name string
	Desc bool
}

// ToYAML returns a generated YAML content
func (i *Index) ToYAML() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "- kind: %s\n", i.Kind)
	if i.Ancestor {
		fmt.Fprintf(&buff, "  ancestor: yes\n")
	}
	fmt.Fprintf(&buff, "  properties:\n")
	for _, p := range i.Properties {
		fmt.Fprintf(&buff, "  - name: %s\n", p.Name)
		if p.Desc {
			fmt.Fprintf(&buff, "    direction: desc\n")
		}
	}
	return buff.String()
}

// AddIndex adds the index definition on the service
func (s *Service) AddIndex(kind string, ancestor bool, properties []*IndexProperty) {
	i := &Index{
		Kind:       kind,
		Ancestor:   ancestor,
		Properties: properties,
	}
	s.indexes = append(s.indexes, i)
}

// GenIndexYaml generates index yaml content to `w`
func (s *Service) GenIndexYAML(w io.Writer) {
	fmt.Fprintf(w, "# Service -- %s\n", s.Key())
	for _, i := range s.indexes {
		fmt.Fprintf(w, "%s", i.ToYAML())
	}
}

// GetIndexes returns a list of indexes defined in the service
func (s *Service) GetIndexes() []*Index {
	return s.indexes
}
