package scraper

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/yssk22/go/x/xstrings"
)

type Head struct {
	Title     string     `json:"title"`
	Canonical string     `json:"canonical"`
	Meta      url.Values `json:"meta"`
	Rel       url.Values `json:"rel"`
}

// htmlScraper is a function wrapper for Scraper interface using goquery.
type htmlScraper func(*goquery.Document, *Head) (interface{}, error)

// Scrape implements Scraper#Scrape
func (f htmlScraper) Scrape(r io.Reader) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	head := &Head{
		Meta: url.Values(make(map[string][]string)),
	}
	head.Title = strings.TrimSpace(doc.Find("head title").Text())
	doc.Find("head meta").Each(func(i int, s *goquery.Selection) {
		key := xstrings.Or(s.AttrOr("name", ""), s.AttrOr("property", ""))
		value := s.AttrOr("content", "")
		if key != "" {
			head.Meta.Add(key, value)
			return
		}
	})
	doc.Find("head link").Each(func(i int, s *goquery.Selection) {
		key := s.AttrOr("rel", "")
		if key != "" {
			head.Meta.Add(key, s.AttrOr("href", ""))
			if key == "canonical" {
				head.Canonical = s.AttrOr("href", "")
			}
			return
		}
	})
	return f(doc, head)
}
