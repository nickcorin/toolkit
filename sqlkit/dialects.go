package sqlkit

import (
	"fmt"
	"net/url"
	"os"
)

// Dialect represents a type of SQL database.
//
// E.g. Postgres, MySQL, SQLite.
type Dialect string

const (
	Postgres Dialect = "postgres"
)

// ErrUnsupportedDialect is returned when the dialect is not supported.
var ErrUnsupportedDialect = fmt.Errorf("unsupported dialect")

// GetConnector returns a Connector for the given dialect.
func GetConnector(d Dialect) (Connector, error) {
	switch d {
	case "postgres":
		return postgres{}, nil
	default:
		return nil, ErrUnsupportedDialect
	}
}

// Connector is an interface that provides methods for configuring a database connection.
type Connector interface {
	Defaults() Config
	Driver() string
	DSN(cfg Config) string
}

type postgres struct{}

func (p postgres) Defaults() Config {
	return Config{
		Host:     "localhost",
		User:     os.Getenv("USER"),
		Port:     5432,
		Database: "postgres",
		Flags:    make(url.Values),
	}
}

func (p postgres) Driver() string {
	return "pgx"
}

func (p postgres) DSN(c Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", c.User, c.Password, c.Host, c.Port, c.Database, c.Flags.Encode())
}
