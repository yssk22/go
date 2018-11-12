package generator

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/yssk22/go/ansi"
	"github.com/yssk22/go/number"
	"github.com/yssk22/go/x/xerrors"
)

// ResultFileType is a type of the result file
type ResultFileType int

// Available ResultFileType
const (
	ResultFileTypeGo ResultFileType = iota
	ResultFileTypeFlow
)

// Generator is an interface to implement generator command
type Generator interface {
	Run(*PackageInfo, []*AnnotatedNode) ([]*Result, error)
	GetAnnotation() *Annotation
}

// Result represents a result of Generator#Run
type Result struct {
	Filename string
	Source   string
	FileType ResultFileType
}

func (gr *Result) write(dir string) (string, error) {
	switch gr.FileType {
	case ResultFileTypeGo:
		return gr.writeGo(dir)
	case ResultFileTypeFlow:
		return gr.writeFlow(dir)
	default:
		panic(fmt.Errorf("unknown file type value %d", gr.FileType))
	}
}

func (gr *Result) writeGo(dir string) (string, error) {
	formatted, err := format.Source([]byte(gr.Source))
	if err != nil {
		return "", &InvalidSourceError{
			Source: gr.Source,
			err:    err,
		}
	}
	filename := filepath.Join(
		dir,
		gr.Filename,
	)
	if err = ioutil.WriteFile(filename, formatted, 0644); err != nil {
		return "", xerrors.Wrap(err, "failed to write the generated source on %s", filename)
	}
	return filename, nil
}

func (gr *Result) writeFlow(dir string) (string, error) {
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
	hasAnnotated, err := r.hasAnnotated(dir)
	if err != nil {
		return err
	}
	if !hasAnnotated {
		return nil
	}
	log.Printf("INFO: parsing %s", dir)
	pkg, err := parsePackage(dir)
	if err != nil {
		absPath, _ := filepath.Abs(dir)
		return xerrors.Wrap(err, "failed to parse package %q", absPath)
	}
	for _, g := range r.generators {
		nodes := g.GetAnnotation().Collect(pkg)
		generated, err := g.Run(pkg, nodes)
		if err != nil {
			return err
		}
		for _, result := range generated {
			filename, err := result.write(dir)
			if err != nil {
				return err
			}
			log.Printf("INFO: Generated: %s", ansi.Blue.Sprintf(filename))
		}
	}
	return nil
}

func (r *Runner) hasAnnotated(dir string) (bool, error) {
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
				if g.GetAnnotation().MaybeMarkedIn(buff) {
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
