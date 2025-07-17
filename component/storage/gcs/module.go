package storagegcs

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/component/storage"
)

func Module() fx.Option {
	return fx.Options(
		fx.Decorate(func(cfg *storage.Config) *storage.GCSConfig {
			return cfg.GCS
		}),
		fx.Provide(NewGCSDriverFactory),
	)
}
