package react

import (
	"html/template"

	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/middleware/session"
	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xerrors"
)

// PageVarsGenerator is an interface to generate a *PageVars
type PageVarsGenerator interface {
	Gen(req *web.Request) (*PageVars, error)
}

// PageVarsGeneratorFunc is a func to ConfigGenerator conversion
type PageVarsGeneratorFunc func(req *web.Request) (*PageVars, error)

// Gen implements ConfigGenerator#Fetch
func (f PageVarsGeneratorFunc) Gen(req *web.Request) (*PageVars, error) {
	return f(req)
}

// Page is a view.Page implementation for react applications.
// The fields in this object is used as a default values of PageVars
type Page struct {
	title          string
	metaNames      map[string]string
	metaProperties map[string]string
	template       *template.Template
	appData        map[string]interface{}
	body           []byte
	stylesheets    []string
	javascripts    []string
	favicon        string
	config         *PageConfig
	generator      PageVarsGenerator
}

var DefaultPage = Must(New())

func Must(p *Page, e error) *Page {
	if e != nil {
		panic(e)
	}
	return p
}

type PageOption func(p *Page) (*Page, error)

// Title returns a PageOption to set the title
func Title(title string) PageOption {
	return func(p *Page) (*Page, error) {
		p.title = title
		return p, nil
	}
}

// MetaProperty returns a PageOption to set the meta props.
func MetaProperty(key string, value string) PageOption {
	return func(p *Page) (*Page, error) {
		p.metaProperties[key] = value
		return p, nil
	}
}

// MetaName returns a PageOption to set the meta name.
func MetaName(key string, value string) PageOption {
	return func(p *Page) (*Page, error) {
		p.metaNames[key] = value
		return p, nil
	}
}

// AppData returns a PageOption to set the AppData field.
func AppData(key string, value interface{}) PageOption {
	return func(p *Page) (*Page, error) {
		p.appData[key] = value
		return p, nil
	}
}

// Body returns a PageOption to add to set body
func Body(b []byte) PageOption {
	return func(p *Page) (*Page, error) {
		p.body = b
		return p, nil
	}
}

// Stylesheets returns a PageOption to add urls into the stylesheet list.
func Stylesheets(urls ...string) PageOption {
	return func(p *Page) (*Page, error) {
		p.stylesheets = append(p.stylesheets, urls...)
		return p, nil
	}
}

// JavasScripts returns a PageOption to add urls into the javascript list.
func JavaScripts(urls ...string) PageOption {
	return func(p *Page) (*Page, error) {
		p.javascripts = append(p.javascripts, urls...)
		return p, nil
	}
}

// Config returns a PageOption to set PageConfig
func Config(c *PageConfig) PageOption {
	return func(p *Page) (*Page, error) {
		p.config = mergeObject(p.config, c).(*PageConfig)
		return p, nil
	}
}

// Faviconreturns a PageOption to set the favicon path
func Favicon(favicon string) PageOption {
	return func(p *Page) (*Page, error) {
		p.favicon = favicon
		return p, nil
	}
}

// GeneratorFunc returns a PageOption to set a *PageVar generator function
func GeneratorFunc(f func(req *web.Request) (*PageVars, error)) PageOption {
	return func(p *Page) (*Page, error) {
		p.generator = PageVarsGeneratorFunc(f)
		return p, nil
	}
}

// Template returns a PageOption to overwrite the default template
func Template(str string) PageOption {
	return func(p *Page) (*Page, error) {
		var err error
		p.template, err = template.New("page-template").Parse(str)
		return p, err
	}
}

func New(options ...PageOption) (*Page, error) {
	p := &Page{
		metaNames:      make(map[string]string),
		metaProperties: make(map[string]string),
		appData:        make(map[string]interface{}),
		config:         &PageConfig{},
		template:       defaultPageTemplate,
	}
	return p.Configure(options...)
}

func (p *Page) Configure(options ...PageOption) (*Page, error) {
	for _, opt := range options {
		var err error
		p, err = opt(p)
		if err != nil {
			return nil, err
		}
	}
	return p, nil
}

// Render implements view.Page#Render
func (p *Page) Render(req *web.Request) *response.Response {
	data, err := DefaultPage.genVar(req)
	if err != nil {
		panic(xerrors.Wrap(err, "genVar error on DefaultPage"))
	}
	d1, err := p.genVar(req)
	if err != nil {
		panic(xerrors.Wrap(err, "genVar error on "))
	}
	data.Merge(d1)
	return response.NewHTMLWithStatus(
		p.template,
		data,
		data.Status,
	)
}

func (p *Page) genVar(req *web.Request) (*PageVars, error) {
	ctx := req.Context()
	// initiate PageVars from scratch
	data := &PageVars{
		Status:         response.HTTPStatusOK,
		MetaNames:      make(map[string]string),
		MetaProperties: make(map[string]string),
		AppData:        make(map[string]interface{}),
		Config:         &PageConfig{},
	}
	if p.title != "" {
		data.Title = p.title
	}
	if p.favicon != "" {
		data.Favicon = p.favicon
	}
	if p.body != nil {
		data.Body = template.HTML(string(p.body))
	}
	data.MetaProperties = mergeStringMap(data.MetaProperties, p.metaProperties)
	data.AppData = mergeObjectMap(data.AppData, p.appData)
	data.Config = mergeObject(data.Config, p.config).(*PageConfig)
	data.Javascripts = mergeStringList(data.Javascripts, p.javascripts)
	data.Stylesheets = mergeStringList(data.Stylesheets, p.stylesheets)
	if p.generator != nil {
		var genData *PageVars
		var err error
		if genData, err = p.generator.Gen(req); err != nil {
			return nil, err
		}
		data.Merge(genData)
	}
	if sess := session.FromContext(ctx); sess != nil {
		data.CSRFToken = sess.CSRFSecret.String()
	}
	return data, nil
}
