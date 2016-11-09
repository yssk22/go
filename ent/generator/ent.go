// Package generator provides types and functions for ent generator.
package generator

import "regexp"

var tagRegexp = regexp.MustCompile(`([a-z0-9A-Z]+):"([^"]+)"`)

const (
	tagNameDefault         = "default"
	tagNameEnt             = "ent"
	tagValueID             = "id"
	tagValueResetIfMissing = "resetifmissing"
	tagValueForm           = "form"
	tagValueTimestamp      = "timestamp"
)
