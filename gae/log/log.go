package log

import (
	"fmt"

	"github.com/speedland/go/x/xlog"
	"google.golang.org/appengine/log"
)

type LogSink struct {
	formatter *xlog.TextFormatter
}

// NewLogSink returns a new xlog.Sink for GAE.
func NewLogSink() xlog.Sink {
	return &LogSink{
		formatter: xlog.NewTextFormatter(
			`{{.Data}}`,
		),
	}
}

// NewLogSinkWithFormatter returns a new xlog.Sink for GAE with the given text formatter.
func NewLogSinkWithFormatter(f *xlog.TextFormatter) xlog.Sink {
	return &LogSink{
		formatter: f,
	}
}

func (s *LogSink) Write(r *xlog.Record) error {
	fmt.Println("FOOO")
	ctx := r.Context()
	if ctx != nil {
		return fmt.Errorf("log context is nil")
	}
	buff, err := s.formatter.Format(r)
	if err != nil {
		return err
	}
	switch r.Level {
	case xlog.LevelDebug, xlog.LevelTrace:
		log.Debugf(ctx, string(buff))
	case xlog.LevelInfo:
		log.Infof(ctx, string(buff))
	case xlog.LevelWarn:
		log.Warningf(ctx, string(buff))
	case xlog.LevelError:
		log.Errorf(ctx, string(buff))
	case xlog.LevelFatal:
		log.Criticalf(ctx, string(buff))
	default:
		return fmt.Errorf("unsupported log level: %q", r.Level)
	}
	return nil
}