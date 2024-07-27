package redis

import (
	"go.uber.org/fx"
)

func Module(config *Config) fx.Option {
	return fx.Module("redis",
		fx.Supply(config),
		fx.Provide(New),
	)
}
