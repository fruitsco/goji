package core

import (
	"github.com/fruitsco/goji/core/database"
	"github.com/fruitsco/goji/core/email"
	"github.com/fruitsco/goji/core/notification"
	"github.com/fruitsco/goji/core/queue"
	"github.com/fruitsco/goji/core/redis"
	"github.com/fruitsco/goji/core/storage"
	"github.com/fruitsco/goji/x"
)

type Config struct {
	Database     *database.Config     `conf:"db"`
	Email        *email.Config        `conf:"email"`
	Notification *notification.Config `conf:"notification"`
	Queue        *queue.Config        `conf:"queue"`
	Redis        *redis.Config        `conf:"redis"`
	Storage      *storage.Config      `conf:"storage"`
}

var DefaultConfig = x.MergeMap(
	database.DefaultConfig,
	email.DefaultConfig,
	notification.DefaultConfig,
	queue.DefaultConfig,
	redis.DefaultConfig,
	storage.DefaultConfig,
)
