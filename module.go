package goji

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewRootModule[C any](
	ctx context.Context,
	config *RootConfig[C],
	log *zap.Logger,
	options ...fx.Option,
) fx.Option {
	return fx.Options(
		// 1. inject global execution context
		fx.Supply(fx.Annotate(ctx, fx.As(new(context.Context)))),

		// 2. inject the logger
		fx.Supply(log),

		// 3. inject the app config
		fx.Supply(config.App),

		// 4. inject the log config
		fx.Supply(config.Log),

		// 6. inject the child config
		fx.Supply(&config.Child),

		// 7. use the logger also for fx' logs
		fx.WithLogger(func() fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log.Named("fx")}
		}),

		// 8. provide user-provided run options
		fx.Options(options...),
	)
}
