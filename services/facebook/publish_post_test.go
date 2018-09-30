package facebook

import (
	"context"
	"fmt"
	"testing"

	"github.com/yssk22/go/x/xtime"

	"github.com/yssk22/go/x/xtesting/assert"
)

func Test_PublishPostMessage(t *testing.T) {
	a := assert.New(t)
	c, pageID := newTestClientWithPage(t)
	if c == nil {
		return
	}
	res, err := c.PublishPost(context.Background(), pageID, &PagePostParams{
		Message: fmt.Sprintf("This is a test at %s", xtime.Now()),
	})
	a.Nil(err)
	a.EqStr(pageID, res.PageID)
	a.OK(res.PageID != "")
}

func Test_PublishPostLink(t *testing.T) {
	a := assert.New(t)
	c, pageID := newTestClientWithPage(t)
	if c == nil {
		return
	}
	res, err := c.PublishPost(context.Background(), pageID, &PagePostParams{
		Message: fmt.Sprintf("This is a test at %s", xtime.Now()),
		Link:    "http://www.example.com/",
	})
	a.Nil(err)
	a.EqStr(pageID, res.PageID)
	a.OK(res.PageID != "")
}

func Test_PublishPostLinkCarousel(t *testing.T) {
	a := assert.New(t)
	c, pageID := newTestClientWithPage(t)
	if c == nil {
		return
	}
	res, err := c.PublishPost(context.Background(), pageID, &PagePostParams{
		Message: fmt.Sprintf("This is a test at %s", xtime.Now()),
		Link:    "http://www.example.com/",
		ChildAttachments: []*Attachment{
			&Attachment{
				Name:        "Same Link",
				Description: "same link description",
				Link:        "http://www.example.com/",
				Picture:     "https://scontent.xx.fbcdn.net/hads-xtf1/t45.1600-4/11410027_6032434826375_425068598_n.png",
			},
			&Attachment{
				Name:        "Google",
				Description: "linked to google",
				Link:        "http://www.google.com/",
				Picture:     "https://scontent.xx.fbcdn.net/hads-xtp1/t45.1600-4/11410105_6031519058975_1161644941_n.png",
			},
		},
	})
	a.Nil(err)
	a.EqStr(pageID, res.PageID)
	a.OK(res.PageID != "")
}
