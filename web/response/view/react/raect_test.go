package react

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/speedland/go/web"
	"github.com/speedland/go/x/xerrors"
	"github.com/speedland/go/x/xtesting/assert"
)

func genResponse(p *Page) (*goquery.Document, *httptest.ResponseRecorder) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", nil)
	xerrors.MustNil(err)
	req := web.NewRequest(r, nil)
	p.Render(req).Render(req.Context(), w)
	doc, err := goquery.NewDocumentFromReader(w.Body)
	xerrors.MustNil(err)
	return doc, w
}

type html interface {
	Html() (string, error)
}

func dumpHtml(h html) string {
	s, err := h.Html()
	if err != nil {
		return fmt.Sprintf("error: %v", err)
	}
	return s
}

func getConfig(doc *goquery.Document) (*PageConfig, error) {
	s := doc.Find("#main").AttrOr("data-config", "")
	if s == "" {
		return nil, fmt.Errorf("no data-config attribute")
	}
	var c PageConfig
	err := json.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to parse %q", s)
	}
	return &c, nil
}

func getAppData(doc *goquery.Document) (map[string]interface{}, error) {
	s := doc.Find("#main").AttrOr("data-app", "")
	if s == "" {
		return nil, fmt.Errorf("no data-data attribute")
	}
	data := make(map[string]interface{})
	err := json.Unmarshal([]byte(s), &data)
	if err != nil {
		return nil, xerrors.Wrap(err, "failed to parse %q", s)
	}
	return data, nil
}

func Test_Page_Render_Static(t *testing.T) {
	a := assert.New(t)
	p, _ := New(
		Title("タイトル"),
		MetaProperty("og:title", "OGタイトル"),
		AppData("foo", "bar"),
		Config(&PageConfig{
			FacebookAppID:   "fb12345",
			FacebookPixelID: "fbp12345",
		}),
		JavaScripts("/path/to/foo.js"),
		Stylesheets("/path/to/foo.css"),
		GeneratorFunc(func(req *web.Request) (*PageVars, error) {
			return &PageVars{
				Config: &PageConfig{
					FacebookPageID: "mypage",
				},
				AppData: map[string]interface{}{
					"url": req.URL.Path,
				},
			}, nil
		}))

	doc, s := genResponse(p)
	a.EqInt(200, s.Code)
	appData, err := getAppData(doc)
	a.Nil(err)
	cfg, err := getConfig(doc)
	a.Nil(err)

	a.EqStr("タイトル", doc.Find("title").Text())
	a.EqStr("OGタイトル", doc.Find("meta[property='og:title']").AttrOr("content", ""))
	a.EqStr("fb12345", doc.Find("meta[property='fb:app_id']").AttrOr("content", ""), dumpHtml(doc))
	a.EqStr("fbp12345", cfg.FacebookPixelID)
	a.EqStr("bar", appData["foo"].(string))
	a.EqStr("/", appData["url"].(string))
	a.EqInt(1, doc.Find("script[src='/path/to/foo.js']").Length())
	a.EqInt(1, doc.Find("link[href='/path/to/foo.css']").Length(), dumpHtml(doc))
}
