package emailmailgun

import (
	"github.com/fruitsco/goji/component/email"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg *email.Config) *email.MailgunConfig {
			return cfg.Mailgun
		}),
		fx.Provide(NewMailgun),
		fx.Provide(NewMailgunDriverFactory),
	)
}
