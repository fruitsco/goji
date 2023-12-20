package database

import (
	"database/sql"
	"fmt"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/fx"
)

type EntDB struct {
	driver dialect.Driver
	config *Config
}

type DBParams struct {
	fx.In

	Config *Config
}

func NewDB(params DBParams) (*EntDB, error) {
	if params.Config == nil {
		return nil, fmt.Errorf("no db config provided")
	}

	mainInstanceDriver, err := createInstanceDriver(DBInstanceParams{
		Username:      params.Config.Username,
		Password:      params.Config.Password,
		Name:          params.Config.Name,
		Host:          params.Config.Host,
		Port:          params.Config.Port,
		Schema:        params.Config.Schema,
		Ssl:           params.Config.Ssl,
		SslRootCert:   params.Config.SslRootCert,
		SslClientCert: params.Config.SslClientCert,
		SslClientKey:  params.Config.SslClientKey,
	})
	if err != nil {
		return nil, err
	}

	driver := mainInstanceDriver

	if params.Config.Replica {
		replicaInstanceDriver, err := createInstanceDriver(DBInstanceParams{
			Username:      params.Config.Username,
			Password:      params.Config.Password,
			Name:          params.Config.Name,
			Host:          params.Config.ReplicaHost,
			Port:          params.Config.ReplicaPort,
			Schema:        params.Config.Schema,
			Ssl:           params.Config.ReplicaSsl,
			SslRootCert:   params.Config.ReplicaSslRootCert,
			SslClientCert: params.Config.ReplicaSslClientCert,
			SslClientKey:  params.Config.ReplicaSslClientKey,
		})
		if err != nil {
			return nil, err
		}

		driver = &multiDriver{
			r: replicaInstanceDriver,
			w: mainInstanceDriver,
		}
	}

	return &EntDB{
		driver: driver,
		config: params.Config,
	}, nil
}

func (db *EntDB) Driver() dialect.Driver {
	return db.driver
}

type DBInstanceParams struct {
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
}

func createInstanceDriver(params DBInstanceParams) (dialect.Driver, error) {
	dsn := buildDSN(params)

	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	return entsql.OpenDB(dialect.Postgres, db), nil
}

func buildDSN(params DBInstanceParams) string {
	dbURI := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d search_path=%s",
		params.Username, params.Password, params.Name, params.Host, params.Port, params.Schema,
	)

	if params.Ssl {
		dbURI += fmt.Sprintf(
			" sslmode=require sslrootcert=%s sslcert=%s sslkey=%s",
			params.SslRootCert, params.SslClientCert, params.SslClientKey,
		)
	}

	return dbURI
}
