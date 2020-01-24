package scraper

import (
	"encoding/xml"
	"io"
)

// Rss2Doc is a struct to represent RSS 2.0 document
type Rss2Doc struct {
	XMLName xml.Name     `xml:"rss"`
	Channel *Rss2Channel `xml:"channel"`
}

// Rss2Channel represents `rss>channel` doc.
type Rss2Channel struct {
	Title          string       `xml:"title"`
	Link           string       `xml:"link"`
	Description    string       `xml:"description"`
	Language       string       `xml:"language"`
	Copyright      string       `xml:"copyright"`
	ManagingEditor string       `xml:"managingEditor"`
	WebMaster      string       `xml:"webMaster"`
	Images         []*Rss2Image `xml:"images"`
	LastBuildDate  Time         `xml:"lastBuildDate"`
	Category       string       `xml:"category"`
	Generator      string       `xml:"generator"`
	Items          []*Rss2Item  `xml:"item"`
}

// Rss2Image represents `rss>channel>image` doc.
type Rss2Image struct {
	URL    string `xml:"url"`
	Title  string `xml:"title"`
	Link   string `xml:"link"`
	Width  int64  `xml:"width"`
	Height int64  `xml:"height"`
}

// Rss2Item represents `rss>item` doc.
type Rss2Item struct {
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

// rss2Scraper is a function wrapper for Scraper interface using goquery.
type rss2Scraper func(*Rss2Doc) (interface{}, error)

func (f rss2Scraper) Scrape(r io.Reader) (interface{}, error) {
	decoder := xml.NewDecoder(r)
	var doc Rss2Doc
	err := decoder.Decode(&doc)
	if err != nil {
		return nil, err
	}
	return f(&doc)
}
