package scraper

import (
	"os"
	"testing"
	"time"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestRss2Scraper(t *testing.T) {
	a := assert.New(t)
	s := Rss2(func(rss2 *Rss2Doc) (interface{}, error) {
		return rss2, nil
	})
	f, err := os.Open("./fixtures/rss2.xml")
	a.Nil(err)
	defer f.Close()
	data, err := s.Scrape(f)
	a.Nil(err)
	feed := data.(*Rss2Doc)
	a.NotNil(feed.Channel)
	a.EqStr("Liftoff News", feed.Channel.Title)
	a.EqStr("http://liftoff.msfc.nasa.gov/", feed.Channel.Link)
	a.EqStr("Liftoff to Space Exploration.", feed.Channel.Description)
	a.EqTime(time.Date(2003, 6, 10, 9, 41, 1, 0, time.UTC), time.Time(feed.Channel.LastBuildDate))
	a.EqInt(4, len(feed.Channel.Items))
}
