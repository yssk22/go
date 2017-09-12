package builtin

import (
	"github.com/speedland/go/gae/service"
	"github.com/speedland/go/web/response/view/react"
)

func setupAdminConfigPages(s *service.Service) {
	if s.APIConfig == nil || s.APIConfig.ConfigAPIBasePath == "" {
		return
	}
	if s.PageConfig == nil || s.PageConfig.AdminConfigPath == "" {
		return
	}
	s.Page(s.PageConfig.AdminConfigPath,
		react.Must(react.New(
			react.Title("Service Configurations"),
			react.Config(&react.PageConfig{
				ReactModulePath: "builtins/configs/",
			}),
			react.AppData("apiBasePath", s.Path(s.APIConfig.ConfigAPIBasePath)),
		)))
}

func setupAdminAsyncTaskPages(s *service.Service) {
	if s.APIConfig == nil || s.APIConfig.AsyncTaskListAPIPath == "" {
		return
	}
	if s.PageConfig == nil || s.PageConfig.AdminAsyncTaskPath == "" {
		return
	}
	s.Page(s.PageConfig.AdminAsyncTaskPath,
		react.Must(react.New(
			react.Title("Async Tasks"),
			react.Config(&react.PageConfig{
				ReactModulePath: "builtins/asynctasks/",
			}),
			react.AppData("asyncTaskListAPIPath", s.Path(s.APIConfig.AsyncTaskListAPIPath)),
		)))
}
