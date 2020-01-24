package fetcher

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// ErrHTTP is an error when fetch fails.
type ErrHTTP struct {
	Response *http.Response
	Content  []byte
}

func (e *ErrHTTP) Error() string {
	return fmt.Sprintf("HTTPError (status: %s)", e.Response.Status)
}

// httpFetcher is an implementation to fetch a content from an url.
type httpFetcher struct {
	url    string
	client *http.Client
}

// NewHTTPFetcher returns a new *HTTPFetcher for the given url with http.Client.
// `client`` can be nil, then http.DefaultClient is used.
func NewHTTPFetcher(url string, client *http.Client) Fetcher {
	if client == nil {
		client = http.DefaultClient
	}
	return &httpFetcher{
		url:    url,
		client: client,
	}
}

// Fetch implements Fetcher#Fetch
func (f *httpFetcher) Fetch() (io.ReadCloser, error) {
	resp, err := f.client.Get(f.url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		buff, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, &ErrHTTP{Response: resp, Content: buff}
	}
	return resp.Body, nil
}
