package database

import (
	"database/sql"
	"fmt"

	"cloud.google.com/go/cloudsqlconn"
	"cloud.google.com/go/cloudsqlconn/postgres/pgxv5"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type CleanupFn func() error

type ConnectionParams struct {
	Username      string
	Password      string
	Name          string
	Host          string
	Port          int
	Schema        string
	Ssl           bool
	SslRootCert   string
	SslClientCert string
	SslClientKey  string

	CloudSQL CloudSQLConnectorParams
}

type CloudSQLConnectorParams struct {
	Enabled   bool
	IAM       bool
	PrivateIP bool
}

type Connection interface {
	Driver() dialect.Driver
	Close() error
	DB() *sql.DB
}

// MARK: - Single Connection

type singleConnection struct {
	db      *sql.DB
	driver  *entsql.Driver
	cleanup CleanupFn
}

var _ Connection = (*singleConnection)(nil)

func NewConnection(params ConnectionParams) (*singleConnection, error) {
	db, cleanup, err := createDB(params)
	if err != nil {
		return nil, err
	}

	driver := entsql.OpenDB(dialect.Postgres, db)

	return &singleConnection{
		db:      db,
		driver:  driver,
		cleanup: cleanup,
	}, nil
}

func (c *singleConnection) DB() *sql.DB {
	return c.db
}

func (c *singleConnection) Driver() dialect.Driver {
	return c.driver
}

func (c *singleConnection) Close() error {
	cErr := c.cleanup()
	dErr := c.driver.Close()

	if cErr != nil {
		return cErr
	}

	if dErr != nil {
		return dErr
	}

	return nil
}

func createDB(params ConnectionParams) (*sql.DB, CleanupFn, error) {
	if params.CloudSQL.Enabled {
		return createCloudSQLConnectorDB(params)
	}

	return createBasicDB(params)
}

var dummyCleanup = func() error { return nil }

// createBasicDB creates a basic DB connection using the pgx driver.
// This method does not use any connector. SSL certificates must be
// provided manually. IAM authentication is not supported.
func createBasicDB(params ConnectionParams) (*sql.DB, CleanupFn, error) {
	dbURI := DsnForConnection(params)

	db, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, nil, err
	}

	return db, dummyCleanup, nil
}

// MARK: - Multi Connection

type multiConnection struct {
	r, w Connection
}

var _ Connection = (*multiConnection)(nil)

func (c *multiConnection) DB() *sql.DB {
	// return the write connection for multi connection
	return c.w.DB()
}

func (c *multiConnection) Driver() dialect.Driver {
	return &multiDriver{
		r: c.r.Driver(),
		w: c.w.Driver(),
	}
}

func (c *multiConnection) Close() error {
	rerr := c.r.Close()
	werr := c.w.Close()

	if rerr != nil {
		return rerr
	}

	if werr != nil {
		return werr
	}

	return nil
}

// MARK: - Create Connections

// createCloudSQLConnectorDB creates a DB connection using the pgx driver
// and the underlying Cloud SQL Connector. SSL certificates are not required
// as they are managed by the connector. IAM authentication is supported.
// See https://github.com/GoogleCloudPlatform/cloud-sql-go-connector
func createCloudSQLConnectorDB(params ConnectionParams) (*sql.DB, CleanupFn, error) {
	dialOpts := []cloudsqlconn.DialOption{}

	if params.CloudSQL.PrivateIP {
		dialOpts = append(dialOpts, cloudsqlconn.WithPrivateIP())
	}

	dialerOpts := []cloudsqlconn.Option{}

	if len(dialOpts) > 0 {
		dialerOpts = append(dialerOpts, cloudsqlconn.WithDefaultDialOptions(dialOpts...))
	}

	if params.CloudSQL.IAM {
		dialerOpts = append(dialerOpts, cloudsqlconn.WithIAMAuthN())
	}

	close, err := pgxv5.RegisterDriver("cloudsql-postgres", dialerOpts...)
	if err != nil {
		return nil, nil, err
	}

	dbURI := DsnForConnection(params)

	db, err := sql.Open("cloudsql-postgres", dbURI)
	if err != nil {
		close()
		return nil, nil, err
	}

	return db, close, nil
}

func DsnForConnection(params ConnectionParams) string {
	dbURI := fmt.Sprintf(
		"user=%s dbname=%s host=%s port=%d search_path=%s",
		params.Username, params.Name, params.Host, params.Port, params.Schema,
	)

	// Do not include password if IAM authentication is enabled
	if !(params.CloudSQL.Enabled && params.CloudSQL.IAM) {
		dbURI += fmt.Sprintf(" password=%s", params.Password)
	}

	if params.Ssl {
		dbURI += fmt.Sprintf(
			" sslmode=require sslrootcert=%s sslcert=%s sslkey=%s",
			params.SslRootCert, params.SslClientCert, params.SslClientKey,
		)
	} else {
		dbURI += " sslmode=disable"
	}

	return dbURI
}
