package builtin

import "github.com/yssk22/go/gae/service"

// Setup sets up builtin APIs and Pages
func Setup(s *service.Service) *service.Service {
	SetupAPI(s)
	setupAdminConfigPages(s)
	setupAdminAsyncTaskPages(s)
	return s
}
