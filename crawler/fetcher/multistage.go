package fetcher

import (
	"io"

	"github.com/yssk22/go/x/xerrors"
)

// multistageFetcher is an implementation to fetch a content from priority orderd fetchers.
type multistageFetcher struct {
	fetchers []Fetcher
}

// NewMultistageFetcher returns a new *FileFetcher for the given file path.
func NewMultistageFetcher(fetchers ...Fetcher) Fetcher {
	return &multistageFetcher{
		fetchers: fetchers,
	}
}

func (f *multistageFetcher) Fetch() (io.ReadCloser, error) {
	errors := xerrors.NewMultiError(len(f.fetchers))
	for i, f := range f.fetchers {
		r, err := f.Fetch()
		if err == nil {
			return r, nil
		}
		errors[i] = err
	}
	return nil, errors
}
