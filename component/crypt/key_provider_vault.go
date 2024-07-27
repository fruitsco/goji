package crypt

import (
	"context"

	"github.com/fruitsco/goji/component/vault"
	"go.uber.org/fx"
)

type VaultKeyProvider struct {
	vault *vault.Vault
}

type VaultKeyProviderParams struct {
	fx.In

	Vault *vault.Vault
}

func NewVaultKeyProvider(params VaultKeyProviderParams) KeyProvider {
	return &VaultKeyProvider{
		vault: params.Vault,
	}
}

var _ = KeyProvider(&VaultKeyProvider{})

func (v *VaultKeyProvider) GetKey(ctx context.Context, name string) (Key, error) {
	secret, err := v.vault.GetLatestVersion(ctx, name)
	if err != nil {
		return Key{}, err
	}

	return Key{
		Name:    secret.Name,
		Version: secret.Version,
		Data:    secret.Payload,
	}, nil
}

func (v *VaultKeyProvider) GetKeyVersion(ctx context.Context, name string, version int) (Key, error) {
	secret, err := v.vault.GetVersion(ctx, name, version)
	if err != nil {
		return Key{}, err
	}

	return Key{
		Name:    secret.Name,
		Version: secret.Version,
		Data:    secret.Payload,
	}, nil
}
