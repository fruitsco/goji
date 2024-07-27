package queue

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("queue",
		fx.Decorate(logging.NamedLogger("queue")),

		fx.Provide(NewNoOpDriverFactory),

		fx.Supply(cfg.PubSub),
		fx.Provide(NewPubSubDriverFactory),

		fx.Supply(cfg),
		fx.Provide(New),
	)
}
