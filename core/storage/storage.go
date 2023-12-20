package storage

import (
	"context"
	"net/http"
	"net/url"

	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SignedUploadOptions struct {
	Size     int64
	MimeType string
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
	Download(ctx context.Context, bucketName string, name string) ([]byte, error)
	Upload(ctx context.Context, bucketName string, name string, data []byte) error
}

type StorageParams struct {
	fx.In

	Drivers []*driver.Factory[StorageDriver, Driver] `group:"drivers"`
	Config  *Config
	Log     *zap.Logger
}

type Storage struct {
	drivers *driver.Pool[StorageDriver, Driver]
	config  *Config
	log     *zap.Logger
}

func New(params StorageParams) *Storage {
	return &Storage{
		drivers: driver.NewPool(params.Drivers),
		config:  params.Config,
		log:     params.Log.Named("storage"),
	}
}

func (s *Storage) resolveDriver() (Driver, error) {

	return s.drivers.Resolve(s.config.Driver)
}

func (s *Storage) Exists(ctx context.Context, bucketName string, name string) (bool, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return false, err
	}

	return driver.Exists(ctx, bucketName, name)
}

func (s *Storage) Delete(ctx context.Context, bucketName string, name string) error {
	driver, err := s.resolveDriver()

	if err != nil {
		return err
	}

	return driver.Delete(ctx, bucketName, name)
}

func (s *Storage) SignedUpload(ctx context.Context, bucketName string, name string, options *SignedUploadOptions) (*SignResult, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.SignedUpload(ctx, bucketName, name, options)
}

func (s *Storage) SignedDownload(ctx context.Context, bucketName string, name string) (*SignResult, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.SignedDownload(ctx, bucketName, name)
}
func (s *Storage) Download(ctx context.Context, bucketName string, name string) ([]byte, error) {
	driver, err := s.resolveDriver()

	if err != nil {
		return nil, err
	}

	return driver.Download(ctx, bucketName, name)
}

// Upload uploads a file to the storage
func (s *Storage) Upload(ctx context.Context, bucketName string, name string, data []byte) error {
	driver, err := s.resolveDriver()

	if err != nil {
		return err
	}

	return driver.Upload(ctx, bucketName, name, data)
}
