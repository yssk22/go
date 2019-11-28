package html

import (
	"io"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Head struct {
	Title     string     `json:"title"`
	Canonical string     `json:"canonical"`
	Meta      url.Values `json:"meta"`
	Rel       url.Values `json:"rel"`
}

// HTMLScraper is a function wrapper for Scraper interface using goquery.
type HTMLScraper func(*goquery.Document, *Head) (interface{}, error)

// Scrape implements Scraper#Scrape
func (f HTMLScraper) Scrape(r io.Reader) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	head := &Head{
		Meta: url.Values(make(map[string][]string)),
	}
	head.Title = strings.TrimSpace(doc.Find("head title").Text())
	doc.Find("head meta").Each(func(i int, s *goquery.Selection) {
		key := s.AttrOr("property", "")
		if key != "" {
			head.Meta.Add(key, s.AttrOr("content", ""))
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
