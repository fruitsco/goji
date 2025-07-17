package email

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("email",
		fx.Decorate(logging.NamedLogger("email")),

		// noop
		fx.Provide(NewNoOpDriverFactory),

		// service
		fx.Supply(cfg),
		fx.Provide(New),
	)
}
