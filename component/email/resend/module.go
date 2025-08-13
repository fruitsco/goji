package emailresend

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/component/email"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg *email.Config) *email.ResendConfig {
			return cfg.Resend
		}),
		fx.Provide(NewResend),
		fx.Provide(NewResendDriverFactory),
	)
}
