package sqlkit

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
)

// Config represents the configuration for a database connection.
type Config struct {
	Dialect  string
	Host     string `default:"localhost"`
	Port     int    `default:"5432"`
	User     string
	Password string
	Database string
	Flags    url.Values
}

func Connect(ctx context.Context, config Config) (*sql.DB, error) {
	connector, err := getConnector(Dialect(config.Dialect))
	if err != nil {
		return nil, fmt.Errorf("get connector: %w", err)
	}

	conn, err := connectWithDSN(ctx, connector.Driver(), connector.DSN(config))
	if err != nil {
		return nil, fmt.Errorf("get DSN: %w", err)
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
