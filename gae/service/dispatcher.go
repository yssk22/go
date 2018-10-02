package service

import (
	"net/http"
	"strings"
)

// Dispatcher is a struct to dispach an http request to one *Service instance
// This class is used when a single GAE service needs to host multiple services.
type Dispatcher struct {
	rules map[string]*Service
}

// NewDispatcher returns a new *Dispacher object
func NewDispatcher(services ...*Service) *Dispatcher {
	rules := make(map[string]*Service)
	for _, s := range services {
		rules["/"+s.URLPrefix()] = s
	}
	return &Dispatcher{
		rules: rules,
	}
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for prefix, service := range d.rules {
		if strings.HasPrefix(r.URL.EscapedPath(), prefix) {
			service.router.ServeHTTP(w, r)
			return
		}
	}
}

// Run initiate http handler with the Dispatcher
func (d *Dispatcher) Run() {
	http.Handle("/", d)
}
