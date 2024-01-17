package database

import "github.com/fruitsco/goji/x/conf"

type Config struct {
	Host          string `conf:"host"`
	Port          int    `conf:"port"`
	Name          string `conf:"name"`
	Username      string `conf:"username"`
	Password      string `conf:"password"`
	Schema        string `conf:"schema"`
	Ssl           bool   `conf:"ssl"`
	SslRootCert   string `conf:"ssl_root_cert"`
	SslClientKey  string `conf:"ssl_client_key"`
	SslClientCert string `conf:"ssl_client_cert"`

	Replica              bool   `conf:"replica"`
	ReplicaHost          string `conf:"replica_host"`
	ReplicaPort          int    `conf:"replica_port"`
	ReplicaSsl           bool   `conf:"replica_ssl"`
	ReplicaSslRootCert   string `conf:"replica_ssl_root_cert"`
	ReplicaSslClientKey  string `conf:"replica_ssl_client_key"`
	ReplicaSslClientCert string `conf:"replica_ssl_client_cert"`

	MaxIdleConnections int `conf:"max_idle_connections"`
	MaxOpenConnections int `conf:"max_open_connections"`

	MigrationPath        string `conf:"migration_path"`
	MigrationDevDatabase string `conf:"migration_dev_database"`

	CloudSQL          bool `conf:"cloudsql_enabled"`
	CloudSQLIAM       bool `conf:"cloudsql_iam"`
	CloudSQLPrivateIP bool `conf:"cloudsql_private_ip"`
}

var DefaultConfig = conf.DefaultConfig{
	"db.host":                   "127.0.0.1",
	"db.port":                   "5432",
	"db.username":               "fruits",
	"db.password":               "fruits",
	"db.name":                   "fruits_roma",
	"db.schema":                 "public",
	"db.migration_dev_database": "fruits_roma_dev",
	"db.migration_path":         "internal/db/migrations",
	"db.replica_host":           "127.0.0.1",
	"db.replica_port":           "5432",
}
