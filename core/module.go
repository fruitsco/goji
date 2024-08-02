package core

import (
	"github.com/fruitsco/goji/component/crypt"
	"github.com/fruitsco/goji/component/database"
	"github.com/fruitsco/goji/component/email"
	"github.com/fruitsco/goji/component/notification"
	"github.com/fruitsco/goji/component/queue"
	"github.com/fruitsco/goji/component/redis"
	"github.com/fruitsco/goji/component/storage"
	"github.com/fruitsco/goji/component/tasks"
	"github.com/fruitsco/goji/component/validation"
	"github.com/fruitsco/goji/component/vault"
	"go.uber.org/fx"
)

func Module(config *Config) fx.Option {
	return fx.Module("core",
		database.Module(config.Database),
		email.Module(config.Email),
		notification.Module(config.Notification),
		queue.Module(config.Queue),
		redis.Module(config.Redis),
		storage.Module(config.Storage),
		validation.Module(),
		vault.Module(config.Vault),
		crypt.Module(config.Crypt),
		tasks.Module(config.Tasks),
	)
}
