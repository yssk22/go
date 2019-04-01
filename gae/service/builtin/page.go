package builtin

import (
	"github.com/yssk22/go/gae/service"
	"github.com/yssk22/go/web/response/view/reactapp"
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
		reactapp.ReactAppPath("/static/hatachi.js"),
	))
}
