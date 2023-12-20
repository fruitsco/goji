package storage

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"net/http/httptrace"

	"github.com/fruitsco/goji/x/driver"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type MinioDriver struct {
	config *MinioConfig
	client *minio.Client
	log    *zap.Logger
}

var _ = Driver(&MinioDriver{})

type MinioDriverParams struct {
	fx.In

	Config *MinioConfig
	Log    *zap.Logger
}

func NewMinioDriverFactory(params MinioDriverParams) driver.FactoryResult[StorageDriver, Driver] {
	return driver.NewFactory(Minio, func() (Driver, error) {
		return NewMinioDriver(params)
	})
}

// NewMinioDriver creates a new storage base struct
func NewMinioDriver(params MinioDriverParams) (*MinioDriver, error) {
	log := params.Log.Named("minio")

	transport, err := minio.DefaultTransport(params.Config.Secure)
	if err != nil {
		return nil, err
	}

	// If we are using a secure connection w/ a self-signed certificate,
	// we need to skip verification.
	if params.Config.Secure && params.Config.TLSSkipVerify {
		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	// If the proxy url is set, set the proxy env variables for the minio client.
	// We use minio's default transport and set the proxy url on it.
	if params.Config.ProxyUrl != nil && *params.Config.ProxyUrl != "" {
		proxyUrl, err := url.Parse(*params.Config.ProxyUrl)
		if err != nil {
			return nil, err
		}
		transport.Proxy = http.ProxyURL(proxyUrl)
	}

	var trace *httptrace.ClientTrace
	if params.Config.Trace {
		trace = createClientTrace(log.Named("trace"))
	}

	client, err := minio.New(params.Config.Host, &minio.Options{
		Creds:     credentials.NewStaticV4(params.Config.AccessKey, params.Config.SecretKey, ""),
		Secure:    params.Config.Secure,
		Region:    params.Config.Region,
		Transport: transport,
		Trace:     trace,
	})

	if err != nil {
		return nil, err
	}

	return &MinioDriver{
		client: client,
		config: params.Config,
		log:    log,
	}, nil
}

func (s *MinioDriver) Exists(ctx context.Context, bucketName string, name string) (bool, error) {
	_, err := s.client.StatObject(ctx, bucketName, name, minio.StatObjectOptions{})

	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *MinioDriver) Delete(ctx context.Context, bucketName string, name string) error {
	return s.client.RemoveObject(ctx, bucketName, name, minio.RemoveObjectOptions{})
}

func (s *MinioDriver) SignedUpload(
	ctx context.Context,
	bucketName string,
	name string,
	options *SignedUploadOptions,
) (*SignResult, error) {
	// Set request parameters for content-type and content-length-range
	headers := http.Header{
		"Content-Type": []string{options.MimeType},
		// ISSUE: content length??
	}

	expires := time.Duration(s.config.Expiration) * time.Second

	url, err := s.client.PresignHeader(ctx, "PUT", bucketName, name, expires, nil, headers)

	if err != nil {
		log.Printf("error presigning put request: %v", err)
		return nil, err
	}

	return &SignResult{
		URL:     url,
		Method:  "PUT",
		Headers: headers,
	}, nil
}

func (s *MinioDriver) SignedDownload(
	ctx context.Context,
	bucketName string,
	name string,
) (*SignResult, error) {
	expires := time.Duration(s.config.Expiration) * time.Second

	url, err := s.client.PresignedGetObject(ctx, bucketName, name, expires, nil)

	if err != nil {
		return nil, err
	}

	return &SignResult{
		URL:    url,
		Method: "GET",
	}, nil
}

func (s *MinioDriver) Download(ctx context.Context, bucketName string, name string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Upload uploads a file to the storage
func (s *MinioDriver) Upload(ctx context.Context, bucketName string, name string, data []byte) error {
	_, err := s.client.PutObject(ctx, bucketName, name, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{})

	return err
}

func createClientTrace(log *zap.Logger) *httptrace.ClientTrace {
	return &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			log.Debug("starting to create conn", zap.String("hostPort", hostPort))
		},
		DNSStart: func(info httptrace.DNSStartInfo) {
			log.Debug("starting to look up dns", zap.String("host", info.Host))
		},
		DNSDone: func(info httptrace.DNSDoneInfo) {
			log.Debug("done looking up dns", zap.Any("info", info))
		},
		ConnectStart: func(network, addr string) {
			log.Debug("starting tcp connection", zap.String("network", network), zap.String("addr", addr))
		},
		ConnectDone: func(network, addr string, err error) {
			log.Debug("tcp connection created", zap.String("network", network), zap.String("addr", addr), zap.Error(err))
		},
		GotConn: func(info httptrace.GotConnInfo) {
			log.Debug("connection established", zap.Bool("reused", info.Reused), zap.Bool("wasIdle", info.WasIdle), zap.Duration("idleTime", info.IdleTime))
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			log.Debug("wrote request", zap.Error(info.Err))
		},
		WroteHeaderField: func(key string, value []string) {
			log.Debug("wrote header field", zap.String("key", key), zap.Strings("value", value))
		},
	}
}
