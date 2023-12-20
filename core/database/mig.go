package database

import (
	"context"
	"fmt"

	"go.uber.org/fx"

	atlas "ariga.io/atlas/sql/migrate"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql/schema"

	_ "ariga.io/atlas/sql/postgres"
	_ "github.com/lib/pq"
)

type MigratorParams struct {
	fx.In

	Config *Config
}

type Mig struct {
	config *Config
}

func NewMig(params MigratorParams) (*Mig, error) {
	if params.Config == nil {
		return nil, fmt.Errorf("no db config provided")
	}

	return &Mig{
		config: params.Config,
	}, nil
}

type Differ = func(ctx context.Context, url, name string, opts ...schema.MigrateOption) error

// TODO: move into migrator struct
func (m *Mig) Diff(ctx context.Context, name string, differ Differ) error {
	dir, err := atlas.NewLocalDir(m.config.MigrationPath)
	if err != nil {
		return fmt.Errorf("failed creating atlas migration directory: %v", err)
	}

	opts := []schema.MigrateOption{
		schema.WithDir(dir),
		schema.WithMigrationMode(schema.ModeReplay),
		schema.WithDialect(dialect.Postgres),
		schema.WithFormatter(atlas.DefaultFormatter),
	}

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?search_path=%s&sslmode=disable",
		m.config.Username, m.config.Password, m.config.Host, m.config.Port, m.config.MigrationDevDatabase, m.config.Schema,
	)

	err = differ(ctx, dsn, name, opts...)
	if err != nil {
		return fmt.Errorf("failed generating migration: %v", err)
	}

	return nil
}
