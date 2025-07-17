package emailsmtp

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/component/email"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg *email.Config) *email.SMTPConfig {
			return cfg.SMTP
		}),
		fx.Provide(NewSmtpDriverFactory),
	)
}
