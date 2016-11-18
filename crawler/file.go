package crawler

import (
	"bytes"
	"io"
	"os"
)

// FileFetcher is an implementation to fetch a content from a local file path.
type FileFetcher struct {
	path string
}

// NewFileFetcher returns a new *FileFetcher for the given file path.
func NewFileFetcher(path string) *FileFetcher {
	return &FileFetcher{
		path: path,
	}
}

// Fetch implements Fetcher#Fetch
func (f *FileFetcher) Fetch() (io.Reader, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	defer file.Close()
	io.Copy(&buff, file)
	return &buff, nil
}
