package crypt

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/x/logging"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("crypt",
		fx.Decorate(logging.NamedLogger("crypt")),

		fx.Provide(NewVaultKeyProvider),

		fx.Supply(cfg),
		fx.Provide(New),
	)
}
