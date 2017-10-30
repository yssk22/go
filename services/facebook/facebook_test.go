package facebook

import (
	"net/http"
	"testing"

	"os"

	"context"

	"github.com/speedland/go/x/xtesting/assert"
)

func newTestClient(t *testing.T) *Client {
	token := os.Getenv("TEST_FACEBOOK_ACCESS_TOKEN")
	if token == "" {
		t.Skipf("needs TEST_FACEBOOK_ACCESS_TOKEN envvar for this test.")
		return nil
	}
	return NewClient(http.DefaultClient, token)
}

func newTestClientWithPage(t *testing.T) (*Client, string) {
	token := os.Getenv("TEST_FACEBOOK_ACCESS_TOKEN")
	if token == "" {
		t.Skipf("needs TEST_FACEBOOK_ACCESS_TOKEN envvar for this test.")
		return nil, ""
	}
	pageID := os.Getenv("TEST_FACEBOOK_PAGE")
	if pageID == "" {
		t.Skipf("needs TEST_FACEBOOK_PAGE envvar for this test.")
		return nil, ""
	}
	return NewClient(http.DefaultClient, token), pageID
}

func Test_GetMe(t *testing.T) {
	a := assert.New(t)
	c := newTestClient(t)
	if c == nil {
		return
	}
	me, err := c.GetMe(context.Background())
	a.Nil(err)
	a.OK(me.ID != "")
}

func Test_ScrapeURL(t *testing.T) {
	a := assert.New(t)
	c := newTestClient(t)
	if c == nil {
		return
	}
	urlobj, err := c.ScrapeURL(context.Background(), "http://www.example.com/")
	a.Nil(err)
	a.EqStr("http://www.example.com/", urlobj.URL)
}
