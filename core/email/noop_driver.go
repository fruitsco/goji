package email

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

func NewNoOpDriverFactory(params NoOpDriverParams) driver.FactoryResult[MailDriver, Driver] {
	return driver.NewFactory(NoOp, func() (Driver, error) {
		return NewNoOpDriver(params), nil
	})
}

func NewNoOpDriver(params NoOpDriverParams) *NoOpDriver {
	return &NoOpDriver{
		log: params.Log,
	}
}

// Send a message using postino
func (m *NoOpDriver) Send(ctx context.Context, message Message) error {
	m.log.With(zap.Any("message", message)).Info("sending message")
	return nil
}
