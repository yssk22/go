// Package visualrecognition provides types and functions for
// IBM Watson Visual Recognition API
package visualrecognition

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
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

func (c *Client) get(path string, query url.Values, dst interface{}) error {
	if query == nil {
		query = url.Values(make(map[string][]string))
	}
	query.Set("api_key", c.APIKey)
	query.Set("version", c.Version)
	resp, err := c.client.Get(fmt.Sprintf("%s%s?%s", c.Endpoint, path, query.Encode()))
	return handleResponse(resp, err, dst)
}

func (c *Client) postFiles(path string, query url.Values, params url.Values, files map[string]*file, dst interface{}) error {
	var buff bytes.Buffer
	writer := multipart.NewWriter(&buff)
	if params != nil {
		for k, vv := range params {
			for _, v := range vv {
				if err := writer.WriteField(k, v); err != nil {
					return err
				}
			}
		}
	}
	for k, f := range files {
		w, err := writer.CreateFormFile(k, f.name)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, f.reader); err != nil {
			return err
		}
	}
	writer.Close()
	url := c.buildURL(path, query)
	resp, err := c.client.Post(url, writer.FormDataContentType(), &buff)
	return handleResponse(resp, err, dst)
	// all, _ := ioutil.ReadAll(&buff)
	// fmt.Println(string(all))
	// return fmt.Errorf("ERR")
}

func (c *Client) delete(path string, query url.Values) error {
	if query == nil {
		query = url.Values(make(map[string][]string))
	}
	query.Set("api_key", c.APIKey)
	query.Set("version", c.Version)
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s?%s", c.Endpoint, path, query.Encode()), nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		return nil
	}
	return handleErrorResponse(resp)
}

func (c *Client) buildURL(path string, query url.Values) string {
	if query == nil {
		query = url.Values(make(map[string][]string))
	}
	query.Set("api_key", c.APIKey)
	query.Set("version", c.Version)
	return fmt.Sprintf("%s%s?%s", c.Endpoint, path, query.Encode())
}

func handleResponse(resp *http.Response, err error, dst interface{}) error {
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		if resp.StatusCode == 413 {
			// Watson seems return HTML file in this case
			return &ErrorResponse{
				Code: resp.StatusCode,
				Err:  "entity too large",
			}
		}
		return handleErrorResponse(resp)
	}
	decorder := json.NewDecoder(resp.Body)
	if err := decorder.Decode(dst); err != nil {
		return err
	}
	return nil
}

func handleErrorResponse(resp *http.Response) error {
	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var errResp ErrorResponse
	json.Unmarshal(buff, &errResp)
	// unespected error response case.
	if errResp.Code == 0 {
		errResp.Code = resp.StatusCode
		errResp.Err = string(buff)
	}
	return &errResp
}

type file struct {
	name   string
	reader io.Reader
}
