package database

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"net"
	"sync"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// TODO: we have some code duplication here, we can do better!

type DriverParams struct {
	CloudSQL CloudSQLDriverParams
}

func registerDriver(name string, params *DriverParams) (CleanupFn, error) {
	if params != nil && params.CloudSQL.Enabled {
		return registerCloudSQLDriverWithParams(name, params.CloudSQL)
	}

	return registerDefaultDriver(name)
}

// MARK: - Default

type CleanupFn func() error

var dummyCleanup = func() error { return nil }

func registerDefaultDriver(name string) (CleanupFn, error) {
	sql.Register(name, &defaultDriver{
		dbURIs: make(map[string]string),
	})
	return dummyCleanup, nil
}

type defaultDriver struct {
	mu     sync.RWMutex
	dbURIs map[string]string
}

func (d *defaultDriver) Open(name string) (driver.Conn, error) {
	dbURI, err := d.dbURI(name)
	if err != nil {
		return nil, err
	}

	return stdlib.GetDefaultDriver().Open(dbURI)
}

// dbURI registers a driver using the provided DSN. If the name has already
// been registered, dbURI returns the existing registration.
func (d *defaultDriver) dbURI(name string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	dbURI, ok := d.dbURIs[name]
	if ok {
		return dbURI, nil
	}

	config, err := createConfig(name)
	if err != nil {
		return "", err
	}

	dbURI = stdlib.RegisterConnConfig(config)
	d.dbURIs[name] = dbURI

	return dbURI, nil
}

// MARK: - CloudSQL

type CloudSQLDriverParams struct {
	Enabled   bool
	IAM       bool
	PrivateIP bool
}

func registerCloudSQLDriverWithParams(name string, params CloudSQLDriverParams) (CleanupFn, error) {
	dialOpts := []cloudsqlconn.DialOption{}

	if params.PrivateIP {
		dialOpts = append(dialOpts, cloudsqlconn.WithPrivateIP())
	}

	dialerOpts := []cloudsqlconn.Option{}

	if len(dialOpts) > 0 {
		dialerOpts = append(dialerOpts, cloudsqlconn.WithDefaultDialOptions(dialOpts...))
	}

	if params.IAM {
		dialerOpts = append(dialerOpts, cloudsqlconn.WithIAMAuthN())
	}

	return registerCloudSQLDriver(name, dialerOpts...)
}

// registerCloudSQLDriver registers a Postgres driver that uses the cloudsqlconn.Dialer
// configured with the provided options. The choice of name is entirely up to
// the caller and may be used to distinguish between multiple registrations of
// differently configured Dialers. The driver uses pgx/v5 internally.
// RegisterDriver returns a cleanup function that should be called one the
// database connection is no longer needed.
func registerCloudSQLDriver(name string, opts ...cloudsqlconn.Option) (func() error, error) {
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return dummyCleanup, err
	}
	sql.Register(name, &cloudSqlPgDriver{
		d:      d,
		dbURIs: make(map[string]string),
	})
	return func() error { return d.Close() }, nil
}

type cloudSqlPgDriver struct {
	d  *cloudsqlconn.Dialer
	mu sync.RWMutex
	// dbURIs is a map of DSN to DB URI for registered connection names.
	dbURIs map[string]string
}

// Open accepts a keyword/value formatted connection string and returns a
// connection to the database using cloudsqlconn.Dialer. The Cloud SQL instance
// connection name should be specified in the host field. For example:
//
// "host=my-project:us-central1:my-db-instance user=myuser password=mypass"
func (p *cloudSqlPgDriver) Open(name string) (driver.Conn, error) {
	dbURI, err := p.dbURI(name)
	if err != nil {
		return nil, err
	}
	return stdlib.GetDefaultDriver().Open(dbURI)
}

// dbURI registers a driver using the provided DSN. If the name has already
// been registered, dbURI returns the existing registration.
func (p *cloudSqlPgDriver) dbURI(name string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	dbURI, ok := p.dbURIs[name]
	if ok {
		return dbURI, nil
	}

	config, err := createConfig(name)
	if err != nil {
		return "", err
	}
	instConnName := config.Config.Host // Extract instance connection name
	config.Config.Host = "localhost"   // Replace it with a default value
	config.DialFunc = func(ctx context.Context, _, _ string) (net.Conn, error) {
		return p.d.Dial(ctx, instConnName)
	}

	dbURI = stdlib.RegisterConnConfig(config)
	p.dbURIs[name] = dbURI

	return dbURI, nil
}

// MARK: - helpers

func createConfig(name string) (*pgx.ConnConfig, error) {
	config, err := pgx.ParseConfig(name)
	if err != nil {
		return nil, err
	}

	// TODO: otel configuration
	config.Tracer = otelpgx.NewTracer()

	return config, nil
}
