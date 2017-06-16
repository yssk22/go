package react

import (
	"fmt"

	"github.com/speedland/go/lazy"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/middleware/session"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xtime"
	"golang.org/x/net/context"
)

var processStartAt = fmt.Sprintf("%d", xtime.Now().Unix())

var ReactPageDefaults = &ReactPage{
	stylesheets: []interface{}{
	// fmt.Sprintf("/static/page.css?%s", processStartAt),
	},
	javascripts: []interface{}{
		fmt.Sprintf("/static/page.js?%s", processStartAt),
	},
}

// ReactPage is a Page implementation for react applications
type ReactPage struct {
	title          interface{}
	serviceData    map[string]interface{}
	metaProperties map[string]interface{}
	appData        map[string]interface{}
	stylesheets    []interface{}
	javascripts    []interface{}
}

func New() *ReactPage {
	return &ReactPage{
		serviceData:    make(map[string]interface{}),
		metaProperties: make(map[string]interface{}),
		appData:        make(map[string]interface{}),
	}
}

// Title sets the title
func (rp *ReactPage) Title(title interface{}) *ReactPage {
	validateType(title, false, "title")
	rp.title = title
	return rp
}

const (
	serviceDataKeyReactModulePath = "reactModulePath"
	serviceDataKeyCSRFToken       = "csrfToken"
)

// ReactModulePath sets the react module path
func (rp *ReactPage) ReactModulePath(modulePath interface{}) *ReactPage {
	validateType(modulePath, false, "ReactModulePath")
	rp.serviceData[serviceDataKeyReactModulePath] = modulePath
	return rp
}

// MetaProperty sets the meta tag key value pairs
func (rp *ReactPage) MetaProperty(key string, value interface{}) *ReactPage {
	validateType(value, true, fmt.Sprintf("MetaProperty[%q]", key))
	rp.metaProperties[key] = value
	return rp
}

// AppData sets the app data passed to data-{key} attribute on react module.
func (rp *ReactPage) AppData(key string, value interface{}) *ReactPage {
	rp.appData[key] = value
	return rp
}

// Stylesheets add stylesheet on the page.
func (rp *ReactPage) Stylesheets(urls ...interface{}) *ReactPage {
	for _, v := range urls {
		validateType(v, true, "Stylesheets")
	}
	rp.stylesheets = append(rp.stylesheets, urls...)
	return rp
}

// Javascripts add javascript on the page.
func (rp *ReactPage) Javascripts(urls ...interface{}) *ReactPage {
	for _, v := range urls {
		validateType(v, true, "Javascripts")
	}
	rp.javascripts = append(rp.javascripts, urls...)
	return rp
}

// Render implements view.Page#Render
func (rp *ReactPage) Render(req *web.Request) *response.Response {
	ctx := req.Context()
	// s := service.FromContext(ctx)
	data := ReactPageDefaults.genVar(req)
	data.Merge(rp.genVar(req))
	if sess := session.FromContext(ctx); sess != nil {
		data.ServiceData[serviceDataKeyCSRFToken] = sess.CSRFSecret.String()
	}
	return response.NewHTMLWithStatus(
		ReactPageTemplate,
		data,
		data.Status,
	)
}

func (rp *ReactPage) genVar(req *web.Request) *ReactPageVars {
	ctx := req.Context()
	data := &ReactPageVars{
		Status:         response.HTTPStatusOK,
		Title:          genString(ctx, rp.title),
		ServiceData:    make(map[string]interface{}),
		MetaProperties: make(map[string]string),
		AppData:        make(map[string]interface{}),
	}
	for key, val := range rp.serviceData {
		data.ServiceData[key] = genObject(ctx, val)
	}
	for key, val := range rp.metaProperties {
		data.MetaProperties[key] = genString(ctx, val)
	}
	for key, val := range rp.appData {
		data.AppData[key] = genObject(ctx, val)
	}
	for _, val := range rp.stylesheets {
		data.Stylesheets = append(data.Stylesheets, genString(ctx, val))
	}
	for _, val := range rp.javascripts {
		data.Javascripts = append(data.Javascripts, genString(ctx, val))
	}
	return data
}

func validateType(v interface{}, panicIfNil bool, fieldName string) {
	if v == nil && panicIfNil {
		panic(fmt.Sprintf("%s: nil is not allowd", fieldName))
	}
	switch t := v.(type) {
	case string:
		return
	case lazy.Value:
		return
	default:
		panic(fmt.Sprintf("%s: %s is not allowd", fieldName, t))
	}
}

func genString(ctx context.Context, v interface{}) string {
	if v == nil {
		return ""
	}
	switch v.(type) {
	case string:
		return v.(string)
	case lazy.Value:
		str, _ := v.(lazy.Value).Eval(ctx)
		return str.(string)
	default:
		return ""
	}
}

func genObject(ctx context.Context, v interface{}) interface{} {
	if v == nil {
		return nil
	}
	if lv, ok := v.(lazy.Value); ok {
		evaled, _ := lv.Eval(ctx)
		return evaled
	}
	return v
}
