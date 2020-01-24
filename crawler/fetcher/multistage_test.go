package fetcher

import (
	"io/ioutil"
	"testing"

	"github.com/yssk22/go/x/xtesting"
	"github.com/yssk22/go/x/xtesting/assert"
)

func TestMultiStage(t *testing.T) {
	r := xtesting.NewRunner(t)
	r.Run("single", func(a *assert.Assert) {
		fetcher := NewMultistageFetcher(NewFileFetcher("./testdata/file.txt"))
		r, err := fetcher.Fetch()
		a.Nil(err)
		defer r.Close()
		buff, err := ioutil.ReadAll(r)
		a.Nil(err)
		a.EqByteString("This is a dummy file to test", buff)
	})

	r.Run("first fail", func(a *assert.Assert) {
		fetcher := NewMultistageFetcher(NewFileFetcher("./testdata/not-exists.txt"), NewFileFetcher("./testdata/file.txt"))
		r, err := fetcher.Fetch()
		a.Nil(err)
		defer r.Close()
		buff, err := ioutil.ReadAll(r)
		a.Nil(err)
		a.EqByteString("This is a dummy file to test", buff)
	})

	r.Run("all fail", func(a *assert.Assert) {
		fetcher := NewMultistageFetcher(NewFileFetcher("./testdata/not-exists.txt"), NewFileFetcher("./testdata/not-exists2.txt"))
		_, err := fetcher.Fetch()
		a.EqStr("open ./testdata/not-exists.txt: no such file or directory (and 1 other errors)", err.Error())
	})
}
