package vaultinfisical

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/component/vault"
)

func Module() fx.Option {
	return fx.Options(
		fx.Decorate(func(cfg *vault.Config) *vault.InfisicalConfig {
			return cfg.Infisical
		}),
		fx.Provide(NewInfisicalDriverFactory),
	)
}
