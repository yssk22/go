package crawler

import (
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
func (f *FileFetcher) Fetch() (io.ReadCloser, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
