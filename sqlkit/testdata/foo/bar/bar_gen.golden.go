// Code generated by scangen. DO NOT EDIT.
package bar

import (
	"context"
	sql "database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/nickcorin/toolkit/sqlkit"
	foo "github.com/nickcorin/toolkit/sqlkit/testdata/foo"
)

// ErrFooNotFound is returned when a query for a foo.Foo returns no results.
var ErrFooNotFound = errors.New("foo not found")

type PostgresRespository struct {
	conn      *sql.DB
	tableName string
	cols      []string
}

func NewPostgresRespository(conn *sql.DB) *PostgresRespositoryRepository {
	return &PostgresRespositoryRepository{
		conn:      conn,
		tableName: "foos",
		cols:      []string{"a", "b", "c", "d_override", "e", "f"},
	}
}

func (r *PostgresRespository) selectPrefix() string {
	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(r.cols, ", "), r.tableName)
}

func (r *PostgresRespository) lookupWhere(ctx context.Context, where string, args ...any) (*foo.Foo, error) {
	row := r.conn.QueryRowContext(ctx, fmt.Sprintf(r.selectPrefix()+" WHERE %s", where), args...)
	return r.scan(row)
}

func (r *PostgresRespository) listWhere(ctx context.Context, where string, args ...any) ([]*foo.Foo, error) {
	rows, err := r.conn.QueryContext(ctx, fmt.Sprintf(r.selectPrefix()+" WHERE %s", where), args...)
	if err != nil {
		return nil, fmt.Errorf("list foo: %w", err)
	}
	return r.list(rows)
}

func (r *PostgresRespository) list(rows *sql.Rows) ([]*foo.Foo, error) {
	ret := make([]*foo.Foo, 0)
	for rows.Next() {
		item, err := r.scan(rows)
		if err != nil {
			return nil, err
		}

		ret = append(ret, item)
	}

	return ret, nil
}

func (r *PostgresRespository) scan(row sqlkit.Scannable) (*foo.Foo, error) {
	var scan bar

	err := row.Scan(&scan.A, &scan.B, &scan.C, &scan.D, &scan.E, &scan.F)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFooNotFound
		}

		return nil, fmt.Errorf("scan foo: %w", err)
	}

	var ret foo.Foo

	ret.A = scan.A
	ret.B = scan.B
	ret.C = foo.Baz(scan.C)
	ret.D = scan.D.Time
	ret.E = scan.E.V
	ret.F = scan.F

	return &ret, nil
}
