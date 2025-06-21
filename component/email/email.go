package email

import (
	"context"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Driver interface {
	Send(ctx context.Context, msg Message) error
	SendID(ctx context.Context, msg Message) (string, error)
}

type Email interface {
	Driver
}

type EmailParams struct {
	fx.In

	Drivers []*driver.Factory[MailDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Manager struct {
	drivers *driver.Pool[MailDriver, Driver]
	config  *Config
	log     *zap.Logger
}

var _ = Email(&Manager{})

func New(params EmailParams) Email {
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

func (q *Manager) SendID(ctx context.Context, message Message) (string, error) {
	driver, err := q.resolveDriver()
	if err != nil {
		return "", err
	}

	return driver.SendID(ctx, message)
}
