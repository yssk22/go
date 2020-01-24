package fetcher

import (
	"io"
	"os"
)

// fileFetcher is an implementation to fetch a content from a local file path.
type fileFetcher struct {
	path string
}

// NewFileFetcher returns a new *FileFetcher for the given file path.
func NewFileFetcher(path string) Fetcher {
	return &fileFetcher{
		path: path,
	}
}

// Fetch implements Fetcher#Fetch
func (f *fileFetcher) Fetch() (io.ReadCloser, error) {
	file, err := os.Open(f.path)
	if err != nil {
		return nil, err
	}
	return file, nil
}
