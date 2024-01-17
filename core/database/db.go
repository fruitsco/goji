package database

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect"
	"go.uber.org/fx"
)

type EntDB struct {
	connection Connection
}

type DBParams struct {
	fx.In

	Config *Config
}

func NewLifecycleDB(lc fx.Lifecycle, params DBParams) (*EntDB, error) {
	db, err := NewDB(params)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return db.Close()
		},
	})

	return db, nil
}

func NewDB(params DBParams) (*EntDB, error) {
	if params.Config == nil {
		return nil, fmt.Errorf("no db config provided")
	}

	cloudSql := CloudSQLDriverParams{
		Enabled:   params.Config.CloudSQL,
		IAM:       params.Config.CloudSQLIAM,
		PrivateIP: params.Config.CloudSQLPrivateIP,
	}

	mainConnection, err := NewConnection(ConnectionParams{
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
		MaxIdleConns:  params.Config.MaxIdleConnections,
		MaxOpenConns:  params.Config.MaxOpenConnections,
		CloudSQL:      cloudSql,
	})
	if err != nil {
		return nil, err
	}

	var connection Connection = mainConnection

	if params.Config.Replica {
		replicaConnection, err := NewConnection(ConnectionParams{
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
			MaxIdleConns:  params.Config.MaxIdleConnections,
			MaxOpenConns:  params.Config.MaxOpenConnections,
			CloudSQL:      cloudSql,
		})
		if err != nil {
			return nil, err
		}

		connection = &multiConnection{
			w: mainConnection,
			r: replicaConnection,
		}
	}

	return &EntDB{
		connection: connection,
	}, nil
}

func (db *EntDB) Driver() dialect.Driver {
	return db.connection.Driver()
}

func (db *EntDB) Close() error {
	return db.connection.Close()
}
