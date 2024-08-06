package search

import (
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/x/driver"
	"github.com/typesense/typesense-go/v2/typesense"
)

type TypesenseConfig struct {
	Url    string   `conf:"url"`
	ApiKey string   `conf:"api_key"`
	Nodes  []string `conf:"nodes"`
}

type TypesenseDriver struct {
	client *typesense.Client
	log    *zap.Logger
}

var _ = Driver(&TypesenseDriver{})

type TypesenseDriverParams struct {
	fx.In

	Config *TypesenseConfig
	Log    *zap.Logger
}

func NewTypesenseDriverFactory(params TypesenseDriverParams) driver.FactoryResult[SearchDriver, Driver] {
	return driver.NewFactory(Typesense, func() (Driver, error) {
		return NewTypesenseDriver(params), nil
	})
}

func NewTypesenseDriver(params TypesenseDriverParams) *TypesenseDriver {
	options := []typesense.ClientOption{
		typesense.WithServer(params.Config.Url),
		typesense.WithAPIKey(params.Config.ApiKey),
	}

	client := typesense.NewClient(options...)

	return &TypesenseDriver{
		client: client,
		log:    params.Log.Named("noop"),
	}
}

func (q *TypesenseDriver) Name() SearchDriver {
	return Typesense
}
