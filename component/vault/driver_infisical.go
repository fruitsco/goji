package vault

import (
	"context"
	"fmt"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/x/driver"
	infisical "github.com/infisical/go-sdk"
	infisicalModels "github.com/infisical/go-sdk/packages/models"
)

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

// InfisicalDriver is the vault driver for Infisical
type InfisicalDriver struct {
	config *InfisicalConfig
	client infisical.InfisicalClientInterface
	log    *zap.Logger
}

// InfisicalDriverParams is the parameters for the Infisical driver
type InfisicalDriverParams struct {
	fx.In

	// Config is the configuration for the Infisical driver
	Config *InfisicalConfig

	// Log is the logger for the Infisical driver
	Log *zap.Logger
}

// NewInfisicalDriverFactory creates a new Infisical driver factory
func NewInfisicalDriverFactory(
	params InfisicalDriverParams,
	lc fx.Lifecycle,
) driver.FactoryResult[DriverName, Driver] {
	return driver.NewFactory(Infisical, func() (Driver, error) {
		return NewInfisicalDriver(params, lc)
	})
}

// NewInfisicalDriver creates a new Infisical driver
func NewInfisicalDriver(
	params InfisicalDriverParams,
	lc fx.Lifecycle,
) (Driver, error) {
	if params.Config == nil {
		return nil, fmt.Errorf("config is required for Infisical driver")
	}

	if params.Config.SiteURL == "" {
		return nil, fmt.Errorf("site URL is required for Infisical driver")
	}

	if params.Config.ProjectID == "" {
		return nil, fmt.Errorf("project ID is required for Infisical driver")
	}

	if params.Config.Environment == "" {
		return nil, fmt.Errorf("environment is required for Infisical driver")
	}

	if params.Config.Auth.Strategy == "" {
		params.Config.Auth.Strategy = InfisicalAuthStrategyUniversal
	}

	client := infisical.NewInfisicalClient(infisical.Config{
		SiteUrl: params.Config.SiteURL,
	})

	if err := infisicalAuth(client.Auth(), params.Config.Auth); err != nil {
		return nil, err
	}

	return &InfisicalDriver{
		client: client,
		config: params.Config,
		log:    params.Log.Named("infisical"),
	}, nil
}

var _ = Driver(&InfisicalDriver{})

// CreateSecret creates a new secret in Infisical
func (d *InfisicalDriver) CreateSecret(
	ctx context.Context,
	name string,
	payload []byte,
) (Secret, error) {
	path, key := d.getSecretPathAndKey(name)

	secret, err := d.client.Secrets().Create(infisical.CreateSecretOptions{
		ProjectID:   d.config.ProjectID,
		Environment: d.config.Environment,

		SecretKey:   key,
		SecretPath:  path,
		SecretValue: string(payload),
	})
	if err != nil {
		return Secret{}, err
	}

	return d.mapSecret(secret), nil
}

// AddVersion adds a new version to a secret in Infisical
func (d *InfisicalDriver) AddVersion(
	ctx context.Context,
	name string,
	payload []byte,
) (Secret, error) {
	path, key := d.getSecretPathAndKey(name)

	secret, err := d.client.Secrets().Update(infisical.UpdateSecretOptions{
		ProjectID:   d.config.ProjectID,
		Environment: d.config.Environment,

		SecretKey:      key,
		SecretPath:     path,
		NewSecretValue: string(payload),
	})
	if err != nil {
		return Secret{}, err
	}

	return d.mapSecret(secret), nil
}

// GetVersion retrieves a specific version of a secret from Infisical
func (d *InfisicalDriver) GetVersion(
	ctx context.Context,
	name string,
	version int,
) (Secret, error) {
	path, key := d.getSecretPathAndKey(name)

	secret, err := d.client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
		ProjectID:   d.config.ProjectID,
		Environment: d.config.Environment,

		SecretKey:  key,
		SecretPath: path,
		// TODO: add version as soon as it's supported by the sdk
		// Version:     version,
	})
	if err != nil {
		return Secret{}, err
	}

	return d.mapSecret(secret), nil
}

// GetLatestVersion retrieves the latest version of a secret from Infisical
func (d *InfisicalDriver) GetLatestVersion(
	ctx context.Context,
	name string,
) (Secret, error) {
	path, key := d.getSecretPathAndKey(name)

	secret, err := d.client.Secrets().Retrieve(infisical.RetrieveSecretOptions{
		ProjectID:   d.config.ProjectID,
		Environment: d.config.Environment,

		SecretKey:  key,
		SecretPath: path,
	})
	if err != nil {
		return Secret{}, err
	}

	return d.mapSecret(secret), nil
}

// DeleteSecret deletes a secret from Infisical
func (d *InfisicalDriver) DeleteSecret(
	ctx context.Context,
	name string,
) error {
	path, key := d.getSecretPathAndKey(name)

	_, err := d.client.Secrets().Delete(infisical.DeleteSecretOptions{
		ProjectID:   d.config.ProjectID,
		Environment: d.config.Environment,

		SecretKey:  key,
		SecretPath: path,
	})
	return err
}

// getSecretPathAndKey splits the secret name into path and key
func (d *InfisicalDriver) getSecretPathAndKey(name string) (string, string) {
	lastSlash := strings.LastIndex(name, "/")
	if lastSlash == -1 {
		return "", name
	}

	return name[:lastSlash], name[lastSlash+1:]
}

// mapSecret maps an Infisical secret to a vault secret
func (d *InfisicalDriver) mapSecret(secret infisicalModels.Secret) Secret {
	return Secret{
		Name:    fmt.Sprintf("%s/%s", secret.SecretPath, secret.SecretKey),
		Version: secret.Version,
		Payload: []byte(secret.SecretValue),
	}
}

// infisicalAuth authenticates the infisical client based on the configuration
func infisicalAuth(
	client infisical.AuthInterface,
	config InfisicalAuthConfig,
) error {
	switch config.Strategy {
	case InfisicalAuthStrategyUniversal:
		if config.Universal == nil {
			return fmt.Errorf("universal auth strategy requires configuration")
		}

		_, err := client.UniversalAuthLogin(
			config.Universal.ClientID,
			config.Universal.ClientSecret,
		)
		return err
	case InfisicalAuthStrategyGCPIam:
		if config.GCP == nil {
			return fmt.Errorf("gcp iam auth strategy requires configuration")
		}

		_, err := client.GcpIamAuthLogin(
			config.GCP.IdentityID,
			config.GCP.ServiceAccountKeyFilePath,
		)
		return err
	case InfisicalAuthStrategyGCPIdToken:
		if config.GCP == nil {
			return fmt.Errorf("gcp id token auth strategy requires configuration")
		}

		_, err := client.GcpIdTokenAuthLogin(
			config.GCP.IdentityID,
		)
		return err
	}

	return nil
}
