package builtin

import (
	"context"

	"github.com/yssk22/go/gae/service"
	"github.com/yssk22/go/gae/service/config"
)

// @api path=/admin/api/configs/
func listConfigs(ctx context.Context) ([]*config.ServiceConfig, error) {
	s := service.FromContext(ctx)
	return s.Config.All(ctx), nil
}

// @api path=/admin/api/configs/:key.json
func getConfig(ctx context.Context, key string) (*config.ServiceConfig, error) {
	s := service.FromContext(ctx)
	return s.Config.Get(ctx, key), nil
}

// @api path=/admin/api/configs/:key.json
func updateConfig(ctx context.Context, key string, newConfig *config.ServiceConfig) (*config.ServiceConfig, error) {
	s := service.FromContext(ctx)
	cfg := s.Config.Get(ctx, key)
	if cfg == nil {
		return nil, nil
	}
	cfg.Value = newConfig.Value
	s.Config.Set(ctx, cfg)
	return cfg, nil
}
