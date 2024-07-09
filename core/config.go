package core

import (
	"github.com/fruitsco/goji/core/crypt"
	"github.com/fruitsco/goji/core/database"
	"github.com/fruitsco/goji/core/email"
	"github.com/fruitsco/goji/core/notification"
	"github.com/fruitsco/goji/core/queue"
	"github.com/fruitsco/goji/core/redis"
	"github.com/fruitsco/goji/core/storage"
	"github.com/fruitsco/goji/core/vault"
	"github.com/fruitsco/goji/util"
)

type Config struct {
	Database     *database.Config     `conf:"db"`
	Email        *email.Config        `conf:"email"`
	Notification *notification.Config `conf:"notification"`
	Queue        *queue.Config        `conf:"queue"`
	Redis        *redis.Config        `conf:"redis"`
	Storage      *storage.Config      `conf:"storage"`
	Vault        *vault.Config        `conf:"vault"`
	Crypt        *crypt.Config        `conf:"crypt"`
}

var DefaultConfig = util.MergeMap(
	database.DefaultConfig,
	email.DefaultConfig,
	notification.DefaultConfig,
	queue.DefaultConfig,
	redis.DefaultConfig,
	storage.DefaultConfig,
	vault.DefaultConfig,
	crypt.DefaultConfig,
)
