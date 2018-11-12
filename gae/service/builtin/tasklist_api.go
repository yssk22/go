package builtin

import (
	"context"

	"github.com/yssk22/go/gae/service"
)

// @api path=/admin/api/tasks/
func listAsyncTasks(ctx context.Context) ([]*service.Task, error) {
	s := service.FromContext(ctx)
	return s.GetTasks(), nil
}
