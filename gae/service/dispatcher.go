package service

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
)

// Dispatcher is a struct to dispach an http request to one *Service instance
// This class is used when a single GAE service needs to host multiple services.
type Dispatcher struct {
	rules       map[string]*Service
	prefixes    []string
	numServices int
}

// NewDispatcher returns a new *Dispacher object
func NewDispatcher(services ...*Service) *Dispatcher {
	var prefixes []string
	rules := make(map[string]*Service)
	for _, s := range services {
		prefix := "/" + s.URLPrefix()
		if _, ok := rules[prefix]; ok {
			panic(fmt.Errorf("service prefix /%q is duplicated", prefix))
		}
		rules[prefix] = s
		prefixes = append(prefixes, prefix)
	}
	sort.StringSlice(prefixes).Sort()
	return &Dispatcher{
		rules:       rules,
		prefixes:    prefixes,
		numServices: len(prefixes),
	}
}

func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// check prefixes in the reverse order.
	// TODO: can be optimized not to check all but some.
	for i := d.numServices; i > 0; i-- {
		prefix := d.prefixes[i-1]
		if strings.HasPrefix(r.URL.EscapedPath(), prefix) {
			d.rules[prefix].ServeHTTP(w, r)
			return
		}
	}
}

// Run initiate http handler with the Dispatcher
func (d *Dispatcher) Run() {
	http.Handle("/", d)
}
