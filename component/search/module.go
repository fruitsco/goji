package search

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/x/logging"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("search",
		fx.Decorate(logging.NamedLogger("search")),

		fx.Supply(cfg.Typesense),
		fx.Provide(NewTypesenseDriverFactory),

		fx.Supply(cfg),
		fx.Provide(New),
	)
}
