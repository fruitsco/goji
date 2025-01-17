package vault

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/x/driver"
)

type Driver interface {
	CreateSecret(context.Context, string, []byte) (Secret, error)
	AddVersion(context.Context, string, []byte) (Secret, error)
	GetLatestVersion(context.Context, string) (Secret, error)
	GetVersion(context.Context, string, int) (Secret, error)
	DeleteSecret(context.Context, string) error
}

type Vault interface {
	Driver

	Close() error
	Driver(name DriverName) (Driver, error)
}

type Closer interface {
	Close() error
}

type VaultParams struct {
	fx.In

	Drivers []*driver.Factory[DriverName, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Manager struct {
	drivers *driver.Pool[DriverName, Driver]
	config  *Config
	log     *zap.Logger
}

var _ = Vault(&Manager{})

func New(params VaultParams) Vault {
	return &Manager{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log.Named("vault"),
	}
}

func (v *Manager) resolveDriver() (Driver, error) {
	return v.drivers.Resolve(v.config.Driver)
}

func (v *Manager) CreateSecret(ctx context.Context, name string, payload []byte) (Secret, error) {
	driver, err := v.resolveDriver()
	if err != nil {
		return Secret{}, err
	}

	return driver.CreateSecret(ctx, name, payload)
}

func (v *Manager) AddVersion(ctx context.Context, name string, payload []byte) (Secret, error) {
	driver, err := v.resolveDriver()
	if err != nil {
		return Secret{}, err
	}

	return driver.AddVersion(ctx, name, payload)
}

func (v *Manager) GetLatestVersion(ctx context.Context, name string) (Secret, error) {
	driver, err := v.resolveDriver()
	if err != nil {
		return Secret{}, err
	}

	return driver.GetLatestVersion(ctx, name)
}

func (v *Manager) GetVersion(ctx context.Context, name string, version int) (Secret, error) {
	driver, err := v.resolveDriver()
	if err != nil {
		return Secret{}, err
	}

	return driver.GetVersion(ctx, name, version)
}

func (v *Manager) DeleteSecret(ctx context.Context, name string) error {
	driver, err := v.resolveDriver()
	if err != nil {
		return err
	}

	return driver.DeleteSecret(ctx, name)
}

func (v *Manager) Close() error {
	driver, err := v.resolveDriver()
	if err != nil {
		return err
	}

	if closer, ok := driver.(Closer); ok {
		return closer.Close()
	}

	return nil
}

func (v *Manager) Driver(name DriverName) (Driver, error) {
	return v.drivers.Resolve(name)
}
