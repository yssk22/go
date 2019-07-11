package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/yssk22/go/ansi"
	"github.com/yssk22/go/number"
	"github.com/yssk22/go/x/xerrors"
)

// Generator is an interface to implement generator command
type Generator interface {
	Run(*PackageInfo, []*AnnotatedNode) ([]*Result, error)
	GetAnnotationSymbol() AnnotationSymbol
	GetFormatter() Formatter
}

// Result represents a result of Generator#Run
type Result struct {
	Filename string
	Source   string
}

func (gr *Result) write(dir string) (string, error) {
	filename := filepath.Join(
		dir,
		gr.Filename,
	)
	if err := ioutil.WriteFile(filename, []byte(gr.Source), 0644); err != nil {
		return "", xerrors.Wrap(err, "failed to write the generated source on %s", filename)
	}
	return filename, nil
}

// Runner is a struct to run generators
type Runner struct {
	generators []Generator
}

// NewRunner returns a *Runner
func NewRunner(generators ...Generator) *Runner {
	return &Runner{
		generators: generators,
	}
}

// Run executes Generator#run for all generators.
func (r *Runner) Run(dir string) error {
	hasAnnotations, err := r.hasAnnotations(dir)
	if err != nil {
		return err
	}
	if !hasAnnotations {
		return nil
	}

	log.Printf("INFO: parsing %s", dir)
	pkg, err := parsePackage(dir)
	if err != nil {
		absPath, _ := filepath.Abs(dir)
		return xerrors.Wrap(err, "failed to parse package %q", absPath)
	}
	annotatedNodes := CollectAnnotatedNodes(pkg)
	gerrors := xerrors.NewMultiError(len(r.generators))
	for i, g := range r.generators {
		var nodes []*AnnotatedNode
		symbol := g.GetAnnotationSymbol()
		for _, n := range annotatedNodes {
			if n.IsAnnotated(symbol) {
				nodes = append(nodes, n)
			}
		}
		if len(nodes) == 0 {
			continue
		}
		generated, err := g.Run(pkg, nodes)
		if err != nil {
			gerrors[i] = err
			continue
		}
		for _, result := range generated {
			result.Source, err = g.GetFormatter().Format(result.Source)
			if err != nil {
				return err
			}
			filename, err := result.write(dir)
			if err != nil {
				return err
			}
			log.Printf("INFO: Generated: %s (by %s)", ansi.Blue.Sprintf(filename), symbol)
		}
	}
	if gerrors.HasError() {
		return gerrors
	}
	return nil
}

func (r *Runner) hasAnnotations(dir string) (bool, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return false, err
	}
	for _, f := range files {
		if !f.IsDir() {
			buff, err := ioutil.ReadFile(path.Join(dir, f.Name()))
			if err != nil {
				return false, err
			}
			for _, g := range r.generators {
				if g.GetAnnotationSymbol().MaybeMarkedIn(buff) {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

// InvalidSourceError is an error if generated source is unable to compile.
// Use SourceWithLine() to debug the generated code.
type InvalidSourceError struct {
	Source string
	err    error
}

// SourceWithLine returns the source code with line numbers
func (e *InvalidSourceError) SourceWithLine(all bool) string {
	lines := strings.Split(e.Source, "\n")
	max := len(fmt.Sprintf("%d", len(lines)))
	format := fmt.Sprintf("%%%dd: %%s", max+1)
	errLineFormat := fmt.Sprintf("!%%%dd: %%s", max)
	n := number.ParseIntOr(strings.Split(e.err.Error(), ":")[0], 0)
	for i := range lines {
		if n != 0 && n == i+1 {
			lines[i] = fmt.Sprintf(errLineFormat, i+1, lines[i])
		} else {
			lines[i] = fmt.Sprintf(format, i+1, lines[i])
		}
	}
	if all || n == 0 {
		return strings.Join(lines, "\n")
	}
	n++
	if n < 5 {
		lines = lines[:n+4]
	} else if n < len(lines)-3 {
		lines = lines[n-6 : n+3]
	} else {
		lines = lines[n-6:]
	}
	return strings.Join(lines, "\n")
}

func (e *InvalidSourceError) Error() string {
	return fmt.Sprintf("%s\n%s", e.err.Error(), e.SourceWithLine(false))
}
