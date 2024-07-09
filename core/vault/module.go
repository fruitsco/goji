package vault

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(config *Config) fx.Option {
	return fx.Module("vault",
		fx.Decorate(logging.NamedLogger("vault")),
		fx.Supply(config),
		fx.Provide(New),
	)
}
