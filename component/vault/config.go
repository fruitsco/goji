package vault

import (
	"github.com/fruitsco/goji/component/redis"
	"github.com/fruitsco/goji/conf"
)

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
	Redis *RedisConfig `conf:"redis"`
}

// DefaultConfig is the default configuration for the vault
var DefaultConfig = conf.DefaultConfig{
	"vault.driver":                "redis",
	"vault.redis.connection_name": "default",
}

// MARK: - GCP

// GCPSecretManagerConfig is the configuration for Google Cloud Secret Manager
type GCPSecretManagerConfig struct {
	// ProjectID is the project ID for Google Cloud Secret Manager
	ProjectID string `conf:"project_id"`
}

// MARK: - Infisical

// InfisicalAuthStrategy is the strategy for authentication
type InfisicalAuthStrategy string

const (
	// InfisicalAuthStrategyUniversal is the universal auth strategy
	InfisicalAuthStrategyUniversal InfisicalAuthStrategy = "universal"

	// InfisicalAuthStrategyGCP is the GCP auth strategy
	InfisicalAuthStrategyGCPIam InfisicalAuthStrategy = "gcp_iam"

	// InfisicalAuthStrategyGCPIdToken is the GCP token auth strategy
	InfisicalAuthStrategyGCPIdToken InfisicalAuthStrategy = "gcp_id_token"
)

// InfisicalUniversalAuthConfig is the configuration for the universal auth strategy
type InfisicalUniversalAuthConfig struct {
	// ClientID is the client ID for the universal auth strategy
	ClientID string `conf:"client_id"`

	// ClientSecret is the client secret for the universal auth strategy
	ClientSecret string `conf:"client_secret"`
}

// InfisicalGCPAuthConfig is the configuration for the GCP auth strategy
type InfisicalGCPAuthConfig struct {
	// IdentityID is the identity ID for the GCP auth strategy
	IdentityID string `conf:"identity_id"`

	// ServiceAccountKeyFilePath is the path to the service account key file
	ServiceAccountKeyFilePath string `conf:"service_account_key_file_path"`
}

// InfisicalAuthConfig is the configuration for the Infisical auth strategy
type InfisicalAuthConfig struct {
	// Strategy is the strategy to use for authentication
	Strategy InfisicalAuthStrategy `conf:"strategy"`

	// Universal is the configuration for the universal auth strategy
	Universal *InfisicalUniversalAuthConfig `conf:"universal"`

	// GCP is the configuration for the GCP auth strategies
	GCP *InfisicalGCPAuthConfig `conf:"gcp"`
}

// InfisicalConfig is the configuration for the Infisical driver
type InfisicalConfig struct {
	// SiteURL is the URL of the Infisical site
	SiteURL string `conf:"site_url"`

	// ProjectID is the project ID to use for the Infisical driver
	ProjectID string `conf:"project_id"`

	// Environment is the environment to use for the Infisical driver
	Environment string `conf:"environment"`

	// Auth is the configuration for authentication
	Auth InfisicalAuthConfig `conf:"auth"`
}

// MARK: - HCP Vault

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

// MARK: - Redis

// RedisConfig is the configuration for the Redis driver
type RedisConfig struct {
	// ConnectionName is the name of the Redis connection to use
	ConnectionName redis.ConnectionName `conf:"connection_name"`

	// EncryptionKey is the key to use for encryption
	EncryptionKey string `conf:"encryption_key"`
}
