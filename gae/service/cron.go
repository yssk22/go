package service

import (
	"bytes"
	"fmt"
	"io"

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

// GenCronYAML generates cron yaml content to `w`
func (s *Service) GenCronYAML(w io.Writer) {
	fmt.Fprintf(w, "# Service -- %s\n", s.Key())
	for _, c := range s.crons {
		fmt.Fprintf(w, "%s", c.ToYAML())
	}
}
