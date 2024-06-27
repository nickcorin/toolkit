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
	Host     string `envconfig:"HOST"`
	Port     int    `envconfig:"PORT"`
	User     string `envconfig:"USER"`
	Password string `envconfig:"PASSWORD"`
	Database string `envconfig:"DATABASE"`
	Flags    Flags  `envconfig:"FLAGS"`
}

// Flags is an alias for url.Values.
type Flags url.Values

// Set implements the flag.Value interface.
func (f Flags) Set(value string) error {
	query, err := url.ParseQuery(value)
	if err != nil {
		return err
	}
	f = Flags(query)

	return nil
}

func (c *Config) OverrideWith(custom *Config) {
	if c.Host == "" {
		c.Host = custom.Host
	}
	if c.Port == 0 {
		c.Port = custom.Port
	}
	if c.User == "" {
		c.User = custom.User
	}
	if c.Password == "" {
		c.Password = custom.Password
	}
	if c.Database == "" {
		c.Database = custom.Database
	}
	if c.Flags == nil {
		c.Flags = custom.Flags
	} else {
		for k, v := range custom.Flags {
			if _, ok := c.Flags[k]; !ok {
				c.Flags[k] = v
			}
		}
	}
}

func Connect(ctx context.Context, dialect Dialect, config *Config) (*sql.DB, error) {
	connector, err := GetConnector(dialect)
	if err != nil {
		return nil, fmt.Errorf("get connector: %w", err)
	}

	config.OverrideWith(connector.Defaults())

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
