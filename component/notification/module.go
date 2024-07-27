package notification

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(config *Config) fx.Option {
	return fx.Module("notifications",
		fx.Decorate(logging.NamedLogger("notifications")),

		// slack driver
		fx.Supply(config.Slack),
		fx.Provide(NewSlackDriverFactory),

		// base
		fx.Supply(config),
		fx.Provide(New),
	)
}
