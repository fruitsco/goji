package vault

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

// GCPSecretManagerConfig is the configuration for Google Cloud Secret Manager
type GCPSecretManagerConfig struct {
	// ProjectID is the project ID for Google Cloud Secret Manager
	ProjectID string `conf:"project_id"`
}

// GCPSecretManagerDriver is the driver for Google Cloud Secret Manager
type GCPSecretManagerDriver struct {
	config *GCPSecretManagerConfig
	client *secretmanager.Client
	log    *zap.Logger
}

// GCPSecretManagerDriverParams is the parameters for the Google Cloud Secret Manager driver
type GCPSecretManagerDriverParams struct {
	fx.In

	// Context is the context for the Google Cloud Secret Manager driver
	Context context.Context

	// Config is the configuration for the Google Cloud Secret Manager driver
	Config *GCPSecretManagerConfig

	// Log is the logger for the Google Cloud Secret Manager driver
	Log *zap.Logger
}

// NewGCPSecretManagerDriverFactory creates a new Google Cloud Secret Manager driver factory
func NewGCPSecretManagerDriverFactory(
	params GCPSecretManagerDriverParams,
	lc fx.Lifecycle,
) driver.FactoryResult[DriverName, Driver] {
	return driver.NewFactory(GCPSecretManager, func() (Driver, error) {
		return NewGCPSecretManagerDriver(params, lc)
	})
}

// NewGCPSecretManagerDriver creates a new Google Cloud Secret Manager driver
func NewGCPSecretManagerDriver(
	params GCPSecretManagerDriverParams,
	lc fx.Lifecycle,
) (Driver, error) {
	if params.Config == nil || params.Config.ProjectID == "" {
		return nil, fmt.Errorf("project ID is required for Google Cloud Secret Manager")
	}

	client, err := secretmanager.NewClient(params.Context)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			return client.Close()
		},
	})

	return &GCPSecretManagerDriver{
		client: client,
		config: params.Config,
		log:    params.Log.Named("gcp_secret_manager"),
	}, nil
}

var _ = Driver(&GCPSecretManagerDriver{})

var _ = Closer(&GCPSecretManagerDriver{})

// CreateSecret creates a new secret in Google Cloud Secret Manager
func (d *GCPSecretManagerDriver) CreateSecret(
	ctx context.Context,
	name string,
	payload []byte,
) (Secret, error) {
	createSecretReq := &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", d.config.ProjectID),
		SecretId: name,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{
					Automatic: &secretmanagerpb.Replication_Automatic{},
				},
			},
			// TODO: implement rotation
		},
	}

	_, err := d.client.CreateSecret(ctx, createSecretReq)
	if err != nil {
		return Secret{}, err
	}

	return d.AddVersion(ctx, name, payload)
}

// AddVersion adds a new version to a secret in Google Cloud Secret Manager
func (d *GCPSecretManagerDriver) AddVersion(
	ctx context.Context,
	name string,
	payload []byte,
) (Secret, error) {
	addSecretVersionReq := &secretmanagerpb.AddSecretVersionRequest{
		Parent: fmt.Sprintf("projects/%s/secrets/%s", d.config.ProjectID, name),
		Payload: &secretmanagerpb.SecretPayload{
			Data: payload,
		},
	}

	addSecretVersionResp, err := d.client.AddSecretVersion(ctx, addSecretVersionReq)
	if err != nil {
		return Secret{}, err
	}

	version, err := d.getVersionFromName(addSecretVersionResp.Name)
	if err != nil {
		return Secret{}, err
	}

	return Secret{
		Name:    name,
		Version: version,
		Payload: payload,
	}, nil
}

// GetVersion retrieves a specific version of a secret from Google Cloud Secret Manager
func (d *GCPSecretManagerDriver) GetVersion(
	ctx context.Context,
	name string,
	version int,
) (Secret, error) {
	return d.getSecretVersion(ctx, name, fmt.Sprintf("%d", version))
}

// GetLatestVersion retrieves the latest version of a secret from Google Cloud Secret Manager
func (d *GCPSecretManagerDriver) GetLatestVersion(
	ctx context.Context,
	name string,
) (Secret, error) {
	return d.getSecretVersion(ctx, name, "latest")
}

// DeleteSecret deletes a secret from Google Cloud Secret Manager
func (d *GCPSecretManagerDriver) DeleteSecret(ctx context.Context, name string) error {
	deleteSecretReq := &secretmanagerpb.DeleteSecretRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s", d.config.ProjectID, name),
	}

	return d.client.DeleteSecret(ctx, deleteSecretReq)
}

// Close closes the Google Cloud Secret Manager driver
func (d *GCPSecretManagerDriver) Close() error {
	return d.client.Close()
}

// getSecretVersion retrieves a specific version of a secret from Google Cloud Secret Manager
func (d *GCPSecretManagerDriver) getSecretVersion(
	ctx context.Context,
	name string,
	version string,
) (Secret, error) {
	accessSecretVersionReq := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", d.config.ProjectID, name, version),
	}

	accessSecretVersionResp, err := d.client.AccessSecretVersion(ctx, accessSecretVersionReq)
	if err != nil {
		return Secret{}, err
	}

	versionParsed, err := d.getVersionFromName(accessSecretVersionResp.Name)
	if err != nil {
		return Secret{}, err
	}

	return Secret{
		Name:    name,
		Version: versionParsed,
		Payload: accessSecretVersionResp.Payload.Data,
	}, nil
}

// getVersionFromName extracts the version from a secret version name
func (d *GCPSecretManagerDriver) getVersionFromName(name string) (int, error) {
	nameParts := strings.Split(name, "/")
	versionStr := nameParts[len(nameParts)-1]
	versionParsed, err := strconv.Atoi(versionStr)
	if err != nil {
		return 0, err
	}

	return versionParsed, nil
}
