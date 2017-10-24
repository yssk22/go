package facebook

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/speedland/go/x/xtime"

	"github.com/speedland/go/x/xtesting/assert"
)

func Test_PublishPostMessage(t *testing.T) {
	a := assert.New(t)
	c := newTestClient(t)
	if c == nil {
		return
	}
	pageID := os.Getenv("TEST_FACEBOOK_PAGE")
	if pageID == "" {
		t.Skipf("needs TEST_FACEBOOK_PAGE envvar for this test.")
		return
	}
	id, err := c.PublishPost(context.Background(), pageID, &PagePostParams{
		Message: fmt.Sprintf("This is a test at %s", xtime.Now()),
	})
	a.Nil(err)
	a.OK(id != "")
	t.Logf("PostURL: https://www.facebook.com/permalink.php?story_fbid=%s&id=%s", id, pageID)
}
