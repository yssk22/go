package xfmt

import "strings"

func PaddingLeft(s string, l int) string {
	sz := len(s)
	if sz >= l {
		return s
	}
	padding := l - sz
	return strings.Repeat(" ", padding) + s
}

func PaddingRight(s string, l int) string {
	sz := len(s)
	if sz >= l {
		return s
	}
	padding := l - sz
	return s + strings.Repeat(" ", padding)
}
