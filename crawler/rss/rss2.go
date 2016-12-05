package rss

import (
	"encoding/xml"
	"io"
)

// Rss2 is a struct to represent RSS 2.0 document
type Rss2 struct {
	XMLName xml.Name `xml:"rss"`
	Channel *Channel `xml:"channel"`
}

// Channel represents `rss>channel` doc.
type Channel struct {
	Title          string   `xml:"title"`
	Link           string   `xml:"link"`
	Description    string   `xml:"description"`
	Language       string   `xml:"language"`
	Copyright      string   `xml:"copyright"`
	ManagingEditor string   `xml:"managingEditor"`
	WebMaster      string   `xml:"webMaster"`
	Images         []*Image `xml:"images"`
	LastBuildDate  Time     `xml:"lastBuildDate"`
	Category       string   `xml:"category"`
	Generator      string   `xml:"generator"`
	Items          []*Item  `xml:"item"`
}

// Image represents `rss>channel>image` doc.
type Image struct {
	URL    string `xml:"url"`
	Title  string `xml:"title"`
	Link   string `xml:"link"`
	Width  int64  `xml:"width"`
	Height int64  `xml:"height"`
}

// Item represents `rss>item` doc.
type Item struct {
	Title          string `xml:"title"`
	Link           string `xml:"link"`
	Description    string `xml:"description"`
	Author         string `xml:"author"`
	Category       string `xml:"category"`
	Comments       string `xml:"comments"`
	GUID           string `xml:"guid"`
	PubDate        Time   `xml:"pubDate"`
	ContentEncoded string `xml:"encoded"`
}

type Rss2Scraper struct{}

func (rss2 *Rss2Scraper) Scrape(r io.Reader) (interface{}, error) {
	decoder := xml.NewDecoder(r)
	var doc Rss2
	err := decoder.Decode(&doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}
