// Package xstrings provides extended utility functions for strings
package xstrings

import "unicode"

// ToSnakeCase converts the string to the one by snake case.
func ToSnakeCase(s string) string {
	if len(s) == 0 {
		return s
	}
	var runes = []rune(s)
	var str = []rune{unicode.ToLower(runes[0])}
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

// var specialUpperCaseWords = []string{
// 	"API", "URL", "ID",
// }

// func getSpecialUpperCaseWordLen(runes []rune, i int) int {
// 	for _, word := range specialUpperCaseWords {
// 		l := len(word)
// 		if i+l < len(runes) {
// 			if string(runes[i:i+l]) == word {
// 				return l
// 			}
// 		}
// 	}
// 	return -1
// }
