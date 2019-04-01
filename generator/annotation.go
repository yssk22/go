package generator

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/x/xerrors"
)

var (
	annotationSymbolRegexp = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_]*$")
	annotationRegexp       = regexp.MustCompile("\\s*(@[a-zA-Z0-9][a-zA-Z0-9_]*)")
)

// Annotation is a pair of AnnotationSymbol and it's params.
type Annotation struct {
	Symbol AnnotationSymbol
	Params map[string]interface{}
}

// AnnotationSymbol is a special string used for generator annotations and must comply with ^[a_zA_Z0-9][a_zA_Z0-9_]*$
type AnnotationSymbol string

// NewAnnotationSymbol returns a Annotation instance from a string
func NewAnnotationSymbol(s string) AnnotationSymbol {
	if !annotationSymbolRegexp.Copy().MatchString(s) {
		panic(fmt.Errorf("annotation must comply with ^[a_zA_Z0-9][a_zA_Z0-9_]*$ but %q", s))
	}
	return AnnotationSymbol(s)
}

// Is returns if the `s` matches with the annotaiton string.
func (a AnnotationSymbol) Is(s string) bool {
	return string(a) == s
}

// IsValid returns if the annotaiton is valid syntax or not
func (a AnnotationSymbol) IsValid() bool {
	return annotationRegexp.Copy().Match([]byte(a))
}

// MaybeMarkedIn returns if @annotation appears in the `contents`
func (a AnnotationSymbol) MaybeMarkedIn(contents []byte) bool {
	return regexp.MustCompile(fmt.Sprintf("\\s*@%s\\s*", a)).Match(contents)
}

// CollectAnnotatedNodes returns a list of node with annotations
func CollectAnnotatedNodes(p *PackageInfo) []*AnnotatedNode {
	var nodes []*AnnotatedNode
	var re = annotationRegexp.Copy()
	offset := 0
	for _, f := range p.Files {
		for node, commentGroups := range f.CommentMap {
			annotations := make(map[AnnotationSymbol]*Annotation)
			for _, c := range commentGroups {
				for _, line := range c.List {
					matched := re.FindAllStringSubmatchIndex(line.Text, -1)
					if len(matched) == 0 {
						continue
					}
					matchedIndex := matched[0][2]
					ann := parseAnnotation(line.Text[matchedIndex:])
					if _, ok := annotations[ann.Symbol]; ok {
						panic(fmt.Errorf("duplicated annotation %q", ann.Symbol))
					}
					annotations[ann.Symbol] = ann
				}
			}
			nodes = append(nodes, &AnnotatedNode{
				annotations: annotations,
				Node:        node,
				Source:      f,
			})
		}
		offset += len(f.Source)
	}
	return nodes
}

// AnnotatedNode represents a Node annotated by Annotation
type AnnotatedNode struct {
	Node        ast.Node
	Source      *FileInfo
	annotations map[AnnotationSymbol]*Annotation
}

// IsAnnotated returns the node is annotated by the symbol `s`
func (a *AnnotatedNode) IsAnnotated(s AnnotationSymbol) bool {
	_, ok := a.annotations[s]
	return ok
}

// GetParamsBy returns the key-vale pairs for annotation parameters.
func (a *AnnotatedNode) GetParamsBy(s AnnotationSymbol) keyvalue.Getter {
	if m, ok := a.annotations[s]; ok {
		return keyvalue.StringKeyMap(m.Params)
	}
	return nil
}

// GenError returns an wrapped error object with the node information
// n must be the ancestor of Signature node or nil
func (a *AnnotatedNode) GenError(e error, n ast.Node) error {
	i := a.Source.GetNodeInfo(a.Node)
	if n != nil {
		ii := a.Source.GetNodeInfo(n)
		return xerrors.Wrap(e, ii.String())
	}
	return xerrors.Wrap(e, i.String())
}

func parseAnnotation(s string) *Annotation {
	if !strings.HasPrefix(s, "@") {
		panic(fmt.Errorf("not an annotation (given %q)", s))
	}
	parts := strings.Split(s, " ")
	a := Annotation{
		Symbol: NewAnnotationSymbol(parts[0][1:]),
		Params: make(map[string]interface{}),
	}
	for i := 1; i < len(parts); i++ {
		var key, value string
		arg := parts[i]
		idx := strings.Index(arg, "=")
		if idx < 0 {
			key = s
			value = ""
		} else {
			key = arg[:idx]
			value = arg[idx+1:]
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		if key != "" {
			a.Params[key] = value
		}
	}
	return &a
}
