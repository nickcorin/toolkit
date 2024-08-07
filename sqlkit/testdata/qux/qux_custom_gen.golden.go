// Code generated by scangen; DO NOT EDIT.
package qux

import (
	context "context"
	sql "database/sql"
	errors "errors"
	fmt "fmt"
	strings "strings"

	sqlkit "github.com/nickcorin/toolkit/sqlkit"
)

// ErrQuxNotFound is returned when a query for a Qux returns no results.
var ErrQuxNotFound = errors.New("qux not found")

type QuxRepository struct {
	conn      *sql.DB
	tableName string
	cols      []string
}

func NewQuxRepository(conn *sql.DB) *QuxRepository {
	return &QuxRepository{
		conn:      conn,
		tableName: "quxes",
		cols:      []string{"a", "b", "c"},
	}
}

func (r *QuxRepository) selectPrefix() string {
	return fmt.Sprintf("SELECT %s FROM %s", strings.Join(r.cols, ", "), r.tableName)
}

func (r *QuxRepository) selectDistinctPrefix() string {
	return fmt.Sprintf("SELECT DISTINCT %s FROM %s", strings.Join(r.cols, ", "), r.tableName)
}

func (r *QuxRepository) lookupWhere(ctx context.Context, where string, args ...any) (*Qux, error) {
	row := r.conn.QueryRowContext(ctx, fmt.Sprintf(r.selectPrefix()+" WHERE %s", where), args...)
	return r.scan(row)
}

func (r *QuxRepository) listWhere(ctx context.Context, where string, args ...any) ([]*Qux, error) {
	rows, err := r.conn.QueryContext(ctx, fmt.Sprintf(r.selectPrefix()+" WHERE %s", where), args...)
	if err != nil {
		return nil, fmt.Errorf("list qux: %w", err)
	}
	return r.list(rows)
}

func (r *QuxRepository) listDistinctWhere(ctx context.Context, where string, args ...any) ([]*Qux, error) {
	rows, err := r.conn.QueryContext(ctx, fmt.Sprintf(r.selectDistinctPrefix()+" WHERE %s", where), args...)
	if err != nil {
		return nil, fmt.Errorf("list qux: %w", err)
	}
	return r.list(rows)
}

func (r *QuxRepository) list(rows *sql.Rows) ([]*Qux, error) {
	defer rows.Close()
	ret := make([]*Qux, 0)
	for rows.Next() {
		item, err := r.scan(rows)
		if err != nil {
			return nil, err
		}

		ret = append(ret, item)
	}

	return ret, nil
}

func (r *QuxRepository) scan(row sqlkit.Scannable) (*Qux, error) {
	var scan gen

	err := row.Scan(&scan.A, &scan.B, &scan.C)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrQuxNotFound
		}

		return nil, fmt.Errorf("scan qux: %w", err)
	}

	var ret Qux

	ret.A = scan.A
	ret.B = scan.B
	ret.C = scan.C

	return &ret, nil
}
