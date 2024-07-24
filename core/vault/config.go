package vault

import "github.com/fruitsco/goji/conf"

type DriverName string

const (
	// GCPSecretManager is a driver for Google Cloud Secret Manager
	GCPSecretManager DriverName = "gcp_secret_manager"

	// InfisicalVault is a driver for Infisical Vault
	Infisical DriverName = "infisical"

	// Vault is a driver for HashiCorp Vault
	HCPVault DriverName = "hcp_vault"

	// Redis is a driver for Redis
	Redis DriverName = "redis"
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

	// Redis is the configuration for Redis
	Redis *RedisDriverConfig `conf:"redis"`
}

// DefaultConfig is the default configuration for the vault
var DefaultConfig = conf.DefaultConfig{
	"vault.driver":                "redis",
	"vault.redis.connection_name": "default",
}
