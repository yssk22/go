package visualrecognition

import (
	"bytes"
	"encoding/json"
	neturl "net/url"
	"strconv"
	"strings"

	"github.com/yssk22/go/x/xarchive/xzip"

	"context"
)

// ClassifyResponse is a object returned by Datect* functions
type ClassifyResponse struct {
	Images         []*ClassifiedImage `json:"images,omitempty"`
	ImageProcessed int                `json:"image_processed"`
}

type ClassifiedImage struct {
	Classifiers []*ClassifiedResult `json:"classifiers,omitempty"`
	Image       string              `json:"image"`
}

type ClassifiedResult struct {
	Classes      []*ClassScore `json:"classes"`
	ClassifierID string        `json:"classifier_id"`
	Name         string        `json:"name"`
}

type ClassScore struct {
	Class string  `json:"class"`
	Score float64 `json:"score"`
}

type ClassifyParams struct {
	ClassifierIDs []string `json:"classifier_ids"`
	Threshold     float64  `json:"threshold"`
}

// ClassifyURL make a GET request to /v3/classify endpoint with the `url` parameter.
func (c *Client) ClassifyURL(ctx context.Context, url string, params *ClassifyParams) (*ClassifyResponse, error) {
	var resp ClassifyResponse
	var urlParams = neturl.Values(map[string][]string{
		"url": []string{url},
	})
	if params != nil {
		if len(params.ClassifierIDs) > 0 {
			urlParams.Set("classifier_ids", strings.Join(params.ClassifierIDs, ","))
		}
		if params.Threshold >= 0.0 {
			urlParams.Set("threshold", strconv.FormatFloat(params.Threshold, 'f', -1, 64))
		}
	}
	if err := c.get("/v3/classify", urlParams, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// ClassifyImages make a POST request to /v3/classify endpoint with the `images_file` parameter.
func (c *Client) ClassifyImages(ctx context.Context, source *xzip.Archiver, params *ClassifyParams) (*ClassifyResponse, error) {
	var resp ClassifyResponse
	files := make(map[string]*file)
	files["images_file"] = &file{
		name:   "image.zip",
		reader: source,
	}
	if params != nil {
		var buff bytes.Buffer
		encoder := json.NewEncoder(&buff)
		if err := encoder.Encode(params); err != nil {
			return nil, err
		}
		files["parameters"] = &file{
			name:   "parameters.json",
			reader: &buff,
		}
	}
	if err := c.postFiles("/v3/classify", nil, nil, files, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}
