package queue

import (
	"context"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Driver interface {
	Publish(ctx context.Context, message Message) error
	RecievePush(ctx context.Context, topic string, data []byte) (Message, error)
}

type Queue interface {
	Driver
}

type QueueParams struct {
	fx.In

	Drivers []*driver.Factory[QueueDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Manager struct {
	drivers *driver.Pool[QueueDriver, Driver]
	config  *Config
	log     *zap.Logger
}

var _ = Queue(&Manager{})

func New(params QueueParams) Queue {
	return &Manager{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log.Named("queue"),
	}
}

func (q *Manager) resolveDriver() (Driver, error) {
	return q.drivers.Resolve(q.config.Driver)
}

func (q *Manager) Publish(ctx context.Context, message Message) error {
	driver, err := q.resolveDriver()

	if err != nil {
		return err
	}

	return driver.Publish(ctx, message)
}

func (q *Manager) RecievePush(ctx context.Context, topic string, data []byte) (Message, error) {
	driver, err := q.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.RecievePush(ctx, topic, data)
}
