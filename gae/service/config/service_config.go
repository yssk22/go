package config

import (
	"time"

	"github.com/yssk22/go/x/xtime"
)

// ServiceConfig is a configuration object for a service
//go:generate ent -type=ServiceConfig
type ServiceConfig struct {
	Key       string    `json:"key" ent:"id"`
	Value     string    `json:"value" ent:"form" datastore:",noindex"`
	UpdatedAt time.Time `json:"updated_at" ent:"timestamp" datastore:",noindex"`

	Description  string `json:"description" datastore:"-"`
	DefaultValue string `json:"default_value" datastore:"-"`
	GlobalValue  string `json:"global_value" datastore:"-"`
	isGlobal     bool   // if true, ServiceConfig will be stored in a service namespace and global namespace
}

func newServiceConfig(key string, value string, description string) *ServiceConfig {
	return &ServiceConfig{
		Key:         key,
		Value:       value,
		Description: description,
		UpdatedAt:   xtime.Now(),
	}
}
