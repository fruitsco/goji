package queue

import (
	"context"
	"errors"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NoOpDriver struct {
	log *zap.Logger
}

var _ = Driver(&NoOpDriver{})

type NoOpDriverParams struct {
	fx.In

	Log *zap.Logger
}

func NewNoOpDriverFactory(params NoOpDriverParams) driver.FactoryResult[QueueDriver, Driver] {
	return driver.NewFactory(NoOp, func() (Driver, error) {
		return NewNoOpDriver(params), nil
	})
}

func NewNoOpDriver(params NoOpDriverParams) *NoOpDriver {
	return &NoOpDriver{
		log: params.Log.Named("noop"),
	}
}

func (q *NoOpDriver) Name() QueueDriver {
	return NoOp
}

func (q *NoOpDriver) Publish(ctx context.Context, message Message) error {
	q.log.With(zap.String("topic", message.GetTopic())).Info("publishing message")
	return nil
}
func (q *NoOpDriver) RecievePush(ctx context.Context, topic string, data []byte) (Message, error) {
	return nil, errors.New("Not implemented")
}
