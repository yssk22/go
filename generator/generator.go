package generator

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/yssk22/go/x/xerrors"

	"github.com/yssk22/go/number"
)

// Generator is an interface to implement generator command
type Generator interface {
	Run(*PackageInfo) ([]*Result, error)
}

// Result represents a result of Generator#Run
type Result struct {
	Filename string
	Source   string
}

func (gr *Result) write(dir string) (string, error) {
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
	pkg, err := parsePackage(dir)
	if err != nil {
		absPath, _ := filepath.Abs(dir)
		return xerrors.Wrap(err, "failed to parse package %q", absPath)
	}
	for _, g := range r.generators {
		generated, err := g.Run(pkg)
		if err != nil {
			return err
		}
		for _, result := range generated {
			filename, err := result.write(dir)
			if err != nil {
				return err
			}
			log.Printf("INFO: Generated: %s", filename)
		}
	}
	return nil
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
	return e.err.Error()
}
