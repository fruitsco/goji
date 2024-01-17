package database

import (
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

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

	MaxIdleConns int
	MaxOpenConns int

	CloudSQL CloudSQLDriverParams
}

type Connection interface {
	Driver() dialect.Driver
	Close() error
}

// MARK: - Single Connection

type singleConnection struct {
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
		driver:  driver,
		cleanup: cleanup,
	}, nil
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

// MARK: - Multi Connection

type multiConnection struct {
	r, w Connection
}

var _ Connection = (*multiConnection)(nil)

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

// createDB creates a DB connection using the pgx driver.
//
// In default mode, no connector is used. SSL certificates must be
// provided manually. IAM authentication is not supported.
//
// If CloudSQL is enabled, the connection is made using an underlying
// Cloud SQL Connector. SSL certificates are not required as they are
// managed by the connector. IAM authentication is supported.
// See https://github.com/GoogleCloudPlatform/cloud-sql-go-connector
func createDB(params ConnectionParams) (*sql.DB, CleanupFn, error) {
	dbURI := buildDSN(params)

	driverParams := &DriverParams{
		CloudSQL: params.CloudSQL,
	}

	cleanup, err := registerDriver("goji", driverParams)
	if err != nil {
		return nil, nil, err
	}

	db, err := sql.Open("goji", dbURI)
	if err != nil {
		return nil, nil, err
	}

	db.SetMaxIdleConns(params.MaxIdleConns)
	db.SetMaxOpenConns(params.MaxOpenConns)

	return db, cleanup, nil
}

func buildDSN(params ConnectionParams) string {
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
