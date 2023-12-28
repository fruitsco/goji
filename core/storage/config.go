package storage

import "github.com/fruitsco/goji/x/conf"

type StorageDriver string

const (
	GCS   StorageDriver = "gcs"
	Minio StorageDriver = "minio"
	NoOp  StorageDriver = "noop"
)

type GCSConfig struct {
	Region     string `conf:"region"`
	Expiration int    `conf:"signed_url_expiration"`
}

type MinioConfig struct {
	// The host:port of the minio server to connect to
	Host string `conf:"host"`

	// The access key to use when connecting to the minio server
	AccessKey string `conf:"access_key"`

	// The secret key to use when connecting to the minio server
	SecretKey string `conf:"secret_key"`

	// Whether to use HTTPS when connecting to the minio server
	Secure bool `conf:"secure"`

	// The URL of the proxy to use when connecting to the minio server
	ProxyUrl string `conf:"proxy_url"`

	// The region of the minio server
	Region string `conf:"region"`

	// The expiration time of signed URLs
	Expiration int `conf:"signed_url_expiration"`

	// Whether to skip TLS verification
	TLSSkipVerify bool `conf:"tls_skip_verify"`

	// Whether to enable HTTP request tracing
	Trace bool `conf:"http_trace"`
}

type Config struct {
	Driver StorageDriver `conf:"driver"`
	GCS    *GCSConfig    `conf:"gcs"`
	Minio  *MinioConfig  `conf:"minio"`
}

var DefaultConfig = conf.DefaultConfig{
	"storage.driver": "noop",

	// minio
	"storage.minio.host":                  "localhost:9000",
	"storage.minio.https":                 "false",
	"storage.minio.secure":                "false",
	"storage.minio.signed_url_expiration": "3600",

	// gcs
	"storage.gcs.signed_url_expiration": "3600",
}
