package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"cloud.google.com/go/storage"
	"github.com/fruitsco/goji/x/driver"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type GCSDriver struct {
	config *GCSConfig
	client *storage.Client
	log    *zap.Logger
}

var _ = Driver(&GCSDriver{})

type GCSDriverParams struct {
	fx.In

	Context context.Context
	Config  *GCSConfig
	Log     *zap.Logger
}

func NewGCSDriverFactory(params GCSDriverParams) driver.FactoryResult[StorageDriver, Driver] {
	return driver.NewFactory(GCS, func() (Driver, error) {
		return NewGCSDriver(params)
	})
}

// NewGCSDriver creates a new storage base struct
func NewGCSDriver(params GCSDriverParams) (*GCSDriver, error) {
	client, err := storage.NewClient(params.Context)

	if err != nil {
		return nil, err
	}

	return &GCSDriver{
		config: params.Config,
		client: client,
		log:    params.Log.Named("gcs"),
	}, nil
}

// ObjectExists checks if an object exists in the bucket
func (s *GCSDriver) Exists(ctx context.Context, bucketName string, name string) (bool, error) {
	bucket := s.client.Bucket(bucketName)

	obj := bucket.Object(name)

	if _, err := obj.Attrs(ctx); err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}

// DeleteFileFromBucket deletes a file from the bucket
func (s *GCSDriver) Delete(ctx context.Context, bucketName string, name string) error {
	bucket := s.client.Bucket(bucketName)

	err := bucket.Object(name).Delete(ctx)

	if err != nil {
		return err
	}

	return nil
}

// GetPresignedUpload returns a presigned url for uploading a file
func (s *GCSDriver) SignedUpload(
	context context.Context,
	bucketName string,
	name string,
	options *SignedUploadOptions,
) (*SignResult, error) {
	bucket := s.client.Bucket(bucketName)

	expires := time.Duration(s.config.Expiration) * time.Second

	headers := http.Header{
		"X-Goog-Content-Length-Range": []string{
			fmt.Sprintf("%d,%d", 0, options.Size),
		},
	}

	reqHeaders := make([]string, len(headers))
	for k, v := range headers {
		reqHeaders = append(reqHeaders, fmt.Sprintf("%s: %s", k, v[0]))
	}

	method := http.MethodPut

	opts := &storage.SignedURLOptions{
		Scheme:      storage.SigningSchemeV4,
		Method:      method,
		ContentType: options.MimeType,
		Headers:     reqHeaders,
		Expires:     time.Now().Add(expires),
	}

	signedUrlString, err := bucket.SignedURL(name, opts)

	if err != nil {
		return nil, err
	}

	signedUrl, err := url.Parse(signedUrlString)

	if err != nil {
		return nil, err
	}

	return &SignResult{
		Method:  method,
		URL:     signedUrl,
		Headers: headers,
	}, nil
}

// PresignedDownload returns a presigned url for downloading a file
func (s *GCSDriver) SignedDownload(
	ctx context.Context,
	bucketName string,
	name string,
) (*SignResult, error) {
	bucket := s.client.Bucket(bucketName)

	expires := time.Duration(s.config.Expiration) * time.Second

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  http.MethodGet,
		Expires: time.Now().Add(expires),
	}

	signedUrlString, err := bucket.SignedURL(name, opts)

	if err != nil {
		return nil, err
	}

	signedUrl, err := url.Parse(signedUrlString)

	if err != nil {
		return nil, err
	}

	return &SignResult{
		URL: signedUrl,
	}, nil
}
func (s *GCSDriver) Download(ctx context.Context, bucketName string, name string) ([]byte, error) {
	bucket := s.client.Bucket(bucketName)
	obj := bucket.Object(name)

	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Upload uploads a file to the bucket
func (s *GCSDriver) Upload(ctx context.Context, bucketName string, name string, data []byte) error {
	bucket := s.client.Bucket(bucketName)

	obj := bucket.Object(name)

	w := obj.NewWriter(ctx)

	_, err := w.Write(data)

	if err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}
