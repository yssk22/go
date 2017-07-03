package builtin

import "github.com/speedland/go/gae/service"

// Setup sets up builtin APIs and Pages
func Setup(s *service.Service) *service.Service {
	// APIs
	setupConfigAPIs(s)
	setupAsyncTaskListAPIs(s)
	setupAuthAPIs(s)

	// pages
	setupAdminConfigPages(s)
	setupAdminAsyncTaskPages(s)
	return s
}
