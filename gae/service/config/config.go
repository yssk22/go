package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/speedland/go/keyvalue"
	"github.com/speedland/go/x/xerrors"
	"context"
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
	serviceConfigs := DefaultServiceConfigKind.MustGetMulti(
		ctx,
		c.defaultKeys,
	)
	globalCtx, err := appengine.Namespace(ctx, "")
	xerrors.MustNil(err)
	globalConfigs := DefaultServiceConfigKind.MustGetMulti(
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
	serviceCfg := DefaultServiceConfigKind.MustGet(ctx, key)
	defaultCfg := c.getDefault(key)
	globalCtx, err := appengine.Namespace(ctx, "")
	xerrors.MustNil(err)
	globalCfg := DefaultServiceConfigKind.MustGet(globalCtx, key)
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

// Set sets the *ServiceConfig
func (c *Config) Set(ctx context.Context, cfg *ServiceConfig) {
	DefaultServiceConfigKind.MustPut(ctx, cfg)
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
	DefaultServiceConfigKind.MustPutMulti(ctx, updates)
}

func (c *Config) normalize(s *ServiceConfig, global *ServiceConfig, default_ *ServiceConfig) *ServiceConfig {
	// normalize *ServiceConfig object with fallback if it's nil.
	// global can be nil if it is not a global config.
	// default_ must not be nil.
	if global == nil {
		global = default_
	}
	if s == nil {
		if default_.isGlobal {
			s = &ServiceConfig{
				Key:   global.Key,
				Value: global.Value,
			}
		} else {
			s = &ServiceConfig{
				Key:   default_.Key,
				Value: default_.Value,
			}
		}
	} else {
		s = &ServiceConfig{
			Key:   s.Key,
			Value: s.Value,
		}
	}
	s.Description = default_.Description
	s.DefaultValue = default_.Value
	if default_.isGlobal {
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
