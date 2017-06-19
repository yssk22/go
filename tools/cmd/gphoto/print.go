// +build linux darwin freebsd netbsd openbsd dragonfly
// +build !appengine

package main

import (
	"fmt"
	"os"
)

func printTable(t [][]string) {
	maxChars := make([]int, len(t[0]))
	for _, row := range t {
		for j, col := range row {
			charLen := len(col) + 4
			if maxChars[j] < charLen {
				maxChars[j] = charLen
			}
		}
	}
	for _, row := range t {
		for j, col := range row {
			fmt.Printf("%s", col)
			for k := len(col); k < maxChars[j]; k++ {
				fmt.Printf(" ")
			}
		}
		fmt.Printf("\n")
	}
}

func printError(err error) {
	if err == nil {
		return
	}
	os.Stderr.WriteString(fmt.Sprintf("[Error] %s\n", err))
}
