package email

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(config *Config) fx.Option {
	return fx.Module("email",
		fx.Decorate(logging.NamedLogger("email")),

		// mailgun
		fx.Supply(config.Mailgun),
		fx.Provide(NewMailgunDriverFactory),

		// smtp
		fx.Supply(config.Smtp),
		fx.Provide(NewSmtpDriverFactory),

		// noop
		fx.Provide(NewNoOpDriverFactory),

		// service
		fx.Supply(config),
		fx.Provide(New),
	)
}
