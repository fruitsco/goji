package search

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/x/driver"
)

type Driver interface {
}

type Search interface {
	Driver
}

type SearchParams struct {
	fx.In

	Drivers []*driver.Factory[SearchDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Manager struct {
	drivers *driver.Pool[SearchDriver, Driver]
	config  *Config
	log     *zap.Logger
}

var _ = Search(&Manager{})

func New(params SearchParams) Search {
	return &Manager{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log.Named("search"),
	}
}

func (q *Manager) resolveDriver() (Driver, error) {
	return q.drivers.Resolve(q.config.Driver)
}
