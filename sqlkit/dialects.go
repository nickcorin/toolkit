package sqlkit

import (
	"fmt"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

//go:generate stringer -type=Dialect -linecomment

// Dialect represents a type of SQL database.
//
// E.g. Postgres, MySQL, SQLite.
type Dialect int

const (
	Unspecified Dialect = iota // unspecified
	Postgres                   // postgres
	sentinel
)

// GetDialectFromString returns a Dialect from a string. It may not necessarily return a valid dialect.
func GetDialectFromString(s string) Dialect {
	var d Dialect
	for i := Unspecified; i < sentinel; i++ {
		if strings.EqualFold(i.String(), s) {
			d = i
			break
		}
	}

	return d
}

// Decode implements the envconfig.Decoder interface.
func (d *Dialect) Decode(value string) error {
	*d = GetDialectFromString(value)
	return nil
}

// Valid returns true if the dialect is a valid dialect.
func (d Dialect) Valid() bool {
	return d > Unspecified && d < sentinel
}

// ErrUnsupportedDialect is returned when the dialect is not supported.
var ErrUnsupportedDialect = fmt.Errorf("unsupported dialect")

// GetConnector returns a Connector for the given dialect.
func GetConnector(d Dialect) (Connector, error) {
	switch d {
	case Postgres:
		return new(PostgresConnector), nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnsupportedDialect, d)
	}
}

// Connector is an interface that provides methods for configuring a database connection.
type Connector interface {
	Defaults() (*Config, error)
	Driver() string
	DSN(cfg *Config) string
}

type PostgresConnector struct{}

func (p *PostgresConnector) Defaults() (*Config, error) {
	var c Config
	if err := envconfig.Process("", &c); err != nil {
		return nil, fmt.Errorf("failed to process env variables: %w", err)
	}

	c.OverrideWith(&Config{
		Dialect:  Postgres,
		Host:     "localhost",
		User:     os.Getenv("USER"),
		Port:     5432,
		Database: "postgres",
		Flags:    make(Flags),
	})

	return &c, nil
}

func (p *PostgresConnector) Driver() string {
	return "pgx"
}

func (p *PostgresConnector) DSN(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", c.User, c.Password, c.Host, c.Port, c.Database, c.Flags.Encode())
}
