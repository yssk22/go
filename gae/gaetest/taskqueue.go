package gaetest

import (
	"golang.org/x/net/context"

	xtaskqueue "github.com/speedland/go/gae/taskqueue"
	"google.golang.org/appengine/taskqueue"
)

func prepareTaskQueueInTest() {
	xtaskqueue.Add = func(ctx context.Context, t *taskqueue.Task, queueName string) (*taskqueue.Task, error) {
		return taskqueue.Add(ctx, t, "default")
	}
}
