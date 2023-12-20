package validation

import (
	"github.com/fruitsco/goji/x"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("validation",
		fx.Decorate(x.NamedLogger("validation")),
		fx.Provide(New),
	)
}
