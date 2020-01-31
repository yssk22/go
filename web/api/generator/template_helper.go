package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	"github.com/yssk22/go/x/xerrors"
)

const alwaysOKResponseStatements = `return api.OK(req.Context())
`

const errorResponseStatements = `if err != nil {
	return api.NewErrorResponse(req.Context(), err)
}
`

var templateHelper = template.FuncMap(map[string]interface{}{
	"serialize": func(v interface{}) string {
		// generate a go statement for the json serialized value for v
		buff, err := json.Marshal(v)
		xerrors.MustNil(err)
		return fmt.Sprintf("[]byte(`%s`)", string(buff))
	},

	"genExecMethodAndReturn": func(v *Spec) string {
		var buff bytes.Buffer
		if v.ReceiverName != "" {
			buff.WriteString(fmt.Sprintf("%s.%s", v.ReceiverName, v.FuncName))
		} else {
			buff.WriteString(v.FuncName)
		}
		buff.WriteString("(\n")
		buff.WriteString("req.Context(),\n")
		for _, p := range v.PathParameters {
			buff.WriteString(fmt.Sprintf("req.Params.GetStringOr(%q, \"\"),\n", p))
		}
		if v.StructuredParameter != nil {
			buff.WriteString("&sp,\n")
		}
		buff.WriteString(")")

		switch v.ReturnType {
		case returnTypeNone:
			return fmt.Sprintf("%s\n%s", buff.String(), alwaysOKResponseStatements)
		case returnTypeObject:
			return fmt.Sprintf("return response.NewJSON(req.Context(), %s)", buff.String())
		case returnTypeError:
			return fmt.Sprintf("err := %s\n%s%s", buff.String(), errorResponseStatements, alwaysOKResponseStatements)
		case returnTypeObjectAndError:
			return fmt.Sprintf("obj, err := %s\n%sreturn response.NewJSON(req.Context(), obj)", buff.String(), errorResponseStatements)
		}
		return ""
	},
})
