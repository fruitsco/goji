package vault

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
	gcpAuth "github.com/hashicorp/vault/api/auth/gcp"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/x/driver"
)

// HCPVaultAuthStrategy is the strategy for authentication
type HCPVaultAuthStrategy string

const (
	// HCPVaultAuthStrategyToken is the token auth strategy
	HCPVaultAuthStrategyToken HCPVaultAuthStrategy = "token"

	// HCPVaultAuthStrategyGCP is the GCP auth strategy
	HCPVaultAuthStrategyGCP HCPVaultAuthStrategy = "gcp"
)

// HCPVaultTokenAuthConfig is the configuration for the token auth strategy
type HCPVaultTokenAuthConfig struct {
	// Token is the token for the token auth strategy
	Token string `conf:"token"`
}

// HCPVaultGCPAuthConfig is the configuration for the GCP auth strategy
type HCPVaultGCPAuthConfig struct {
	// RoleName is the role name for the GCP auth strategy
	RoleName string `conf:"role_name"`

	// ServiceAccountEmail is the service account email for the GCP auth strategy
	ServiceAccountEmail string `conf:"service_account_email"`
}

// HCPVaultAuthConfig is the configuration for the HashiCorp Vault auth strategy
type HCPVaultAuthConfig struct {
	// Strategy is the strategy to use for authentication
	Strategy HCPVaultAuthStrategy `conf:"strategy"`

	// Token is the configuration for the token auth strategy
	Token *HCPVaultTokenAuthConfig `conf:"token"`

	// GCP is the configuration for the GCP auth strategy
	GCP *HCPVaultGCPAuthConfig `conf:"gcp"`
}

// HCPVaultConfig is the configuration for the HashiCorp Vault driver
type HCPVaultConfig struct {
	// Address is the address of the HashiCorp Vault server
	Address string `conf:"address"`

	// MountPath is the mount path for the HashiCorp Vault server
	MountPath string `conf:"mount_path"`

	// Auth is the configuration for the HashiCorp Vault auth strategy
	Auth HCPVaultAuthConfig `conf:"auth"`
}

// HCPVaultDriver is the HashiCorp Vault driver
type HCPVaultDriver struct {
	config *HCPVaultConfig
	client *vault.Client
	log    *zap.Logger
}

// HCPVaultDriverParams is the parameters for the HashiCorp Vault driver
type HCPVaultDriverParams struct {
	fx.In

	// Config is the configuration for the HashiCorp Vault driver
	Config *HCPVaultConfig

	// Log is the logger for the HashiCorp Vault driver
	Log *zap.Logger
}

// NewHCPVaultDriverFactory creates a new HashiCorp Vault driver factory
func NewHCPVaultDriverFactory(
	params HCPVaultDriverParams,
	lc fx.Lifecycle,
) driver.FactoryResult[DriverName, Driver] {
	return driver.NewFactory(HCPVault, func() (Driver, error) {
		return NewHCPVaultDriver(params, lc)
	})
}

// NewHCPVaultDriver creates a new HashiCorp Vault driver
func NewHCPVaultDriver(
	params HCPVaultDriverParams,
	lc fx.Lifecycle,
) (Driver, error) {
	if params.Config == nil {
		return nil, fmt.Errorf("config is required for HashiCorp Vault driver")
	}

	if params.Config.Address == "" {
		return nil, fmt.Errorf("address is required for HashiCorp Vault driver")
	}

	if params.Config.Auth.Strategy == "" {
		params.Config.Auth.Strategy = HCPVaultAuthStrategyToken
	}

	// TODO: advanced config
	config := vault.DefaultConfig()
	config.Address = params.Config.Address

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, err
	}

	if params.Config.MountPath == "" {
		params.Config.MountPath = "secret"
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return hcpAuth(ctx, client, params.Config.Auth)
		},
	})

	return &HCPVaultDriver{
		client: client,
		config: params.Config,
		log:    params.Log.Named("hcp_vault"),
	}, nil
}

var _ = Driver(&HCPVaultDriver{})

func (d *HCPVaultDriver) kv() *vault.KVv2 {
	return d.client.KVv2(d.config.MountPath)
}

// CreateSecret creates a new secret in HashiCorp Vault
func (d *HCPVaultDriver) CreateSecret(
	ctx context.Context,
	name string,
	payload []byte,
) (Secret, error) {
	secret, err := d.kv().Put(ctx, name, map[string]interface{}{
		"data": payload,
	})
	if err != nil {
		return Secret{}, err
	}

	return d.mapSecret(name, secret)
}

// AddVersion adds a new version to a secret in HashiCorp Vault
func (d *HCPVaultDriver) AddVersion(
	ctx context.Context,
	name string,
	payload []byte,
) (Secret, error) {
	secret, err := d.kv().Put(ctx, name, map[string]interface{}{
		"data": payload,
	})
	if err != nil {
		return Secret{}, err
	}

	return d.mapSecret(name, secret)
}

// GetVersion gets a specific version of a secret from HashiCorp Vault
func (d *HCPVaultDriver) GetVersion(
	ctx context.Context,
	name string,
	version int,
) (Secret, error) {
	secret, err := d.kv().GetVersion(ctx, name, version)
	if err != nil {
		return Secret{}, err
	}

	return d.mapSecret(name, secret)
}

// GetLatestVersion gets the latest version of a secret from HashiCorp Vault
func (d *HCPVaultDriver) GetLatestVersion(
	ctx context.Context,
	name string,
) (Secret, error) {
	secret, err := d.kv().Get(ctx, name)
	if err != nil {
		return Secret{}, err
	}

	return d.mapSecret(name, secret)
}

// DeleteSecret deletes a secret from HashiCorp Vault
func (d *HCPVaultDriver) DeleteSecret(
	ctx context.Context,
	name string,
) error {
	return d.kv().Delete(ctx, name)
}

func (d *HCPVaultDriver) mapSecret(name string, secret *vault.KVSecret) (Secret, error) {
	payload, ok := secret.Data["data"].([]byte)
	if !ok {
		return Secret{}, fmt.Errorf("unable to parse secret data")
	}

	return Secret{
		Name:    name,
		Version: secret.VersionMetadata.Version,
		Payload: payload,
	}, nil
}

// hcpAuth authenticates to HashiCorp Vault
func hcpAuth(
	ctx context.Context,
	client *vault.Client,
	config HCPVaultAuthConfig,
) error {
	switch config.Strategy {
	case HCPVaultAuthStrategyToken:
		if config.Token == nil {
			return fmt.Errorf("token auth strategy requires token configuration")
		}

		client.SetToken(config.Token.Token)
	case HCPVaultAuthStrategyGCP:
		if config.GCP == nil {
			return fmt.Errorf("GCP auth strategy requires GCP configuration")
		}

		loginOption := gcpAuth.WithGCEAuth()

		if config.GCP.ServiceAccountEmail != "" {
			loginOption = gcpAuth.WithIAMAuth(config.GCP.ServiceAccountEmail)
		}

		auth, err := gcpAuth.NewGCPAuth(config.GCP.RoleName, loginOption)
		if err != nil {
			return err
		}

		authInfo, err := client.Auth().Login(ctx, auth)
		if err != nil {
			return fmt.Errorf("unable to login to GCP auth method: %w", err)
		}
		if authInfo == nil {
			return fmt.Errorf("login response did not return client token")
		}
	}

	return nil
}
