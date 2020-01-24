package fetcher

import "io"

// Fetcher is an interface to get a raw resource for crawled targed.
type Fetcher interface {
	Fetch() (io.ReadCloser, error)
}
