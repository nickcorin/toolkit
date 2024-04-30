package sqlkit

import (
	"fmt"
	"net/url"
	"os"
)

type Dialect string

const (
	Postgres Dialect = "postgres"
)

// ErrUnsupportedDialect is returned when the dialect is not supported.
var ErrUnsupportedDialect = fmt.Errorf("unsupported dialect")

// getConnector returns a dialect based on the given string.
func getConnector(d Dialect) (Connector, error) {
	switch d {
	case "postgres":
		return postgres{}, nil
	default:
		return nil, ErrUnsupportedDialect
	}
}

type Connector interface {
	Defaults() Config
	Driver() string
	DSN(cfg Config) string
}

type postgres struct{}

func (p postgres) Defaults() Config {
	return Config{
		Dialect:  "postgres",
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
