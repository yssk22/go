package builtin

import (
	"github.com/speedland/go/gae/service"
	"github.com/speedland/go/web/response/view/reactapp"
)

const appName = "admin"

func setupAdminConfigPages(s *service.Service) {
	s.Page("/admin/configs/", newAdminReactApp(s, "Configs"))
}

func setupAdminAsyncTaskPages(s *service.Service) {
	s.Page("/admin/tasks/", newAdminReactApp(s, "Tasks"))
}

func newAdminReactApp(s *service.Service, title string) *reactapp.Page {
	return reactapp.Must(reactapp.New(
		appName,
		reactapp.Title(title),
		reactapp.AppData("urlprefix", s.Path("/")),
	))
}
