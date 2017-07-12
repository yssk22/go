package service

// APIConfig is a configuration object for setupBuiltInAPIs, which actiate the following endpoints
//
// [config]
//
//    - GET /{ConfigAPIBasePath}/
//    - GET /{ConfigAPIBasePath}/:key.json
//    - PUT /{ConfigAPIBasePath}/:key.json
//
// [asynctask]
//    - GET /{AsyncTaskListAPIPath}/
//
// [auth]
//    - GET /{AuthAPIBasePath}/me.json
//    - POST /{AuthAPIBasePath}/login/facebook/
//
// [webhook]
//    - GET /{WebhookBasePath}/messenger/
//
type BuiltInAPIConfig struct {
	ConfigAPIBasePath    string
	AsyncTaskListAPIPath string
	AuthAPIBasePath      string
	AuthNamespace        string
	WebhookBasePath      string
}

// BuiltInPageConfig is a configuration object for setupBuiltinPages, which actiate the following pages.
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
