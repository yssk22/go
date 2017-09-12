package react

import (
	"html/template"

	"github.com/speedland/go/web"
	"github.com/speedland/go/web/middleware/session"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xerrors"
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
	metaProperties map[string]string
	appData        map[string]interface{}
	body           []byte
	stylesheets    []string
	javascripts    []string
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

func Config(c *PageConfig) PageOption {
	return func(p *Page) (*Page, error) {
		p.config = mergeObject(p.config, c).(*PageConfig)
		return p, nil
	}
}

func New(options ...PageOption) (*Page, error) {
	p := &Page{
		metaProperties: make(map[string]string),
		appData:        make(map[string]interface{}),
		config:         &PageConfig{},
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

	// if fetcher := rp.ServiceDataFetcher; fetcher != nil {
	// 	sData, err := fetcher.Fetch(ctx)

	// }
	// if fetcher != nil {
	// 	if err == nil {
	// 		if sData.FacebookAppID != "" {
	// 			data.ServiceData[serviceDataKeyFacebookAppID] = fb.ClientID
	// 			data.MetaProperties["fb:app_id"] = fb.ClientID
	// 		}
	// 	}
	// }
	// if s := service.FromContext(ctx); s != nil {
	// 	if fb := s.Config.GetFacebookConfig(ctx); fb != nil {
	// 		data.ServiceData[serviceDataKeyFacebookAppID] = fb.ClientID
	// 		data.MetaProperties["fb:app_id"] = fb.ClientID
	// 	}
	// 	if apiConfig := s.APIConfig; apiConfig != nil {
	// 		data.ServiceData[serviceDataKeyAuthAPIBasePath] = apiConfig.AuthAPIBasePath
	// 	}
	// }
	return response.NewHTMLWithStatus(
		pageTemplate,
		data,
		data.Status,
	)
}

func (p *Page) genVar(req *web.Request) (*PageVars, error) {
	ctx := req.Context()
	// initiate PageVars from scratch
	data := &PageVars{
		Status:         response.HTTPStatusOK,
		MetaProperties: make(map[string]string),
		AppData:        make(map[string]interface{}),
		Config:         &PageConfig{},
	}
	if p.title != "" {
		data.Title = p.title
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
