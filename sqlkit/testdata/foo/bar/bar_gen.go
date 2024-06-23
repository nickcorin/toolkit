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

// ErrFooTypeNotFound is returned when a query for a foo.FooType returns no results.
var ErrFooTypeNotFound = errors.New("fooType not found")

type PostgresRepository struct {
	conn      *sql.DB
	tableName string
	cols      []string
}

func NewPostgresRepository(conn *sql.DB) *PostgresRepository {
	return &postgresRepository{
		conn:      conn,
		tableName: "foos",
		cols:      []string{"a", "b", "c", "d_override", "e", "f", "g"},
	}
}

func (r *PostgresRepository) selectPrefix() string {
	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(r.cols, ", "), r.tableName)
}

func (r *PostgresRepository) lookupWhere(ctx context.Context, where string, args ...any) (*foo.FooType, error) {
	row := r.conn.QueryRowContext(ctx, fmt.Sprintf(r.selectPrefix()+" WHERE %s", where), args...)
	return r.scan(row)
}

func (r *PostgresRepository) listWhere(ctx context.Context, where string, args ...any) ([]*foo.FooType, error) {
	rows := r.conn.QueryRowContext(ctx, fmt.Sprintf(r.selectPrefix()+" WHERE %s", where), args...)
	return r.list(rows)
}

func (r *PostgresRepository) list(rows sqlkit.Scannable) ([]*foo.FooType, error) {
	ret := make([]*foo.FooType, 0)
	for rows.Next() {
		item, err := r.scan(rows)
		if err != nil {
			return nil, err
		}

		ret = append(ret, item)
	}

	return ret, nil
}

func (r *PostgresRepository) scan(row sqlkit.Scannable) (*foo.FooType, error) {
	var scan barType

	err := row.Scan(&scan.A, &scan.B, &scan.C, &scan.D, &scan.E, &scan.F, &scan.G)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrFooTypeNotFound
		}

		return nil, fmt.Errorf("scan fooType: %w", err)
	}

	var ret foo.FooType

	ret.A = scan.A
	ret.B = scan.B
	ret.C = foo.BazType(scan.C)
	ret.D = scan.D
	ret.E = scan.E.V
	ret.F = scan.F
	ret.G = scan.G

	return &ret, nil
}