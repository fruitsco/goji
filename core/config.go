package core

import (
	"github.com/fruitsco/goji/component/crypt"
	"github.com/fruitsco/goji/component/database"
	"github.com/fruitsco/goji/component/email"
	"github.com/fruitsco/goji/component/queue"
	"github.com/fruitsco/goji/component/redis"
	"github.com/fruitsco/goji/component/storage"
	"github.com/fruitsco/goji/component/tasks"
	"github.com/fruitsco/goji/component/vault"
	"github.com/fruitsco/goji/util"
)

type Config struct {
	Database *database.Config `conf:"db"`
	Email    *email.Config    `conf:"email"`
	Queue    *queue.Config    `conf:"queue"`
	Redis    *redis.Config    `conf:"redis"`
	Storage  *storage.Config  `conf:"storage"`
	Vault    *vault.Config    `conf:"vault"`
	Crypt    *crypt.Config    `conf:"crypt"`
	Tasks    *tasks.Config    `conf:"tasks"`
}

var DefaultConfig = util.MergeMap(
	database.DefaultConfig,
	email.DefaultConfig,
	queue.DefaultConfig,
	redis.DefaultConfig,
	storage.DefaultConfig,
	vault.DefaultConfig,
	crypt.DefaultConfig,
	tasks.DefaultConfig,
)
