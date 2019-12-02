package gcp

import "os"

// IsOnDevAppServer returns if an app is running on devappserver
func IsOnDevAppServer() bool {
	return os.Getenv("GAE_PARTITION") == "dev"
}

// IsOnAppEngine returns if an app is running on app engine environment
func IsOnAppEngine(projectName string) bool {
	return os.Getenv("GAE_ENV") == "standard" && os.Getenv("GOOGLE_CLOUD_PROJECT") == projectName
}
