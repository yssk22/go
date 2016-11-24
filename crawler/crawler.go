// Package crawler provides types and functions to run crawler.
package crawler

import (
	"fmt"
	"io"
)

// Fetcher is an interface to get a raw resource for crawled targed.
type Fetcher interface {
	Fetch() (io.Reader, error)
}

// Scraper is an interface to scrape a content
type Scraper interface {
	Scrape(r io.Reader) (interface{}, error)
}

// Run execute the fetcher and pass the content to the scraper.
func Run(fetcher Fetcher, scraper Scraper) (interface{}, error) {
	content, err := fetcher.Fetch()
	if err != nil {
		return nil, fmt.Errorf("fetch error: %v", err)
	}
	return scraper.Scrape(content)
}
