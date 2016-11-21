// Package visualrecognition provides types and functions for
// IBM Watson Visual Recognition API
package visualrecognition

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	apiVersion = "2016-05-20"
	endpoint   = "https://gateway-a.watsonplatform.net/visual-recognition/api"
)

// Client is an API client to access IBM Watson Visual Recognition API
type Client struct {
	APIKey   string
	Version  string
	Endpoint string
	client   *http.Client
}

func NewClient(apiKey string, client *http.Client) *Client {
	return &Client{
		APIKey:   apiKey,
		client:   client,
		Version:  apiVersion,
		Endpoint: endpoint,
	}
}

type ErrorResponse struct {
	Err  string `json:"error"`
	Code int    `json:"code"`
}

func (er *ErrorResponse) Error() string {
	return er.Err
}

func (c *Client) get(path string, params url.Values, dst interface{}) error {
	if params == nil {
		params = url.Values(make(map[string][]string))
	}
	params.Set("api_key", c.APIKey)
	params.Set("version", c.Version)
	resp, err := c.client.Get(fmt.Sprintf("%s%s?%s", c.Endpoint, path, params.Encode()))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	decorder := json.NewDecoder(resp.Body)
	if resp.StatusCode != 200 {
		var errResp ErrorResponse
		if err := decorder.Decode(&errResp); err != nil {
			return err
		}
		return &errResp
	}
	if err := decorder.Decode(dst); err != nil {
		return err
	}
	return nil
}
