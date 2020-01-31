package appengine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/yssk22/go/gcp"
	"github.com/yssk22/go/gcp/stackdriver"
	"github.com/yssk22/go/x/xlog"
	"github.com/yssk22/go/x/xtime"

	"go.opencensus.io/exporter/stackdriver/propagation"
)

var httpFormat = &propagation.HTTPFormat{}

// Server represents appengine web server
type Server struct {
	ProjectID string
}

// ListenAndServe serves appengine server with h
func (s *Server) ListenAndServe(h http.Handler) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(
		fmt.Sprintf(":%s", port),
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ww := &responseWriterWithLog{
				inner:        w,
				requestStart: xtime.Now(),
			}
			h.ServeHTTP(ww, req)
			if !gcp.IsOnDevAppServer() {
				s.logRequest(ww, req)
			}
		}),
	))
}

// NewLogSink returns xlog.Sink to write logs combined with request logs.
func NewLogSink(projectID string) xlog.Sink {
	return &logSink{
		projectID: projectID,
	}
}

type logSink struct {
	projectID string
}

type appLog struct {
	Time     string      `json:"time"`
	Message  interface{} `json:"message"`
	Severity string      `json:"severity"`
	Trace    string      `json:"logging.googleapis.com/trace,omitempty"`
}

type requestLog struct {
	Time        string               `json:"time"`
	HTTPRequest *httpRequestLogField `json:"httpRequest"`
	Severity    string               `json:"severity"`
	Trace       string               `json:"logging.googleapis.com/trace,omitempty"`
}

type httpRequestLogField struct {
	Status       int    `json:"status"`
	Method       string `json:"requestMethod"`
	URL          string `json:"requestUrl"`
	ResponseSize int    `json:"responseSize"`
	Latency      string `json:"latency"`
	UserAgent    string `json:"userAgent"`
}

func (s *logSink) Write(r *xlog.Record) error {
	obj := &appLog{
		Time:     r.Timestamp.Format(time.RFC3339Nano),
		Message:  r.Data,
		Severity: stackdriver.LogLevelToSeverityString(r.Level),
	}
	ctx := r.Context()
	if ctx != nil {
		traceID := stackdriver.GetTraceID(ctx)
		if traceID != "" {
			obj.Trace = fmt.Sprintf("projects/%s/traces/%s", s.projectID, traceID)
		}
	}
	line, _ := json.Marshal(obj)
	fmt.Println(string(line))
	return nil
}

func (s *Server) logRequest(w *responseWriterWithLog, req *http.Request) {
	t1 := xtime.Now()
	elapsed := t1.Sub(w.requestStart)
	obj := &requestLog{
		Time: t1.Format(time.RFC3339Nano),
		HTTPRequest: &httpRequestLogField{
			Status:       int(w.statusCode),
			Method:       req.Method,
			URL:          req.URL.RequestURI(),
			ResponseSize: w.bytesWritten,
			UserAgent:    req.UserAgent(),
			Latency:      fmt.Sprintf("%fs", elapsed.Seconds()),
		},
		Severity: "INFO",
	}
	sc, ok := httpFormat.SpanContextFromRequest(req)
	if ok {
		obj.Trace = fmt.Sprintf("projects/%s/traces/%s", s.ProjectID, sc.TraceID.String())
	}
	line, _ := json.Marshal(obj)
	fmt.Fprintln(os.Stderr, string(line))
}

type responseWriterWithLog struct {
	inner        http.ResponseWriter
	statusCode   int
	bytesWritten int
	requestStart time.Time
}

func (w *responseWriterWithLog) Header() http.Header {
	return w.inner.Header()
}

func (w *responseWriterWithLog) Write(b []byte) (int, error) {
	written, err := w.inner.Write(b)
	w.bytesWritten += written
	return written, err
}

func (w *responseWriterWithLog) WriteHeader(statusCode int) {
	w.inner.WriteHeader(statusCode)
	w.statusCode = statusCode
}
