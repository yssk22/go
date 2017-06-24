// Package service provides a gae service instance framework on top of
// github.com/speedland/go/web package.
//
// Using this package, what you need to do in your GAE app looks like this:
//
//     // app.go
//
//     func init(){
//         s := service.New("serviceKey")
//         s.Get("/path/to/endpoint/", web.HandlerFunc(...))  // register /serviceKey/path/to/endpoint handler
//         ...
//         s.Run()
//     }
//
package service

import (
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/speedland/go/gae/service/apierrors"

	"sync"

	"os"

	"github.com/speedland/go/gae/service/asynctask"
	"github.com/speedland/go/gae/service/config"
	"github.com/speedland/go/gae/service/view"
	xtaskqueue "github.com/speedland/go/gae/taskqueue"
	"github.com/speedland/go/uuid"
	"github.com/speedland/go/web"
	"github.com/speedland/go/web/middleware/session"
	"github.com/speedland/go/web/response"
	"github.com/speedland/go/x/xcontext"
	"github.com/speedland/go/x/xerrors"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// ContextKey is a key to get a service.
var ContextKey = xcontext.NewKey("service")

// Service is a set of endpoints
type Service struct {
	Init      func(*web.Request)
	Every     func(*web.Request)
	OnError   func(*web.Request, error) *response.Response
	Config    *config.Config
	once      sync.Once // for Init control
	key       string    // service key
	urlPrefix string    // url base path
	namespace string    // datastore/memcache namespace for services
	crons     []*cron
	queues    []*xtaskqueue.PushQueue
	tasks     []*asynctaskEndpoint
	router    *web.Router // service router
}

// FromContext returns a service object associated with the context
func FromContext(ctx context.Context) *Service {
	service, ok := ctx.Value(ContextKey).(*Service)
	if ok {
		return service
	}
	return nil
}

// WithContext returns a new context.Context associated with the service
func WithContext(ctx context.Context, s *Service) context.Context {
	return context.WithValue(ctx, ContextKey, s)
}

// MustFromContext is like FromContext but panics if a service is not in the context
func MustFromContext(ctx context.Context) *Service {
	service, ok := ctx.Value(ContextKey).(*Service)
	if !ok {
		panic(fmt.Errorf("not a service context"))
	}
	return service
}

// New returns a new *Service instance
func New(key string) *Service {
	return NewWithURLAndNamespace(
		key,
		strings.Replace(key, "-", "/", -1),
		strings.Replace(key, "-", ".", -1),
	)
}

// NewWithURLAndNamespace is like New with using the given url prefix and namespece instead of 'key' value.
func NewWithURLAndNamespace(key string, url string, namespace string) *Service {
	if key == "" {
		panic(fmt.Errorf("service key must not be nil"))
	}
	s := &Service{
		key:       key,
		urlPrefix: url,
		namespace: namespace,
		router:    web.NewRouter(nil),
		Config:    config.New(),
	}
	s.router.Use(namespaceMiddleware(s))
	s.router.Use(errorMiddleware)
	s.Use(initMiddleware)
	s.Use(session.Default)
	s.Use(everyMiddleware)
	return s
}

// Run register the service on http.Hander
func (s *Service) Run() {
	http.Handle("/", s.router)
}

// ServeHTTP implements http.Handler#ServeHTTP
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Key returns a key string
func (s *Service) Key() string {
	return s.key
}

// URLPrefix returns a url prefix string
func (s *Service) URLPrefix() string {
	return s.urlPrefix
}

// Namespace returns a namespace string
func (s *Service) Namespace() string {
	return s.namespace
}

// WithNamespace sets the namespace of the given context
func (s *Service) WithNamespace(ctx context.Context) context.Context {
	ctx, err := appengine.Namespace(ctx, s.namespace)
	if err != nil {
		panic(err)
	}
	return ctx
}

// Use adds the middleware onto the service router
func (s *Service) Use(handlers ...web.Handler) {
	s.router.Use(handlers...)
}

// Get defines an endpoint for GET
func (s *Service) Get(path string, handlers ...web.Handler) {
	s.router.Get(s.Path(path), handlers...)
}

// Post defines an endpoint for POST
func (s *Service) Post(path string, handlers ...web.Handler) {
	s.router.Post(s.Path(path), handlers...)
}

// Put defines an endpoint for PUT
func (s *Service) Put(path string, handlers ...web.Handler) {
	s.router.Put(s.Path(path), handlers...)
}

// Delete defines an endpoint for DELETE
func (s *Service) Delete(path string, handlers ...web.Handler) {
	s.router.Delete(s.Path(path), handlers...)
}

// Page defines an endpoint for view.Page interaface
func (s *Service) Page(path string, p view.Page) {
	s.Get(path, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return p.Render(req)
	}))
}

