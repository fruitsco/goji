package database

import (
	"context"
	"database/sql"
	"sync/atomic"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
)

type multiDriver struct {
	r, w dialect.Driver

	// queryWriteFraction is the fraction of total queries that should be sent
	// to the write replica. 0 < queryWriteFraction < queryTotal.
	//
	// Example: If queryWriteFraction = 1, queryTotal = 10. 1/10 queries will
	// be sent to the write replica, 9/10 to the read replica.
	queryWriteFraction int64

	// queryTotal is the total number of queries that are sent to both the read
	// and write replicas.
	queryTotal int64

	// queryCounter is the number of queries sent to either replica.
	queryCounter int64
}

var _ dialect.Driver = (*multiDriver)(nil)

func (d *multiDriver) Query(ctx context.Context, query string, args, v any) error {
	e := d.r

	n := atomic.AddInt64(&d.queryCounter, 1)

	if ent.QueryFromContext(ctx) == nil {
		// Mutation statements that use the RETURNING clause.
		e = d.w
	} else if (n % d.queryTotal) < d.queryWriteFraction {
		// Round-robin between read and write based on writeFraction
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
