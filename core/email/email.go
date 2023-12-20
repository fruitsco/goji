package email

import (
	"context"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Driver interface {
	Send(ctx context.Context, msg Message) error
}

type Email struct {
	drivers *driver.Pool[MailDriver, Driver]
	config  *Config
	log     *zap.Logger
}

type EmailParams struct {
	fx.In

	Drivers []*driver.Factory[MailDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

func New(params EmailParams) *Email {
	return &Email{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log,
	}
}

func (q *Email) resolveDriver() (Driver, error) {
	return q.drivers.Resolve(q.config.Driver)
}

func (q *Email) Send(ctx context.Context, message Message) error {
	driver, err := q.resolveDriver()

	if err != nil {
		return err
	}

	return driver.Send(ctx, message)
}
