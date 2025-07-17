package storagegcs

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	gcs "cloud.google.com/go/storage"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/storage"
	"github.com/fruitsco/goji/x/driver"
)

type GCSDriver struct {
	config *storage.GCSConfig
	client *gcs.Client
	log    *zap.Logger
}

var _ = storage.Driver(&GCSDriver{})

type GCSDriverParams struct {
	fx.In

	Context context.Context
	Config  *storage.GCSConfig
	Log     *zap.Logger
}

func NewGCSDriverFactory(params GCSDriverParams) driver.FactoryResult[storage.StorageDriver, storage.Driver] {
	return driver.NewFactory(storage.GCS, func() (storage.Driver, error) {
		return NewGCSDriver(params)
	})
}

// NewGCSDriver creates a new storage base struct
func NewGCSDriver(params GCSDriverParams) (*GCSDriver, error) {
	client, err := gcs.NewClient(params.Context)
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
		if errors.Is(err, gcs.ErrObjectNotExist) {
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

	if err := bucket.Object(name).Delete(ctx); err != nil {
		return err
	}

	return nil
}

// GetPresignedUpload returns a presigned url for uploading a file
func (s *GCSDriver) SignedUpload(
	context context.Context,
	bucketName string,
	name string,
	options *storage.SignedUploadOptions,
) (*storage.SignResult, error) {
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

	opts := &gcs.SignedURLOptions{
		Scheme:      gcs.SigningSchemeV4,
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

	return &storage.SignResult{
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
) (*storage.SignResult, error) {
	bucket := s.client.Bucket(bucketName)

	expires := time.Duration(s.config.Expiration) * time.Second

	opts := &gcs.SignedURLOptions{
		Scheme:  gcs.SigningSchemeV4,
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

	return &storage.SignResult{
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

	if _, err := w.Write(data); err != nil {
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

func (s *GCSDriver) Copy(
	ctx context.Context,
	srcBucketName string,
	srcName string,
	dstBucketName string,
	dstName string,
) error {
	srcBucketObj := s.client.Bucket(srcBucketName)
	dstBucketObj := s.client.Bucket(dstBucketName)

	srcObj := srcBucketObj.Object(srcName)
	dstObj := dstBucketObj.Object(dstName)

	if _, err := dstObj.CopierFrom(srcObj).Run(ctx); err != nil {
		return err
	}

	return nil
}
