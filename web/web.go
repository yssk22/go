// Package web provides simple web application framework for golang
package web

import "github.com/speedland/go/web/response"

// LoggerKeys
const (
	LoggerKeyRouter = "web.router"
)

// NotFound is the default response for 404
var NotFound = response.NewTextWithStatus("not found", response.HTTPStatusNotFound)
