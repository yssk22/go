// Package xstrings provides extended utility functions for strings
package xstrings

import (
	"strings"
	"unicode"

	"golang.org/x/text/width"
)

// Or returns the first element of string array which is not empty
func Or(arr ...string) string {
	for _, s := range arr {
		if s != "" {
			return s
		}
	}
	return ""
}

// SplitAndTrim is like strings.Split but spaces in each of item are trimmed
func SplitAndTrim(s string, sep string) []string {
	var list []string
	splitted := strings.Split(s, sep)
	for _, v := range splitted {
		s := strings.TrimSpace(v)
		if s != "" {
			list = append(list, s)
		}
	}
	return list
}

// SplitAndTrimAsMap is like SplitAndTrim but returns an map[string]bool
// where the map key is the string flagment included in
func SplitAndTrimAsMap(s string, sep string) map[string]struct{} {
	var m = make(map[string]struct{})
	splitted := strings.Split(s, sep)
	for _, v := range splitted {
		s := strings.TrimSpace(v)
		if s != "" {
			m[s] = struct{}{}
		}
	}
	return m
}

// ToSnakeCase converts the string to the one by snake case.
func ToSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}
	var runes = []rune(s)
	var str = []rune{unicode.ToLower(runes[0])}
	if len(runes) == 1 {
		return string(str)
	}
	for i := 1; i < len(runes)-1; i++ {
		previous := runes[i-1]
		current := runes[i]
		next := runes[i+1]
		if unicode.IsUpper(current) {
			if !unicode.IsUpper(next) || !unicode.IsUpper(previous) {
				str = append(str, '_')
			}
			str = append(str, unicode.ToLower(current))
		} else {
			str = append(str, runes[i])
		}
	}
	str = append(str, unicode.ToLower(runes[len(runes)-1]))
	return string(str)
}

// Fold transforms all runes to their canonical width. unicode string
func Fold(s string) string {
	return width.Fold.String(s)
}

// Narrow transforms all runes to the narrow unicode string
func Narrow(s string) string {
	return width.Narrow.String(s)
}

// Widen transforms all runes to the widen unicode string
func Widen(s string) string {
	return width.Widen.String(s)
}