// AsyncTask defines endpoints for asynctask execution
func (s *Service) AsyncTask(path string, taskConfig *asynctask.Config) {
	if !strings.HasSuffix(path, "/") {
		panic(fmt.Errorf("AsyncTask path must ends with '/' (got %q)", path))
	}
	fullPath := s.Path(path)
	s.Get(fmt.Sprintf("%s:taskid.json", path), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		status := taskConfig.GetStatus(req.Context(), req.Params.GetStringOr("taskid", ""))
		if status == nil {
			return nil
		}
		return response.NewJSON(status)
	}))

	const TaskQueueHeader = "X-AppEngine-TaskName"
	s.Post(fmt.Sprintf("%s:taskid.json", path), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		if req.Header.Get(TaskQueueHeader) == "" {
			return apierrors.Forbidden.ToResponse()
		}
		taskID := req.Params.GetStringOr("taskid", "")
		progress, err := taskConfig.Process(req.Context(), taskID, fmt.Sprintf("%s%s.json", fullPath, taskID), req.Request.URL.Query())
		if err != nil {
			if err == asynctask.ErrNoTaskInstance {
				return nil
			}
			return (&apierrors.Error{
				Code:    "invalid_asynctask_call",
				Message: err.Error(),
				Status:  response.HTTPStatusBadRequest,
			}).ToResponse()
		}
		if progress == nil {
			return response.NewJSON(true)
		}
		return response.NewJSON(progress.Next)
	}))

	s.Post(path, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		taskID := uuid.New().String()
		status, err := taskConfig.Prepare(req.Context(), taskID, fmt.Sprintf("%s%s.json", fullPath, taskID), req.Request.URL.Query())
		xerrors.MustNil(err)
		return response.NewJSONWithStatus(status, response.HTTPStatusCreated)
	}))

	s.Get(path, web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
		return response.NewJSON(taskConfig.GetRecentTasks(req.Context(), req.Query.GetIntOr("n", 5)))
	}))

	if schedule := taskConfig.GetSchedule(); schedule != "" {
		const CronHeader = "X-AppEngine-Cron"
		s.AddCron(fmt.Sprintf("%s/cron/", path), schedule, taskConfig.GetDescription(), web.HandlerFunc(func(req *web.Request, next web.NextHandler) *response.Response {
			if req.Header.Get(CronHeader) == "" {
				return apierrors.Forbidden.ToResponse()
			}
			taskID := uuid.New().String()
			params := req.Request.URL.Query()
			params.Set("cron", "true")
			status, err := taskConfig.Prepare(req.Context(), taskID, fmt.Sprintf("%s:%s.json", fullPath, taskID), params)
			xerrors.MustNil(err)
			return response.NewJSONWithStatus(status, response.HTTPStatusCreated)
		}))
	}
	s.tasks = append(s.tasks, &asynctaskEndpoint{
		path:   s.Path(path),
		config: taskConfig,
	})
}

// Path returns an absolute path for this s.
func (s *Service) Path(p string) string {
	if s.urlPrefix != "" {
		if path.Ext(p) == "" {
			return path.Join("/", s.urlPrefix, p) + "/"
		}
		fmt.Fprintln(os.Stderr, path.Join("/", s.urlPrefix, p))
		return path.Join("/", s.urlPrefix, p)
	}
	if p == "/" {
		return p
	}
	if path.Ext(p) == "" {
		return path.Join(p) + "/"
	}
	return p
}
