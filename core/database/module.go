package database

import (
	"go.uber.org/fx"
)

func Module(config *Config) fx.Option {
	return fx.Module("database",
		fx.Supply(config),
		fx.Provide(NewDB),
		fx.Provide(NewMig),
	)
}
