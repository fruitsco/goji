package notification

import (
	"github.com/fruitsco/goji/x"
	"go.uber.org/fx"
)

func Module(config *Config) fx.Option {
	return fx.Module("notifications",
		fx.Decorate(x.NamedLogger("notifications")),

		// slack driver
		fx.Supply(config.Slack),
		fx.Provide(NewSlackDriverFactory),

		// base
		fx.Supply(config),
		fx.Provide(New),
	)
}
