package photos

import (
	"encoding/xml"
	"fmt"
	"io"
	"time"
)

// Access is alias for access scope string
type Access string

// Available Access values
const (
	AccessPublic    Access = "public"
	AccessPrivate          = "private"
	AccessProtected        = "protected"
)

// Album represents an album data
type Album struct {
	XMLName            xml.Name   `xml:"http://www.w3.org/2005/Atom entry"`
	XMLNSMedia         string     `xml:"xmlns:media,attr"`
	XMLNSGPhoto        string     `xml:"xmlns:gphoto,attr"`
	ID                 string     `xml:"http://schemas.google.com/photos/2007 id,omitempty"`
	PublishedAt        *time.Time `xml:"published,omitempty"` // use pointer to make 'omitempty' work
	UpdatedAt          *time.Time `xml:"updated,omitempty"`   // use pointer to make 'omitempty' work
	Title              string     `xml:"title"`
	Summary            string     `xml:"summary"`
	Category           *Category  `xml:"category"`
	Timestamp          int64      `xml:"http://schemas.google.com/photos/2007 timestamp,omitempty"`
	Access             Access     `xml:"http://schemas.google.com/photos/2007 access,omitempty"`
	Location           string     `xml:"http://schemas.google.com/photos/2007 location,omitempty"`
	AuthorID           string     `xml:"http://schemas.google.com/photos/2007 user,omitempty"`
	AuthorName         string     `xml:"http://schemas.google.com/photos/2007 nickname,omitempty"`
	NumPhotos          int        `xml:"http://schemas.google.com/photos/2007 numphotos,omitempty"`
	NumPhotosRemaining int        `xml:"http://schemas.google.com/photos/2007 numphotosremaining,omitempty"`
	BytesUsed          int64      `xml:"http://schemas.google.com/photos/2007 bytesUsed,omitempty"`
}

// Category represents the category field in an album
type Category struct {
	Scheme string `xml:"scheme,attr"`
	Term   string `xml:"term,attr"`
}

// NewAlbum returns a new *Album
func NewAlbum() *Album {
	a := &Album{}
	a.XMLNSMedia = "http://search.yahoo.com/mrss/"
	a.XMLNSGPhoto = "http://schemas.google.com/photos/2007"
	a.Category = &Category{
		Scheme: "http://schemas.google.com/g/2005#kind",
		Term:   "http://schemas.google.com/photos/2007#album",
	}
	a.Access = AccessPrivate
	return a
}

func (a *Album) String() string {
	return fmt.Sprintf("<%s:%s>", a.ID, a.Title)
}

type xmlUserFeed struct {
	Entries []*Album `xml:"entry"`
}

func parseUserFeed(r io.Reader) ([]*Album, error) {
	decoder := xml.NewDecoder(r)
	feed := xmlUserFeed{}
	err := decoder.Decode(&feed)
	if err != nil {
		return nil, err
	}
	return feed.Entries, nil
}
