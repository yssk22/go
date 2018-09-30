// Package taskqueue provides a utility for taskqueue
package taskqueue

import (
	"bytes"
	"fmt"
	"net/url"

	"context"

	"github.com/yssk22/go/web"
	"github.com/yssk22/go/web/response"
	"github.com/yssk22/go/x/xlog"
	"google.golang.org/appengine"
	"google.golang.org/appengine/taskqueue"
)

const LoggerKey = "gae.taskqueue"

// QueueMode is a type alias for queue mode string
type QueueMode string

// Available constants for QueueMode.
const (
	QueueModePuth QueueMode = "push"
	QueueModePull           = "pull"
)

// PushQueue is a struct to define push queue.
type PushQueue struct {
	Name                  string
	BucketSize            string
	MaxConcurrentRequests string
	Rate                  string
	RetryLimit            string
	AgeLimit              string
	MinBackoffSeconds     string
	MaxBackoffSeconds     string
	MaxDoubling           string
	Target                string
}

// ToYAML generates the string for the queue
func (queue *PushQueue) ToYAML() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "- name: %s\n", queue.Name)
	if queue.BucketSize != "" {
		fmt.Fprintf(&buff, "  bucket_size: %s\n", queue.BucketSize)
	}
	if queue.MaxConcurrentRequests != "" {
		fmt.Fprintf(&buff, "  max_concurrent_requests: %s\n", queue.MaxConcurrentRequests)
	}
	if queue.Rate != "" {
		fmt.Fprintf(&buff, "  rate: %s\n", queue.Rate)
	}
	if queue.RetryLimit != "" || queue.AgeLimit != "" ||
		queue.MinBackoffSeconds != "" || queue.MaxBackoffSeconds != "" ||
		queue.MaxDoubling != "" {
		fmt.Fprintf(&buff, "  retry_parameters:\n")
		if queue.RetryLimit != "" {
			fmt.Fprintf(&buff, "    task_retry_limit: %s\n", queue.RetryLimit)
		}
		if queue.AgeLimit != "" {
			fmt.Fprintf(&buff, "    task_age_limit: %s\n", queue.AgeLimit)
		}
		if queue.MinBackoffSeconds != "" {
			fmt.Fprintf(&buff, "    min_backoff_seconds: %s\n", queue.MinBackoffSeconds)
		}
		if queue.MaxBackoffSeconds != "" {
			fmt.Fprintf(&buff, "    max_backoff_seconds: %s\n", queue.MaxBackoffSeconds)
		}
		if queue.MaxDoubling != "" {
			fmt.Fprintf(&buff, "    max_doubling: %s\n", queue.MaxDoubling)
		}
	}
	return buff.String()
}

// PushTask to push a task into the queue.
func (queue *PushQueue) PushTask(ctx context.Context, urlPath string, form url.Values) error {
	// aetest environment does not support non-default queue
	// https://code.google.com/p/googleappengine/issues/detail?id=10771
	var queueName = queue.Name
	if appengine.IsDevAppServer() {
		queueName = "default"
	}
	if _, err := Add(ctx, taskqueue.NewPOSTTask(urlPath, form), queueName); err != nil {
		return fmt.Errorf("failed to push queue: %s, url: %s, - %v", queueName, urlPath, err)
	}
	return nil
}

// RequestValidator returns a web.Handler to validate the request targets the task or not.
// This is useful for handlers throttled via PushTask.
func (queue *PushQueue) RequestValidator() web.Handler {
	return web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		if !appengine.IsDevAppServer() {
			name := req.Header.Get("X-AppEngine-QueueName")
			if queue.Name != name {
				_, logger := xlog.WithContextAndKey(req.Context(), "", LoggerKey)
				logger.Warnf("Task Queue invalidation: %q != %q", queue.Name, name)
				return response.NewErrorWithStatus(
					fmt.Errorf("task queue validation failed"),
					response.HTTPStatusBadRequest,
				)
			}
		}
		return next(req)
	})
}
