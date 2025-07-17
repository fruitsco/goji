package vaulthcp

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
	gcpAuth "github.com/hashicorp/vault/api/auth/gcp"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/vault"
	"github.com/fruitsco/goji/x/driver"
)

// HCPVaultDriver is the HashiCorp Vault driver
type HCPVaultDriver struct {
	config *vault.HCPVaultConfig
	client *vaultapi.Client
	log    *zap.Logger
}

// HCPVaultDriverParams is the parameters for the HashiCorp Vault driver
type HCPVaultDriverParams struct {
	fx.In

	// Config is the configuration for the HashiCorp Vault driver
	Config *vault.HCPVaultConfig

	// Log is the logger for the HashiCorp Vault driver
	Log *zap.Logger
}

// NewHCPVaultDriverFactory creates a new HashiCorp Vault driver factory
func NewHCPVaultDriverFactory(
	params HCPVaultDriverParams,
	lc fx.Lifecycle,
) driver.FactoryResult[vault.DriverName, vault.Driver] {
	return driver.NewFactory(vault.HCPVault, func() (vault.Driver, error) {
		return NewHCPVaultDriver(params, lc)
	})
}

// NewHCPVaultDriver creates a new HashiCorp Vault driver
func NewHCPVaultDriver(
	params HCPVaultDriverParams,
	lc fx.Lifecycle,
) (vault.Driver, error) {
	if params.Config == nil {
		return nil, fmt.Errorf("config is required for HashiCorp Vault driver")
	}

	if params.Config.Address == "" {
		return nil, fmt.Errorf("address is required for HashiCorp Vault driver")
	}

	if params.Config.Auth.Strategy == "" {
		params.Config.Auth.Strategy = vault.HCPVaultAuthStrategyToken
	}

	// TODO: advanced config
	config := vaultapi.DefaultConfig()
	config.Address = params.Config.Address

	client, err := vaultapi.NewClient(config)
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

var _ = vault.Driver(&HCPVaultDriver{})

func (d *HCPVaultDriver) kv() *vaultapi.KVv2 {
	return d.client.KVv2(d.config.MountPath)
}

// CreateSecret creates a new secret in HashiCorp Vault
func (d *HCPVaultDriver) CreateSecret(
	ctx context.Context,
	name string,
	payload []byte,
) (vault.Secret, error) {
	secret, err := d.kv().Put(ctx, name, map[string]interface{}{
		"data": payload,
	})
	if err != nil {
		return vault.Secret{}, err
	}

	return d.mapSecret(name, secret)
}

// AddVersion adds a new version to a secret in HashiCorp Vault
func (d *HCPVaultDriver) AddVersion(
	ctx context.Context,
	name string,
	payload []byte,
) (vault.Secret, error) {
	secret, err := d.kv().Put(ctx, name, map[string]interface{}{
		"data": payload,
	})
	if err != nil {
		return vault.Secret{}, err
	}

	return d.mapSecret(name, secret)
}

// GetVersion gets a specific version of a secret from HashiCorp Vault
func (d *HCPVaultDriver) GetVersion(
	ctx context.Context,
	name string,
	version int,
) (vault.Secret, error) {
	secret, err := d.kv().GetVersion(ctx, name, version)
	if err != nil {
		return vault.Secret{}, err
	}

	return d.mapSecret(name, secret)
}

// GetLatestVersion gets the latest version of a secret from HashiCorp Vault
func (d *HCPVaultDriver) GetLatestVersion(
	ctx context.Context,
	name string,
) (vault.Secret, error) {
	secret, err := d.kv().Get(ctx, name)
	if err != nil {
		return vault.Secret{}, err
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

func (d *HCPVaultDriver) mapSecret(name string, secret *vaultapi.KVSecret) (vault.Secret, error) {
	payload, ok := secret.Data["data"].([]byte)
	if !ok {
		return vault.Secret{}, fmt.Errorf("unable to parse secret data")
	}

	return vault.Secret{
		Name:    name,
		Version: secret.VersionMetadata.Version,
		Payload: payload,
	}, nil
}

// hcpAuth authenticates to HashiCorp Vault
func hcpAuth(
	ctx context.Context,
	client *vaultapi.Client,
	config vault.HCPVaultAuthConfig,
) error {
	switch config.Strategy {
	case vault.HCPVaultAuthStrategyToken:
		if config.Token == nil {
			return fmt.Errorf("token auth strategy requires token configuration")
		}

		client.SetToken(config.Token.Token)
	case vault.HCPVaultAuthStrategyGCP:
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
