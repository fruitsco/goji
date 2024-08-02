package tasks

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/x/logging"
)

func Module(cfg *Config) fx.Option {
	return fx.Module("tasks",
		fx.Decorate(logging.NamedLogger("tasks")),

		fx.Provide(NewNoOpDriverFactory),

		fx.Provide(NewQueueDriverFactory),

		fx.Supply(cfg.CloudTasks),
		fx.Provide(NewCloudTasksDriverFactory),

		fx.Supply(cfg),
		fx.Provide(New),
	)
}
