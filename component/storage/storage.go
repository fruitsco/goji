package storage

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/fruitsco/goji/x/driver"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SignedUploadOptions struct {
	Size     int64
	MimeType string
	Expires  time.Duration
}

type SignedDownloadOptions struct {
	Expires time.Duration
}

type SignResult struct {
	Method  string
	URL     *url.URL
	Headers http.Header
}

type Driver interface {
	Exists(ctx context.Context, bucketName string, name string) (bool, error)
	Delete(ctx context.Context, bucketName string, name string) error
	SignedUpload(ctx context.Context, bucketName string, name string, options *SignedUploadOptions) (*SignResult, error)
	SignedDownload(ctx context.Context, bucketName string, name string) (*SignResult, error)
	SignedDownloadWithOptions(ctx context.Context, bucketName string, name string, options *SignedDownloadOptions) (*SignResult, error)
	Download(ctx context.Context, bucketName string, name string) ([]byte, error)
	Upload(ctx context.Context, bucketName string, name string, data []byte) error
	Copy(ctx context.Context, srcBucket string, srcName string, dstBucket string, dstName string) error
}

type Storage interface {
	Driver

	Driver(name StorageDriver) (Driver, error)
}

type StorageParams struct {
	fx.In

	Drivers []*driver.Factory[StorageDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Manager struct {
	drivers *driver.Pool[StorageDriver, Driver]
	config  *Config
	log     *zap.Logger
}

var _ = Storage(&Manager{})

func New(params StorageParams) Storage {
	return &Manager{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log.Named("storage"),
	}
}

func (s *Manager) Driver(name StorageDriver) (Driver, error) {
	return s.drivers.Resolve(name)
}

// MARK: - Default Driver

func (s *Manager) defaultDriver() (Driver, error) {
	return s.drivers.Resolve(s.config.Driver)
}

func (s *Manager) Exists(ctx context.Context, bucketName string, name string) (bool, error) {
	driver, err := s.defaultDriver()
	if err != nil {
		return false, err
	}

	return driver.Exists(ctx, bucketName, name)
}

func (s *Manager) Delete(ctx context.Context, bucketName string, name string) error {
	driver, err := s.defaultDriver()
	if err != nil {
		return err
	}

	return driver.Delete(ctx, bucketName, name)
}

func (s *Manager) SignedUpload(ctx context.Context, bucketName string, name string, options *SignedUploadOptions) (*SignResult, error) {
	driver, err := s.defaultDriver()
	if err != nil {
		return nil, err
	}

	return driver.SignedUpload(ctx, bucketName, name, options)
}

func (s *Manager) SignedDownload(ctx context.Context, bucketName string, name string) (*SignResult, error) {
	driver, err := s.defaultDriver()
	if err != nil {
		return nil, err
	}

	return driver.SignedDownload(ctx, bucketName, name)
}

func (s *Manager) SignedDownloadWithOptions(ctx context.Context, bucketName string, name string, options *SignedDownloadOptions) (*SignResult, error) {
	driver, err := s.defaultDriver()
	if err != nil {
		return nil, err
	}

	return driver.SignedDownloadWithOptions(ctx, bucketName, name, options)
}

func (s *Manager) Download(ctx context.Context, bucketName string, name string) ([]byte, error) {
	driver, err := s.defaultDriver()
	if err != nil {
		return nil, err
	}

	return driver.Download(ctx, bucketName, name)
}

func (s *Manager) Upload(ctx context.Context, bucketName string, name string, data []byte) error {
	driver, err := s.defaultDriver()
	if err != nil {
		return err
	}

	return driver.Upload(ctx, bucketName, name, data)
}

func (s *Manager) Copy(ctx context.Context, srcBucket string, srcName string, dstBucket string, dstName string) error {
	driver, err := s.defaultDriver()
	if err != nil {
		return err
	}

	return driver.Copy(ctx, srcBucket, srcName, dstBucket, dstName)
}
