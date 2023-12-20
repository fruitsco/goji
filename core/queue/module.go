package queue

import (
	"github.com/fruitsco/goji/x"
	"go.uber.org/fx"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("queue",
		fx.Decorate(x.NamedLogger("queue")),

		fx.Provide(NewNoOpDriverFactory),

		fx.Supply(cfg.PubSub),
		fx.Provide(NewPubSubDriverFactory),

		fx.Supply(cfg),
		fx.Provide(New),
	)
}
