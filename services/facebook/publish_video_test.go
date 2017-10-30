package facebook

import (
	"context"
	"fmt"
	"testing"

	"github.com/speedland/go/x/xtesting/assert"
	"github.com/speedland/go/x/xtime"
)

func Test_CreateSlideshow(t *testing.T) {
	a := assert.New(t)
	c, pageID := newTestClientWithPage(t)
	if c == nil {
		return
	}
	params := NewSlideshowParams(
		[]string{
			"https://scontent.xx.fbcdn.net/hads-xtf1/t45.1600-4/11410027_6032434826375_425068598_n.png",
			"https://scontent.xx.fbcdn.net/hads-xtp1/t45.1600-4/11410105_6031519058975_1161644941_n.png",
			"http://vignette1.wikia.nocookie.net/parody/images/2/27/Minions_bob_and_his_teddy_bear_2.jpg",
		},
	)
	params.Title = "Foo"
	params.Description = fmt.Sprintf("This is a test at %s", xtime.Now())
	videoID, err := c.CreateSlideshow(context.Background(), pageID, params)
	a.Nil(err)
	a.OK(videoID != "")
}
