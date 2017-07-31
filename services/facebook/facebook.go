package facebook

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"io/ioutil"

	"bytes"

	"github.com/speedland/go/x/xlog"
	"context"
)

// LoggerKey is a key for this package
const LoggerKey = "speedland.net.services.facebook"

// Client is a API client for facebook graph API.
type Client struct {
	client      *http.Client
	accessToken string
	version     string
}

// NewClient returns a new facebook client
func NewClient(client *http.Client, accessToken string) *Client {
	return &Client{
		client:      client,
		accessToken: accessToken,
		version:     "v2.9",
	}
}

// SetVersion specifies the version used for the client.
func (c *Client) SetVersion(v string) {
	c.version = v
}

// Me is an struct returned by /me endpoint.
type Me struct {
	ID string `json:"id"`
}

// GetMe gets the *Profile of an authorized user.
func (c *Client) GetMe(ctx context.Context) (*Me, error) {
	var m Me
	if err := c.Get(ctx, "/me", url.Values{
		"fields": []string{"id"},
	}, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

const (
	errCouldNotGetResponse string = "could not get response from %s: %s"
	errInvalidStatusCode          = "could not get response from %s: got %d, expected: %d"
	errReadingResponse            = "could not read response from %s: %s"
	errInvalidJSONBody            = "could not parse json from %s: %s"
)

// Get to call GET request on the given path
func (c *Client) Get(ctx context.Context, path string, query url.Values, v interface{}) error {
	url := c.prepareURL(path, query, c.accessToken)
	logger := xlog.WithContext(ctx).WithKey(LoggerKey).WithPrefix(fmt.Sprintf("[GraphAPI:GET %s] ", url))
	resp, err := c.client.Get(url)
	return processResponse(logger, v, url, nil, resp, err)
}

// Post to call POST request on the given path with json body specified by `content` argument.
func (c *Client) Post(ctx context.Context, path string, query url.Values, content interface{}, v interface{}) error {
	url := c.prepareURL(path, query, c.accessToken)
	logger := xlog.WithContext(ctx).WithKey(LoggerKey).WithPrefix(fmt.Sprintf("[GraphAPI:POST %s] ", url))
	var body bytes.Buffer
	err := json.NewEncoder(&body).Encode(content)
	if err != nil {
		return err
	}
	resp, err := c.client.Post(url, "application/json; charset=utf-8", &body)
	return processResponse(logger, v, url, content, resp, err)
}

func (c *Client) prepareURL(path string, query url.Values, accessToken string) string {
	const baseEndpoint = "https://graph.facebook.com/"
	if query == nil {
		query = url.Values{}
	}
	query.Set("access_token", accessToken)
	return fmt.Sprintf("%s%s%s?%s", baseEndpoint, c.version, path, query.Encode())
}

func processResponse(logger *xlog.Logger, v interface{}, url string, content interface{}, resp *http.Response, err error) error {
	if err != nil {
		return logAndError(
			logger,
			"",
			fmt.Errorf(errCouldNotGetResponse, url, err),
		)
	}
	defer resp.Body.Close()
	// expected json bytes are not so big so buffer here.
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return logAndError(
			logger,
			"",
			fmt.Errorf(errReadingResponse, url, err),
		)
	}
	if resp.StatusCode != 200 {
		var message = fmt.Sprintf("response body: %s", string(buff))
		if resp.StatusCode == 400 && content != nil {
			reqBody, _ := json.Marshal(content)
			message = fmt.Sprintf("%s, request body: %s", message, string(reqBody))
		}
		return logAndError(
			logger,
			message,
			fmt.Errorf(errInvalidStatusCode, url, resp.StatusCode, 200),
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
