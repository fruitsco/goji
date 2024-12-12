package tasks

import (
	"context"
	"errors"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/x/driver"
)

type NoOpDriver struct {
	log *zap.Logger
}

var _ = Driver(&NoOpDriver{})

type NoOpDriverParams struct {
	fx.In

	Log *zap.Logger
}

func NewNoOpDriverFactory(params NoOpDriverParams) driver.FactoryResult[TaskDriver, Driver] {
	return driver.NewFactory(NoOp, func() (Driver, error) {
		return NewNoOpDriver(params), nil
	})
}

func NewNoOpDriver(params NoOpDriverParams) *NoOpDriver {
	return &NoOpDriver{
		log: params.Log.Named("noop"),
	}
}

func (q *NoOpDriver) Name() TaskDriver {
	return NoOp
}

func (q *NoOpDriver) Submit(context.Context, *CreateTaskRequest) error {
	return errors.New("not implemented")
}

func (q *NoOpDriver) Receive(context.Context, RawTask) (*Task, error) {
	return nil, errors.New("not implemented")
}
