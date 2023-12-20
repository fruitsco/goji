package database

import (
	"context"
	"database/sql"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
)

type multiDriver struct {
	r, w dialect.Driver
}

var _ dialect.Driver = (*multiDriver)(nil)

func (d *multiDriver) Query(ctx context.Context, query string, args, v any) error {
	e := d.r
	// Mutation statements that use the RETURNING clause.
	if ent.QueryFromContext(ctx) == nil {
		e = d.w
	}

	return e.Query(ctx, query, args, v)
}

func (d *multiDriver) Exec(ctx context.Context, query string, args, v any) error {
	return d.w.Exec(ctx, query, args, v)
}

func (d *multiDriver) Tx(ctx context.Context) (dialect.Tx, error) {
	return d.w.Tx(ctx)
}

func (d *multiDriver) BeginTx(ctx context.Context, opts *sql.TxOptions) (dialect.Tx, error) {
	return d.w.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
}

func (d *multiDriver) Close() error {
	rerr := d.r.Close()
	werr := d.w.Close()

	if rerr != nil {
		return rerr
	}

	if werr != nil {
		return werr
	}

	return nil
}

func (d *multiDriver) Dialect() string {
	return d.r.Dialect()
}
