package service

import "github.com/speedland/go/gae/service/view/react"

// BuiltInPageConfig is a configuration object for ActivateBuiltinPages, which actiate the following endpoints
//
// [config]
//
//    - GET /{ConfigAPIBasePath}/
//    - GET /{ConfigAPIBasePath}/:key.json
//    - PUT /{ConfigAPIBasePath}/:key.json
//
type BuiltInPageConfig struct {
	AdminConfigPath string
}

// ActivateEndpoints sets up builtin API endpoints on the *Service
func (bc *BuiltInPageConfig) ActivateEndpoints(s *Service, apiConfig *BuiltInAPIConfig) *Service {
	bc.activateAdminConfig(s, apiConfig.ConfigAPIBasePath)
	return s
}

// DefaultBuiltinPageConfig is a default object of BuiltInPageConfig
var DefaultBuiltinPageConfig = &BuiltInPageConfig{
	AdminConfigPath: "/admin/configs/",
}

func (bc *BuiltInPageConfig) activateAdminConfig(s *Service, apiBasePath string) {
	s.Page(bc.AdminConfigPath,
		react.New().
			Title("サービス設定").
			ReactModulePath("builtins/configs/").
			AppData("apiBasePath", s.Path(apiBasePath)))
}
