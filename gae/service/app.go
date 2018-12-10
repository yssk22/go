package service

import (
	"bytes"
	"fmt"
	"io"
)

// DefaultExpirationValue is a default value for service.DefaultExpiration
const DefaultExpirationValue = "10m"

// HandlerOption to define handlers
type HandlerOption struct {
	URL         string
	Expiration  string
	Login       LoginOption
	StaticFiles string
	StaticDir   string
	Upload      string
	Script      string
}

// ToYAML returns a handler element content for app.yaml
func (option *HandlerOption) ToYAML() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "- url: %s\n", option.URL)
	fmt.Fprintf(&buff, "  secure: always\n") // always secure
	if option.Expiration != "" {
		fmt.Fprintf(&buff, "  expiration: %s\n", option.Expiration)
	}
	if option.Login != "" {
		fmt.Fprintf(&buff, "  login: %s\n", option.Login)
	}
	if option.StaticFiles != "" {
		fmt.Fprintf(&buff, "  static_files: %s\n", option.StaticFiles)
	}
	if option.StaticDir != "" {
		fmt.Fprintf(&buff, "  static_dir: %s\n", option.StaticDir)
	}
	if option.Upload != "" {
		fmt.Fprintf(&buff, "  upload: %s\n", option.Upload)
	}
	if option.StaticFiles == "" && option.StaticDir == "" {
		fmt.Fprintf(&buff, "  script: _go_app\n")
	}
	return buff.String()
}

// LoginOption is a option string for handler login
type LoginOption string

// LoginOption constants
const (
	LoginOptionRequired LoginOption = "required"
	LoginOptionAdmin    LoginOption = "admin"
)

// ToYAML returns a generated app.yaml
func (s *Service) ToYAML() string {
	var buff bytes.Buffer
	fmt.Fprintf(&buff, "service: %s\n", s.key)
	fmt.Fprintf(&buff, "runtime: go111\n")
	fmt.Fprintf(&buff, "handlers:\n")
	for _, option := range s.handlerOptions {
		fmt.Fprintf(&buff, option.ToYAML())
	}
	return buff.String()
}

// EnableRemoteAPI enables remote API via /_ah/remote_api
func (s *Service) EnableRemoteAPI() {
	s.handlerOptions = append(s.handlerOptions, &HandlerOption{
		URL:    "/_ah/remote_api",
		Script: "_go_app",
	})
}

// AddHandlerOption adds an option for handler
func (s *Service) AddHandlerOption(option *HandlerOption) {
	s.handlerOptions = append(s.handlerOptions, option)
}

// GenAppYAML generates app yaml content to `w`
func (s *Service) GenAppYAML(w io.Writer) {
	fmt.Fprintf(w, "# Service -- %s\n", s.Key())
	fmt.Fprintf(w, "%s", s.ToYAML())
}

// GenHandlersYAML write handler definitions into `w`
func (s *Service) GenHandlersYAML(w io.Writer) {
	fmt.Fprintf(w, "# Service -- %s\n", s.Key())
	for _, option := range s.handlerOptions {
		fmt.Fprintf(w, option.ToYAML())
	}
}
