package html

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/speedland/go/x/xtesting/assert"
)

func TestHTMLScraper(t *testing.T) {
	a := assert.New(t)
	s := HTMLScraper(func(doc *goquery.Document, h *Head) (interface{}, error) {
		return h, nil
	})
	f, err := os.Open("./fixtures/sample.html")
	a.Nil(err)
	defer f.Close()
	data, err := s.Scrape(f)
	a.Nil(err)
	h := data.(*Head)
	a.EqStr("The Rock (1996)", h.Title)
	a.EqStr("The Rock", h.Meta.Get("og:title"))
}
