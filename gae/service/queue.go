package service

import (
	"bytes"
	"fmt"

	"github.com/speedland/go/gae/taskqueue"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/response"
)

// AddPushQueue adds a named queue into the service.
func (s *Service) AddPushQueue(name string) *taskqueue.PushQueue {
	queue := &taskqueue.PushQueue{}
	queue.Name = fmt.Sprintf("%s-%s", s.Key(), name)
	// default value, see the doc: https://cloud.google.com/appengine/docs/go/config/queueref
	queue.BucketSize = "5"
	queue.MaxConcurrentRequests = "1000"
	queue.Rate = "1/s"
	s.queues = append(s.queues, queue)
	return queue
}

func (s *Service) serveQueueYAML(r *web.Request, _ web.NextHandler) *response.Response {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "# Service -- %s\n", s.Key())
	for _, q := range s.queues {
		fmt.Fprintf(&buff, "%s", q.ToYAML())
	}
	res := response.NewText(buff.String())
	res.Header.Set(response.ContentType, "application/yaml; charset=utf-8")
	return res
}
