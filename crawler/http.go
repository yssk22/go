package crawler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ErrHTTP is an error when fetch fails.
type ErrHTTP struct {
	Response *http.Response
	Content  []byte
}

func (e *ErrHTTP) Error() string {
	return fmt.Sprintf("HTTPError (status: %s)", e.Response.Status)
}

// HTTPFetcher is an implementation to fetch a content from an url.
type HTTPFetcher struct {
	url        string
	client     *http.Client
	MaxRetries int
	Backoff    func(int, *http.Response, error) time.Duration
}

// DefaultHTTPMaxRetries is a default value of MaxRetries field for a newly created *HTTPFetcher
var DefaultHTTPMaxRetries = 3

// DefaultHTTPBackoff is a default value of Backoff function for a newly created *HTTPFetcher
var DefaultHTTPBackoff = func(i int, resp *http.Response, err error) time.Duration {
	return time.Duration(0)
}

// NewHTTPFetcher returns a new *HTTPFetcher for the given url with http.Client.
// `client`` can be nil, then http.DefaultClient is used.
func NewHTTPFetcher(url string, client *http.Client) *HTTPFetcher {
	if client == nil {
		client = http.DefaultClient
	}
	return &HTTPFetcher{
		url:        url,
		client:     client,
		MaxRetries: DefaultHTTPMaxRetries,
		Backoff:    DefaultHTTPBackoff,
	}
}

// Fetch implements Fetcher#Fetch
func (f *HTTPFetcher) Fetch() (io.Reader, error) {
	var numTries = 0
	for {
		resp, err := f.client.Get(f.url)
		numTries++
		if err != nil {
			if numTries <= f.MaxRetries {
				time.Sleep(f.Backoff(numTries, resp, err))
				continue
			}
			return nil, err
		}
		if resp.StatusCode >= 500 && numTries <= f.MaxRetries {
			resp.Body.Close()
			time.Sleep(f.Backoff(numTries, resp, err))
			continue
		}
		var buff bytes.Buffer
		defer resp.Body.Close()
		io.Copy(&buff, resp.Body)
		if resp.StatusCode != 200 {
			return nil, &ErrHTTP{Response: resp, Content: buff.Bytes()}
		}
		return &buff, nil
	}
}
