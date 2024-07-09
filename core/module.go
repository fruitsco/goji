package core

import (
	"github.com/fruitsco/goji/core/crypt"
	"github.com/fruitsco/goji/core/database"
	"github.com/fruitsco/goji/core/email"
	"github.com/fruitsco/goji/core/notification"
	"github.com/fruitsco/goji/core/queue"
	"github.com/fruitsco/goji/core/redis"
	"github.com/fruitsco/goji/core/storage"
	"github.com/fruitsco/goji/core/validation"
	"github.com/fruitsco/goji/core/vault"
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
	)
}
