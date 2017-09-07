// Package crawler provides types and functions to run crawler.
package crawler

import (
	"fmt"
	"io"

	"github.com/PuerkitoBio/goquery"
)

// Fetcher is an interface to get a raw resource for crawled targed.
type Fetcher interface {
	Fetch() (io.Reader, error)
}

// Scraper is an interface to scrape a content
type Scraper interface {
	Scrape(r io.Reader) (interface{}, error)
}

// ScraperFunc is a function wrapper for Scraper interface
type ScraperFunc func(r io.Reader) (interface{}, error)

// Scrape implements Scraper#Scrape
func (f ScraperFunc) Scrape(r io.Reader) (interface{}, error) {
	return f(r)
}

// HTMLScraper is a function wrapper for Scraper interface using goquery.
// [Deprecation] use github.com/speedland/go/crawler/html.HTMLScraper instead
type HTMLScraper func(doc *goquery.Document) (interface{}, error)

// Scrape implements Scraper#Scrape
// [Deprecation] use github.com/speedland/go/crawler/html.HTMLScraper#Scrape instead
func (f HTMLScraper) Scrape(r io.Reader) (interface{}, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	return f(doc)
}

// Run execute the fetcher and pass the content to the scraper.
func Run(fetcher Fetcher, scraper Scraper) (interface{}, error) {
	content, err := fetcher.Fetch()
	if err != nil {
		return nil, fmt.Errorf("fetch error: %v", err)
	}
	return scraper.Scrape(content)
}
