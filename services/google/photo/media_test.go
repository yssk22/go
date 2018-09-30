package photo

import (
	"os"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_parseAlbumFeed(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./fixtures/Test_parseAlbumFeed.xml")
	a.Nil(err, "Open fixture file")
	defer f.Close()
	media, err := parseAlbumFeed(f)
	a.Nil(err, "parseMedia should return no error.")
	a.EqInt(3, len(media), "len(media)")
	a.EqStr("6296695022750349698", media[0].ID, "media[0].ID")
	a.EqStr("6296694673225137745", media[0].AlbumID, "media[0].AlbumID")
	a.EqStr("VIDEO0016.3gp", media[0].Title, "media[0].Title")
	a.EqInt64(6323709, media[0].Size, "media[0].Size")
	a.EqInt(800, media[0].Width, "media[0].Width")
	a.EqInt(480, media[0].Height, "media[0].Height")
	a.EqStr("MOV", media[0].OriginalVideo.Type, "media[0].OriginalVideo.Type")
	a.EqInt(30, media[0].OriginalVideo.Duration, "media[0].OriginalVideo.Duration")

	a.EqInt(2, len(media[0].Contents), "len(media[0].Contents)")
	a.EqStr("image/gif", media[0].Contents[0].Type, "media[0].Contents[0].Type")
	a.EqStr("image", media[0].Contents[0].Medium, "media[0].Contents[0].Medium")
	a.EqInt(512, media[0].Contents[0].Width, "media[0].Contents[0].Width")
	a.EqInt(308, media[0].Contents[0].Height, "media[0].Contents[0].Height")
	a.EqInt(3, len(media[0].Thumbnails), "len(media[0].Thumbnails)")
}

func Test_parsePhotoFeed(t *testing.T) {
	a := assert.New(t)
	f, err := os.Open("./fixtures/Test_parsePhotoFeed.xml")
	a.Nil(err, "Open fixture file")
	defer f.Close()
	media, err := parsePhotoFeed(f)
	a.Nil(err, "parseMedia should return no error.")
	a.EqInt(4, len(media.Contents), "len(media[0].Contents)")
}
