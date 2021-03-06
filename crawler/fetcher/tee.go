package fetcher

import (
	"io"
)

// teeFetcher is an implementation to `tee` a content from a fetcher to another destination
type teeFetcher struct {
	src Fetcher
	dst io.WriteCloser
}

// NewTeeFetcher returns a new *HTTPFetcher for the given url with http.Client.
// `client`` can be nil, then http.DefaultClient is used.
func NewTeeFetcher(src Fetcher, dst io.WriteCloser) Fetcher {
	return &teeFetcher{
		src: src,
		dst: dst,
	}
}

type teeReadCloser struct {
	src io.ReadCloser
	dst io.WriteCloser
	tee io.Reader
}

func (tee *teeReadCloser) Read(p []byte) (int, error) {
	return tee.tee.Read(p)
}

func (tee *teeReadCloser) Close() error {
	tee.src.Close()
	return tee.dst.Close()
}

// Fetch implements Fetcher#Fetch
func (f *teeFetcher) Fetch() (io.ReadCloser, error) {
	r, err := f.src.Fetch()
	if err != nil {
		return nil, err
	}
	return &teeReadCloser{
		src: r,
		dst: f.dst,
		tee: io.TeeReader(r, f.dst),
	}, nil
}
