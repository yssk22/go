// Package web provides simple web application framework for golang
package web

import (
	"fmt"

	"github.com/speedland/go/web/response"
)

// LoggerKeys
const (
	LoggerKeyRouter = "web.router"
)

type contextKey struct {
	key string
}

func (c *contextKey) String() string {
	return fmt.Sprintf("ContextKey(github.com/speedland/go/web) %s", c.key)
}

// NotFound is the default response for 404
var NotFound = response.NewTextWithStatus("not found", response.HTTPStatusNotFound)
