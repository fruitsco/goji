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

type Queue struct {
	drivers *driver.Pool[QueueDriver, Driver]
	config  *Config
	log     *zap.Logger
}

type QueueParams struct {
	fx.In

	Drivers []*driver.Factory[QueueDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

func New(params QueueParams) *Queue {
	return &Queue{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log.Named("queue"),
	}
}

func (q *Queue) resolveDriver() (Driver, error) {
	return q.drivers.Resolve(q.config.Driver)
}

func (q *Queue) Publish(ctx context.Context, message Message) error {
	driver, err := q.resolveDriver()

	if err != nil {
		return err
	}

	return driver.Publish(ctx, message)
}

func (q *Queue) RecievePush(ctx context.Context, topic string, data []byte) (Message, error) {
	driver, err := q.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.RecievePush(ctx, topic, data)
}
