package storage

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/x/driver"
)

type NoOpDriver struct {
	log *zap.Logger
}

var _ = Driver(&NoOpDriver{})

type NoOpDriverParams struct {
	fx.In

	Context context.Context
	Log     *zap.Logger
}

func NewNoOpDriverFactory(params NoOpDriverParams) driver.FactoryResult[StorageDriver, Driver] {
	return driver.NewFactory(NoOp, func() (Driver, error) {
		return NewNoOpDriver(params), nil
	})
}

// NewNoOpDriver creates a new storage base struct
func NewNoOpDriver(params NoOpDriverParams) *NoOpDriver {
	return &NoOpDriver{
		log: params.Log.Named("noop"),
	}
}

// Exists checks if an object exists in the bucket
func (s *NoOpDriver) Exists(ctx context.Context, bucketName string, name string) (bool, error) {
	return false, nil
}

// Delete deletes a file from the bucket
func (s *NoOpDriver) Delete(ctx context.Context, bucketName string, name string) error {
	return nil
}

// SignedUpload returns a presigned url for uploading a file
func (s *NoOpDriver) SignedUpload(
	context context.Context,
	bucketName string,
	name string,
	options *SignedUploadOptions,
) (*SignResult, error) {
	return nil, nil
}

// SignedDownload returns a presigned url for downloading a file
func (s *NoOpDriver) SignedDownload(
	ctx context.Context,
	bucketName string,
	name string,
) (*SignResult, error) {
	return nil, nil
}

func (s *NoOpDriver) SignedDownloadWithOptions(
	ctx context.Context,
	bucketName string,
	name string,
	options *SignedDownloadOptions,
) (*SignResult, error) {
	return nil, nil
}

func (s *NoOpDriver) Download(ctx context.Context, bucketName string, name string) ([]byte, error) {
	return nil, nil
}

func (s *NoOpDriver) Upload(ctx context.Context, bucketName string, name string, data []byte) error {
	return nil
}

func (s *NoOpDriver) Copy(ctx context.Context, srcBucket string, srcName string, dstBucket string, dstName string) error {
	return nil
}
