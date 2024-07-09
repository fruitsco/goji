package vault

import "github.com/fruitsco/goji/x/conf"

type DriverName string

const (
	// GCPSecretManager is a driver for Google Cloud Secret Manager
	GCPSecretManager DriverName = "gcp_secret_manager"

	// InfisicalVault is a driver for Infisical Vault
	Infisical DriverName = "infisical"

	// Vault is a driver for HashiCorp Vault
	HCPVault DriverName = "vault"
)

type Config struct {
	// Driver is the driver to use for the vault
	Driver DriverName `conf:"driver"`

	// GCP is the configuration for Google Cloud Platform
	GCPSecretManager *GCPSecretManagerConfig `conf:"gcp_secret_manager"`

	// Infisical is the configuration for Infisical Vault
	Infisical *InfisicalConfig `conf:"infisical"`

	// HCPVault is the configuration for HashiCorp Vault
	HCPVault *HCPVaultConfig `conf:"vault"`
}

// DefaultConfig is the default configuration for the vault
var DefaultConfig = conf.DefaultConfig{
	"driver": Infisical,
}
