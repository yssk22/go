package xruntime

import "fmt"

func ExampleCaptureStack() {
	frames := CaptureStack(50)
	fmt.Println(frames[0])
	// Output:
	// github.com/yssk22/go/x/xruntime.ExampleCaptureStack (at github.com/yssk22/go/x/xruntime/xruntime_test.go#6)
}

func ExampleCaptureFrame() {
	f := CaptureFrame()
	fmt.Printf("PackageName: %s\n", f.PackageName)
	fmt.Printf("FunctionName: %s\n", f.FunctionName)
	fmt.Printf("ShortFilePath: %s\n", f.ShortFilePath)
	fmt.Printf("LineNumber: %d\n", f.LineNumber)
	// Output:
	// PackageName: github.com/yssk22/go/x/xruntime
	// FunctionName: ExampleCaptureFrame
	// ShortFilePath: github.com/yssk22/go/x/xruntime/xruntime_test.go
	// LineNumber: 13
}

type T struct{}

func (*T) F() *Frame {
	return CaptureFrame()
}

func ExampleCaptureFrame_forStructFunc() {
	f := (&T{}).F()
	fmt.Printf("PackageName: %s\n", f.PackageName)
	fmt.Printf("FunctionName: %s\n", f.FunctionName)
	fmt.Printf("ShortFilePath: %s\n", f.ShortFilePath)
	fmt.Printf("LineNumber: %d\n", f.LineNumber)
	// Output:
	// PackageName: github.com/yssk22/go/x/xruntime
	// FunctionName: (*T).F
	// ShortFilePath: github.com/yssk22/go/x/xruntime/xruntime_test.go
	// LineNumber: 28
}
