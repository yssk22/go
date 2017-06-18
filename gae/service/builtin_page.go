package service

import "github.com/speedland/go/gae/service/view/react"

// BuiltInPageConfig is a configuration object for ActivateBuiltinPages, which actiate the following pages.
//
// [config]
//
//    - /{ConfigAPIBasePath}/
//
// [asynctask]
//
//    - /{AdminAsyncTaskPath}/
//
type BuiltInPageConfig struct {
	AdminConfigPath    string
	AdminAsyncTaskPath string
}

// ActivateEndpoints sets up builtin API endpoints on the *Service
func (bc *BuiltInPageConfig) ActivateEndpoints(s *Service, apiConfig *BuiltInAPIConfig) *Service {
	bc.activateAdminConfig(s, apiConfig.ConfigAPIBasePath)
	bc.activateAdminAsyncTask(s, apiConfig.AsyncTaskListAPIPath)
	return s
}

// DefaultBuiltinPageConfig is a default object of BuiltInPageConfig
var DefaultBuiltinPageConfig = &BuiltInPageConfig{
	AdminConfigPath:    "/admin/configs/",
	AdminAsyncTaskPath: "/admin/asynctasks/",
}

func (bc *BuiltInPageConfig) activateAdminConfig(s *Service, apiBasePath string) {
	s.Page(bc.AdminConfigPath,
		react.New().
			Title("サービス設定").
			ReactModulePath("builtins/configs/").
			AppData("apiBasePath", s.Path(apiBasePath)))
}

func (bc *BuiltInPageConfig) activateAdminAsyncTask(s *Service, asyncTaskListAPIPath string) {
	s.Page(bc.AdminAsyncTaskPath,
		react.New().
			Title("AsyncTask一覧").
			ReactModulePath("builtins/asynctasks/").
			AppData("asyncTaskListAPIPath", s.Path(asyncTaskListAPIPath)))
}
