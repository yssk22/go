package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"context"

	"github.com/yssk22/go/keyvalue"
	"github.com/yssk22/go/x/xerrors"
	"google.golang.org/appengine"
)

// Config is a struct to access configurations represented by *ServiceConfig
type Config struct {
	defaultMap  keyvalue.StringKeyMap
	defaultKeys []string
}

// New returns a new *Config object
func New() *Config {
	c := &Config{
		defaultMap:  keyvalue.NewStringKeyMap(),
		defaultKeys: make([]string, 0),
	}
	for _, sc := range globalDefaults {
		c.register(sc)
	}
	return c
}

// Register registers configuration variable in application
func (c *Config) Register(key string, defaultValue string, description string) {
	cfg := newServiceConfig(key, defaultValue, description)
	c.register(cfg)
}

func (c *Config) register(cfg *ServiceConfig) {
	ok, err := c.defaultMap.Get(cfg.Key)
	if err == nil && ok != nil {
		panic(fmt.Errorf("key %q is already registred (global=%t)", cfg.Key, (ok.(*ServiceConfig)).isGlobal))
	}
	c.defaultMap.Set(cfg.Key, cfg)
	c.defaultKeys = append(c.defaultKeys, cfg.Key)
}

// All returns all *ServiceConfig in app.
func (c *Config) All(ctx context.Context) []*ServiceConfig {
	datastore := NewServiceConfigKind()
	_, serviceConfigs := datastore.MustGetMulti(
		ctx,
		c.defaultKeys,
	)
	globalCtx, err := appengine.Namespace(ctx, "")
	xerrors.MustNil(err)
	_, globalConfigs := datastore.MustGetMulti(
		globalCtx,
		c.defaultKeys,
	)
	for i := range serviceConfigs {
		defaultCfg := c.getDefault(c.defaultKeys[i])
		globalCfg := globalConfigs[i]
		serviceConfigs[i] = c.normalize(serviceConfigs[i], globalCfg, defaultCfg)
	}
	return serviceConfigs
}

// Get gets the *ServiceConfig
func (c *Config) Get(ctx context.Context, key string) *ServiceConfig {
	datastore := NewServiceConfigKind()
	_, serviceCfg := datastore.MustGet(ctx, key)
	defaultCfg := c.getDefault(key)
	globalCtx, err := appengine.Namespace(ctx, "")
	xerrors.MustNil(err)
	_, globalCfg := datastore.MustGet(globalCtx, key)
	return c.normalize(serviceCfg, globalCfg, defaultCfg)
}

// GetValue is like Get but returns only the value as string.
func (c *Config) GetValue(ctx context.Context, key string) string {
	return c.Get(ctx, key).Value
}

// GetDefaultValue returns the default value by `key`
func (c *Config) GetDefaultValue(key string) string {
	return c.getDefault(key).Value
}

// GetIntValue is like GetValue and return the value as int. If invalid int value is set on `key`
// this will return a default value of `key`.
func (c *Config) GetIntValue(ctx context.Context, key string) int {
	v := c.GetValue(ctx, key)
	vv, err := strconv.Atoi(v)
	if err == nil {
		return vv
	}
	return c.GetIntDefaultValue(key)
}

// GetIntDefaultValue returns the default value by `key` as int
func (c *Config) GetIntDefaultValue(key string) int {
	vv, err := strconv.Atoi(c.GetDefaultValue(key))
	xerrors.MustNil(err)
	return vv
}

// GetFloatValue is like GetValue and return the value as float64. If invalid int value is set on `key`
// this will return a default value of `key`.
func (c *Config) GetFloatValue(ctx context.Context, key string) float64 {
	v := c.GetValue(ctx, key)
	vv, err := strconv.ParseFloat(v, 64)
	if err == nil {
		return vv
	}
	return c.GetFloatDefaultValue(key)
}

// GetFloatDefaultValue returns the default value by `key` as float64
func (c *Config) GetFloatDefaultValue(key string) float64 {
	vv, err := strconv.ParseFloat(c.GetDefaultValue(key), 64)
	xerrors.MustNil(err)
	return vv
}

// GetBoolValue is like GetValue and return the value as bool. If invalid int value is set on `key`
// this will return a default value of `key`.
func (c *Config) GetBoolValue(ctx context.Context, key string) bool {
	v := c.GetValue(ctx, key)
	vv, err := strconv.ParseBool(v)
	if err == nil {
		return vv
	}
	return c.GetBoolDefaultValue(key)
}

// GetBoolDefaultValue returns the default value by `key` as bool
func (c *Config) GetBoolDefaultValue(key string) bool {
	vv, err := strconv.ParseBool(c.GetDefaultValue(key))
	xerrors.MustNil(err)
	return vv
}

// Set sets the *ServiceConfig
func (c *Config) Set(ctx context.Context, cfg *ServiceConfig) {
	NewServiceConfigKind().MustPut(ctx, cfg)
}

// SetValue set the new value on key
func (c *Config) SetValue(ctx context.Context, key string, value string) {
	cfg := c.Get(ctx, key)
	cfg.Value = value
	NewServiceConfigKind().MustPut(ctx, cfg)
}

// LoadFromJSON loads config values from a given file path
func (c *Config) LoadFromJSON(ctx context.Context, path string) {
	f, err := os.Open(path)
	xerrors.MustNil(err)
	defer f.Close()
	m := make(map[string]string)
	d := json.NewDecoder(f)
	xerrors.MustNil(d.Decode(m))
	var updates []*ServiceConfig
	for key, value := range m {
		cfg := c.Get(ctx, key)
		cfg.Value = value
		updates = append(updates, cfg)
	}
	NewServiceConfigKind().MustPutMulti(ctx, updates)
}

func (c *Config) normalize(s *ServiceConfig, global *ServiceConfig, defaultCfg *ServiceConfig) *ServiceConfig {
	// normalize *ServiceConfig object with fallback if it's nil.
	// global can be nil if it is not a global config.
	// default_ must not be nil.
	if global == nil {
		global = defaultCfg
	}
	if s == nil {
		if defaultCfg.isGlobal {
			s = &ServiceConfig{
				Key:   global.Key,
				Value: global.Value,
			}
		} else {
			s = &ServiceConfig{
				Key:   defaultCfg.Key,
				Value: defaultCfg.Value,
			}
		}
	} else {
		s = &ServiceConfig{
			Key:   s.Key,
			Value: s.Value,
		}
	}
	s.Description = defaultCfg.Description
	s.DefaultValue = defaultCfg.Value
	if defaultCfg.isGlobal {
		s.GlobalValue = global.Value
	}
	return s
}

func (c *Config) getDefault(key string) *ServiceConfig {
	cfg, err := c.defaultMap.Get(key)
	if err != nil {
		panic(fmt.Errorf("unexpected error while accessing %q (available keys: %q): %v", key, c.defaultKeys, err))
	}
	return cfg.(*ServiceConfig)
}
