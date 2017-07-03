package facebook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"io/ioutil"

	"github.com/speedland/go/x/xlog"
	"golang.org/x/net/context"
)

// LoggerKey is a key for this package
const LoggerKey = "speedland.net.services.facebook"

// Client is a API client for facebook graph API.
type Client struct {
	client      *http.Client
	accessToken string
}

// NewClient returns a new facebook client
func NewClient(client *http.Client, accessToken string) *Client {
	return &Client{
		client:      client,
		accessToken: accessToken,
	}
}

// Me is an struct returned by /me endpoint.
type Me struct {
	ID string `json:"id"`
}

// GetMe gets the *Profile of an authorized user.
func (c *Client) GetMe(ctx context.Context) (*Me, error) {
	const endpoint = "https://graph.facebook.com/me"
	var m Me
	if err := c.get(ctx, endpoint, url.Values{
		"fields": []string{"id"},
	}, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

const (
	errCouldNotGetResponse string = "could not get response from %s: %s"
	errInvalidStatusCode          = "could not get response from %s: %d (expected: %d)"
	errReadingResponse            = "could not read response from %s: %s"
	errInvalidJSONBody            = "could not parse json from %s: %s"
)

func (c *Client) get(ctx context.Context, endpoint string, query url.Values, v interface{}) error {
	logger := xlog.WithContext(ctx).WithKey(LoggerKey)
	query.Set("access_token", c.accessToken)
	url := fmt.Sprintf("%s?%s", endpoint, query.Encode())
	resp, err := c.client.Get(url)
	if err != nil {
		return logAndError(
			logger,
			"",
			fmt.Errorf(errCouldNotGetResponse, url, err),
		)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return logAndError(
			logger,
			"",
			fmt.Errorf(errInvalidStatusCode, url, resp.StatusCode, 200),
		)
	}
	// expected json bytes are not so big so buffer here.
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return logAndError(
			logger,
			"",
			fmt.Errorf(errReadingResponse, url, err),
		)
	}
	if err := json.Unmarshal(buff, v); err != nil {
		return logAndError(
			logger,
			fmt.Sprintf("response body: %s", string(buff)),
			fmt.Errorf(errInvalidJSONBody, url, err),
		)
	}
	return nil
}

func logAndError(logger *xlog.Logger, log string, e error) error {
	if log != "" {
		logger.Errorf("%s --- %s", e.Error(), log)
	} else {
		logger.Errorf(e.Error())
	}
	return e
}
