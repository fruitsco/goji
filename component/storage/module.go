package storage

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("storage",
		fx.Decorate(logging.NamedLogger("storage")),

		// noop driver
		fx.Provide(NewNoOpDriverFactory),

		// base
		fx.Supply(cfg),
		fx.Provide(New),
	)
}
