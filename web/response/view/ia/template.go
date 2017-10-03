package ia

import (
	"html/template"
	"time"
)

type PageVars struct {
	CanonicalURL       string
	Style              string
	AutomaticPlacement bool
	CoverImage         string
	CoverCaption       string
	Title              string
	Subtitle           string
	Kicker             string
	AuthorURL          string
	AuthorName         string
	PublishedAt        time.Time
	ModifiedAt         time.Time
	Body               template.HTML
}

var iaMarkupTemplate = template.Must(template.New("ia").Funcs(iaPageTemplateFuncs).Parse(`<!DOCTYPE html>
<html>
	<head>
		<meta charse="utf-8">
		<meta property="op:markup_version" content="v1.0">
		<link rel="canonical" href="{{.CanonicalURL}}">
		<meta property="fb:article_style" content="{{.Style}}">
		<meta property="fb:use_automatic_ad_placement" content="{{.AutomaticPlacement}}">
	</head>
	<body>
	<article>
		<header>
			{{if .CoverImage}}
			<figure>
				<img src="{{.CoverImage}}" />
				{{with .CoverCaption}}<figcaption>{{.}}</figcaption>{{end}}
			</figure>
			{{end}}
			<h1>{{.Title}}</h1>
			{{with .Subtitle}}<h2>{{.}}</h2>{{end}}
			{{with .Kicker}}<h3 class="op-kicker">{{.}}</h3>{{end}}
			<address><a href="{{.AuthorURL}}">{{.AuthorName}}</a></address>
			<time class="op-published" dateTime="{{opDate .PublishedAt}}">{{opDateStr .PublishedAt}}</time>
			<time class="op-modified" dateTime="{{opDate .ModifiedAt}}">{{opDateStr .ModifiedAt}}</time>
		</header>
		{{.Body}}
	</article>
</body>
</html>
`))

var iaPageTemplateFuncs = template.FuncMap(map[string]interface{}{
	"opDate": func(v interface{}) string {
		return (v.(time.Time)).Format(time.RFC3339)
	},
	"opDateStr": func(v interface{}) string {
		return (v.(time.Time)).Format(time.RFC822)
	},
})
