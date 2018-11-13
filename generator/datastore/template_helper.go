package datastore

import (
	"text/template"

	"github.com/yssk22/go/x/xstrings"
)

var templateHelper = template.FuncMap(map[string]interface{}{
	"snakecase": func(s string) string {
		return xstrings.ToSnakeCase(s)
	},
})
