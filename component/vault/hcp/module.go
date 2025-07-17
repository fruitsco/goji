package vaulthcp

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/component/vault"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg *vault.Config) *vault.HCPVaultConfig {
			return cfg.HCPVault
		}),
		fx.Provide(NewHCPVaultDriverFactory),
	)
}
