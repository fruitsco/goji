package storage

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("storage",
		fx.Decorate(logging.NamedLogger("storage")),

		// noop driver
		fx.Provide(NewNoOpDriverFactory),

		// gcs driver
		fx.Supply(cfg.GCS),
		fx.Provide(NewGCSDriverFactory),

		// minio driver
		fx.Supply(cfg.Minio),
		fx.Provide(NewMinioDriverFactory),

		// base
		fx.Supply(cfg),
		fx.Provide(New),
	)
}
