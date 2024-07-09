package vault

import (
	"context"

	"github.com/fruitsco/goji/core/crypt"
	"github.com/fruitsco/goji/core/vault"
	"go.uber.org/fx"
)

type VaultKeyProvider struct {
	vault *vault.Vault
}

type VaultKeyProviderParams struct {
	fx.In

	Vault *vault.Vault
}

func NewVaultKeyProvider(params VaultKeyProviderParams) crypt.KeyProvider {
	return &VaultKeyProvider{
		vault: params.Vault,
	}
}

var _ = crypt.KeyProvider(&VaultKeyProvider{})

func (v *VaultKeyProvider) GetKey(ctx context.Context, name string) (crypt.Key, error) {
	secret, err := v.vault.GetLatestVersion(ctx, name)
	if err != nil {
		return crypt.Key{}, err
	}

	return crypt.Key{
		Name:    secret.Name,
		Version: secret.Version,
		Data:    secret.Payload,
	}, nil
}

func (v *VaultKeyProvider) GetKeyVersion(ctx context.Context, name string, version int) (crypt.Key, error) {
	secret, err := v.vault.GetVersion(ctx, name, version)
	if err != nil {
		return crypt.Key{}, err
	}

	return crypt.Key{
		Name:    secret.Name,
		Version: secret.Version,
		Data:    secret.Payload,
	}, nil
}
