package react

import (
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/speedland/go/web/response"
)

// ReactPageVars is a page data used for the `React`` page.
type ReactPageVars struct {
	Status         response.HTTPStatus
	Title          string
	Body           template.HTML
	MetaProperties map[string]string
	ServiceData    map[string]interface{}
	AppData        map[string]interface{}
	Stylesheets    []string
	Javascripts    []string
}

func (rpv *ReactPageVars) Merge(rpv2 *ReactPageVars) {
	rpv.Status = rpv2.Status
	rpv.Body = rpv2.Body
	if rpv2.Title != "" {
		if rpv.Title == "" {
			rpv.Title = fmt.Sprintf("%s", rpv2.Title)
		} else {
			rpv.Title = fmt.Sprintf("%s - %s", rpv2.Title, rpv.Title)
		}
	}
	for key, val := range rpv2.ServiceData {
		rpv.ServiceData[key] = val
	}
	for key, val := range rpv2.MetaProperties {
		rpv.MetaProperties[key] = val
	}
	for key, val := range rpv2.AppData {
		rpv.AppData[key] = val
	}
	for _, val := range rpv2.Stylesheets {
		rpv.Stylesheets = append(rpv.Stylesheets, val)
	}
	for _, val := range rpv2.Javascripts {
		rpv.Javascripts = append(rpv.Javascripts, val)
	}
}

// Default is a default template
var ReactPageTemplate = template.Must(template.New("react").Funcs(reactPageTemplateFuncs).Parse(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{.Title}}</title>
  {{range $key, $val := .MetaProperties -}}
  <meta property="{{$key}}" content="{{$val}}">
  {{- end -}}
  <!--[if lt IE 9]>
  <script src="//oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
  <script src="//oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
  <![endif]-->
  <style>
  body {
      padding: 0;
      margin: 0;
  }
  </style>
  {{range .Stylesheets -}}
  <link rel="stylesheet" type="text/css" href="{{.}}" />
  {{end}}
</head>
<body>
  <div id="fb-root"></div>
  <div id="main"
  	data-app="{{json .AppData}}"
	data-service="{{json .ServiceData}}">
  	<div class="body">{{.Body}}</div>
  </div>
  {{range .Javascripts -}}
  <script type="text/javascript" src="{{.}}"></script>
  {{end}}
</body>
</html>
`))

var reactPageTemplateFuncs = template.FuncMap(map[string]interface{}{
	"json": func(v interface{}) string {
		buff, _ := json.Marshal(v)
		return string(buff)
	},
	"safeHtml": func(v interface{}) template.HTML {
		return template.HTML(fmt.Sprintf("%s", v))
	},
})
