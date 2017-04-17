package service

import (
	"fmt"
	"io"

	"github.com/speedland/go/gae/taskqueue"
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

// GenQueueYAML generates cron yaml content to `w`
func (s *Service) GenQueueYAML(w io.Writer) {
	fmt.Fprintf(w, "# Service -- %s\n", s.Key())
	for _, q := range s.queues {
		fmt.Fprintf(w, "%s", q.ToYAML())
	}
}
