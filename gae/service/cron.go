package service

import (
	"bytes"
	"fmt"

	"github.com/speedland/go/web/response"

	"github.com/speedland/go/web"
)

// CronTimezone is a timezone used in cron.yaml
var CronTimezone = "Asia/Tokyo"

type cron struct {
	path        string
	time        string
	description string
	handlers    []web.Handler
}

func (c *cron) ToYAML() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "- url:         %s\n", c.path)
	fmt.Fprintf(&buff, "  schedule:    %s\n", c.time)
	fmt.Fprintf(&buff, "  description: %s\n", c.description)
	fmt.Fprintf(&buff, "  timezone:    %s\n", CronTimezone)
	return buff.String()
}

func (s *Service) AddCron(path, time, desc string, handlers ...web.Handler) {
	c := &cron{
		path:        s.Path(path),
		time:        time,
		description: desc,
		handlers:    handlers,
	}
	s.crons = append(s.crons, c)
	s.Get(path, handlers...)
}

func (s *Service) serveCronYAML(r *web.Request, _ web.NextHandler) *response.Response {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "# Service -- %s\n", s.Key())
	for _, c := range s.crons {
		fmt.Fprintf(&buff, "%s", c.ToYAML())
	}
	res := response.NewText(buff.String())
	res.Header.Set(response.ContentType, "application/yaml; charset=utf-8")
	return res
}
