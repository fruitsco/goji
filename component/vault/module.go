package vault

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/x/logging"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("vault",
		fx.Decorate(logging.NamedLogger("vault")),

		fx.Supply(cfg),
		fx.Provide(New),
	)
}
