package vault

import (
	"context"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Driver interface {
	CreateSecret(context.Context, string, []byte) (Secret, error)
	AddVersion(context.Context, string, []byte) (Secret, error)
	GetLatestVersion(context.Context, string) (Secret, error)
	GetVersion(context.Context, string, int) (Secret, error)
	DeleteSecret(context.Context, string) error
}

type Closer interface {
	Close() error
}

type Vault struct {
	drivers *driver.Pool[DriverName, Driver]
	config  *Config
	log     *zap.Logger
}

type VaultParams struct {
	fx.In

	Drivers []*driver.Factory[DriverName, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

func New(params VaultParams) *Vault {
	return &Vault{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log.Named("vault"),
	}
}

func (v *Vault) resolveDriver() (Driver, error) {
	return v.drivers.Resolve(v.config.Driver)
}

func (v *Vault) CreateSecret(ctx context.Context, name string, payload []byte) (Secret, error) {
	driver, err := v.resolveDriver()
	if err != nil {
		return Secret{}, err
	}

	return driver.CreateSecret(ctx, name, payload)
}

func (v *Vault) AddVersion(ctx context.Context, name string, payload []byte) (Secret, error) {
	driver, err := v.resolveDriver()
	if err != nil {
		return Secret{}, err
	}

	return driver.AddVersion(ctx, name, payload)
}

func (v *Vault) GetLatestVersion(ctx context.Context, name string) (Secret, error) {
	driver, err := v.resolveDriver()
	if err != nil {
		return Secret{}, err
	}

	return driver.GetLatestVersion(ctx, name)
}

func (v *Vault) GetVersion(ctx context.Context, name string, version int) (Secret, error) {
	driver, err := v.resolveDriver()
	if err != nil {
		return Secret{}, err
	}

	return driver.GetVersion(ctx, name, version)
}

func (v *Vault) DeleteSecret(ctx context.Context, name string) error {
	driver, err := v.resolveDriver()
	if err != nil {
		return err
	}

	return driver.DeleteSecret(ctx, name)
}

func (v *Vault) Close() error {
	driver, err := v.resolveDriver()
	if err != nil {
		return err
	}

	if closer, ok := driver.(Closer); ok {
		return closer.Close()
	}

	return nil
}

func (v *Vault) Driver(name DriverName) (Driver, error) {
	return v.drivers.Resolve(name)
}
