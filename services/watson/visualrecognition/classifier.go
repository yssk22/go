package visualrecognition

import (
	"fmt"
	"net/url"
	"time"

	"github.com/speedland/go/x/xarchive/xzip"

	"context"
)

// Status is a type alias to represent classifier status string.
type Status string

// Available (known) values of ClassifierStatus.
const (
	StatusAvailable   Status = "available"
	StatusReady       Status = "ready"
	StatusUnavailable Status = "unavailable"
	StatusFailed      Status = "failed"
)

// ClassifierListResponse is a object returned by ListClassifiers function.
type ClassifierListResponse struct {
	Classifiers []*Classifier `json:"classifiers"`
}

// Classifier represents a classifier on Watson
type Classifier struct {
	ClassifierID string    `json:"classifier_id"`
	Name         string    `json:"name"`
	Owner        string    `json:"owner"`
	Status       Status    `json:"status"`
	Created      time.Time `json:"created"`
	Explanation  string    `json:"explanation"`
	Classes      []*Class  `json:"classes"`
}

// Class is a class of classifiers on Watson.
type Class struct {
	Class string `json:"class"`
}

// CreateClassifier makes a POST request to '/v3/classifiers' endpoint with the *_positive_examples files.
// the source map should be a map from class name to *xzip.Archiver so that `{class}_positive_examples` are registered.
// A class name "negative" is registered as "negative_examples"
func (c *Client) CreateClassifier(ctx context.Context, name string, positives map[string]*xzip.Archiver, negative *xzip.Archiver) (*Classifier, error) {
	files := make(map[string]*file)
	for cls, source := range positives {
		files[fmt.Sprintf("%s_positive_examples", cls)] = &file{
			name:   fmt.Sprintf("%s.zip", cls),
			reader: source,
		}
	}
	if negative != nil {
		files["negative_examples"] = &file{
			name:   "negative.zip",
			reader: negative,
		}
	}
	var classifier Classifier
	if err := c.postFiles("/v3/classifiers", nil, map[string][]string{
		"name": []string{name},
	}, files, &classifier); err != nil {
		return nil, err
	}
	return &classifier, nil
}

// UpdateClassifier makes a POST request to '/v3/classifiers/{classifier_id}' endpoint with the *_positive_examples files.
// the source map should be a map from class name to *xzip.Archiver so that `{class}_positive_examples` are registered.
// A class name "negative" is registered as "negative_examples"
func (c *Client) UpdateClassifier(ctx context.Context, id string, positives map[string]*xzip.Archiver, negative *xzip.Archiver) (*Classifier, error) {
	files := make(map[string]*file)
	for cls, source := range positives {
		files[fmt.Sprintf("%s_positive_examples", cls)] = &file{
			name:   fmt.Sprintf("%s.zip", cls),
			reader: source,
		}
	}
	if negative != nil {
		files["negative_examples"] = &file{
			name:   "negative.zip",
			reader: negative,
		}
	}
	var classifier Classifier
	if err := c.postFiles(fmt.Sprintf("/v3/classifiers/%s", id), nil, nil, files, &classifier); err != nil {
		return nil, err
	}
	return &classifier, nil
}

// ListClassifiers makes a GET request to '/v3/classifiers' endpoint
func (c *Client) ListClassifiers(ctx context.Context, verbose bool) (*ClassifierListResponse, error) {
	var clr ClassifierListResponse
	var params = url.Values(map[string][]string{})
	if verbose {
		params.Set("verbose", "true")
	}
	if err := c.get("/v3/classifiers", params, &clr); err != nil {
		return nil, err
	}
	return &clr, nil
}

// GetClassifier makes a GET request to '/v3/classifiers/{id}`` endpoint.
func (c *Client) GetClassifier(ctx context.Context, id string) (*Classifier, error) {
	var cls Classifier
	var params = url.Values(map[string][]string{})
	if err := c.get(fmt.Sprintf("/v3/classifiers/%s", id), params, &cls); err != nil {
		return nil, err
	}
	return &cls, nil
}

// DeleteClassifier makes a DELETE request to '/v3/classifiers/{id}` endpoint.
func (c *Client) DeleteClassifier(ctx context.Context, id string) (bool, error) {
	if err := c.delete(fmt.Sprintf("/v3/classifiers/%s", id), nil); err != nil {
		return false, err
	}
	return true, nil
}
