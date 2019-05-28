package service

import (
	"encoding/json"
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

type dispatchAPIResponse struct {
	ID     string `json:"id"`
	Prefix string `json:"prefix"`
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
	path := r.URL.EscapedPath()
	if strings.HasPrefix(path, "/__services/") {
		var list []*dispatchAPIResponse
		for _, p := range d.prefixes {
			s := d.rules[p]
			list = append(list, &dispatchAPIResponse{
				ID:     s.Key(),
				Prefix: fmt.Sprintf("/%s", s.URLPrefix()),
			})
		}
		buff, _ := json.Marshal(list)
		w.WriteHeader(200)
		w.Write(buff)
		return
	}

	// check prefixes in the reverse order.
	// TODO: can be optimized not to check all but some.
	for i := d.numServices; i > 0; i-- {
		prefix := d.prefixes[i-1]
		if strings.HasPrefix(path, prefix) {
			d.rules[prefix].ServeHTTP(w, r)
			return
		}
	}
}

// Run initiate http handler with the Dispatcher
func (d *Dispatcher) Run() {
	http.Handle("/", d)
}
