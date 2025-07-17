package tasksgcp

import (
	"go.uber.org/fx"

	"github.com/fruitsco/goji/component/tasks"
)

func Module() fx.Option {
	return fx.Options(
		fx.Provide(func(cfg *tasks.Config) *tasks.CloudTasksConfig {
			return cfg.CloudTasks
		}),
		fx.Provide(NewCloudTasksDriverFactory),
	)
}
