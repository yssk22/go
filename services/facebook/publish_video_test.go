package facebook

import (
	"context"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
)

func Test_CreateSlideshow(t *testing.T) {
	a := assert.New(t)
	c, pageID := newTestClientWithPage(t)
	if c == nil {
		return
	}
	videoID, err := c.CreateSlideshow(context.Background(), pageID, NewSlideshowParams(
		[]string{
			"https://scontent.xx.fbcdn.net/hads-xtf1/t45.1600-4/11410027_6032434826375_425068598_n.png",
			"https://scontent.xx.fbcdn.net/hads-xtp1/t45.1600-4/11410105_6031519058975_1161644941_n.png",
			"http://vignette1.wikia.nocookie.net/parody/images/2/27/Minions_bob_and_his_teddy_bear_2.jpg",
		},
	))
	a.Nil(err)
	a.OK(videoID != "")
}
