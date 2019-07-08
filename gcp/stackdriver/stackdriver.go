package stackdriver

import (
	"context"

	"cloud.google.com/go/logging"
	"github.com/yssk22/go/web"
	"github.com/yssk22/go/x/xlog"
	"go.opencensus.io/exporter/stackdriver/propagation"
	"go.opencensus.io/trace"
)

var levelMap = map[xlog.Level]logging.Severity{
	xlog.LevelTrace: logging.Debug,
	xlog.LevelDebug: logging.Debug,
	xlog.LevelInfo:  logging.Info,
	xlog.LevelWarn:  logging.Warning,
	xlog.LevelError: logging.Error,
	xlog.LevelFatal: logging.Critical,
}

var levelMapString = map[xlog.Level]string{
	xlog.LevelTrace: "DEBUG",
	xlog.LevelDebug: "DEBUG",
	xlog.LevelInfo:  "INFO",
	xlog.LevelWarn:  "WARNING",
	xlog.LevelError: "ERROR",
	xlog.LevelFatal: "CRITICAL",
}

// LogLevelToSeverity converts xlog.Level to logging.Severity
func LogLevelToSeverity(l xlog.Level) logging.Severity {
	return levelMap[l]
}

// LogLevelToSeverityString converts xlog.Level to logging.Severity string
func LogLevelToSeverityString(l xlog.Level) string {
	return levelMapString[l]
}

var httpFormat = &propagation.HTTPFormat{}

// GetTraceID returns a trace identifier in the context
func GetTraceID(ctx context.Context) string {
	span := trace.FromContext(ctx)
	if span != nil {
		return span.SpanContext().TraceID.String()
	}
	// try http format
	req := web.FromContext(ctx)
	if req != nil {
		if sc, ok := httpFormat.SpanContextFromRequest(req.Request); ok {
			return sc.TraceID.String()
		}
	}
	return ""
}
