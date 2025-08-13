package email

import (
	"context"
	"fmt"
	"strings"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Driver interface {
	Send(ctx context.Context, msg Message) error
	SendID(ctx context.Context, msg Message) (string, error)
}

type ConnectionFactory func(ConnectionConfig) (Driver, error)

func NewConnectionFactory(name MailDriver, f ConnectionFactory) driver.FactoryResult[MailDriver, ConnectionFactory] {
	return driver.NewFactory(name, func() (ConnectionFactory, error) {
		return f, nil
	})
}

type Email interface {
	Driver
}

type EmailParams struct {
	fx.In

	Drivers driver.Factories[MailDriver, ConnectionFactory] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Manager struct {
	drivers *driver.Pool[MailDriver, ConnectionFactory]
	config  *Config
	log     *zap.Logger
}

var _ = Email(&Manager{})

func New(params EmailParams) Email {
	if params.Config == nil {
		params.Config = &Config{}
	}

	sanitizedConnections := make(map[string]ConnectionConfig)
	for name, cfg := range params.Config.Connections {
		sanitizedConnections[strings.ToLower(name)] = cfg
	}
	params.Config.Connections = sanitizedConnections

	return &Manager{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log,
	}
}

func (q *Manager) resolveDefaultDriver() (Driver, error) {
	f, err := q.drivers.Resolve(q.config.Driver)
	if err != nil {
		return nil, fmt.Errorf("could not resolve default connection: %s", q.config.Driver)
	}

	// create artificial connection config from legacy default connection
	cfg := ConnectionConfig{
		Driver:  q.config.Driver,
		Mailgun: q.config.Mailgun,
		SMTP:    q.config.SMTP,
		Resend:  q.config.Resend,
	}

	c, err := f(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create default connection %s: %w", q.config.Driver, err)
	}

	return c, nil
}

func (q *Manager) Send(ctx context.Context, message Message) error {
	driver, err := q.resolveDefaultDriver()
	if err != nil {
		return err
	}

	return driver.Send(ctx, message)
}

func (q *Manager) SendID(ctx context.Context, message Message) (string, error) {
	driver, err := q.resolveDefaultDriver()
	if err != nil {
		return "", err
	}

	return driver.SendID(ctx, message)
}

func (q *Manager) Connection(name string) (Driver, error) {
	if name == "default" {
		return q.resolveDefaultDriver()
	}

	cfg, ok := q.config.Connections[name]
	if !ok {
		return nil, fmt.Errorf("could not find connection: %s", name)
	}

	f, err := q.drivers.Resolve(cfg.Driver)
	if err != nil {
		return nil, fmt.Errorf("could not resolve connection %s: %w", name, err)
	}

	c, err := f(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create connection %s: %w", name, err)
	}

	return c, nil
}
