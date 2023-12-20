package payment

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("payment",
		fx.Decorate(logging.NamedLogger("payment")),

		// stripe driver
		fx.Supply(cfg.Stripe),
		fx.Provide(NewStripeDriverFactory),

		// base
		fx.Supply(cfg),
		fx.Provide(New),
	)
}
