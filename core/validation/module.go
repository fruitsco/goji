package validation

import (
	"github.com/fruitsco/goji/x/logging"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("validation",
		fx.Decorate(logging.NamedLogger("validation")),
		fx.Provide(New),
	)
}
