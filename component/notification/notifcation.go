package notification

import (
	"context"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Driver interface {
	Send(ctx context.Context, msg Message) error
}

type Notification interface {
	Driver
}

type NotifcationParams struct {
	fx.In

	Drivers []*driver.Factory[NotificationDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Manager struct {
	drivers *driver.Pool[NotificationDriver, Driver]
	config  *Config
	log     *zap.Logger
}

var _ = Notification(&Manager{})

func New(params NotifcationParams) Notification {
	return &Manager{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log,
	}
}

func (q *Manager) resolveDriver() (Driver, error) {
	return q.drivers.Resolve(q.config.Driver)
}

func (q *Manager) Send(ctx context.Context, message Message) error {
	driver, err := q.resolveDriver()

	if err != nil {
		return err
	}

	return driver.Send(ctx, message)
}
