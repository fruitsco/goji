package notification

import (
	"context"

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

func NewNoOpDriverFactory(params NoOpDriverParams) driver.FactoryResult[NotificationDriver, Driver] {
	return driver.NewFactory(NoOp, func() (Driver, error) {
		return NewNoOpDriver(params), nil
	})
}

// NewPostinoDriver returns a new queue mailer implementation
func NewNoOpDriver(params NoOpDriverParams) *NoOpDriver {
	return &NoOpDriver{
		log: params.Log,
	}
}

func (m *NoOpDriver) Send(ctx context.Context, message Message) error {
	return nil
}
