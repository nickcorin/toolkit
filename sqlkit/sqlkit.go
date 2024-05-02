package sqlkit

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5/pgconn"
)

// Scannable is an interface that wraps the Scan method.
type Scannable interface {
	Scan(dest ...interface{}) error
}

var (
	_ Scannable = (*sql.Row)(nil)
	_ Scannable = (*sql.Rows)(nil)
)

// Config represents the configuration for a database connection.
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
	Flags    url.Values
}

func Connect(ctx context.Context, dialect Dialect, config Config) (*sql.DB, error) {
	connector, err := GetConnector(dialect)
	if err != nil {
		return nil, fmt.Errorf("get connector: %w", err)
	}

	conn, err := connectWithDSN(ctx, connector.Driver(), connector.DSN(config))
	if err != nil {
		return nil, fmt.Errorf("connect with dsn: %w", err)
	}

	return conn, nil
}

func connectWithDSN(ctx context.Context, driver, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("open connection: %w", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping connection: %w", err)
	}

	return db, nil
}

func PgErrorIs(err error, target string) bool {
	if err == nil {
		return false
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == target
	}

	return false
}
