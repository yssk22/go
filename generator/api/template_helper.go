package api

import (
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/yssk22/go/x/xerrors"
)

var templateHelper = template.FuncMap(map[string]interface{}{
	"serialize": func(v interface{}) string {
		// generate a go statement for the json serialized value for v
		buff, err := json.Marshal(v)
		xerrors.MustNil(err)
		return fmt.Sprintf("[]byte(`%s`)", string(buff))
	},
})
