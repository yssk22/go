package scraper

import (
	"github.com/PuerkitoBio/goquery"
	"io"
)

// Scraper is an interface to scrape a content
type Scraper interface {
	Scrape(r io.Reader) (interface{}, error)
}

// Func is a function wrapper for Scraper interface
type Func func(r io.Reader) (interface{}, error)

// Scrape implements Scraper#Scrape
func (f Func) Scrape(r io.Reader) (interface{}, error) {
	return f(r)
}

// Html returns a scraper for HTML with a custom logic on top of goquery
func Html(f func(*goquery.Document, *Head) (interface{}, error)) Scraper {
	return htmlScraper(f)
}

// Rss2 returns a scraper for Rss2 with a custom logic on top of Rss2 struct
func Rss2(f func(*Rss2Doc) (interface{}, error)) Scraper {
	return rss2Scraper(f)
}
