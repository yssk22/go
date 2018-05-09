package reactapp

import (
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xtime"
)

var processStartAt = fmt.Sprintf("%d", xtime.Now().Unix())

// Config represents the built-in data items used in react pages configurations
type PageConfig struct {
	AuthAPIBasePath   string `json:"auth_api_base_path"`
	FacebookAppID     string `json:"facebook_app_id"`
	FacebookPageID    string `json:"facebook_page_id"`
	FacebookPixelID   string `json:"facebook_pixel_id"`
	GoogleAnalyticsID string `json:"google_analytics_id"`
	TwitterID         string `json:"twitter_id"`
	InstagramID       string `json:"instagram_id"`
}

func (c *PageConfig) Merge(obj interface{}) interface{} {
	c1, ok := obj.(*PageConfig)
	if !ok {
		return c
	}
	if c1 == nil {
		return c
	}
	c.AuthAPIBasePath = mergeString(c.AuthAPIBasePath, c1.AuthAPIBasePath)
	c.FacebookAppID = mergeString(c.FacebookAppID, c1.FacebookAppID)
	c.FacebookPageID = mergeString(c.FacebookPageID, c1.FacebookPageID)
	c.FacebookPixelID = mergeString(c.FacebookPixelID, c1.FacebookPixelID)
	c.GoogleAnalyticsID = mergeString(c.GoogleAnalyticsID, c1.GoogleAnalyticsID)
	c.TwitterID = mergeString(c.TwitterID, c1.TwitterID)
	c.InstagramID = mergeString(c.InstagramID, c1.InstagramID)
	return c
}

// Status returns a blank *PageVars with a status code
func Status(s response.HTTPStatus) *PageVars {
	return &PageVars{
		Status: s,
	}
}

// PageVars is a page data generated per an http request from a Page object.
type PageVars struct {
	Status         response.HTTPStatus
	Title          string
	CanonicalURL   string
	Favicon        string
	BasePath       string
	Body           template.HTML
	MetaNames      map[string]string
	MetaProperties map[string]string
	CSRFToken      string
	Config         *PageConfig
	AppName        string
	AppData        map[string]interface{}
	Auth           interface{}
}

func (pv *PageVars) Merge(pv2 *PageVars) {
	if pv2.Title != "" {
		if pv.Title == "" {
			pv.Title = fmt.Sprintf("%s", pv2.Title)
		} else {
			pv.Title = fmt.Sprintf("%s - %s", pv2.Title, pv.Title)
		}
	}
	if pv.Config == nil {
		pv.Config = &PageConfig{}
	}
	if pv2.Status >= 200 {
		pv.Status = pv2.Status
	}
	pv.CanonicalURL = mergeString(pv.CanonicalURL, pv2.CanonicalURL)
	pv.BasePath = mergeString(pv.BasePath, pv2.BasePath)
	pv.Favicon = mergeString(pv.Favicon, pv2.Favicon)
	pv.Body = template.HTML(mergeString(string(pv.Body), string(pv2.Body)))
	pv.Config = mergeObject(pv.Config, pv2.Config).(*PageConfig)
	pv.MetaNames = mergeStringMap(pv.MetaNames, pv2.MetaNames)
	pv.MetaProperties = mergeStringMap(pv.MetaProperties, pv2.MetaProperties)
	pv.AppData = mergeObjectMap(pv.AppData, pv2.AppData)
	pv.Auth = mergeObject(pv.Auth, pv2.Auth)
}

// Default is a default template
var defaultPageTemplate = template.Must(template.New("react").Funcs(reactPageTemplateFuncs).Parse(`<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  {{with .CanonicalURL}}<link rel="canonical" href="{{.}}">{{end}}
  {{with .Favicon}}<link rel="icon" href="{{.}}">{{end}}
  <title>{{.Title}}</title>{{range $key, $val := .MetaNames }}
  <meta name="{{$key}}" content="{{$val}}">
  {{- end -}}{{range $key, $val := .MetaProperties }}
  <meta property="{{$key}}" content="{{$val}}">
  {{- end -}}
  {{with .Config.FacebookAppID }}<meta property="fb:app_id" content="{{.}}">{{ end }}
  {{with .Config.FacebookPageID }}<meta property="fb:pages" content="{{.}}">{{ end }}
  <!--[if lt IE 9]>
  <script src="//oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
  <script src="//oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
  <![endif]-->
</head>
<body>
  <div id="fb-root"></div>
  <div id="root"
	data-auth="{{json .Auth}}"
  	data-app="{{json .AppData}}"
	data-config="{{json .Config}}">
  	<div class="body">{{.Body}}</div>
  </div>
  <script type="text/javascript" src="{{.BasePath}}/{{.AppName}}/static/js/main.js"></script>
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
