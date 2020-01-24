package scraper

import (
	"encoding/xml"
	"time"
)

var timeFormats = []string{
	time.RFC822,
	time.RFC822Z,
	time.RFC1123,
	time.RFC1123Z,
}

// Time is an alias type for time.Time to unmarshal XML doc.
type Time time.Time

// UnmarshalXML to implement xml unmarshalization
func (t *Time) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	var v string
	if err := d.DecodeElement(&v, &start); err != nil {
		return err
	}
	var tt time.Time
	var err error
	for _, format := range timeFormats {
		tt, err = time.Parse(format, v)
		if err == nil {
			*t = Time(tt)
			return nil
		}
	}
	return err
}
