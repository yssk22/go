package generator

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/yssk22/go/x/xerrors"
)

var (
	annotationRegexp = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_]*$")
)

// Annotation is a special string used for generator annotations and must comply with ^[a_zA_Z0-9][a_zA_Z0-9_]*$
type Annotation struct {
	str string
}

// NewAnnotation returns a Annotation instance from a string
func NewAnnotation(s string) *Annotation {
	if !annotationRegexp.Copy().MatchString(s) {
		panic(fmt.Errorf("annotation must comply with ^[a_zA_Z0-9][a_zA_Z0-9_]*$ but %q", s))
	}
	return &Annotation{
		str: s,
	}
}

// Is returns if the `s` matches with the annotaiton string.
func (a *Annotation) Is(s string) bool {
	return a.str == s
}

func (a *Annotation) String() string {
	return fmt.Sprintf("@%s", a.str)
}

// IsValid returns if the annotaiton is valid syntax or not
func (a *Annotation) IsValid() bool {
	return annotationRegexp.Copy().Match([]byte(a.str))
}

// MaybeMarkedIn returns if @annotation appears in the `contents`
func (a *Annotation) MaybeMarkedIn(contents []byte) bool {
	return regexp.MustCompile(fmt.Sprintf("\\s*@%s\\s*", a.str)).Match(contents)
}

// Collect returns a list of annotated signature
func (a *Annotation) Collect(p *PackageInfo) []*AnnotatedNode {
	var nodes []*AnnotatedNode
	var re = regexp.MustCompile(fmt.Sprintf("\\s*@%s\\s*", a.str))
	for _, f := range p.Files {
		for node, commentGroups := range f.CommentMap {
			for _, c := range commentGroups {
				for _, line := range c.List {
					idx := re.FindStringIndex(line.Text)
					if len(idx) > 0 {
						commentFlag := strings.Index("#", line.Text)
						if commentFlag > 0 && commentFlag < idx[0] {
							continue
						}
						remains := line.Text[idx[0]+len(a.str)+1:]
						params := parseSignatureParams(remains)
						nodes = append(nodes, &AnnotatedNode{
							annotation: a,
							Params:     params,
							Node:       node,
							Source:     f,
						})
					}
				}
			}
		}
	}
	return nodes
}

// AnnotatedNode represents a Node annotated by Annotation
type AnnotatedNode struct {
	annotation *Annotation
	Params     map[string]string
	Node       ast.Node
	Source     *FileInfo
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

func parseSignatureParams(s string) map[string]string {
	params := make(map[string]string)
	arguments := strings.Split(s, " ")
	for _, arg := range arguments {
		var key, value string
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
			params[key] = value
		}
	}
	return params
}
