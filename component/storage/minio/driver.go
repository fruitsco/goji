package storageminio

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/fruitsco/goji/component/storage"
	"github.com/fruitsco/goji/x/driver"
)

type MinioDriver struct {
	config *storage.MinioConfig
	client *minio.Client
	log    *zap.Logger
}

var _ = storage.Driver(&MinioDriver{})

type MinioDriverParams struct {
	fx.In

	Config *storage.MinioConfig
	Log    *zap.Logger
}

func NewMinioDriverFactory(params MinioDriverParams) driver.FactoryResult[storage.StorageDriver, storage.Driver] {
	return driver.NewFactory(storage.Minio, func() (storage.Driver, error) {
		return NewMinioDriver(params)
	})
}

// NewMinioDriver creates a new storage base struct
func NewMinioDriver(params MinioDriverParams) (*MinioDriver, error) {
	log := params.Log.Named("minio")

	// use secure transport if the secure option is set, and we're either not using a proxy or the proxy itself is secure
	secureTransport := params.Config.Secure && (params.Config.ProxyURL == "" || strings.HasPrefix(params.Config.ProxyURL, "https://"))

	transport, err := minio.DefaultTransport(secureTransport)
	if err != nil {
		return nil, err
	}

	// If we are using a secure transport w/ a self-signed certificate,
	// we need to skip verification.
	if secureTransport && params.Config.TLSSkipVerify {
		transport.TLSClientConfig.InsecureSkipVerify = true
	}

	// If the proxy url is set, set the proxy env variables for the minio client.
	// We use minio's default transport and set the proxy url on it.
	if params.Config.ProxyURL != "" {
		proxyUrl, err := url.Parse(params.Config.ProxyURL)
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
	options *storage.SignedUploadOptions,
) (*storage.SignResult, error) {
	// Set request parameters for content-type and content-length-range
	headers := http.Header{
		"Content-Type": []string{options.MimeType},
		// ISSUE: content length??
	}

	expires := time.Duration(s.config.Expires) * time.Second

	url, err := s.client.PresignHeader(ctx, "PUT", bucketName, name, expires, nil, headers)

	if err != nil {
		log.Printf("error presigning put request: %v", err)
		return nil, err
	}

	return &storage.SignResult{
		URL:     url,
		Method:  "PUT",
		Headers: headers,
	}, nil
}

func (s *MinioDriver) SignedDownload(
	ctx context.Context,
	bucketName string,
	name string,
) (*storage.SignResult, error) {
	return s.SignedDownloadWithOptions(ctx, bucketName, name, nil)
}

func (s *MinioDriver) SignedDownloadWithOptions(
	ctx context.Context,
	bucketName string,
	name string,
	options *storage.SignedDownloadOptions,
) (*storage.SignResult, error) {
	expires := time.Duration(s.config.Expires) * time.Second
	if options != nil && options.Expires > 0 {
		expires = options.Expires
	}

	url, err := s.client.PresignedGetObject(ctx, bucketName, name, expires, nil)
	if err != nil {
		return nil, err
	}

	return &storage.SignResult{
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

func (s *MinioDriver) Copy(ctx context.Context, srcBucket string, srcName string, dstBucket string, dstName string) error {
	dest := minio.CopyDestOptions{Bucket: dstBucket, Object: dstName}
	src := minio.CopySrcOptions{Bucket: srcBucket, Object: srcName}
	_, err := s.client.CopyObject(ctx, dest, src)
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
