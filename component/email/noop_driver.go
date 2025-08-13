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

func NewNoOpDriverFactory(params NoOpDriverParams) driver.FactoryResult[MailDriver, ConnectionFactory] {
	return NewConnectionFactory(NoOp, func(ConnectionConfig) (Driver, error) {
		return NewNoOpDriver(params), nil
	})
}

func NewNoOpDriver(params NoOpDriverParams) *NoOpDriver {
	return &NoOpDriver{
		log: params.Log,
	}
}

func (m *NoOpDriver) Send(ctx context.Context, message Message) error {
	m.log.With(zap.Any("message", message)).Info("sending message")
	return nil
}

func (m *NoOpDriver) SendID(ctx context.Context, message Message) (string, error) {
	m.log.With(zap.Any("message", message)).Info("sending message")
	return "", nil
}
