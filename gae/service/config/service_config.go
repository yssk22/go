package config

import (
	"time"

	"github.com/yssk22/go/x/xtime"
)

// ServiceConfig is a configuration object for a service
// @flow
// @datastore
type ServiceConfig struct {
	ID        string    `json:"id" ent:"id"`
	Key       string    `json:"key"` // deprecated
	Value     string    `json:"value" ent:"form" datastore:",noindex"`
	UpdatedAt time.Time `json:"updated_at" ent:"timestamp" datastore:",noindex"`

	Description  string `json:"description" datastore:"-"`
	DefaultValue string `json:"default_value" datastore:"-"`
	GlobalValue  string `json:"global_value" datastore:"-"`
	isGlobal     bool   // if true, ServiceConfig will be stored in a service namespace and global namespace
}

func newServiceConfig(key string, value string, description string) *ServiceConfig {
	return &ServiceConfig{
		ID:          key,
		Key:         key,
		Value:       value,
		Description: description,
		UpdatedAt:   xtime.Now(),
	}
}
