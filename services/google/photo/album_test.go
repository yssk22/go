package photo

import (
	"os"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
	"github.com/speedland/go/x/xtime"
)

func Test_parseUserFeed(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./fixtures/Test_parseUserFeed.xml")
	a.Nil(err, "Open fixture file")
	defer f.Close()
	albums, err := parseUserFeed(f)
	a.Nil(err, "parseAlbums should return no error.")
	a.EqInt(2, len(albums), "len(albums)")
	a.EqStr("1000000487409559", albums[0].ID, "albums[0].ID")
	a.EqStr("自動バックアップ", albums[0].Title, "albums[0].Title")
	a.EqStr("109784491874587409559", albums[0].AuthorID, "albums[0].AuthorID")
	a.EqStr("Yohei Sasaki", albums[0].AuthorName, "albums[0].AuthorName")
	publishedAt, _ := xtime.Parse("2015-11-29T11:03:42.000Z")
	updatedAt, _ := xtime.Parse("2015-11-29T15:32:17.582Z")
	a.EqTime(publishedAt, *albums[0].PublishedAt, "albums[0].PublishedAt")
	a.EqTime(updatedAt, *albums[0].UpdatedAt, "albums[0].UpdatedAt")

	a.EqStr("6173825975323480657", albums[1].ID, "albums[0].ID")
	a.EqStr("2015/07/21", albums[1].Title, "albums[0].Title")
	a.EqStr("109784491874587409559", albums[1].AuthorID, "albums[0].AuthorID")
	a.EqStr("Yohei Sasaki", albums[1].AuthorName, "albums[0].AuthorName")
	publishedAt, _ = xtime.Parse("2015-07-21T05:17:55.000Z")
	updatedAt, _ = xtime.Parse("2015-07-21T05:38:48.016Z")
	a.EqTime(publishedAt, *albums[1].PublishedAt, "albums[0].PublishedAt")
	a.EqTime(updatedAt, *albums[1].UpdatedAt, "albums[0].UpdatedAt")
}
