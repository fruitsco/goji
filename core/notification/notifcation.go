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

type Notifcation struct {
	drivers *driver.Pool[NotificationDriver, Driver]
	config  *Config
	log     *zap.Logger
}

type NotifcationParams struct {
	fx.In

	Drivers []*driver.Factory[NotificationDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

func New(params NotifcationParams) *Notifcation {
	return &Notifcation{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log,
	}
}

func (q *Notifcation) resolveDriver() (Driver, error) {
	return q.drivers.Resolve(q.config.Driver)
}

func (q *Notifcation) Send(ctx context.Context, message Message) error {
	driver, err := q.resolveDriver()

	if err != nil {
		return err
	}

	return driver.Send(ctx, message)
}
