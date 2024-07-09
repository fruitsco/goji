package crypt

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("crypt",
		fx.Decorate(logging.NamedLogger("crypt")),

		fx.Supply(cfg),
		fx.Provide(New),
	)
}
