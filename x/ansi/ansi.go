// Package ansi supports ANSI color decorated functions for fmt package
// You can decorate output text with colors by replacing fmt.Print function with ansi.{ColorName}.Print
package ansi

import (
	"fmt"
	"io"
)

// Code implements ANSI code support for fmt.Print, fmt.FPrint, and fmt.Sprint functions.
type Code int

// Available ANSI decorators
const (
	Bold Code = 1 + iota
	Faint
	Italic
	Underline
	BlinkSlow
	BlinkRapid
	ReverseVideo
	Concealed
	CrossedOut
)

// Available ANSI Color
const (
	Black Code = 30 + iota
	Red
	Green
	Yellow
	Blue
	Megenta
	Cyan
	White
)

// Available ANSI Background Color
const (
	BgBlack Code = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgMegenta
	BgCyan
	BgWhite
)

const colorEnd = "\x1b[0m"

// Sprint returns the decorated string with `fmt.Sprint(...interface{})`
func (i Code) Sprint(a ...interface{}) string {
	return fmt.Sprintf("%s%s%s",
		fmt.Sprintf("\x1b[%dm", i),
		fmt.Sprint(a...),
		colorEnd,
	)
}

// Sprintf returns the decorated string with `fmt.Sprint(string, ...interface{})`
func (i Code) Sprintf(s string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s",
		fmt.Sprintf("\x1b[%dm", i),
		fmt.Sprintf(s, a...),
		colorEnd,
	)
}

// Sprintln returns the decorated string with `fmt.Sprint(...interface{})`
func (i Code) Sprintln(s string, a ...interface{}) string {
	return fmt.Sprintf("%s%s%s\n",
		fmt.Sprintf("\x1b[%dm", i),
		fmt.Sprint(a...),
		colorEnd,
	)
}

// Print outputs the decorated string with `fmt.Print(...interface{})`
func (i Code) Print(a ...interface{}) {
	fmt.Printf("%s%s%s",
		fmt.Sprintf("\x1b[%dm", i),
		fmt.Sprint(a...),
		colorEnd,
	)
}

// Printf outputs the decorated string with `fmt.Printf(string, ...interface{})`
func (i Code) Printf(s string, a ...interface{}) {
	fmt.Printf("%s%s%s",
		fmt.Sprintf("\x1b[%dm", i),
		fmt.Sprintf(s, a...),
		colorEnd,
	)
}

// Println outputs the decorated string with `fmt.Println(...interface{})`
func (i Code) Println(a ...interface{}) {
	fmt.Printf("%s%s%s\n",
		fmt.Sprintf("\x1b[%dm", i),
		fmt.Sprint(a...),
		colorEnd,
	)
}

// FPrint outputs the decorated string with `fmt.FPrint(io.Writer, ...interface{})`
func (i Code) FPrint(w io.Writer, a ...interface{}) {
	fmt.Fprintf(w,
		"%s%s%s",
		fmt.Sprintf("\x1b[%dm", i),
		fmt.Sprint(a...),
		colorEnd,
	)
}

// FPrintf outputs the decorated string with `fmt.Printf(string, ...interface{})`
func (i Code) FPrintf(w io.Writer, s string, a ...interface{}) {
	fmt.Fprintf(
		w,
		"%s%s%s",
		fmt.Sprintf("\x1b[%dm", i),
		fmt.Sprintf(s, a...),
		colorEnd,
	)
}

// FPrintln outputs the decorated string with `fmt.Println(...interface{})`
func (i Code) FPrintln(w io.Writer, a ...interface{}) {
	fmt.Fprintf(
		w,
		"%s%s%s\n",
		fmt.Sprintf("\x1b[%dm", i),
		fmt.Sprint(a...),
		colorEnd,
	)
}
