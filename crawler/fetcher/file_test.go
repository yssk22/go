package fetcher

import (
	"io/ioutil"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

func TestFile(t *testing.T) {
	a := assert.New(t)
	fetcher := NewFileFetcher("./testdata/file.txt")
	r, err := fetcher.Fetch()
	a.Nil(err)
	defer r.Close()
	buff, err := ioutil.ReadAll(r)
	a.Nil(err)
	a.EqByteString("This is a dummy file to test", buff)
}
