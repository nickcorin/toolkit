package sqlkit

import (
	"fmt"
	"os"
)

//go:generate stringer -type=Dialect -linecomment

// Dialect represents a type of SQL database.
//
// E.g. Postgres, MySQL, SQLite.
type Dialect int

const (
	UnknownDialect Dialect = iota // unknown
	Postgres                      // postgres
	MySQL                         // mysql
	sentinel
)

// GetDialectFromString returns a Dialect from a string. It may not necessarily return a valid dialect.
func GetDialectFromString(s string) Dialect {
	var d Dialect
	for i := UnknownDialect; i < sentinel; i++ {
		if i.String() == s {
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
	return d > UnknownDialect && d < sentinel
}

// ErrUnsupportedDialect is returned when the dialect is not supported.
var ErrUnsupportedDialect = fmt.Errorf("unsupported dialect")

// GetConnector returns a Connector for the given dialect.
func GetConnector(d Dialect) (Connector, error) {
	switch d {
	case Postgres:
		return postgres{}, nil
	default:
		return nil, ErrUnsupportedDialect
	}
}

// Connector is an interface that provides methods for configuring a database connection.
type Connector interface {
	Defaults() *Config
	Driver() string
	DSN(cfg *Config) string
}

type postgres struct{}

func (p postgres) Defaults() *Config {
	return &Config{
		Dialect:  Postgres,
		Host:     "localhost",
		User:     os.Getenv("USER"),
		Port:     5432,
		Database: "postgres",
		Flags:    make(Flags),
	}
}

func (p postgres) Driver() string {
	return "pgx"
}

func (p postgres) DSN(c *Config) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?%s", c.User, c.Password, c.Host, c.Port, c.Database, c.Flags.Encode())
}
