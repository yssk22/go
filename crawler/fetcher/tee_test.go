package fetcher

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/yssk22/go/x/xtesting/assert"
)

type buffCloser struct {
	bytes.Buffer
}

func (bc *buffCloser) Close() error {
	return nil
}

func TestTee(t *testing.T) {
	a := assert.New(t)
	var buff buffCloser
	fetcher := NewTeeFetcher(NewFileFetcher("./testdata/file.txt"), &buff)
	r, err := fetcher.Fetch()
	a.Nil(err)
	defer r.Close()
	buff2, err := ioutil.ReadAll(r)
	a.Nil(err)
	a.EqByteString("This is a dummy file to test", buff.Bytes())
	a.EqByteString("This is a dummy file to test", buff2)
}
