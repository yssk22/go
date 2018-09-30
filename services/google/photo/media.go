package photo

import (
	"encoding/xml"
	"io"
	"time"
)

// UploadMediaInfo is a struct used for uploading a media
type UploadMediaInfo struct {
	XMLName     xml.Name  `xml:"http://www.w3.org/2005/Atom entry"`
	XMLNSMedia  string    `xml:"xmlns:media,attr"`
	XMLNSGPhoto string    `xml:"xmlns:gphoto,attr"`
	Title       string    `xml:"title"`
	Summary     string    `xml:"summary"`
	Category    *Category `xml:"category"`
	Client      string    `xml:"gphoto:client"`
}

// Media represents media data
type Media struct {
	ID            string         `xml:"http://schemas.google.com/photos/2007 id,omitempty"`
	AlbumID       string         `xml:"http://schemas.google.com/photos/2007 albumid,omitempty"`
	Title         string         `xml:"title"`
	PublishedAt   *time.Time     `xml:"published,omitempty"`
	UpdatedAt     *time.Time     `xml:"updated,omitempty"`
	Summary       string         `xml:"summary"`
	Category      *Category      `xml:"category"`
	Size          int64          `xml:"http://schemas.google.com/photos/2007 size,omitempty"`
	Width         int            `xml:"http://schemas.google.com/photos/2007 width,omitempty"`
	Height        int            `xml:"http://schemas.google.com/photos/2007 height,omitempty"`
	VideoStatus   string         `xml:"http://schemas.google.com/photos/2007 videostatus,omitempty"`
	OriginalVideo *OriginalVideo `xml:"http://schemas.google.com/photos/2007 originalvideo,omitempty"`
	Contents      []*Content     `xml:"http://search.yahoo.com/mrss/ group>content,omitempty"`
	Thumbnails    []*Thumbnail   `xml:"http://search.yahoo.com/mrss/ group>thumbnail,omitempty"`
}

// OriginalVideo represents original_video field data.
type OriginalVideo struct {
	Width        int     `xml:"width,attr"`
	Height       int     `xml:"height,attr"`
	Duration     int     `xml:"duration,attr"`
	Type         string  `xml:"type,attr"`
	Channels     int     `xml:"channels,attr"`
	SamplingRate float32 `xml:"samplingrate,attr"`
	VideoCodec   string  `xml:"videoCodec,attr"`
	AudioCodec   string  `xml:"audioCodec,attr"`
	Fps          float32 `xml:"fps,attr"`
}

// Content represents content field data
type Content struct {
	URL    string `xml:"url,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
	Type   string `xml:"type,attr"`
	Medium string `xml:"medium,attr"`
}

// Thumbnail represents thumbnail field data
type Thumbnail struct {
	URL    string `xml:"url,attr"`
	Width  int    `xml:"width,attr"`
	Height int    `xml:"height,attr"`
}

// NewUploadMediaInfo returns a new *UploadMediaInfo
func NewUploadMediaInfo(title string) *UploadMediaInfo {
	m := &UploadMediaInfo{}
	m.Title = title
	m.Client = "gphoto client (github.com/yssk22)"
	m.Category = &Category{
		Scheme: "http://schemas.google.com/g/2005#kind",
		Term:   "http://schemas.google.com/photos/2007#photo",
	}
	m.XMLNSMedia = "http://search.yahoo.com/mrss/"
	m.XMLNSGPhoto = "http://schemas.google.com/photos/2007"
	return m
}

type xmlAlbumFeed struct {
	Entries []*Media `xml:"entry"`
}

func parseAlbumFeed(r io.Reader) ([]*Media, error) {
	decoder := xml.NewDecoder(r)
	feed := xmlAlbumFeed{}
	err := decoder.Decode(&feed)
	if err != nil {
		return nil, err
	}
	return feed.Entries, nil
}

func parsePhotoFeed(r io.Reader) (*Media, error) {
	decoder := xml.NewDecoder(r)
	media := Media{}
	err := decoder.Decode(&media)
	if err != nil {
		return nil, err
	}
	return &media, nil
}
