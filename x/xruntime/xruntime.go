// Package xruntime provides extended utility functions for runtime
package xruntime

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

// Frame is a stack frame
type Frame struct {
	ShortFilePath string
	FullFilePath  string
	LineNumber    int
	PackageName   string
	FunctionName  string
	pc            uintptr
}

// CaptureCaller returns a frame where a caller is.
func CaptureCaller() *Frame {
	return CaptureStackFrom(2, 1)[0]
}

// CaptureStack returns a list of stack frames as []*Frame
func CaptureStack(maxDepth int) []*Frame {
	return captureFrames(1, maxDepth)
}

// CaptureStackFrom is like CaptureStack but skip the given number of frames
func CaptureStackFrom(skip int, maxDepth int) []*Frame {
	return captureFrames(1+skip, maxDepth)
}

// CaptureFrame to capture the current stack frame
func CaptureFrame() *Frame {
	return captureFrames(1, 1)[0]
}

func (f *Frame) String() string {
	if f.PackageName == "" {
		return fmt.Sprintf("%s (at %s#%d)", f.FunctionName, f.ShortFilePath, f.LineNumber)
	}
	return fmt.Sprintf("%s.%s (at %s#%d)", f.PackageName, f.FunctionName, f.ShortFilePath, f.LineNumber)
}

func captureFrames(skip int, maxDepth int) []*Frame {
	counters := make([]uintptr, maxDepth)
	stack := make([]*Frame, maxDepth)
	runtime.Callers(2+skip, counters)
	for i, pc := range counters {
		f := runtime.FuncForPC(pc)
		if f == nil {
			if i > 0 {
				return stack[:(i - 1)]
			}
			return make([]*Frame, 0)
		}
		frame := &Frame{
			pc: pc,
		}
		frame.FullFilePath, frame.LineNumber = f.FileLine(pc)
		frame.PackageName, frame.FunctionName = getPackageAndFunction(f)
		if idx := strings.LastIndex(frame.FullFilePath, frame.PackageName); idx >= 0 {
			frame.ShortFilePath = frame.FullFilePath[idx:]
		} else {
			// source is in module so lookup go.mod file
			moduleRoot, moduleName := lookupGoModuleInfoFromFilePath(frame.FullFilePath)
			s := strings.Replace(frame.FullFilePath, moduleRoot, moduleName, 1)
			if idx := strings.LastIndex(s, frame.PackageName); idx >= 0 {
				frame.ShortFilePath = s[idx:]
			} else {
				frame.ShortFilePath = frame.FullFilePath
			}
		}
		stack[i] = frame
	}
	return stack
}

var (
	moduleDefRe = regexp.MustCompile("module\\s+(\\S+)\n")
)

func lookupGoModuleInfoFromFilePath(path string) (string, string) {
	dirname, err := filepath.Abs(filepath.Dir(path))
	if err != nil {
		panic(err)
	}
	gomod := filepath.Join(dirname, "go.mod")
	_, err = os.Stat(gomod)
	if err != nil {
		if os.IsNotExist(err) {
			if dirname == "/" {
				return "", ""
			}
			return lookupGoModuleInfoFromFilePath(filepath.Join(dirname, "..") + "/")
		}
		panic(err)
	}
	contents, err := ioutil.ReadFile(gomod)
	if err != nil {
		panic(err)
	}
	found := moduleDefRe.Copy().FindSubmatch(contents)
	if len(found) == 0 {
		panic(fmt.Errorf("could not find module declaration in %s", gomod))
	}
	return dirname, string(found[1])
}

func getPackageAndFunction(f *runtime.Func) (string, string) {
	// we'll see example.com/package/name.Struct.Func by *runtime.Func
	// so return example.com/package/name and Struct.Func
	fullName := f.Name()
	if lastSlashAt := strings.LastIndex(fullName, "/"); lastSlashAt >= 0 {
		pkgPrefix := fullName[:lastSlashAt]
		shortName := fullName[lastSlashAt:]
		if firstDotAt := strings.Index(shortName, "."); firstDotAt >= 0 {
			return fmt.Sprintf("%s%s", pkgPrefix, shortName[:firstDotAt]), shortName[firstDotAt+1:]
		}
	}
	return "", fullName
}

// CollectAllStacksSimple returns a list of stack frame with "{source}:{line}" format
func CollectAllStacksSimple() []string {
	var stack []string
	c := 0
	for {
		_, src, line, ok := runtime.Caller(c)
		if !ok {
			return stack[1:]
		}
		stack = append(stack, fmt.Sprintf("%s:%d", src, line))
		c++
	}
}
