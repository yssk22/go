package generator

import (
	"regexp"

	"github.com/yssk22/go/keyvalue"
)

var tagRegexp = regexp.MustCompile(`([a-z0-9A-Z]+):"([^"]+)"`)

// ParseTag returns a keyvalue.Getter for the defined tag like `json:"foo"`.
func ParseTag(tag string) keyvalue.Getter {
	m := keyvalue.NewStringKeyMap()
	if tags := tagRegexp.FindAllStringSubmatch(tag, -1); tags != nil {
		for _, tag := range tags {
			tagName := tag[1]
			tagValue := tag[2]
			m.Set(tagName, tagValue)
		}
	}
	return m
}
