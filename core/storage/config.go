package storage

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
	Host          string  `conf:"host"`
	AccessKey     string  `conf:"access_key"`
	SecretKey     string  `conf:"secret_key"`
	Secure        bool    `conf:"secure"`
	ProxyUrl      *string `conf:"proxy_url"`
	Region        string  `conf:"region"`
	Expiration    int     `conf:"signed_url_expiration"`
	TLSSkipVerify bool    `conf:"tls_skip_verify"`
	Trace         bool    `conf:"http_trace"`
}

type Config struct {
	Driver StorageDriver `conf:"driver"`
	GCS    *GCSConfig    `conf:"gcs"`
	Minio  *MinioConfig  `conf:"minio"`
}

var DefaultConfig = map[string]any{
	"storage.driver": "noop",

	// minio
	"storage.minio.host":                  "localhost:9000",
	"storage.minio.secure":                "false",
	"storage.minio.signed_url_expiration": "3600",

	// gcs
	"storage.gcs.signed_url_expiration": "3600",
}
