package facebook

import (
	"net/http"
	"testing"

	"os"

	"github.com/speedland/go/x/xtesting/assert"
	"context"
)

func Test_GetMe(t *testing.T) {
	a := assert.New(t)
	token := os.Getenv("TEST_FACEBOOK_ACCESS_TOKEN")
	if token == "" {
		t.Skipf("needs TEST_FACEBOOK_ACCESS_TOKEN envvar for this test.")
		return
	}
	c := NewClient(http.DefaultClient, token)
	me, err := c.GetMe(context.Background())
	a.Nil(err)
	a.OK(me.ID != "")
}
