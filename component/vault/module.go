package vault

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("vault",
		fx.Decorate(logging.NamedLogger("vault")),

		fx.Supply(cfg.GCPSecretManager),
		fx.Provide(NewGCPSecretManagerDriverFactory),

		fx.Supply(cfg.Redis),
		fx.Provide(NewRedisDriverFactory),

		fx.Supply(cfg.Infisical),
		fx.Provide(NewInfisicalDriverFactory),

		fx.Supply(cfg.HCPVault),
		fx.Provide(NewHCPVaultDriverFactory),

		fx.Supply(cfg),
		fx.Provide(New),
	)
}