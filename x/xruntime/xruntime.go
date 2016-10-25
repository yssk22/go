// Package xruntime provides extended utility functions for runtime
package xruntime

import (
	"fmt"
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
		}
		stack[i] = frame
	}
	return stack
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
