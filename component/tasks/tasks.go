package tasks

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/x/driver"
)

type Driver interface {
	Submit(context.Context, *CreateTaskRequest) error
	Receive(context.Context, RawTask) (*Task, error)
}

type Tasks interface {
	Driver
}

type TaskParams struct {
	fx.In

	Drivers []*driver.Factory[TaskDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Manager struct {
	drivers *driver.Pool[TaskDriver, Driver]
	config  *Config
	log     *zap.Logger
}

var _ = Tasks(&Manager{})

func New(params TaskParams) Tasks {
	return &Manager{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log.Named("tasks"),
	}
}

func (q *Manager) resolveDriver() (Driver, error) {
	return q.drivers.Resolve(q.config.Driver)
}

func (q *Manager) Submit(ctx context.Context, req *CreateTaskRequest) error {
	driver, err := q.resolveDriver()
	if err != nil {
		return err
	}

	return driver.Submit(ctx, req)
}

func (q *Manager) Receive(ctx context.Context, raw RawTask) (*Task, error) {
	driver, err := q.resolveDriver()
	if err != nil {
		return nil, err
	}

	return driver.Receive(ctx, raw)
}
